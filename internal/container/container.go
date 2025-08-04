package container

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
)

// AppConfig represents the root application configuration.
type AppConfig struct {
	ShutdownTimeout time.Duration    `yaml:"shutdown_timeout"`
	Logger          log.LoggerConfig `yaml:"logger"`
	RestServer      rest.Config      `yaml:"rest_server"`
	Postgres        postgres.Config  `yaml:"postgres"`
}

// WebContainer holds all dependencies for the web application.
type WebContainer struct {
	Config     *AppConfig
	Logger     *slog.Logger
	DB         *sql.DB
	WebService protocol.WebService
}

// Start initializes and starts all services.
func (c *WebContainer) Start() error {
	return c.WebService.Start()
}

// Stop gracefully shuts down all services.
func (c *WebContainer) Stop(ctx context.Context) error {
	return c.WebService.Shutdown(ctx)
}

// Close closes all resources.
func (c *WebContainer) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
