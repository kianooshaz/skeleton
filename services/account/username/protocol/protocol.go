package aup

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/stat"
	"github.com/kianooshaz/skeleton/foundation/types"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

type Username struct {
	ID             uuid.UUID          `json:"id" bson:"id"`
	Username       iunp.Username      `json:"username" bson:"username"`
	UserID         iup.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID iop.OrganizationID `json:"organization_id" bson:"organization_id"`
	Status         stat.Status        `json:"status" bson:"status"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type ListUsername struct {
	ID             uuid.UUID          `json:"id" bson:"id"`
	Username       iunp.Username      `json:"username" bson:"username"`
	UserID         iup.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID iop.OrganizationID `json:"organization_id" bson:"organization_id"`
	Primary        bool               `json:"primary" bson:"primary"`
	Locked         bool               `json:"locked" bson:"locked"`
	Blocked        bool               `json:"blocked" bson:"blocked"`
	Reserved       bool               `json:"reserved" bson:"reserved"`
}

type UsernameService interface {
	// Add creates a new username.
	Assign(ctx context.Context, req AssignRequest) error

	ListAssigned(ctx context.Context, req ListAssignedRequest) (ListAssignedResponse, error)

	// Get returns the username with the specified ID.
	Get(ctx context.Context, id uuid.UUID) (Username, error)

	// Search returns usernames with the specified search criteria.
	List(ctx context.Context, req ListRequest) (ListResponse, error)

	// BePrimary sets the username with the specified ID as the primary username.
	BePrimary(context.Context, uuid.UUID, uuid.UUID, string) error

	// Hidden hides the username with the specified ID.
	Hidden(ctx context.Context, id uuid.UUID) error

	// Unhidden makes the username with the specified ID visible.
	Unhidden(ctx context.Context, id uuid.UUID) error

	// Reserve reserves the username with the specified ID.
	Reserve(ctx context.Context, username iunp.Username) (Username, error)

	// Unreserve unreserves the username with the specified ID.
	Unreserve(ctx context.Context, id uuid.UUID) error
}

type ListAssignedResponse pagination.Response[ListUsername]

type ListRequest struct {
	Username       types.Nullable[iunp.Username]      `query:"username"`
	UserID         types.Nullable[iup.UserID]         `query:"user_id"`
	OrganizationID types.Nullable[iop.OrganizationID] `query:"organization_id"`
	Status         types.Nullable[stat.Status]        `query:"status"`
	pagination.Page
	order.OrderBy
}

type ListResponse pagination.Response[ListUsername]

type GetRequest struct {
	
}
