package aup

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/stat"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

type Username struct {
	ID             int                `json:"id" bson:"id"`
	Identifier     iunp.Username      `json:"identifier" bson:"identifier"`
	UserID         iup.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID iop.OrganizationID `json:"organization_id" bson:"organization_id"`
	Status         stat.Status        `json:"status" bson:"status"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type AddRequest struct {
	UserID         iup.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID iop.OrganizationID `json:"organization_id" bson:"organization_id"`
	Identifier     iunp.Username      `json:"identifier" bson:"identifier"`
}

type ListRequest struct {
	Identifier     *iunp.Username `json:"identifier" bson:"identifier"`
	UserID         *string        `json:"user_id" bson:"user_id"`
	OrganizationID *string        `json:"organization_id" bson:"organization_id"`
	Status         *stat.Status   `json:"status" bson:"status"`
	pagination.Page
	order.OrderBy
}

type UsernameService interface {
	// Add creates a new username.
	Add(ctx context.Context, req AddRequest) error

	// Get returns the username with the specified ID.
	Get(ctx context.Context, id string) (Username, error)

	// Search returns usernames with the specified search criteria.
	List(ctx context.Context, req ListRequest) (pagination.Response[Username], error)

	// BePrimary sets the username with the specified ID as the primary username.
	BePrimary(context.Context, uuid.UUID, uuid.UUID, string) error

	// Hidden hides the username with the specified ID.
	Hidden(ctx context.Context, id string) error

	// Unhidden makes the username with the specified ID visible.
	Unhidden(ctx context.Context, id string) error

	// Reserve reserves the username with the specified ID.
	Reserve(ctx context.Context, id string) (Username, error)

	// Unreserve unreserves the username with the specified ID.
	Unreserve(ctx context.Context, id string, userID, organizationID uuid.UUID) error
}
