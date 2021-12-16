package uuid

import (
	"regexp"
	"testing"
)

func TestNewV4(t *testing.T) {
	id, err := NewV4()
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	uuidV4Regex := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
	if !uuidV4Regex.MatchString(id.String()) {
		t.Fatalf("expecting pattern is match")
	}
}
