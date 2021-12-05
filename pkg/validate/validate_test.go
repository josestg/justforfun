package validate_test

import (
	"context"
	"errors"
	"github.com/josestg/justforfun/pkg/validate"
	"reflect"
	"regexp"
	"testing"
)

const (
	emailRegexString = "^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

var Email = regexp.MustCompile(emailRegexString)

var (
	errRequired  = errors.New("is required")
	errZeroValue = errors.New("is zero value")
	errMinLength = errors.New("is smaller than min")
	errMaxLength = errors.New("is larger than max")
	errNotEmail  = errors.New("invalid email format")
)

var required = validate.RuleFunc(func(ctx context.Context, v interface{}) error {
	if reflect.ValueOf(v).IsNil() {
		return validate.BreakerError(errRequired)
	}

	return nil
})

var notZeroValue = validate.RuleFunc(func(ctx context.Context, v interface{}) error {
	if reflect.ValueOf(v).Elem().IsZero() {
		return errZeroValue
	}

	return nil
})

var isEmail = validate.RuleFunc(func(ctx context.Context, v interface{}) error {
	s := reflect.ValueOf(v).Elem().String()
	if !Email.MatchString(s) {
		return errNotEmail
	}

	return nil
})

var length = func(min, max int) validate.Rule {
	return validate.RuleFunc(func(ctx context.Context, v interface{}) error {
		l := reflect.ValueOf(v).Elem().Len()
		if l < min {
			return errMinLength
		}

		if l > max {
			return errMaxLength
		}

		return nil
	})
}

func TestField(t *testing.T) {
	t.Run("if given nil value, the next rule will be not executed", func(t *testing.T) {
		var name *string = nil
		validator := validate.Field(name, required, notZeroValue)
		_errors := validator.Validate(context.Background())
		if len(_errors) != 1 {
			t.Fatal("expecting contains an error")
		}

		if _errors[0] != errRequired {
			t.Fatalf("expecting %v but bot %v", errRequired, _errors[0])
		}
	})

	t.Run("if given en empty value the next rule will be executed unless a breaker error is found", func(t *testing.T) {
		name := ""

		// nil rule will be omitted.
		validator := validate.Field(&name, required, notZeroValue, nil, length(2, 10))
		_errors := validator.Validate(context.Background())
		if len(_errors) != 2 {
			t.Fatal("expecting contains an error")
		}

		if _errors[0] != errZeroValue {
			t.Fatalf("expecting %v but bot %v", errZeroValue, _errors[0])
		}

		if _errors[1] != errMinLength {
			t.Fatalf("expecting %v but bot %v", errMinLength, _errors[0])
		}
	})
}

func TestSchema_Valid(t *testing.T) {
	t.Run("schema with errors", func(t *testing.T) {
		var (
			name  = ""
			email = "abc"
		)

		schema := validate.Schema{
			"name":  validate.Field(&name, required, notZeroValue, nil, length(2, 10)),
			"email": validate.Field(&email, required, notZeroValue, isEmail),
		}

		_errors := schema.Valid(context.Background(), validate.TransformFunc(func(ctx context.Context, err error) string {
			field, ok := validate.FieldFromContext(ctx)
			if !ok {
				return "expecting context contains field name"
			}

			if field != "name" && field != "email" {
				return "expecting field value must be email or name"
			}

			return err.Error()
		}))

		vErrors, ok := _errors.(validate.Errors)
		if !ok {
			t.Fatalf("expecting validation.Error")
		}

		if len(vErrors) != 2 {
			t.Fatalf("expecting 2 errors")
		}

		nameErrors := []string{errZeroValue.Error(), errMinLength.Error()}
		if !reflect.DeepEqual(nameErrors, vErrors["name"]) {
			t.Fatalf("expecting vErros[name] contains %v", nameErrors)
		}

		emailErrors := []string{errNotEmail.Error()}
		if !reflect.DeepEqual(emailErrors, vErrors["email"]) {
			t.Fatalf("expecting vErros[email] contains %v", emailErrors)
		}

		s := vErrors.Error()
		expected := `{"email":["invalid email format"],"name":["is zero value","is smaller than min"]}`
		if s != expected {
			t.Fatalf("expecting %v but got %v", expected, s)
		}
	})

	t.Run("schema without errors", func(t *testing.T) {
		var (
			name  = "bob"
			email = "abc@mail.com"
		)

		schema := validate.Schema{
			"name":  validate.Field(&name, required, notZeroValue, nil, length(2, 10)),
			"email": validate.Field(&email, required, notZeroValue, isEmail),
		}

		_errors := schema.Valid(context.Background(), validate.TransformFunc(func(ctx context.Context, err error) string {
			field, ok := validate.FieldFromContext(ctx)
			if !ok {
				return "expecting context contains field name"
			}

			if field != "name" && field != "email" {
				return "expecting field value must be email or name"
			}

			return err.Error()
		}))

		if _errors != nil {
			t.Fatalf("expecting error nil")
		}
	})
}
