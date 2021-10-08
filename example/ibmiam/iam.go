package main

import (
	"context"
	"log"
	"os"

	"github.com/trinchan/sysdig-go/sysdig"
	"github.com/trinchan/sysdig-go/sysdig/authentication/ibmiam"
)

func main() {
	apiKey := os.Getenv("IBM_API_KEY")
	instanceID := os.Getenv("SYSDIG_INSTANCE_ID")
	teamID := os.Getenv("SYSDIG_TEAM_ID")
	ctx := context.Background()
	authenticator, err := ibmiam.Authenticator(
		apiKey,
		ibmiam.WithIBMInstanceID(instanceID),
		ibmiam.WithSysdigTeamID(teamID),
	)
	if err != nil {
		panic(err)
	}
	client, err := sysdig.NewClient(
		sysdig.WithAuthenticator(authenticator),
		sysdig.WithIBMBaseURL(sysdig.RegionUSSouth, false),
		// sysdig.WithDebug(false), // Enable to see requests/responses.
		// sysdig.WithLogger(log.Default()), // Uncomment to set the logger when setting debug to true.
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
