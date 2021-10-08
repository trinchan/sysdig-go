package sysdig

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDashboardsService_List(t *testing.T) {
	methodName := "List"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"dashboards":[{"id":1}]}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			dashboard, _, err := client.Dashboards.List(ctx)
			if err != nil {
				t.Errorf("Dashboards.List returned error: %v", err)
			}
			want := &ListDashboardsResponse{
				Dashboards: []Dashboard{{ID: 1}},
			}
			if !cmp.Equal(dashboard, want) {
				t.Errorf("Dashboards.List returned %+v, want %+v", dashboard, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, ferr := client.Dashboards.List(context.Background())
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, ferr
	})
}

func TestDashboardsService_Delete(t *testing.T) {
	methodName := "Delete"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards/", func(w http.ResponseWriter, r *http.Request) {
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
			_, _, err := client.Dashboards.Delete(ctx, test.id)
			if err != nil {
				t.Errorf("Dashboards.Delete returned error: %v", err)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Dashboards.Delete(context.Background(), tests[0].id)
		return resp, err
	})
}

func TestDashboardsService_Get(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards/", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprint(w, `{"dashboard":{"id":1}}`)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			dashboard, _, err := client.Dashboards.Get(ctx, 1)
			if err != nil {
				t.Errorf("Dashboards.Get returned error: %v", err)
			}
			want := &DashboardResponse{Dashboard: Dashboard{ID: 1}}
			if !cmp.Equal(dashboard, want) {
				t.Errorf("Dashboards.Get returned %+v, want %+v", dashboard, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Dashboards.Get(context.Background(), 1)
		return resp, err
	})
}

func TestNewDashboard(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"test",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewDashboard(test.name)
			if got.Name != test.name {
				t.Errorf("got name: %s, want: %s", got.Name, test.name)
			}
			if got.Schema != 3 {
				t.Errorf("got version: %d, want: 3", got.Version)
			}
		})
	}
}

func TestDashboardService_Create(t *testing.T) {
	methodName := "Create"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	type dashboardRequest struct {
		Dashboard Dashboard `json:"dashboard"`
	}

	handlerFor := func(options Dashboard, output string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, http.MethodPost)
			var v dashboardRequest
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				t.Errorf("failed to decode test case: %v", err)
				t.FailNow()
			}
			options.Schema = 3
			if options.Panels == nil {
				options.Panels = make([]Panel, 0)
			}
			if options.Layout == nil {
				options.Layout = make([]Layout, 0)
			}
			ereq := dashboardRequest{options}
			if !cmp.Equal(v, ereq) {
				t.Errorf("Request body = %+v, want %+v", v, ereq)
			}
			fmt.Fprint(w, output)
		}
	}

	tests := []struct {
		name    string
		options Dashboard
		output  string
		want    *DashboardResponse
	}{
		{
			name:    "test",
			options: Dashboard{Name: "test"},
			output:  `{"dashboard":{"id":1,"schema":3}}`,
			want:    &DashboardResponse{Dashboard: Dashboard{ID: 1, Schema: 3}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = handlerFor(test.options, test.output)
			ctx := context.Background()
			event, _, err := client.Dashboards.Create(ctx, test.options)
			if err != nil {
				t.Errorf("Dashboards.Create returned error: %v", err)
			}

			if !cmp.Equal(event, test.want) {
				t.Errorf("Dashboards.Create returned %+v, want %+v", event, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		got, resp, derr := client.Dashboards.Create(context.Background(), tests[0].options)
		if got != nil {
			t.Errorf("testNewRequestAndDoFailure %v = %#v, want nil", methodName, got)
		}
		return resp, derr
	})
}

func TestDashboardsService_Update(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		options Dashboard
		handler http.HandlerFunc
	}{
		{
			name: "test",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPut)
				fmt.Fprint(w, `{"dashboard":{"id":1,"schema":3}}`)
			},
			options: Dashboard{Name: "test", Schema: 3, Panels: []Panel{}, Layout: []Layout{}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			dashboard, _, err := client.Dashboards.Update(ctx, test.options)
			if err != nil {
				t.Errorf("Dashboards.Update returned error: %v", err)
			}
			want := &DashboardResponse{Dashboard: Dashboard{ID: 1, Schema: 3}}
			if !cmp.Equal(dashboard, want) {
				t.Errorf("Dashboards.Update returned %+v, want %+v", dashboard, want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Dashboards.Update(context.Background(), tests[0].options)
		return resp, err
	})
}

func TestDashboardsService_Favorite(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name     string
		id       int
		favorite bool
		handler  http.HandlerFunc
		want     *DashboardResponse
	}{
		{
			name: "test favorite",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPatch)
				fmt.Fprint(w, `{"dashboard":{"id":1,"schema":3,"favorite":true}}`)
			},
			id:       1,
			favorite: true,
			want:     &DashboardResponse{Dashboard: Dashboard{ID: 1, Schema: 3, Favorite: true}},
		},
		{
			name: "test unfavorite",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPatch)
				fmt.Fprint(w, `{"dashboard":{"id":1,"schema":3,"favorite":false}}`)
			},
			id:       1,
			favorite: false,
			want:     &DashboardResponse{Dashboard: Dashboard{ID: 1, Schema: 3, Favorite: false}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			dashboard, _, err := client.Dashboards.Favorite(ctx, test.id, test.favorite)
			if err != nil {
				t.Errorf("Dashboards.Favorite returned error: %v", err)
			}
			if !cmp.Equal(dashboard, test.want) {
				t.Errorf("Dashboards.Favorite returned %+v, want %+v", dashboard, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Dashboards.Favorite(context.Background(), tests[0].id, tests[0].favorite)
		return resp, err
	})
}

func TestDashboardsService_Transfer(t *testing.T) {
	methodName := "Get"
	client, mux, _, teardown := setup()
	var h http.HandlerFunc
	mux.HandleFunc("/api/v3/dashboards/", func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	})
	defer teardown()

	tests := []struct {
		name    string
		handler http.HandlerFunc
		ids     []int
		want    *DashboardTransferResponse
		wantErr bool
	}{
		{
			name: "test transfer",
			handler: func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPost)
				fmt.Fprint(w, `{"results":{"id":1}}`)
			},
			ids:     []int{1},
			want:    &DashboardTransferResponse{DashboardTransferResults{ID: 1}},
			wantErr: false,
		},
		{
			name:    "test no ids",
			handler: nil,
			ids:     nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h = test.handler
			ctx := context.Background()
			dashboard, _, err := client.Dashboards.Transfer(ctx, 0, 0, true, test.ids...)
			if test.wantErr != (err != nil) {
				t.Errorf("Dashboards.Transfer returned error: %v", err)
			}
			if !cmp.Equal(dashboard, test.want) {
				t.Errorf("Dashboards.Transfer returned %+v, want %+v", dashboard, test.want)
			}
		})
	}

	testNewRequestAndDoFailure(t, methodName, client, func() (*http.Response, error) {
		_, resp, err := client.Dashboards.Transfer(context.Background(), 0, 0, true, 0)
		return resp, err
	})
}
