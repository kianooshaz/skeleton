package usernamesrv

import (
	"context"

	"github.com/kianooshaz/skeleton/protocol"
)

// Count implements protocol.ServiceUsername.
func (s *Service) Count(ctx context.Context, userID protocol.ID) (int64, error) {
	panic("unimplemented")
}
