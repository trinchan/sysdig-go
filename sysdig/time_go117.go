//go:build go1.17
// +build go1.17

package sysdig

import (
	"time"
)

// UnixMilli converts the time to Unix milliseconds for Sysdig.
// Uses the built-in time.Time UnixMilli to implement in go >=1.17.
func (t *Time) UnixMilli() int64 {
	return t.Time.UnixMilli()
}

// UnixMilli returns the local Time corresponding to the given Unix time in milliseconds.
// Uses the built-in time.UnixMilli to implement in go >=1.17.
func UnixMilli(msec int64) Time {
	return *NewTime(time.UnixMilli(msec))
}
