package passwordservice

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/proto"
	"github.com/kianooshaz/skeleton/services/authentication/password/persistence"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
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

var _ passwordproto.PasswordService = (*Service)(nil)

// New creates a new password service instance.
func New(cfg Config, db *sql.DB, logger *slog.Logger) *Service {
	serviceLogger := *logger.With(
		slog.Group("package_info",
			slog.String("module", "password"),
			slog.String("service", "authentication"),
		),
	)

	// TODO load common passwords from assets
	// var commonPasswordsMap = make(map[string]bool)
	// for _, password := range commonPasswords {
	// 	commonPasswordsMap[password] = true
	// }

	return &Service{
		config: cfg,
		logger: serviceLogger,
		storage: &persistence.PasswordStorage{
			Conn: db,
		},
		storageConn: db,
	}
}
