package userproto

import (
	"context"
	"time"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
)

type User struct {
	ID        UserID    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserService interface {
	Create(ctx context.Context) (CreateResponse, error)
	Get(ctx context.Context, req GetRequest) (GetResponse, error)
	List(ctx context.Context, req ListRequest) (ListResponse, error)
}

type CreateResponse struct {
	Data User `json:"data"`
}

type GetRequest struct {
	ID UserID `json:"query"`
}

type GetResponse struct {
	Data User `json:"data"`
}

type ListRequest struct {
	pagination.Page
	order.OrderBy
}

type ListResponse pagination.Response[User]
