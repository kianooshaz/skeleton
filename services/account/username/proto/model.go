package usernameproto

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/stat"
	"github.com/kianooshaz/skeleton/foundation/types"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/protocol"
)

type Username struct {
	ID        uuid.UUID             `json:"id" bson:"id"`
	Username  string                `json:"username" bson:"username"`
	AccountID accprotocol.AccountID `json:"account_id" bson:"account_id"`
	Status    stat.Status           `json:"status" bson:"status"`
	CreatedAt time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" bson:"updated_at"`
}

type ListUsername struct {
	ID        uuid.UUID             `json:"id" bson:"id"`
	Username  string                `json:"username" bson:"username"`
	AccountID accprotocol.AccountID `json:"account_id" bson:"account_id"`
	Primary   bool                  `json:"primary" bson:"primary"`
	Locked    bool                  `json:"locked" bson:"locked"`
	Blocked   bool                  `json:"blocked" bson:"blocked"`
	Reserved  bool                  `json:"reserved" bson:"reserved"`
}

type UsernameService interface {
	// Add creates a new username.
	Assign(ctx context.Context, req AssignRequest) (Username, error)

	ListAssigned(ctx context.Context, req ListAssignedRequest) (ListAssignedResponse, error)

	// Get returns the username with the specified ID.
	Get(ctx context.Context, id uuid.UUID) (Username, error)

	// Search returns usernames with the specified search criteria.
	List(ctx context.Context, req ListRequest) (ListResponse, error)

	// BePrimary sets the username with the specified ID as the primary username.
	BePrimary(context.Context, uuid.UUID) error
}

type ListAssignedResponse pagination.Response[ListUsername]

type ListRequest struct {
	Username  types.Nullable[Username]              `query:"username"`
	AccountID types.Nullable[accprotocol.AccountID] `query:"account_id"`
	Status    types.Nullable[stat.Status]           `query:"status"`
	pagination.Page
	order.OrderBy
}

type ListResponse pagination.Response[ListUsername]

type GetRequest struct {
}
