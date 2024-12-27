package service

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

func (u Username) GetID() uuid.UUID {
	return u.ID
}

func (u Username) GetUserID() uuid.UUID {
	return u.UserID
}

func (u Username) GetValue() string {
	return u.Value
}

func (u Username) GetStatus() status.Status {
	return u.Status
}

func (u Username) IsPrimary() bool {
	return u.Primary
}

func (u Username) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u Username) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
