package storage

import "github.com/kianooshaz/skeleton/foundation/order"

var oderStringer order.StringerFunc = func(orderBy order.OrderBy) string {
	var field string
	switch orderBy.Field {
	case "created_at":
		field = "created_at"
	case "organization_id":
		field = "organization_id"
	default:
		field = "created_at"
	}

	var direction string
	switch orderBy.Direction {
	case order.ASC:
		direction = "ASC"
	case order.DESC:
		direction = "DESC"
	default:
		direction = "ASC"
	}

	return field + " " + direction
}
