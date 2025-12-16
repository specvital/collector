package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/specvital/collector/internal/domain/analysis"
	uc "github.com/specvital/collector/internal/usecase/analysis"
)

const TypeAnalyze = "analysis:analyze"

type AnalyzePayload struct {
	Owner  string  `json:"owner"`
	Repo   string  `json:"repo"`
	UserID *string `json:"user_id,omitempty"`
}

type AnalyzeHandler struct {
	analyzeUC *uc.AnalyzeUseCase
}

func NewAnalyzeHandler(analyzeUC *uc.AnalyzeUseCase) *AnalyzeHandler {
	return &AnalyzeHandler{analyzeUC: analyzeUC}
}

func (h *AnalyzeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var payload AnalyzePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	slog.InfoContext(ctx, "processing analyze task",
		"owner", payload.Owner,
		"repo", payload.Repo,
	)

	req := analysis.AnalyzeRequest{
		Owner:  payload.Owner,
		Repo:   payload.Repo,
		UserID: payload.UserID,
	}

	if err := h.analyzeUC.Execute(ctx, req); err != nil {
		slog.ErrorContext(ctx, "analyze task failed",
			"owner", payload.Owner,
			"repo", payload.Repo,
			"error", err,
		)
		return err
	}

	slog.InfoContext(ctx, "analyze task completed",
		"owner", payload.Owner,
		"repo", payload.Repo,
	)

	return nil
}
