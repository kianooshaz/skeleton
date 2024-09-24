package usersrv

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/protocol"
)

func (m *Service) New(ctx context.Context) (protocol.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return userModel{}, fmt.Errorf("new uuid: %w", err)
	}

	user, err := m.queries.Create(ctx, m.pool, id)
	if err != nil {
		return userModel{}, err
	}

	return userModel{
		id:        user.ID,
		createdAt: user.CreatedAt.Time,
	}, nil
}
