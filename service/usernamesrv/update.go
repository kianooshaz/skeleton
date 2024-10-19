package usernamesrv

import (
	"context"

	"github.com/kianooshaz/skeleton/protocol"
) // Update implements protocol.ServiceUsername.
func (s *Service) Update(ctx context.Context, id protocol.ID, value string) error {
	panic("unimplemented")
}
