package ibmiam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/trinchan/sysdig-go/sysdig/authentication"
)

const (
	// DefaultIAMEndpoint is the default production IAM endpoint for IBM Cloud.
	DefaultIAMEndpoint = "https://iam.cloud.ibm.com/identity/token"
	// TestIAMEndpoint is the test IAM endpoint for IBM Cloud.
	TestIAMEndpoint = "https://iam.test.cloud.ibm.com/identity/token"

	// DefaultRefreshBeforeExpirationDuration is the default duration before expected expiration to refresh the IAM token.
	DefaultRefreshBeforeExpirationDuration = 5 * time.Minute
	// tokenValidDuration is the default validity period for IBM Cloud IAM tokens.
	tokenValidDuration   = time.Hour
	defaultRefreshBefore = tokenValidDuration - DefaultRefreshBeforeExpirationDuration
)

type iamTokenResponse struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	UAAAccessToken  string `json:"uaa_token"`
	UAARefreshToken string `json:"uaa_refresh_token"`
	TokenType       string `json:"token_type"`
}

type authenticator struct {
	httpClient    *http.Client
	iamEndpoint   string
	apiKey        string
	ibmInstanceID string
	sysdigTeamID  string
	refreshBefore time.Duration

	lock        sync.RWMutex
	lastRefresh time.Time
	token       iamTokenResponse
}

// Authenticate implements authentication.Authenticator using IBM Cloud IAM.
func (a *authenticator) Authenticate(req *http.Request) error {
	if time.Since(a.lastRefresh) > a.refreshBefore {
		if err := a.Refresh(); err != nil {
			return err
		}
	}
	a.lock.RLock()
	at := a.token.AccessToken
	a.lock.RUnlock()

	req.Header.Set(authentication.AuthorizationHeader, authentication.AuthorizationHeaderFor(at))
	req.Header.Set(authentication.IBMInstanceIDHeader, a.ibmInstanceID)
	if a.sysdigTeamID != "" {
		req.Header.Set(authentication.SysdigTeamIDHeader, a.sysdigTeamID)
	}
	return nil
}

// Refresh implements Refreshable for the Authenticator.
func (a *authenticator) Refresh() error {
	return a.refreshAccessToken()
}

func (a *authenticator) refreshAccessToken() error {
	a.lock.Lock()
	defer a.lock.Unlock()
	v := url.Values{
		"grant_type":    []string{"urn:ibm:params:oauth:grant-type:apikey"},
		"response_type": []string{"cloud_iam"},
		"apikey":        []string{a.apiKey},
	}
	req, err := http.NewRequest(http.MethodPost, a.iamEndpoint, bytes.NewBufferString(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to refresh token: %d: %v", resp.StatusCode, err)
	}
	err = json.NewDecoder(resp.Body).Decode(&a.token)
	if err != nil {
		return err
	}
	a.lastRefresh = time.Now()
	return nil
}

// AuthenticatorOption defines options for the IBM IAM authentication.Authenticator.
type AuthenticatorOption func(*authenticator) error

// WithHTTPClient sets the http.Client to be used for IAM authentication.
func WithHTTPClient(c *http.Client) AuthenticatorOption {
	return func(a *authenticator) error {
		a.httpClient = c
		return nil
	}
}

// WithRefreshBeforeDuration sets the duration before expiration to trigger a token refresh.
func WithRefreshBeforeDuration(duration time.Duration) AuthenticatorOption {
	return func(a *authenticator) error {
		if duration > tokenValidDuration {
			return fmt.Errorf(
				"invalid refresh before duration: %s, must be less than expiration time: %s",
				duration,
				tokenValidDuration,
			)
		}
		a.refreshBefore = tokenValidDuration - duration
		return nil
	}
}

// WithIAMEndpoint sets the IAM endpoint to be used for IAM authentication.
func WithIAMEndpoint(iamEndpoint string) AuthenticatorOption {
	return func(a *authenticator) error {
		a.iamEndpoint = iamEndpoint
		return nil
	}
}

// WithIBMInstanceID sets the instance ID to be set for IBM Sysdig requests.
// See: https://cloud.ibm.com/docs/monitoring?topic=monitoring-mon-curl#mon-curl-headers-iam
func WithIBMInstanceID(ibmInstanceID string) AuthenticatorOption {
	return func(a *authenticator) error {
		a.ibmInstanceID = ibmInstanceID
		return nil
	}
}

// WithSysdigTeamID sets the TeamID to be set for IBM Sysdig requests.
// May not be required. TODO: check if this is still required.
// See: https://cloud.ibm.com/docs/monitoring?topic=monitoring-team_id
func WithSysdigTeamID(sysdigTeamID string) AuthenticatorOption {
	return func(a *authenticator) error {
		a.sysdigTeamID = sysdigTeamID
		return nil
	}
}

// Authenticator returns an authentication.Authenticator for IBM Cloud IAM.
func Authenticator(apiKey string, options ...AuthenticatorOption) (authentication.Authenticator, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apikey cannot be blank")
	}
	a := &authenticator{
		httpClient:    http.DefaultClient,
		iamEndpoint:   DefaultIAMEndpoint,
		refreshBefore: defaultRefreshBefore,
		apiKey:        apiKey,
	}
	for _, o := range options {
		if err := o(a); err != nil {
			return nil, err
		}
	}
	return a, nil
}
