package up

import (
	"context"
	"time"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
)

type User struct {
	ID        iup.UserID `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
}

type GetUserRequest struct {
	ID iup.UserID `json:"id" bson:"id"`
}

type ListUserRequest struct {
	pagination.Page
	order.OrderBy
}

type UserService interface {
	Create(ctx context.Context) (User, error)
	Get(ctx context.Context, req GetUserRequest) (User, error)
	List(ctx context.Context, req ListUserRequest) (pagination.Response[User], error)
}
