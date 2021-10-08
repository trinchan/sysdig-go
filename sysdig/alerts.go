package sysdig

import (
	"context"
	"fmt"
	"net/http"
)

// AlertService is the Service for communicating with the Sysdig Monitor Alert related API.
type AlertService service

type AlertType string

const (
	// AlertTypeEvent means the Alert is from an Event.
	AlertTypeEvent AlertType = "EVENT"
)

type Alert struct {
	ID                 int                     `json:"id,omitempty"`
	Version            int                     `json:"version,omitempty"`
	CreatedOn          MilliTime               `json:"createdOn,omitempty"`
	ModifiedOn         MilliTime               `json:"modifiedOn,omitempty"`
	Type               AlertType               `json:"type"`
	Name               string                  `json:"name"`
	Description        string                  `json:"description"`
	Enabled            bool                    `json:"enabled"`
	Criteria           AlertCriteria           `json:"criteria"`
	Severity           Severity                `json:"severity"`
	Timespan           MicroDuration           `json:"timespan"`
	CustomNotification AlertCustomNotification `json:"customNotification"`
	NotificationCount  int                     `json:"notificationCount"`
	TeamID             int                     `json:"teamId"`
	AutoCreated        bool                    `json:"autoCreated"`
	RateOfChange       bool                    `json:"rateOfChange"`
	ReNotifyMinutes    int                     `json:"reNotifyMinutes"`
	ReNotify           bool                    `json:"reNotify"`
	// TODO what is the format of this?
	InvalidMetrics []interface{} `json:"invalidMetrics"`
	GroupName      string        `json:"groupName"`
	Valid          bool          `json:"valid"`
	SeverityLabel  SeverityLabel `json:"severityLabel"`
	Condition      string        `json:"condition"`
	CustomerID     int           `json:"customerId"`
}

type AlertCustomNotification struct {
	TitleTemplate  string `json:"titleTemplate"`
	UseNewTemplate bool   `json:"useNewTemplate"`
}

type AlertCriteria struct {
	Text string `json:"text"`
	// TODO what are the format of these?
	Source   interface{} `json:"source"`
	Severity interface{} `json:"severity"`
	Query    interface{} `json:"query"`
	Scope    interface{} `json:"scope"`
}

type AlertResponse struct {
	Alert Alert `json:"alert"`
}

// Get retrieves a Alert.
func (s *AlertService) Get(ctx context.Context, alertID int) (*AlertResponse, *http.Response, error) {
	u := fmt.Sprintf("api/alerts/%d", alertID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(AlertResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// ListAlertConfigurationsResponse describes the response for AlertsService.List.
type ListAlertConfigurationsResponse struct {
	Alerts []Alert `json:"alerts"`
}

// List lists all Alerts.
func (s *AlertService) List(ctx context.Context) (*ListAlertConfigurationsResponse, *http.Response, error) {
	u := "api/alerts"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ListAlertConfigurationsResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Delete deletes an Alert.
func (s *AlertService) Delete(ctx context.Context, alertID int) (*http.Response, error) {
	u := fmt.Sprintf("api/alerts/%d", alertID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
