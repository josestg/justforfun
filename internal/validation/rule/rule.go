package rule

import (
	"context"
	"fmt"
	"reflect"

	"github.com/josestg/justforfun/internal/wording"

	"github.com/josestg/justforfun/pkg/validate"
)

var TagPrefix = "validation"

const (
	argKeyActual = "actual"
	argKeyExpect = "expect"
)

func expectActualArg(expect, actual interface{}) map[string]interface{} {
	return map[string]interface{}{
		argKeyActual: actual,
		argKeyExpect: expect,
	}
}

func tag(name string, args map[string]interface{}) error {
	key := fmt.Sprintf("%s.%s", TagPrefix, name)
	err := wording.NewError(key, args)
	return err
}

func Required(allowZeroValue ...bool) validate.Rule {
	return validate.RuleFunc(func(ctx context.Context, v interface{}) error {
		err := tag("required", nil)
		rv := reflect.ValueOf(v)
		if rv.IsNil() {
			return err
		}

		if len(allowZeroValue) > 0 && allowZeroValue[0] {
			return nil
		}

		e := rv.Elem()
		if e.IsZero() {
			return err
		}

		return nil
	})
}

func Len(min, max int) validate.Rule {
	return validate.RuleFunc(func(ctx context.Context, v interface{}) error {
		rv := reflect.ValueOf(v)
		e := rv.Elem()
		switch e.Kind() {
		default:
			return nil
		case reflect.Slice, reflect.Array, reflect.String, reflect.Map, reflect.Chan:
			if e.Len() < min && min != -1 {
				return tag("min_length", expectActualArg(min, e.Len()))
			}

			if e.Len() > max && max != -1 {
				return tag("max_length", expectActualArg(max, e.Len()))
			}
		}

		return nil
	})
}

func Email() validate.Rule {
	return validate.RuleFunc(func(ctx context.Context, v interface{}) error {
		rv := reflect.ValueOf(v)
		e := rv.Elem()
		switch e.Kind() {
		default:
			return nil
		case reflect.String:
			s := e.String()
			if !emailRegex.MatchString(s) {
				return tag("invalid_email_format", nil)
			}
		}

		return nil
	})
}
