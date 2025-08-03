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
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
)

type PasswordStorage struct {
	Conn dbproto.QueryExecutor
}

const defaultPageSize = 20

//go:embed queries/create.sql
var createQuery string

//go:embed queries/delete.sql
var deletePasswordQuery string

//go:embed queries/get.sql
var getQuery string

//go:embed queries/get_by_account_id.sql
var getByAccountIDQuery string

//go:embed queries/list_by_account.sql
var listByAccountQuery string

//go:embed queries/count_with_search.sql
var countWithSearchQuery string

//go:embed queries/history.sql
var historyQuery string

func (ps *PasswordStorage) Create(ctx context.Context, password passwordproto.Password) error {
	conn := session.GetDBConnection(ctx, ps.Conn)

	_, err := conn.ExecContext(ctx, createQuery, password.ID, password.AccountID, password.PasswordHash)
	return err
}

func (ps *PasswordStorage) Delete(ctx context.Context, id uuid.UUID) error {
	conn := session.GetDBConnection(ctx, ps.Conn)

	_, err := conn.ExecContext(ctx, deletePasswordQuery, id)
	return err
}

func (ps *PasswordStorage) Get(ctx context.Context, id uuid.UUID) (passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	var password passwordproto.Password
	err := conn.QueryRowContext(ctx, getQuery, id).
		Scan(&password.ID, &password.AccountID, &password.PasswordHash, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return passwordproto.Password{}, dbproto.ErrRowNotFound
		}
		return passwordproto.Password{}, err
	}

	return password, nil
}

func (ps *PasswordStorage) GetByAccountID(ctx context.Context, accountID accprotocol.AccountID) (passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	var password passwordproto.Password
	err := conn.QueryRowContext(ctx, getByAccountIDQuery, accountID).
		Scan(&password.ID, &password.AccountID, &password.PasswordHash, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return passwordproto.Password{}, dbproto.ErrRowNotFound
		}
		return passwordproto.Password{}, err
	}

	return password, nil
}

func (ps *PasswordStorage) ListWithSearch(ctx context.Context, req passwordproto.ListRequest) ([]passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	query := listByAccountQuery + req.OrderBy.String(oderStringer) + req.Page.String(pagination.SQLStringer(defaultPageSize))

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

func (ps *PasswordStorage) CountWithSearch(ctx context.Context, req passwordproto.ListRequest) (int64, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

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

func (ps *PasswordStorage) History(ctx context.Context, accountID accprotocol.AccountID, limit int32) ([]passwordproto.Password, error) {
	conn := session.GetDBConnection(ctx, ps.Conn)

	rows, err := conn.QueryContext(ctx, historyQuery, accountID, limit)
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
