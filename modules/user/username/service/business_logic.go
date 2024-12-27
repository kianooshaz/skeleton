package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/status"
	"github.com/kianooshaz/skeleton/modules/user/username/protocol"
	"github.com/kianooshaz/skeleton/modules/user/username/service/stores/db"
)

func (s *Service) Add(ctx context.Context, userID uuid.UUID, value string) (protocol.Username, error) {
	id, err := uuid.NewV7()
	if err != nil {
		s.logger.Error("error at creating new uuid", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if len(value) < int(s.config.MinLength) || len(value) > int(s.config.MaxLength) {
		return &Username{}, derror.ErrUsernameLength
	}

	countValue, err := s.db.CountByUsername(ctx, value)
	if err != nil {
		s.logger.Error("error at getting count by username", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if countValue > 0 {
		return &Username{}, derror.ErrUsernameAlreadyExists
	}

	countForUser, err := s.db.CountByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		s.logger.Error("error at getting count by user id", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	if countForUser > int64(s.config.MaxPerUser) {
		return &Username{}, derror.ErrUsernameMaxPerUser
	}

	username, err := s.db.Create(ctx, db.CreateParams{
		ID:            id,
		UsernameValue: value,
		UserID:        userID,
		IsPrimary:     countForUser == 0,
		Status:        int64(status.Unset),
	})
	if err != nil {
		s.logger.Error("error at creating username in database", "error", err)

		return &Username{}, derror.ErrInternalSystem
	}

	return &Username{
		ID:      username.ID,
		Value:   username.UsernameValue,
		UserID:  username.UserID,
		Primary: username.IsPrimary,
		Status:  status.Status(username.Status),
	}, nil
}

// List implements protocol.ServiceUsername.
func (s *Service) Get(ctx context.Context, usernameValue string) (protocol.Username, error) {
	username, err := s.db.GetByUsername(ctx, usernameValue)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &Username{}, derror.ErrUserNotFound
		}
		s.logger.Error("error at getting username from database", "username", usernameValue, "error", err)

		return Username{}, derror.ErrInternalSystem
	}

	return &Username{
		ID:        username.ID,
		Value:     username.UsernameValue,
		UserID:    username.UserID,
		Primary:   username.IsPrimary,
		Status:    status.Status(username.Status),
		CreatedAt: username.CreatedAt.Time,
		UpdatedAt: username.UpdatedAt.Time,
	}, nil
}

// Search implements protocol.ServiceUsername.
func (s *Service) Search(ctx context.Context, value string) (Username, error) {
	panic("unimplemented")
}

// Count implements protocol.ServiceUsername.
func (s *Service) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	panic("unimplemented")
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, value string) error {
	panic("unimplemented")
}

// BePrimary implements protocol.ServiceUsername.
func (s *Service) BePrimary(ctx context.Context, id uuid.UUID) error {
	username, err := s.db.Get(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return derror.ErrUsernameNotFound
		}
		s.logger.Error("error at getting username from database", "id", id, "error", err)

		return derror.ErrInternalSystem
	}

	if username.IsPrimary {
		return nil
	}

	usernames, err := s.db.List(ctx, username.UserID)
	if err != nil {
		s.logger.Error("error at listing usernames", "error", err)

		return derror.ErrInternalSystem
	}

	tx, err := s._pdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		s.logger.Error("error at beginning transaction", "error", err)

		return derror.ErrInternalSystem
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			s.logger.Error("error at rolling back transaction", "error", err)
		}
	}()

	queries := db.New(tx)

	for _, u := range usernames {
		if u.IsPrimary {
			if err := queries.Update(ctx, db.UpdateParams{
				ID:        u.ID,
				IsPrimary: false,
				Status:    u.Status,
			}); err != nil {
				return fmt.Errorf("update username to be non primary: %w", err)
			}
		}
	}

	if err := queries.Update(ctx, db.UpdateParams{
		ID:        username.ID,
		IsPrimary: true,
		Status:    username.Status,
	}); err != nil {
		return fmt.Errorf("update username to be primary: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// Hidden implements protocol.ServiceUsername.
func (s *Service) Hidden(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// Unhidden implements protocol.ServiceUsername.
func (s *Service) Unhidden(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// Reserve implements protocol.ServiceUsername.
func (s *Service) Reserve(ctx context.Context, value string) error {
	panic("unimplemented")
}

// Unreserve implements protocol.ServiceUsername.
func (s *Service) Unreserve(ctx context.Context, value string) error {
	panic("unimplemented")
}

func (s *Service) isValidUsername(value string) bool {
	for _, char := range value {
		if !strings.ContainsRune(s.config.AllowCharacters, char) {
			return false
		}
	}

	return true
}
