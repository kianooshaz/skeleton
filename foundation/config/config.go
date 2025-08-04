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
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// k is the global koanf instance used for config loading and unmarshaling.
var k = koanf.New(".")

// loadConfig loads the configuration file from the specified path.
func loadConfig(path string) error {
	// Load the configuration file using koanf and YAML parser.
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return fmt.Errorf("error loading config from %s: %w", path, err)
	}
	return nil
}

// LoadConfigWithDefaults loads the configuration file using dependency injection patterns.
// The default path is "config.yaml", but can be overridden by the CONFIG_PATH environment variable.
// Returns a koanf instance that can be used for configuration loading.
func LoadConfigWithDefaults() (*koanf.Koanf, error) {
	path := "config.yaml"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		path = envPath
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config from %s: %w", path, err)
	}
	return k, nil
}

// LoadFromKoanf unmarshals and validates a config section from a provided koanf instance.
// The k parameter is the koanf instance to load from.
// The path argument specifies the config section to unmarshal.
// Returns an error if unmarshaling or validation fails.
func LoadFromKoanf[T any](k *koanf.Koanf, path string) (T, error) {
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

// Load unmarshals and validates a config section into the provided struct type T.
// This function uses the global koanf instance for backward compatibility.
// The path argument specifies the config section to unmarshal.
// Returns an error if unmarshaling or validation fails.
// Deprecated: Use LoadFromKoanf with dependency injection instead.
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

// LoadWithFile loads a config file and returns a typed configuration section.
// This is useful for dependency injection where you want to explicitly control config loading.
func LoadWithFile[T any](configPath, section string) (T, error) {
	var out T

	// Create a new koanf instance for this specific load
	localK := koanf.New(".")

	// Load the configuration file
	if err := localK.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return out, fmt.Errorf("error loading config from %s: %w", configPath, err)
	}

	// Unmarshal the config section
	err := localK.UnmarshalWithConf(section, &out, koanf.UnmarshalConf{Tag: "yaml"})
	if err != nil {
		return out, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate the unmarshaled struct
	if err := validator.New().Struct(out); err != nil {
		return out, fmt.Errorf("error validating config: %w", err)
	}

	return out, nil
}
