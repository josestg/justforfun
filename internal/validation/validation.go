package validation

import (
	"context"
	"encoding/json"

	"github.com/josestg/justforfun/pkg/validate"
	"github.com/josestg/justforfun/pkg/xerrs"
)

type FieldMessages struct {
	Field    string   `json:"field"`
	Messages []string `json:"messages"`
}

type Error struct {
	fields []FieldMessages
}

func NewError() *Error {
	return &Error{
		fields: make([]FieldMessages, 0),
	}
}

func (e *Error) Push(fields ...FieldMessages) {
	e.fields = append(e.fields, fields...)
}

func (e *Error) String() string {
	b, err := json.Marshal(e.fields)
	if err != nil {
		return err.Error()
	}

	return string(b)
}

func (e *Error) Error() string {
	return e.String()
}

type Validator struct {
	transformer validate.ErrorTransformer
}

func NewValidator(transformer validate.ErrorTransformer) *Validator {
	return &Validator{
		transformer: transformer,
	}
}

func (v *Validator) Validate(ctx context.Context, schema validate.Schema) error {
	err := schema.Valid(ctx, v.transformer)
	switch et := xerrs.Cause(err).(type) {
	default:
		return err
	case validate.Errors:
		ve := NewError()
		for k, v := range et {
			ve.Push(FieldMessages{
				Field:    k,
				Messages: v,
			})
		}

		return ve
	}
}
