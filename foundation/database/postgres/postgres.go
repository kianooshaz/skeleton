package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kianooshaz/skeleton/foundation/config"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
)

type Config struct {
	Port        int           `yaml:"port" validate:"required"`
	Name        string        `yaml:"name" validate:"required"`
	Host        string        `yaml:"host" validate:"required"`
	User        string        `yaml:"user" validate:"required"`
	Password    string        `yaml:"password" validate:"required"`
	SSLMode     string        `yaml:"ssl_mode" validate:"required"`
	PingTimeout time.Duration `yaml:"ping_timeout"`
}

var ConnectionPool *sql.DB

func init() {
	cfg, err := config.Load[Config]("database.postgres")
	if err != nil {
		panic(fmt.Sprintf("error loading config: %v", err))
	}

	connectionPool, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		panic(fmt.Sprintf("error opening database connection: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()

	if err = connectionPool.PingContext(ctx); err != nil {
		panic(fmt.Sprintf("error pinging database: %v", err))
	}

	ConnectionPool = connectionPool
}

func dsn(cfg Config) string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
}
