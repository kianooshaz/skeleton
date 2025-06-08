package rest

func (s *server) registerRoutes() {
	s.core.GET("/health", HealthCheck)

	// s.core.GET("/user", registerHandler(userService.Service.Get))
	// s.core.GET("/user/list", registerHandler(userService.Service.List))
}
