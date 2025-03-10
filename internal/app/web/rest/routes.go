package rest

import (
	"github.com/kianooshaz/skeleton/internal/app/web/rest/handler"
	userService "github.com/kianooshaz/skeleton/services/user/user/service"
)

func (s *server) registerRoutes() {
	s.core.GET("/health", handler.HealthCheck)

	s.core.GET("/user", registerHandler(userService.Service.Get))
	s.core.GET("/user/list", registerHandler(userService.Service.List))
}
