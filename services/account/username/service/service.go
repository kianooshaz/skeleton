package service

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/database/postgres"
	aup "github.com/kianooshaz/skeleton/services/account/username/protocol"
	"github.com/kianooshaz/skeleton/services/account/username/service/storage"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

var UsernameService aup.UsernameService = &Service{}

type (
	Config struct {
		MaxPerUser         uint   `yaml:"max_per_user"`
		MaxPerOrganization uint   `yaml:"max_per_user"`
		MinLength          uint   `yaml:"min_length"`
		MaxLength          uint   `yaml:"max_length"`
		AllowCharacters    string `yaml:"allow_characters"`
	}

	Storage interface {
		Create(ctx context.Context, username aup.Username) error
		Delete(ctx context.Context, id uuid.UUID) error
		Get(ctx context.Context, id uuid.UUID) (aup.Username, error)
		ListWithSearch(ctx context.Context, req aup.ListRequest) ([]aup.Username, error)
		CountWithSearch(ctx context.Context, req aup.ListRequest) (int64, error)

		ListByUserAndOrganization(ctx context.Context, userID iup.UserID, organizationID iop.OrganizationID) ([]aup.Username, error)
		UpdateStatus(ctx context.Context, username aup.Username) error
		Count(ctx context.Context, username iunp.Username) (int64, error)
		CountByUser(ctx context.Context, userID iup.UserID) (int64, error)
		CountByUserAndOrganization(ctx context.Context, userID iup.UserID, organization iop.OrganizationID) (int64, error)
	}

	Service struct {
		config      Config
		logger      slog.Logger
		storage     Storage
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
		storage: &storage.UsernameStorage{
			Conn: postgres.ConnectionPool,
		},
	}
}
