// Package container provides dependency injection container using Google Wire.
//
//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build wireinject
// +build wireinject

package container

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/google/wire"
	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/protocol"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	usernameservice "github.com/kianooshaz/skeleton/services/account/username/service"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	passwordservice "github.com/kianooshaz/skeleton/services/authentication/password/service"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	orgservice "github.com/kianooshaz/skeleton/services/organization/organization/service"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	auditservice "github.com/kianooshaz/skeleton/services/risk/audit/service"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
	userservice "github.com/kianooshaz/skeleton/services/user/user/service"
)

// ProvideConfig loads the application configuration using dependency injection.
func ProvideConfig() (*Config, error) {
	configPath := "config.yml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadWithFile[Config](configPath, "app")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ProvideLogger creates a structured logger using dependency injection.
func ProvideLogger(cfg *Config) *slog.Logger {
	return log.NewLogger(cfg.Logger)
}

// ProvideDatabase creates a database connection.
func ProvideDatabase(cfg *Config) (*sql.DB, error) {
	return postgres.NewConnection(cfg.Database)
}

// ProvideRestServer creates the REST server.
func ProvideRestServer(cfg *Config, logger *slog.Logger) (protocol.WebService, error) {
	return rest.New(cfg.RestServer, logger)
}

// ProvideUserService creates the user service.
func ProvideUserService(db *sql.DB, logger *slog.Logger) userproto.UserService {
	return userservice.New(db, logger)
}

// ProvideOrgService creates the organization service.
func ProvideOrgService(db *sql.DB, logger *slog.Logger) orgproto.OrganizationService {
	return orgservice.New(db, logger)
}

// ProvideUsernameService creates the username service.
func ProvideUsernameService(cfg *Config, db *sql.DB, logger *slog.Logger) usernameproto.UsernameService {
	return usernameservice.New(cfg.Account.Username, db, logger)
}

// ProvideAuditService creates the audit service.
func ProvideAuditService(cfg *Config, db *sql.DB, logger *slog.Logger) auditproto.AuditService {
	return auditservice.New(cfg.Risk.Audit, db, logger)
}

// ProvidePasswordService creates a password service instance.
func ProvidePasswordService(cfg *Config, db *sql.DB, logger *slog.Logger) passwordproto.PasswordService {
	return passwordservice.New(cfg.Authentication.Password, db, logger)
}

// Wire set for all providers.
var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideDatabase,
	ProvideRestServer,
	ProvideUserService,
	ProvideOrgService,
	ProvideUsernameService,
	ProvideAuditService,
	ProvidePasswordService,
)

// NewContainer creates a new dependency injection container.
func NewContainer() (*Container, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Container), "*"),
	)
	return &Container{}, nil
}
