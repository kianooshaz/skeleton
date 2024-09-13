package web

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/kianooshaz/skeleton/foundation/postgres"
	"github.com/kianooshaz/skeleton/internal/app/web/rest"
)

type Config struct {
	Version         string           `yaml:"version" env:"SKELETON_VERSION" validate:"required"`
	ShutdownTimeout time.Duration    `yaml:"shutdown_timeout"`
	Rest            *rest.Config     `yaml:"rest" validate:"required"`
	Postgres        *postgres.Config `yaml:"postgres" validate:"required"`
}

func newConfig(configPath string) (*Config, error) {
	c := config.New()

	c.AddFeeder(feeder.Env{})
	if configPath != "" {
		c.AddFeeder(feeder.Yaml{Path: configPath})
	}

	var cfg Config
	c.AddStruct(&cfg)
	if err := c.Feed(); err != nil {
		return nil, fmt.Errorf("error at feed config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error at validate config: %w", err)
	}

	backupCfg := bytes.Buffer{}
	if err := gob.NewEncoder(&backupCfg).Encode(cfg); err != nil {
		return nil, fmt.Errorf("cant encode config for backup: %w", err)
	}

	c.SetupListener(func(err error) {
		if err != nil {
			log.Printf("error at setup listener: %v", err)
		}

		if err := gob.NewEncoder(&backupCfg).Encode(cfg); err != nil {
			fmt.Printf("cant encode config for backup: %v\n", err.Error())
		}
	})

	return &cfg, nil
}

func (c *Config) Validate() error {
	return validator.New().Struct(c)
}

func (c *Config) Setup() error {
	return nil
}
