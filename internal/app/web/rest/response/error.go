package response

import "github.com/kianooshaz/skeleton/internal/app/web/rest/response/code"

func NewError(c code.Code) *Error {
	return &Error{
		Code: c,
	}
}

type Error struct {
	Code code.Code
}
