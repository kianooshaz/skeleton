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
	"github.com/kianooshaz/skeleton/modules/user/user/protocol"
)

func (s *userService) Create(ctx context.Context) (protocol.User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.ErrorContext(ctx, "error at generating user id", slog.Any("error", err))

		return protocol.User{}, derror.ErrInternalSystem
	}

	user := protocol.User{
		ID:        types.UserID(id),
		CreatedAt: time.Now(),
	}

	if err = s.storage.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "error at creating user in database", slog.Any("error", err))

		return protocol.User{}, derror.ErrInternalSystem
	}

	return user, nil
}

func (s *userService) Get(ctx context.Context, req protocol.GetUserRequest) (protocol.User, error) {
	user, err := s.storage.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, derror.ErrUserNotFound) {
			return protocol.User{}, err
		}

		s.logger.ErrorContext(ctx, "error at getting user from database", slog.Any("error", err))

		return protocol.User{}, derror.ErrInternalSystem
	}

	return user, nil

}

func (s *userService) List(ctx context.Context, req protocol.ListUserRequest) (pagination.Response[protocol.User], error) {
	users, err := s.storage.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		s.logger.ErrorContext(ctx, "error at listing users from database", slog.Any("error", err))

		return pagination.Response[protocol.User]{}, derror.ErrInternalSystem
	}

	totalCount, err := s.storage.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "error at counting users from database", slog.Any("error", err))

		return pagination.Response[protocol.User]{}, derror.ErrInternalSystem
	}

	return pagination.NewResponse(req.Page, totalCount, users), nil
}
