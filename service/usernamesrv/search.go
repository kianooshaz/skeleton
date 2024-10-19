package usernamesrv

import (
	"context"

	"github.com/kianooshaz/skeleton/protocol"
)

// Search implements protocol.ServiceUsername.
func (s *Service) Search(ctx context.Context, value string) (protocol.Username, error) {
	panic("unimplemented")
}
