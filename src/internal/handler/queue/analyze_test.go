package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/specvital/collector/internal/domain/analysis"
	uc "github.com/specvital/collector/internal/usecase/analysis"
)

// Mock implementations for testing

type mockVCS struct {
	cloneFn func(ctx context.Context, url string, token *string) (analysis.Source, error)
}

func (m *mockVCS) Clone(ctx context.Context, url string, token *string) (analysis.Source, error) {
	if m.cloneFn != nil {
		return m.cloneFn(ctx, url, token)
	}
	return nil, nil
}

type mockSource struct {
	branchFn    func() string
	commitSHAFn func() string
	closeFn     func(ctx context.Context) error
}

func (m *mockSource) Branch() string {
	if m.branchFn != nil {
		return m.branchFn()
	}
	return "main"
}

func (m *mockSource) CommitSHA() string {
	if m.commitSHAFn != nil {
		return m.commitSHAFn()
	}
	return "abc123"
}

func (m *mockSource) Close(ctx context.Context) error {
	if m.closeFn != nil {
		return m.closeFn(ctx)
	}
	return nil
}

type mockParser struct {
	scanFn func(ctx context.Context, src analysis.Source) (*analysis.Inventory, error)
}

func (m *mockParser) Scan(ctx context.Context, src analysis.Source) (*analysis.Inventory, error) {
	if m.scanFn != nil {
		return m.scanFn(ctx, src)
	}
	return &analysis.Inventory{Files: []analysis.TestFile{}}, nil
}

type mockRepository struct {
	createAnalysisRecordFn  func(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error)
	recordFailureFn         func(ctx context.Context, analysisID analysis.UUID, errMessage string) error
	saveAnalysisInventoryFn func(ctx context.Context, params analysis.SaveAnalysisInventoryParams) error
}

func (m *mockRepository) CreateAnalysisRecord(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error) {
	if m.createAnalysisRecordFn != nil {
		return m.createAnalysisRecordFn(ctx, params)
	}
	return analysis.NewUUID(), nil
}

func (m *mockRepository) RecordFailure(ctx context.Context, analysisID analysis.UUID, errMessage string) error {
	if m.recordFailureFn != nil {
		return m.recordFailureFn(ctx, analysisID, errMessage)
	}
	return nil
}

func (m *mockRepository) SaveAnalysisInventory(ctx context.Context, params analysis.SaveAnalysisInventoryParams) error {
	if m.saveAnalysisInventoryFn != nil {
		return m.saveAnalysisInventoryFn(ctx, params)
	}
	return nil
}

// Test helper functions

func newSuccessfulMocks() (*mockRepository, *mockVCS, *mockParser) {
	src := &mockSource{
		branchFn:    func() string { return "main" },
		commitSHAFn: func() string { return "abc123" },
		closeFn:     func(ctx context.Context) error { return nil },
	}

	vcs := &mockVCS{
		cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
			return src, nil
		},
	}

	repo := &mockRepository{
		createAnalysisRecordFn: func(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error) {
			return analysis.NewUUID(), nil
		},
		saveAnalysisInventoryFn: func(ctx context.Context, params analysis.SaveAnalysisInventoryParams) error {
			return nil
		},
	}

	parser := &mockParser{
		scanFn: func(ctx context.Context, src analysis.Source) (*analysis.Inventory, error) {
			return &analysis.Inventory{Files: []analysis.TestFile{}}, nil
		},
	}

	return repo, vcs, parser
}

// Tests

func TestNewAnalyzeHandler(t *testing.T) {
	repo, vcs, parser := newSuccessfulMocks()
	analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)

	handler := NewAnalyzeHandler(analyzeUC)

	if handler == nil {
		t.Error("expected handler, got nil")
	}
	if handler.analyzeUC == nil {
		t.Error("expected handler.analyzeUC to be set, got nil")
	}
}

