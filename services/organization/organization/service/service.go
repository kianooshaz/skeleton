// Package service provides the implementation of the OrganizationService interface.
package orgservice

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/services/organization/organization/persistence"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
)

type (
	persister interface {
		Create(ctx context.Context, organization orgproto.Organization) error
		Get(ctx context.Context, id orgproto.OrganizationID) (orgproto.Organization, error)
		List(ctx context.Context, page pagination.Page, orderBy order.OrderBy) ([]orgproto.Organization, error)
		Count(ctx context.Context) (int, error)
	}

	Service struct {
		logger    *slog.Logger
		persister persister
		dbConn    *sql.DB
	}
)

var _ orgproto.OrganizationService = (*Service)(nil)

// New creates a new organization service instance.
func New(db *sql.DB, logger *slog.Logger) *Service {
	serviceLogger := logger.With(
		slog.Group("package_info",
			slog.String("module", "organization"),
			slog.String("service", "organization"),
		),
	)

	svc := &Service{
		logger: serviceLogger,
		persister: &persistence.OrganizationStorage{
			Conn: db,
		},
		dbConn: db,
	}

	return svc
}
