package suu

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	suup "github.com/kianooshaz/skeleton/services/user/user/protocol"
)

func (s *service) Create(ctx context.Context) (suup.CreateResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while generating user id", slog.String("error", err.Error()))

		return suup.CreateResponse{}, derror.ErrInternalSystem
	}

	user := suup.User{
		ID:        iup.UserID(id),
		CreatedAt: time.Now(),
	}

	if err = s.storage.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while creating user in storage", slog.String("error", err.Error()))

		return suup.CreateResponse{}, derror.ErrInternalSystem
	}

	return suup.CreateResponse{
		Data: user,
	}, nil
}

func (s *service) Get(ctx context.Context, req suup.GetRequest) (suup.GetResponse, error) {
	user, err := s.storage.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, derror.ErrUserNotFound) {
			return suup.GetResponse{}, err
		}

		s.logger.ErrorContext(ctx, "Error encountered while getting user from storage", slog.String("error", err.Error()), slog.Any("req", req))

		return suup.GetResponse{}, derror.ErrInternalSystem
	}

	return suup.GetResponse{Data: user}, nil

}

func (s *service) List(ctx context.Context, req suup.ListRequest) (suup.ListResponse, error) {
	users, err := s.storage.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while listing users from storage", slog.String("error", err.Error()), slog.Any("req", req))

		return suup.ListResponse{}, derror.ErrInternalSystem
	}

	totalCount, err := s.storage.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while counting users from storage", slog.String("error", err.Error()), slog.Any("req", req))

		return suup.ListResponse{}, derror.ErrInternalSystem
	}

	return suup.ListResponse(pagination.NewResponse(req.Page, totalCount, users)), nil
}
