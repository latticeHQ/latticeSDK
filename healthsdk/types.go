// Package healthsdk provides health check and diagnostics types for Lattice deployments.
package healthsdk

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// HealthSection identifies a section of the health report.
type HealthSection string

const (
	HealthSectionDERP               HealthSection = "DERP"
	HealthSectionAccessURL          HealthSection = "AccessURL"
	HealthSectionWebsocket          HealthSection = "Websocket"
	HealthSectionDatabase           HealthSection = "Database"
	HealthSectionAgentProxy         HealthSection = "AgentProxy"
	HealthSectionProvisionerDaemons HealthSection = "ProvisionerDaemons"
)

// HealthSections contains all known health sections.
var HealthSections = []HealthSection{
	HealthSectionDERP,
	HealthSectionAccessURL,
	HealthSectionWebsocket,
	HealthSectionDatabase,
	HealthSectionAgentProxy,
	HealthSectionProvisionerDaemons,
}

// Severity indicates the severity level of a health report.
type Severity string

const (
	SeverityOK      Severity = "ok"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// HealthMessage is a message associated with a health check.
type HealthMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HealthSettings contains the settings for health checks.
type HealthSettings struct {
	DismissedHealthchecks []HealthSection `json:"dismissed_healthchecks"`
}

// UpdateHealthSettings is the request body for updating health settings.
type UpdateHealthSettings struct {
	DismissedHealthchecks []HealthSection `json:"dismissed_healthchecks"`
}

// HealthcheckReport is the full health report for a Lattice deployment.
type HealthcheckReport struct {
	Time               time.Time                `json:"time"`
	Healthy            bool                     `json:"healthy"` // Deprecated: Use Severity instead.
	Severity           Severity                 `json:"severity"`
	DERP               DERPHealthReport         `json:"derp"`
	AccessURL          AccessURLReport          `json:"access_url"`
	Websocket          WebsocketReport          `json:"websocket"`
	Database           DatabaseReport           `json:"database"`
	AgentProxy         AgentProxyReport         `json:"agent_proxy"`
	ProvisionerDaemons ProvisionerDaemonsReport `json:"provisioner_daemons"`
	LatticeVersion     string                   `json:"lattice_version"`
}

// Summarize returns a human-readable summary of the entire health report.
// Each sub-report is summarized with the given documentation URL prefix.
func (r HealthcheckReport) Summarize(docsURL string) []string {
	var msgs []string
	msgs = append(msgs, r.DERP.Summarize("DERP: ", docsURL)...)
	msgs = append(msgs, r.AccessURL.Summarize("AccessURL: ", docsURL)...)
	msgs = append(msgs, r.Websocket.Summarize("Websocket: ", docsURL)...)
	msgs = append(msgs, r.Database.Summarize("Database: ", docsURL)...)
	msgs = append(msgs, r.AgentProxy.Summarize("AgentProxy: ", docsURL)...)
	msgs = append(msgs, r.ProvisionerDaemons.Summarize("ProvisionerDaemons: ", docsURL)...)
	return msgs
}

// BaseReport contains fields common to all health sub-reports.
type BaseReport struct {
	Error     *string         `json:"error"`
	Severity  Severity        `json:"severity"`
	Warnings  []HealthMessage `json:"warnings"`
	Dismissed bool            `json:"dismissed"`
}

// Summarize returns a human-readable summary of the base report.
func (b BaseReport) Summarize(prefix, docsURL string) []string {
	var msgs []string
	if b.Error != nil {
		msgs = append(msgs, fmt.Sprintf("%sERROR: %s", prefix, *b.Error))
	}
	for _, w := range b.Warnings {
		msg := fmt.Sprintf("%sWARN: %s", prefix, w.Message)
		if docsURL != "" && w.Code != "" {
			msg += fmt.Sprintf(" (%s#%s)", docsURL, w.Code)
		}
		msgs = append(msgs, msg)
	}
	if b.Dismissed {
		msgs = append(msgs, fmt.Sprintf("%s(dismissed)", prefix))
	}
	return msgs
}

// AccessURLReport is the health report for the access URL.
type AccessURLReport struct {
	BaseReport
	Healthy        bool   `json:"healthy"`
	AccessURL      string `json:"access_url"`
	Reachable      bool   `json:"reachable"`
	StatusCode     int    `json:"status_code"`
	HealthzResponse string `json:"healthz_response"`
}

// DERPHealthReport is the health report for the DERP mesh.
type DERPHealthReport struct {
	BaseReport
	Healthy      bool                       `json:"healthy"`
	Regions      map[int]*DERPRegionReport  `json:"regions"`
	NetcheckErr  *string                    `json:"netcheck_err"`
	NetcheckLogs []string                   `json:"netcheck_logs"`
}

// DERPRegionReport is the health report for a single DERP region.
type DERPRegionReport struct {
	Healthy     bool              `json:"healthy"`
	Severity    Severity          `json:"severity"`
	Warnings    []HealthMessage   `json:"warnings"`
	Error       *string           `json:"error"`
	Region      *DERPRegion       `json:"region"`
	NodeReports []*DERPNodeReport `json:"node_reports"`
}

// DERPRegion describes a DERP region in the health report.
type DERPRegion struct {
	RegionID   int         `json:"region_id"`
	RegionCode string      `json:"region_code"`
	RegionName string      `json:"region_name"`
	Nodes      []*DERPNode `json:"nodes"`
}

// DERPNode describes a single DERP node in a region.
type DERPNode struct {
	Name     string `json:"name"`
	RegionID int    `json:"region_id"`
	HostName string `json:"host_name"`
	IPv4     string `json:"ipv4"`
	IPv6     string `json:"ipv6"`
	STUNPort int    `json:"stun_port"`
	STUNOnly bool   `json:"stun_only"`
	DERPPort int    `json:"derp_port"`
}

// DERPNodeReport is the health report for a single DERP node.
type DERPNodeReport struct {
	Healthy             bool        `json:"healthy"`
	Severity            Severity    `json:"severity"`
	Warnings            []HealthMessage `json:"warnings"`
	Error               *string     `json:"error"`
	Node                *DERPNode   `json:"node"`
	CanExchangeMessages bool        `json:"can_exchange_messages"`
	RoundTripPing       string      `json:"round_trip_ping"`
	RoundTripPingMs     int         `json:"round_trip_ping_ms"`
	UsesWebsocket       bool        `json:"uses_websocket"`
	ClientLogs          [][]string  `json:"client_logs"`
	ClientErrs          [][]string  `json:"client_errs"`
	STUN                STUNReport  `json:"stun"`
}

// STUNReport is the health report for STUN connectivity.
type STUNReport struct {
	Enabled bool    `json:"enabled"`
	CanSTUN bool    `json:"can_stun"`
	Error   *string `json:"error"`
}

// DatabaseReport is the health report for the database.
type DatabaseReport struct {
	BaseReport
	Healthy     bool   `json:"healthy"`
	Reachable   bool   `json:"reachable"`
	Latency     string `json:"latency"`
	LatencyMS   int64  `json:"latency_ms"`
	ThresholdMS int64  `json:"threshold_ms"`
}

// WebsocketReport is the health report for websocket connectivity.
type WebsocketReport struct {
	BaseReport
	Healthy bool   `json:"healthy"`
	Body    string `json:"body"`
	Code    int    `json:"code"`
}

// AgentProxyReport is the health report for agent proxies.
type AgentProxyReport struct {
	BaseReport
	Healthy      bool         `json:"healthy"`
	AgentProxies []AgentProxy `json:"agent_proxies"`
}

// AgentProxy describes an agent proxy in the health report.
type AgentProxy struct {
	ID               uuid.UUID        `json:"id"`
	Name             string           `json:"name"`
	DisplayName      string           `json:"display_name"`
	IconURL          string           `json:"icon_url"`
	Healthy          bool             `json:"healthy"`
	PathAppURL       string           `json:"path_app_url"`
	WildcardHostname string           `json:"wildcard_hostname"`
	DERPEnabled      bool             `json:"derp_enabled"`
	DERPOnly         bool             `json:"derp_only"`
	Status           AgentProxyStatus `json:"status"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Deleted          bool             `json:"deleted"`
	Version          string           `json:"version"`
}

// AgentProxyStatus contains the latest health check result for a proxy.
type AgentProxyStatus struct {
	Status    string                 `json:"status"`
	Report    AgentProxyStatusReport `json:"report"`
	CheckedAt time.Time             `json:"checked_at"`
}

// AgentProxyStatusReport contains errors and warnings from a proxy health check.
type AgentProxyStatusReport struct {
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// ProvisionerDaemonsReport is the health report for provisioner daemons.
type ProvisionerDaemonsReport struct {
	BaseReport
	Items []ProvisionerDaemonsReportItem `json:"items"`
}

// ProvisionerDaemonsReportItem is a single provisioner daemon in the health report.
type ProvisionerDaemonsReportItem struct {
	ProvisionerDaemon ProvisionerDaemonInfo `json:"provisioner_daemon"`
	Warnings          []HealthMessage       `json:"warnings"`
}

// ProvisionerDaemonInfo describes a provisioner daemon in the health report.
type ProvisionerDaemonInfo struct {
	ID         uuid.UUID          `json:"id"`
	Name       string             `json:"name"`
	CreatedAt  time.Time          `json:"created_at"`
	LastSeenAt *time.Time         `json:"last_seen_at"`
	Version    string             `json:"version"`
	APIVersion string             `json:"api_version"`
	Tags       map[string]string  `json:"tags"`
}

// InterfacesReport is the health report for network interfaces.
type InterfacesReport struct {
	BaseReport
	Interfaces []NetworkInterface `json:"interfaces"`
}

// NetworkInterface describes a network interface.
type NetworkInterface struct {
	Name      string   `json:"name"`
	MTU       int      `json:"mtu"`
	Addresses []string `json:"addresses"`
}

// ClientNetcheckReport is the report from a client-side netcheck.
type ClientNetcheckReport struct {
	DERP       DERPHealthReport `json:"derp"`
	Interfaces InterfacesReport `json:"interfaces"`
}

// SidecarNetcheckReport is the report from a sidecar-side netcheck.
type SidecarNetcheckReport struct {
	BaseReport
	Interfaces InterfacesReport `json:"interfaces"`
}
