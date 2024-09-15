package protocol

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/order"
)

type User interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
}

type ServiceUser interface {
	New(ctx context.Context) (User, error)
	Get(ctx context.Context, id uuid.UUID) (User, error)
	List(ctx context.Context, orderBy order.By) ([]User, error)
	Count(ctx context.Context) (int64, error)
}
