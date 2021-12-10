package env

import (
	"os"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	const (
		key     = "TESTING_ENV_STRING"
		initial = "STRING INITIAL VALUE"
		val     = "STRING VALUE"
	)

	if got := String(key, initial); got != initial {
		t.Errorf("expecting using the initial value")
	}

	if err := os.Setenv(key, val); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	if got := String(key, initial); got != val {
		t.Errorf("expcting using the env value")
	}

}

func TestInt(t *testing.T) {
	const (
		key       = "TESTING_ENV_INT"
		initial   = 1
		val       = 2
		valString = "2"
	)

	if got := Int(key, initial); got != initial {
		t.Errorf("expecting using the initial value")
	}

	if err := os.Setenv(key, valString); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	if got := Int(key, initial); got != val {
		t.Errorf("expcting using the env value")
	}

	// Test: Expecting panic if the env value not a number.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expecting panic")
		}
	}()

	if err := os.Setenv(key, "INVALID NUMBER FORMAT"); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	// this should be panic.
	_ = Int(key, initial)

}

func TestDuration(t *testing.T) {
	const (
		key       = "TESTING_ENV_DURATION"
		initial   = 1 * time.Second
		val       = 2 * time.Second
		valString = "2s"
	)

	if got := Duration(key, initial); got != initial {
		t.Errorf("expecting using the initial value")
	}

	if err := os.Setenv(key, valString); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	if got := Duration(key, initial); got != val {
		t.Errorf("expcting using the env value")
	}

	// Test: Expecting panic if the env value not a duration.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expecting panic")
		}
	}()

	if err := os.Setenv(key, "INVALID DURATION FORMAT"); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	// this should be panic.
	_ = Duration(key, initial)
}

func TestBool(t *testing.T) {
	const (
		key       = "TESTING_ENV_BOOL"
		initial   = true
		val       = false
		valString = "false"
	)

	if got := Bool(key, initial); got != initial {
		t.Errorf("expecting using the initial value")
	}

	if err := os.Setenv(key, valString); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	if got := Bool(key, initial); got != val {
		t.Errorf("expcting using the env value")
	}

	// Test: Expecting panic if the env value not a duration.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expecting panic")
		}
	}()

	if err := os.Setenv(key, "INVALID BOOLEAN FORMAT"); err != nil {
		t.Errorf("expecting error nil but got %v", err)
	}

	// this should be panic.
	_ = Bool(key, initial)
}
