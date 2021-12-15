package sqlize

import (
	"strings"
	"testing"
)

func TestTemplating_Read(t *testing.T) {
	tpl := DefaultTemplating.Template()
	script, err := DefaultTemplating.Read(strings.NewReader(tpl), MigrationUp)
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	const upScriptExpected = `-- up script here...`

	if script != upScriptExpected {
		t.Fatalf("expecting script %v but got %v", upScriptExpected, script)
	}

	const downScriptExpected = `-- down script here...`
	script, err = DefaultTemplating.Read(strings.NewReader(tpl), MigrationDown)
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	if script != downScriptExpected {
		t.Fatalf("expecting script %v but got %v", downScriptExpected, script)
	}

	script, err = DefaultTemplating.Read(strings.NewReader(tpl), 111)
	if err == nil {
		t.Fatalf("expecting got an error")
	}

}
