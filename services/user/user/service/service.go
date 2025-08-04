// Package service provides the implementation of the UserService interface.
package uus

import (
	"context"
	"database/sql"
	"log/slog"

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

// Service is the global service instance for backward compatibility
// TODO: Remove this after all dependencies are migrated to DI
var Service userproto.UserService

// New creates a new user service instance.
func New(db *sql.DB, logger *slog.Logger) userproto.UserService {
	serviceLogger := logger.With(
		slog.Group("package_info",
			slog.String("module", "user"),
			slog.String("service", "user"),
		),
	)

	svc := &service{
		logger: serviceLogger,
		persister: &persistence.UserStorage{
			Conn: db,
		},
		dbConn: db,
	}

	// Set global service for backward compatibility
	if Service == nil {
		Service = svc
	}

	return svc
}
