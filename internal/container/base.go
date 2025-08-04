// Package container provides dependency injection containers using Google Wire.
package container

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	usernameservice "github.com/kianooshaz/skeleton/services/account/username/service"
	passwordservice "github.com/kianooshaz/skeleton/services/authentication/password/service"
	auditservice "github.com/kianooshaz/skeleton/services/risk/audit/service"
)

// IContainer represents a base interface for all containers.
type IContainer interface {
	// Close gracefully shuts down the container and its dependencies.
	Close() error
}

// BaseConfig represents the base configuration shared across all containers.
type BaseConfig struct {
	ShutdownTimeout time.Duration    `yaml:"shutdown_timeout"`
	Logger          log.LoggerConfig `yaml:"logger"`
	Database        postgres.Config  `yaml:"postgres"`
}

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

// Close gracefully shuts down the base container.
func (c *BaseContainer) Close() error {
	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}
