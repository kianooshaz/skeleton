// Package config provides configuration loading and validation utilities for the application.
//
// It loads configuration from a YAML file (default: config.yaml),
// with the path optionally overridden by the CONFIG_PATH environment variable.
// The configuration is validated using struct tags and go-playground/validator.
//
// Usage:
//   - The configuration is loaded automatically at package initialization.
//   - Use Load[T any](path string) to unmarshal and validate a config section into a struct.
//
// Example:
//
//	type MyConfig struct {
//	    Name string `yaml:"name" validate:"required"`
//	}
//	cfg, err := config.Load[MyConfig]("parentKey")
//
// Errors in loading or validation will cause the application to log.Fatal or return an error.
package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// k is the global koanf instance used for config loading and unmarshaling.
var k = koanf.New(".")

// init loads the configuration file at package initialization.
// The default path is "config.yaml", but can be overridden by the CONFIG_PATH environment variable.
func init() {
	path := "config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		path = envPath
	}

	// Load the configuration file using koanf and YAML parser.
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		slog.Error("error loading config", "error", err)
		os.Exit(1)
	}
}

// Load unmarshals and validates a config section into the provided struct type T.
// The path argument specifies the config section to unmarshal.
// Returns an error if unmarshaling or validation fails.
func Load[T any](path string) (T, error) {
	var out T
	// Unmarshal the config section at the given path into out, using YAML struct tags.
	err := k.UnmarshalWithConf(path, &out, koanf.UnmarshalConf{Tag: "yaml"})
	if err != nil {
		return out, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate the unmarshaled struct using go-playground/validator.
	if err := validator.New().Struct(out); err != nil {
		return out, fmt.Errorf("error validating config: %w", err)
	}

	return out, nil
}
