package uus

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	uup "github.com/kianooshaz/skeleton/services/user/user/protocol"
)

func (s *service) Create(ctx context.Context) (uup.CreateResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while generating user id", slog.String("error", err.Error()))

		return uup.CreateResponse{}, derror.ErrInternalSystem
	}

	user := uup.User{
		ID:        iup.UserID(id),
		CreatedAt: time.Now(),
	}

	if err = s.storage.Create(ctx, user); err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while creating user in storage", slog.String("error", err.Error()))

		return uup.CreateResponse{}, derror.ErrInternalSystem
	}

	return uup.CreateResponse{
		Data: user,
	}, nil
}

func (s *service) Get(ctx context.Context, req uup.GetRequest) (uup.GetResponse, error) {
	user, err := s.storage.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, derror.ErrUserNotFound) {
			return uup.GetResponse{}, err
		}

		s.logger.ErrorContext(ctx, "Error encountered while getting user from storage", slog.String("error", err.Error()), slog.Any("req", req))

		return uup.GetResponse{}, derror.ErrInternalSystem
	}

	return uup.GetResponse{Data: user}, nil

}

func (s *service) List(ctx context.Context, req uup.ListRequest) (uup.ListResponse, error) {
	users, err := s.storage.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while listing users from storage", slog.String("error", err.Error()), slog.Any("req", req))

		return uup.ListResponse{}, derror.ErrInternalSystem
	}

	totalCount, err := s.storage.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(ctx, "Error encountered while counting users from storage", slog.String("error", err.Error()), slog.Any("req", req))

		return uup.ListResponse{}, derror.ErrInternalSystem
	}

	return uup.ListResponse(pagination.NewResponse(req.Page, totalCount, users)), nil
}
