package aunp

import (
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	iop "github.com/kianooshaz/skeleton/services/identify/organization/protocol"
	iup "github.com/kianooshaz/skeleton/services/identify/user/protocol"
	iunp "github.com/kianooshaz/skeleton/services/identify/username/protocol"
)

type AssignRequest struct {
	UserID         iup.UserID         `json:"user_id" bson:"user_id"`
	OrganizationID iop.OrganizationID `json:"organization_id" bson:"organization_id"`
	Username       iunp.Username      `json:"username" bson:"username"`
}

func (req *AssignRequest) Validate(minLen, MaxLen int) error {
	if len(req.UserID.String()) == 0 {
		return derror.ErrUserIDRequired

	}

	if len(req.OrganizationID.String()) == 0 {
		return derror.ErrOrganizationIDRequired
	}

	if len(req.Username) == 0 {
		return derror.ErrUsernameRequired
	}

	if len(req.Username) < minLen || len(req.Username) > MaxLen {
		return derror.ErrUsernameInvalid
	}

	return nil
}

type ListAssignedRequest struct {
	UserID         iup.UserID         `query:"user_id"`
	OrganizationID iop.OrganizationID `query:"organization_id"`
	pagination.Page
	order.OrderBy
}
