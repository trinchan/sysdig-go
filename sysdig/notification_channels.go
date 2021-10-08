package sysdig

import (
	"context"
	"fmt"
	"net/http"
)

// NotificationChannelsService is the Service for communicating with the Sysdig Monitor Notification Channel related API.
type NotificationChannelsService service

// NotificationChannelType is a type of Sysdig notification channel.
type NotificationChannelType string

const (
	// NotificationChannelTypeEmail is a notification channel which sends to email.
	NotificationChannelTypeEmail NotificationChannelType = "EMAIL"
	// NotificationChannelTypeSNS is a notification channel which sends to Amazon SNS.
	NotificationChannelTypeSNS NotificationChannelType = "SNS"
	// NotificationChannelTypePagerDuty is a notification channel which sends to Pagerduty.
	NotificationChannelTypePagerDuty NotificationChannelType = "PAGER_DUTY"
	// NotificationChannelTypeSlack is a notification channel which sends to Slack.
	NotificationChannelTypeSlack NotificationChannelType = "SLACK"
	// NotificationChannelTypeOpsGenie is a notification channel which sends to OpsGenie.
	NotificationChannelTypeOpsGenie NotificationChannelType = "OPSGENIE"
	// NotificationChannelTypeVictorOps is a notification channel which sends to VictorOps.
	NotificationChannelTypeVictorOps NotificationChannelType = "VICTOROPS"
	// NotificationChannelTypeWebhook is a notification channel which sends to a webhook.
	NotificationChannelTypeWebhook NotificationChannelType = "WEBHOOK"
)

// NotificationChannel describes a Sysdig notification channel. Used to direct alerts or notifications.
type NotificationChannel struct {
	Type    NotificationChannelType    `json:"type"`
	Name    string                     `json:"name"`
	Enabled bool                       `json:"enabled"`
	Options NotificationChannelOptions `json:"options"`

	ID         string     `json:"id,omitempty"`
	Version    int        `json:"version,omitempty"`
	CreatedOn  *MilliTime `json:"createdOn,omitempty"`
	ModifiedOn *MilliTime `json:"modifiedOn,omitempty"`
}

// NotificationChannelOptions describes the options for a NotificationChannel.
// See: https://cloud.ibm.com/docs/monitoring?topic=monitoring-notifications_api#notifications-api-parm
type NotificationChannelOptions struct {
	NotifyOnOK      bool     `json:"notifyOnOk"`
	NotifyOnResolve bool     `json:"notifyOnResolve"`
	ResolveOnOK     bool     `json:"resolveOnOk"`
	Channel         string   `json:"channel"`
	EmailRecipients []string `json:"emailRecipients"`
	URL             string   `json:"url"`
	APIKey          string   `json:"apiKey"`
	RoutingKey      string   `json:"routingKey"`
	Account         string   `json:"account"`
	ServiceKey      string   `json:"serviceKey"`
	ServiceName     string   `json:"serviceName"`
}

// NotificationChannelResponse describes the response for a NotificationChannel from the NotificationChannelsService API.
type NotificationChannelResponse struct {
	NotificationChannel NotificationChannel `json:"notificationChannel"`
}

// Get retrieves a NotificationChannel.
func (s *NotificationChannelsService) Get(ctx context.Context, id string) (*NotificationChannelResponse, *http.Response, error) {
	u := fmt.Sprintf("api/notificationChannels/%s", id)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(NotificationChannelResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// ListNotificationChannelsResponse describes the response for NotificationChannelsService.List.
type ListNotificationChannelsResponse struct {
	NotificationChannels []NotificationChannel `json:"notificationChannels"`
}

// List lists all NotificationChannels.
func (s *NotificationChannelsService) List(
	ctx context.Context,
	from, to MilliTime) (*ListNotificationChannelsResponse, *http.Response, error) {
	u := "api/notificationChannels"
	type listOptions struct {
		From MilliTime `url:"from"`
		To   MilliTime `url:"to"`
	}
	o := listOptions{
		From: from,
		To:   to,
	}
	uWithOpts, err := addOptions(u, o)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, uWithOpts, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ListNotificationChannelsResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Create creates a new NotificationChannel.
func (s *NotificationChannelsService) Create(
	ctx context.Context,
	t NotificationChannelType,
	name string,
	options NotificationChannelOptions) (*NotificationChannelResponse, *http.Response, error) {
	u := "api/notificationChannels"
	channel := NotificationChannel{
		Type:    t,
		Name:    name,
		Enabled: true,
		Options: options,
	}
	type notificationChannelRequest struct {
		NotificationChannel NotificationChannel `json:"notificationChannel"`
	}
	r := notificationChannelRequest{
		NotificationChannel: channel,
	}
	req, err := s.client.NewRequest(http.MethodPost, u, r)
	if err != nil {
		return nil, nil, err
	}
	c := new(NotificationChannelResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Delete deletes a NotificationChannel.
func (s *NotificationChannelsService) Delete(ctx context.Context, notificationChannelID string) (*http.Response, error) {
	u := fmt.Sprintf("api/notificationChannels/%s", notificationChannelID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
