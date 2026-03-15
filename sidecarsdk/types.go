package sidecarsdk

import (
	"time"

	"github.com/google/uuid"
)

// AuthenticateResponse is returned after a successful instance identity authentication.
type AuthenticateResponse struct {
	SessionToken string `json:"session_token"`
}

// GoogleInstanceIdentityToken is the payload for Google Cloud instance identity authentication.
type GoogleInstanceIdentityToken struct {
	JSONWebToken string `json:"json_web_token"`
}

// AWSInstanceIdentityToken is the payload for AWS instance identity authentication.
type AWSInstanceIdentityToken struct {
	Signature string `json:"signature"`
	Document  string `json:"document"`
}

// AzureInstanceIdentityToken is the payload for Azure instance identity authentication.
type AzureInstanceIdentityToken struct {
	Signature string `json:"signature"`
	Encoding  string `json:"encoding"`
}

// PostStartupRequest reports the sidecar's startup information to Runtime.
type PostStartupRequest struct {
	Version           string `json:"version"`
	ExpandedDirectory string `json:"expanded_directory"`
	Subsystems        []string `json:"subsystems"`
}

// AgentSidecarLifecycle represents the lifecycle state of a sidecar.
type AgentSidecarLifecycle string

const (
	AgentSidecarLifecycleCreated         AgentSidecarLifecycle = "created"
	AgentSidecarLifecycleStarting        AgentSidecarLifecycle = "starting"
	AgentSidecarLifecycleStartTimeout    AgentSidecarLifecycle = "start_timeout"
	AgentSidecarLifecycleStartError      AgentSidecarLifecycle = "start_error"
	AgentSidecarLifecycleReady           AgentSidecarLifecycle = "ready"
	AgentSidecarLifecycleShuttingDown    AgentSidecarLifecycle = "shutting_down"
	AgentSidecarLifecycleShutdownTimeout AgentSidecarLifecycle = "shutdown_timeout"
	AgentSidecarLifecycleShutdownError   AgentSidecarLifecycle = "shutdown_error"
	AgentSidecarLifecycleOff             AgentSidecarLifecycle = "off"
)

// PostLifecycleRequest reports a lifecycle state change to Runtime.
type PostLifecycleRequest struct {
	State     AgentSidecarLifecycle `json:"state"`
	ChangedAt time.Time             `json:"changed_at"`
}

// PostMetadataRequest sends metadata key-value pairs to Runtime.
type PostMetadataRequest struct {
	Metadata []MetadataItem `json:"metadata"`
}

// MetadataItem is a single metadata key-value pair with collection time.
type MetadataItem struct {
	Key         string    `json:"key"`
	CollectedAt time.Time `json:"collected_at"`
	Value       string    `json:"value"`
	Error       string    `json:"error"`
}

// Stats contains connection and session statistics reported by the sidecar.
type Stats struct {
	ConnectionsByProto      map[string]int64 `json:"conns_by_proto"`
	ConnectionCount         int64            `json:"connection_count"`
	ConnectionMedianLatencyMS float64        `json:"connection_median_latency_ms"`
	RxPackets               int64            `json:"rx_packets"`
	RxBytes                 int64            `json:"rx_bytes"`
	TxPackets               int64            `json:"tx_packets"`
	TxBytes                 int64            `json:"tx_bytes"`
	SessionCountVSCode      int64            `json:"session_count_vscode"`
	SessionCountJetBrains   int64            `json:"session_count_jetbrains"`
	SessionCountReconnectingPTY int64        `json:"session_count_reconnecting_pty"`
	SessionCountSSH         int64            `json:"session_count_ssh"`
	Metrics                 []SidecarMetric  `json:"metrics"`
}

// SidecarMetric is a single metric reported by the sidecar.
type SidecarMetric struct {
	Name   string              `json:"name"`
	Type   SidecarMetricType   `json:"type"`
	Value  float64             `json:"value"`
	Labels []SidecarMetricLabel `json:"labels"`
}

// SidecarMetricType identifies the type of a sidecar metric.
type SidecarMetricType string

const (
	SidecarMetricTypeCounter SidecarMetricType = "counter"
	SidecarMetricTypeGauge   SidecarMetricType = "gauge"
)

