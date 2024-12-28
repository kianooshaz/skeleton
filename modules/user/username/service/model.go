package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/stat"
)

type Username struct {
	ID             string
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	Status         stat.Status
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (u Username) GetID() string {
	return u.ID
}

func (u Username) GetUserID() uuid.UUID {
	return u.UserID
}

func (u Username) GetOrganizationID() uuid.UUID {
	return u.OrganizationID
}

func (u Username) GetStatus() stat.Status {
	return u.Status
}

func (u Username) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u Username) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
