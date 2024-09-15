package usersrv

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kianooshaz/skeleton/service/usersrv/stores/userdb"
)

type (
	Module struct {
		pool    *pgxpool.Pool
		queries *userdb.Queries
	}
)

func New(pool *pgxpool.Pool) *Module {
	return &Module{
		pool:    pool,
		queries: userdb.New(),
	}
}

type UserI interface {
	New(ctx context.Context) (User, error)
}
