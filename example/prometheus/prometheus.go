package main

import (
	"context"
	"log"
	"os"
	"time"

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
	values, _, err := client.Prometheus.Query(ctx, "sysdig_host_cpu_used_percent", time.Now())
	if err != nil {
		panic(err)
	}
	log.Printf("Values: %+v", values)
	alerts, err := client.Prometheus.Alerts(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Found %d alerts", len(alerts.Alerts))
	for _, alert := range alerts.Alerts {
		log.Printf("Alert: %+v", alert)
	}
}
