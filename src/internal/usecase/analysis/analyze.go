package analysis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/specvital/collector/internal/domain/analysis"
)

const (
	DefaultMaxConcurrentClones = 2
	DefaultAnalysisTimeout     = 15 * time.Minute
	// DefaultOAuthProvider is the OAuth provider for VCS authentication.
	// Currently only GitHub is supported as the VCS provider (see repoURL construction in Execute).
	DefaultOAuthProvider = "github"
)

// AnalyzeUseCase orchestrates repository analysis workflow.
type AnalyzeUseCase struct {
	cloneSem    *semaphore.Weighted
	parser      analysis.Parser
	repository  analysis.Repository
	timeout     time.Duration
	tokenLookup analysis.TokenLookup
	vcs         analysis.VCS
}

// Config holds configuration for AnalyzeUseCase.
type Config struct {
	AnalysisTimeout     time.Duration
	MaxConcurrentClones int64
}

// Option is a functional option for configuring AnalyzeUseCase.
type Option func(*Config)

// WithAnalysisTimeout sets the timeout for analysis operations.
// Zero or negative values are ignored and the default timeout is used.
func WithAnalysisTimeout(d time.Duration) Option {
	return func(cfg *Config) {
		if d > 0 {
			cfg.AnalysisTimeout = d
		}
	}
}

// WithMaxConcurrentClones sets the maximum number of concurrent clone operations.
// Zero or negative values are ignored and the default value is used.
func WithMaxConcurrentClones(n int64) Option {
	return func(cfg *Config) {
		if n > 0 {
			cfg.MaxConcurrentClones = n
		}
	}
}

// NewAnalyzeUseCase creates a new AnalyzeUseCase with given dependencies.
// tokenLookup is optional - if nil, all clones use public access (token=nil).
func NewAnalyzeUseCase(
	repository analysis.Repository,
	vcs analysis.VCS,
	parser analysis.Parser,
	tokenLookup analysis.TokenLookup,
	opts ...Option,
) *AnalyzeUseCase {
	cfg := Config{
		AnalysisTimeout:     DefaultAnalysisTimeout,
		MaxConcurrentClones: DefaultMaxConcurrentClones,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &AnalyzeUseCase{
		cloneSem:    semaphore.NewWeighted(cfg.MaxConcurrentClones),
		parser:      parser,
		repository:  repository,
		timeout:     cfg.AnalysisTimeout,
		tokenLookup: tokenLookup,
		vcs:         vcs,
	}
}

// Execute performs the complete analysis workflow:
// 1. Validates input
// 2. Clones repository (with concurrency control)
// 3. Creates analysis record
// 4. Scans for test inventory
// 5. Saves analysis results
// On any error after record creation, RecordFailure is called.
func (uc *AnalyzeUseCase) Execute(ctx context.Context, req analysis.AnalyzeRequest) (err error) {
	if err = req.Validate(); err != nil {
		return err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	repoURL := fmt.Sprintf("https://github.com/%s/%s", req.Owner, req.Repo)

	token, err := uc.lookupToken(timeoutCtx, req.UserID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrTokenLookupFailed, err)
	}

	src, err := uc.cloneWithSemaphore(timeoutCtx, repoURL, token)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCloneFailed, err)
	}
	defer uc.closeSource(src, req.Owner, req.Repo)

	createParams := analysis.CreateAnalysisRecordParams{
		Branch:    src.Branch(),
		CommitSHA: src.CommitSHA(),
		Owner:     req.Owner,
		Repo:      req.Repo,
	}
	if err = createParams.Validate(); err != nil {
		return fmt.Errorf("%w: %w", ErrSaveFailed, err)
	}

	analysisID, err := uc.repository.CreateAnalysisRecord(timeoutCtx, createParams)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSaveFailed, err)
	}

	// Record failure on any error after analysis record creation.
	// Uses context.Background() to ensure RecordFailure completes even if parent context
	// is cancelled (timeout, shutdown). Failure recording is critical for data integrity
	// and should complete independently of the analysis workflow lifecycle.
	defer func() {
		if err != nil {
			if recordErr := uc.repository.RecordFailure(context.Background(), analysisID, err.Error()); recordErr != nil {
				slog.ErrorContext(context.Background(), "failed to record analysis failure",
					"error", recordErr,
					"analysis_id", analysisID,
					"original_error", err,
				)
			}
		}
	}()

	inventory, err := uc.parser.Scan(timeoutCtx, src)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrScanFailed, err)
		return err
	}

	if inventory == nil {
		slog.WarnContext(ctx, "scan result has no inventory",
			"owner", req.Owner,
			"repo", req.Repo,
			"commit", src.CommitSHA(),
		)
		inventory = &analysis.Inventory{Files: []analysis.TestFile{}}
	}

	saveParams := analysis.SaveAnalysisInventoryParams{
		AnalysisID: analysisID,
		Inventory:  inventory,
	}
	if err = saveParams.Validate(); err != nil {
		err = fmt.Errorf("%w: %w", ErrSaveFailed, err)
		return err
	}

	if err = uc.repository.SaveAnalysisInventory(timeoutCtx, saveParams); err != nil {
		err = fmt.Errorf("%w: %w", ErrSaveFailed, err)
		return err
	}

	return nil
}

func (uc *AnalyzeUseCase) cloneWithSemaphore(ctx context.Context, url string, token *string) (analysis.Source, error) {
	if err := uc.cloneSem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer uc.cloneSem.Release(1)

	return uc.vcs.Clone(ctx, url, token)
}

// lookupToken retrieves OAuth token for the given user.
//
// Returns:
//   - (nil, nil): no userID provided, tokenLookup not configured, or token not found (graceful degradation)
//   - (*token, nil): token found successfully
//   - (nil, error): infrastructure error (should fail the operation)
//
// Token not found (analysis.ErrTokenNotFound) triggers graceful degradation and is logged at INFO level.
// Infrastructure errors are returned to fail the operation.
func (uc *AnalyzeUseCase) lookupToken(ctx context.Context, userID *string) (*string, error) {
	if userID == nil || uc.tokenLookup == nil {
		return nil, nil
	}

	token, err := uc.tokenLookup.GetOAuthToken(ctx, *userID, DefaultOAuthProvider)
	if err != nil {
		if errors.Is(err, analysis.ErrTokenNotFound) {
			slog.InfoContext(ctx, "no OAuth token found, using public access",
				"user_id", *userID,
			)
			return nil, nil
		}
		return nil, fmt.Errorf("failed to lookup OAuth token for user %s: %w", *userID, err)
	}

	if token == "" {
		slog.WarnContext(ctx, "empty token returned, using public access",
			"user_id", *userID,
		)
		return nil, nil
	}

	return &token, nil
}

func (uc *AnalyzeUseCase) closeSource(src analysis.Source, owner, repo string) {
	// Use background context for cleanup operations
	ctx := context.Background()
	if closeErr := src.Close(ctx); closeErr != nil {
		slog.Error("failed to close source",
			"error", closeErr,
			"owner", owner,
			"repo", repo,
		)
	}
}
