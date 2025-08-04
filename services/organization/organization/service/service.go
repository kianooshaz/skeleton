// Package service provides the implementation of the OrganizationService interface.
package oos

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

	service struct {
		logger    *slog.Logger
		persister persister
		dbConn    *sql.DB
	}
)

// Service is the global service instance for backward compatibility
// TODO: Remove this after all dependencies are migrated to DI
var Service orgproto.OrganizationService

// New creates a new organization service instance
func New(db *sql.DB, logger *slog.Logger) orgproto.OrganizationService {
	serviceLogger := logger.With(
		slog.Group("package_info",
			slog.String("module", "organization"),
			slog.String("service", "organization"),
		),
	)

	svc := &service{
		logger: serviceLogger,
		persister: &persistence.OrganizationStorage{
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
