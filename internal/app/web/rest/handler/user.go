package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) NewUser(c echo.Context) error {
	user, err := h.User.New(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}
