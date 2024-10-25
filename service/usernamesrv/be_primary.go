package usernamesrv

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/service/usernamesrv/stores/usernamedb"
)

// BePrimary implements protocol.ServiceUsername.
func (s *Service) BePrimary(ctx context.Context, id uuid.UUID) (err error) {
	username, err := s.queries.Get(ctx, uuid.UUID(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("get username: %w", err)
	}

	if username.IsPrimary {
		return nil
	}

	usernames, err := s.queries.List(ctx, username.UserID)
	if err != nil {
		return fmt.Errorf("list of usernames: %w", err)
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func(tx pgx.Tx) {
		if terr := tx.Rollback(ctx); terr != nil && err == nil && !errors.Is(terr, pgx.ErrTxClosed) {
			err = fmt.Errorf("rollback transaction: %w", terr)
		}
	}(tx)

	queries := usernamedb.New(tx)

	for _, u := range usernames {
		if u.IsPrimary {
			if err := queries.Update(ctx, usernamedb.UpdateParams{
				ID:        u.ID,
				IsPrimary: false,
				Status:    u.Status,
			}); err != nil {
				return fmt.Errorf("update username to be non primary: %w", err)
			}
		}
	}

	if err := queries.Update(ctx, usernamedb.UpdateParams{
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
