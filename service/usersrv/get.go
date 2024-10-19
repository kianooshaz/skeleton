package usersrv

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/protocol"
)

func (m *Service) Get(ctx context.Context, id protocol.ID) (protocol.User, error) {
	user, err := m.queries.Get(ctx, m.pool, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userModel{}, protocol.ErrNotFound
		}

		return userModel{}, fmt.Errorf("get user with id %v: %w", id, err)
	}

	return userModel{
		id:        user.ID,
		createdAt: user.CreatedAt.Time,
	}, nil
}
