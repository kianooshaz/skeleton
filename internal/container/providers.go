package container

import (
	"database/sql"
	"log/slog"
	"os"

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

// ProvideAppConfig loads the application configuration using dependency injection.
func ProvideAppConfig() (*AppConfig, error) {
	configPath := "config.yml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadWithFile[AppConfig](configPath, "app")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ProvideBaseContainer creates a base container with common dependencies.
func ProvideBaseContainer(cfg *AppConfig, db *sql.DB, logger *slog.Logger) *BaseContainer {
	return &BaseContainer{
		Config: cfg,
		DB:     db,
		Logger: logger,
	}
}

// ProvideLogger creates a structured logger using dependency injection.
func ProvideLogger(cfg *AppConfig) *slog.Logger {
	return log.NewLogger(cfg.Logger)
}

// ProvideDatabase creates a database connection.
func ProvideDatabase(cfg *AppConfig) (*sql.DB, error) {
	return postgres.NewConnection(cfg.Database)
}

// ProvideRestServer creates a REST server instance.
func ProvideRestServer(cfg *AppConfig, logger *slog.Logger) (protocol.WebService, error) {
	return rest.New(cfg.RestServer, logger)
}

// ProvideUserService creates a user service instance.
func ProvideUserService(cfg *AppConfig, db *sql.DB, logger *slog.Logger) userproto.UserService {
	return userservice.New(db, logger)
}

// ProvideOrgService creates an organization service instance.
func ProvideOrgService(cfg *AppConfig, db *sql.DB, logger *slog.Logger) orgproto.OrganizationService {
	return orgservice.New(db, logger)
}

// ProvideUsernameService creates a username service instance.
func ProvideUsernameService(cfg *AppConfig, db *sql.DB, logger *slog.Logger) usernameproto.UsernameService {
	return usernameservice.New(cfg.Account.Username, db, logger)
}

// ProvideAuditService creates an audit service instance.
func ProvideAuditService(cfg *AppConfig, db *sql.DB, logger *slog.Logger) auditproto.AuditService {
	return auditservice.New(cfg.Risk.Audit, db, logger)
}

// ProvidePasswordService creates a password service instance.
func ProvidePasswordService(cfg *AppConfig, db *sql.DB, logger *slog.Logger) passwordproto.PasswordService {
	return passwordservice.New(cfg.Authentication.Password, db, logger)
}
