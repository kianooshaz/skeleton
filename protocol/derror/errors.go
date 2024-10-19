package derror

import "errors"

// Define a set of common error variables for reuse throughout the application.
// This helps in consistent error handling and comparison.

var (
	// ErrNotFound represents an error when a requested resource is not found.
	NotFound       = errors.New("not found error")
	Duplicate      = errors.New("duplicate error")
	InvalidRequest = errors.New("invalid request error")
)
