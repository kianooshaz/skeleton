package service

import (
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/modules/user/username/service/stores/db"
)

type (
	Config struct {
		MaxPerUser      uint   `yaml:"max_per_user"`
		MinLength       uint   `yaml:"min_length"`
		MaxLength       uint   `yaml:"max_length"`
		AllowCharacters string `yaml:"allow_characters"`
	}

	Service struct {
		config Config
		logger log.Logger
		_pdb   postgres.DB
		db     *db.Queries
	}
)

func New(config Config, logger log.Logger, pdb postgres.DB) *Service {
	return &Service{
		config: config,
		logger: logger,
		_pdb:   pdb,
		db:     db.New(pdb),
	}
}

// NewTx implements protocol.ServiceUser.
func (m *Service) NewWithTx(tx pgx.Tx) *Service {
	return &Service{
		config: m.config,
		_pdb:   m._pdb,
		db:     db.New(tx),
	}
}
