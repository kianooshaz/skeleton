package session

import (
	"context"

	"github.com/kianooshaz/skeleton/foundation/database"
)

type queryExecutorKey struct{}

// GetQueryExecutor retrieves the query executor stored in the context.
// If no query executor is found in the context, it returns the provided fallback value.
func GetQueryExecutor(ctx context.Context, fallback database.QueryExecutorProtocol) database.QueryExecutorProtocol {
	tx := ctx.Value(queryExecutorKey{})
	if tx == nil {
		return fallback
	}
	return tx.(database.QueryExecutorProtocol)
}

// SetQueryExecutor stores the provided query executor in the context.
func SetQueryExecutor(ctx context.Context, tx database.QueryExecutorProtocol) context.Context {
	return context.WithValue(ctx, queryExecutorKey{}, tx)
}
