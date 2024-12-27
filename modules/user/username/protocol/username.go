package protocol

import (
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/status"
)

type Username interface {
	GetID() uuid.UUID
	GetUserID() uuid.UUID
	GetValue() string
	GetStatus() status.Status
	IsPrimary() bool
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
