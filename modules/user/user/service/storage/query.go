package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kianooshaz/skeleton/foundation/database"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/session"
	"github.com/kianooshaz/skeleton/foundation/types"
	"github.com/kianooshaz/skeleton/modules/user/user/protocol"
)

type UserStorage struct {
	Conn database.ConnectionProtocol
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

func (us *UserStorage) Create(ctx context.Context, user protocol.User) error {
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

func (us *UserStorage) Get(ctx context.Context, id types.UserID) (protocol.User, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	row := conn.QueryRowContext(ctx, get, id)

	var user protocol.User
	err := row.Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return protocol.User{}, derror.ErrUserNotFound
		}

		return protocol.User{}, err
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

func (us *UserStorage) List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]protocol.User, error) {
	conn := session.GetDBConnection(ctx, us.Conn)

	list := list + page.String(pagination.SQLStringer(20)) + orderBy.String(oderStringer)

	rows, err := conn.QueryContext(ctx, list)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []protocol.User
	for rows.Next() {
		var user protocol.User
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

func (us *UserStorage) Count(ctx context.Context) (int64, error) {
	row := us.Conn.QueryRowContext(ctx, count)
	var count int64
	err := row.Scan(&count)
	return count, err
}
