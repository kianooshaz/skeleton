package usernamesrv

import (
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/protocol"
	"github.com/kianooshaz/skeleton/protocol/status"
)

var _ protocol.Username = usernameModel{}

type usernameModel struct {
	id        uuid.UUID
	userID    uuid.UUID
	value     string
	status    status.Status
	primary   bool
	createdAt time.Time
	updatedAt time.Time
}

// ID implements protocol.Username.
func (u usernameModel) ID() protocol.ID {
	return protocol.ID(u.id)
}

// Value implements protocol.Username.
func (u usernameModel) Value() string {
	return u.value
}

// Status implements protocol.Username.
func (u usernameModel) Status() status.Status {
	return u.status
}

// Primary implements protocol.Username.
func (u usernameModel) Primary() bool {
	return u.primary
}

// CreatedAt implements protocol.Username.
func (u usernameModel) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt implements protocol.Username.
func (u usernameModel) UpdatedAt() time.Time {
	return u.updatedAt
}
