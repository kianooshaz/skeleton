// Package container provides dependency injection container using Google Wire.
package container

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	usernameservice "github.com/kianooshaz/skeleton/services/account/username/service"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	passwordservice "github.com/kianooshaz/skeleton/services/authentication/password/service"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	auditservice "github.com/kianooshaz/skeleton/services/risk/audit/service"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

// Container holds all the application dependencies.
type Container struct {
	Config          *Config
	DB              *sql.DB
	Logger          *slog.Logger
	RestServer      protocol.WebService
	UserService     userproto.UserService
	OrgService      orgproto.OrganizationService
	UsernameService usernameproto.UsernameService
	AuditService    auditproto.AuditService
	PasswordService passwordproto.PasswordService
}

// Config represents the application configuration.
type Config struct {
	ShutdownTimeout time.Duration    `yaml:"shutdown_timeout"`
	Logger          log.LoggerConfig `yaml:"logger"`
	RestServer      rest.Config      `yaml:"rest_server"`
	Database        postgres.Config  `yaml:"postgres"`
	Account         struct {
		Username usernameservice.Config `yaml:"username"`
	} `yaml:"account"`
	Authentication struct {
		Password passwordservice.Config `yaml:"password"`
	} `yaml:"authentication"`
	Risk struct {
		Audit auditservice.Config `yaml:"audit"`
	} `yaml:"risk"`
}

// Close gracefully shuts down all services and closes connections.
func (c *Container) Close() error {
	var errs []error

	// Shutdown audit service workers
	if c.AuditService != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		c.AuditService.Shutdown(ctx)
		cancel()
	}

	// Close database connection
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing database: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
