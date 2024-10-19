package usernamesrv

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kianooshaz/skeleton/protocol"
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
		queries *usernamedb.Queries
	}
)

func New(config Config, pool *pgxpool.Pool) protocol.ServiceUsername {
	return &Service{
		config:  config,
		queries: usernamedb.New(pool),
	}
}

// NewTx implements protocol.ServiceUser.
func (m *Service) NewWithTx(tx pgx.Tx) protocol.ServiceUsername {
	return &Service{
		config:  m.config,
		queries: usernamedb.New(tx),
	}
}
