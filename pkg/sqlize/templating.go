package sqlize

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const migrationTemplate = `
-- up script here...

---+split+---

-- down script here...
`

const actionScriptSeparator = "---+split+---"

type templating struct{}

func (t *templating) Template() string {
	return migrationTemplate
}

func (t *templating) Read(r io.Reader, action Action) (string, error) {
	buf := strings.Builder{}
	if _, err := io.Copy(&buf, r); err != nil {
		return "", fmt.Errorf("%w: read migration template", err)
	}

	parts := strings.SplitN(buf.String(), actionScriptSeparator, 2)
	if len(parts) != 2 {
		return "", errors.New("invalid migration template")
	}

	switch action {
	default:
		return "", errors.New("unknown action")
	case MigrationUp:
		return strings.TrimSpace(parts[0]), nil
	case MigrationDown:
		return strings.TrimSpace(parts[1]), nil
	}
}
