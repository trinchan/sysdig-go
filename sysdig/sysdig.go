package sysdig

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/trinchan/sysdig-go/sysdig/authentication"

	"github.com/google/go-querystring/query"
)

// Logger is the interface for logging used by the Client.
type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}

var noopLog = &noopLogger{}

type noopLogger struct{}

func (l *noopLogger) Print(args ...interface{})                 {}
func (l *noopLogger) Printf(format string, args ...interface{}) {}

const (
	defaultBaseURL = "https://app.sysdigcloud.com/"
	ibmBaseURL     = "monitoring.cloud.ibm.com/"

	userAgent = "sysdig-go"
)

// Region is a type for defining available IBM regions for Sysdig.
type Region string

const (
	// RegionUSSouth is the IBM region us-south, located in Dallas.
	RegionUSSouth Region = "us-south"
	// RegionEUDE is the IBM region eu-de, located in Frankfurt.
	RegionEUDE Region = "eu-de"
	// RegionJPOSA is the IBM region jp-osa, located in Osaka.
	RegionJPOSA Region = "jp-osa"
	// RegionJPTOK is the IBM region jp-tok, located in Tokyo.
	RegionJPTOK Region = "jp-tok"
	// RegionUSEast is the IBM region us-east, located in Washington, DC.
	RegionUSEast Region = "us-east"
	// RegionAUSYD is the IBM region au-syd, located in Sydney.
	RegionAUSYD Region = "au-syd"
	// RegionCATOR is the IBM region ca-tor, located in Toronto.
	RegionCATOR Region = "ca-tor"
	// RegionBRSAO is the IBM region br-sao, located in São Paulo.
	RegionBRSAO Region = "br-sao"
)

type service struct {
	client *Client
}

// Client manages communication with the Sysdig API.
type Client struct {
	// Base URL for API requests. Defaults to the public Sysdig API, but can be
	// set to a domain endpoint to use with IBM or on-premise. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL
	// User agent used when communicating with the Sysdig API.
	UserAgent string

	httpClient             *http.Client // HTTP client used to communicate with the API.
	logger                 Logger
	debug                  bool
	shouldCompressResponse bool
	authenticator          authentication.Authenticator

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	// Services used for talking to different parts of the Sysdig API.
	Events               *EventsService
	Users                *UsersService
	NotificationChannels *NotificationChannelsService
	Alerts               *AlertService
	Dashboards           *DashboardService
	Teams                *TeamsService
}

// ClientOption defines the options for a Sysdig Client.
type ClientOption func(*Client) error

// NewClient creates a new Sysdig Client and applies all provided ClientOption.
func NewClient(options ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{
		BaseURL:    baseURL,
		UserAgent:  userAgent,
		httpClient: http.DefaultClient,
		logger:     noopLog,
	}
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}
	c.common.client = c
	c.Events = (*EventsService)(&c.common)
	c.Users = (*UsersService)(&c.common)
	c.NotificationChannels = (*NotificationChannelsService)(&c.common)
	c.Alerts = (*AlertService)(&c.common)
	c.Dashboards = (*DashboardService)(&c.common)
	c.Teams = (*TeamsService)(&c.common)
	return c, nil
}

// WithHTTPClient sets the HTTP client for the Sysdig client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = client
		return nil
	}
}

// WithBaseURL sets the Client.BaseURL to the provided URL. BaseURLs should have a trailing slash.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		url, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.BaseURL = url
		return nil
	}
}

// WithIBMBaseURL sets the Client.BaseURL to the BaseURL associated with the provided IBM Region and network.
func WithIBMBaseURL(ibmRegion Region, privateEndpoint bool) ClientOption {
	return func(c *Client) error {
		var err error
		if privateEndpoint {
			c.BaseURL, err = url.Parse(fmt.Sprintf("https://%s.private.%s", ibmRegion, ibmBaseURL))
		} else {
			c.BaseURL, err = url.Parse(fmt.Sprintf("https://%s.%s", ibmRegion, ibmBaseURL))
		}
		return err
	}
}

// WithUserAgent sets the User Agent to be sent to Sysdig.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithAuthenticator sets the authentication.Authenticator for the Sysdig Client.
func WithAuthenticator(a authentication.Authenticator) ClientOption {
	return func(c *Client) error {
		c.authenticator = a
		return nil
	}
}

// WithResponseCompression sets whether to set compression headers in requests to Sysdig.
func WithResponseCompression(shouldCompressResponse bool) ClientOption {
	return func(c *Client) error {
		c.shouldCompressResponse = shouldCompressResponse
		return nil
	}
}

// WithLogger sets the default logger for the Client.
func WithLogger(l Logger) ClientOption {
	return func(c *Client) error {
		c.logger = l
		return nil
	}
}

// WithDebug sets whether to print debug information about requests and responses.
func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

// Client returns the http.Client used by this Sysdig client.
func (c *Client) Client() *http.Client {
	clientCopy := *c.httpClient
	return &clientCopy
}

