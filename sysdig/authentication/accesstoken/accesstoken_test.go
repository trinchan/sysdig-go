package accesstoken

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestAuthenticator(t *testing.T) {
	wantIBMIDHeader := "ibmiam"
	wantSysdigTeamID := "ibm"
	wantAccessToken := "foo"
	a, err := Authenticator(wantAccessToken, WithIBMInstanceID(wantIBMIDHeader), WithSysdigTeamID(wantSysdigTeamID))
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	if aerr := a.Authenticate(req); aerr != nil {
		t.Fatal(aerr)
	}
	gotIBMIDHeader := req.Header.Get(ibmInstanceIDHeader)
	if gotIBMIDHeader != wantIBMIDHeader {
		t.Errorf("got IBMInstanceID header: %s, want: %s", gotIBMIDHeader, wantIBMIDHeader)
	}
	gotSysdigTeamID := req.Header.Get(sysdigTeamIDHeader)
	if gotSysdigTeamID != wantSysdigTeamID {
		t.Errorf("got SysdigTeamID header: %s, want: %s", gotSysdigTeamID, wantSysdigTeamID)
	}
	gotAccessToken := strings.TrimPrefix(req.Header.Get(authorizationHeader), "Bearer ")
	if gotAccessToken != wantAccessToken {
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
	_, err := Authenticator("foo", func(a *authenticator) error { return errors.New("test error") })
	if err == nil {
		t.Fatal("did not return an expected error")
	}
}
