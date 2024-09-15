package order

import (
	"errors"
	"strings"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

var ErrUnkownOrder = errors.New("unknown order")
var ErrUnkownDirection = errors.New("unknown direction")

type By struct {
	Field     string
	Direction string
}

func NewBy(field string, direction string) By {
	if _, exists := directions[direction]; !exists {
		return By{
			Field:     field,
			Direction: ASC,
		}
	}

	return By{
		Field:     field,
		Direction: direction,
	}
}

func Parse(fieldMappings map[string]string, orderBy string, defaultOrder By) (By, error) {
	if orderBy == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(orderBy, ",")

	orgFieldName := strings.TrimSpace(orderParts[0])
	fieldName, exists := fieldMappings[orgFieldName]
	if !exists {
		return By{}, ErrUnkownOrder
	}

	switch len(orderParts) {
	case 1:
		return NewBy(fieldName, ASC), nil

	case 2:
		direction := strings.TrimSpace(orderParts[1])
		if _, exists := directions[direction]; !exists {
			return By{}, ErrUnkownDirection
		}

		return NewBy(fieldName, direction), nil

	default:
		return By{}, ErrUnkownOrder
	}
}

func MustParse(fieldMappings map[string]string, orderBy string, defaultOrder By) By {
	by, err := Parse(fieldMappings, orderBy, defaultOrder)
	if err != nil {
		panic(err)
	}
	return by
}

func (b By) PGX() string {
	return b.Field + "_" + b.Direction
}
