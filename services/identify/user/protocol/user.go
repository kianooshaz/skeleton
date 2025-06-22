package iup

import "github.com/google/uuid"

type UserID uuid.UUID

func (u UserID) String() string {
	return uuid.UUID(u).String()
}
