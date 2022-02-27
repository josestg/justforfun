package wording

import (
	"context"
	"fmt"

	"github.com/josestg/justforfun/pkg/validate"

	"github.com/josestg/justforfun/pkg/xerrs"
)

type Error struct {
	key  string
	args map[string]interface{}
}

func NewError(key string, args map[string]interface{}) error {
	return &Error{
		key:  key,
		args: args,
	}
}

func (e *Error) Error() string {
	return e.key
}

func (e *Error) Args() map[string]interface{} {
	return e.args
}

type Translator interface {
	Translate(language, messageID string, args interface{}) (string, error)
}

type Wording struct {
	translator Translator
}

func NewWording(translator Translator) *Wording {
	return &Wording{
		translator: translator,
	}
}

func (w *Wording) Transform(ctx context.Context, err error) string {
	switch et := xerrs.Cause(err).(type) {
	default:
		return err.Error()
	case *Error:
		field, exist := validate.FieldFromContext(ctx)
		if !exist {
			return err.Error()
		}

		args := struct {
			Field string
			Args  map[string]interface{}
		}{
			Field: field,
			Args:  et.Args(),
		}

		messageID := err.Error()
		language := ""
		translated, err := w.translator.Translate(language, messageID, args)
		if err != nil {
			return fmt.Sprintf("message_id '%s' is not found at language '%s'", messageID, language)
		}

		return translated
	}
}
