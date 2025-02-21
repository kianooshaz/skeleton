package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kianooshaz/skeleton/foundation/config"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
)

type Config struct {
	Name     string `yaml:"name" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	SSLMode  string `yaml:"ssl_mode" validate:"required"`
}

var ConnectionPool *sql.DB

func Init(ctx context.Context) error {
	cfg, err := config.Load[Config]("postgres")
	if err != nil {
		return fmt.Errorf("unable to parse environment variables: %w", err)
	}

	connectionPool, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		return fmt.Errorf("unable to open database connection: %w", err)
	}

	if err = connectionPool.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	ConnectionPool = connectionPool

	return nil
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
