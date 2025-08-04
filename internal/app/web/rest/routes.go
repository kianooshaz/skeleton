package rest

import (
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

func (s *server) registerRoutes(
	UserService userproto.UserService,
	OrganizationService orgproto.OrganizationService,
	PasswordService passwordproto.PasswordService,
	UsernameService usernameproto.UsernameService,
	AuditService auditproto.AuditService,
) {
	s.core.GET("/health", HealthCheck)

	// s.core.GET("/user", registerHandler(userService.Service.Get))
	// s.core.GET("/user/list", registerHandler(userService.Service.List))
}
