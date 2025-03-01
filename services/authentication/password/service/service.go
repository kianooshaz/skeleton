package service

import (
	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/log"
	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/services/authentication/password/service/stores/db"
)

type (
	Config struct {
		MinLength                 uint   `yaml:"min_length"`
		AllowCharacters           string `yaml:"allow_characters"`
		Cost                      int    `yaml:"cost"`
		CheckPasswordHistoryLimit int32  `yaml:"check_password_history_limit"`
	}

	PasswordService struct {
		config          Config
		commonPasswords map[string]bool
		logger          log.Logger
		_pdb            postgres.DB
		db              *db.Queries
	}
)

func New(
	config Config,
	commonPasswords []string,
	logger log.Logger,
	pdb postgres.DB,
	db *db.Queries,
) *PasswordService {

	var commonPasswordsMap = make(map[string]bool)
	for _, password := range commonPasswords {
		commonPasswordsMap[password] = true
	}

	return &PasswordService{
		config:          config,
		commonPasswords: commonPasswordsMap,
		logger:          logger,
		_pdb:            pdb,
		db:              db,
	}
}

func (s *PasswordService) NewWithTx(tx pgx.Tx) *PasswordService {
	return &PasswordService{
		config:          s.config,
		commonPasswords: s.commonPasswords,
		logger:          s.logger,
		_pdb:            s._pdb,
		db:              db.New(tx),
	}
}
