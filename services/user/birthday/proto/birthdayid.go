package birthdayproto

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// BirthdayID represents a unique identifier for a birthday record.
type BirthdayID struct {
	uuid.UUID
}

// NewBirthdayID creates a new BirthdayID.
func NewBirthdayID() BirthdayID {
	return BirthdayID{UUID: uuid.New()}
}

// ParseBirthdayID parses a string into a BirthdayID.
func ParseBirthdayID(s string) (BirthdayID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return BirthdayID{}, fmt.Errorf("invalid birthday ID format: %w", err)
	}
	return BirthdayID{UUID: u}, nil
}

// MustParseBirthdayID parses a string into a BirthdayID and panics on error.
func MustParseBirthdayID(s string) BirthdayID {
	id, err := ParseBirthdayID(s)
	if err != nil {
		panic(err)
	}
	return id
}

// String returns the string representation of the BirthdayID.
func (id BirthdayID) String() string {
	return id.UUID.String()
}

// Value implements the driver.Valuer interface for database operations.
func (id BirthdayID) Value() (driver.Value, error) {
	return id.UUID.String(), nil
}

// Scan implements the sql.Scanner interface for database operations.
func (id *BirthdayID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("cannot scan %v into BirthdayID: %w", value, err)
		}
		id.UUID = u
	case []byte:
		u, err := uuid.ParseBytes(v)
		if err != nil {
			return fmt.Errorf("cannot scan %v into BirthdayID: %w", value, err)
		}
		id.UUID = u
	default:
		return fmt.Errorf("cannot scan %T into BirthdayID", value)
	}

	return nil
}

// IsZero returns true if the BirthdayID is zero (empty).
func (id BirthdayID) IsZero() bool {
	return id.UUID == uuid.Nil
}
