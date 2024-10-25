package usernamesrv

import (
	"context"

	"github.com/google/uuid"
) // Update implements protocol.ServiceUsername.
func (s *Service) Update(ctx context.Context, id uuid.UUID, value string) error {
	panic("unimplemented")
}
