package container

import (
	"time"

	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
	usernameservice "github.com/kianooshaz/skeleton/services/account/username/service"
	passwordservice "github.com/kianooshaz/skeleton/services/authentication/password/service"
	auditservice "github.com/kianooshaz/skeleton/services/risk/audit/service"
	birthdayservice "github.com/kianooshaz/skeleton/services/user/birthday/service"
)

// AppConfig represents the root application configuration.
type AppConfig struct {
	ShutdownTimeout time.Duration          `yaml:"shutdown_timeout"`
	Logger          log.LoggerConfig       `yaml:"logger"`
	RestServer      rest.Config            `yaml:"rest_server"`
	Postgres        postgres.Config        `yaml:"postgres"`
	Password        passwordservice.Config `yaml:"password"`
	Username        usernameservice.Config `yaml:"username"`
	Audit           auditservice.Config    `yaml:"audit"`
	Birthday        birthdayservice.Config `yaml:"birthday"`
}
