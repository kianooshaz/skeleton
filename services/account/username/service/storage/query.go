package storage

import (
	"context"
	"database/sql"
	"errors"

	dp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	fdp "github.com/kianooshaz/skeleton/foundation/database/protocol"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	aunp "github.com/kianooshaz/skeleton/services/account/username/protocol"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

type UsernameStorage struct {
	Conn dp.QueryExecutor
}

const create = `
	INSERT INTO usernames (
    	id, username, user_id, organization_id, status, created_at, updated_at, deleted_at
	) VALUES (
        $1, $2, $3, $4, $5, NOW(), NOW(), NULL
    )
`

func (us *UsernameStorage) Create(ctx context.Context, username aunp.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, create, username.ID, username.Username, username.UserID, username.OrganizationID, username.Status)
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

func (us *UsernameStorage) Delete(ctx context.Context, id iunp.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, delete, id)
	return err
}

const get = `
	SELECT
		id, username, user_id, organization_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) Get(ctx context.Context, id iunp.Username) (aunp.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var username aunp.Username
	err := conn.QueryRowContext(ctx, get, id).
		Scan(&username.ID, &username.Username, &username.UserID, &username.OrganizationID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return aunp.Username{}, fdp.ErrRowNotFound
		}
		return aunp.Username{}, err
	}

	return username, nil
}

const listByUser = `
	SELECT
		id, username, user_id, organization_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		user_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) ListWithSearch(ctx context.Context, req aunp.ListRequest) ([]aunp.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listByUser + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(20))

	rows, err := conn.QueryContext(ctx, query, req.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []aunp.Username
	for rows.Next() {
		var username aunp.Username
		if err := rows.Scan(
			&username.ID,
			&username.Username,
			&username.UserID,
			&username.OrganizationID,
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
		user_id = $1 AND organization_id = $2 AND deleted_at IS NULL
`

func (us *UsernameStorage) CountWithSearch(ctx context.Context, req aunp.ListRequest) (int64, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, CountWithSearch, req.UserID, req.OrganizationID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fdp.ErrRowNotFound
		}
		return 0, err
	}

	return count, nil
}

const listByUserAndOrganization = `
	SELECT
		id, username, user_id, organization_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		user_id = $1 AND organization_id = $2 AND deleted_at IS NULL
`

func (us *UsernameStorage) ListByUserAndOrganization(ctx context.Context, req aunp.ListAssignedRequest) ([]aunp.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listByUserAndOrganization + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(20))

	rows, err := conn.QueryContext(ctx, query, req.UserID, req.OrganizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usernames := make([]aunp.Username, 0, req.Page.PageRows)
	for rows.Next() {
		var username aunp.Username
		err := rows.Scan(&username.ID, &username.Username, &username.UserID, &username.OrganizationID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
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

func (us *UsernameStorage) UpdateStatus(ctx context.Context, username aunp.Username) error {
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

func (us *UsernameStorage) Exist(ctx context.Context, username iunp.Username) (bool, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var exists bool
	err := conn.QueryRowContext(ctx, exist, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

const countByUserAndOrganization = `
	SELECT
		COUNT(id)
	FROM
		usernames
	WHERE
		user_id = $1 AND organization_id = $2 deleted_at IS NULL
`

func (us *UsernameStorage) CountByUserAndOrganization(ctx context.Context, userID iup.UserID, organizationID iop.OrganizationID) (int64, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, countByUserAndOrganization, userID, organizationID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
