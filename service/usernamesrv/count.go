package usernamesrv

import (
	"context"

	"github.com/google/uuid"
)

// Count implements protocol.ServiceUsername.
func (s *Service) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	panic("unimplemented")
}
