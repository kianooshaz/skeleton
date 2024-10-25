package usersrv

import (
	"context"
	"fmt"

	"github.com/kianooshaz/skeleton/foundation/order"
)

func (m *Service) List(ctx context.Context, orderBy order.By) ([]User, error) {
	users, err := m.queries.List(ctx, orderBy.PGX())
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
