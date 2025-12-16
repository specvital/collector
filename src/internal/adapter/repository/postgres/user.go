package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/specvital/collector/internal/infra/db"
)

var ErrTokenNotFound = errors.New("oauth token not found")

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) GetOAuthToken(ctx context.Context, userID string, provider string) (string, error) {
	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}
	if provider == "" {
		return "", fmt.Errorf("provider is required")
	}

	var pgUserID pgtype.UUID
	if err := pgUserID.Scan(userID); err != nil {
		return "", fmt.Errorf("invalid user ID format: %w", err)
	}

	queries := db.New(r.pool)
	account, err := queries.GetOAuthAccountByUserAndProvider(ctx, db.GetOAuthAccountByUserAndProviderParams{
		UserID:   pgUserID,
		Provider: db.OauthProvider(provider),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrTokenNotFound
		}
		return "", fmt.Errorf("query oauth account: %w", err)
	}

	if !account.AccessToken.Valid || account.AccessToken.String == "" {
		return "", ErrTokenNotFound
	}

	return account.AccessToken.String, nil
}
