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
	"github.com/kianooshaz/skeleton/modules/user/user/protocol"
	"github.com/kianooshaz/skeleton/modules/user/user/service/storage"
)

type Storage interface {
	Create(ctx context.Context, user protocol.User) error
	Get(ctx context.Context, id types.UserID) (protocol.User, error)
	List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]protocol.User, error)
	Count(ctx context.Context) (int, error)
}

var Service protocol.UserService = &userService{}

type userService struct {
	logger  slog.Logger
	storage Storage
	dbConn  *sql.DB
}

func Init() {
	Service = &userService{
		logger: *slog.With(
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
