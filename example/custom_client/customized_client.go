package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/trinchan/sysdig-go/sysdig"
	"github.com/trinchan/sysdig-go/sysdig/authentication/ibmiam"
)

func main() {
	apiKey := os.Getenv("IBM_API_KEY")
	instanceID := os.Getenv("SYSDIG_INSTANCE_ID")
	teamID := os.Getenv("SYSDIG_TEAM_ID")
	ctx := context.Background()
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 2500 * time.Millisecond,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				MinVersion:         tls.VersionTLS12,
			},
			TLSHandshakeTimeout: 2500 * time.Millisecond,
			MaxIdleConns:        200,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
		Timeout: 10 * time.Second,
	}
	authenticator, err := ibmiam.Authenticator(
		apiKey,
		ibmiam.WithIBMInstanceID(instanceID),
		ibmiam.WithSysdigTeamID(teamID),
		ibmiam.WithHTTPClient(httpClient),
	)
	if err != nil {
		panic(err)
	}
	client, err := sysdig.NewClient(
		sysdig.WithAuthenticator(authenticator),
		sysdig.WithHTTPClient(httpClient),
		sysdig.WithResponseCompression(true),
		sysdig.WithUserAgent("sysdig-go; custom client example"),
		// sysdig.WithDebug(false), // Enable to see requests/responses.
		// sysdig.WithLogger(log.Default()), // Uncomment to set the logger when setting debug to true.
		sysdig.WithIBMBaseURL(sysdig.RegionUSSouth, false), // Uncomment to use an IBM Cloud Monitoring instance
	)
	if err != nil {
		panic(err)
	}
	me, _, err := client.Users.Me(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Logged in as %s %s", me.User.FirstName, me.User.LastName)
}
