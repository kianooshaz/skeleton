package service

import (
	"time"

	"github.com/google/uuid"
)

// type for hashed password
type hash string

type Password struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	PasswordHash hash
	CreatedAt    time.Time
}

func (p Password) GetID() uuid.UUID {
	return p.ID
}

func (p Password) GetUserID() uuid.UUID {
	return p.UserID
}

func (p Password) GetPasswordHash() string {
	return string(p.PasswordHash)
}

func (p Password) GetCreatedAt() time.Time {
	return p.CreatedAt
}
