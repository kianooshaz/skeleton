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
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
)

func (s *Service) Assign(ctx context.Context, req usernameproto.AssignRequest) (usernameproto.Username, error) {
	if len(req.Username) < int(s.config.MinLength) || len(req.Username) > int(s.config.MaxLength) {
		return usernameproto.Username{}, derror.ErrUsernameInvalid
	}

	if !s.isValidUsername(req.Username) {
		return usernameproto.Username{}, derror.ErrUsernameInvalid
	}

	exist, err := s.storage.Exist(ctx, req.Username)
	if err != nil {
		s.logger.Error("Error encountered while getting count by username", slog.String("error", err.Error()))

		return usernameproto.Username{}, derror.ErrInternalSystem
	}

	if exist {
		return usernameproto.Username{}, derror.ErrUsernameCannotBeAssigned
	}

	shouldBePrimary, err := s.checkAccountMax(ctx, req.AccountID)
	if err != nil {
		return usernameproto.Username{}, err
	}

	status := stat.Unset
	if shouldBePrimary {
		status = stat.Primary
	}

	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("Error encountered while generating username id", slog.String("error", err.Error()))
		return usernameproto.Username{}, derror.ErrInternalSystem
	}

	username := usernameproto.Username{
		ID:        id,
		Username:  req.Username,
		AccountID: req.AccountID,
		Status:    status,
	}

	err = s.storage.Create(ctx, username)
	if err != nil {
		s.logger.Error("Error encountered while creating username in database", slog.String("error", err.Error()))

		return usernameproto.Username{}, derror.ErrInternalSystem
	}

	return username, nil
}

func (s *Service) isValidUsername(value string) bool {
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

func (s *Service) checkAccountMax(ctx context.Context, accountID accprotocol.AccountID) (bool, error) {
	countByAccount, err := s.storage.CountByAccount(ctx, accountID)
	if err != nil {
		s.logger.Error("Error encountered while getting count by account", slog.String("error", err.Error()))
		return false, derror.ErrInternalSystem
	}
	if countByAccount > int64(s.config.MaxUserUsernamePerOrganization) {
		return false, derror.ErrUsernameMaxPerOrganization
	}

	shouldBePrimary := countByAccount == 0

	return shouldBePrimary, nil
}

func (s *Service) Unassigned(ctx context.Context, id uuid.UUID) error {
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

func (s *Service) ListAssigned(ctx context.Context, req usernameproto.ListAssignedRequest) (usernameproto.ListAssignedResponse, error) {
	usernames, err := s.storage.ListByUserAndOrganization(ctx, req)
	if err != nil {
		s.logger.Error("Error encountered while listing assigned usernames", "accountID", req.AccountID, slog.String("error", err.Error()))
		return usernameproto.ListAssignedResponse{}, derror.ErrInternalSystem
	}

	count, err := s.storage.CountByAccount(ctx, req.AccountID)
	if err != nil {
		s.logger.Error("Error encountered while counting assigned usernames", "accountID", req.AccountID, slog.String("error", err.Error()))
		return usernameproto.ListAssignedResponse{}, derror.ErrInternalSystem
	}

	result := make([]usernameproto.ListUsername, 0, len(usernames))
	for _, username := range usernames {
		if username.Status.Has(stat.Blocked) {
			count--
			continue
		}

		result = append(result, usernameproto.ListUsername{
			ID:        username.ID,
			Username:  username.Username,
			AccountID: username.AccountID,
			Primary:   username.Status.Has(stat.Primary),
			Locked:    username.Status.Has(stat.Locked),
			Blocked:   username.Status.Has(stat.Blocked),
			Reserved:  username.Status.Has(stat.Reserved),
		})
	}

	return usernameproto.ListAssignedResponse(pagination.NewResponse(req.Page, int(count), result)), nil
}

func (s *Service) BePrimary(ctx context.Context, id uuid.UUID) error {
	shouldBePrimary, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, fdp.ErrRowNotFound) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while fetching username from db", slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	usernames, err := s.storage.ListByUserAndOrganization(ctx, usernameproto.ListAssignedRequest{
		AccountID: shouldBePrimary.AccountID,
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
