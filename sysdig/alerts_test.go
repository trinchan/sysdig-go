package sysdig

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAlertsService_List(t *testing.T) {
	methodName := "List"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"alerts":[{"id":1}]}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			alert, _, err := client.Alerts.List(ctx)
			if err != nil {
				t.Errorf("Alerts.List returned error: %v", err)
			}
			want := &ListAlertConfigurationsResponse{
				Alerts: []Alert{{ID: 1}},
			}
			if !cmp.Equal(alert, want) {
				t.Errorf("Alerts.List returned %+v, want %+v", alert, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, ferr := client.Alerts.List(context.Background())
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, ferr
	})
}

func TestAlertsService_Delete(t *testing.T) {
	methodName := "Delete"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/alerts/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		id      int
		handler http.HandlerFunc
	}{
		{
			name: "test",
			id:   1,
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
			_, err := client.Alerts.Delete(ctx, test.id)
			if err != nil {
				t.Errorf("Alerts.Delete returned error: %v", err)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		return client.Alerts.Delete(context.Background(), tests[0].id)
	})
}

func TestAlertsService_Get(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/alerts/", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"alert":{"id":1}}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			alert, _, err := client.Alerts.Get(ctx, 1)
			if err != nil {
				t.Errorf("Alerts.Get returned error: %v", err)
			}
			want := &AlertResponse{Alert: Alert{ID: 1}}
			if !cmp.Equal(alert, want) {
				t.Errorf("Alerts.Get returned %+v, want %+v", alert, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Alerts.Get(context.Background(), 1)
		return resp, err
	})
}
