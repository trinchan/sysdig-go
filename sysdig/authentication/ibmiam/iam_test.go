package ibmiam

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/trinchan/sysdig-go/sysdig/authentication"
)

func TestAuthenticator(t *testing.T) {
	wantIBMIDHeader := "ibmiam"
	wantSysdigTeamID := "ibm"
	accessToken := "foo"
	wantAccessToken := iamTokenResponse{
		AccessToken: "bar",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(wantAccessToken)
		if err != nil {
			t.Fatal(err)
		}
		if _, werr := w.Write(b); werr != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	a, err := Authenticator(accessToken,
		WithIBMInstanceID(wantIBMIDHeader),
		WithSysdigTeamID(wantSysdigTeamID),
		WithIAMEndpoint(server.URL),
		WithRefreshBeforeDuration(DefaultRefreshBeforeExpirationDuration),
		WithHTTPClient(server.Client()),
	)
	if err != nil {
		t.Fatal(err)
	}
	if aerr := a.Authenticate(req); aerr != nil {
		t.Fatal(aerr)
	}
	gotIBMIDHeader := req.Header.Get(authentication.IBMInstanceIDHeader)
	if gotIBMIDHeader != wantIBMIDHeader {
		t.Errorf("got IBMInstanceID header: %s, want: %s", gotIBMIDHeader, wantIBMIDHeader)
	}
	gotSysdigTeamID := req.Header.Get(authentication.SysdigTeamIDHeader)
	if gotSysdigTeamID != wantSysdigTeamID {
		t.Errorf("got SysdigTeamID header: %s, want: %s", gotSysdigTeamID, wantSysdigTeamID)
	}
	gotAccessToken := strings.TrimPrefix(req.Header.Get(authentication.AuthorizationHeader), "Bearer ")
	if gotAccessToken != wantAccessToken.AccessToken {
		t.Errorf("got access token header: %s, want: %s", gotAccessToken, wantAccessToken)
	}
}

func TestAuthenticatorEmpty(t *testing.T) {
	_, err := Authenticator("")
	if err == nil {
		t.Fatal("did not return an expected error")
	}
}

func TestAuthenticatorBadOption(t *testing.T) {
	_, err := Authenticator("foo", WithRefreshBeforeDuration(tokenValidDuration+time.Minute))
	if err == nil {
		t.Fatal("did not return an expected error")
	}
}
