package session

import (
	"context"
	"database/sql"

	dp "github.com/kianooshaz/skeleton/foundation/database/protocol"
)

type dbConnectionKey struct{}

// GetDBConnection retrieves the db connection stored in the context.
// If no db connection is found in the context, it returns the provided fallback value.
func GetDBConnection(ctx context.Context, fallback dp.QueryExecutor) dp.QueryExecutor {
	tx := ctx.Value(dbConnectionKey{})
	if tx == nil {
		return fallback
	}
	return tx.(dp.QueryExecutor)
}

// SetDBConnection stores the provided db connection in the context.
func SetDBConnection(ctx context.Context, tx dp.QueryExecutor) context.Context {
	return context.WithValue(ctx, dbConnectionKey{}, tx)
}

func BeginTransaction(ctx context.Context, db *sql.DB) (*sql.Tx, context.Context, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	return tx, SetDBConnection(ctx, tx), nil
}
