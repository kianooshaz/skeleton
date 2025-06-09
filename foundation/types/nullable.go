package types

import (
	"database/sql"
	"fmt"
	"log/slog"
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
