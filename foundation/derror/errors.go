// Package derror defines errors used throughout the application.
package derror

import "errors"

// system errors
var ErrInternalSystem = errors.New("100000")
var ErrUndefinedPathAndMethod = errors.New("100001")
var ErrInvalidJsonFormat = errors.New("100002")
var ErrInvalidQueryParameter = errors.New("100003")

// user errors
var ErrUserNotFound = errors.New("100100")
var ErrUserAlreadyExists = errors.New("100101")

// username errors =
var ErrUsernameNotFound = errors.New("100200")
var ErrUsernameAlreadyExists = errors.New("100201")
var ErrUsernameInvalid = errors.New("100202")
var ErrUsernameLength = errors.New("100203")
var ErrUsernameMaxPerUser = errors.New("100204")
