package sysdig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestEventsService_CreateEvent(t *testing.T) {
	methodName := "CreateEvent"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/v2/events", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	type eventRequest struct {
		Event EventOptions `json:"event"`
	}

	handlerFor := func(options EventOptions, output string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			var v eventRequest
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				t.Errorf("failed to decode test case: %v", err)
				t.FailNow()
			}
			ereq := eventRequest{options}
			if !cmp.Equal(v, ereq) {
				t.Errorf("Request body = %+v, want %+v", v, ereq)
			}
			fmt.Fprint(w, output)
		}
	}

	tests := []struct {
		name    string
		options EventOptions
		output  string
		want    *EventResponse
	}{
		{
			name:    "test",
			options: EventOptions{Name: "test"},
			output:  `{"event":{"id":"1"}}`,
			want:    &EventResponse{Event: Event{ID: "1"}},
		},
		{
			name:    "test w/ time",
			options: EventOptions{Name: "test", Timestamp: NewTime(time.UnixMilli(1))},
			output:  `{"event":{"id":"1","timestamp":1}}`,
			want:    &EventResponse{Event: Event{ID: "1", Timestamp: *NewTime(time.UnixMilli(1))}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = handlerFor(test.options, test.output)
			ctx := context.Background()
			event, _, err := client.Events.CreateEvent(ctx, &test.options)
			if err != nil {
				t.Errorf("Events.CreateEvent returned error: %v", err)
			}

			if !cmp.Equal(event, test.want) {
				t.Errorf("Events.CreateEvent returned %+v, want %+v", event, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, derr := client.Events.CreateEvent(context.Background(), &tests[0].options)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, derr
	})
}

func TestEventsService_ListEvents(t *testing.T) {
	methodName := "ListEvents"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/v2/events", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		options *ListEventOptions
		handler http.HandlerFunc
	}{
		{
			name:    "test",
			options: &ListEventOptions{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"total":1,"matched":1,"events":[{"id":"1"}]}`)
			},
		},
		{
			name:    "test",
			options: &ListEventOptions{
				Categories: Categories{CategoryCustom},
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"total":1,"matched":1,"events":[{"id":"1"}]}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			event, _, err := client.Events.ListEvents(ctx, test.options)
			if err != nil {
				t.Errorf("Events.ListEvents returned error: %v", err)
			}
			want := &ListEventsResponse{
				Total:   1,
				Matched: 1,
				Events:  []Event{{ID: "1"}},
			}
			if !cmp.Equal(event, want) {
				t.Errorf("Events.ListEvents returned %+v, want %+v", event, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, ferr := client.Events.ListEvents(context.Background(), tests[0].options)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, ferr
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Events.ListEvents(context.Background(), &ListEventOptions{
			Filter: "\n",
		})
		return err
	})

}

func TestEventsService_DeleteEvent(t *testing.T) {
	methodName := "DeleteEvent"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/v2/events/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		id      string
		handler http.HandlerFunc
	}{
		{
			name: "test",
			id:   "1",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodDelete)
				w.WriteHeader(http.StatusNoContent)
			},
		},
	}
	for _, test := range tests {
		h = test.handler
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := client.Events.DeleteEvent(ctx, test.id)
			if err != nil {
				t.Errorf("Events.DeleteEvent returned error: %v", err)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		return client.Events.DeleteEvent(context.Background(), tests[0].id)
	})

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.Events.DeleteEvent(context.Background(), "\n")
		return err
	})
}

func TestEventsService_GetEvent(t *testing.T) {
	methodName := "GetEvent"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/v2/events/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"event":{"id":"1"}}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			event, _, err := client.Events.GetEvent(ctx, "1")
			if err != nil {
				t.Errorf("Events.GetEvent returned error: %v", err)
			}
			want := &EventResponse{Event: Event{ID: "1"}}
			if !cmp.Equal(event, want) {
				t.Errorf("Events.GetEvent returned %+v, want %+v", event, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Events.GetEvent(context.Background(), tests[0].name)
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Events.GetEvent(context.Background(), "\n")
		return err
	})
}
