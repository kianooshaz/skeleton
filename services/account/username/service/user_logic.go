package auns

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	fdp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	"github.com/kianooshaz/skeleton/foundation/stat"
	aunp "github.com/kianooshaz/skeleton/services/account/username/protocol"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

func (s *Service) Assign(ctx context.Context, req aunp.AssignRequest) (aunp.Username, error) {
	if len(req.Username) < int(s.config.MinLength) || len(req.Username) > int(s.config.MaxLength) {
		return aunp.Username{}, derror.ErrUsernameInvalid
	}

	if !s.isValidUsername(req.Username) {
		return aunp.Username{}, derror.ErrUsernameInvalid
	}

	countValue, err := s.storage.Count(ctx, req.Username)
	if err != nil {
		s.logger.Error("Error encountered while getting count by username", slog.String("error", err.Error()))

		return aunp.Username{}, derror.ErrInternalSystem
	}

	if countValue > 0 {
		return aunp.Username{}, derror.ErrUsernameCannotBeAssigned
	}

	shouldBePrimary, err := s.checkUserOrganizationMax(ctx, req.UserID, req.OrganizationID)
	if err != nil {
		return aunp.Username{}, err
	}

	status := stat.Unset
	if shouldBePrimary {
		status = stat.Primary
	}

	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("Error encountered while generating username id", slog.String("error", err.Error()))
		return aunp.Username{}, derror.ErrInternalSystem
	}

	username := aunp.Username{
		ID:             id,
		Username:       req.Username,
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		Status:         status,
	}

	err = s.storage.Create(ctx, username)
	if err != nil {
		s.logger.Error("Error encountered while creating username in database", slog.String("error", err.Error()))

		return aunp.Username{}, derror.ErrInternalSystem
	}

	return username, nil
}

func (s *Service) isValidUsername(value iunp.Username) bool {
	for _, char := range value {
		if !strings.ContainsRune(s.config.AllowCharacters, char) {
			return false
		}

		if char == ' ' || char == '-' || char == '_' {
			return false
		}
	}

	return true
}

func (s *Service) checkUserOrganizationMax(ctx context.Context, userID iup.UserID, organizationID iop.OrganizationID) (bool, error) {
	countByOrganization, err := s.storage.CountByUserAndOrganization(ctx, userID, organizationID)
	if err != nil {
		s.logger.Error("Error encountered while getting count by user and organization", slog.String("error", err.Error()))
		return false, derror.ErrInternalSystem
	}
	if countByOrganization > int64(s.config.MaxUserUsernamePerOrganization) {
		return false, derror.ErrUsernameMaxPerOrganization
	}

	shouldBePrimary := countByOrganization == 0

	return shouldBePrimary, nil
}

func (s *Service) Unassigned(ctx context.Context, id iunp.Username) error {
	username, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while getting username from database", "username", id, slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	if username.Status.Has(stat.Locked) {
		s.logger.Error("username is locked and cannot be unassigned", "username", id)
		return derror.ErrUsernameLocked
	}

	err = s.storage.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Error encountered while unassigning username", "id", id, slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) ListAssigned(ctx context.Context, req aunp.ListAssignedRequest) (aunp.ListAssignedResponse, error) {
	usernames, err := s.storage.ListByUserAndOrganization(ctx, req)
	if err != nil {
		s.logger.Error("Error encountered while listing assigned usernames", "userID", req.UserID, "organizationID", req.OrganizationID, slog.String("error", err.Error()))
		return aunp.ListAssignedResponse{}, derror.ErrInternalSystem
	}

	count, err := s.storage.CountByUserAndOrganization(ctx, req.UserID, req.OrganizationID)
	if err != nil {
		s.logger.Error("Error encountered while counting assigned usernames", "userID", req.UserID, "organizationID", req.OrganizationID, slog.String("error", err.Error()))
		return aunp.ListAssignedResponse{}, derror.ErrInternalSystem
	}

	result := make([]aunp.ListUsername, 0, len(usernames))
	for _, username := range usernames {
		if username.Status.Has(stat.Blocked) {
			count--
			continue
		}

		result = append(result, aunp.ListUsername{
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

	return aunp.ListAssignedResponse(pagination.NewResponse(req.Page, int(count), result)), nil
}

func (s *Service) BePrimary(ctx context.Context, id iunp.Username) error {
	shouldBePrimary, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, fdp.ErrRowNotFound) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while fetching username from db", slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	usernames, err := s.storage.ListByUserAndOrganization(ctx, aunp.ListAssignedRequest{
		UserID:         shouldBePrimary.UserID,
		OrganizationID: shouldBePrimary.OrganizationID,
		Page: pagination.Page{
			PageNumber: s.config.MaxUserUsernamePerOrganization, // We only need to fetch the usernames to update their status
		},
	})
	if err != nil {
		s.logger.Error("Error encountered while fetching usernames from db", slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	tx, ctx, err := session.BeginTransaction(ctx, s.storageConn)
	if err != nil {
		s.logger.Error("error")
	}

	defer func() {
		// TODO pgx.ErrTxClosed should change not depend database
		if err := tx.Rollback(); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			s.logger.Error("Error encountered while rolling back transaction", slog.String("error", err.Error()))
		}
	}()

	for _, username := range usernames {
		switch {
		case username.ID == shouldBePrimary.ID:
			if username.Status.Has(Primary) {
				return nil
			}

			username.Status.Add(Primary)

			if err := s.storage.UpdateStatus(ctx, username); err != nil {
				s.logger.Error("Error encountered while adding primary status to username", slog.String("error", err.Error()))

				return derror.ErrInternalSystem
			}

		case username.Status.Has(Primary):

			username.Status.Remove(Primary)

			if err := s.storage.UpdateStatus(ctx, username); err != nil {
				s.logger.Error("Error encountered while adding deleting status to username", slog.String("error", err.Error()))

				return derror.ErrInternalSystem
			}

		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("Error encountered while committing transaction", slog.String("error", err.Error()))

		return derror.ErrInternalSystem
	}

	return nil
}
