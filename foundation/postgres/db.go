package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Name     string `yaml:"name" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" env:"ROSTAM_POSTGRES_PASSWORD" validate:"required"`
}

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
)

func NewPool(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	var err error

	poolOnce.Do(func() {
		p, errPgx := pgxpool.New(ctx, dsn(cfg))
		if errPgx != nil {
			err = fmt.Errorf("unable to create connection pool: %w", errPgx)

			return
		}

		if errPgx = pool.Ping(ctx); err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", errPgx)

			return
		}

		pool = p
	})

	if err != nil {
		return nil, err
	}

	return pool, nil
}

func dsn(cfg *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}