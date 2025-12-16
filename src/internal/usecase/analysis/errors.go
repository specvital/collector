package analysis

import "errors"

var (
	// ErrCloneFailed indicates VCS clone operation failed.
	ErrCloneFailed = errors.New("clone failed")

	// ErrSaveFailed indicates repository save operation failed.
	ErrSaveFailed = errors.New("save failed")

	// ErrScanFailed indicates parser scan operation failed.
	ErrScanFailed = errors.New("scan failed")

	// ErrTokenLookupFailed indicates OAuth token lookup failed due to infrastructure error.
	// This is different from token not found (which triggers graceful degradation).
	ErrTokenLookupFailed = errors.New("token lookup failed")
)
