package usermod

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/order"
)

var (
	ErrNotFound = errors.New("user not found")
)

func (m *Module) New(ctx context.Context) (User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return User{}, fmt.Errorf("new uuid: %w", err)
	}

	user, err := m.queries.Create(ctx, m.pool, id)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (m *Module) Get(ctx context.Context, id uuid.UUID) (User, error) {
	user, err := m.queries.Get(ctx, m.pool, id)
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

func (m *Module) List(ctx context.Context, orderBy order.By) ([]User, error) {
	users, err := m.queries.List(ctx, m.pool, orderBy.PGX())
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	list := make([]User, 0, len(users))
	for _, u := range users {
		list = append(list, User{
			ID:        u.ID,
			CreatedAt: u.CreatedAt.Time,
		})
	}

	return list, nil
}

func (m *Module) Count(ctx context.Context) (int64, error) {
	count, err := m.queries.Count(ctx, m.pool)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}
