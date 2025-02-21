package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/kianooshaz/skeleton/foundation/database"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	"github.com/kianooshaz/skeleton/foundation/types"
	"github.com/kianooshaz/skeleton/modules/user/username/protocol"
)

type UsernameStorage struct {
	Conn database.ConnectionProtocol
}

const create = `
	INSERT INTO usernames (
    	id, user_id, organization_id, status, created_at, updated_at, deleted_at
	) VALUES (
        $1, $2, $3, $4, NOW(), NOW(), NULL
    )
`

func (us *UsernameStorage) Create(ctx context.Context, username protocol.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, row := conn.ExecContext(ctx, create, username.ID, username.UserID, username.OrganizationID, username.Status)
	return row
}

const delete = `
	UPDATE
		usernames
	SET 
    	deleted_at = NOW()
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) Delete(ctx context.Context, id string) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, delete, id)
	return err
}

const get = `
	SELECT
		id, user_id, organization_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) Get(ctx context.Context, id string) (protocol.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	row := conn.QueryRowContext(ctx, get, id)

	var username protocol.Username
	err := row.Scan(&username.ID, &username.UserID, &username.OrganizationID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.Username{}, derror.ErrUsernameNotFound
		}

		return protocol.Username{}, err
	}

	return username, err
}

const listByUser = `
	SELECT
		id, user_id, organization_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		user_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) ListByUser(
	ctx context.Context,
	userID types.UserID,
	orderBy order.OrderBy,
	page pagination.Page,
	isPrimary bool,
) ([]protocol.Username, error) {

	var query string

	if isPrimary {
		query = listByUser + fmt.Sprintf(" AND status & %d = %d ", types.Primary, types.Primary)
	}

	query = query + orderBy.String(oderStringer) + page.String(pagination.SQLStringer(20))

	rows, err := us.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usernames := make([]protocol.Username, 0, page.PageRows)
	for rows.Next() {
		var username protocol.Username
		err := rows.Scan(&username.ID, &username.UserID, &username.OrganizationID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
		if err != nil {
			return nil, err
		}

		usernames = append(usernames, username)
	}

	return usernames, nil
}

const listByUserAndOrganization = `
	SELECT
		id, user_id, organization_id, status, created_at, updated_at
	FROM
		usernames
	WHERE
		user_id = $1 AND organization_id = $2 AND deleted_at IS NULL
`

func (us *UsernameStorage) ListByUserAndOrganization(
	ctx context.Context,
	userID types.UserID,
	organizationID types.OrganizationID,
	orderBy order.OrderBy,
	page pagination.Page,
	isPrimary bool,
) ([]protocol.Username, error) {
	var query string

	if isPrimary {
		query = listByUserAndOrganization + fmt.Sprintf(" AND status & %d = %d ", types.Primary, types.Primary)
	}

	query = query + orderBy.String(oderStringer) + page.String(pagination.SQLStringer(20))

	rows, err := us.Conn.QueryContext(ctx, query, userID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usernames := make([]protocol.Username, 0, page.PageRows)
	for rows.Next() {
		var username protocol.Username
		err := rows.Scan(&username.ID, &username.UserID, &username.OrganizationID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
		if err != nil {
			return nil, err
		}

		usernames = append(usernames, username)
	}

	return usernames, nil
}

const update = `
	UPDATE
		usernames
	SET 
		user_id = $2,
		organization_id = $3,
		status = $4,
		updated_at = NOW()
	WHERE
		id = $1 AND deleted_at IS NULL
`

type UpdateRequest struct {
	ID             string
	UserID         types.UserID
	OrganizationID types.OrganizationID
	Status         types.Status
}

func (us *UsernameStorage) UpdateStatus(ctx context.Context, req protocol.Username) error {
	_, err := us.Conn.ExecContext(ctx, update, req.ID, req.UserID, req.OrganizationID, req.Status)
	return err
}

const count = `
	SELECT
		COUNT(id)
	FROM
		usernames
	WHERE
		id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) Count(ctx context.Context, id string) (int64, error) {
	row := us.Conn.QueryRowContext(ctx, count, id)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const countByUser = `
	SELECT
		COUNT(id)
	FROM
		usernames
	WHERE
		user_id = $1 AND deleted_at IS NULL
`

func (us *UsernameStorage) CountByUser(ctx context.Context, userID types.UserID) (int64, error) {
	row := us.Conn.QueryRowContext(ctx, countByUser, userID)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

const countByUserAndOrganization = `
	SELECT
		COUNT(id)
	FROM
		usernames
	WHERE
		user_id = $1 AND organization_id = $2 deleted_at IS NULL
`

func (us *UsernameStorage) CountByUserAndOrganization(ctx context.Context, userID types.UserID, organization types.OrganizationID) (int64, error) {
	row := us.Conn.QueryRowContext(ctx, countByUserAndOrganization, userID, organization)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
