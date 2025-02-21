package log

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
)

type config struct {
	AppEnv string `env:"APP_ENV" envDefault:"development"`
}

func Init() {
	cfg, err := env.ParseAs[config]()
	if err != nil {
		slog.Error("failed to parse env", slog.Any("error", err))
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
}
