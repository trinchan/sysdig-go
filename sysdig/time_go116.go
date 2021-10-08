//go:build !go1.17
// +build !go1.17

package sysdig

import "time"

// UnixMilli converts the time to Unix milliseconds for Sysdig.
func (t *MilliTime) UnixMilli() int64 {
	return t.Unix()*1e3 + int64(t.Nanosecond())/1e6
}

// UnixMilli returns the local MilliTime corresponding to the given Unix time,
func UnixMilli(msec int64) MilliTime {
	return NewMilliTime(time.Unix(msec/1e3, (msec%1e3)*1e6))
}
