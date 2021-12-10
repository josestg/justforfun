package jwt

import (
	"reflect"
	"testing"
	"time"
)

func TestNewTime(t *testing.T) {

	t1 := NewTime(time.Now())
	t2 := new(Time)

	b, err := t1.MarshalJSON()
	if err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	if err := t2.UnmarshalJSON(b); err != nil {
		t.Fatalf("expecting error nil but got %v", err)
	}

	if !reflect.DeepEqual(t1, t2) {
		t.Fatalf("expecting t1 and t2 are equal")
	}

}
