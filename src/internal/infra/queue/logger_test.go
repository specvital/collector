package queue

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestSlogAdapter_LogsToStdout(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	adapter := NewSlogAdapter(logger)

	tests := []struct {
		name      string
		logFunc   func(args ...any)
		expected  string
		level     string
		checkAttr bool
		attrKey   string
		attrValue string
	}{
		{
			name:     "debug level",
			logFunc:  adapter.Debug,
			expected: "test debug message",
			level:    "DEBUG",
		},
		{
			name:     "info level",
			logFunc:  adapter.Info,
			expected: "test info message",
			level:    "INFO",
		},
		{
			name:     "warn level",
			logFunc:  adapter.Warn,
			expected: "test warn message",
			level:    "WARN",
		},
		{
			name:     "error level",
			logFunc:  adapter.Error,
			expected: "test error message",
			level:    "ERROR",
		},
		{
			name:      "fatal level",
			logFunc:   adapter.Fatal,
			expected:  "test fatal message",
			level:     "ERROR",
			checkAttr: true,
			attrKey:   "severity",
			attrValue: "fatal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.expected)

			output := buf.String()
			if output == "" {
				t.Fatal("expected log output, got empty string")
			}

			var logEntry map[string]any
			if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
				t.Fatalf("failed to parse log as JSON: %v\nOutput: %s", err, output)
			}

			msg, ok := logEntry["msg"].(string)
			if !ok {
				t.Fatalf("log entry missing 'msg' field: %v", logEntry)
			}
			if msg != tt.expected {
				t.Errorf("expected message %q, got %q", tt.expected, msg)
			}

			level, ok := logEntry["level"].(string)
			if !ok {
				t.Fatalf("log entry missing 'level' field: %v", logEntry)
			}
			if !strings.EqualFold(level, tt.level) {
				t.Errorf("expected level %q, got %q", tt.level, level)
			}

			if tt.checkAttr {
				attrVal, exists := logEntry[tt.attrKey]
				if !exists {
					t.Errorf("expected attribute %q not found in log entry", tt.attrKey)
				} else if attrVal != tt.attrValue {
					t.Errorf("expected attribute %q=%q, got %q", tt.attrKey, tt.attrValue, attrVal)
				}
			}
		})
	}
}

func TestSlogAdapter_MultipleArguments(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	adapter := NewSlogAdapter(logger)

	adapter.Info("processing task", " ", "task_id=123")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log as JSON: %v", err)
	}

	msg := logEntry["msg"].(string)
	if !strings.Contains(msg, "processing task") {
		t.Errorf("expected message to contain 'processing task', got %q", msg)
	}
	if !strings.Contains(msg, "task_id=123") {
		t.Errorf("expected message to contain 'task_id=123', got %q", msg)
	}
}
