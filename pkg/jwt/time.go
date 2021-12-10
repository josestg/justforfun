package jwt

import (
	"encoding/json"
	"time"
)

// Time represents a JWT time.
// This Time truncates the time.Time with time.Millisecond.
//
// Time overrides the MarshalJSON and UnmarshalJSON of time.Time
// to make the marshaled time as plaintext number instead
// of formatted-string like time.RFC3339.
type Time struct {
	time.Time
}

// NewTime creates a new time at given time.
func NewTime(at time.Time) *Time {
	return &Time{at.Truncate(time.Millisecond)}
}

func (t *Time) MarshalJSON() ([]byte, error) {
	ms := t.Truncate(time.Millisecond).UnixMilli()
	return json.Marshal(ms)
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var ms json.Number
	if err := json.Unmarshal(b, &ms); err != nil {
		return err
	}

	msi64, err := ms.Int64()
	if err != nil {
		return err
	}

	mst := time.UnixMilli(msi64)
	*t = Time{mst}
	return nil
}
