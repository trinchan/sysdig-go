package sysdig

import (
	"context"
	"net/http"
)

// UsersService is the Service for communicating with the Sysdig Monitor User related API.
type UsersService service

// User defines a Sysdig User.
type User struct {
	TermsAndConditions   bool               `json:"termsAndConditions"`
	Timezone             string             `json:"timezone"`
	PictureURL           string             `json:"pictureUrl"`
	CustomerSettings     CustomerSettings   `json:"customerSettings"`
	Customer             Customer           `json:"customer"`
	Oauth                bool               `json:"oauth"`
	AgentInstallParams   AgentInstallParams `json:"agentInstallParams"`
	Properties           UserProperties     `json:"properties"`
	ResetPassword        bool               `json:"resetPassword"`
	AdditionalRoles      []interface{}      `json:"additionalRoles"` // TODO What is the format of this...
	TeamRoles            []TeamRole         `json:"teamRoles"`
	LastUpdated          MilliTime          `json:"lastUpdated"`
	AccessKey            string             `json:"accessKey"`
	IntercomUserIDHash   string             `json:"intercomUserIdHash"`
	UniqueIntercomUserID string             `json:"uniqueIntercomUserId"`
	CurrentTeam          int                `json:"currentTeam"`
	Enabled              bool               `json:"enabled"`
	Version              int                `json:"version"`
	DateCreated          MilliTime          `json:"dateCreated"`
	Status               string             `json:"status"`
	Products             []string           `json:"products"`
	FirstName            string             `json:"firstName"`
	LastName             string             `json:"lastName"`
	SystemRole           string             `json:"systemRole"`
	Username             string             `json:"username"`
	LastSeen             int64              `json:"lastSeen"`
	Name                 string             `json:"name"`
	ID                   int                `json:"id"`
}

// UserProperties are the properties for a User.
type UserProperties struct {
	ResetPassword          bool   `json:"resetPassword"`
	OpenIDConnectProfileID string `json:"OpenID Connect profile id"`
	IAMID                  string `json:"iamId"`
	OpenID                 bool   `json:"openid"`
	UserEmailAlias         string `json:"user_email_alias"`
	HasBeenInvited         bool   `json:"has_been_invited"`
}

// AgentInstallParams are the agent installation parameters for a User.
type AgentInstallParams struct {
	AccessKey        string `json:"accessKey"`
	CollectorAddress string `json:"collectorAddress"`
	CollectorPort    int    `json:"collectorPort"`
	CheckCertificate bool   `json:"checkCertificate"`
	SSLEnabled       bool   `json:"sslEnabled"`
}

// Customer is the customer information for a User.
type Customer struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	AccessKey   string    `json:"accessKey"`
	ExternalID  string    `json:"externalId"`
	DateCreated MilliTime `json:"dateCreated"`
}

// TeamRole is the role for a User in a Team.
type TeamRole struct {
	TeamID    int    `json:"teamId"`
	TeamName  string `json:"teamName"`
	TeamTheme string `json:"teamTheme"`
	UserID    int    `json:"userId"`
	UserName  string `json:"userName"`
	Role      string `json:"role"`
	Admin     bool   `json:"admin"`
}

// CustomerSettings are the customer related settings for a user.
type CustomerSettings struct {
	Sysdig      UserSysdigSettings `json:"sysdig"`
	Plan        Plan               `json:"plan"`
	Environment interface{}        `json:"environment"` // TODO What is the format of this...
}

// UserSysdigSettings are the Sysdig settings for a user.
type UserSysdigSettings struct {
	Enabled    bool          `json:"enabled"`
	EnabledSSE bool          `json:"enabledSSE"`
	Buckets    []interface{} `json:"buckets"` // TODO What is the format of this...
}

// Plan is the plan for a User.
type Plan struct {
	MaxAgents                 int                   `json:"maxAgents"`
	OnDemandAgents            int                   `json:"onDemandAgents"`
	MaxTeams                  int                   `json:"maxTeams"`
	Timelines                 []Timeline            `json:"timelines"`
	MetricsSettings           MetricsSettings       `json:"metricsSettings"`
	SecureEnabled             bool                  `json:"secureEnabled"`
	MonitorEnabled            bool                  `json:"monitorEnabled"`
	AllocatedAgentsCount      int                   `json:"allocatedAgentsCount"`
	PaymentsIntegrationID     PaymentsIntegrationID `json:"paymentsIntegrationId"`
	PricingPlan               string                `json:"pricingPlan"`
	IndirectCustomer          bool                  `json:"indirectCustomer"`
	TrialPlanName             string                `json:"trialPlanName"`
	Partner                   string                `json:"partner"`
	MigratedToV2Direct        bool                  `json:"migratedToV2Direct"`
	OverageAssessmentEligible bool                  `json:"overageAssessmentEligible"`
}

