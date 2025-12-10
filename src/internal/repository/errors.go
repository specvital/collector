package repository

import "errors"

// Sentinel errors for repository operations.
var (
	ErrInvalidInput = errors.New("invalid input")
)
