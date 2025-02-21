package session

import (
	"context"

	"github.com/kianooshaz/skeleton/foundation/database"
)

type dbConnectionKey struct{}

// GetDBConnection retrieves the db connection stored in the context.
// If no db connection is found in the context, it returns the provided fallback value.
func GetDBConnection(ctx context.Context, fallback database.ConnectionProtocol) database.ConnectionProtocol {
	tx := ctx.Value(dbConnectionKey{})
	if tx == nil {
		return fallback
	}
	return tx.(database.ConnectionProtocol)
}

// SetDBConnection stores the provided db connection in the context.
func SetDBConnection(ctx context.Context, tx database.ConnectionProtocol) context.Context {
	return context.WithValue(ctx, dbConnectionKey{}, tx)
}
