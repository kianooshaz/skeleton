package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	dbproto "github.com/kianooshaz/skeleton/foundation/database/proto"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
)

type UsernameStorage struct {
	Conn dbproto.QueryExecutor
}

const create = `
	INSERT INTO usernames (
    	id, username, account_id, status, created_at, updated_at, deleted_at
	) VALUES (
        $1, $2, $3, $4, NOW(), NOW(), NULL
    )
`

func (us *UsernameStorage) Create(ctx context.Context, username usernameproto.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, create, username.ID, username.Username, username.AccountID, username.Status)
	return err
}

const delete = `
	UPDATE
		usernames
	SET 
    	deleted_at = NOW()
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) Delete(ctx context.Context, id uuid.UUID) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, delete, id)
	return err
}

const get = `
	SELECT
		id, username, account_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) Get(ctx context.Context, id uuid.UUID) (usernameproto.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var username usernameproto.Username
	err := conn.QueryRowContext(ctx, get, id).
		Scan(&username.ID, &username.Username, &username.AccountID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usernameproto.Username{}, dbproto.ErrRowNotFound
		}
		return usernameproto.Username{}, err
	}

	return username, nil
}

const listByAccount = `
	SELECT
		id, username, account_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) ListWithSearch(ctx context.Context, req usernameproto.ListRequest) ([]usernameproto.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listByAccount + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(20))

	rows, err := conn.QueryContext(ctx, query, req.AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []usernameproto.Username
	for rows.Next() {
		var username usernameproto.Username
		if err := rows.Scan(
			&username.ID,
			&username.Username,
			&username.AccountID,
			&username.Status,
			&username.CreatedAt,
			&username.UpdatedAt,
		); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}

const CountWithSearch = `
	SELECT
		COUNT(id)
	FROM	
		usernames
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) CountWithSearch(ctx context.Context, req usernameproto.ListRequest) (int64, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, CountWithSearch, req.AccountID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, dbproto.ErrRowNotFound
		}
		return 0, err
	}

	return count, nil
}

const listByAccount2 = `
	SELECT
		id, username, account_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) ListByUserAndOrganization(ctx context.Context, req usernameproto.ListAssignedRequest) ([]usernameproto.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listByAccount2 + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(20))

	rows, err := conn.QueryContext(ctx, query, req.AccountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usernames := make([]usernameproto.Username, 0, req.Page.PageRows)
	for rows.Next() {
		var username usernameproto.Username
		err := rows.Scan(&username.ID, &username.Username, &username.AccountID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
		if err != nil {
			return nil, err
		}

		usernames = append(usernames, username)
	}

	return usernames, nil
}

const updateStatus = `
	UPDATE
		usernames
	SET 
		status = $2,
		updated_at = NOW()
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) UpdateStatus(ctx context.Context, username usernameproto.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, updateStatus, username.ID, username.Status)
	return err
}

const exist = `
	SELECT
		EXISTS(
			SELECT 1
			FROM usernames
			WHERE username = $1 AND deleted_at IS NULL
		)
`

func (us *UsernameStorage) Exist(ctx context.Context, username string) (bool, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var exists bool
	err := conn.QueryRowContext(ctx, exist, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

const countByAccount = `
	SELECT
		COUNT(id)
	FROM
		usernames
	WHERE
		account_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) CountByAccount(ctx context.Context, accountID accprotocol.AccountID) (int64, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, countByAccount, accountID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
