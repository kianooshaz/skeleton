package database

import (
	"context"
	"database/sql"

	"github.com/kianooshaz/skeleton/foundation/session"
)

type Database struct {
	cli *sql.DB
}

func (db *Database) Insert(ctx context.Context, id int) error {
	cli := session.GetQueryExecutor(ctx, db.cli)

}
