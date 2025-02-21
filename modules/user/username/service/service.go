package service

import (
	"context"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/config"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/sqlservice"
	"github.com/kianooshaz/skeleton/foundation/storage/postgres"
	"github.com/kianooshaz/skeleton/foundation/types"
	"github.com/kianooshaz/skeleton/modules/user/username/protocol"
	"github.com/kianooshaz/skeleton/modules/user/username/service/storage"
)

var UsernameService protocol.UsernameService = &Service{}

type (
	Config struct {
		MaxPerUser         uint   `yaml:"max_per_user"`
		MaxPerOrganization uint   `yaml:"max_per_user"`
		MinLength          uint   `yaml:"min_length"`
		MaxLength          uint   `yaml:"max_length"`
		AllowCharacters    string `yaml:"allow_characters"`
	}

	Storage interface {
		Create(ctx context.Context, username protocol.Username) error
		Delete(ctx context.Context, id string) error
		Get(ctx context.Context, id string) (protocol.Username, error)
		ListByUser(ctx context.Context, userID types.UserID, orderBy order.OrderBy, page pagination.Page,
			isPrimary bool) ([]protocol.Username, error)
		ListByUserAndOrganization(ctx context.Context, userID types.UserID, organizationID types.OrganizationID,
			orderBy order.OrderBy, page pagination.Page, isPrimary bool) ([]protocol.Username, error)
		UpdateStatus(ctx context.Context, req protocol.Username) error
		Count(ctx context.Context, id string) (int64, error)
		CountByUser(ctx context.Context, userID types.UserID) (int64, error)
		CountByUserAndOrganization(ctx context.Context, userID types.UserID, organization types.OrganizationID) (int64, error)
	}

	Service struct {
		config  Config
		logger  slog.Logger
		storage Storage
	}
)

func Init() {
	cfg, err := config.Load[Config]("user.username")
	if err != nil {
		panic(err)
	}

	UsernameService = &Service{
		config: cfg,
		logger: *slog.With(
			slog.Group("package_info",
				slog.String("module", "user"),
				slog.String("service", "user"),
			),
		),
		storage: &storage.UsernameStorage{
			Conn: postgres.ConnectionPool,
		},
	}
}

// NewTx implements protocol.ServiceUser.
func (s *Service) NewWithTx(sqlConn sqlservice.ConnectionPool) protocol.UsernameService {
	return &Service{
		config: s.config,
		logger: s.logger,
		storage: &storage.UsernameStorage{
			Conn: sqlConn,
		},
	}
}
