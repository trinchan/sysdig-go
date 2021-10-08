package sysdig

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTeamsService_Get(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/team/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		id      int
		handler http.HandlerFunc
		want    *TeamResponse
	}{
		{
			name: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"team":{"id":1}}`)
			},
			want: &TeamResponse{Team: Team{ID: 1}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			team, _, err := client.Teams.Get(ctx, test.id)
			if err != nil {
				t.Errorf("Teams.Get returned error: %v", err)
			}
			if !cmp.Equal(team, test.want) {
				t.Errorf("Teams.Get returned %+v, want %+v", team, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Users.Me(context.Background())
		return resp, err
	})
}

func TestTeamsService_List(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/team", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name        string
		handler     http.HandlerFunc
		productType ProductType
		want        *ListTeamsResponse
	}{
		{
			name: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"teams":[{"id":1}]}`)
			},
			productType: ProductTypeAny,
			want:        &ListTeamsResponse{Teams: []Team{{ID: 1}}},
		},
		{
			name: "test sdc",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				testFormValues(t, r, values{"product": string(ProductTypeSDC)})
				fmt.Fprint(w, `{"teams":[{"id":1}]}`)
			},
			productType: ProductTypeSDC,
			want:        &ListTeamsResponse{Teams: []Team{{ID: 1}}},
		},
		{
			name: "test sds",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				testFormValues(t, r, values{"product": string(ProductTypeSDS)})
				fmt.Fprint(w, `{"teams":[{"id":1}]}`)
			},
			productType: ProductTypeSDS,
			want:        &ListTeamsResponse{Teams: []Team{{ID: 1}}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			team, _, err := client.Teams.List(ctx, test.productType)
			if err != nil {
				t.Errorf("Teams.List returned error: %v", err)
			}
			if !cmp.Equal(team, test.want) {
				t.Errorf("Teams.List returned %+v, want %+v", team, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Users.Me(context.Background())
		return resp, err
	})
}

func TestTeamsService_ListUsers(t *testing.T) {
	methodName := "Get"
	teamID := 1
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc(fmt.Sprintf("/api/team/%d/users", teamID), func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    *ListUsersResponse
	}{
		{
			name: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"total":2,"offset":1,"users":[{"id":1},{"id":2}]}`)
			},
			want: &ListUsersResponse{
				Offset: 1,
				Total:  2,
				Users:  []User{{ID: 1}, {ID: 2}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			team, _, err := client.Teams.ListUsers(ctx, teamID)
			if err != nil {
				t.Errorf("Teams.ListUsers returned error: %v", err)
			}
			if !cmp.Equal(team, test.want) {
				t.Errorf("Teams.ListUsers returned %+v, want %+v", team, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Users.Me(context.Background())
		return resp, err
	})
}

func TestTeamsService_Delete(t *testing.T) {
	methodName := "Delete"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/team/", func(w http.ResponseWriter, r *http.Request) {
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
			_, err := client.Teams.Delete(ctx, test.id)
			if err != nil {
				t.Errorf("Teams.Delete returned error: %v", err)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		return client.Teams.Delete(context.Background(), tests[0].id)
	})
}

func TestTeamsService_Infrastructure(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/team/infrastructure", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		want    *InfrastructureResponse
	}{
		{
			name: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				fmt.Fprint(w, `{"infrastructure":{"metricCount":{"total":1}}}`)
			},
			want: &InfrastructureResponse{Infrastructure: Infrastructure{
				MetricCount: InfrastructureMetricCount{
					Total: 1,
				},
			}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			team, _, err := client.Teams.Infrastructure(ctx)
			if err != nil {
				t.Errorf("Teams.Infrastructure returned error: %v", err)
			}
			if !cmp.Equal(team, test.want) {
				t.Errorf("Teams.Infrastructure returned %+v, want %+v", team, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Teams.Infrastructure(context.Background())
		return resp, err
	})
}
