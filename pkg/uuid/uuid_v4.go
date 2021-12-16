package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

// NewV4 creates a new uuid version 4.
// This function uses the rand.Reader as random generator.
func NewV4() (UUID, error) {
	return NewV4WithReader(rand.Reader)
}

// NewV4WithReader creates a new uuid version 4 with a given random reader.
func NewV4WithReader(reader io.Reader) (UUID, error) {
	var id UUID
	_, err := io.ReadFull(reader, id[:])
	if err != nil {
		return Nil, fmt.Errorf("%w: reading random byte", err)
	}

	id[6] = (id[6] & 0x0f) | 0x40 // Version 4
	id[8] = (id[8] & 0x3f) | 0x80 // Variant is 10
	return id, nil
}