func TestAnalyzeHandler_ProcessTask(t *testing.T) {
	tests := []struct {
		name        string
		payload     any
		setupMocks  func() (*mockRepository, *mockVCS, *mockParser)
		wantErr     bool
		errContains string
	}{
		{
			name: "success case - valid payload and use case succeeds",
			payload: AnalyzePayload{
				Owner: "octocat",
				Repo:  "Hello-World",
			},
			setupMocks: func() (*mockRepository, *mockVCS, *mockParser) {
				return newSuccessfulMocks()
			},
			wantErr: false,
		},
		{
			name: "clone failed - VCS clone returns error",
			payload: AnalyzePayload{
				Owner: "testowner",
				Repo:  "testrepo",
			},
			setupMocks: func() (*mockRepository, *mockVCS, *mockParser) {
				repo, _, parser := newSuccessfulMocks()
				vcs := &mockVCS{
					cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
						return nil, errors.New("git clone failed")
					},
				}
				return repo, vcs, parser
			},
			wantErr: true,
		},
		{
			name: "scan failed - parser returns error",
			payload: AnalyzePayload{
				Owner: "testowner",
				Repo:  "testrepo",
			},
			setupMocks: func() (*mockRepository, *mockVCS, *mockParser) {
				repo, vcs, _ := newSuccessfulMocks()

				testAnalysisID := analysis.NewUUID()
				repo.createAnalysisRecordFn = func(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error) {
					return testAnalysisID, nil
				}
				repo.recordFailureFn = func(ctx context.Context, analysisID analysis.UUID, errMessage string) error {
					return nil
				}

				parser := &mockParser{
					scanFn: func(ctx context.Context, src analysis.Source) (*analysis.Inventory, error) {
						return nil, errors.New("parser error")
					},
				}

				return repo, vcs, parser
			},
			wantErr: true,
		},
		{
			name: "save failed - repository save returns error",
			payload: AnalyzePayload{
				Owner: "testowner",
				Repo:  "testrepo",
			},
			setupMocks: func() (*mockRepository, *mockVCS, *mockParser) {
				repo, vcs, parser := newSuccessfulMocks()

				testAnalysisID := analysis.NewUUID()
				repo.createAnalysisRecordFn = func(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error) {
					return testAnalysisID, nil
				}
				repo.recordFailureFn = func(ctx context.Context, analysisID analysis.UUID, errMessage string) error {
					return nil
				}
				repo.saveAnalysisInventoryFn = func(ctx context.Context, params analysis.SaveAnalysisInventoryParams) error {
					return errors.New("database save error")
				}

				return repo, vcs, parser
			},
			wantErr: true,
		},
		{
			name: "invalid input - empty owner",
			payload: AnalyzePayload{
				Owner: "",
				Repo:  "testrepo",
			},
			setupMocks: func() (*mockRepository, *mockVCS, *mockParser) {
				return newSuccessfulMocks()
			},
			wantErr: true,
		},
		{
			name: "invalid input - empty repo",
			payload: AnalyzePayload{
				Owner: "testowner",
				Repo:  "",
			},
			setupMocks: func() (*mockRepository, *mockVCS, *mockParser) {
				return newSuccessfulMocks()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, vcs, parser := tt.setupMocks()
			analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
			handler := NewAnalyzeHandler(analyzeUC)

			payloadBytes, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			task := asynq.NewTask(TypeAnalyze, payloadBytes)

			err = handler.ProcessTask(context.Background(), task)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if tt.errContains != "" && err != nil {
					if !containsString(err.Error(), tt.errContains) {
						t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestAnalyzeHandler_ProcessTask_InvalidPayload(t *testing.T) {
	tests := []struct {
		name        string
		payload     []byte
		errContains string
	}{
		{
			name:        "malformed JSON - invalid syntax",
			payload:     []byte(`{invalid json`),
			errContains: "unmarshal payload",
		},
		{
			name:        "invalid JSON type - string instead of object",
			payload:     []byte(`"just a string"`),
			errContains: "unmarshal payload",
		},
		{
			name:        "empty payload - empty bytes",
			payload:     []byte(``),
			errContains: "unmarshal payload",
		},
		{
			name:        "null payload - JSON null value",
			payload:     []byte(`null`),
			errContains: "",
		},
		{
			name:        "wrong field types - integer and boolean",
			payload:     []byte(`{"owner": 123, "repo": true}`),
			errContains: "unmarshal payload",
		},
		{
			name:        "array instead of object",
			payload:     []byte(`["owner", "repo"]`),
			errContains: "unmarshal payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, vcs, parser := newSuccessfulMocks()
			analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
			handler := NewAnalyzeHandler(analyzeUC)

			task := asynq.NewTask(TypeAnalyze, tt.payload)

			err := handler.ProcessTask(context.Background(), task)

			if err == nil && tt.errContains != "" {
				t.Error("expected error, got nil")
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
			}
		})
	}
}

func TestAnalyzeHandler_ProcessTask_ServiceInvocation(t *testing.T) {
	t.Run("should pass correct parameters to use case", func(t *testing.T) {
		var capturedOwner, capturedRepo string
		repo, _, parser := newSuccessfulMocks()

		src := &mockSource{
			branchFn:    func() string { return "main" },
			commitSHAFn: func() string { return "abc123" },
			closeFn:     func(ctx context.Context) error { return nil },
		}

		vcs := &mockVCS{
			cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
				// Extract owner/repo from URL
				if containsString(url, "test-owner/test-repo") {
					capturedOwner = "test-owner"
					capturedRepo = "test-repo"
				}
				return src, nil
			},
		}

		analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
		handler := NewAnalyzeHandler(analyzeUC)

		payload := AnalyzePayload{
			Owner: "test-owner",
			Repo:  "test-repo",
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		task := asynq.NewTask(TypeAnalyze, payloadBytes)

		err = handler.ProcessTask(context.Background(), task)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedOwner != "test-owner" {
			t.Errorf("expected owner 'test-owner', got '%s'", capturedOwner)
		}
		if capturedRepo != "test-repo" {
			t.Errorf("expected repo 'test-repo', got '%s'", capturedRepo)
		}
	})

	t.Run("should not call use case when unmarshal fails", func(t *testing.T) {
		useCaseCalled := false
		repo, _, parser := newSuccessfulMocks()

		vcs := &mockVCS{
			cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
				useCaseCalled = true
				return nil, errors.New("should not be called")
			},
		}

		analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
		handler := NewAnalyzeHandler(analyzeUC)

		task := asynq.NewTask(TypeAnalyze, []byte(`invalid json`))

		err := handler.ProcessTask(context.Background(), task)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if useCaseCalled {
			t.Error("use case should not be called when payload unmarshal fails")
		}
	})
}

func TestAnalyzeHandler_ProcessTask_ContextPropagation(t *testing.T) {
	t.Run("should propagate context to use case", func(t *testing.T) {
		type ctxKey string
		testKey := ctxKey("test-key")
		testValue := "test-value"

		var capturedCtx context.Context
		repo, _, parser := newSuccessfulMocks()

		src := &mockSource{
			branchFn:    func() string { return "main" },
			commitSHAFn: func() string { return "abc123" },
			closeFn:     func(ctx context.Context) error { return nil },
		}

		vcs := &mockVCS{
			cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
				capturedCtx = ctx
				return src, nil
			},
		}

		analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
		handler := NewAnalyzeHandler(analyzeUC)

		payload := AnalyzePayload{
			Owner: "owner",
			Repo:  "repo",
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		task := asynq.NewTask(TypeAnalyze, payloadBytes)
		ctx := context.WithValue(context.Background(), testKey, testValue)

		err = handler.ProcessTask(ctx, task)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedCtx == nil {
			t.Fatal("context was not propagated to use case")
		}
		if capturedCtx.Value(testKey) != testValue {
			t.Errorf("expected context value '%s', got '%v'", testValue, capturedCtx.Value(testKey))
		}
	})

	t.Run("should propagate cancelled context", func(t *testing.T) {
		repo, _, parser := newSuccessfulMocks()
		vcs := &mockVCS{
			cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
				// Should not reach here because semaphore.Acquire will fail first
				return nil, ctx.Err()
			},
		}

		analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
		handler := NewAnalyzeHandler(analyzeUC)

		payload := AnalyzePayload{
			Owner: "owner",
			Repo:  "repo",
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		task := asynq.NewTask(TypeAnalyze, payloadBytes)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = handler.ProcessTask(ctx, task)

		if err == nil {
			t.Error("expected error from cancelled context, got nil")
		}
		// Verify that the error is related to context cancellation
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected error to wrap context.Canceled, got %v", err)
		}
	})
}

func TestAnalyzeHandler_ProcessTask_ErrorPropagation(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() (*mockRepository, *mockVCS, *mockParser)
		wantError error
	}{
		{
			name: "clone failed error",
			setupMock: func() (*mockRepository, *mockVCS, *mockParser) {
				repo, _, parser := newSuccessfulMocks()
				vcs := &mockVCS{
					cloneFn: func(ctx context.Context, url string, token *string) (analysis.Source, error) {
						return nil, errors.New("clone error")
					},
				}
				return repo, vcs, parser
			},
			wantError: uc.ErrCloneFailed,
		},
		{
			name: "scan failed error",
			setupMock: func() (*mockRepository, *mockVCS, *mockParser) {
				repo, vcs, _ := newSuccessfulMocks()

				testAnalysisID := analysis.NewUUID()
				repo.createAnalysisRecordFn = func(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error) {
					return testAnalysisID, nil
				}
				repo.recordFailureFn = func(ctx context.Context, analysisID analysis.UUID, errMessage string) error {
					return nil
				}

				parser := &mockParser{
					scanFn: func(ctx context.Context, src analysis.Source) (*analysis.Inventory, error) {
						return nil, errors.New("scan error")
					},
				}

				return repo, vcs, parser
			},
			wantError: uc.ErrScanFailed,
		},
		{
			name: "save failed error",
			setupMock: func() (*mockRepository, *mockVCS, *mockParser) {
				repo, vcs, parser := newSuccessfulMocks()

				testAnalysisID := analysis.NewUUID()
				repo.createAnalysisRecordFn = func(ctx context.Context, params analysis.CreateAnalysisRecordParams) (analysis.UUID, error) {
					return testAnalysisID, nil
				}
				repo.recordFailureFn = func(ctx context.Context, analysisID analysis.UUID, errMessage string) error {
					return nil
				}
				repo.saveAnalysisInventoryFn = func(ctx context.Context, params analysis.SaveAnalysisInventoryParams) error {
					return errors.New("save error")
				}

				return repo, vcs, parser
			},
			wantError: uc.ErrSaveFailed,
		},
		{
			name: "invalid input error",
			setupMock: func() (*mockRepository, *mockVCS, *mockParser) {
				return newSuccessfulMocks()
			},
			wantError: analysis.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, vcs, parser := tt.setupMock()
			analyzeUC := uc.NewAnalyzeUseCase(repo, vcs, parser)
			handler := NewAnalyzeHandler(analyzeUC)

			var payload AnalyzePayload
			if tt.wantError == analysis.ErrInvalidInput {
				payload = AnalyzePayload{Owner: "", Repo: "repo"}
			} else {
				payload = AnalyzePayload{Owner: "owner", Repo: "repo"}
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				t.Fatalf("failed to marshal payload: %v", err)
			}
			task := asynq.NewTask(TypeAnalyze, payloadBytes)

			err = handler.ProcessTask(context.Background(), task)

			if err == nil {
				t.Errorf("expected error %v, got nil", tt.wantError)
				return
			}
			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error to wrap %v, got %v", tt.wantError, err)
			}
		})
	}
}

// containsString checks if s contains substr
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
