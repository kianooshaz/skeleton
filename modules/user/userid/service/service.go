package service

import (
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/modules/user/userid/protocol"
	"github.com/kianooshaz/skeleton/modules/user/userid/service/stores/db"
)

type (
	UserService struct {
		logger log.Logger
		_pdb   postgres.DB
		db     *db.Queries
	}
)

func New(logger log.Logger, pdb postgres.DB) protocol.UserService {
	return &UserService{
		logger: logger,
		_pdb:   pdb,
		db:     db.New(pdb),
	}

}

// NewTx implements protocol.ServiceUser.
func (m *UserService) NewWithTx(tx pgx.Tx) protocol.UserService {
	return &UserService{
		logger: m.logger,
		_pdb:   m._pdb,
		db:     db.New(tx),
	}
}
