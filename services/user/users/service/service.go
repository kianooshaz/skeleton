// Package service provides the implementation of the UserService interface.
package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/types"
	up "github.com/kianooshaz/skeleton/services/user/user/protocol"
	"github.com/kianooshaz/skeleton/services/user/user/service/storage"
)

type (
	storer interface {
		Create(ctx context.Context, user up.User) error
		Get(ctx context.Context, id types.UserID) (up.User, error)
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]up.User, error)
		Count(ctx context.Context) (int, error)
	}

	service struct {
		logger  *slog.Logger
		storage storer
		dbConn  *sql.DB
	}
)

var Service up.UserService = &service{}

func init() {
	Service = &service{
		logger: slog.With(
			slog.Group("package_info",
				slog.String("module", "user"),
				slog.String("service", "user"),
			),
		),
		storage: &storage.UserStorage{
			Conn: postgres.ConnectionPool,
		},
		dbConn: postgres.ConnectionPool,
	}
}
