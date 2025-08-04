package log

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
)

// LoggerConfig represents the logger configuration.
type LoggerConfig struct {
	Environment string `yaml:"environment" validate:"required,oneof=development production staging"`
	Level       string `yaml:"level" validate:"required,oneof=debug info warn error"`
	AddSource   bool   `yaml:"add_source"`
	Format      string `yaml:"format" validate:"required,oneof=json text"`
}

type config struct {
	AppEnv string `env:"APP_ENV" envDefault:"development"`
}

// NewLogger creates a new logger instance with proper configuration using dependency injection.
// The cfg parameter contains the logger configuration.
// Returns a configured logger instance.
func NewLogger(cfg LoggerConfig) *slog.Logger {
	// Set log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		AddSource: cfg.AddSource,
		Level:     level,
	}

	var handler slog.Handler
	if cfg.Format == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(&SessionHandler{handler})
}

// NewLoggerFromAppEnv creates a new logger instance based on app environment.
// Deprecated: Use NewLogger with LoggerConfig from dependency injection instead.
func NewLoggerFromAppEnv(appEnv string) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	if appEnv == "development" {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	}

	handler := &SessionHandler{slog.NewJSONHandler(os.Stdout, opts)}
	return slog.New(handler)
}

// NewLoggerFromEnv creates a new logger instance using environment variables.
// This is a convenience function that reads APP_ENV environment variable.
// Returns a configured logger instance.
// Deprecated: Use NewLogger with LoggerConfig from dependency injection instead.
func NewLoggerFromEnv() (*slog.Logger, error) {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return nil, err
	}
	return NewLoggerFromAppEnv(cfg.AppEnv), nil
}

// InitLogger initializes the default logger with proper configuration.
// This should be called explicitly during application startup.
// Deprecated: Use NewLogger or NewLoggerFromEnv with dependency injection instead.
func InitLogger() error {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		return err
	}

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	if cfg.AppEnv == "development" {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
	}

	handler := &SessionHandler{slog.NewJSONHandler(os.Stdout, opts)}

	slog.SetDefault(slog.New(handler))
	return nil
}
