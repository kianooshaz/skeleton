package usersrv

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/protocol"
	"github.com/kianooshaz/skeleton/protocol/derror"
)

func (m *Service) Get(ctx context.Context, id protocol.ID) (protocol.User, error) {
	user, err := m.queries.Get(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userModel{}, derror.NotFound
		}

		return userModel{}, fmt.Errorf("get user with id %v: %w", id, err)
	}

	return userModel{
		id:        user.ID,
		createdAt: user.CreatedAt.Time,
	}, nil
}
