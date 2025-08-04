//go:build wireinject
// +build wireinject

package container

import (
	"database/sql"
	"log/slog"

	"github.com/google/wire"
	"github.com/knadh/koanf/v2"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
)

// ProvideAppConfig provides the application configuration
func ProvideAppConfig(k *koanf.Koanf) (*AppConfig, error) {
	cfg, err := config.LoadFromKoanf[AppConfig](k, "app")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ProvideLogger provides a configured logger
func ProvideLogger(cfg *AppConfig) *slog.Logger {
	return log.NewLogger(cfg.Logger)
}

// ProvideDatabase provides a database connection
func ProvideDatabase(cfg *AppConfig) (*sql.DB, error) {
	return postgres.NewConnection(cfg.Postgres)
}

// ProvideWebService provides the REST web service
func ProvideWebService(cfg *AppConfig, logger *slog.Logger, passwordService passwordproto.PasswordService) (protocol.WebService, error) {
	return rest.New(cfg.RestServer, logger, passwordService)
}

// ProvideWebContainer provides the complete web container
func ProvideWebContainer(cfg *AppConfig, logger *slog.Logger, db *sql.DB, webService protocol.WebService) *WebContainer {
	return &WebContainer{
		Config:     cfg,
		Logger:     logger,
		DB:         db,
		WebService: webService,
	}
}

// Wire sets define the dependency injection graph
var ConfigSet = wire.NewSet(
	config.LoadConfigWithDefaults,
	ProvideAppConfig,
)

var LoggerSet = wire.NewSet(
	ProvideLogger,
)

var DatabaseSet = wire.NewSet(
	ProvideDatabase,
)

var WebServiceSet = wire.NewSet(
	ProvideWebService,
)

var WebContainerSet = wire.NewSet(
	ConfigSet,
	LoggerSet,
	DatabaseSet,
	WebServiceSet,
	ProvideWebContainer,
)

// NewWebContainer creates a new web container with all dependencies wired
func NewWebContainer() (*WebContainer, error) {
	wire.Build(WebContainerSet)
	return nil, nil
}
