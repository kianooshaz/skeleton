package container

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

// WebContainer holds all dependencies for the web application.
type WebContainer struct {
	Config              *AppConfig
	Logger              *slog.Logger
	DB                  *sql.DB
	WebService          protocol.WebService
	UserService         userproto.UserService
	OrganizationService orgproto.OrganizationService
	PasswordService     passwordproto.PasswordService
	UsernameService     usernameproto.UsernameService
	AuditService        auditproto.AuditService
}

// Start initializes and starts all services.
func (c *WebContainer) Start() error {
	return c.WebService.Start()
}

// Stop gracefully shuts down all services.
func (c *WebContainer) Stop(ctx context.Context) error {
	// Stop audit service first to finish processing records
	c.AuditService.Shutdown(ctx)

	return c.WebService.Shutdown(ctx)
}

// Close closes all resources.
func (c *WebContainer) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
