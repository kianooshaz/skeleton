package code

import "fmt"

type Code int

const (
	Success                 Code = 1000
	ErrInternalServerError  Code = 1001
	ErrInvalidPayloadSyntax Code = 1002
	ErrInvalidPagination    Code = 1003
	ErrInvalidOrderBy       Code = 1004

	ErrRequiredUserID Code = 1011
	ErrInvalidUserID  Code = 1012
	ErrNotFoundUser   Code = 1013
)

func (c Code) Error() string {
	return fmt.Sprintf("error code (%d)", c)
}
