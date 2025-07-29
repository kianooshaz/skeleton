package auns

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
	"github.com/kianooshaz/skeleton/services/account/username/persistence"
	aunp "github.com/kianooshaz/skeleton/services/account/username/proto"
)

var UsernameService aunp.UsernameService = &Service{}

type (
	Config struct {
		MaxUserUsernamePerOrganization uint   `yaml:"max_user_username_per_organization"`
		MinLength                      uint   `yaml:"min_length"`
		MaxLength                      uint   `yaml:"max_length"`
		AllowCharacters                string `yaml:"allow_characters"`
	}

	Storer interface {
		Create(ctx context.Context, username aunp.Username) error
		Delete(ctx context.Context, id uuid.UUID) error
		Get(ctx context.Context, id uuid.UUID) (aunp.Username, error)
		ListWithSearch(ctx context.Context, req aunp.ListRequest) ([]aunp.Username, error)
		CountWithSearch(ctx context.Context, req aunp.ListRequest) (int64, error)

		ListByUserAndOrganization(ctx context.Context, req aunp.ListAssignedRequest) ([]aunp.Username, error)
		UpdateStatus(ctx context.Context, username aunp.Username) error
		Exist(ctx context.Context, username string) (bool, error)
		CountByAccount(ctx context.Context, accountID accprotocol.AccountID) (int64, error)
	}

	Service struct {
		config      Config
		logger      slog.Logger
		storage     Storer
		storageConn *sql.DB
	}
)

func init() {
	cfg, err := config.Load[Config]("account.username")
	if err != nil {
		panic(err)
	}

	UsernameService = &Service{
		config: cfg,
		logger: *slog.With(
			slog.Group("package_info",
				slog.String("module", "username"),
				slog.String("service", "account"),
			),
		),
		storage: &persistence.UsernameStorage{
			Conn: postgres.ConnectionPool,
		},
	}
}
