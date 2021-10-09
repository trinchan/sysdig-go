package sysdig

import (
	"context"
	"fmt"
	"net/http"
)

// DashboardService is the Service for communicating with the Sysdig Monitor Dashboard related API.
type DashboardService service

// Dashboard is the structure for a Sysdig Dashboard.
// See: https://docs.sysdig.com/en/docs/sysdig-monitor/dashboards/
type Dashboard struct {
	ID                   int                  `json:"id,omitempty"`
	TeamID               int                  `json:"teamId"`
	UserID               int                  `json:"userId,omitempty"`
	Name                 string               `json:"name"`
	Panels               []Panel              `json:"panels"`
	EventDisplaySettings EventDisplaySettings `json:"eventDisplaySettings"`
	Shared               bool                 `json:"shared"`
	Public               bool                 `json:"public"`
	Version              int                  `json:"version,omitempty"`
	CreatedOn            MilliTime            `json:"createdOn"`
	ModifiedOn           MilliTime            `json:"modifiedOn"`
	Description          string               `json:"description"`
	Layout               []Layout             `json:"layout"`
	SharingSettings      []SharingSetting     `json:"sharingSettings"`
	PublicNotation       bool                 `json:"publicNotation"`
	PublicToken          string               `json:"publicToken"`
	Favorite             bool                 `json:"favorite"`
	Schema               int                  `json:"schema"`
	Username             string               `json:"username"`
	Permissions          []string             `json:"permissions"`
	ScopeExpressionList  []ScopeExpression    `json:"scopeExpressionList,omitempty"`
}

// NewDashboard constructs an empty Dashboard with the given name.
func NewDashboard(name string) *Dashboard {
	return &Dashboard{
		Name:   name,
		Schema: 3,
	}
}

// DashboardResponse is a container for a Dashboard returned by the DashboardService API.
type DashboardResponse struct {
	Dashboard Dashboard `json:"dashboard"`
}