// Timeline are the sampling timelines for a Plan.
type Timeline struct {
	From     *MilliTime `json:"from"`
	To       *MilliTime `json:"to"`
	Sampling int64      `json:"sampling"`
}

// MetricsSettings are the metrics settings for a User.
type MetricsSettings struct {
	Enforce                       bool         `json:"enforce"`
	ShowExperimentals             bool         `json:"showExperimentals"`
	Limits                        Limits       `json:"limits"`
	LegacyLimits                  LegacyLimits `json:"legacyLimits"`
	EnforceAgentAggregation       bool         `json:"enforceAgentAggregation"`
	EnablePromCalculatedIngestion bool         `json:"enablePromCalculatedIngestion"`
}

// Limits are limits for a User.
type Limits struct {
	JMX                      int     `json:"jmx"`
	Statsd                   int     `json:"statsd"`
	AppCheck                 int     `json:"appCheck"`
	Prometheus               int     `json:"prometheus"`
	PrometheusPerProcess     int     `json:"prometheusPerProcess"`
	Connections              int     `json:"connections"`
	ProgAggregationCount     int     `json:"progAggregationCount"`
	AppCheckAggregationCount int     `json:"appCheckAggregationCount"`
	PromMetricsWeight        float64 `json:"promMetricsWeight"`
	TopFilesCount            int     `json:"topFilesCount"`
	TopDevicesCount          int     `json:"topDevicesCount"`
	HostServerPorts          int     `json:"hostServerPorts"`
	ContainerServerPorts     int     `json:"containerServerPorts"`
	LimitKubernetesResources bool    `json:"limitKubernetesResources"`
	KubernetesPods           int     `json:"kubernetesPods"`
	KubernetesJobs           int     `json:"kubernetesJobs"`
	ContainerDensity         int     `json:"containerDensity"`
	MeerkatSuited            bool    `json:"meerkatSuited"`
}

// LegacyLimits are legacy limits for a User.
type LegacyLimits struct {
	JMX                      int     `json:"jmx"`
	Statsd                   int     `json:"statsd"`
	AppCheck                 int     `json:"appCheck"`
	Prometheus               int     `json:"prometheus"`
	PrometheusPerProcess     int     `json:"prometheusPerProcess"`
	Connections              int     `json:"connections"`
	ProgAggregationCount     int     `json:"progAggregationCount"`
	AppCheckAggregationCount int     `json:"appCheckAggregationCount"`
	PromMetricsWeight        float64 `json:"promMetricsWeight"`
	TopFilesCount            int     `json:"topFilesCount"`
	TopDevicesCount          int     `json:"topDevicesCount"`
	HostServerPorts          int     `json:"hostServerPorts"`
	ContainerServerPorts     int     `json:"containerServerPorts"`
	LimitKubernetesResources bool    `json:"limitKubernetesResources"`
	KubernetesPods           int     `json:"kubernetesPods"`
	KubernetesJobs           int     `json:"kubernetesJobs"`
	ContainerDensity         int     `json:"containerDensity"`
	MeerkatSuited            bool    `json:"meerkatSuited"`
}

// PaymentsIntegrationID is the ID of the payment integration for this User.
type PaymentsIntegrationID struct {
	ID  string `json:"id"`
	TTL TTL    `json:"ttl"`
}

// TTL is the PaymentsIntegrationID TTL.
type TTL struct {
	TTL int `json:"ttl"`
}

// MeResponse describes the response for UsersService.Me
type MeResponse struct {
	User User `json:"user"`
}

// Me returns information about the current User.
func (s *UsersService) Me(ctx context.Context) (*MeResponse, *http.Response, error) {
	u := "api/user/me"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(MeResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// TokenResponse describes the response for UsersService.Token.
type TokenResponse struct {
	Token Token `json:"token"`
}

// Token is a Sysdig token.
type Token struct {
	Key string `json:"key"`
}

// Token fetches the API token for this User.
func (s *UsersService) Token(ctx context.Context) (*TokenResponse, *http.Response, error) {
	u := "api/token"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(TokenResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}

// ConnectedAgentsResponse describes the response for UsersService.ConnectedAgents.
type ConnectedAgentsResponse struct {
	Total  int     `json:"total"`
	Agents []Agent `json:"agents"`
}

// Agent is a Sysdig Agent. // TODO link to docs? What is the structure of this?
type Agent struct {
	ID string `json:"id"`
}

// ConnectedAgents lists the connected agents for the user.
func (s *UsersService) ConnectedAgents(ctx context.Context) (*ConnectedAgentsResponse, *http.Response, error) {
	u := "api/agents/connected"
	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	c := new(ConnectedAgentsResponse)
	resp, err := s.client.Do(ctx, req, c)
	return c, resp, err
}
