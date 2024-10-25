package usernamesrv

import (
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/service/usernamesrv/stores/usernamedb"
)

type (
	Config struct {
		MaxPerUser      uint   `yaml:"max_per_user"`
		MinLength       uint   `yaml:"min_length"`
		MaxLength       uint   `yaml:"max_length"`
		AllowCharacters string `yaml:"allow_characters"`
	}

	Service struct {
		config  Config
		db      postgres.DB
		queries *usernamedb.Queries
	}
)

func New(config Config, db postgres.DB) *Service {
	return &Service{
		config:  config,
		db:      db,
		queries: usernamedb.New(db),
	}
}

// NewTx implements protocol.ServiceUser.
func (m *Service) NewWithTx(tx pgx.Tx) *Service {
	return &Service{
		config:  m.config,
		db:      m.db,
		queries: usernamedb.New(tx),
	}
}
