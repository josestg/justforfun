package x

import "time"

// LocalClock is an OO-like time.Now() with location.
type LocalClock struct {
	loc *time.Location
}

// NewLocalClock creates a new LocalClock.
func NewLocalClock(loc *time.Location) *LocalClock {
	return &LocalClock{
		loc: loc,
	}
}

// Now returns the current time at some location.
func (l *LocalClock) Now() time.Time {
	return time.Now().In(l.loc)
}
