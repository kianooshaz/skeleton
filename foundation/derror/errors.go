// Package derror defines errors used throughout the application.
package derror

import "errors"

// system errors
var ErrInternalSystem = errors.New("100000")
var ErrUndefinedPathAndMethod = errors.New("100001")
var ErrInvalidJsonFormat = errors.New("100002")
var ErrInvalidQueryParameter = errors.New("100003")
var ErrUnknownOrder = errors.New("100004")
var ErrUnknownOrderDirection = errors.New("100005")
var ErrInvalidPage = errors.New("100006")
var ErrInvalidRows = errors.New("100007")
var ErrPageValueTooSmall = errors.New("100008")
var ErrRowsValueTooSmall = errors.New("100009")
var ErrRowsValueTooLarge = errors.New("100010")

// user errors
var ErrUserNotFound = errors.New("100100")
var ErrUserAlreadyExists = errors.New("100101")
var ErrUsernameNotFound = errors.New("100102")
var ErrUsernameAlreadyExists = errors.New("100103")
var ErrUsernameInvalid = errors.New("100104")
var ErrUsernameLength = errors.New("100105")
var ErrUsernameMaxPerUser = errors.New("100106")
var ErrUsernameMaxPerOrganization = errors.New("100107")
var ErrUsernameInvalidCharacters = errors.New("100108")
var ErrUsernameNotReserved = errors.New("100109")

var ErrPasswordInvalid = errors.New("100200")
var ErrPasswordIsWeak = errors.New("100201")
var ErrPasswordIsCommon = errors.New("100202")
var ErrPasswordIsInHistory = errors.New("100203")
