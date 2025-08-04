package userservice

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

func (s *Service) Create(ctx context.Context) (userproto.CreateResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while generating user id", slog.String("error", err.Error()))

		return userproto.CreateResponse{}, derror.ErrInternalSystem
	}

	user := userproto.User{
		ID:        userproto.UserID(id),
		CreatedAt: time.Now(),
	}

	if err = s.persister.Create(ctx, user); err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while creating user in storage",
			slog.String("error", err.Error()),
			slog.Any("user", user),
		)

		return userproto.CreateResponse{}, derror.ErrInternalSystem
	}

	return userproto.CreateResponse{
		Data: user,
	}, nil
}

func (s *Service) Get(ctx context.Context, req userproto.GetRequest) (userproto.GetResponse, error) {
	user, err := s.persister.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, derror.ErrUserNotFound) {
			return userproto.GetResponse{}, err
		}

		s.logger.ErrorContext(
			ctx,
			"Error encountered while getting user from storage",
			slog.String("error", err.Error()),
			slog.Any("req", req),
		)

		return userproto.GetResponse{}, derror.ErrInternalSystem
	}

	return userproto.GetResponse{Data: user}, nil
}

func (s *Service) List(ctx context.Context, req userproto.ListRequest) (userproto.ListResponse, error) {
	users, err := s.persister.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while listing users from storage",
			slog.String("error", err.Error()),
			slog.Any("req", req),
		)

		return userproto.ListResponse{}, derror.ErrInternalSystem
	}

	totalCount, err := s.persister.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while counting users from storage",
			slog.String("error", err.Error()),
			slog.Any("req", req),
		)

		return userproto.ListResponse{}, derror.ErrInternalSystem
	}

	return userproto.ListResponse(pagination.NewResponse(req.Page, totalCount, users)), nil
}
