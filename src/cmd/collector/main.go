package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/specvital/collector/internal/infra/config"
	"github.com/specvital/collector/internal/infra/db"
	"github.com/specvital/collector/internal/infra/queue"
	"github.com/specvital/collector/internal/jobs"
	"github.com/specvital/collector/internal/repository"
	"github.com/specvital/collector/internal/service"
)

const (
	workerConcurrency       = 5
	gracefulShutdownTimeout = 60 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting collector")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	slog.Info("config loaded",
		"database_url", maskURL(cfg.DatabaseURL),
		"redis_url", maskURL(cfg.RedisURL),
	)

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	slog.Info("postgres connected")

	srv, err := queue.NewServer(queue.ServerConfig{
		RedisURL:        cfg.RedisURL,
		Concurrency:     workerConcurrency,
		ShutdownTimeout: gracefulShutdownTimeout,
	})
	if err != nil {
		slog.Error("failed to create queue server", "error", err)
		pool.Close()
		os.Exit(1)
	}

	analysisRepo := repository.NewPostgresAnalysisRepository(pool)
	analysisSvc := service.NewAnalysisService(analysisRepo)

	mux := queue.NewServeMux()
	analyzeHandler := jobs.NewAnalyzeHandler(analysisSvc)
	mux.HandleFunc(jobs.TypeAnalyze, analyzeHandler.ProcessTask)

	slog.Info("worker starting", "concurrency", workerConcurrency)
	if err := srv.Start(mux); err != nil {
		slog.Error("failed to start worker", "error", err)
		pool.Close()
		os.Exit(1)
	}
	slog.Info("worker ready", "concurrency", workerConcurrency)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	sig := <-shutdown
	slog.Info("shutdown signal received, waiting for in-flight jobs...", "signal", sig.String())

	srv.Shutdown()
	slog.Info("worker shutdown complete")

	pool.Close()
	slog.Info("collector shutdown complete")
}

func maskURL(url string) string {
	if len(url) > 20 {
		return url[:20] + "..."
	}
	return url
}
