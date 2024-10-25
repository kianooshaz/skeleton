package usernamesrv

import (
	"context"

	"github.com/google/uuid"
)

// Hidden implements protocol.ServiceUsername.
func (s *Service) Hidden(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// Unhidden implements protocol.ServiceUsername.
func (s *Service) Unhidden(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}
