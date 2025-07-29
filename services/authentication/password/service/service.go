package authpass

import (
	"database/sql"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	authpassp "github.com/kianooshaz/skeleton/services/authentication/password/protocol"
	"github.com/kianooshaz/skeleton/services/user/user/service/storage"
)

type (
	Config struct {
		MinLength                 uint     `yaml:"min_length"`
		AllowCharacters           string   `yaml:"allow_characters"`
		Cost                      int      `yaml:"cost"`
		CheckPasswordHistoryLimit int32    `yaml:"check_password_history_limit"`
		RequiredGuidelines        []string `yaml:"required_guidelines"`
		BetterHave                []string `yaml:"better_have"`
	}

	storer interface {
	}

	service struct {
		config          Config
		commonPasswords map[string]bool
		logger          *slog.Logger
		storage         storer
		storageConn     *sql.DB
	}
)

var Service authpassp.PasswordServices

func init() {
	cfg, err := config.Load[Config]("authentication.password")
	if err != nil {

	}

	// TODO load common passwords from assets

	// var commonPasswordsMap = make(map[string]bool)
	// for _, password := range commonPasswords {
	// 	commonPasswordsMap[password] = true
	// }

	s := &service{
		config: cfg,
		logger: slog.With(
			slog.Group("package_info",
				slog.String("module", "user"),
				slog.String("service", "user"),
			),
		),
		storage: &storage.UserStorage{
			Conn: postgres.ConnectionPool,
		},
		storageConn: postgres.ConnectionPool,
	}

	Service = s
}
