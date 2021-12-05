package validate

import (
	"context"
	"encoding/json"
	"errors"
)

// ctxType is a type for context.
type ctxType int

// breakerError is a breaker error.
type breakerError struct {
	err error
}

// BreakerError marks the given error as breaker error.
func BreakerError(err error) error {
	return &breakerError{err: err}
}

func (b *breakerError) Error() string {
	return b.err.Error()
}

const (
	fieldContextKey ctxType = 0
)

func contextWithField(ctx context.Context, field string) context.Context {
	return context.WithValue(ctx, fieldContextKey, field)
}

// FieldFromContext gets a field name for context.
func FieldFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(fieldContextKey).(string)
	return v, ok
}

// ErrorTransformer is a contract for error transformer.
type ErrorTransformer interface {
	// Transform transforms the given error into anything you want.
	Transform(ctx context.Context, err error) string
}

// TransformFunc is an adapter to allow the use of ordinary functions as ErrorTransformer.
type TransformFunc func(ctx context.Context, err error) string

func (t TransformFunc) Transform(ctx context.Context, err error) string {
	return t(ctx, err)
}

// Rule is a validation rule.
type Rule interface {
	// Valid returns an error if v not satisfied the rule.
	Valid(ctx context.Context, v interface{}) error
}

// RuleFunc is an adapter to allow the use of ordinary functions as Rule.
type RuleFunc func(ctx context.Context, v interface{}) error

func (r RuleFunc) Valid(ctx context.Context, v interface{}) error {
	return r(ctx, v)
}

// Validator knows how to validate the Schema.
type Validator interface {
	// Validate validates the Schema.
	Validate(ctx context.Context) []error
}

// Func is an adapter to allow the use of ordinary functions as Validator.
type Func func(ctx context.Context) []error

func (f Func) Validate(ctx context.Context) []error {
	return f(ctx)
}

// Field creates a new Field validator.
func Field(value interface{}, rules ...Rule) Validator {
	fn := func(ctx context.Context) []error {
		_errors := make([]error, 0)
		for _, rule := range rules {
			if rule == nil {
				continue
			}

			if err := rule.Valid(ctx, value); err != nil {
				var bErr *breakerError
				if breakable := errors.As(err, &bErr); breakable {
					_errors = append(_errors, bErr.err)
					break
				}

				_errors = append(_errors, err)
			}
		}

		return _errors
	}

	return Func(fn)
}

// Errors represents validation errors.
type Errors map[string][]string

func (e Errors) Error() string {
	b, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

// Schema is a validation schema.
type Schema map[string]Validator

// Valid returns the Errors if Schema not valid, otherwise returns nil.
func (s Schema) Valid(ctx context.Context, transformer ErrorTransformer) error {
	_errors := make(Errors)

	for name, validator := range s {
		messages := make([]string, 0)

		fieldContext := contextWithField(ctx, name)

		errs := validator.Validate(fieldContext)
		for _, err := range errs {
			messages = append(messages, transformer.Transform(fieldContext, err))
		}

		if len(messages) != 0 {
			_errors[name] = messages
		}
	}

	if len(_errors) != 0 {
		return _errors
	}

	return nil
}
