package usernamesrv

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/protocol"
	"github.com/kianooshaz/skeleton/protocol/derror"
	"github.com/kianooshaz/skeleton/protocol/status"
	"github.com/kianooshaz/skeleton/service/usernamesrv/stores/usernamedb"
)

func (s *Service) Add(ctx context.Context, userID protocol.ID, value string) (protocol.Username, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return usernameModel{}, fmt.Errorf("new uuid: %w", err)
	}

	if len(value) < int(s.config.MinLength) || len(value) > int(s.config.MaxLength) {
		return usernameModel{}, derror.InvalidRequest
	}

	countValue, err := s.queries.CountByUsername(ctx, value)
	if err != nil {
		return usernameModel{}, fmt.Errorf("get count by username: %w", err)
	}

	if countValue > 0 {
		return usernameModel{}, derror.Duplicate
	}

	countForUser, err := s.queries.CountByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return usernameModel{}, fmt.Errorf("get count by user id: %w", err)
	}

	if countForUser > int64(s.config.MaxPerUser) {
		return usernameModel{}, derror.InvalidRequest
	}

	username, err := s.queries.Create(ctx, usernamedb.CreateParams{
		ID:            id,
		UsernameValue: value,
		UserID:        uuid.UUID(userID),
		IsPrimary:     countForUser == 0,
		Status:        int64(status.Registered),
	})
	if err != nil {
		return usernameModel{}, fmt.Errorf("create username on db: %w", err)
	}

	return usernameModel{
		id:      username.ID,
		value:   username.UsernameValue,
		userID:  username.UserID,
		primary: username.IsPrimary,
		status:  status.Status(username.Status),
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
