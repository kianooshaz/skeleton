package usersrv

import (
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/service/usersrv/stores/userdb"
)

type (
	Service struct {
		db      postgres.DB
		queries *userdb.Queries
	}
)

func New(pool postgres.DB) *Service {
	return &Service{
		db:      pool,
		queries: userdb.New(pool),
	}
}

// NewTx implements protocol.ServiceUser.
func (m *Service) NewWithTx(tx pgx.Tx) *Service {
	return &Service{
		db:      m.db,
		queries: userdb.New(tx),
	}
}
