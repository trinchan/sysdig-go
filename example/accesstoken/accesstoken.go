package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/trinchan/sysdig-go/sysdig"
	"github.com/trinchan/sysdig-go/sysdig/authentication/accesstoken"
	"github.com/trinchan/sysdig-go/sysdig/scope"
)

func main() {
	accessToken := os.Getenv("SYSDIG_ACCESS_TOKEN")
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
	log.Printf("instance id: %q", instanceID)
	log.Printf("team id: %q", teamID)
	accessTokenAuthenticator, err := accesstoken.Authenticator(
		accessToken,
		accesstoken.WithIBMInstanceID(instanceID),
		accesstoken.WithSysdigTeamID(teamID),
	)
	if err != nil {
		panic(err)
	}
	client, err := sysdig.NewClient(
		sysdig.WithLogger(log.Default()),
		sysdig.WithHTTPClient(httpClient),
		sysdig.WithAuthenticator(accessTokenAuthenticator),
		sysdig.WithIBMBaseURL(sysdig.RegionUSSouth, false),
		sysdig.WithResponseCompression(true),
		sysdig.WithUserAgent("sysdig-go; access token example"),
	)
	if err != nil {
		panic(err)
	}
	s := scope.NewEventScope().
		AddIsSelection("app", "example-app").
		AddIsSelection("host", "my-laptop")
	event, _, err := client.Events.CreateEvent(ctx, &sysdig.EventOptions{
		Name:        "example event",
		Description: "This is an example event",
		Severity:    sysdig.SeverityMedium,
		Scope:       s.String(),
		Tags:        map[string]string{"app": "example"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("created event: %+v\n", event)
	fmt.Println("sleeping a few seconds to ensure events sync")
	time.Sleep(5 * time.Second)
	events, _, err := client.Events.ListEvents(ctx, &sysdig.ListEventOptions{
		Scope:        s.String(),
		IncludeTotal: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("retrieved %d matches of %d total\n", events.Matched, events.Total)
	for i, e := range events.Events {
		fmt.Printf("[%d/%d] deleting event %s - %s\n", i+1, len(events.Events), e.ID, e.Name)
		if _, derr := client.Events.DeleteEvent(ctx, e.ID); derr != nil {
			panic(derr)
		}
	}

	fmt.Printf("deleted %d events\n", events.Matched)
}