// SidecarMetricLabel is a key-value label attached to a metric.
type SidecarMetricLabel struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// StatsResponse is returned after posting stats, indicating the next report interval.
type StatsResponse struct {
	// ReportInterval is the recommended interval between stats reports, in nanoseconds.
	ReportInterval int64 `json:"report_interval"`
}

// PostLogSourceRequest registers a new log source with Runtime.
type PostLogSourceRequest struct {
	ID          uuid.UUID `json:"id"`
	DisplayName string    `json:"display_name"`
	Icon        string    `json:"icon"`
}

// LogSource is the response after registering a log source.
type LogSource struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Icon        string    `json:"icon"`
}

// PatchLogsRequest sends log entries for a given log source.
type PatchLogsRequest struct {
	LogSourceID uuid.UUID `json:"log_source_id"`
	Logs        []Log     `json:"logs"`
}

// Log is a single log entry.
type Log struct {
	CreatedAt time.Time `json:"created_at"`
	Output    string    `json:"output"`
	Level     LogLevel  `json:"level"`
}

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// GitSSHKey contains the SSH key pair used for Git operations.
type GitSSHKey struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// ExternalAuthRequest contains parameters for requesting external authentication.
type ExternalAuthRequest struct {
	ID    string `json:"id"`
	Match string `json:"match"`
	Listen bool  `json:"listen"`
}

// ExternalAuthResponse contains the result of an external authentication request.
type ExternalAuthResponse struct {
	AccessToken string                 `json:"access_token"`
	TokenExtra  map[string]interface{} `json:"token_extra"`
	URL         string                 `json:"url"`
	Type        string                 `json:"type"`
}

// Manifest contains the full configuration manifest for a sidecar.
type Manifest struct {
	SidecarID                uuid.UUID         `json:"sidecar_id"`
	SidecarName              string            `json:"sidecar_name"`
	OwnerName                string            `json:"owner_name"`
	AgentID                  uuid.UUID         `json:"agent_id"`
	AgentName                string            `json:"agent_name"`
	GitAuthConfigs           int               `json:"git_auth_configs"`
	VSCodePortProxyURI       string            `json:"vscode_port_proxy_uri"`
	EnvironmentVariables     map[string]string `json:"environment_variables"`
	Directory                string            `json:"directory"`
	MOTDFile                 string            `json:"motd_file"`
	DisableDirectConnections bool              `json:"disable_direct_connections"`
	Scripts                  []AgentScript     `json:"scripts"`
	Apps                     []AgentApp        `json:"apps"`
}

// AgentScript represents a script configured to run on the sidecar.
type AgentScript struct {
	ID               uuid.UUID `json:"id"`
	LogSourceID      uuid.UUID `json:"log_source_id"`
	LogPath          string    `json:"log_path"`
	Script           string    `json:"script"`
	Cron             string    `json:"cron"`
	RunOnStart       bool      `json:"run_on_start"`
	RunOnStop        bool      `json:"run_on_stop"`
	StartBlocksLogin bool      `json:"start_blocks_login"`
	// Timeout is the script execution timeout in nanoseconds.
	Timeout int64 `json:"timeout"`
}

// AgentApp represents an application running on the sidecar.
type AgentApp struct {
	ID            uuid.UUID   `json:"id"`
	URL           string      `json:"url"`
	External      bool        `json:"external"`
	Slug          string      `json:"slug"`
	DisplayName   string      `json:"display_name"`
	Command       string      `json:"command,omitempty"`
	Icon          string      `json:"icon,omitempty"`
	Subdomain     bool        `json:"subdomain"`
	SubdomainName string      `json:"subdomain_name,omitempty"`
	SharingLevel  string      `json:"sharing_level"`
	Healthcheck   Healthcheck `json:"healthcheck"`
	Health        string      `json:"health"`
}

// Healthcheck configures health checking for an agent app.
type Healthcheck struct {
	URL       string `json:"url"`
	Interval  int32  `json:"interval"`
	Threshold int32  `json:"threshold"`
}

// SidecarSubsystem identifies a sidecar subsystem.
type SidecarSubsystem string
