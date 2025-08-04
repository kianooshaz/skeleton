// Package container provides web-specific dependency injection container using Google Wire.
package container

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

const (
	// DefaultShutdownTimeout is the default timeout for graceful shutdown operations.
	DefaultShutdownTimeout = 30 * time.Second
)

// WebContainer holds web-specific dependencies.
type WebContainer struct {
	*BaseContainer
	RestServer      protocol.WebService
	UserService     userproto.UserService
	OrgService      orgproto.OrganizationService
	UsernameService usernameproto.UsernameService
	AuditService    auditproto.AuditService
	PasswordService passwordproto.PasswordService
}

// Start initializes and starts the web container and all its services.
func (c *WebContainer) Start() error {
	// Start base container first
	if c.BaseContainer != nil {
		if err := c.BaseContainer.Start(); err != nil {
			return fmt.Errorf("starting base container: %w", err)
		}
	}

	// Start the REST server
	if c.RestServer != nil {
		if err := c.RestServer.Start(); err != nil {
			return fmt.Errorf("starting REST server: %w", err)
		}
	}

	return nil
}

// Stop gracefully stops the web container and all its services.
func (c *WebContainer) Stop(ctx context.Context) error {
	var errs []error

	// Stop the REST server first
	if c.RestServer != nil {
		if err := c.RestServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("shutting down REST server: %w", err))
		}
	}

	// Shutdown audit service workers
	if c.AuditService != nil {
		auditCtx, cancel := context.WithTimeout(ctx, DefaultShutdownTimeout)
		c.AuditService.Shutdown(auditCtx)
		cancel()
	}

	// Stop base container
	if c.BaseContainer != nil {
		if err := c.BaseContainer.Stop(ctx); err != nil {
			errs = append(errs, fmt.Errorf("stopping base container: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Close gracefully shuts down the web container and all its services.
func (c *WebContainer) Close() error {
	// Close delegates to Stop with a default timeout
	ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
	defer cancel()
	return c.Stop(ctx)
}
