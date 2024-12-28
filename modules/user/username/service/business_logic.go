package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/stat"
	"github.com/kianooshaz/skeleton/modules/user/username/protocol"
	"github.com/kianooshaz/skeleton/modules/user/username/service/stores/db"
)

func (s *Service) Add(ctx context.Context, userID, OrganizationID uuid.UUID, id string) (protocol.Username, error) {
	if len(id) < int(s.config.MinLength) || len(id) > int(s.config.MaxLength) {
		return &Username{}, derror.ErrUsernameLength
	}

	if !s.isValidUsername(id) {
		return &Username{}, derror.ErrUsernameInvalidCharacters
	}

	countValue, err := s.db.Count(ctx, id)
	if err != nil {
		s.logger.Error("error at getting count by username", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if countValue > 0 {
		return &Username{}, derror.ErrUsernameAlreadyExists
	}

	countByUser, err := s.db.CountByUser(ctx, userID)
	if err != nil {
		s.logger.Error("error at getting count by user id", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if countByUser > int64(s.config.MaxPerUser) {
		return &Username{}, derror.ErrUsernameMaxPerUser
	}

	countByOrganization, err := s.db.CountByUserAndOrganization(ctx, userID)
	if err != nil {
		s.logger.Error("error at getting count by user and organization", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if countByOrganization > int64(s.config.MaxPerOrganization) {
		return &Username{}, derror.ErrUsernameMaxPerOrganization
	}

	status := stat.Unset
	if countByOrganization == 0 {
		status = stat.Primary
	}

	row, err := s.db.Create(ctx, db.CreateParams{
		ID:     id,
		UserID: userID,
		Status: int64(status),
	})
	if err != nil {
		s.logger.Error("error at creating username in database", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	return &Username{
		ID:             row.ID,
		UserID:         row.UserID,
		OrganizationID: row.OrganizationID,
		Status:         stat.Status(row.Status),
	}, nil
}

func (s *Service) Get(ctx context.Context, id string) (protocol.Username, error) {
	row, err := s.db.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &Username{}, derror.ErrUserNotFound
		}
		s.logger.Error("error at getting username from database", "username", id, "error", err)

		return Username{}, derror.ErrInternalSystem
	}

	return &Username{
		ID:             row.ID,
		UserID:         row.UserID,
		OrganizationID: row.OrganizationID,
		Status:         stat.Status(row.Status),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}, nil
}

// Search implements protocol.ServiceUsername.
func (s *Service) Search(ctx context.Context, userID *string, organizationID *string, status *stat.Status, limit int, offset int) ([]protocol.Username, error) {
	panic("unimplemented")
}

func (s *Service) Count(ctx context.Context, userID *string, organizationID *string, status *stat.Status) (int64, error) {
	panic("unimplemented")
}

func (s *Service) BePrimary(ctx context.Context, userID, organizationID uuid.UUID, id string) error {
	rows, err := s.db.ListByUserAndOrganization(ctx, db.ListByUserAndOrganizationParams{
		UserID:         userID,
		OrganizationID: organizationID,
	})
	if err != nil {
		s.logger.Error("error at listing usernames", "error", err)

		return derror.ErrInternalSystem
	}

	if len(rows) == 0 {
		return derror.ErrUsernameNotFound
	}

	tx, err := s._pdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		s.logger.Error("error at beginning transaction", "error", err)

		return derror.ErrInternalSystem
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("error at rolling back transaction", "error", err)
		}
	}()

	queries := db.New(tx)

	for _, row := range rows {
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
				s.logger.Error("error at updating username", "error", err)

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
				s.logger.Error("error at updating username", "error", err)

				return derror.ErrInternalSystem
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("error at committing transaction", "error", err)

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
		s.logger.Error("error at getting username from database", "username", id, "error", err)

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
		s.logger.Error("error at updating username", "error", err)

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
		s.logger.Error("error at getting username from database", "username", id, "error", err)

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
		s.logger.Error("error at updating username", "error", err)

		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) Reserve(ctx context.Context, id string) (protocol.Username, error) {
	if len(id) < int(s.config.MinLength) || len(id) > int(s.config.MaxLength) {
		return &Username{}, derror.ErrUsernameLength
	}

	if !s.isValidUsername(id) {
		return &Username{}, derror.ErrUsernameInvalidCharacters
	}

	count, err := s.db.Count(ctx, id)
	if err != nil {
		s.logger.Error("error at getting username from database", "username", id, "error", err)

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
		s.logger.Error("error at creating username in database", "error", err)

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
		s.logger.Error("error at getting username from database", "username", id, "error", err)

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
		s.logger.Error("error at updating username", "error", err)

		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) isValidUsername(value string) bool {
	for _, char := range value {
		if !strings.ContainsRune(s.config.AllowCharacters, char) {
			return false
		}
	}

	return true
}
