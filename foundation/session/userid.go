package session

import (
	"context"
)

type userIDKey struct{}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey{}).(int)
	return userID, ok
}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}
