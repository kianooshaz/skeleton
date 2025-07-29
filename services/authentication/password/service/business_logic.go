package authpass

import (
	"context"
	"errors"
	"log/slog"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/derror"
	authpassp "github.com/kianooshaz/skeleton/services/authentication/password/protocol"
	"github.com/kianooshaz/skeleton/services/authentication/password/service/stores/db"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	"golang.org/x/crypto/bcrypt"
)


func  (s *service) Verify(ctx, password string) error {
	return nil
}

// SavePassword
func (s *service) Update(ctx, req authpassp.UpdateRequest) error {
	// TODO check otp

	if !s.evaluatePasswordStrength(req.NewPassword) {
		return derror.ErrPasswordIsWeak
	}

	passwordHash, err := s.hashPassword(req.NewPassword)
	if err != nil {
		s.logger.Error("failed to hash password", slog.String("error", err.Error()))

		return err
	}

	used, err := s.usedBefore(ctx, req.UserID, req.OrganizationID, passwordHash)
	if err != nil {
		return err
	}

	if used {
		return  derror.ErrPasswordUsedBefore
	}

	
		if err := s.db.Delete(ctx, row.ID); err != nil {
			s.logger.Error("failed to delete old password", slog.String("error", err.Error()))

			return &Password{}, err
		}

	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("Error encountered while creating new uuid", slog.String("error", err.Error()))

		return &Password{}, derror.ErrInternalSystem
	}

	row, err = s.db.Create(ctx, db.CreateParams{
		ID:           id,
		UserID:       userID,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		s.logger.Error("failed to save password", slog.String("error", err.Error()))

		return &Password{}, err
	}

	return &Password{
		ID:           row.ID,
		UserID:       row.UserID,
		PasswordHash: hash(row.PasswordHash),
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (s *service) hashPassword(password string) (hash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Cost)
	return hash(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func (s *service) verifyPassword(storedPassword hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	return err == nil
}

func (s *service) usedBefore(ctx context.Context, userID iup.UserID, organizationID iop.OrganizationID, hashed string) (bool, error) {

	rows, err := s.db.History(ctx, db.HistoryParams{
		UserID: userID,
		Limit:  s.config.CheckPasswordHistoryLimit,
	})
	if err != nil {
		s.logger.Error("failed to get password history", slog.String("error", err.Error()))

		return false, derror.ErrInternalSystem
	}

	for _, row := range rows {
		if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) {
			return true, nil
		}
	}

	return false, nil
}

func (s *service) evaluatePasswordStrength(password string) bool {
	if len(password) < 8 {
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

func (s *service) Guidelines() (authpassp.GuidelinesResponse, error){
	return authpassp.GuidelinesResponse{
     Data: authpassp.Guidelines{
		Required: s.config.RequiredGuidelines,
		BetterHave: s.config.BetterHave,
	 }
	}
}

func (s *service) isPasswordCommon(password string) bool {
	return s.commonPasswords[password]
}

