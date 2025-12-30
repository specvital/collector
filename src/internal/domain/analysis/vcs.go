package analysis

import (
	"context"
	"time"
)

type VCS interface {
	Clone(ctx context.Context, url string, token *string) (Source, error)
	// GetHeadCommit returns the HEAD commit SHA of the default branch without cloning.
	GetHeadCommit(ctx context.Context, url string, token *string) (string, error)
}

type Source interface {
	Branch() string
	CommitSHA() string
	CommittedAt() time.Time
	Close(ctx context.Context) error
	// VerifyCommitExists checks if a commit SHA exists in the remote repository
	// by running "git fetch --depth 1 origin <sha>" on the cloned repository.
	// Returns true if the commit exists, false if not found (e.g., "not our ref").
	// This enables reanalysis verification without API calls.
	VerifyCommitExists(ctx context.Context, sha string) (bool, error)
}

type RepoInfo struct {
	ExternalRepoID string
	Name           string
	Owner          string
}

type VCSAPIClient interface {
	// Returns ErrRepoNotFound if the repository does not exist.
	GetRepoInfo(ctx context.Context, host, owner, repo string, token *string) (RepoInfo, error)
}
