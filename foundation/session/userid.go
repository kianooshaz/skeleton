package session

import (
	"context"

	"github.com/google/uuid"
)

type userIDKey struct{}

// GetUserID retrieves the user ID stored in the context.
// It returns the user ID and a boolean indicating whether the user ID was found.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(userIDKey{}).(uuid.UUID)
	return userID, ok
}

// SetUserID stores the provided user ID in the context.
func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}
