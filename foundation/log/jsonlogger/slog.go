package jsonlogger

import (
	"context"
	"io"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/log"
)

func New(w io.Writer, lvl log.Level) log.Protocol {
	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true})
	handler.Enabled(context.Background(), slog.Level(lvl))
	return slog.New(handler)
}
