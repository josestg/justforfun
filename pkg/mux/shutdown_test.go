package mux

import (
	"fmt"
	"testing"
)

func TestErrorShutdown(t *testing.T) {
	const message = "A SHUTDOWN ERROR"
	shutdownError := NewShutdownError(message)

	if shutdownError == nil {
		t.Fatalf("Expecting error not nil")
	}

	wrappedErr := fmt.Errorf("a: %w", shutdownError)
	wrappedErr = fmt.Errorf("b: %w", wrappedErr)
	wrappedErr = fmt.Errorf("c: %w", wrappedErr)

	if !IsShutdownError(wrappedErr) {
		t.Fatalf("Expecting shutdown error")
	}

	// use %v instead %w.
	wrappedErr = fmt.Errorf("a: %v", shutdownError)

	if IsShutdownError(wrappedErr) {
		t.Fatalf("Expecting not shutdown error")
	}
}
