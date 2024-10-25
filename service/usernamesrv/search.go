package usernamesrv

import (
	"context"
)

// Search implements protocol.ServiceUsername.
func (s *Service) Search(ctx context.Context, value string) (Username, error) {
	panic("unimplemented")
}
