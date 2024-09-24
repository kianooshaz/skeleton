package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kianooshaz/skeleton/internal/app/web/rest/response"
	"github.com/kianooshaz/skeleton/internal/app/web/rest/response/code"
	"github.com/kianooshaz/skeleton/protocol"
	"github.com/labstack/echo/v4"
)

func (h *Handler) NewUser(c echo.Context) error {
	user, err := h.User.New(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError(code.InternalServerError))
	}

	data := response.User{
		ID:        user.ID().String(),
		CreatedAt: user.CreatedAt().Unix(),
	}

	return c.JSON(http.StatusOK, response.New(data, nil))
}

func (h *Handler) GetUser(c echo.Context) error {
	idStr := c.QueryParam("user_id")

	if idStr == "" {
		return c.JSON(http.StatusBadRequest, response.NewError(code.RequiredUserID))
	}

	id, err := uuid.FromBytes([]byte(idStr))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.NewError(code.InvalidUserID))
	}

	user, err := h.User.Get(c.Request().Context(), protocol.ID(id))
	if err != nil {

		return c.JSON(http.StatusInternalServerError, response.NewError(code.InternalServerError))
	}

	return c.JSON(http.StatusOK, response.New(user, nil))
}
