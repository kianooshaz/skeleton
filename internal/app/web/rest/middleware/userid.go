package middleware

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/session"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
	"github.com/labstack/echo/v4"
)

func UserID(userService userproto.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDStr := c.Request().Header.Get("X-User-ID")
			var userID uuid.UUID
			var err error

			if userIDStr != "" {
				userID, err = uuid.Parse(userIDStr)
				if err == nil {
					ctx := session.SetUserID(c.Request().Context(), userID)
					c.SetRequest(c.Request().WithContext(ctx))
					return next(c)
				}

				slog.Error(
					"Error encountered while converting userID string to uuid",
					slog.String("error", err.Error()),
				)
			}

			user, err := userService.Create(c.Request().Context())
			if err != nil {
				slog.Error(
					"Error encountered while creating user in UserID middleware",
					slog.String("error", err.Error()),
				)
				return derror.ErrInternalSystem
			}
			userID = uuid.UUID(user.Data.ID)

			ctx := session.SetUserID(c.Request().Context(), userID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
