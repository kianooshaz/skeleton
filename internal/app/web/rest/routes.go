package rest

import (
	"github.com/kianooshaz/skeleton/internal/app/web/rest/handler"
	userProtocol "github.com/kianooshaz/skeleton/modules/user/user/protocol"
	userService "github.com/kianooshaz/skeleton/modules/user/user/service"
)

func (s *server) registerRoutes() {
	s.core.GET("/health", handler.HealthCheck)

	s.core.GET("/user", registerHandler[userProtocol.GetUserRequest, userProtocol.User](userService.Service.Get))
	s.core.GET("/user/list", registerHandler[userProtocol.ListUserRequest, userProtocol.ListUserResponse](userService.Service.List))
}
