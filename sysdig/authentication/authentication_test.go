package authentication

import (
	"net/http"
	"testing"
)

func TestAuthorizationHeaderFor(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "sysdig-token",
			in:   "123e4567-e89b-12d3-a456-426614174000",
			want: "Bearer 123e4567-e89b-12d3-a456-426614174000",
		},
		{
			name: "jwt-token",
			in: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY" +
				"3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
				"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			want: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM" +
				"0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
				"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := AuthorizationHeaderFor(test.in)
			if got != test.want {
				t.Errorf("got: %q, want: %q", got, test.want)
			}
		})
	}
}

func TestAuthenticateFunc(t *testing.T) {
	ran := false
	authFunc := AuthenticatorFunc(func(req *http.Request) error { ran = true; return nil })
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	err = authFunc.Authenticate(req)
	if err != nil {
		t.Errorf("AuthenticatorFunc returned unexpected error: %v", err)
	}
	if !ran {
		t.Error("AuthenticatorFunc did not run when authenticating")
	}
}
