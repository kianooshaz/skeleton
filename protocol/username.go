// Package protocol defines interfaces and types for handling usernames and their management within the system.
package protocol

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/protocol/status"
)

// Username represents a username entity with a unique ID, value, and associated status.
type Username interface {
	// ID returns the unique identifier for the username.
	ID() ID

	// Value returns the string representation of the username.
	Value() string

	// Status returns the current status of the username, such as active, inactive, or pending.
	Status() status.Status

	// Primary indicates whether this username is the primary (or main) identifier for the user.
	// A primary username is typically the default or preferred username for the user.
	Primary() bool

	// UpdatedAt returns the timestamp indicating the last time the username was modified.
	UpdatedAt() time.Time

	// CreatedAt returns the timestamp indicating when the username was first created.
	CreatedAt() time.Time
}

type ServiceUsername interface {
	// NewWithTx creates a new instance of the service with the provided transaction.
	NewWithTx(tx pgx.Tx) ServiceUsername

	// Add assigns a new username value for the given user ID.
	Add(ctx context.Context, userID ID, value string) (Username, error)

	// Update modifies the value of an existing username identified by its unique ID.
	Update(ctx context.Context, id ID, value string) error

	// Search looks up a username by its value and returns the corresponding Username entity.
	// Returns an error if the username does not exist.
	Search(ctx context.Context, value string) (Username, error)

	// Count returns the number of usernames associated with a specific user ID.
	Count(ctx context.Context, userID ID) (int64, error)

	// List retrieves all usernames associated with a specific user ID.
	List(ctx context.Context, userID ID) []Username

	// Reserve marks a username as reserved, preventing others from claiming it.
	Reserve(ctx context.Context, value string) error

	// Unreserve removes the reservation status of a username, making it available for others to use.
	Unreserve(ctx context.Context, value string) error

	// Hidden sets a username to a hidden state, making it invisible to regular users.
	Hidden(ctx context.Context, id ID) error

	// Unhidden restores the visibility of a hidden username, making it available to regular users again.
	Unhidden(ctx context.Context, id ID) error

	// BePrimary sets the username identified by the given ID as the primary username for the user.
	// This method ensures the username becomes the user's main or default identifier.
	BePrimary(ctx context.Context, id ID) error
}
