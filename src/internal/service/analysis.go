package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/semaphore"

	"github.com/specvital/collector/internal/repository"
	"github.com/specvital/core/pkg/parser"
	"github.com/specvital/core/pkg/source"
)

const (
	DefaultMaxConcurrentClones = 2
	DefaultAnalysisTimeout     = 15 * time.Minute
)

type AnalyzeRequest struct {
	Owner string
	Repo  string
}

func (r AnalyzeRequest) Validate() error {
	if r.Owner == "" {
		return fmt.Errorf("%w: owner is required", ErrInvalidInput)
	}
	if r.Repo == "" {
		return fmt.Errorf("%w: repo is required", ErrInvalidInput)
	}
	if len(r.Owner) > 39 || len(r.Repo) > 100 {
		return fmt.Errorf("%w: owner/repo exceeds length limit", ErrInvalidInput)
	}
	if !isValidGitHubName(r.Owner) || !isValidGitHubName(r.Repo) {
		return fmt.Errorf("%w: invalid characters in owner/repo", ErrInvalidInput)
	}
	return nil
}

func isValidGitHubName(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

type AnalysisService interface {
	Analyze(ctx context.Context, req AnalyzeRequest) error
}

type AnalysisServiceConfig struct {
	AnalysisTimeout     time.Duration
	MaxConcurrentClones int64
}

type analysisService struct {
	analysisRepo repository.AnalysisRepository
	cloneSem     *semaphore.Weighted
	timeout      time.Duration
}

func NewAnalysisService(repo repository.AnalysisRepository, opts ...AnalysisServiceOption) AnalysisService {
	cfg := AnalysisServiceConfig{
		AnalysisTimeout:     DefaultAnalysisTimeout,
		MaxConcurrentClones: DefaultMaxConcurrentClones,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &analysisService{
		analysisRepo: repo,
		cloneSem:     semaphore.NewWeighted(cfg.MaxConcurrentClones),
		timeout:      cfg.AnalysisTimeout,
	}
}

type AnalysisServiceOption func(*AnalysisServiceConfig)

func WithAnalysisTimeout(d time.Duration) AnalysisServiceOption {
	return func(cfg *AnalysisServiceConfig) {
		if d > 0 {
			cfg.AnalysisTimeout = d
		}
	}
}

func WithMaxConcurrentClones(n int64) AnalysisServiceOption {
	return func(cfg *AnalysisServiceConfig) {
		if n > 0 {
			cfg.MaxConcurrentClones = n
		}
	}
}

func (s *analysisService) Analyze(ctx context.Context, req AnalyzeRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// Apply total timeout for entire analysis operation
	timeoutCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	repoURL := fmt.Sprintf("https://github.com/%s/%s", req.Owner, req.Repo)

	gitSrc, err := func() (*source.GitSource, error) {
		if err := s.cloneSem.Acquire(timeoutCtx, 1); err != nil {
			return nil, err
		}
		defer s.cloneSem.Release(1)
		return source.NewGitSource(timeoutCtx, repoURL, nil)
	}()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCloneFailed, err)
	}
	defer gitSrc.Close()

	analysisID, err := s.analysisRepo.CreateAnalysisRecord(timeoutCtx, repository.CreateAnalysisRecordParams{
		Branch:    gitSrc.Branch(),
		CommitSHA: gitSrc.CommitSHA(),
		Owner:     req.Owner,
		Repo:      req.Repo,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSaveFailed, err)
	}

	result, err := parser.Scan(timeoutCtx, gitSrc)
	if err != nil {
		// Use parent ctx for RecordFailure to ensure DB write succeeds even if timeout fired
		if recordErr := s.analysisRepo.RecordFailure(ctx, analysisID, err.Error()); recordErr != nil {
			slog.ErrorContext(ctx, "failed to record scan failure",
				"error", recordErr,
				"analysis_id", analysisID,
				"original_error", err,
			)
		}
		return fmt.Errorf("%w: %w", ErrScanFailed, err)
	}

	if result.Inventory == nil {
		slog.WarnContext(ctx, "scan result has no inventory",
			"owner", req.Owner,
			"repo", req.Repo,
			"commit", gitSrc.CommitSHA(),
		)
	}

	if err := s.analysisRepo.SaveAnalysisInventory(timeoutCtx, repository.SaveAnalysisInventoryParams{
		AnalysisID: analysisID,
		Result:     result,
	}); err != nil {
		// Use parent ctx for RecordFailure to ensure DB write succeeds even if timeout fired
		if recordErr := s.analysisRepo.RecordFailure(ctx, analysisID, err.Error()); recordErr != nil {
			slog.ErrorContext(ctx, "failed to record save failure",
				"error", recordErr,
				"analysis_id", analysisID,
				"original_error", err,
			)
		}
		return fmt.Errorf("%w: %w", ErrSaveFailed, err)
	}

	return nil
}
