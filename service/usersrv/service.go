package usersrv

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kianooshaz/skeleton/protocol"
	"github.com/kianooshaz/skeleton/service/usersrv/stores/userdb"
)

type (
	Service struct {
		pool    *pgxpool.Pool
		queries *userdb.Queries
	}
)

func New(pool *pgxpool.Pool) protocol.ServiceUser {
	return &Service{
		pool:    pool,
		queries: userdb.New(),
	}
}
