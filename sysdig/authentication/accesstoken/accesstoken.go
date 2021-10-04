package accesstoken

import (
	"fmt"
	"net/http"

	"github.com/trinchan/sysdig-go/sysdig/authentication"
)

const (
	ibmInstanceIDHeader = "IBMInstanceID"
	sysdigTeamIDHeader  = "TeamID"
	authorizationHeader = "Authorization"
)

type authenticator struct {
	token string

	ibmInstanceID string
	sysdigTeamID  string
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

// Authenticate implements the authentication.Authenticator interface using a Sysdig Access Key.
func (a *authenticator) Authenticate(req *http.Request) error {
	req.Header.Set(authorizationHeader, fmt.Sprintf("Bearer %s", a.token))
	if a.ibmInstanceID != "" {
		req.Header.Set(ibmInstanceIDHeader, a.ibmInstanceID)
	}
	if a.sysdigTeamID != "" {
		req.Header.Set(sysdigTeamIDHeader, a.sysdigTeamID)
	}
	return nil
}

// AuthenticatorOption defines the type for passing options to the Authenticator constructor.
type AuthenticatorOption func(*authenticator) error

// Authenticator implements the authentication.Authenticator interface using a Sysdig Access Key.
// See: https://cloud.ibm.com/docs/monitoring?topic=monitoring-api_monitoring_token for IBM Cloud
// See: https://sysdig.gitbooks.io/sysdig-cloud-api/content/rest_api/getting_started.html for Sysdig Cloud.
func Authenticator(accessToken string, options ...AuthenticatorOption) (authentication.Authenticator, error) {
	a := &authenticator{
		token: accessToken,
	}
	for _, o := range options {
		if err := o(a); err != nil {
			return nil, err
		}
	}
	if accessToken == "" {
		return nil, fmt.Errorf("access token must be set")
	}
	return a, nil
}
