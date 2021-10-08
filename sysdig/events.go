package sysdig

import (
	"context"
	"fmt"
	"net/http"
)

// EventsService is the Service for communicating with the Sysdig Events API.
type EventsService service

// Severity is the severity level for the Event. The severity is in syslog
// style and ranges from 0 (high) to 7 (low).
type Severity int

const (
	// SeverityEmergency is a "panic" condition - notify all tech staff
	// on call? (Earthquake? Tornado?) - affects multiple apps/servers/sites.
	SeverityEmergency Severity = iota
	// SeverityAlert should be corrected immediately - notify staff who can fix
	// the problem - example is loss of backup ISP connection.
	SeverityAlert
	// SeverityCritical should be corrected immediately, but indicates failure
	// in a primary system - fix CRITICAL problems before ALERT - example is loss
	// of primary ISP connection.
	SeverityCritical
	// SeverityError is a non-urgent failure - these should be relayed to
	// developers or admins; each item must be resolved within a given time.
	SeverityError
	// SeverityWarning are warning messages - not an error, but indication that
	// an error will occur if action is not taken, e.g. file system 85% full -
	// each item must be resolved within a given time.
	SeverityWarning
	// SeverityNotice are events that are unusual but not error conditions -
	// might be summarized in an email to developers or admins to spot potential
	// problems - no immediate action required.
	SeverityNotice
	// SeverityInformational are normal operational messages - may be harvested
	// for reporting, measuring throughput, etc. - no action required.
	SeverityInformational
	// SeverityDebug is info useful to developers for debugging the app, not
	// useful during operations.
	SeverityDebug
)

// SeverityLabel is the severity level label for an Event.
type SeverityLabel string

const (
	// SeverityHigh is a high severity alert. It should be corrected immediately.
	SeverityHigh SeverityLabel = "HIGH"
	// SeverityMedium is a medium severity alert. It should be corrected immediately, but indicates failure
	// in a primary system - fix SeverityHigh problems before SeverityMedium.
	SeverityMedium SeverityLabel = "MEDIUM"
	// SeverityLow is a low severity alert. It is a non-urgent failure - these should be relayed to
	// developers or admins; each item must be resolved within a given time.
	SeverityLow SeverityLabel = "LOW"
	// SeverityInfo are normal operational messages - may be harvested
	// for reporting, measuring throughput, etc. - no action required.
	SeverityInfo SeverityLabel = "INFO"
	// SeverityNone is an event without any specified SeverityLabel.
	SeverityNone SeverityLabel = "NONE"
)

// Direction defines the ordering of a list of events. (?) TODO figure out what this parameter actually does
type Direction string

const (
	// DirectionBefore will order a list of events by oldest age first (?).
	DirectionBefore Direction = "before"
	// DirectionAfter will order a list of events by newest age first (?).
	DirectionAfter Direction = "after"
)

// Category is an event category. Can be used as a filter in EventsService.List.
type Category string

const (
	// CategoryAlert are Events coming from Alerts.
	CategoryAlert Category = "ALERT"
	// CategoryCustom are custom events sent by the user.
	CategoryCustom Category = "CUSTOM"
	// CategoryDocker are events emitted by Docker.
	CategoryDocker Category = "DOCKER"
	// CategoryContainerd are events emitted by containerd.
	CategoryContainerd Category = "CONTAINERD"
	// CategoryKubernetes are events emitted by Kubernetes.
	CategoryKubernetes Category = "KUBERNETES"
)

// Categories is a type encapsulating a slice of Category to allow for easy
// marshalling into the proper JSON field.
type Categories []Category

// Status is an event status. Can be used as a filter in EventsService.List.
type Status string

const (
	// StatusTriggered is an event status indicating the event has not
	// been acknowledged.
	StatusTriggered Status = "triggered"
	// StatusResolved is an event status indicating that the event has
	// been resolved.
	StatusResolved Status = "resolved"
	// StatusAcknowledged is an event status indicating that the event
	// has been acknowledged by the user.
	StatusAcknowledged Status = "acknowledged"
	// StatusUnacknowledged is an event status indicating that the event
	// has been unacknowledged by the user.
	StatusUnacknowledged Status = "unacknowledged"
)

