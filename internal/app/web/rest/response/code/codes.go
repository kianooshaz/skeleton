package code

type Code int

const (
	Success Code = 1000

	InternalServerError  Code = 1011
	InvalidPayloadSyntax Code = 1012
	InvalidPagination    Code = 1013
	InvalidOrderBy       Code = 1014

	RequiredUserID Code = 1021
	InvalidUserID  Code = 1022
)
