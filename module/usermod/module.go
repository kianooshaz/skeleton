package usermod

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kianooshaz/skeleton/module/usermod/stores/userdb"
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
