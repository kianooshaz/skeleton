package persistence

import "github.com/kianooshaz/skeleton/foundation/order"

var orderStringer order.StringerFunc = func(orderBy order.OrderBy) string {
	var field string
	switch orderBy.Field {
	case "id":
		field = "id"
	case "user_id":
		field = "user_id"
	case "date_of_birth":
		field = "date_of_birth"
	case "age":
		field = "age"
	case "created_at":
		field = "created_at"
	case "updated_at":
		field = "updated_at"
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

	return field + " " + direction
}
