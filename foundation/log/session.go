package log

import (
	"context"
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/session"
)

// SessionHandler is a custom slog.Handler that enriches log records with
// contextual information from the session. It wraps an existing handler
// and automatically adds session-specific attributes to all log entries.
type SessionHandler struct {
	slog.Handler
}

// Handle processes a log record by enriching it with session context before
// delegating to the underlying handler. It extracts session attributes and
// request ID from the context and adds them to the log record.
//
// The method performs the following operations:
//   - Retrieves all session log attributes from the context
//   - Adds each attribute to the log record
//   - Adds the request ID as a separate attribute
//   - Delegates to the wrapped handler for final processing
//
// Parameters:
//   - ctx: The context containing session information
//   - r: The log record to be enriched and processed
//
// Returns an error if the underlying handler fails to process the record.
func (h SessionHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add all session-specific log attributes to the record
	for _, v := range session.GetLogAttributes(ctx) {
		r.AddAttrs(v)
	}

	// Add the request ID as a dedicated attribute for request tracing
	r.Add(slog.String("request_id", session.GetRequestID(ctx)))

	// Delegate to the wrapped handler for final processing
	return h.Handler.Handle(ctx, r)
}
