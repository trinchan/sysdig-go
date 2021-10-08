package sysdig

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUsersService_Me(t *testing.T) {
	methodName := "Me"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/user/", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"user":{"id":1}}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			user, _, err := client.Users.Me(ctx)
			if err != nil {
				t.Errorf("Users.Me returned error: %v", err)
			}
			want := &MeResponse{User: User{ID: 1}}
			if !cmp.Equal(user, want) {
				t.Errorf("Users.Me returned %+v, want %+v", user, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Users.Me(context.Background())
		return resp, err
	})
}

func TestUsersService_Token(t *testing.T) {
	methodName := "Token"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/token", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"token":{"key":"1"}}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			user, _, err := client.Users.Token(ctx)
			if err != nil {
				t.Errorf("Users.Token returned error: %v", err)
			}
			want := &TokenResponse{Token: Token{Key: "1"}}
			if !cmp.Equal(user, want) {
				t.Errorf("Users.Token returned %+v, want %+v", user, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Users.Me(context.Background())
		return resp, err
	})
}

func TestUsersService_ConnectedAgents(t *testing.T) {
	methodName := "ConnectedAgents"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/agents/", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"total":1,"agents":[{"id":"1"}]}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			user, _, err := client.Users.ConnectedAgents(ctx)
			if err != nil {
				t.Errorf("Users.Token returned error: %v", err)
			}
			want := &ConnectedAgentsResponse{Total: 1, Agents: []Agent{{ID: "1"}}}
			if !cmp.Equal(user, want) {
				t.Errorf("Users.Token returned %+v, want %+v", user, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Users.Me(context.Background())
		return resp, err
	})
}
