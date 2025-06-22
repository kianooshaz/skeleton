package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/stat"
	aup "github.com/kianooshaz/skeleton/services/account/username/protocol"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

func (s *Service) Assign(ctx context.Context, req aup.AssignRequest) (aup.Username, error) {
	if len(req.Username) < int(s.config.MinLength) || len(req.Username) > int(s.config.MaxLength) {
		return aup.Username{}, derror.ErrUsernameInvalid
	}

	if !s.isValidUsername(req.Username) {
		return aup.Username{}, derror.ErrUsernameInvalid
	}

	countValue, err := s.storage.Count(ctx, req.Username)
	if err != nil {
		s.logger.Error("Error encountered while getting count by username", "error", err)

		return aup.Username{}, derror.ErrInternalSystem
	}

	if countValue > 0 {
		return aup.Username{}, derror.ErrUsernameCannotBeAssigned
	}

	err = s.checkUserMax(ctx, req.UserID)
	if err != nil {
		return aup.Username{}, err
	}

	shouldBePrimary, err := s.checkUserOrganizationMax(ctx, req.UserID, req.OrganizationID)
	if err != nil {
		return aup.Username{}, err
	}

	status := stat.Unset
	if shouldBePrimary {
		status = stat.Primary
	}

	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("Error encountered while generating username id", "error", err)
		return aup.Username{}, derror.ErrInternalSystem
	}

	username := aup.Username{
		ID:             id,
		Username:       req.Username,
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		Status:         status,
	}

	err = s.storage.Create(ctx, username)
	if err != nil {
		s.logger.Error("Error encountered while creating username in database", "error", err)

		return aup.Username{}, derror.ErrInternalSystem
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

// Extracted function for user max check
func (s *Service) checkUserMax(ctx context.Context, userID iup.UserID) error {
	countByUser, err := s.storage.CountByUser(ctx, userID)
	if err != nil {
		s.logger.Error("Error encountered while getting count by user id", "error", err)
		return derror.ErrInternalSystem
	}
	if countByUser > int64(s.config.MaxPerUser) {
		return derror.ErrUsernameMaxPerUser
	}
	return nil
}

func (s *Service) checkUserOrganizationMax(ctx context.Context, userID iup.UserID, organizationID iop.OrganizationID) (bool, error) {
	countByOrganization, err := s.storage.CountByUserAndOrganization(ctx, userID, organizationID)
	if err != nil {
		s.logger.Error("Error encountered while getting count by user and organization", "error", err)
		return false, derror.ErrInternalSystem
	}
	if countByOrganization > int64(s.config.MaxPerOrganization) {
		return false, derror.ErrUsernameMaxPerOrganization
	}

	shouldBePrimary := countByOrganization == 0

	return shouldBePrimary, nil
}

func (s *Service) Unassign(ctx context.Context, id uuid.UUID) error {
	username, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("Error encountered while getting username from database", "username", id, "error", err)
		return derror.ErrInternalSystem
	}

	if username.Status.Has(stat.Locked) {
		s.logger.Error("username is locked and cannot be unassigned", "username", id)
		return derror.ErrUsernameLocked
	}

	err = s.storage.Delete(ctx, id)
	if err != nil {
		s.logger.Error("Error encountered while unassigning username", "id", id, "error", err)
		return derror.ErrInternalSystem
	}

	return nil
}

func (s *Service) ListAssigned(ctx context.Context, userID iup.UserID, organizationID iop.OrganizationID) ([]aup.Username, error) {
	usernames, err := s.storage.ListByUserAndOrganization(ctx, userID, organizationID)
	if err != nil {
		s.logger.Error("Error encountered while listing assigned usernames", "userID", userID, "organizationID", organizationID, "error", err)
		return nil, derror.ErrInternalSystem
	}

	result := make([]aup.Username, 0, len(usernames))
	for _, username := range usernames {
		if username.Status.Has(stat.Blocked) {
			continue
		}

		result = append(result, username)
	}

	return result, nil
}
