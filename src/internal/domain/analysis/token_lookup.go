package analysis

import (
	"context"
	"errors"
)

// ErrTokenNotFound indicates the OAuth token does not exist for the user/provider.
// This is an expected condition (user hasn't connected OAuth) and should trigger graceful degradation.
var ErrTokenNotFound = errors.New("oauth token not found")

// TokenLookup retrieves OAuth tokens for repository access.
//
// Implementations should return:
//   - ErrTokenNotFound: when token doesn't exist (expected, triggers graceful degradation)
//   - Other errors: infrastructure failures (should fail the operation)
type TokenLookup interface {
	GetOAuthToken(ctx context.Context, userID string, provider string) (string, error)
}
