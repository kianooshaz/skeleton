package usernamesrv

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/status"
	"github.com/kianooshaz/skeleton/service/usernamesrv/stores/usernamedb"
)

func (s *Service) Add(ctx context.Context, userID uuid.UUID, value string) (Username, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Username{}, fmt.Errorf("new uuid: %w", err)
	}

	if len(value) < int(s.config.MinLength) || len(value) > int(s.config.MaxLength) {
		return Username{}, ErrInvalidRequest
	}

	countValue, err := s.queries.CountByUsername(ctx, value)
	if err != nil {
		return Username{}, fmt.Errorf("get count by username: %w", err)
	}

	if countValue > 0 {
		return Username{}, ErrDuplicate
	}

	countForUser, err := s.queries.CountByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return Username{}, fmt.Errorf("get count by user id: %w", err)
	}

	if countForUser > int64(s.config.MaxPerUser) {
		return Username{}, ErrInvalidRequest
	}

	username, err := s.queries.Create(ctx, usernamedb.CreateParams{
		ID:            id,
		UsernameValue: value,
		UserID:        uuid.UUID(userID),
		IsPrimary:     countForUser == 0,
		Status:        int64(status.Unset),
	})
	if err != nil {
		return Username{}, fmt.Errorf("create username on db: %w", err)
	}

	return Username{
		ID:      username.ID,
		Value:   username.UsernameValue,
		UserID:  username.UserID,
		Primary: username.IsPrimary,
		Status:  status.Status(username.Status),
	}, nil
}

func (s *Service) isValidUsername(value string) bool {
	for _, char := range value {
		if !strings.ContainsRune(s.config.AllowCharacters, char) {
			return false
		}
	}

	return true
}
