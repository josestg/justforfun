package locale

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
	"sync"
	"text/template"

	"github.com/josestg/justforfun/pkg/xerrs"
)

//go:embed bundle
var fSys embed.FS

type Translation map[string]string

type dictionary struct {
	bundle map[string]Translation
}

func (d *dictionary) Translate(lang, messageID string, args interface{}) (string, error) {
	local, exist := d.bundle[lang]
	if !exist {
		return "", fmt.Errorf("unkown language '%s'", lang)
	}

	tpl, exist := local[messageID]
	if !exist {
		return "", fmt.Errorf("message id '%s' not found at language '%s'", messageID, lang)
	}

	t, err := template.New(messageID).Parse(tpl)
	if err != nil {
		return "", xerrs.Wrap(err, "parsing template")
	}
	buf := &bytes.Buffer{}
	err = t.Execute(buf, args)
	if err != nil {
		return "", xerrs.Wrap(err, "executing template")
	}

	return buf.String(), nil
}

var Dictionary = &dictionary{}
var initiator sync.Once

func init() {
	initiator.Do(func() {
		bundle := make(map[string]Translation)
		err := fs.WalkDir(fSys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			if !strings.HasSuffix(d.Name(), ".json") {
				return nil
			}

			b, err := fSys.ReadFile(d.Name())
			if err != nil {
				return err
			}

			var t Translation
			if err := json.Unmarshal(b, &t); err != nil {
				return err
			}

			lang := strings.TrimSuffix(d.Name(), ".json")
			bundle[lang] = t
			return nil
		})

		if err != nil {
			panic(err)
		}

		Dictionary.bundle = bundle
	})
}
