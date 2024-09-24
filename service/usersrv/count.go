package usersrv

import (
	"context"
	"fmt"
)

func (m *Service) Count(ctx context.Context) (int64, error) {
	count, err := m.queries.Count(ctx, m.pool)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}
