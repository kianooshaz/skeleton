package types

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

// Nullable is a generic type that can represent a value that may be null (invalid).
// T is the type of the value being wrapped.
type Nullable[T any] struct {
	Value T    // The actual value
	Valid bool // Valid is true if Value is not null
}

// NewNullable creates a new Nullable with a valid value.
func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{
		Value: value,
		Valid: true,
	}
}

// IsValid returns true if the Nullable contains a valid (non-null) value.
func (n Nullable[T]) IsValid() bool {
	return n.Valid
}

// Get returns the value if valid, otherwise returns the zero value of T.
func (n Nullable[T]) Get() T {
	if !n.Valid {
		var zero T
		return zero
	}
	return n.Value
}

// Set assigns a new value and marks the Nullable as valid.
func (n *Nullable[T]) Set(value T) {
	n.Value = value
	n.Valid = true
}

func (n *Nullable[T]) UnmarshalParam(param string) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in UnmarshalParam",
				slog.String("param", param),
				slog.Any("recovered", r),
			)
			n.Valid = false
		}
	}()

	if len(param) == 0 || string(param) == "null" {
		n.Valid = false
		var zero T
		n.Value = zero
		return nil
	}

	var v T
	var err error

	switch any(v).(type) {
	case string:
		n.Value = any(param).(T)
	case int8:
		var parsed int64
		parsed, err = strconv.ParseInt(param, 10, 8)
		n.Value = any(int8(parsed)).(T)
	case int16:
		var parsed int64
		parsed, err = strconv.ParseInt(param, 10, 16)
		n.Value = any(int16(parsed)).(T)
	case int32:
		var parsed int64
		parsed, err = strconv.ParseInt(param, 10, 32)
		n.Value = any(int32(parsed)).(T)
	case int64:
		var parsed int64
		parsed, err = strconv.ParseInt(param, 10, 64)
		n.Value = any(parsed).(T)
	case int:
		var parsed int64
		parsed, err = strconv.ParseInt(param, 10, 0)
		n.Value = any(int(parsed)).(T)
	case uint8:
		var parsed uint64
		parsed, err = strconv.ParseUint(param, 10, 8)
		n.Value = any(uint8(parsed)).(T)
	case uint16:
		var parsed uint64
		parsed, err = strconv.ParseUint(param, 10, 16)
		n.Value = any(uint16(parsed)).(T)
	case uint32:
		var parsed uint64
		parsed, err = strconv.ParseUint(param, 10, 32)
		n.Value = any(uint32(parsed)).(T)
	case uint64:
		var parsed uint64
		parsed, err = strconv.ParseUint(param, 10, 64)
		n.Value = any(parsed).(T)
	case float32:
		var parsed float64
		parsed, err = strconv.ParseFloat(param, 32)
		n.Value = any(float32(parsed)).(T)
	case float64:
		var parsed float64
		parsed, err = strconv.ParseFloat(param, 64)
		n.Value = any(parsed).(T)
	case bool:
		var parsed bool
		parsed, err = strconv.ParseBool(param)
		n.Value = any(parsed).(T)
	case time.Time:
		var parsed time.Time
		parsed, err = time.Parse(time.RFC3339, param)
		n.Value = any(parsed).(T)
	default:
		return fmt.Errorf("unsupported type %T for Nullable", v)
	}

	if err != nil {
		n.Valid = false
		return fmt.Errorf("failed to parse param %q to %T: %w", param, v, err)
	}

	n.Valid = true
	return nil
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in UnmarshalJSON",
				slog.Any("data", string(data)),
				slog.Any("recovered", r),
			)
			n.Valid = false
		}
	}()

	if len(data) == 0 || string(data) == "null" {
		n.Valid = false
		var zero T
		n.Value = zero
		return nil
	}

	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		n.Valid = false
		return fmt.Errorf("failed to unmarshal JSON to %T: %w", v, err)
	}

	n.Value = v
	n.Valid = true
	return nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		// TODO : Handle null value in a way that fits your application
		return []byte{}, nil
	}

	return fmt.Appendf([]byte{}, "%v", n.Value), nil
}

// FromSQLNullInt64 converts a sql.NullInt64 to a Nullable[int64].
func FromSQLNullString[T any](value sql.NullString) Nullable[T] {
	if !value.Valid {
		var zero T
		return Nullable[T]{Value: zero, Valid: false}
	}

	var val T
	if v, ok := any(value.String).(T); ok {
		return Nullable[T]{Value: v, Valid: true}
	} else {
		slog.Error("Type assertion failed in FromSQLNullString",
			slog.String("value", value.String),
			slog.String("expected_type", fmt.Sprintf("%T", val)),
			slog.String("actual_type", "string"),
		)
		var zero T
		return Nullable[T]{Value: zero, Valid: false}
	}
}

// FromSQLNullTime converts a sql.NullTime to a Nullable[T].
func FromSQLNullTime[T any](value sql.NullTime) Nullable[T] {
	if !value.Valid {
		var zero T
		return Nullable[T]{Value: zero, Valid: false}
	}

	var val T
	if v, ok := any(value.Time).(T); ok {
		return Nullable[T]{Value: v, Valid: true}
	} else {
		slog.Error("Type assertion failed in FromSQLNullTime",
			slog.String("value", value.Time.String()),
			slog.String("expected_type", fmt.Sprintf("%T", val)),
			slog.String("actual_type", "time.Time"),
		)
		var zero T
		return Nullable[T]{Value: zero, Valid: false}
	}
}
