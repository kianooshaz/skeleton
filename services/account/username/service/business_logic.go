package usernameservice

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	dbproto "github.com/kianooshaz/skeleton/foundation/database/proto"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/stat"
	aunp "github.com/kianooshaz/skeleton/services/account/username/proto"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (aunp.Username, error) {
	username, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, dbproto.ErrRowNotFound) {
			return aunp.Username{}, derror.ErrUsernameNotFound
		}

		s.logger.Error("Error encountered while fetching username from database", slog.String("error", err.Error()))
		return aunp.Username{}, derror.ErrInternalSystem
	}

	return username, nil
}

func (s *Service) List(ctx context.Context, req aunp.ListRequest) (aunp.ListResponse, error) {
	usernames, err := s.storage.ListWithSearch(ctx, req)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while searching usernames",
			slog.String("error", err.Error()),
			slog.Any("request", req),
		)

		return aunp.ListResponse{}, derror.ErrInternalSystem
	}

	count, err := s.storage.CountWithSearch(ctx, req)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while counting usernames",
			slog.String("error", err.Error()),
			slog.Any("request", req),
		)

		return aunp.ListResponse{}, derror.ErrInternalSystem
	}

	result := make([]aunp.ListUsername, 0, len(usernames))
	for _, username := range usernames {
		result = append(result, aunp.ListUsername{
			ID:        username.ID,
			Username:  username.Username,
			AccountID: username.AccountID,
			Primary:   username.Status.Has(stat.Primary),
			Locked:    username.Status.Has(stat.Locked),
			Blocked:   username.Status.Has(stat.Blocked),
			Reserved:  username.Status.Has(stat.Reserved),
		})
	}

	return aunp.ListResponse(pagination.NewResponse(req.Page, int(count), result)), nil
}
