package authpass

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
	"github.com/kianooshaz/skeleton/services/authentication/password/persistence"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
)

var PasswordService passwordproto.PasswordService = &Service{}

type (
	Config struct {
		MinLength                 uint     `yaml:"min_length"`
		AllowCharacters           string   `yaml:"allow_characters"`
		Cost                      int      `yaml:"cost"`
		CheckPasswordHistoryLimit int32    `yaml:"check_password_history_limit"`
		RequiredGuidelines        []string `yaml:"required_guidelines"`
		BetterHave                []string `yaml:"better_have"`
	}

	Storer interface {
		Create(ctx context.Context, password passwordproto.Password) error
		Delete(ctx context.Context, id uuid.UUID) error
		Get(ctx context.Context, id uuid.UUID) (passwordproto.Password, error)
		GetByAccountID(ctx context.Context, accountID accprotocol.AccountID) (passwordproto.Password, error)
		ListWithSearch(ctx context.Context, req passwordproto.ListRequest) ([]passwordproto.Password, error)
		CountWithSearch(ctx context.Context, req passwordproto.ListRequest) (int64, error)
		History(ctx context.Context, accountID accprotocol.AccountID, limit int32) ([]passwordproto.Password, error)
	}

	Service struct {
		config          Config
		commonPasswords map[string]bool
		logger          slog.Logger
		storage         Storer
		storageConn     *sql.DB
	}
)

func init() {
	cfg, err := config.Load[Config]("authentication.password")
	if err != nil {
		panic(err)
	}

	// TODO load common passwords from assets

	// var commonPasswordsMap = make(map[string]bool)
	// for _, password := range commonPasswords {
	// 	commonPasswordsMap[password] = true
	// }

	PasswordService = &Service{
		config: cfg,
		logger: *slog.With(
			slog.Group("package_info",
				slog.String("module", "password"),
				slog.String("service", "authentication"),
			),
		),
		storage: &persistence.PasswordStorage{
			Conn: postgres.ConnectionPool,
		},
	}
}
