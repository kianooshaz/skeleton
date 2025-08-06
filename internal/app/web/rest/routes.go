package rest

import (
	usernameproto "github.com/kianooshaz/skeleton/services/account/username/proto"
	passwordproto "github.com/kianooshaz/skeleton/services/authentication/password/proto"
	orgproto "github.com/kianooshaz/skeleton/services/organization/organization/proto"
	auditproto "github.com/kianooshaz/skeleton/services/risk/audit/proto"
	userproto "github.com/kianooshaz/skeleton/services/user/user/proto"
)

func (s *server) registerRoutes(
	userService userproto.UserService,
	organizationService orgproto.OrganizationService,
	passwordService passwordproto.PasswordService,
	usernameService usernameproto.UsernameService,
	auditService auditproto.AuditService,
) {
	s.core.GET("/health", HealthCheck)

	s.core.GET("/user", registerHandler(userService.Get))
	// s.core.GET("/user/list", registerHandler(userService.Service.List))
}
