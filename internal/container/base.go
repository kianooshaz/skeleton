// Package container provides dependency injection containers using Google Wire.
package container

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/wire"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	usernameservice "github.com/kianooshaz/skeleton/services/account/username/service"
	passwordservice "github.com/kianooshaz/skeleton/services/authentication/password/service"
	auditservice "github.com/kianooshaz/skeleton/services/risk/audit/service"
)

// BaseProviderSet contains providers for base container dependencies.
var BaseProviderSet = wire.NewSet(
	ProvideAppConfig,
	ProvideLogger,
	ProvideDatabase,
	ProvideBaseContainer,
)

// AppConfig represents the full application configuration.
type AppConfig struct {
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

// BaseContainer holds common dependencies shared across all containers.
type BaseContainer struct {
	Config *AppConfig
	DB     *sql.DB
	Logger *slog.Logger
}

// Start initializes the base container (currently no startup logic needed).
func (c *BaseContainer) Start() error {
	// Base container doesn't need specific startup logic
	// Database connections are established during provider initialization
	return nil
}

// Stop gracefully stops the base container (currently delegates to Close).
func (c *BaseContainer) Stop(ctx context.Context) error {
	// For base container, stop is equivalent to close
	return c.Close()
}

// Close gracefully shuts down the base container.
func (c *BaseContainer) Close() error {
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}
