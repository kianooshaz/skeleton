package protocol

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kianooshaz/skeleton/foundation/order"
)

// User represents a user entity with a unique ID and creation timestamp.
type User interface {
	// ID returns the unique identifier for the user.
	ID() ID
	// CreatedAt returns the timestamp indicating when the user was created.
	CreatedAt() time.Time
}

// ServiceUser defines the behavior for managing users within a service context.
type ServiceUser interface {
	// NewWithTx creates a new instance of the service with the provided transaction.
	NewWithTx(tx pgx.Tx) ServiceUser

	// New creates a new user and returns the created user or an error.
	// Context allows for timeout and cancellation control.
	New(ctx context.Context) (User, error)

	// Get retrieves a user by their unique ID.
	// If the user does not exist, ErrNotFound is returned.
	Get(ctx context.Context, id ID) (User, error)

	// List retrieves a list of users ordered by the specified criteria.
	// The order parameter determines how the list is sorted.
	List(ctx context.Context, orderBy order.By) ([]User, error)

	// Count returns the total number of users within the system.
	Count(ctx context.Context) (int64, error)
}
