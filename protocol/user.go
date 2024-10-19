package protocol

import (
	"context"
	"time"

	"github.com/kianooshaz/skeleton/foundation/order"
)

type User interface {
	ID() ID
	CreatedAt() time.Time
}

type ServiceUser interface {
	// New creates a new user and returns the created user.
	New(ctx context.Context) (User, error)
	// Get retrieves a user by ID or returns ErrNotFound if the user does not exist.
	Get(ctx context.Context, id ID) (User, error)
	// List retrieves a list of users, ordered by the specified criteria.
	List(ctx context.Context, orderBy order.By) ([]User, error)
	// Count returns the total number of users.
	Count(ctx context.Context) (int64, error)
}
