package analysis

import "context"

type VCS interface {
	Clone(ctx context.Context, url string, token *string) (Source, error)
	// GetHeadCommit returns the HEAD commit SHA of the default branch without cloning.
	GetHeadCommit(ctx context.Context, url string, token *string) (string, error)
}

type Source interface {
	Branch() string
	CommitSHA() string
	Close(ctx context.Context) error
	// VerifyCommitExists checks if a commit SHA exists in the remote repository
	// by running "git fetch --depth 1 origin <sha>" on the cloned repository.
	// Returns true if the commit exists, false if not found (e.g., "not our ref").
	// This enables reanalysis verification without API calls.
	VerifyCommitExists(ctx context.Context, sha string) (bool, error)
}

// VCSAPIClient retrieves repository metadata from VCS platform APIs.
type VCSAPIClient interface {
	// GetRepoID returns the platform-specific repository ID.
	// Returns ErrRepoNotFound if the repository does not exist.
	GetRepoID(ctx context.Context, host, owner, repo string, token *string) (string, error)
}
