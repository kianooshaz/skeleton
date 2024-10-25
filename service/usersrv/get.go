package usersrv

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (m *Service) Get(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := m.queries.Get(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrNotFound
		}

		return User{}, fmt.Errorf("get user with id %v: %w", id, err)
	}

	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}
