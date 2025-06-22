// Package service provides the implementation of the UserService interface.
package uus

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	uup "github.com/kianooshaz/skeleton/services/user/user/protocol"
	"github.com/kianooshaz/skeleton/services/user/user/service/storage"
)

type (
	storer interface {
		Create(ctx context.Context, user uup.User) error
		Get(ctx context.Context, id iup.UserID) (uup.User, error)
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]uup.User, error)
		Count(ctx context.Context) (int, error)
	}

	service struct {
		logger      *slog.Logger
		storage     storer
		storageConn *sql.DB
	}
)

var Service uup.UserService = &service{}

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
		storageConn: postgres.ConnectionPool,
	}
}
