-- name: UpsertCodebase :one
INSERT INTO codebases (host, owner, name, default_branch)
VALUES ($1, $2, $3, $4)
ON CONFLICT (host, owner, name)
DO UPDATE SET
    default_branch = COALESCE(EXCLUDED.default_branch, codebases.default_branch),
    updated_at = now()
RETURNING *;

-- name: GetCodebaseByID :one
SELECT * FROM codebases WHERE id = $1;

-- name: CreateAnalysis :one
INSERT INTO analyses (codebase_id, commit_sha, branch_name, status, started_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateAnalysisCompleted :exec
UPDATE analyses
SET status = 'completed', total_suites = $2, total_tests = $3, completed_at = $4
WHERE id = $1;

-- name: UpdateAnalysisFailed :exec
UPDATE analyses
SET status = 'failed', error_message = $2, completed_at = $3
WHERE id = $1;

-- name: CreateTestSuite :one
INSERT INTO test_suites (analysis_id, parent_id, name, file_path, line_number, framework, depth)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: CreateTestCase :one
INSERT INTO test_cases (suite_id, name, line_number, status, tags, modifier)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTestSuitesByAnalysisID :many
SELECT * FROM test_suites WHERE analysis_id = $1 ORDER BY file_path, line_number;

-- name: GetTestCasesBySuiteID :many
SELECT * FROM test_cases WHERE suite_id = $1 ORDER BY line_number;
