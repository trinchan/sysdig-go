package sysdig

import (
	"context"
	"fmt"
	"net/http"
)

// TeamsService is the Service for communicating with the Sysdig Monitor Team related API.
type TeamsService service

type ProductType string

const (
	ProductTypeSDC ProductType = "SDC"
	ProductTypeSDS ProductType = "SDS"
	ProductTypeAny ProductType = ""
)

type Team struct {
	Version     int       `json:"version"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"`
	LastUpdated MilliTime `json:"lastUpdated"`
	DateCreated MilliTime `json:"dateCreated"`
	// TODO what is this structure?
	NamespaceFilters    interface{}    `json:"namespaceFilters"`
	CustomerId          int            `json:"customerId"`
	Show                string         `json:"show"`
	Products            []string       `json:"products"`
	Theme               string         `json:"theme"`
	EntryPoint          TeamEntryPoint `json:"entryPoint"`
	DefaultTeamRole     string         `json:"defaultTeamRole"`
	Immutable           bool           `json:"immutable"`
	CanUseSysdigCapture bool           `json:"canUseSysdigCapture"`
	CanUseAgentCli      bool           `json:"canUseAgentCli"`
	CanUseCustomEvents  bool           `json:"canUseCustomEvents"`
	CanUseAwsMetrics    bool           `json:"canUseAwsMetrics"`
	CanUseBeaconMetrics bool           `json:"canUseBeaconMetrics"`
	CanUseRapidResponse bool           `json:"canUseRapidResponse"`
	UserCount           int            `json:"userCount"`
	Name                string         `json:"name"`
	// TODO what is this structure?
	Properties interface{} `json:"properties"`
	ID         int         `json:"id"`
	Default    bool        `json:"default"`
}

type TeamEntryPoint struct {
	Module string `json:"module"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

// Get gets a Team.
func (s *TeamsService) Get(ctx context.Context, teamID int) (*TeamResponse, *http.Response, error) {
	u := fmt.Sprintf("api/team/%d", teamID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(TeamResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

type ListTeamsResponse struct {
	Teams []Team `json:"teams"`
}

// List returns the list of Teams for the given ProductType.
func (s *TeamsService) List(ctx context.Context, product ProductType) (*ListTeamsResponse, *http.Response, error) {
	u := "api/team"
	type listOptions struct {
		Product ProductType `url:"product"`
	}
	uWithOpts, err := addOptions(u, listOptions{Product: product})
	if err != nil {
		return nil, nil, err
	}
	req, err := s.client.NewRequest(http.MethodGet, uWithOpts, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ListTeamsResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

type ListUsersResponse struct {
	Offset int    `json:"offset"`
	Total  int    `json:"total"`
	Users  []User `json:"users"`
}

// ListUsers returns the list of Users for the given Team.
func (s *TeamsService) ListUsers(ctx context.Context, teamID int) (*ListUsersResponse, *http.Response, error) {
	u := fmt.Sprintf("api/team/%d/users", teamID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ListUsersResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

func (s *TeamsService) Delete(ctx context.Context, teamID int) (*http.Response, error) {
	u := fmt.Sprintf("api/team/%d", teamID)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(ctx, req, nil)
	return resp, err
}

type InfrastructureResponse struct {
	Infrastructure Infrastructure `json:"infrastructure"`
}

type Infrastructure struct {
	HostCount        int `json:"hostCount"`
	ContainerCount   int `json:"containerCount"`
	UnresolvedEvents int `json:"unresolvedEvents"`
	// TODO figure out the structure of these things
	Orchestrations      []interface{}                      `json:"orchestrations"`
	Platforms           []interface{}                      `json:"platforms"`
	ContainerTypes      []interface{}                      `json:"containerTypes"`
	MetricCount         InfrastructureMetricCount          `json:"metricCount"`
	OnPremOverview      OnPremOverview                     `json:"onPremOverview"`
	AgentMetricOverview InfrastructureAgentMetricOverviews `json:"agentMetricOverview"`
}

type InfrastructureMetricCount struct {
	Total    int `json:"total"`
	JMX      int `json:"jmx"`
	StatsD   int `json:"statsD"`
	AppCheck int `json:"appCheck"`
}

type InfrastructureAgentMetricOverviews struct {
	ExceedingLimitCount int     `json:"exceedingLimitCount"`
	TotalAgents         int     `json:"totalAgents"`
	ExceedingLimitPct   float64 `json:"exceedingLimitPct"`
}

type OnPremOverview struct {
	LatestVersion   string `json:"latestVersion"`
	CustomerVersion string `json:"customerVersion"`
	ShowPlanInfo    bool   `json:"showPlanInfo"`
}

func (s *TeamsService) Infrastructure(ctx context.Context) (*InfrastructureResponse, *http.Response, error) {
	u := "api/team/infrastructure"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(InfrastructureResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}
