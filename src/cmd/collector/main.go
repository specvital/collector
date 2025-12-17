// Package main provides the collector entry point.
// Deprecated: Use cmd/worker instead. This entry point is maintained for backward compatibility.
package main

import (
	"log/slog"
	"os"

	"github.com/specvital/collector/internal/app/bootstrap"
	"github.com/specvital/collector/internal/infra/config"

	_ "github.com/specvital/core/pkg/parser/strategies/all"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := bootstrap.StartWorker(bootstrap.WorkerConfig{
		ServiceName:   "collector",
		DatabaseURL:   cfg.DatabaseURL,
		EncryptionKey: cfg.EncryptionKey,
		RedisURL:      cfg.RedisURL,
	}); err != nil {
		slog.Error("collector failed", "error", err)
		os.Exit(1)
	}
}
