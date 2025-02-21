package service

import (
	"database/sql"

	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/modules/user/userid/protocol"
)

type (
	UserService struct {
		logger log.Logger
		db     *sql.DB
	}
)

func New(logger log.Logger, db *sql.DB) protocol.UserService {
	return &UserService{
		logger: logger,
		db:     db,
	}
}
