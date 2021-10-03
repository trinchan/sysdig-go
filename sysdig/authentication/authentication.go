package authentication

import (
	"net/http"
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