// SetLogger sets the logger to be used by this Sysdig client.
func (c *Client) SetLogger(l Logger) {
	c.logger = l
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		eerr := enc.Encode(body)
		if eerr != nil {
			return nil, eerr
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.shouldCompressResponse {
		req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	}
	return req, nil
}

// BareDo sends an API request and lets you handle the api response. If an error
// or API Error occurs, the error will contain more information. Otherwise, you
// are supposed to read and close the response's Body.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it is
// canceled or times out, ctx.Err() will be returned.
func (c *Client) BareDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.bareDo(ctx, req)
}

func (c *Client) bareDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	if ctx == nil {
		return nil, fmt.Errorf("cannot pass a nil-context")
	}
	if c.authenticator != nil {
		if c.debug {
			c.logger.Printf("authenticating with %T", c.authenticator)
		}
		if err := c.authenticator.Authenticate(req); err != nil {
			return nil, err
		}
		if c.debug {
			c.logger.Print("authentication succeeded")
		}
	}
	if c.debug {
		if req != nil {
			var data []byte
			if req.Body != nil {
				var rerr error
				data, rerr = ioutil.ReadAll(req.Body)
				if rerr != nil {
					c.logger.Printf("failed to read request body for debugging: %v", rerr)
				} else {
					req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
				}
			}
			if req.URL != nil {
				c.logger.Printf("-> request: %s %s\n%s", req.Method, req.URL.String(), string(data))
			}
		}
	}
	cReq := req.Clone(ctx)
	resp, err := c.httpClient.Do(cReq)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	if c.authenticator != nil && isAuthenticationError(resp) {
		if refreshableAuthenticator, ok := c.authenticator.(authentication.Refreshable); ok {
			if rerr := refreshableAuthenticator.Refresh(); rerr != nil {
				c.logger.Printf("error refreshing authenticator: %v", rerr)
				return nil, rerr
			}
			// Retry one time after a successful refresh
			return c.bareDo(ctx, req)
		}
	}
	if gzipped(resp) {
		var gerr error
		resp.Body, gerr = gzip.NewReader(resp.Body)
		if gerr != nil {
			c.logger.Printf("failed to inflate gzipped response: %v", gerr)
			return nil, gerr
		}
	}
	if c.debug {
		data, rerr := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.logger.Printf("failed to read response body for debugging: %v", rerr)
		} else {
			c.logger.Printf("<- response: %d\n%s", resp.StatusCode, string(data))
			for k, v := range resp.Header {
				c.logger.Printf("%s: %s", k, strings.Join(v, ","))
			}
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
	}
	err = c.CheckResponse(resp)
	return resp, err
}

func gzipped(resp *http.Response) bool {
	return strings.Contains(resp.Header.Get("Content-Encoding"), "gzip")
}

func isAuthenticationError(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it. If v is nil, and no error happens, the response is returned as is.
//
// The provided ctx must be non-nil, if it is nil an error is returned. If it
// is canceled or times out, ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.BareDo(ctx, req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	switch iv := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(iv, resp.Body)
	default:
		decErr := json.NewDecoder(resp.Body).Decode(v)
		if decErr == io.EOF {
			decErr = nil // ignore EOF errors caused by empty response body
		}
		if decErr != nil {
			err = decErr
		}
	}
	return resp, err
}

// PrometheusClient creates a Prometheus Client using the Sysdig Client as a base.
// Note: only a subset of the PrometheusClient APIs are implemented by Sysdig.
// Known working:
// - Query
// - QueryRange
// - Alerts.
func (c *Client) PrometheusClient() v1.API {
	return v1.NewAPI(&prometheusClient{client: c})
}

type prometheusClient struct {
	client *Client
}

// URL implements URL for the Prometheus API Client interface.
// See: https://github.com/prometheus/client_golang/blob/v1.9.0/api/client.go
func (c *prometheusClient) URL(endpoint string, args map[string]string) *url.URL {
	p := path.Join("prometheus", endpoint)
	for arg, val := range args {
		arg = ":" + arg
		p = strings.ReplaceAll(p, arg, val)
	}
	u, err := c.client.BaseURL.Parse(p)
	if err != nil {
		c.client.logger.Printf("invalid prometheus endpoint %q: %v", endpoint, err)
		return nil
	}
	return u
}

// Do implements Do for the Prometheus Client API client.
// See: https://github.com/prometheus/client_golang/blob/v1.9.0/api/client.go
func (c *prometheusClient) Do(ctx context.Context, request *http.Request) (*http.Response, []byte, error) {
	resp, err := c.client.BareDo(ctx, request)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	return resp, b, err
}

// ErrorResponse reports one or more errors caused by an API request.
type ErrorResponse struct {
	Response *http.Response
	Message  string  `json:"message,omitempty"`
	Errors   []Error `json:"errors,omitempty"`
}

// Error contains a further explanation for the reason of an error..
type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

// Error implements the error interface for ErrorResponse.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Errors)
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
// API error responses are expected to have response
// body, and a JSON response body that maps to ErrorResponse.
func (c *Client) CheckResponse(r *http.Response) error {
	if c := r.StatusCode; http.StatusOK <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		if uerr := json.Unmarshal(data, errorResponse); uerr != nil {
			errorResponse.Message = fmt.Sprintf("error unmarshaling error response: %v", uerr)
		}
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return errorResponse
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}
	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}
	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}
	u.RawQuery = qs.Encode()
	return u.String(), nil
}
