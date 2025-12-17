package queue

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/hibiken/asynq"
)

const (
	DefaultConcurrency     = 5
	DefaultShutdownTimeout = 30 * time.Second
)

type ServerConfig struct {
	Concurrency     int
	RedisURL        string
	ShutdownTimeout time.Duration
	Logger          *slog.Logger
}

func NewServer(cfg ServerConfig) (*asynq.Server, error) {
	concurrency := cfg.Concurrency
	if concurrency <= 0 {
		concurrency = DefaultConcurrency
	}

	shutdownTimeout := cfg.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = DefaultShutdownTimeout
	}

	opt, err := asynq.ParseRedisURI(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis URI: %w", err)
	}

	asynqConfig := asynq.Config{
		Concurrency:     concurrency,
		ShutdownTimeout: shutdownTimeout,
	}

	// Use custom slog adapter to ensure logs go to stdout as JSON
	if cfg.Logger != nil {
		asynqConfig.Logger = NewSlogAdapter(cfg.Logger)
	}

	return asynq.NewServer(opt, asynqConfig), nil
}

func NewServeMux() *asynq.ServeMux {
	return asynq.NewServeMux()
}
