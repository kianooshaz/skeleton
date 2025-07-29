package persistence

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	dbproto "github.com/kianooshaz/skeleton/foundation/database/proto"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

type UserStorage struct {
	Conn dbproto.QueryExecutor
}

const defaultPageSize = 20

//go:embed queries/create.sql
var createQuery string

//go:embed queries/get.sql
var getQuery string

//go:embed queries/list.sql
var listQuery string

//go:embed queries/count.sql
var countQuery string

func (us *UserStorage) Create(ctx context.Context, user userproto.User) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, createQuery, user.ID, user.CreatedAt)
	return err
}

func (us *UserStorage) Get(ctx context.Context, id userproto.UserID) (userproto.User, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	row := conn.QueryRowContext(ctx, getQuery, id)

	var user userproto.User
	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userproto.User{}, derror.ErrUserNotFound
		}

		return userproto.User{}, err
	}

	return user, err
}

func (us *UserStorage) List(ctx context.Context, page pagination.Page,
	orderBy order.OrderBy) ([]userproto.User, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	query := listQuery + page.String(pagination.SQLStringer(defaultPageSize)) + orderBy.String(oderStringer)

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []userproto.User
	for rows.Next() {
		var user userproto.User
		err := rows.Scan(&user.ID, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (us *UserStorage) Count(ctx context.Context) (int, error) {
	row := us.Conn.QueryRowContext(ctx, countQuery)
	var count int
	err := row.Scan(&count)
	return count, err
}
