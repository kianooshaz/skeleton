package usersrv

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (m *Service) New(ctx context.Context) (User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return User{}, fmt.Errorf("new uuid: %w", err)
	}

	user, err := m.queries.Create(ctx, id)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
