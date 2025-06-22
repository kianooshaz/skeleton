package middleware

import (
	"log/slog"

	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/ratelimit"
	"github.com/labstack/echo/v4"
)

// This middleware enforces rate limiting per user. For each request, a userID is determined: if the user is signed in, their userID is used; otherwise, a guest user is created with status 'guest'. Rate limiting is then applied based on this userID.
func RateLimit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// TODO: Replace c.RealIP() with logic to extract userID if signed in, otherwise assign/create a guest user with status 'guest'.
			allow, err := ratelimit.SlidingWindowRateLimiter.Allow(c.Request().Context(), c.RealIP())
			if err != nil {
				slog.Error("failed to check rate limit", slog.String("error", err.Error()), slog.String("ip", c.RealIP()))
				return err
			}

			if !allow {
				return derror.ErrRateLimitExceeded
			}

			return next(c)
		}
	}
}
