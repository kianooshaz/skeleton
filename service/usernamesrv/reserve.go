package usernamesrv

import (
	"context"
)

// Reserve implements protocol.ServiceUsername.
func (s *Service) Reserve(ctx context.Context, value string) error {
	panic("unimplemented")
}

// Unreserve implements protocol.ServiceUsername.
func (s *Service) Unreserve(ctx context.Context, value string) error {
	panic("unimplemented")
}
