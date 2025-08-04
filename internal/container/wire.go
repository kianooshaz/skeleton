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
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	usernameservice "github.com/kianooshaz/skeleton/services/account/username/service"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	passwordservice "github.com/kianooshaz/skeleton/services/authentication/password/service"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	orgservice "github.com/kianooshaz/skeleton/services/organization/organization/service"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	auditservice "github.com/kianooshaz/skeleton/services/risk/audit/service"
	birthdayproto "github.com/kianooshaz/skeleton/services/user/birthday/proto"
	birthdayservice "github.com/kianooshaz/skeleton/services/user/birthday/service"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
	userservice "github.com/kianooshaz/skeleton/services/user/user/service"
)

// Simple config extraction functions - Wire can use these automatically.
func ProvideAppConfig(k *koanf.Koanf) (*AppConfig, error) {
	cfg, err := config.LoadFromKoanf[AppConfig](k, "app")
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ProvidePasswordConfig(cfg *AppConfig) passwordservice.Config { return cfg.Password }
func ProvideUsernameConfig(cfg *AppConfig) usernameservice.Config { return cfg.Username }
func ProvideAuditConfig(cfg *AppConfig) auditservice.Config       { return cfg.Audit }
func ProvideBirthdayConfig(cfg *AppConfig) birthdayservice.Config { return cfg.Birthday }
func ProvideRestConfig(cfg *AppConfig) rest.Config                { return cfg.RestServer }
func ProvideLoggerConfig(cfg *AppConfig) log.LoggerConfig         { return cfg.Logger }
func ProvidePostgresConfig(cfg *AppConfig) postgres.Config        { return cfg.Postgres }

// ProvideWebContainer provides the complete web container.
func ProvideWebContainer(
	cfg *AppConfig,
	logger *slog.Logger,
	db *sql.DB,
	webService protocol.WebService,
	userService userproto.UserService,
	orgService orgproto.OrganizationService,
	passwordService passwordproto.PasswordService,
	usernameService usernameproto.UsernameService,
	auditService auditproto.AuditService,
	birthdayService birthdayproto.BirthdayService,
) Container {
	return &WebContainer{
		config:              cfg,
		logger:              logger,
		db:                  db,
		webService:          webService,
		userService:         userService,
		organizationService: orgService,
		passwordService:     passwordService,
		usernameService:     usernameService,
		auditService:        auditService,
		birthdayService:     birthdayService,
	}
}

// Wire sets define the dependency injection graph.
var ConfigSet = wire.NewSet(
	config.LoadConfigWithDefaults,
	ProvideAppConfig,
	ProvidePasswordConfig,
	ProvideUsernameConfig,
	ProvideAuditConfig,
	ProvideBirthdayConfig,
	ProvideRestConfig,
	ProvideLoggerConfig,
	ProvidePostgresConfig,
)

var LoggerSet = wire.NewSet(
	log.NewLogger,
)

var DatabaseSet = wire.NewSet(
	postgres.NewConnection,
)

var WebContainerSet = wire.NewSet(
	ConfigSet,
	LoggerSet,
	DatabaseSet,
	userservice.New,
	orgservice.New,
	passwordservice.New,
	usernameservice.New,
	auditservice.New,
	birthdayservice.New,
	rest.New,
	ProvideWebContainer,
)

// NewWebContainer creates a new web container with all dependencies wired.
func NewWebContainer() (Container, error) {
	wire.Build(WebContainerSet)
	return nil, nil
}
