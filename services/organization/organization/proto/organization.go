package orgproto

import (
	"context"
	"time"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
)

type Organization struct {
	ID        OrganizationID `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
}

type OrganizationService interface {
	Create(ctx context.Context) (CreateResponse, error)
	Get(ctx context.Context, req GetRequest) (GetResponse, error)
	List(ctx context.Context, req ListRequest) (ListResponse, error)
}

type CreateResponse struct {
	Data Organization `json:"data"`
}

type GetRequest struct {
	ID OrganizationID `json:"query"`
}

type GetResponse struct {
	Data Organization `json:"data"`
}

type ListRequest struct {
	pagination.Page
	order.OrderBy
}

type ListResponse pagination.Response[Organization]
