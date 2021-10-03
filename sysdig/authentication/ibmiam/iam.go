package ibmiam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/trinchan/sysdig-go/sysdig/authentication"
)

const (
	// DefaultIAMEndpoint is the default production IAM endpoint for IBM Cloud.
	DefaultIAMEndpoint = "https://iam.cloud.ibm.com/identity/token"
	// TestIAMEndpoint is the test IAM endpoint for IBM Cloud.
	TestIAMEndpoint = "https://iam.test.cloud.ibm.com/identity/token"
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

	lock  sync.RWMutex
	token iamTokenResponse
}

// Authenticate implements authentication.Authenticator using IBM Cloud IAM.
func (a *authenticator) Authenticate(req *http.Request) error {
	a.lock.RLock()
	at := a.token.AccessToken
	a.lock.RUnlock()

	if at == "" {
		if err := a.Refresh(); err != nil {
			return err
		}
		a.lock.RLock()
		at = a.token.AccessToken
		a.lock.RUnlock()
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at))
	req.Header.Set("IBMInstanceID", a.ibmInstanceID)
	if a.sysdigTeamID != "" {
		req.Header.Set("TeamID", a.sysdigTeamID)
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
		return fmt.Errorf("failed to refresh token: %d: %w", resp.StatusCode, err)
	}
	err = json.NewDecoder(resp.Body).Decode(&a.token)
	if err != nil {
		return err
	}
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

// WithIAMEndpoint sets the IAM endpoint to be used for IAM authentication.
func WithIAMEndpoint(iamEndpoint string) AuthenticatorOption {
	return func(a *authenticator) error {
		a.iamEndpoint = iamEndpoint
		return nil
	}
}

// WithAPIKey sets the IBM Cloud API key to be used for IAM authentication.
func WithAPIKey(apiKey string) AuthenticatorOption {
	return func(a *authenticator) error {
		a.apiKey = apiKey
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
// See: https://cloud.ibm.com/docs/monitoring?topic=monitoring-team_id
func WithSysdigTeamID(sysdigTeamID string) AuthenticatorOption {
	return func(a *authenticator) error {
		a.sysdigTeamID = sysdigTeamID
		return nil
	}
}

// Authenticator returns an authentication.Authenticator for IBM Cloud IAM.
func Authenticator(options ...AuthenticatorOption) (authentication.Authenticator, error) {
	a := &authenticator{
		httpClient:  http.DefaultClient,
		iamEndpoint: DefaultIAMEndpoint,
	}
	for _, o := range options {
		if err := o(a); err != nil {
			return nil, err
		}
	}
	return a, nil
}
