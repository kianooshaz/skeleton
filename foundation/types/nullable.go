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

	value, err := n.parseParam(param)
	if err != nil {
		n.Valid = false
		return err
	}

	n.Value = value
	n.Valid = true
	return nil
}

func (n *Nullable[T]) parseParam(param string) (T, error) {
	var v T

	switch any(v).(type) {
	case string:
		return any(param).(T), nil
	case int8:
		return n.parseInt8(param)
	case int16:
		return n.parseInt16(param)
	case int32:
		return n.parseInt32(param)
	case int64:
		return n.parseInt64(param)
	case int:
		return n.parseInt(param)
	case uint8:
		return n.parseUint8(param)
	case uint16:
		return n.parseUint16(param)
	case uint32:
		return n.parseUint32(param)
	case uint64:
		return n.parseUint64(param)
	case float32:
		return n.parseFloat32(param)
	case float64:
		return n.parseFloat64(param)
	case bool:
		return n.parseBool(param)
	case time.Time:
		return n.parseTime(param)
	default:
		return v, fmt.Errorf("unsupported type %T for Nullable", v)
	}
}

func (n *Nullable[T]) parseInt8(param string) (T, error) {
	parsed, err := strconv.ParseInt(param, 10, 8)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to int8: %w", param, err)
	}
	return any(int8(parsed)).(T), nil
}

func (n *Nullable[T]) parseInt16(param string) (T, error) {
	parsed, err := strconv.ParseInt(param, 10, 16)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to int16: %w", param, err)
	}
	return any(int16(parsed)).(T), nil
}

func (n *Nullable[T]) parseInt32(param string) (T, error) {
	parsed, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to int32: %w", param, err)
	}
	return any(int32(parsed)).(T), nil
}

func (n *Nullable[T]) parseInt64(param string) (T, error) {
	parsed, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to int64: %w", param, err)
	}
	return any(parsed).(T), nil
}

func (n *Nullable[T]) parseInt(param string) (T, error) {
	parsed, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to int: %w", param, err)
	}
	return any(int(parsed)).(T), nil
}

func (n *Nullable[T]) parseUint8(param string) (T, error) {
	parsed, err := strconv.ParseUint(param, 10, 8)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to uint8: %w", param, err)
	}
	return any(uint8(parsed)).(T), nil
}

func (n *Nullable[T]) parseUint16(param string) (T, error) {
	parsed, err := strconv.ParseUint(param, 10, 16)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to uint16: %w", param, err)
	}
	return any(uint16(parsed)).(T), nil
}

func (n *Nullable[T]) parseUint32(param string) (T, error) {
	parsed, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to uint32: %w", param, err)
	}
	return any(uint32(parsed)).(T), nil
}

func (n *Nullable[T]) parseUint64(param string) (T, error) {
	parsed, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to uint64: %w", param, err)
	}
	return any(parsed).(T), nil
}

func (n *Nullable[T]) parseFloat32(param string) (T, error) {
	parsed, err := strconv.ParseFloat(param, 32)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to float32: %w", param, err)
	}
	return any(float32(parsed)).(T), nil
}

func (n *Nullable[T]) parseFloat64(param string) (T, error) {
	parsed, err := strconv.ParseFloat(param, 64)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to float64: %w", param, err)
	}
	return any(parsed).(T), nil
}

func (n *Nullable[T]) parseBool(param string) (T, error) {
	parsed, err := strconv.ParseBool(param)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to bool: %w", param, err)
	}
	return any(parsed).(T), nil
}

func (n *Nullable[T]) parseTime(param string) (T, error) {
	parsed, err := time.Parse(time.RFC3339, param)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse param %q to time.Time: %w", param, err)
	}
	return any(parsed).(T), nil
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
