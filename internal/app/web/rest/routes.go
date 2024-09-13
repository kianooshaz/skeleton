package rest

import "github.com/kianooshaz/skeleton/internal/app/web/rest/handler"

func (s *Server) registerRoutes() {
	s.core.GET("/health", handler.HealthCheck())
}
