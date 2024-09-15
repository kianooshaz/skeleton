package usersrv

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func (u User) GetID() uuid.UUID {
	return u.ID
}

func (u User) GetCreatedAt() time.Time {
	return u.CreatedAt
}
