package usersrv

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kianooshaz/skeleton/protocol"
	"github.com/kianooshaz/skeleton/service/usersrv/stores/userdb"
)

type (
	Service struct {
		queries *userdb.Queries
	}
)

func New(pool *pgxpool.Pool) protocol.ServiceUser {
	return &Service{
		queries: userdb.New(pool),
	}
}

// NewTx implements protocol.ServiceUser.
func (m *Service) NewWithTx(tx pgx.Tx) protocol.ServiceUser {
	return &Service{
		queries: userdb.New(tx),
	}
}
