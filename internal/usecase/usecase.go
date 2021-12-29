package usecase

import (
	"context"
	"time"

	"github.com/josestg/justforfun/internal/domain/user"

	"github.com/josestg/justforfun/pkg/validate"
)

type Clock interface {
	Now() time.Time
}

type Validator interface {
	Validate(ctx context.Context, schema validate.Schema) error
}

type Identifier interface {
	NewUUID() string
}

type Repository struct {
	User user.Repository
}

type Provider struct {
	Clock      Clock
	Validator  Validator
	Identifier Identifier
	Repository *Repository
}
