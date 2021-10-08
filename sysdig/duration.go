package sysdig

import (
	"encoding/json"
	"strconv"
	"time"
)

// MicroDuration is a custom time.Duration which implements transforming between time.Duration's
// default representation and Sysdig's expected Microsecond duration format for timespans.
type MicroDuration struct {
	time.Duration
}

// NewMicroDuration creates a new MicroDuration with the provided time.Duration.
func NewMicroDuration(t time.Duration) MicroDuration {
	return MicroDuration{t}
}

// MarshalJSON implements the json.Marshaler interface for MilliTime.
func (t MicroDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Microseconds())
}

// UnmarshalJSON implements json.Unmarshaler for MilliTime.
func (t *MicroDuration) UnmarshalJSON(b []byte) error {
	u, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	*t = NewMicroDuration(time.Duration(u) * time.Microsecond)
	return nil
}
