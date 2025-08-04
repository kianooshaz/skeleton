package birthdayproto

import (
	"context"
	"time"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

// Birthday represents a user's birthday information.
type Birthday struct {
	ID          BirthdayID       `json:"id"`
	UserID      userproto.UserID `json:"user_id"`
	DateOfBirth time.Time        `json:"date_of_birth"`
	Age         int              `json:"age"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// BirthdayService defines the interface for birthday operations.
type BirthdayService interface {
	Create(ctx context.Context, req CreateRequest) (CreateResponse, error)
	Get(ctx context.Context, req GetRequest) (GetResponse, error)
	GetByUserID(ctx context.Context, req GetByUserIDRequest) (GetByUserIDResponse, error)
	Update(ctx context.Context, req UpdateRequest) (UpdateResponse, error)
	Delete(ctx context.Context, req DeleteRequest) error
	List(ctx context.Context, req ListRequest) (ListResponse, error)
}

// CreateRequest represents the request to create a birthday.
type CreateRequest struct {
	UserID      userproto.UserID `json:"user_id" validate:"required"`
	DateOfBirth time.Time        `json:"date_of_birth" validate:"required"`
}

// CreateResponse represents the response from creating a birthday.
type CreateResponse struct {
	Data Birthday `json:"data"`
}

// GetRequest represents the request to get a birthday by ID.
type GetRequest struct {
	ID BirthdayID `json:"id" validate:"required"`
}

// GetResponse represents the response from getting a birthday.
type GetResponse struct {
	Data Birthday `json:"data"`
}

// GetByUserIDRequest represents the request to get a birthday by user ID.
type GetByUserIDRequest struct {
	UserID userproto.UserID `json:"user_id" validate:"required"`
}

// GetByUserIDResponse represents the response from getting a birthday by user ID.
type GetByUserIDResponse struct {
	Data Birthday `json:"data"`
}

// UpdateRequest represents the request to update a birthday.
type UpdateRequest struct {
	ID          BirthdayID `json:"id" validate:"required"`
	DateOfBirth time.Time  `json:"date_of_birth" validate:"required"`
}

// UpdateResponse represents the response from updating a birthday.
type UpdateResponse struct {
	Data Birthday `json:"data"`
}

// DeleteRequest represents the request to delete a birthday.
type DeleteRequest struct {
	ID BirthdayID `json:"id" validate:"required"`
}

// ListRequest represents the request to list birthdays.
type ListRequest struct {
	pagination.Page
	order.OrderBy
	UserID     *userproto.UserID `json:"user_id,omitempty"`
	MinAge     *int              `json:"min_age,omitempty"`
	MaxAge     *int              `json:"max_age,omitempty"`
	BirthMonth *int              `json:"birth_month,omitempty"`
}

// ListResponse represents the response from listing birthdays.
type ListResponse pagination.Response[Birthday]
