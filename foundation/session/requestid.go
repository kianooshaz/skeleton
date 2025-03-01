package session

import "context"

type RequestIDKey struct{}

// SetRequestID stores the provided request ID in the context.
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey{}, requestID)
}

// GetRequestID retrieves the request ID stored in the context.
func GetRequestID(ctx context.Context) string {
	requestID := ctx.Value(RequestIDKey{})
	if requestID == nil {
		return ""
	}

	return requestID.(string)
}
