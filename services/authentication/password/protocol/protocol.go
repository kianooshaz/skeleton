package protocol

import (
	"time"

	"github.com/google/uuid"
)

type Password interface {
	GetID() uuid.UUID
	GetUserID() uuid.UUID
	GetPasswordHash() string
	GetCreatedAt() time.Time
}

type PasswordService interface {
}
