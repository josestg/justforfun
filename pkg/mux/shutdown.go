package mux

import (
	"errors"
	"os"
)

// ShutdownChannel is a channel that Mux used to tell the application to
// shut down gracefully by sending a termination signal (syscall.SIGTERM).
type ShutdownChannel chan os.Signal

// shutdownError is a shutdown error type
type shutdownError string

// NewShutdownError creates a new shutdown error with message.
func NewShutdownError(msg string) error {
	return shutdownError(msg)
}

// Error implements the error interface.
func (s shutdownError) Error() string {
	return string(s)
}

// IsShutdownError returns true if the error is a shutdown error.
func IsShutdownError(err error) bool {
	var se shutdownError
	return errors.As(err, &se)
}
