package authpass

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"unicode"

	"github.com/google/uuid"
	fdp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	accproto "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	"golang.org/x/crypto/bcrypt"
)

type hash []byte

func (s *Service) Verify(ctx context.Context, password string) error {
	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (passwordproto.Password, error) {
	password, err := s.storage.Get(ctx, id)
	if err != nil {
		if errors.Is(err, fdp.ErrRowNotFound) {
			return passwordproto.Password{}, derror.ErrPasswordNotFound
		}

		s.logger.Error("Error encountered while fetching password from database", slog.String("error", err.Error()))
		return passwordproto.Password{}, derror.ErrInternalSystem
	}

	return password, nil
}

func (s *Service) List(ctx context.Context, req passwordproto.ListRequest) (passwordproto.ListResponse, error) {
	passwords, err := s.storage.ListWithSearch(ctx, req)
	if err != nil {
		s.logger.Error(
			"Error encountered while searching passwords",
			slog.String("error", err.Error()),
			slog.Any("request", req),
		)

		return passwordproto.ListResponse{}, derror.ErrInternalSystem
	}

	count, err := s.storage.CountWithSearch(ctx, req)
	if err != nil {
		s.logger.Error(
			"Error encountered while counting passwords",
			slog.String("error", err.Error()),
			slog.Any("request", req),
		)

		return passwordproto.ListResponse{}, derror.ErrInternalSystem
	}

	return passwordproto.ListResponse(pagination.NewResponse(req.Page, int(count), passwords)), nil
}

// Update updates the password for a given account.
func (s *Service) Update(ctx context.Context, req passwordproto.UpdateRequest) error {
	// TODO check otp

	if !s.evaluatePasswordStrength(req.NewPassword) {
		return derror.ErrPasswordIsWeak
	}

	passwordHash, err := s.hashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("failed to hash password", slog.String("error", err.Error()))
		return err
	}

	used, err := s.usedBefore(ctx, req.AccountID, string(passwordHash))
	if err != nil {
		return err
	}

	if used {
		return derror.ErrPasswordUsedBefore
	}

	// Get existing password to delete it
	existingPassword, err := s.storage.GetByAccountID(ctx, req.AccountID)
	if err != nil && !errors.Is(err, fdp.ErrRowNotFound) {
		s.logger.Error("failed to get existing password", slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	// Delete existing password if it exists
	if err == nil {
		if err := s.storage.Delete(ctx, existingPassword.ID); err != nil {
			s.logger.Error("failed to delete old password", slog.String("error", err.Error()))
			return err
		}
	}

	// Create new password
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("Error encountered while creating new uuid", slog.String("error", err.Error()))
		return derror.ErrInternalSystem
	}

	newPassword := passwordproto.Password{
		ID:           id,
		AccountID:    req.AccountID,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.storage.Create(ctx, newPassword)
	if err != nil {
		s.logger.Error("failed to save password", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *Service) hashPassword(password string) (hash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Cost)
	return hash(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func (s *Service) verifyPassword(storedPassword hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	return err == nil
}

func (s *Service) usedBefore(ctx context.Context, accountID accproto.AccountID, hashedPassword string) (bool, error) {
	passwords, err := s.storage.History(ctx, accountID, s.config.CheckPasswordHistoryLimit)
	if err != nil {
		s.logger.Error("failed to get password history", slog.String("error", err.Error()))
		return false, derror.ErrInternalSystem
	}

	for _, password := range passwords {
		if err := bcrypt.CompareHashAndPassword([]byte(password.PasswordHash), []byte(hashedPassword)); err == nil {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) evaluatePasswordStrength(password string) bool {
	if len(password) < int(s.config.MinLength) {
		return false
	}

	var hasSpecial bool
	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return false
	}

	if s.isPasswordCommon(password) {
		return false
	}

	return true
}

func (s *Service) Guidelines() (passwordproto.GuidelinesResponse, error) {
	return passwordproto.GuidelinesResponse{
		Data: passwordproto.Guidelines{
			Required:   s.config.RequiredGuidelines,
			BetterHave: s.config.BetterHave,
		},
	}, nil
}

func (s *Service) isPasswordCommon(password string) bool {
	return s.commonPasswords[password]
}
