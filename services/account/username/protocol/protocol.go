package protocol

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/types"
)

type UsernameService interface {
	// Add creates a new username.
	Add(ctx context.Context, req AddRequest) error

	// Get returns the username with the specified ID.
	Get(ctx context.Context, id string) (Username, error)

	// Search returns usernames with the specified search criteria.
	List(ctx context.Context, userID *string, organizationID *string, status *types.Status, limit int, offset int) ([]Username, error)

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

type Username struct {
	ID             string               `json:"id" bson:"id"`
	UserID         types.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID types.OrganizationID `json:"organization_id" bson:"organization_id"`
	Status         types.Status         `json:"status" bson:"status"`
	CreatedAt      time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at" bson:"updated_at"`
}

type AddRequest struct {
	UserID         types.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID types.OrganizationID `json:"organization_id" bson:"organization_id"`
	ID             string               `json:"id" bson:"id"`
}

type ListRequest struct {
	UserID         *string       `json:"user_id" bson:"user_id"`
	OrganizationID *string       `json:"organization_id" bson:"organization_id"`
	Status         *types.Status `json:"status" bson:"status"`
	Limit          int           `json:"limit" bson:"limit"`
	Offset         int           `json:"offset" bson:"offset"`
}

type ListResponse struct {
	Items []Username `json:"items" bson:"items"`
	Total int64      `json:"total" bson:"total"`
}
