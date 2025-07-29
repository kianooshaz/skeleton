package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	dp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	fdp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
)

type PasswordStorage struct {
	Conn dp.QueryExecutor
}

const create = `
	INSERT INTO passwords (
    	id, account_id, password_hash, created_at, updated_at, deleted_at
	) VALUES (
        $1, $2, $3, NOW(), NOW(), NULL
    )
`

func (ps *PasswordStorage) Create(ctx context.Context, password passwordproto.Password) error {
	conn := session.GetDBConnection(ctx, ps.Conn)

	_, err := conn.ExecContext(ctx, create, password.ID, password.AccountID, password.PasswordHash)
	return err
}

const deletePassword = `
	UPDATE
		passwords
	SET 
    	deleted_at = NOW()
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (ps *PasswordStorage) Delete(ctx context.Context, id uuid.UUID) error {
	conn := session.GetDBConnection(ctx, ps.Conn)

	_, err := conn.ExecContext(ctx, deletePassword, id)
	return err
}

const get = `
	SELECT
		id, account_id, password_hash, created_at, updated_at
	FROM
		passwords
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (ps *PasswordStorage) Get(ctx context.Context, id uuid.UUID) (passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	var password passwordproto.Password
	err := conn.QueryRowContext(ctx, get, id).
		Scan(&password.ID, &password.AccountID, &password.PasswordHash, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return passwordproto.Password{}, fdp.ErrRowNotFound
		}
		return passwordproto.Password{}, err
	}

	return password, nil
}

const getByAccountID = `
	SELECT
		id, account_id, password_hash, created_at, updated_at
	FROM
		passwords
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (ps *PasswordStorage) GetByAccountID(ctx context.Context, accountID accprotocol.AccountID) (passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	var password passwordproto.Password
	err := conn.QueryRowContext(ctx, getByAccountID, accountID).
		Scan(&password.ID, &password.AccountID, &password.PasswordHash, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return passwordproto.Password{}, fdp.ErrRowNotFound
		}
		return passwordproto.Password{}, err
	}

	return password, nil
}

const listByAccount = `
	SELECT
		id, account_id, password_hash, created_at, updated_at
	FROM
		passwords
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (ps *PasswordStorage) ListWithSearch(ctx context.Context, req passwordproto.ListRequest) ([]passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	query := listByAccount + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(20))

	rows, err := conn.QueryContext(ctx, query, req.AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []passwordproto.Password
	for rows.Next() {
		var password passwordproto.Password
		if err := rows.Scan(
			&password.ID,
			&password.AccountID,
			&password.PasswordHash,
			&password.CreatedAt,
			&password.UpdatedAt,
		); err != nil {
			return nil, err
		}
		passwords = append(passwords, password)
	}

	return passwords, nil
}

const CountWithSearch = `
	SELECT
		COUNT(id)
	FROM	
		passwords
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (ps *PasswordStorage) CountWithSearch(ctx context.Context, req passwordproto.ListRequest) (int64, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, CountWithSearch, req.AccountID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fdp.ErrRowNotFound
		}
		return 0, err
	}

	return count, nil
}

const history = `
	SELECT
		id, account_id, password_hash, created_at, updated_at
	FROM
		passwords
	WHERE
		account_id = $1 AND deleted_at IS NULL
	ORDER BY created_at DESC
	LIMIT $2
`

func (ps *PasswordStorage) History(ctx context.Context, accountID accprotocol.AccountID, limit int32) ([]passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	rows, err := conn.QueryContext(ctx, history, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []passwordproto.Password
	for rows.Next() {
		var password passwordproto.Password
		if err := rows.Scan(
			&password.ID,
			&password.AccountID,
			&password.PasswordHash,
			&password.CreatedAt,
			&password.UpdatedAt,
		); err != nil {
			return nil, err
		}
		passwords = append(passwords, password)
	}

	return passwords, nil
}
