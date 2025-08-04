package log

import (
	"log/slog"
	"os"
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
