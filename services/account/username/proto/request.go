package usernameproto

import (
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/order"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	accprotocol "github.com/kianooshaz/skeleton/services/account/accounts/proto"
)

type AssignRequest struct {
	AccountID accprotocol.AccountID `json:"account_id" bson:"account_id"`
	Username  string                `json:"username" bson:"username"`
}

func (req *AssignRequest) Validate(minLen, MaxLen int) error {
	if len(req.AccountID.String()) == 0 {
		return derror.ErrAccountIDRequired
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
	AccountID accprotocol.AccountID `query:"account_id"`
	pagination.Page
	order.OrderBy
}
