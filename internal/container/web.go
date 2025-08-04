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

// Close gracefully shuts down the web container and all its services.
func (c *WebContainer) Close() error {
	var errs []error

	// Shutdown audit service workers
	if c.AuditService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		c.AuditService.Shutdown(ctx)
		cancel()
	}

	// Close base container (database, etc.)
	if c.BaseContainer != nil {
		if err := c.BaseContainer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing base container: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
