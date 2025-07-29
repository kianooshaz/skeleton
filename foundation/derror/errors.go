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
var ErrRateLimitExceeded = errors.New("100011")

// user errors
var ErrUserIDRequired = errors.New("100100")
var ErrUserNotFound = errors.New("100101")
var ErrUserAlreadyExists = errors.New("100102")

var ErrPasswordInvalid = errors.New("100200")
var ErrPasswordIsWeak = errors.New("100201")
var ErrPasswordIsCommon = errors.New("100202")
var ErrPasswordUsedBefore = errors.New("100203")
var ErrPasswordNotFound = errors.New("100204")

var ErrUsernameNotFound = errors.New("100300")
var ErrUsernameInvalid = errors.New("100301")
var ErrUsernameMaxPerUser = errors.New("100302")
var ErrUsernameMaxPerOrganization = errors.New("100303")
var ErrUsernameNotReserved = errors.New("100304")
var ErrUsernameCannotBeAssigned = errors.New("100305")
var ErrUsernameLocked = errors.New("100306")
var ErrUsernameAlreadyExists = errors.New("100307")
var ErrUsernameRequired = errors.New("100308")

var ErrOrganizationIDRequired = errors.New("100400")

var ErrAccountIDRequired = errors.New("100500")
