package authentication

import (
	"fmt"
	"net/http"
)

const (
	// IBMInstanceIDHeader is the header used to indicate an IBM Cloud Monitoring instance to target.
	IBMInstanceIDHeader = "IBMInstanceID"
	// SysdigTeamIDHeader is the header used to indicate a Sysdig Team to target. May not be required anymore?
	SysdigTeamIDHeader = "TeamID"
	// AuthorizationHeader is the standard Authorization header used to authorize to the Sysdig API.
	AuthorizationHeader = "Authorization"
)

// Authenticator defines an interface for authenticating a request to the Sysdig API.
// Implementers should add required headers or authorization fields to the request.
type Authenticator interface {
	Authenticate(req *http.Request) error
}

// Refreshable defines an optional interface for Authenticators that can be Refreshed.
// Authentication failures will trigger a Refresh and a retry when implemented.
type Refreshable interface {
	Refresh() error
}

// AuthenticatorFunc defines a function that will authenticate the given Request.
type AuthenticatorFunc func(req *http.Request) error

// Authenticate implements Authenticator using the AuthenticatorFunc.
func (f AuthenticatorFunc) Authenticate(req *http.Request) error {
	return f(req)
}

// AuthorizationHeaderFor returns a formatted Authorization header as a Bearer token.
func AuthorizationHeaderFor(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}
