package log

import (
	"context"
	"io"
	"log/slog"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Level int

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

func New(w io.Writer, lvl Level) Logger {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true})
	handler.Enabled(context.Background(), slog.Level(lvl))
	return slog.New(handler)
}
