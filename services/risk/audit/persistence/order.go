package persistence

import (
	"github.com/kianooshaz/skeleton/foundation/order"
)

var oderStringer order.StringerFunc = func(orderBy order.OrderBy) string {
	var field string
	switch orderBy.Field {
	case "created_at":
		field = "created_at"
	case "action":
		field = "action"
	case "user_id":
		field = "user_id"
	case "resource_type":
		field = "resource_type"
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
		direction = "DESC"
	}

	return " ORDER BY " + field + " " + direction
}
