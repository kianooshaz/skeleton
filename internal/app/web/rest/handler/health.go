package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "up",
	})
}