// Get retrieves a Dashboard.
func (s *DashboardService) Get(ctx context.Context, dashboardID int) (*DashboardResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v3/dashboards/%d", dashboardID)
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(DashboardResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// ListDashboardsResponse is a container for Dashboards returned by the DashboardService.List API.
type ListDashboardsResponse struct {
	Dashboards []Dashboard `json:"dashboards"`
}

// List lists all Dashboards.
func (s *DashboardService) List(ctx context.Context) (*ListDashboardsResponse, *http.Response, error) {
	u := "api/v3/dashboards"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ListDashboardsResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Create creates a new Dashboard.
func (s *DashboardService) Create(ctx context.Context, dashboard Dashboard) (*DashboardResponse, *http.Response, error) {
	type dashboardRequest struct {
		Dashboard Dashboard `json:"dashboard"`
	}
	u := "api/v3/dashboards"
	if dashboard.Panels == nil {
		dashboard.Panels = make([]Panel, 0)
		dashboard.Layout = make([]Layout, 0)
	}
	dashboard.ID = 0
	dashboard.Version = 0
	dashboard.Schema = 3
	req, err := s.client.NewRequest(http.MethodPost, u, dashboardRequest{dashboard})
	if err != nil {
		return nil, nil, err
	}
	c := new(DashboardResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Delete deletes a Dashboard.
func (s *DashboardService) Delete(ctx context.Context, id int) (*DashboardResponse, *http.Response, error) {
	u := fmt.Sprintf("api/v3/dashboards/%d", id)
	req, err := s.client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(DashboardResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Update updates a Dashboard.
func (s *DashboardService) Update(ctx context.Context, dashboard Dashboard) (*DashboardResponse, *http.Response, error) {
	type dashboardRequest struct {
		Dashboard Dashboard `json:"dashboard"`
	}
	u := fmt.Sprintf("api/v3/dashboards/%d", dashboard.ID)
	req, err := s.client.NewRequest(http.MethodPut, u, dashboardRequest{dashboard})
	if err != nil {
		return nil, nil, err
	}
	c := new(DashboardResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// Favorite favorites or unfavorites a Dashboard.
func (s *DashboardService) Favorite(ctx context.Context, id int, favorite bool) (*DashboardResponse, *http.Response, error) {
	type favoriteRequest struct {
		Favorite bool `json:"favorite"`
	}
	u := fmt.Sprintf("api/v3/dashboards/%d", id)
	req, err := s.client.NewRequest(http.MethodPatch, u, favoriteRequest{favorite})
	if err != nil {
		return nil, nil, err
	}
	c := new(DashboardResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// DashboardTransferResponse is a container for a DashboardTransferResults for the DashboardService.Transfer API.
type DashboardTransferResponse struct {
	Results DashboardTransferResults `json:"results"`
}

// DashboardTransferResults is the response structure for the DashboardService.Transfer API.
type DashboardTransferResults struct {
	ID              int              `json:"id"`
	Name            string           `json:"name"`
	Private         bool             `json:"privateDashboard"`
	TargetTeamID    int              `json:"targetTeamId"`
	TargetTeamName  string           `json:"targetTeamName"`
	Excluded        []SharingSetting `json:"sharingSettingsExcluded"`
	Kept            []SharingSetting `json:"sharingSettingsKept"`
	CurrentTeamID   int              `json:"currentTeamId"`
	CurrentTeamName string           `json:"currentTeamName"`
}

// Transfer transfers the ownership of a set of dashboards to another user.
func (s *DashboardService) Transfer(
	ctx context.Context,
	ownerID, targetOwnerID int,
	simulate bool,
	dashboardIDs ...int) (*DashboardTransferResponse, *http.Response, error) {
	if len(dashboardIDs) == 0 {
		return nil, nil, fmt.Errorf("DashboardService.Transfer no dashboard ids specified")
	}
	type transferRequest struct {
		OwnerID                     int   `json:"ownerId"`
		TargetOwnerID               int   `json:"targetOwnerId"`
		Simulate                    bool  `json:"simulate"`
		DashboardIDsToBeTransferred []int `json:"dashboardIdsToBeTransferred"`
	}
	treq := transferRequest{
		OwnerID:                     ownerID,
		TargetOwnerID:               targetOwnerID,
		Simulate:                    simulate,
		DashboardIDsToBeTransferred: dashboardIDs,
	}
	u := "api/v3/dashboards/transfer"
	req, err := s.client.NewRequest(http.MethodPost, u, treq)
	if err != nil {
		return nil, nil, err
	}
	c := new(DashboardTransferResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// SharingSetting defines a sharing setting for a Dashboard.
type SharingSetting struct {
	Role   string        `json:"role"`
	Member SharingMember `json:"member"`
}

// SharingMember defines a sharing member for a Dashboard.
type SharingMember struct {
	Type      string `json:"type"`
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TeamTheme string `json:"teamTheme"`
}

// ScopeExpression is a scope expression used in a Dashboard.
type ScopeExpression struct {
	Operand     string   `json:"operand"`
	Operator    string   `json:"operator"`
	DisplayName string   `json:"displayName"`
	Value       []string `json:"value"`
	Descriptor  *string  `json:"descriptor,omitempty"`
	Variable    bool     `json:"variable"`
	IsVariable  bool     `json:"isVariable"`
}

// Layout defines the Layout of Panels a Dashboard.
type Layout struct {
	PanelID int `json:"panelId"`
	X       int `json:"x"`
	Y       int `json:"y"`
	W       int `json:"w"`
	H       int `json:"h"`
}

// Panel is the structure of a Panel in a Dashboard.
type Panel struct {
	ID                     int                 `json:"id"`
	Type                   string              `json:"type"`
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	NullValueDisplayText   *string             `json:"nullValueDisplayText"`
	BasicQueries           []BasicQuery        `json:"basicQueries,omitempty"`
	NumberThresholds       Thresholds          `json:"numberThresholds,omitempty"`
	ApplyScopeToAll        bool                `json:"applyScopeToAll,omitempty"`
	ApplySegmentationToAll bool                `json:"applySegmentationToAll,omitempty"`
	LegendConfiguration    LegendConfiguration `json:"legendConfiguration,omitempty"`
	AxesConfiguration      AxesConfiguration   `json:"axesConfiguration,omitempty"`
	MarkdownSource         string              `json:"markdownSource,omitempty"`
	TransparentBackground  bool                `json:"transparentBackground,omitempty"`
	PanelTitleVisible      bool                `json:"panelTitleVisible,omitempty"`
	TextAutosized          bool                `json:"textAutosized,omitempty"`
}

// BasicQuery is a basic query type used in a Dashboard.
type BasicQuery struct {
	Enabled      bool                   `json:"enabled"`
	DisplayInfo  BasicQueryDisplayInfo  `json:"displayInfo"`
	Format       BasicQueryFormat       `json:"format"`
	Scope        BasicQueryScope        `json:"scope"`
	CompareTo    BasicQueryCompareTo    `json:"compareTo"`
	Metrics      []BasicQueryMetric     `json:"metrics"`
	Segmentation BasicQuerySegmentation `json:"segmentation,omitempty"`
}

// BasicQueryCompareTo is used in a BasicQuery on a Dashboard.
type BasicQueryCompareTo struct {
	Enabled    bool   `json:"enabled"`
	Delta      int    `json:"delta"`
	TimeFormat string `json:"timeFormat"`
}

// BasicQueryScope is a scope used in a BasicQuery on a Dashboard.
type BasicQueryScope struct {
	Expressions           []string `json:"expressions"`
	ExtendsDashboardScope bool     `json:"extendsDashboardScope"`
}

// BasicQueryMetric is a metric used in a BasicQuery on a Dashboard.
type BasicQueryMetric struct {
	ID               string      `json:"id"`
	TimeAggregation  string      `json:"timeAggregation"`
	GroupAggregation string      `json:"groupAggregation"`
	Descriptor       *string     `json:"descriptor,omitempty"`
	Sorting          interface{} `json:"sorting"`
}

// BasicQueryDisplayInfo is the display info used in a BasicQuery on a Dashboard.
type BasicQueryDisplayInfo struct {
	DisplayName                   string `json:"displayName"`
	TimeSeriesDisplayNameTemplate string `json:"timeSeriesDisplayNameTemplate"`
	Type                          string `json:"type"`
}

// BasicQueryFormat is the metric format used in a BasicQuery on a Dashboard.
type BasicQueryFormat struct {
	Unit                 string `json:"unit"`
	InputFormat          string `json:"inputFormat"`
	DisplayFormat        string `json:"displayFormat"`
	Decimals             *int   `json:"decimals"`
	YAxis                string `json:"yAxis"`
	NullValueDisplayMode string `json:"nullValueDisplayMode"`
}

// BasicQuerySegmentation is the segmentation used in a BasicQuery on a Dashboard.
type BasicQuerySegmentation struct {
	Labels    []BasicQuerySegmentationLabel `json:"labels"`
	Limit     int                           `json:"limit"`
	Direction string                        `json:"direction"`
}

// BasicQuerySegmentationLabel is a segmentation label used in a BasicQuerySegmentation of a BasicQuery on a Dashboard.
type BasicQuerySegmentationLabel struct {
	ID          string  `json:"id"`
	Descriptor  *string `json:"descriptor,omitempty"`
	DisplayName *string `json:"displayName,omitempty"`
	Sorting     *string `json:"sorting,omitempty"`
}

// Thresholds are the threshold for a Panel on a Dashboard.
type Thresholds struct {
	Values      []ThresholdValue `json:"values"`
	Base        ThresholdBase    `json:"base"`
	UseDefaults *bool            `json:"useDefaults"`
}

// ThresholdValue is a threshold value for a Thresholds on a Panel in a Dashboard.
type ThresholdValue struct {
	Severity    string  `json:"severity"`
	Value       float64 `json:"value"`
	InputFormat string  `json:"inputFormat"`
	DisplayText string  `json:"displayText"`
}

// ThresholdBase is a threshold base for a Thresholds on a Panel in a Dashboard.
type ThresholdBase struct {
	Severity    string `json:"severity"`
	DisplayText string `json:"displayText"`
}

// LegendConfiguration is the configuration for a legend on a Panel in a Dashboard.
type LegendConfiguration struct {
	Enabled     bool     `json:"enabled"`
	Position    string   `json:"position"`
	Layout      string   `json:"layout"`
	ShowCurrent bool     `json:"showCurrent"`
	Width       *float64 `json:"width"`
	Height      *float64 `json:"height"`
}

// AxesConfiguration is the configuration for the axes of a Panel in a Dashboard.
type AxesConfiguration struct {
	Bottom struct {
		Enabled bool `json:"enabled"`
	} `json:"bottom"`
	Left  Axis `json:"left"`
	Right Axis `json:"right"`
}

// Axis is an axis configuration used in an AxesConfiguration for the axes of a Panel in a Dashboard.
type Axis struct {
	Enabled        bool        `json:"enabled"`
	DisplayName    *string     `json:"displayName"`
	Unit           string      `json:"unit"`
	DisplayFormat  string      `json:"displayFormat"`
	Decimals       interface{} `json:"decimals"`
	MinValue       *float64    `json:"minValue"`
	MaxValue       *float64    `json:"maxValue"`
	MinInputFormat string      `json:"minInputFormat"`
	MaxInputFormat string      `json:"maxInputFormat"`
	Scale          string      `json:"scale"`
}

// EventDisplaySettings are the event display settings for a Dashboard.
type EventDisplaySettings struct {
	Enabled     bool                            `json:"enabled"`
	QueryParams EventDisplaySettingsQueryParams `json:"queryParams"`
}

// EventDisplaySettingsQueryParams are the query parameters used in an EventDisplaySettings for a Dashboard.
type EventDisplaySettingsQueryParams struct {
	Severities    []Severity `json:"severities"`
	AlertStatuses []Status   `json:"alertStatuses"`
	Categories    []Category `json:"categories"`
	Filter        *string    `json:"filter"`
	TeamScope     bool       `json:"teamScope"`
}
