package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func registerHandler[T any, S any](handler func(ctx context.Context, req T) (S, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return err
		}

		res, err := handler(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}

func registerHandlerNoResponse[T any](handler func(ctx context.Context, req T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req T
		if err := c.Bind(&req); err != nil {
			return err
		}

		err := handler(c.Request().Context(), req)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}

func registerHandlerNoRequest[T any, S any](handler func(ctx context.Context) (S, error)) echo.HandlerFunc {
	return func(c echo.Context) error {
		res, err := handler(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}
}

func registerHandlerNoRequestNoResponse(handler func(ctx context.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := handler(c.Request().Context())
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}

}
