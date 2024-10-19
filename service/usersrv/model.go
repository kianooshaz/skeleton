package usersrv

import (
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/protocol"
)

var _ protocol.User = userModel{}

type userModel struct {
	id        uuid.UUID
	createdAt time.Time
}

func (u userModel) ID() protocol.ID {
	return [16]byte(u.id)
}

func (u userModel) CreatedAt() time.Time {
	return u.createdAt
}
