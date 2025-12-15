package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/collector/internal/infra/db"
	"github.com/specvital/core/pkg/domain"
	"github.com/specvital/core/pkg/parser"
)

const defaultHost = "github.com"
const maxErrorMessageLength = 1000

func validateRepositoryInfo(owner, repo, commitSHA string) error {
	if owner == "" {
		return fmt.Errorf("%w: owner is required", ErrInvalidParams)
	}
	if repo == "" {
		return fmt.Errorf("%w: repo is required", ErrInvalidParams)
	}
	if commitSHA == "" {
		return fmt.Errorf("%w: commit SHA is required", ErrInvalidParams)
	}
	return nil
}

func truncateErrorMessage(msg string) string {
	if len(msg) <= maxErrorMessageLength {
		return msg
	}
	return msg[:maxErrorMessageLength-15] + "... (truncated)"
}

type CreateAnalysisRecordParams struct {
	Branch    string
	CommitSHA string
	Owner     string
	Repo      string
}

func (p CreateAnalysisRecordParams) Validate() error {
	return validateRepositoryInfo(p.Owner, p.Repo, p.CommitSHA)
}

type SaveAnalysisInventoryParams struct {
	AnalysisID pgtype.UUID
	Result     *parser.ScanResult
}

func (p SaveAnalysisInventoryParams) Validate() error {
	if !p.AnalysisID.Valid {
		return fmt.Errorf("%w: analysis ID is required", ErrInvalidParams)
	}
	if p.Result == nil {
		return fmt.Errorf("%w: result is required", ErrInvalidParams)
	}
	return nil
}

type SaveAnalysisResultParams struct {
	Branch    string
	CommitSHA string
	Owner     string
	Repo      string
	Result    *parser.ScanResult
}

func (p SaveAnalysisResultParams) Validate() error {
	if err := validateRepositoryInfo(p.Owner, p.Repo, p.CommitSHA); err != nil {
		return err
	}
	if p.Result == nil {
		return fmt.Errorf("%w: result is required", ErrInvalidParams)
	}
	return nil
}

type AnalysisRepository interface {
	CreateAnalysisRecord(ctx context.Context, params CreateAnalysisRecordParams) (pgtype.UUID, error)
	RecordFailure(ctx context.Context, analysisID pgtype.UUID, errMessage string) error
	SaveAnalysisInventory(ctx context.Context, params SaveAnalysisInventoryParams) error
	SaveAnalysisResult(ctx context.Context, params SaveAnalysisResultParams) error
}

type PostgresAnalysisRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAnalysisRepository(pool *pgxpool.Pool) *PostgresAnalysisRepository {
	return &PostgresAnalysisRepository{pool: pool}
}

func (r *PostgresAnalysisRepository) CreateAnalysisRecord(ctx context.Context, params CreateAnalysisRecordParams) (pgtype.UUID, error) {
	if err := params.Validate(); err != nil {
		return pgtype.UUID{}, err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, "failed to rollback transaction",
				"error", err,
				"owner", params.Owner,
				"repo", params.Repo,
			)
		}
	}()

	queries := db.New(tx)
	startedAt := time.Now()

	codebase, err := queries.UpsertCodebase(ctx, db.UpsertCodebaseParams{
		Host:          defaultHost,
		Owner:         params.Owner,
		Name:          params.Repo,
		DefaultBranch: pgtype.Text{String: params.Branch, Valid: params.Branch != ""},
	})
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("upsert codebase: %w", err)
	}

	analysis, err := queries.CreateAnalysis(ctx, db.CreateAnalysisParams{
		CodebaseID: codebase.ID,
		CommitSha:  params.CommitSHA,
		BranchName: pgtype.Text{String: params.Branch, Valid: params.Branch != ""},
		Status:     db.AnalysisStatusRunning,
		StartedAt:  pgtype.Timestamptz{Time: startedAt, Valid: true},
	})
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("create analysis: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return pgtype.UUID{}, fmt.Errorf("commit transaction: %w", err)
	}

	return analysis.ID, nil
}

