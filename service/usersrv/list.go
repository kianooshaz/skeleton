package usersrv

import (
	"context"
	"fmt"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/protocol"
)

func (m *Service) List(ctx context.Context, orderBy order.By) ([]protocol.User, error) {
	users, err := m.queries.List(ctx, orderBy.PGX())
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}

	list := make([]protocol.User, 0, len(users))
	for _, u := range users {
		list = append(list, userModel{
			id:        u.ID,
			createdAt: u.CreatedAt.Time,
		})
	}

	return list, nil
}
