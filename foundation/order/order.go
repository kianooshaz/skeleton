// Package order provides functionality for parsing and constructing order-by clauses.
// Initially designed for SQL, it allows extension to support other storage types
// by utilizing custom string builder functions.

package order

type Direction string

// Constants for sorting directions
const (
	// ASC represents ascending order.
	ASC Direction = "ascending"
	// DESC represents descending order.
	DESC Direction = "descending"
)

// directions maps valid direction strings for validation purposes.
var directions = map[Direction]bool{
	ASC:  true,
	DESC: false,
}

// OrderBy represents an order-by clause with a field and direction.
type OrderBy struct {
	Field     string    `json:"field" bson:"field"`
	Direction Direction `json:"direction" bson:"direction"`
}

// NewBy creates a new By instance with validation on the direction.
// If the direction is invalid, it defaults to ASC.
func NewBy(field string, direction Direction) OrderBy {
	if !directions[direction] {
		return OrderBy{
			Field:     field,
			Direction: ASC,
		}
	}

	return OrderBy{
		Field:     field,
		Direction: direction,
	}
}
