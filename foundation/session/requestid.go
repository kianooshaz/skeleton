package session

import (
	"context"

	"github.com/labstack/echo/v4"
)

var requestIDKey = "session_request_id"

// SetRequestID stores the provided request ID in the context.
func SetRequestIDEcho() func(c echo.Context, id string) {
	return func(c echo.Context, id string) {
		c.Set(requestIDKey, id)
	}
}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID retrieves the request ID stored in the context.
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}

	return ""
}
