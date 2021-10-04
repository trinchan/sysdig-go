package sysdig

import (
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/trinchan/sysdig-go/sysdig/authentication"
	"github.com/trinchan/sysdig-go/sysdig/authentication/accesstoken"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints. See issue #752.
	baseURLPath = "/api"
)

// setup sets up a test HTTP server along with a sysdig.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup(options ...ClientOption) (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// We want to ensure that tests catch mistakes where the endpoint URL is
	// specified as absolute rather than relative.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the Sysdig client being tested and is
	// configured to use test server.
	client, _ = NewClient(WithDebug(true))
	for _, o := range options {
		_ = o(client)
	}
	url, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("request method: %v, want %v", got, want)
	}
}

// Test how bad options are handled. Method f under test should
// return an error.
func testBadOptions(t *testing.T, methodName string, f func() error) {
	t.Helper()
	if methodName == "" {
		t.Error("testBadOptions: must supply method methodName")
	}
	if err := f(); err == nil {
		t.Errorf("bad options %v err = nil, want error", methodName)
	}
}

// Test function under NewRequest failure and then s.client.Do failure.
// Method f should be a regular call that would normally succeed, but
// should return an error when NewRequest or s.client.Do fails.
func testNewRequestAndDoFailure(t *testing.T, methodName string, client *Client, f func() (*http.Response, error)) {
	t.Helper()
	if methodName == "" {
		t.Error("testNewRequestAndDoFailure: must supply method methodName")
	}

	client.BaseURL.Path = ""
	resp, err := f()
	if resp != nil {
		t.Errorf("client.BaseURL.Path='' %v resp = %#v, want nil", methodName, resp)
	}
	if err == nil {
		t.Errorf("client.BaseURL.Path='' %v err = nil, want error", methodName)
	}
}

func TestClientOptions(t *testing.T) {
	tests := []struct {
		name    string
		option  ClientOption
		wantErr bool
	}{
		{
			name:    "WithLogger",
			option:  WithLogger(noopLog),
			wantErr: false,
		},
		{
			name:    "WithHTTPClient",
			option:  WithHTTPClient(http.DefaultClient),
			wantErr: false,
		},
		{
			name:    "WithBaseURL",
			option:  WithBaseURL(defaultBaseURL),
			wantErr: false,
		},
		{
			name:    "WithBaseURL_Error",
			option:  WithBaseURL("https://:123:weird:url"),
			wantErr: true,
		},
		{
			name:    "WithIBMBaseURL_Private",
			option:  WithIBMBaseURL(RegionUSSouth, true),
			wantErr: false,
		},
		{
			name:    "WithIBMBaseURL_Public",
			option:  WithIBMBaseURL(RegionUSSouth, false),
			wantErr: false,
		},
		{
			name:    "WithUserAgent",
			option:  WithUserAgent(userAgent),
			wantErr: false,
		},
		{
			name:    "WithNoAuthenticator",
			option:  WithAuthenticator(nil),
			wantErr: false,
		},
		{
			name: "WithAuthenticator",
			option: WithAuthenticator(func() authentication.Authenticator {
				a, err := accesstoken.Authenticator("foo")
				if err != nil {
					t.Fatal(err)
				}
				return a
			}()),
			wantErr: false,
		},
		{
			name:    "WithResponseCompression",
			option:  WithResponseCompression(true),
			wantErr: false,
		},
		{
			name:    "WithDebug",
			option:  WithDebug(false),
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := NewClient(test.option)
			if (err != nil) != test.wantErr {
				t.Errorf("got err: %v, want err: %v", err, test.wantErr)
			}
		})
	}
}

func TestClientCopy(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	c2 := c.Client()
	if c.httpClient == c2 {
		t.Error("Client returned same http.Client, but should be different")
	}
}

func TestSetLogger(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	l := log.Default()
	c.SetLogger(l)
	if c.logger != l {
		t.Error("Client returned different logger, but should be same")
	}
}

func TestAuthentication(t *testing.T) {
	a, err := accesstoken.Authenticator("foo")
	if err != nil {
		t.Fatal(err)
	}
	client, mux, baseURL, teardown := setup(WithAuthenticator(a))
	defer teardown()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/foo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := client.BareDo(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to baredo: %v", err)
	}
	defer resp.Body.Close()
}

func TestClientResponses(t *testing.T) {
	client, mux, baseURL, teardown := setup()
	defer teardown()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/foo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	resp, err := client.BareDo(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to baredo: %v", err)
	}
	defer resp.Body.Close()
}

func TestBareDo_Zipped(t *testing.T) {
	client, mux, baseURL, teardown := setup()
	defer teardown()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "text/plain; gzip; charset=utf8")
		w.WriteHeader(http.StatusOK)
		respW := gzip.NewWriter(w)
		if _, err := respW.Write([]byte("ok")); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	})
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/foo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = client.BareDo(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to baredo: %v", err)
	}
}

func TestBareDo_AuthenticationError(t *testing.T) {
	a, err := accesstoken.Authenticator("foo")
	if err != nil {
		t.Fatal(err)
	}
	client, mux, baseURL, teardown := setup(WithAuthenticator(a))
	defer teardown()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/foo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = client.BareDo(context.Background(), req)
	if err == nil {
		t.Errorf("did not get expected err")
	}
}

type refreshableAuthenticationWrapper struct {
	authentication.Authenticator
	Refresher func() error
}

func (r *refreshableAuthenticationWrapper) Refresh() error {
	return r.Refresher()
}

func TestBareDo_AuthenticationRefreshable(t *testing.T) {
	a, err := accesstoken.Authenticator("foo")
	if err != nil {
		t.Fatal(err)
	}
	hit := 0
	client, mux, baseURL, teardown := setup(WithAuthenticator(&refreshableAuthenticationWrapper{
		Authenticator: a,
		Refresher:     func() error { hit++; return nil },
	}))
	defer teardown()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		if hit == 0 {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/foo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	_, err = client.BareDo(context.Background(), req)
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	if hit == 0 {
		t.Errorf("did not run refresher")
	}
}

func TestBareDo_DoError(t *testing.T) {
	client, mux, baseURL, teardown := setup()
	defer teardown()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/foo", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.URL = nil
	_, err = client.BareDo(context.Background(), req)
	if err == nil {
		t.Errorf("did not get expected err")
	}
}
