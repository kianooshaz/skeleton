package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/types"
	up "github.com/kianooshaz/skeleton/services/user/user/protocol"
)

func (s *service) Create(ctx context.Context) (up.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.ErrorContext(ctx, "error at generating user id", slog.Any("error", err))

		return up.User{}, derror.ErrInternalSystem
	}

	user := up.User{
		ID:        types.UserID(id),
		CreatedAt: time.Now(),
	}

	if err = s.storage.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "error at creating user in database", slog.Any("error", err))

		return up.User{}, derror.ErrInternalSystem
	}

	return user, nil
}

func (s *service) Get(ctx context.Context, req up.GetUserRequest) (up.User, error) {
	user, err := s.storage.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, derror.ErrUserNotFound) {
			return up.User{}, err
		}

		s.logger.ErrorContext(ctx, "error at getting user from database", slog.Any("error", err))

		return up.User{}, derror.ErrInternalSystem
	}

	return user, nil

}

func (s *service) List(ctx context.Context, req up.ListUserRequest) (pagination.Response[up.User], error) {
	users, err := s.storage.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		s.logger.ErrorContext(ctx, "error at listing users from database", slog.Any("error", err))

		return pagination.Response[up.User]{}, derror.ErrInternalSystem
	}

	totalCount, err := s.storage.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "error at counting users from database", slog.Any("error", err))

		return pagination.Response[up.User]{}, derror.ErrInternalSystem
	}

	return pagination.NewResponse(req.Page, totalCount, users), nil
}
