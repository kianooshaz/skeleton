package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

type Config struct {
	Name     string `yaml:"name" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

func New(ctx context.Context, cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn(cfg))
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	return db, nil
}

func dsn(cfg *Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}
