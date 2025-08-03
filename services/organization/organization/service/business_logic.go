package oos

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
)

func (s *service) Create(ctx context.Context) (orgproto.CreateResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while generating organization id",
			slog.String("error", err.Error()),
		)

		return orgproto.CreateResponse{}, derror.ErrInternalSystem
	}

	organization := orgproto.Organization{
		ID:        orgproto.OrganizationID(id),
		CreatedAt: time.Now(),
	}

	if err = s.persister.Create(ctx, organization); err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while creating organization in storage",
			slog.String("error", err.Error()),
			slog.Any("organization", organization),
		)

		return orgproto.CreateResponse{}, derror.ErrInternalSystem
	}

	return orgproto.CreateResponse{
		Data: organization,
	}, nil
}

func (s *service) Get(ctx context.Context, req orgproto.GetRequest) (orgproto.GetResponse, error) {
	organization, err := s.persister.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, derror.ErrOrganizationNotFound) {
			return orgproto.GetResponse{}, err
		}

		s.logger.ErrorContext(
			ctx,
			"Error encountered while getting organization from storage",
			slog.String("error", err.Error()),
			slog.Any("req", req),
		)

		return orgproto.GetResponse{}, derror.ErrInternalSystem
	}

	return orgproto.GetResponse{Data: organization}, nil
}

func (s *service) List(ctx context.Context, req orgproto.ListRequest) (orgproto.ListResponse, error) {
	organizations, err := s.persister.List(ctx, req.Page, req.OrderBy)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while listing organizations from storage",
			slog.String("error", err.Error()),
			slog.Any("req", req),
		)

		return orgproto.ListResponse{}, derror.ErrInternalSystem
	}

	totalCount, err := s.persister.Count(ctx)
	if err != nil {
		s.logger.ErrorContext(
			ctx,
			"Error encountered while counting organizations from storage",
			slog.String("error", err.Error()),
			slog.Any("req", req),
		)

		return orgproto.ListResponse{}, derror.ErrInternalSystem
	}

	return orgproto.ListResponse(pagination.NewResponse(req.Page, totalCount, organizations)), nil
}
