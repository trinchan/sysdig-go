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

func TestNotificationChannelsService_Create(t *testing.T) {
	methodName := "Create"
	client, mux, _, teardown := setup(nil)
	var h http.HandlerFunc
	mux.HandleFunc("/api/notificationChannels", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	type notificationChannelRequest struct {
		NotificationChannel NotificationChannel `json:"notificationChannel"`
	}

	handlerFor := func(options NotificationChannel, output string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			var v notificationChannelRequest
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				t.Errorf("failed to decode test case: %v", err)
				t.FailNow()
			}
			ereq := notificationChannelRequest{options}
			if !cmp.Equal(v, ereq) {
				t.Errorf("Request body = %+v, want %+v", v, ereq)
			}
			fmt.Fprint(w, output)
		}
	}

	tests := []struct {
		name                    string
		notificationChannelType NotificationChannelType
		options                 NotificationChannelOptions
		output                  string
		want                    *NotificationChannelResponse
	}{
		{
			name:                    "test",
			notificationChannelType: NotificationChannelTypeEmail,
			options:                 NotificationChannelOptions{},
			output:                  `{"notificationChannel":{"id":"1","type":"EMAIL"}}`,
			want: &NotificationChannelResponse{
				NotificationChannel: NotificationChannel{
					ID:   "1",
					Type: NotificationChannelTypeEmail,
				},
			},
		},
		{
			name:                    "test w/ options",
			notificationChannelType: NotificationChannelTypeSlack,
			options: NotificationChannelOptions{
				EmailRecipients: []string{"example@foo.com"},
			},
			output: `{"notificationChannel":{"id":"1","type":"SLACK","options":{"emailRecipients": ["example@foo.com"]}}}`,
			want: &NotificationChannelResponse{
				NotificationChannel: NotificationChannel{
					Type: NotificationChannelTypeSlack,
					ID:   "1",
					Options: NotificationChannelOptions{
						EmailRecipients: []string{"example@foo.com"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ch := NotificationChannel{
				Type:    test.notificationChannelType,
				Name:    test.name,
				Options: test.options,
				Enabled: true,
			}
			h = handlerFor(ch, test.output)
			ctx := context.Background()
			channel, _, err := client.NotificationChannels.Create(ctx, test.notificationChannelType, test.name, test.options)
			if err != nil {
				t.Errorf("NotificationChannels.Create returned error: %v", err)
			}

			if !cmp.Equal(channel, test.want) {
				t.Errorf("NotificationChannels.Create returned %+v, want %+v", channel, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, derr := client.NotificationChannels.Create(
			context.Background(),
			tests[0].notificationChannelType, tests[0].name, tests[0].options)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, derr
	})
}

func TestNotificationChannelsService_List(t *testing.T) {
	methodName := "List"
	client, mux, _, teardown := setup(nil)
	var h http.HandlerFunc
	mux.HandleFunc("/api/notificationChannels", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"notificationChannels":[{"id":"1"}]}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			channel, _, err := client.NotificationChannels.List(ctx, NewMilliTime(time.Now()), NewMilliTime(time.Now()))
			if err != nil {
				t.Errorf("NotificationChannels.List returned error: %v", err)
			}
			want := &ListNotificationChannelsResponse{
				NotificationChannels: []NotificationChannel{{ID: "1"}},
			}
			if !cmp.Equal(channel, want) {
				t.Errorf("NotificationChannels.List returned %+v, want %+v", channel, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, ferr := client.NotificationChannels.List(context.Background(), NewMilliTime(time.Now()), NewMilliTime(time.Now()))
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, ferr
	})
}

func TestNotificationChannelsService_Delete(t *testing.T) {
	methodName := "Delete"
	client, mux, _, teardown := setup(nil)
	var h http.HandlerFunc
	mux.HandleFunc("/api/notificationChannels/", func(w http.ResponseWriter, r *http.Request) {
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
			_, err := client.NotificationChannels.Delete(ctx, test.id)
			if err != nil {
				t.Errorf("NotificationChannels.Delete returned error: %v", err)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		return client.NotificationChannels.Delete(context.Background(), tests[0].id)
	})

	testBadOptions(t, methodName, func() (err error) {
		_, err = client.NotificationChannels.Delete(context.Background(), "\n")
		return err
	})
}

func TestNotificationChannelsService_Get(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup(nil)
	var h http.HandlerFunc
	mux.HandleFunc("/api/notificationChannels/", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"notificationChannel":{"id":"1"}}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			channel, _, err := client.NotificationChannels.Get(ctx, "1")
			if err != nil {
				t.Errorf("NotificationChannels.Get returned error: %v", err)
			}
			want := &NotificationChannelResponse{NotificationChannel: NotificationChannel{ID: "1"}}
			if !cmp.Equal(channel, want) {
				t.Errorf("NotificationChannels.Get returned %+v, want %+v", channel, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.NotificationChannels.Get(context.Background(), tests[0].name)
		return resp, err
	})

	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.NotificationChannels.Get(context.Background(), "\n")
		return err
	})
}
