package usernamesrv

import "errors"

var (
	ErrNotFound       = errors.New("username not found")
	ErrInvalidRequest = errors.New("invalid request")
	ErrDuplicate      = errors.New("duplicate")
)
