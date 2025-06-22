package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	fdp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	"github.com/kianooshaz/skeleton/foundation/stat"
	"github.com/kianooshaz/skeleton/protocol"
	aup "github.com/kianooshaz/skeleton/services/account/username/protocol"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (aup.Username, error) {
	username, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, fdp.ErrRowNotFound) {
			return aup.Username{}, derror.ErrUsernameNotFound
		}

		s.logger.Error("Error encountered while fetching username from database", slog.String("error", err.Error()))
		return aup.Username{}, derror.ErrInternalSystem
	}

	return username, nil
}

func (s *Service) List(ctx context.Context, req aup.ListRequest) (pagination.Response[aup.ListUsername], error) {
	usernames, err := s.storage.ListWithSearch(ctx, req)
	if err != nil {
		s.logger.Error("Error encountered while searching usernames", slog.String("error", err.Error()), slog.Any("request", req))

		return pagination.Response[aup.ListUsername]{}, derror.ErrInternalSystem
	}

	count, err := s.storage.CountWithSearch(ctx, req)
	if err != nil {
		s.logger.Error("Error encountered while counting usernames", slog.String("error", err.Error()), slog.Any("request", req))

		return pagination.Response[aup.ListUsername]{}, derror.ErrInternalSystem
	}

	result := make([]aup.ListUsername, 0, len(usernames))
	for _, username := range usernames {
		result = append(result, aup.ListUsername{
			ID:             username.ID,
			Username:       username.Username,
			UserID:         username.UserID,
			OrganizationID: username.OrganizationID,
			Primary:        username.Status.Has(stat.Primary),
			Locked:         username.Status.Has(stat.Locked),
			Blocked:        username.Status.Has(stat.Blocked),
			Reserved:       username.Status.Has(stat.Reserved),
		})
	}

	return pagination.NewResponse(req.Page, int(count), result), nil
}

func (s *Service) BePrimary(ctx context.Context, id uuid.UUID) error {
	username, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, fdp.ErrRowNotFound) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while fetching username from db", "error", err)
		return derror.ErrInternalSystem
	}

	tx, ctx, err := session.BeginTransaction(ctx, s.storageConn)
	if err != nil {
		s.logger.Error("error")
	}

	defer func() {
		// TODO pgx.ErrTxClosed should change not depend database
		if err := tx.Rollback(); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("Error encountered while rolling back transaction", "error", err)
		}
	}()

	queries := db.New(tx)

	for _, row := range username {
		status := stat.Status(row.Status)

		switch {
		case row.ID == id:
			if status.Has(stat.Primary) {
				return nil
			}

			status.Add(stat.Primary)

			if err := queries.Update(ctx, db.UpdateParams{
				ID:             id,
				UserID:         row.UserID,
				OrganizationID: row.OrganizationID,
				Status:         int64(status),
			}); err != nil {
				s.logger.Error("Error encountered while updating username", "error", err)

				return derror.ErrInternalSystem
			}

		case status.Has(stat.Primary):
			status.Remove(stat.Primary)

			if err := queries.Update(ctx, db.UpdateParams{
				ID:             row.ID,
				UserID:         row.UserID,
				OrganizationID: row.OrganizationID,
				Status:         int64(status),
			}); err != nil {
				s.logger.Error("Error encountered while updating username", "error", err)

				return derror.ErrInternalSystem
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("Error encountered while committing transaction", "error", err)

		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) Hidden(ctx context.Context, id string) error {
	row, err := s.db.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while getting username from database", "username", id, "error", err)

		return derror.ErrInternalSystem
	}

	status := stat.Status(row.Status)

	if status.Has(stat.Hidden) {
		return nil
	}

	status.Add(stat.Hidden)

	if err := s.db.Update(ctx, db.UpdateParams{
		ID:             id,
		UserID:         row.UserID,
		OrganizationID: row.OrganizationID,
		Status:         int64(status),
	}); err != nil {
		s.logger.Error("Error encountered while updating username", "error", err)

		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) Unhidden(ctx context.Context, id string) error {
	row, err := s.db.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while getting username from database", "username", id, "error", err)

		return derror.ErrInternalSystem
	}

	status := stat.Status(row.Status)

	if !status.Has(stat.Hidden) {
		return nil
	}

	status.Remove(stat.Hidden)

	if err := s.db.Update(ctx, db.UpdateParams{
		ID:             id,
		UserID:         row.UserID,
		OrganizationID: row.OrganizationID,
		Status:         int64(status),
	}); err != nil {
		s.logger.Error("Error encountered while updating username", "error", err)

		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) Reserve(ctx context.Context, id string) (protocol.Username, error) {
	if len(id) < int(s.config.MinLength) || len(id) > int(s.config.MaxLength) {
		return &Username{}, derror.ErrUsernameInvalid
	}

	if !s.isValidUsername(id) {
		return &Username{}, derror.ErrUsernameInvalid
	}

	count, err := s.db.Count(ctx, id)
	if err != nil {
		s.logger.Error("Error encountered while getting username from database", "username", id, "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if count > 0 {
		return &Username{}, derror.ErrUsernameAlreadyExists
	}

	newRow, err := s.db.Create(ctx, db.CreateParams{
		ID:     id,
		Status: int64(stat.Reserved),
	})
	if err != nil {
		s.logger.Error("Error encountered while creating username in database", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	return &Username{
		ID:             newRow.ID,
		UserID:         newRow.UserID,
		OrganizationID: newRow.OrganizationID,
		Status:         stat.Status(newRow.Status),
	}, nil
}

func (s *Service) Unreserve(ctx context.Context, id string, userID, organizationID uuid.UUID) error {
	row, err := s.db.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while getting username from database", "username", id, "error", err)

		return derror.ErrInternalSystem
	}

	if stat.Status(row.Status) != stat.Reserved {
		return derror.ErrUsernameNotReserved
	}

	if err := s.db.Update(ctx, db.UpdateParams{
		ID:             id,
		UserID:         userID,
		OrganizationID: organizationID,
		Status:         int64(stat.Unset),
	}); err != nil {
		s.logger.Error("Error encountered while updating username", "error", err)

		return derror.ErrInternalSystem
	}

	return nil
}
