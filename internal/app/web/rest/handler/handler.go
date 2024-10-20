package handler

import (
	"github.com/kianooshaz/skeleton/protocol"
)

type Handler struct {
	UserService     protocol.ServiceUser
	UsernameService protocol.ServiceUsername
}
