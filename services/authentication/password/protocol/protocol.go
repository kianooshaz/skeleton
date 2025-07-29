package authpassp

import (
	"time"

	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
)

type Password struct {
	UserID         iup.UserID
	OrganizationID iop.OrganizationID
	Password       string
	CreatedAt      time.Time
}

type PasswordServices interface {
	Update(ctx, req UpdateRequest) error
	Verify(ctx, password string) error
	Guidelines() (GuidelinesResponse, error)
}

type UpdateRequest struct {
	OTP            string             `json:"otp"`
	NewPassword    string             `json:"new_password"`
	UserID         iup.UserID         `json:"user_id"`
	OrganizationID iop.OrganizationID `json:"organization_id"`
}

type Guidelines struct {
	Required   []string `json:"required"`
	BetterHave []string `json:"better_have"`
}

type GuidelinesResponse struct {
	Data Guidelines `json:"data"`
}
