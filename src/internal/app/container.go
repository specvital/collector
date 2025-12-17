package app

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/collector/internal/adapter/parser"
	"github.com/specvital/collector/internal/adapter/repository/postgres"
	"github.com/specvital/collector/internal/adapter/vcs"
	"github.com/specvital/collector/internal/handler/queue"
	handlerscheduler "github.com/specvital/collector/internal/handler/scheduler"
	infraqueue "github.com/specvital/collector/internal/infra/queue"
	infrascheduler "github.com/specvital/collector/internal/infra/scheduler"
	uc "github.com/specvital/collector/internal/usecase/analysis"
	"github.com/specvital/collector/internal/usecase/autorefresh"
	"github.com/specvital/core/pkg/crypto"
)

const (
	schedulerLockKey = "scheduler:auto-refresh:lock"
	schedulerLockTTL = 10 * time.Minute
)

type ContainerConfig struct {
	EncryptionKey string
	Pool          *pgxpool.Pool
	RedisURL      string
}

func (c ContainerConfig) Validate() error {
	if c.Pool == nil {
		return fmt.Errorf("pool is required")
	}
	if c.RedisURL == "" {
		return fmt.Errorf("redis URL is required")
	}
	return nil
}

func (c ContainerConfig) ValidateWorker() error {
	if err := c.Validate(); err != nil {
		return err
	}
	if c.EncryptionKey == "" {
		return fmt.Errorf("encryption key is required")
	}
	return nil
}

type WorkerContainer struct {
	AnalyzeHandler *queue.AnalyzeHandler
	QueueClient    *infraqueue.Client
}

func NewWorkerContainer(cfg ContainerConfig) (*WorkerContainer, error) {
	if err := cfg.ValidateWorker(); err != nil {
		return nil, fmt.Errorf("invalid container config: %w", err)
	}

	encryptor, err := crypto.NewEncryptorFromBase64(cfg.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("create encryptor: %w", err)
	}

	analysisRepo := postgres.NewAnalysisRepository(cfg.Pool)
	userRepo := postgres.NewUserRepository(cfg.Pool, encryptor)
	gitVCS := vcs.NewGitVCS()
	coreParser := parser.NewCoreParser()
	analyzeUC := uc.NewAnalyzeUseCase(analysisRepo, gitVCS, coreParser, userRepo)
	analyzeHandler := queue.NewAnalyzeHandler(analyzeUC)

	queueClient, err := infraqueue.NewClient(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("create queue client: %w", err)
	}

	return &WorkerContainer{
		AnalyzeHandler: analyzeHandler,
		QueueClient:    queueClient,
	}, nil
}

func (c *WorkerContainer) Close() error {
	if c.QueueClient != nil {
		if err := c.QueueClient.Close(); err != nil {
			return fmt.Errorf("close queue client: %w", err)
		}
	}
	return nil
}

type SchedulerContainer struct {
	AutoRefreshHandler *handlerscheduler.AutoRefreshHandler
	Scheduler          *infrascheduler.Scheduler
	queueClient        *infraqueue.Client
	schedulerLock      *infrascheduler.DistributedLock
}

func NewSchedulerContainer(cfg ContainerConfig) (*SchedulerContainer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid container config: %w", err)
	}

	analysisRepo := postgres.NewAnalysisRepository(cfg.Pool)

	queueClient, err := infraqueue.NewClient(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("create queue client: %w", err)
	}

	schedulerLock, err := infrascheduler.NewDistributedLock(cfg.RedisURL, schedulerLockKey, schedulerLockTTL)
	if err != nil {
		queueClient.Close()
		return nil, fmt.Errorf("create scheduler lock: %w", err)
	}

	autoRefreshUC := autorefresh.NewAutoRefreshUseCase(analysisRepo, queueClient)
	autoRefreshHandler := handlerscheduler.NewAutoRefreshHandler(autoRefreshUC, schedulerLock)

	scheduler := infrascheduler.New()

	return &SchedulerContainer{
		AutoRefreshHandler: autoRefreshHandler,
		Scheduler:          scheduler,
		queueClient:        queueClient,
		schedulerLock:      schedulerLock,
	}, nil
}

func (c *SchedulerContainer) Close() error {
	var errs []error

	if c.schedulerLock != nil {
		if err := c.schedulerLock.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close scheduler lock: %w", err))
		}
	}

	if c.queueClient != nil {
		if err := c.queueClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close queue client: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close scheduler container: %v", errs)
	}
	return nil
}
