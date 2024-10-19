package usernamesrv

import (
	"context"

	"github.com/kianooshaz/skeleton/protocol"
)

// List implements protocol.ServiceUsername.
func (s *Service) List(ctx context.Context, userID protocol.ID) []protocol.Username {
	panic("unimplemented")
}
