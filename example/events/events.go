package main

import (
	"context"
	"log"
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
	authenticator, err := accesstoken.Authenticator(
		accessToken,
		accesstoken.WithIBMInstanceID(instanceID),
		accesstoken.WithSysdigTeamID(teamID),
	)
	if err != nil {
		panic(err)
	}
	client, err := sysdig.NewClient(
		sysdig.WithAuthenticator(authenticator),
		// sysdig.WithDebug(true), // Enable to see requests/responses.
		// sysdig.WithLogger(log.Default()), // Uncomment to set the logger when setting debug to true.
		// sysdig.WithIBMBaseURL(sysdig.RegionUSSouth, false), // Uncomment to use an IBM Cloud Monitoring instance
	)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	s := scope.NewEventScope().
		AddIsSelection("app", "example-app").
		AddIsSelection("host", "my-laptop")
	event, _, err := client.Events.Create(ctx, sysdig.EventOptions{
		Name:        "example event",
		Description: "This is an example event",
		Severity:    sysdig.SeverityMedium,
		Scope:       s.String(),
		Tags:        map[string]string{"app": "example"},
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Created event: %+v", event)
	log.Print("Sleeping a few seconds to ensure events sync...")
	time.Sleep(5 * time.Second)
	events, _, err := client.Events.List(ctx, sysdig.ListEventOptions{
		Scope:        s.String(),
		IncludeTotal: true,
	})
	if err != nil {
		panic(err)
	}
	log.Printf("retrieved %d matches of %d total", events.Matched, events.Total)
	for i, e := range events.Events {
		log.Printf("[%d/%d] deleting event %s - %s", i+1, len(events.Events), e.ID, e.Name)
		if _, derr := client.Events.Delete(ctx, e.ID); derr != nil {
			panic(derr)
		}
	}

	log.Printf("Deleted %d events", events.Matched)
}
