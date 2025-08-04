// Package birthdayservice provides the implementation of the BirthdayService interface.
package birthdayservice

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/services/user/birthday/persistence"
	birthdayproto "github.com/kianooshaz/skeleton/services/user/birthday/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

type (
	// Config holds the configuration for the birthday service.
	Config struct {
		// MaxAge defines the maximum age allowed for a birthday (default: 150).
		MaxAge int `yaml:"max_age" validate:"min=1,max=200"`
		// MinAge defines the minimum age allowed for a birthday (default: 0).
		MinAge int `yaml:"min_age" validate:"min=0,max=100"`
	}

	persister interface {
		Create(ctx context.Context, birthday birthdayproto.Birthday) error
		Get(ctx context.Context, id birthdayproto.BirthdayID) (birthdayproto.Birthday, error)
		GetByUserID(ctx context.Context, userID userproto.UserID) (birthdayproto.Birthday, error)
		Update(ctx context.Context, birthday birthdayproto.Birthday) error
		Delete(ctx context.Context, id birthdayproto.BirthdayID) error
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy, filters persistence.ListFilters) ([]birthdayproto.Birthday, error)
		Count(ctx context.Context, filters persistence.ListFilters) (int, error)
		ExistsByUserID(ctx context.Context, userID userproto.UserID) (bool, error)
	}

	Service struct {
		config    Config
		logger    *slog.Logger
		persister persister
		dbConn    *sql.DB
	}
)

// New creates a new birthday service instance.
func New(cfg Config, db *sql.DB, logger *slog.Logger) birthdayproto.BirthdayService {
	serviceLogger := logger.With(
		slog.Group("package_info",
			slog.String("module", "user"),
			slog.String("service", "birthday"),
		),
	)

	return &Service{
		config: cfg,
		logger: serviceLogger,
		persister: &persistence.BirthdayStorage{
			Conn: db,
		},
		dbConn: db,
	}
}
