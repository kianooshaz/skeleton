// Package order provides functionality for handling ordering criteria, including sorting fields and directions.
// It defines structures and functions to create, parse, and validate order specifications for sorting.
package order

import (
	"errors"
	"strings"
)

// Constants for sorting directions
const (
	ASC  = "ASC"  // Ascending order
	DESC = "DESC" // Descending order
)

// Error variables for handling unknown order and direction errors
var ErrUnknownOrder = errors.New("unknown order")
var ErrUnknownDirection = errors.New("unknown direction")

// OrderBy struct represents the sorting criteria with a field and direction
type OrderBy struct {
	Field     string // The field name to sort by
	Direction string // The direction of sorting (ASC or DESC)
}

// NewOrderBy creates a new OrderBy instance with a given field and direction.
// If the direction is invalid, it defaults to ASC.
func NewOrderBy(field string, direction string) OrderBy {
	if direction != ASC && direction != DESC {
		direction = ASC // Default to ASC if the direction is invalid
	}

	return OrderBy{
		Field:     field,
		Direction: direction,
	}
}

// Parse function parses a comma-separated order string and returns an OrderBy struct.
// If the order string is invalid, it returns an error.
// If no order string is provided, it returns the default order.
func Parse(fieldMappings map[string]string, orderBy string, defaultOrder OrderBy) (OrderBy, error) {
	if orderBy == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(orderBy, ",")

	orgFieldName := strings.TrimSpace(orderParts[0])

	fieldName, exists := fieldMappings[orgFieldName]
	if !exists {
		return OrderBy{}, ErrUnknownOrder
	}

	switch len(orderParts) {
	// Only the field is provided (default to ASC direction)
	case 1:
		return NewOrderBy(fieldName, ASC), nil

	// Both field and direction are provided
	case 2:
		direction := strings.TrimSpace(orderParts[1])
		if direction != ASC && direction != DESC {
			return OrderBy{}, ErrUnknownDirection
		}

		return NewOrderBy(fieldName, direction), nil

	default:
		return OrderBy{}, ErrUnknownOrder
	}
}

// MustParse is a helper function that calls Parse but panics if an error occurs.
// It ensures that parsing is successful, and if not, the program will terminate with the error.
func MustParse(fieldMappings map[string]string, orderBy string, defaultOrder OrderBy) OrderBy {
	by, err := Parse(fieldMappings, orderBy, defaultOrder)
	if err != nil {
		panic(err)
	}

	return by
}

// PGX method formats the OrderBy struct into a string suitable for database queries
// by concatenating the field and direction with an underscore.
func (o OrderBy) PGX() string {
	return o.Field + "_" + o.Direction
}
