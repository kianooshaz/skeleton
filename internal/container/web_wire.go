//go:build wireinject
// +build wireinject

// Package container provides web container dependency injection using Google Wire.
package container

import (
	"github.com/google/wire"
)

// WebProviderSet contains providers for web container dependencies.
var WebProviderSet = wire.NewSet(
	BaseProviderSet,
	ProvideRestServer,
	ProvideUserService,
	ProvideOrgService,
	ProvideUsernameService,
	ProvideAuditService,
	ProvidePasswordService,
)

// NewWebContainer creates a new web dependency injection container.
func NewWebContainer() (*WebContainer, error) {
	wire.Build(
		WebProviderSet,
		wire.Struct(new(WebContainer), "*"),
	)
	return &WebContainer{}, nil
}
