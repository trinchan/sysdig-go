package main

import (
	"context"
	"log"
	"os"

	"github.com/trinchan/sysdig-go/sysdig"
	"github.com/trinchan/sysdig-go/sysdig/authentication/accesstoken"
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
		authenticator,
		// sysdig.WithDebug(true),                             // Enable to see requests/responses.
		// sysdig.WithLogger(log.Default()),                   // Uncomment to set the logger when setting debug to true.
		// sysdig.WithIBMBaseURL(sysdig.RegionUSSouth, false), // Uncomment to use an IBM Cloud Monitoring instance
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
