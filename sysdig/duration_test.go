package sysdig

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

func TestNewMicroDuration(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want time.Duration
	}{
		{
			name: "1 second",
			in:   time.Second,
			want: time.Second,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewMicroDuration(test.in)
			if got.Duration != test.in {
				t.Errorf("got duration: %v, want: %v", got.Duration, test.in)
			}
		})
	}
}

func TestMicroDuration_JSONMarshal(t *testing.T) {
	tests := []struct {
		name string
		in   MicroDuration
		want []byte
	}{
		{
			name: "1 millisecond",
			in:   NewMicroDuration(time.Millisecond),
			want: []byte(strconv.Itoa(int(time.Millisecond.Microseconds()))),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.in)
			if err != nil {
				t.Fatalf("unexpected json marshal error: %v", err)
			}
			if !bytes.Equal(got, test.want) {
				t.Errorf("got: %s, want: %s", string(got), string(test.want))
			}
		})
	}
}

func TestMicroDuration_JSONUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		in      []byte
		want    MicroDuration
		wantErr bool
	}{
		{
			name:    "1 millisecond",
			in:      []byte(strconv.Itoa(int(time.Millisecond.Microseconds()))),
			want:    NewMicroDuration(time.Millisecond),
			wantErr: false,
		},
		{
			name:    "invalid",
			in:      []byte("not a duration"),
			want:    NewMicroDuration(time.Duration(0)),
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got MicroDuration
			err := json.Unmarshal(test.in, &got)
			if test.wantErr != (err != nil) {
				t.Fatalf("unexpected json marshal error: %v", err)
			}
			if got != test.want {
				t.Errorf("got: %s, want: %s", got, test.want)
			}
		})
	}
}
