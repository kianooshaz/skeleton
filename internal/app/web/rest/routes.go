package rest

func (s *Server) registerRoutes() {
	s.core.GET("/health", s.handler.HealthCheck)

	s.core.POST("/user", s.handler.NewUser)
	s.core.GET("/user", s.handler.GetUser)
}
