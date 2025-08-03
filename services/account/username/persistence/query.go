package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/google/uuid"
	dbproto "github.com/kianooshaz/skeleton/foundation/database/proto"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/proto"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
)

type UsernameStorage struct {
	Conn dbproto.QueryExecutor
}

const defaultPageSize = 20

//go:embed queries/create.sql
var createQuery string

//go:embed queries/delete.sql
var deleteUsernameQuery string

//go:embed queries/get.sql
var getQuery string

//go:embed queries/list_by_account.sql
var listByAccountQuery string

//go:embed queries/count_with_search.sql
var countWithSearchQuery string

//go:embed queries/update_status.sql
var updateStatusQuery string

//go:embed queries/exist.sql
var existQuery string

//go:embed queries/count_by_account.sql
var countByAccountQuery string

func (us *UsernameStorage) Create(ctx context.Context, username usernameproto.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, createQuery, username.ID, username.Username, username.AccountID, username.Status)
	return err
}

func (us *UsernameStorage) Delete(ctx context.Context, id uuid.UUID) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, deleteUsernameQuery, id)
	return err
}

func (us *UsernameStorage) Get(ctx context.Context, id uuid.UUID) (usernameproto.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var username usernameproto.Username
	err := conn.QueryRowContext(ctx, getQuery, id).
		Scan(&username.ID, &username.Username, &username.AccountID, &username.Status, &username.CreatedAt, &username.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return usernameproto.Username{}, dbproto.ErrRowNotFound
		}
		return usernameproto.Username{}, err
	}

	return username, nil
}

func (us *UsernameStorage) ListWithSearch(ctx context.Context, req usernameproto.ListRequest) ([]usernameproto.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listByAccountQuery + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(defaultPageSize))

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

func (us *UsernameStorage) CountWithSearch(ctx context.Context, req usernameproto.ListRequest) (int64, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, countWithSearchQuery, req.AccountID).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, dbproto.ErrRowNotFound
		}
		return 0, err
	}

	return count, nil
}

func (us *UsernameStorage) ListByUserAndOrganization(ctx context.Context, req usernameproto.ListAssignedRequest) ([]usernameproto.Username, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listByAccountQuery + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(defaultPageSize))

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

func (us *UsernameStorage) UpdateStatus(ctx context.Context, username usernameproto.Username) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, updateStatusQuery, username.ID, username.Status)
	return err
}

func (us *UsernameStorage) Exist(ctx context.Context, username string) (bool, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var exists bool
	err := conn.QueryRowContext(ctx, existQuery, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (us *UsernameStorage) CountByAccount(ctx context.Context, accountID accprotocol.AccountID) (int64, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	var count int64
	err := conn.QueryRowContext(ctx, countByAccountQuery, accountID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