// Event describes an event from the Sysdig API.
type Event struct {
	ID          string            `json:"id"`
	Version     int               `json:"version"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Severity    SeverityLabel     `json:"severity"`
	Scope       string            `json:"scope"`
	Timestamp   MilliTime         `json:"timestamp"`
	CreatedOn   MilliTime         `json:"createdOn"`
	ScopeLabels map[string]string `json:"scopeLabels"`
	Tags        map[string]string `json:"tags"`
	Type        Category          `json:"type"`
}

// EventOptions are the parameters that make up an Event. To be used with EventsService.Create.
type EventOptions struct {
	// Name is the name of the event.
	Name string `json:"name"`
	// Description is a description of the event.
	Description string `json:"description,omitempty"`
	// Timestamp is the MilliTime an event occurred.
	Timestamp MilliTime `json:"timestamp,omitempty"`
	// Severity is the Severity to the associated with the event.
	Severity SeverityLabel `json:"severity,omitempty"`
	// Scope defines the scope of the event. Only ScopeSelectionIs ScopeSelections allowed during creation.
	Scope string `json:"scope,omitempty"`
	// Tags are optional tags to be added to the event.
	Tags map[string]string `json:"tags,omitempty"`
}

// EventResponse describes an EventResponse returned from the Sysdig API.
type EventResponse struct {
	Event Event `json:"event"`
}

// Get retrieves an Event.
func (s *EventsService) Get(ctx context.Context, eventID string) (*EventResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v2/events/%s", eventID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(EventResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// ListEventsResponse describes a response returned from the Sysdig List API.
type ListEventsResponse struct {
	Total   int     `json:"total"`
	Matched int     `json:"matched"`
	Events  []Event `json:"events"`
}

// ListEventOptions defines the search parameters for EventsService.List.
type ListEventOptions struct {
	// Filter can filter events by name.
	Filter string
	// AlertStatus filters events to the matching Status.
	AlertStatus Status
	// Categories filters events to the matching Categories.
	Categories Categories
	// Direction orders the list of events.
	Direction Direction
	// Scope filters events based on the Scope
	Scope string
	// Limit limits the number of events to retrieve, default 100.
	Limit int
	// Pivot is the Event ID to be used as a pivot.
	Pivot string
	// From is the timestamp for the beginning of the events to retrieve.
	From MilliTime
	// To is the timestamp for the end of the events to retrieve.
	To MilliTime
	// IncludeTotal determines whether the return the total count of events and not just the matched events.
	IncludeTotal bool
}

// List lists events with the given ListEventOptions.
func (s *EventsService) List(ctx context.Context, options ListEventOptions) (*ListEventsResponse, *http.Response, error) {
	type listEventOptions struct {
		Filter      string     `url:"filter,omitempty"`
		AlertStatus Status     `url:"alertStatus,omitempty"`
		Categories  Categories `url:"category,comma,omitempty"`
		Direction   Direction  `url:"dir,omitempty"`
		Feed        bool       `url:"feed,omitempty"`
		Limit       int        `url:"limit,omitempty"`
		Pivot       string     `url:"pivot,omitempty"`
		From        MilliTime  `url:"from,omitempty"`
		To          MilliTime  `url:"to,omitempty"`
		Scope       string     `url:"scope,omitempty"`

		IncludePivot bool `url:"include_pivot"`
		IncludeTotal bool `url:"include_total"`
	}

	u := "api/v2/events"
	o := listEventOptions{
		Filter:       options.Filter,
		AlertStatus:  options.AlertStatus,
		Categories:   options.Categories,
		Direction:    options.Direction,
		Limit:        options.Limit,
		Pivot:        options.Pivot,
		Scope:        options.Scope,
		IncludeTotal: options.IncludeTotal,
		From:         options.From,
		To:           options.To,
		Feed:         true,
		IncludePivot: true,
	}

	u, err := addOptions(u, o)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ListEventsResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Create creates an event.
func (s *EventsService) Create(ctx context.Context, event EventOptions) (*EventResponse, *http.Response, error) {
	type eventRequest struct {
		Event EventOptions `json:"event"`
	}
	u := "api/v2/events"
	req, err := s.client.NewRequest(http.MethodPost, u, eventRequest{event})
	if err != nil {
		return nil, nil, err
	}
	c := new(EventResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Delete deletes an event.
func (s *EventsService) Delete(ctx context.Context, eventID string) (*http.Response, error) {
	u := fmt.Sprintf("api/v2/events/%s", eventID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}
