package session

import (
	"context"

	"github.com/kianooshaz/skeleton/foundation/database"
)

type queryExecutorKey struct{}

func GetQueryExecutor(ctx context.Context, fallback database.QueryExecutorProtocol) database.QueryExecutorProtocol {
	tx := ctx.Value(queryExecutorKey{})
	if tx == nil {
		return fallback
	}
	return tx.(database.QueryExecutorProtocol)
}

func SetQueryExecutor(ctx context.Context, tx database.QueryExecutorProtocol) context.Context {
	return context.WithValue(ctx, queryExecutorKey{}, tx)
}
