package accproto

import (
	"time"

	"github.com/google/uuid"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

type AccountID uuid.UUID

func (a AccountID) String() string {
	return uuid.UUID(a).String()
}

type Account struct {
	ID             AccountID               `json:"id"`
	UserID         userproto.UserID        `json:"user_id"`
	OrganizationID orgproto.OrganizationID `json:"organization_id"`
	CreatedAt      time.Time               `json:"created_at"`
}
