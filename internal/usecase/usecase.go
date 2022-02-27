package usecase

import (
	"github.com/josestg/justforfun/internal/domain/health"
	"github.com/josestg/justforfun/internal/usecase/internal/user"
	"github.com/josestg/justforfun/internal/usecase/provider"
)

// Container contains all UseCase instances.
type Container struct {
	SystemHealthCheck *health.UseCase
	UserRegistration  *user.Registration
	UserTokenization  *user.Tokenization
}

// NewContainer creates a new Container.
func NewContainer(provider *provider.Provider) *Container {
	return &Container{
		SystemHealthCheck: nil,
		UserRegistration:  user.NewRegistration(provider),
		UserTokenization:  user.NewTokenization(provider),
	}
}
