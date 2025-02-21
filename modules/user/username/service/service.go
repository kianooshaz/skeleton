package service

import (
	"database/sql"

	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/modules/user/username/protocol"
)

type (
	Config struct {
		MaxPerUser         uint   `yaml:"max_per_user"`
		MaxPerOrganization uint   `yaml:"max_per_user"`
		MinLength          uint   `yaml:"min_length"`
		MaxLength          uint   `yaml:"max_length"`
		AllowCharacters    string `yaml:"allow_characters"`
	}

	Service struct {
		config Config
		logger log.Logger
		db     *sql.DB
	}
)

func New(config Config, logger log.Logger, db *sql.DB) protocol.UsernameService {
	return &Service{
		config: config,
		logger: logger,
		db:     db,
	}
}
