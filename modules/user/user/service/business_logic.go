package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/modules/user/user/protocol"
)

func (s *UserService) Create(ctx context.Context) (protocol.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("error at creating new uuid", "error", err)

		return nil, derror.ErrInternalSystem
	}

	user, err := s.db.Create(ctx, id)
	if err != nil {
		s.logger.Error("error at creating user in database", "error", err)

		return nil, derror.ErrInternalSystem
	}

	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *UserService) Get(ctx context.Context, id uuid.UUID) (protocol.User, error) {
	user, err := s.db.Get(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &User{}, derror.ErrUserNotFound
		}
		s.logger.Error("error at getting user from database", "id", id, "error", err)

		return nil, derror.ErrInternalSystem
	}

	return &User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (s *UserService) List(ctx context.Context, orderBy order.By) ([]protocol.User, error) {
	users, err := s.db.List(ctx, orderBy.PGX())
	if err != nil {
		s.logger.Error("error at listing users from database", "error", err)

		return nil, derror.ErrInternalSystem
	}

	list := make([]protocol.User, 0, len(users))
	for _, u := range users {
		list = append(list, &User{
			ID:        u.ID,
			CreatedAt: u.CreatedAt.Time,
		})
	}

	return list, nil
}

func (s *UserService) Count(ctx context.Context) (int64, error) {
	count, err := s.db.Count(ctx)
	if err != nil {
		s.logger.Error("error at counting users in database", "error", err)

		return 0, derror.ErrInternalSystem
	}

	return count, nil
}
