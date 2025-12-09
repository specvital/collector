package jobs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/collector/internal/db"
	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser"
	"github.com/specvital/core/pkg/source"

	_ "github.com/specvital/core/pkg/parser/strategies/all"
)

const (
	TypeAnalyze = "analysis:analyze"
	defaultHost = "github.com"
)

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
	pool *pgxpool.Pool
}

func NewAnalyzeHandler(pool *pgxpool.Pool) *AnalyzeHandler {
	return &AnalyzeHandler{pool: pool}
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

	if err := h.saveResults(ctx, payload, gitSrc.CommitSHA(), gitSrc.Branch(), result); err != nil {
		return fmt.Errorf("save results: %w", err)
	}

	return nil
}

func (h *AnalyzeHandler) saveResults(ctx context.Context, payload AnalyzePayload, commitSHA, branch string, result *parser.ScanResult) error {
	tx, err := h.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(context.Background())

	queries := db.New(tx)
	now := time.Now()

	codebase, err := queries.UpsertCodebase(ctx, db.UpsertCodebaseParams{
		Host:          defaultHost,
		Owner:         payload.Owner,
		Name:          payload.Repo,
		DefaultBranch: pgtype.Text{String: branch, Valid: branch != ""},
	})
	if err != nil {
		return fmt.Errorf("upsert codebase: %w", err)
	}

	analysis, err := queries.CreateAnalysis(ctx, db.CreateAnalysisParams{
		CodebaseID: codebase.ID,
		CommitSha:  commitSHA,
		BranchName: pgtype.Text{String: branch, Valid: branch != ""},
		Status:     db.AnalysisStatusRunning,
		StartedAt:  pgtype.Timestamptz{Time: now, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("create analysis: %w", err)
	}

	var totalSuites, totalTests int
	if result.Inventory != nil {
		for _, file := range result.Inventory.Files {
			suites, tests, err := h.saveTestFile(ctx, queries, analysis.ID, file, 0)
			if err != nil {
				return fmt.Errorf("save test file %s: %w", file.Path, err)
			}
			totalSuites += suites
			totalTests += tests
		}
	} else {
		slog.WarnContext(ctx, "scan result has no inventory",
			"owner", payload.Owner,
			"repo", payload.Repo,
			"commit", commitSHA,
		)
	}

	if err := queries.UpdateAnalysisCompleted(ctx, db.UpdateAnalysisCompletedParams{
		ID:          analysis.ID,
		TotalSuites: int32(totalSuites),
		TotalTests:  int32(totalTests),
		CompletedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}); err != nil {
		return fmt.Errorf("update analysis: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (h *AnalyzeHandler) saveTestFile(ctx context.Context, queries *db.Queries, analysisID pgtype.UUID, file domain.TestFile, depth int) (suites int, tests int, err error) {
	for _, suite := range file.Suites {
		s, t, err := h.saveSuite(ctx, queries, analysisID, pgtype.UUID{}, file, suite, depth)
		if err != nil {
			return 0, 0, err
		}
		suites += s
		tests += t
	}

	if len(file.Tests) > 0 {
		implicitSuite, err := queries.CreateTestSuite(ctx, db.CreateTestSuiteParams{
			AnalysisID: analysisID,
			ParentID:   pgtype.UUID{},
			Name:       file.Path,
			FilePath:   file.Path,
			LineNumber: pgtype.Int4{Int32: 1, Valid: true},
			Framework:  pgtype.Text{String: file.Framework, Valid: file.Framework != ""},
			Depth:      int32(depth),
		})
		if err != nil {
			return 0, 0, fmt.Errorf("create implicit suite: %w", err)
		}
		suites++

		for _, test := range file.Tests {
			if err := h.saveTest(ctx, queries, implicitSuite.ID, test); err != nil {
				return 0, 0, err
			}
			tests++
		}
	}

	return suites, tests, nil
}

func (h *AnalyzeHandler) saveSuite(ctx context.Context, queries *db.Queries, analysisID, parentID pgtype.UUID, file domain.TestFile, suite domain.TestSuite, depth int) (suites int, tests int, err error) {
	created, err := queries.CreateTestSuite(ctx, db.CreateTestSuiteParams{
		AnalysisID: analysisID,
		ParentID:   parentID,
		Name:       suite.Name,
		FilePath:   file.Path,
		LineNumber: pgtype.Int4{Int32: int32(suite.Location.StartLine), Valid: true},
		Framework:  pgtype.Text{String: file.Framework, Valid: file.Framework != ""},
		Depth:      int32(depth),
	})
	if err != nil {
		return 0, 0, fmt.Errorf("create suite: %w", err)
	}
	suites = 1

	for _, test := range suite.Tests {
		if err := h.saveTest(ctx, queries, created.ID, test); err != nil {
			return 0, 0, err
		}
		tests++
	}

	for _, nested := range suite.Suites {
		s, t, err := h.saveSuite(ctx, queries, analysisID, created.ID, file, nested, depth+1)
		if err != nil {
			return 0, 0, err
		}
		suites += s
		tests += t
	}

	return suites, tests, nil
}

func (h *AnalyzeHandler) saveTest(ctx context.Context, queries *db.Queries, suiteID pgtype.UUID, test domain.Test) error {
	status := mapTestStatus(test.Status)
	_, err := queries.CreateTestCase(ctx, db.CreateTestCaseParams{
		SuiteID:    suiteID,
		Name:       test.Name,
		LineNumber: pgtype.Int4{Int32: int32(test.Location.StartLine), Valid: true},
		Status:     status,
		Tags:       []byte("[]"),
	})
	if err != nil {
		return fmt.Errorf("create test case: %w", err)
	}
	return nil
}

func mapTestStatus(status domain.TestStatus) db.TestStatus {
	switch status {
	case domain.TestStatusSkipped:
		return db.TestStatusSkipped
	case domain.TestStatusPending, domain.TestStatusFixme:
		return db.TestStatusTodo
	default:
		return db.TestStatusActive
	}
}