func (r *PostgresAnalysisRepository) RecordFailure(ctx context.Context, analysisID pgtype.UUID, errMessage string) error {
	if !analysisID.Valid {
		return fmt.Errorf("%w: analysis ID is required", ErrInvalidParams)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, "failed to rollback transaction",
				"error", err,
				"analysis_id", analysisID,
			)
		}
	}()

	queries := db.New(tx)
	truncatedMsg := truncateErrorMessage(errMessage)
	if err := queries.UpdateAnalysisFailed(ctx, db.UpdateAnalysisFailedParams{
		ID:           analysisID,
		ErrorMessage: pgtype.Text{String: truncatedMsg, Valid: true},
		CompletedAt:  pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}); err != nil {
		return fmt.Errorf("update analysis failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresAnalysisRepository) SaveAnalysisInventory(ctx context.Context, params SaveAnalysisInventoryParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.ErrorContext(ctx, "failed to rollback transaction",
				"error", err,
				"analysis_id", params.AnalysisID,
			)
		}
	}()

	queries := db.New(tx)

	totalSuites, totalTests, err := r.saveInventory(ctx, queries, params.AnalysisID, params.Result.Inventory)
	if err != nil {
		return fmt.Errorf("save inventory: %w", err)
	}

	if err := queries.UpdateAnalysisCompleted(ctx, db.UpdateAnalysisCompletedParams{
		ID:          params.AnalysisID,
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

func (r *PostgresAnalysisRepository) SaveAnalysisResult(ctx context.Context, params SaveAnalysisResultParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	analysisID, err := r.CreateAnalysisRecord(ctx, CreateAnalysisRecordParams{
		Branch:    params.Branch,
		CommitSHA: params.CommitSHA,
		Owner:     params.Owner,
		Repo:      params.Repo,
	})
	if err != nil {
		return err
	}

	if err := r.SaveAnalysisInventory(ctx, SaveAnalysisInventoryParams{
		AnalysisID: analysisID,
		Result:     params.Result,
	}); err != nil {
		return err
	}

	return nil
}

func (r *PostgresAnalysisRepository) saveInventory(ctx context.Context, queries *db.Queries, analysisID pgtype.UUID, inventory *domain.Inventory) (int, int, error) {
	if inventory == nil {
		return 0, 0, nil
	}

	var totalSuites, totalTests int
	for _, file := range inventory.Files {
		suites, tests, err := r.saveTestFile(ctx, queries, analysisID, file, 0)
		if err != nil {
			return 0, 0, fmt.Errorf("save test file %s: %w", file.Path, err)
		}
		totalSuites += suites
		totalTests += tests
	}

	return totalSuites, totalTests, nil
}

func (r *PostgresAnalysisRepository) saveTestFile(ctx context.Context, queries *db.Queries, analysisID pgtype.UUID, file domain.TestFile, depth int) (int, int, error) {
	var totalSuites, totalTests int

	for _, suite := range file.Suites {
		suites, tests, err := r.saveSuite(ctx, queries, analysisID, pgtype.UUID{}, file, suite, depth)
		if err != nil {
			return 0, 0, err
		}
		totalSuites += suites
		totalTests += tests
	}

	if len(file.Tests) > 0 {
		implicitSuite, err := r.createImplicitSuite(ctx, queries, analysisID, file, depth)
		if err != nil {
			return 0, 0, err
		}
		totalSuites++

		for _, test := range file.Tests {
			if err := r.saveTest(ctx, queries, implicitSuite.ID, test); err != nil {
				return 0, 0, err
			}
			totalTests++
		}
	}

	return totalSuites, totalTests, nil
}

func (r *PostgresAnalysisRepository) createImplicitSuite(ctx context.Context, queries *db.Queries, analysisID pgtype.UUID, file domain.TestFile, depth int) (db.TestSuite, error) {
	suite, err := queries.CreateTestSuite(ctx, db.CreateTestSuiteParams{
		AnalysisID: analysisID,
		ParentID:   pgtype.UUID{},
		Name:       file.Path,
		FilePath:   file.Path,
		LineNumber: pgtype.Int4{Int32: 1, Valid: true},
		Framework:  pgtype.Text{String: file.Framework, Valid: file.Framework != ""},
		Depth:      int32(depth),
	})
	if err != nil {
		return db.TestSuite{}, fmt.Errorf("create implicit suite: %w", err)
	}
	return suite, nil
}

func (r *PostgresAnalysisRepository) saveSuite(ctx context.Context, queries *db.Queries, analysisID, parentID pgtype.UUID, file domain.TestFile, suite domain.TestSuite, depth int) (int, int, error) {
	name := truncateString(suite.Name, maxTestSuiteNameLength)

	created, err := queries.CreateTestSuite(ctx, db.CreateTestSuiteParams{
		AnalysisID: analysisID,
		ParentID:   parentID,
		Name:       name,
		FilePath:   file.Path,
		LineNumber: pgtype.Int4{Int32: int32(suite.Location.StartLine), Valid: true},
		Framework:  pgtype.Text{String: file.Framework, Valid: file.Framework != ""},
		Depth:      int32(depth),
	})
	if err != nil {
		return 0, 0, fmt.Errorf("create suite (name=%q, file=%s, line=%d): %w",
			truncateString(suite.Name, 100), file.Path, suite.Location.StartLine, err)
	}

	totalSuites := 1
	var totalTests int

	for _, test := range suite.Tests {
		if err := r.saveTest(ctx, queries, created.ID, test); err != nil {
			return 0, 0, err
		}
		totalTests++
	}

	for _, nested := range suite.Suites {
		suites, tests, err := r.saveSuite(ctx, queries, analysisID, created.ID, file, nested, depth+1)
		if err != nil {
			return 0, 0, err
		}
		totalSuites += suites
		totalTests += tests
	}

	return totalSuites, totalTests, nil
}

const (
	maxTestCaseNameLength  = 2000
	maxTestSuiteNameLength = 500
)

func (r *PostgresAnalysisRepository) saveTest(ctx context.Context, queries *db.Queries, suiteID pgtype.UUID, test domain.Test) error {
	status := mapTestStatus(test.Status)
	name := truncateString(test.Name, maxTestCaseNameLength)

	_, err := queries.CreateTestCase(ctx, db.CreateTestCaseParams{
		SuiteID:    suiteID,
		Name:       name,
		LineNumber: pgtype.Int4{Int32: int32(test.Location.StartLine), Valid: true},
		Status:     status,
		Tags:       []byte("[]"),
	})
	if err != nil {
		return fmt.Errorf("create test case (name=%q, line=%d): %w",
			truncateString(test.Name, 100), test.Location.StartLine, err)
	}
	return nil
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func mapTestStatus(status domain.TestStatus) db.TestStatus {
	switch status {
	case domain.TestStatusSkipped:
		return db.TestStatusSkipped
	case domain.TestStatusTodo, domain.TestStatusXfail:
		return db.TestStatusTodo
	default:
		return db.TestStatusActive
	}
}
