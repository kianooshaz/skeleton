package protocol

import (
	"context"
	"time"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/foundation/types"
)

type User struct {
	ID        types.UserID `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
}

type GetUserRequest struct {
	ID types.UserID `json:"id" bson:"id"`
}

type ListUserRequest struct {
	pagination.Page
	order.OrderBy
}

type ListUserResponse struct {
	Users []User `json:"users" bson:"users"`
	Total int64  `json:"total" bson:"total"`
}

type UserService interface {
	Create(ctx context.Context) (User, error)
	Get(ctx context.Context, req GetUserRequest) (User, error)
	List(ctx context.Context, req ListUserRequest) (ListUserResponse, error)
}
