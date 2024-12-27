package handler

import (
	userProtocol "github.com/kianooshaz/skeleton/modules/user/user/protocol"
)

type Handler struct {
	UserService userProtocol.UserService
}
