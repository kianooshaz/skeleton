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
	birthdayproto "github.com/kianooshaz/skeleton/services/user/birthday/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

// WebContainer holds all dependencies for the web application.
type WebContainer struct {
	config              *AppConfig
	logger              *slog.Logger
	db                  *sql.DB
	webService          protocol.WebService
	userService         userproto.UserService
	organizationService orgproto.OrganizationService
	passwordService     passwordproto.PasswordService
	usernameService     usernameproto.UsernameService
	auditService        auditproto.AuditService
	birthdayService     birthdayproto.BirthdayService
}

// Start initializes and starts all services.
func (c *WebContainer) Start(cancel context.CancelFunc) error {
	go func() {
		if err := c.webService.Start(); err != nil {
			c.logger.Error("Failed to start web service", "error", err)
			cancel()
		}
	}()

	return nil
}

// Stop gracefully shuts down all services.
func (c *WebContainer) Stop() error {
	c.logger.Info("Starting graceful shutdown of web container")

	// Create a timeout context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), c.config.ShutdownTimeout)
	defer cancel()

	// Stop services in reverse order of dependency
	if c.webService != nil {
		c.logger.Info("Shutting down web service")
		if err := c.webService.Shutdown(ctx); err != nil {
			c.webService.Close()
			c.logger.Error("Failed to shutdown web service", "error", err)
		}
	}

	if c.auditService != nil {
		c.logger.Info("Shutting down audit service")
		c.auditService.Shutdown(ctx)
	}

	if c.db != nil {
		c.logger.Info("Closing database connection")
		if err := c.db.Close(); err != nil {
			c.logger.Error("Failed to close database connection", "error", err)
		}
	}

	return nil
}

func (c *WebContainer) Logger() *slog.Logger {
	return c.logger
}
