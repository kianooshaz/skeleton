package usernamesrv

import (
	"context"

	"github.com/kianooshaz/skeleton/protocol"
)

// BePrimary implements protocol.ServiceUsername.
func (s *Service) BePrimary(ctx context.Context, id protocol.ID) error {
	panic("unimplemented")
}
