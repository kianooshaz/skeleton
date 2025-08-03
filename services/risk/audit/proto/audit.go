package auditproto

import (
	"context"

	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
)

type AuditService interface {
	Record(record Record)
	Get(ctx context.Context, req GetRequest) (GetResponse, error)
	List(ctx context.Context, req ListRequest) (ListResponse, error)
	Shutdown(ctx context.Context)
}

type GetRequest struct {
	ID RecordID `json:"id"`
}

type GetResponse struct {
	Data Record `json:"data"`
}

type ListRequest struct {
	pagination.Page
	order.OrderBy
}

type ListResponse pagination.Response[Record]
