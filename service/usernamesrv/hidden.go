package usernamesrv

import (
	"context"

	"github.com/kianooshaz/skeleton/protocol"
) // Hidden implements protocol.ServiceUsername.
func (s *Service) Hidden(ctx context.Context, id protocol.ID) error {
	panic("unimplemented")
}

// Unhidden implements protocol.ServiceUsername.
func (s *Service) Unhidden(ctx context.Context, id protocol.ID) error {
	panic("unimplemented")
}
