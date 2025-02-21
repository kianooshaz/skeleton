package rest

import (
	"log/slog"
	"net/http"

	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/labstack/echo/v4"
)

func errorResponse(err error, c echo.Context) {
	// TODO: test this function
	status, ok := DerrorToHTTPStatus[err]
	if !ok {
		status = http.StatusInternalServerError
	}

	if err := c.JSON(status, echo.Map{
		"error": err.Error(),
	}); err != nil {
		slog.Error(
			"error at sending error response",
			slog.Any("error", err),
			slog.Any("status", status),
			slog.String("package", "rest"),
		)
	}

}

var DerrorToHTTPStatus = map[error]int{
	derror.ErrInternalSystem:         http.StatusInternalServerError,
	derror.ErrUndefinedPathAndMethod: http.StatusBadRequest,
	derror.ErrInvalidJsonFormat:      http.StatusBadRequest,
	derror.ErrInvalidQueryParameter:  http.StatusBadRequest,
	derror.ErrUnknownOrder:           http.StatusBadRequest,
	derror.ErrUnknownOrderDirection:  http.StatusBadRequest,
	derror.ErrUnknownOrder:           http.StatusBadRequest,
	derror.ErrUnknownOrderDirection:  http.StatusBadRequest,
	derror.ErrInvalidPage:            http.StatusBadRequest,
	derror.ErrInvalidRows:            http.StatusBadRequest,
	derror.ErrPageValueTooSmall:      http.StatusBadRequest,
	derror.ErrRowsValueTooSmall:      http.StatusBadRequest,
	derror.ErrRowsValueTooLarge:      http.StatusBadRequest,

	derror.ErrUserNotFound:               http.StatusBadRequest,
	derror.ErrUserAlreadyExists:          http.StatusBadRequest,
	derror.ErrUsernameNotFound:           http.StatusBadRequest,
	derror.ErrUsernameAlreadyExists:      http.StatusBadRequest,
	derror.ErrUsernameInvalid:            http.StatusBadRequest,
	derror.ErrUsernameLength:             http.StatusBadRequest,
	derror.ErrUsernameMaxPerUser:         http.StatusBadRequest,
	derror.ErrUsernameMaxPerOrganization: http.StatusBadRequest,
	derror.ErrUsernameInvalidCharacters:  http.StatusBadRequest,
	derror.ErrUsernameNotReserved:        http.StatusBadRequest,

	derror.ErrPasswordInvalid:     http.StatusBadRequest,
	derror.ErrPasswordIsWeak:      http.StatusBadRequest,
	derror.ErrPasswordIsCommon:    http.StatusBadRequest,
	derror.ErrPasswordIsInHistory: http.StatusBadRequest,
}
