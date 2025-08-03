package passwordproto

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	accproto "github.com/kianooshaz/skeleton/services/account/accounts/proto"
)

type Password struct {
	ID           uuid.UUID          `json:"id" bson:"id"`
	AccountID    accproto.AccountID `json:"account_id" bson:"account_id"`
	PasswordHash string             `json:"password_hash" bson:"password_hash"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type PasswordService interface {
	Update(ctx context.Context, req UpdateRequest) error
	Verify(ctx context.Context, password string) error
	Guidelines() (GuidelinesResponse, error)
	Get(ctx context.Context, id uuid.UUID) (Password, error)
	List(ctx context.Context, req ListRequest) (ListResponse, error)
}

type UpdateRequest struct {
	OTP         string             `json:"otp"`
	NewPassword string             `json:"new_password"`
	AccountID   accproto.AccountID `json:"account_id"`
}

type ListRequest struct {
	AccountID accproto.AccountID `query:"account_id"`
	pagination.Page
	order.OrderBy
}

type ListResponse pagination.Response[Password]

type Guidelines struct {
	Required   []string `json:"required"`
	BetterHave []string `json:"better_have"`
}

type GuidelinesResponse struct {
	Data Guidelines `json:"data"`
}
