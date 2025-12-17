package queue

import (
	"fmt"
	"log/slog"
)

// SlogAdapter adapts slog.Logger to asynq.Logger interface.
// This ensures Asynq logs are written to stdout as JSON, not stderr.
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter creates a new Asynq logger that wraps slog.
func NewSlogAdapter(logger *slog.Logger) *SlogAdapter {
	return &SlogAdapter{logger: logger}
}

func (l *SlogAdapter) Debug(args ...any) {
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *SlogAdapter) Info(args ...any) {
	l.logger.Info(fmt.Sprint(args...))
}

func (l *SlogAdapter) Warn(args ...any) {
	l.logger.Warn(fmt.Sprint(args...))
}

func (l *SlogAdapter) Error(args ...any) {
	l.logger.Error(fmt.Sprint(args...))
}

// Fatal logs at Error level with severity=fatal attribute.
// Unlike log.Fatal, this does NOT call os.Exit.
func (l *SlogAdapter) Fatal(args ...any) {
	l.logger.Error(fmt.Sprint(args...), "severity", "fatal")
}
