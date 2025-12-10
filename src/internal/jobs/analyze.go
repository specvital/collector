package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/collector/internal/repository"
	"github.com/specvital/core/pkg/parser"
	"github.com/specvital/core/pkg/source"

	_ "github.com/specvital/core/pkg/parser/strategies/all"
)

const TypeAnalyze = "analysis:analyze"

type AnalyzePayload struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

func (p AnalyzePayload) Validate() error {
	if p.Owner == "" {
		return errors.New("owner is required")
	}
	if p.Repo == "" {
		return errors.New("repo is required")
	}
	return nil
}

type AnalyzeHandler struct {
	analysisRepo repository.AnalysisRepository
}

func NewAnalyzeHandler(pool *pgxpool.Pool) *AnalyzeHandler {
	return &AnalyzeHandler{
		analysisRepo: repository.NewPostgresAnalysisRepository(pool),
	}
}

func (h *AnalyzeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload AnalyzePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	if err := payload.Validate(); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	slog.InfoContext(ctx, "processing analyze task",
		"owner", payload.Owner,
		"repo", payload.Repo,
	)

	if err := h.analyze(ctx, payload); err != nil {
		return err
	}

	slog.InfoContext(ctx, "analyze task completed",
		"owner", payload.Owner,
		"repo", payload.Repo,
	)

	return nil
}

func (h *AnalyzeHandler) analyze(ctx context.Context, payload AnalyzePayload) error {
	repoURL := fmt.Sprintf("https://github.com/%s/%s", payload.Owner, payload.Repo)

	gitSrc, err := source.NewGitSource(ctx, repoURL, nil)
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}
	defer gitSrc.Close()

	result, err := parser.Scan(ctx, gitSrc)
	if err != nil {
		return fmt.Errorf("scan repository: %w", err)
	}

	if result.Inventory == nil {
		slog.WarnContext(ctx, "scan result has no inventory",
			"owner", payload.Owner,
			"repo", payload.Repo,
			"commit", gitSrc.CommitSHA(),
		)
	}

	if err := h.analysisRepo.SaveAnalysisResult(ctx, repository.SaveAnalysisResultParams{
		Branch:    gitSrc.Branch(),
		CommitSHA: gitSrc.CommitSHA(),
		Owner:     payload.Owner,
		Repo:      payload.Repo,
		Result:    result,
	}); err != nil {
		return fmt.Errorf("save results: %w", err)
	}

	return nil
}
