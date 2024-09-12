package web

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

type Config struct {
	Version string `yaml:"version" env:"SKELETON_VERSION"`
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

	if err := validator.New().Struct(&cfg); err != nil {
		return nil, fmt.Errorf("error at validate config: %w", err)
	}

	return &cfg, nil
}
