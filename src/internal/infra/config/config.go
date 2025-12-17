package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL   string
	EncryptionKey string
	RedisURL      string
}

func Load() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		return nil, errors.New("ENCRYPTION_KEY is required")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, errors.New("REDIS_URL is required")
	}

	return &Config{
		DatabaseURL:   databaseURL,
		EncryptionKey: encryptionKey,
		RedisURL:      redisURL,
	}, nil
}
