package sysdig

import (
	"encoding/json"
	"strconv"
	"time"
)

// MilliTime is a custom time.Time which implements transforming between time.Time's default representation and Sysdig's
// expected time.UnixMillis representation.
type MilliTime struct {
	time.Time
}

// NewMilliTime creates a new MilliTime with the provided time.Time.
func NewMilliTime(t time.Time) MilliTime {
	return MilliTime{t}
}

// MarshalJSON implements the json.Marshaler interface for MilliTime.
func (t MilliTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.UnixMilli())
}

// UnmarshalJSON implements json.Unmarshaler for MilliTime.
func (t *MilliTime) UnmarshalJSON(b []byte) error {
	u, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = UnixMilli(int64(u))
	return nil
}
