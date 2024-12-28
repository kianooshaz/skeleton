package service

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


func (s *PasswordService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.Cost)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func (s *PasswordService) verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *PasswordService) hasPasswordBeenUsedBefore(ctx context.Context, userID uuid.UUID, passwordHash string) (bool, error) {
	passwords, err := s.db.History(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get password history", "error", err)

		return false, err
	}

	for _, password := range passwords {
		if password.PasswordHash == passwordHash {
			return true, nil
		}
	}

	return false, nil
}
