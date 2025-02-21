package log

import (
	"context"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/session"
)

type SessionHandler struct {
	slog.Handler
}

// Handle adds contextual attributes to the Record before calling the underlying
// handler
func (h SessionHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, v := range session.GetLogAttributes(ctx) {
		r.AddAttrs(v)
	}

	return h.Handler.Handle(ctx, r)
}
