package session

import (
	"context"
)

type userIDKey struct{}

// GetUserID retrieves the user ID stored in the context.
// It returns the user ID and a boolean indicating whether the user ID was found.
func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey{}).(int)
	return userID, ok
}

// SetUserID stores the provided user ID in the context.
func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}
