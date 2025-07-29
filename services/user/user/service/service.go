// Package service provides the implementation of the UserService interface.
package uus

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/services/user/user/persistence"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

type (
	persister interface {
		Create(ctx context.Context, user userproto.User) error
		Get(ctx context.Context, id userproto.UserID) (userproto.User, error)
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]userproto.User, error)
		Count(ctx context.Context) (int, error)
	}

	service struct {
		logger    *slog.Logger
		persister persister
		dbConn    *sql.DB
	}
)

var Service userproto.UserService = &service{}

func init() {
	Service = &service{
		logger: slog.With(
			slog.Group("package_info",
				slog.String("module", "user"),
				slog.String("service", "user"),
			),
		),
		persister: &persistence.UserStorage{
			Conn: postgres.ConnectionPool,
		},
		dbConn: postgres.ConnectionPool,
	}
}
