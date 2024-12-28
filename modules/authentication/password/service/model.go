package service

import (
	"time"

	"github.com/google/uuid"
)

type Password struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	PasswordHash string
	CreatedAt    time.Time
}

func (p Password) GetID() uuid.UUID {
	return p.ID
}

func (p Password) GetUserID() uuid.UUID {
	return p.UserID
}

func (p Password) GetPasswordHash() string {
	return p.PasswordHash
}

func (p Password) GetCreatedAt() time.Time {
	return p.CreatedAt
}
