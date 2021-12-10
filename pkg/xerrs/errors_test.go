package xerrs

import (
	"errors"
	"fmt"
	"testing"
)

func TestCause(t *testing.T) {
	if Wrap(nil, "message") != nil {
		t.Fatalf("expecting nil")
	}

	rootErr := New("root of error")

	const expectedMessage = "int: 4: int: 3: int: 2: int: 1: int: 0: root of error"
	wrappedErr := rootErr
	for i := 0; i < 5; i++ {
		wrappedErr = Wrap(wrappedErr, fmt.Sprintf("int: %d", i))
	}

	causer := Cause(wrappedErr)
	if causer != rootErr {
		t.Fatalf("expecting %v but got %v", rootErr, causer)
	}

	if wrappedErr.Error() != expectedMessage {
		t.Fatalf("expecting %v but got %v", expectedMessage, wrappedErr.Error())
	}

	if !errors.As(wrappedErr, &rootErr) {
		t.Fatalf("expecting equal")
	}
}
