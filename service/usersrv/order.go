package usersrv

import "github.com/kianooshaz/skeleton/foundation/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByCreatedAt, order.DESC)

// Set of fields that the results can be ordered by.

const (
	OrderByCreatedAt = "created_at"
)
