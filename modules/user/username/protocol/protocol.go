package protocol

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/stat"
)

// Username represents the interface for a username entity.
// It provides methods to access various properties of the username.
type Username interface {
	// GetID returns the username ID. It is unique
	GetID() string

	// GetUserID returns the user ID associated with the username
	GetUserID() uuid.UUID

	// GetOrganizationID returns the organization ID associated with the username
	GetOrganizationID() uuid.UUID

	// GetStatus returns the current status of the username
	GetStatus() stat.Status

	// GetCreatedAt returns the creation timestamp of the username. It is a required field.
	GetCreatedAt() time.Time

	// GetUpdatedAt returns the last updated timestamp of the username. It is a required field.
	GetUpdatedAt() time.Time
}

type UsernameService interface {
	// Add creates a new username.
	Add(ctx context.Context, userID, OrganizationID uuid.UUID, id string) (Username, error)

	// Get returns the username with the specified ID.
	Get(ctx context.Context, id string) (Username, error)

	// Search returns usernames with the specified search criteria.
	Search(ctx context.Context, userID *string, organizationID *string, status *stat.Status, limit int, offset int) ([]Username, error)

	Count(ctx context.Context, userID *string, organizationID *string, status *stat.Status) (int64, error)

	// BePrimary sets the username with the specified ID as the primary username.
	BePrimary(context.Context, uuid.UUID, uuid.UUID, string) error

	// Hidden hides the username with the specified ID.
	Hidden(ctx context.Context, id string) error

	// Unhidden makes the username with the specified ID visible.
	Unhidden(ctx context.Context, id string) error

	// Reserve reserves the username with the specified ID.
	Reserve(ctx context.Context, id string) (Username, error)

	// Unreserve unreserves the username with the specified ID.
	Unreserve(ctx context.Context, id string, userID, organizationID uuid.UUID) error
}
