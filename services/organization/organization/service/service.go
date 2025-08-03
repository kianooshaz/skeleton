// Package service provides the implementation of the OrganizationService interface.
package oos

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
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

	service struct {
		logger    *slog.Logger
		persister persister
		dbConn    *sql.DB
	}
)

var Service orgproto.OrganizationService = &service{}

func init() {
	Service = &service{
		logger: slog.With(
			slog.Group("package_info",
				slog.String("module", "organization"),
				slog.String("service", "organization"),
			),
		),
		persister: &persistence.OrganizationStorage{
			Conn: postgres.ConnectionPool,
		},
		dbConn: postgres.ConnectionPool,
	}
}
