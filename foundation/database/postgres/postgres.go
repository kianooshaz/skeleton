package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
)

type Config struct {
	Port        int           `yaml:"port"         validate:"required"`
	Name        string        `yaml:"name"         validate:"required"`
	Host        string        `yaml:"host"         validate:"required"`
	User        string        `yaml:"user"         validate:"required"`
	Password    string        `yaml:"password"     validate:"required"`
	SSLMode     string        `yaml:"ssl_mode"     validate:"required"`
	PingTimeout time.Duration `yaml:"ping_timeout"`
}

// ConnectionPool is the global connection pool for backward compatibility
// TODO: Remove this after all services are refactored
var ConnectionPool *sql.DB

// NewConnection creates a new database connection
func NewConnection(cfg Config) (*sql.DB, error) {
	connectionPool, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	pingTimeout := cfg.PingTimeout
	if pingTimeout == 0 {
		pingTimeout = 10 * time.Second // default timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err = connectionPool.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	// Set the global variable for backward compatibility
	if ConnectionPool == nil {
		ConnectionPool = connectionPool
	}

	return connectionPool, nil
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
