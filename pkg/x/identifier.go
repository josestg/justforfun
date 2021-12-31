package x

import "github.com/josestg/justforfun/pkg/uuid"

// Identifier knows how to generate a unique id.
// The id can be in UUID form, random string, or a number.
type Identifier struct{}

// NewIdentifier creates a new Identifier.
func NewIdentifier() *Identifier {
	return &Identifier{}
}

// NewUUID creates a new UUID.
func (i *Identifier) NewUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return id.String()
}
