package provider

import (
	"context"
	"time"

	"github.com/josestg/justforfun/internal/iam"

	"github.com/josestg/justforfun/internal/repository"
	"github.com/josestg/justforfun/pkg/validate"
)

// Clock is a time capture.
type Clock interface {
	// Now returns the current time at some location.
	Now() time.Time
}

// Validator knows how to validate validation schema.
type Validator interface {
	// Validate validates given schema and transform the error into
	// standardized format.
	Validate(ctx context.Context, schema validate.Schema) error
}

// Identifier knows how to generate a unique id.
type Identifier interface {
	// NewUUID returns a UUID.
	NewUUID() string
}

// Provider holds all dependencies for the UseCase to work with.
type Provider struct {
	Clock      Clock                 // for capturing event timestamp.
	Validator  Validator             // for validation input.
	Identifier Identifier            // for generating unique id.
	Tokenizer  iam.Tokenizer         // for creating and validating token-based identify.
	Repository *repository.Container // provides centralized repository instances.
}
