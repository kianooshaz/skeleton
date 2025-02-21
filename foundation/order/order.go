package order

import (
	"errors"
	"strings"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

var ErrUnknownOrder = errors.New("unknown order")
var ErrUnknownDirection = errors.New("unknown direction")

type OrderBy struct {
	Field     string
	Direction string
}

func NewOrderBy(field string, direction string) OrderBy {
	if direction != ASC && direction != DESC {
		direction = ASC // Default to ASC if the direction is invalid
	}

	return OrderBy{
		Field:     field,
		Direction: direction,
	}
}

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
	case 1:
		return NewOrderBy(fieldName, ASC), nil

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

func MustParse(fieldMappings map[string]string, orderBy string, defaultOrder OrderBy) OrderBy {
	by, err := Parse(fieldMappings, orderBy, defaultOrder)
	if err != nil {
		panic(err)
	}

	return by
}

func (o OrderBy) PGX() string {
	return o.Field + "_" + o.Direction
}
