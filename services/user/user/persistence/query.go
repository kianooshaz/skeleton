package persistence

import (
	"context"
	"database/sql"
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

const create = `
	INSERT INTO users (
		id
		created_at
	) VALUES (
		$1
		$2
	)
`

func (us *UserStorage) Create(ctx context.Context, user userproto.User) error {
	conn := session.GetDBConnection(ctx, us.Conn)

	_, err := conn.ExecContext(ctx, create, user.ID, user.CreatedAt)
	return err
}

const get = `
	SELECT
		id
		created_at
	FROM
		users
	WHERE
		id = $1
`

func (us *UserStorage) Get(ctx context.Context, id userproto.UserID) (userproto.User, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	row := conn.QueryRowContext(ctx, get, id)

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

const list = `
	SELECT
		id
		created_at
	FROM
		users
`

func (us *UserStorage) List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]userproto.User, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	list := list + page.String(pagination.SQLStringer(20)) + orderBy.String(oderStringer)

	rows, err := conn.QueryContext(ctx, list)
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

const count = `
	SELECT
		COUNT(*)
	FROM
		users
`

func (us *UserStorage) Count(ctx context.Context) (int, error) {
	row := us.Conn.QueryRowContext(ctx, count)
	var count int
	err := row.Scan(&count)
	return count, err
}
