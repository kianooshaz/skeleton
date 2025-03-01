package service

import (
	"context"
	"errors"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/services/authentication/password/protocol"
	"github.com/kianooshaz/skeleton/services/authentication/password/service/stores/db"
	"golang.org/x/crypto/bcrypt"
)

// SavePassword
func (s *PasswordService) SavePassword(ctx context.Context, userID uuid.UUID, password string) (protocol.Password, error) {
	if !s.evaluatePasswordStrength(password) {
		return &Password{}, derror.ErrPasswordIsWeak
	}

	passwordHash, err := s.hashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)

		return &Password{}, err
	}

	ok, err := s.checkPasswordHistory(ctx, userID, passwordHash)
	if err != nil {
		return &Password{}, err
	}

	if ok {
		return &Password{}, derror.ErrPasswordIsInHistory
	}

	row, err := s.db.GetByUserID(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {

	}

	if row.PasswordHash != "" {
		if err := s.db.Delete(ctx, row.ID); err != nil {
			s.logger.Error("failed to delete old password", "error", err)

			return &Password{}, err
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("error at creating new uuid", "error", err)

		return &Password{}, derror.ErrInternalSystem
	}

	row, err = s.db.Create(ctx, db.CreateParams{
		ID:           id,
		UserID:       userID,
		PasswordHash: string(passwordHash),
	})
	if err != nil {
		s.logger.Error("failed to save password", "error", err)

		return &Password{}, err
	}

	return &Password{
		ID:           row.ID,
		UserID:       row.UserID,
		PasswordHash: hash(row.PasswordHash),
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (s *PasswordService) hashPassword(password string) (hash, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Cost)
	return hash(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func (s *PasswordService) verifyPassword(storedPassword hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	return err == nil
}

func (s *PasswordService) checkPasswordHistory(ctx context.Context, userID uuid.UUID, password hash) (bool, error) {
	rows, err := s.db.History(ctx, db.HistoryParams{
		UserID: userID,
		Limit:  s.config.CheckPasswordHistoryLimit,
	})
	if err != nil {
		s.logger.Error("failed to get password history", "error", err)

		return false, derror.ErrInternalSystem
	}

	for _, row := range rows {
		if hash(row.PasswordHash) == password {
			return true, nil
		}
	}

	return false, nil
}

func (s *PasswordService) evaluatePasswordStrength(password string) bool {
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

// GetPasswordGuidelines
func (s *PasswordService) GetPasswordGuidelines() string {
	return ""
}

func (s *PasswordService) isPasswordCommon(password string) bool {
	return s.commonPasswords[password]
}
