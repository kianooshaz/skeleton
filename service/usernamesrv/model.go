package usernamesrv

import (
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/status"
)

type Username struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Value     string
	Status    status.Status
	Primary   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
