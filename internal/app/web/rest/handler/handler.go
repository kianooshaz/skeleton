package handler

import (
	"github.com/kianooshaz/skeleton/service/usernamesrv"
	"github.com/kianooshaz/skeleton/service/usersrv"
)

type Handler struct {
	UserService     *usersrv.Service
	UsernameService *usernamesrv.Service
}
