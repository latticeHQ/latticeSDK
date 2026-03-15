// Package agents provides agent lifecycle management for the Lattice Runtime API.
// It covers agents, builds, proxies, sidecars, and port shares.
package agents

import (
	"time"

	"github.com/google/uuid"
)

// AutomaticUpdates controls agent automatic update behavior.
type AutomaticUpdates string

const (
	AutomaticUpdatesAlways AutomaticUpdates = "always"
	AutomaticUpdatesNever  AutomaticUpdates = "never"
)

// AgentTransition represents agent lifecycle transitions.
type AgentTransition string

const (
	AgentTransitionStart  AgentTransition = "start"
	AgentTransitionStop   AgentTransition = "stop"
	AgentTransitionDelete AgentTransition = "delete"
)

// AgentStatus represents the current status of an agent.
type AgentStatus string

const (
	AgentStatusPending   AgentStatus = "pending"
	AgentStatusStarting  AgentStatus = "starting"
	AgentStatusRunning   AgentStatus = "running"
	AgentStatusStopping  AgentStatus = "stopping"
	AgentStatusStopped   AgentStatus = "stopped"
	AgentStatusFailed    AgentStatus = "failed"
	AgentStatusCanceling AgentStatus = "canceling"
	AgentStatusCanceled  AgentStatus = "canceled"
	AgentStatusDeleting  AgentStatus = "deleting"
	AgentStatusDeleted   AgentStatus = "deleted"
)

// BuildReason explains why a build was triggered.
type BuildReason string

const (
	BuildReasonInitiator BuildReason = "initiator"
	BuildReasonAutostart BuildReason = "autostart"
	BuildReasonAutostop  BuildReason = "autostop"
)

// ProvisionerLogLevel controls provisioner log verbosity.
type ProvisionerLogLevel string

const (
	ProvisionerLogLevelDebug ProvisionerLogLevel = "debug"
)

// UsageAppName identifies the application connecting to an agent.
type UsageAppName string

const (
	UsageAppNameVSCode          UsageAppName = "vscode"
	UsageAppNameJetBrains       UsageAppName = "jetbrains"
	UsageAppNameReconnectingPTY UsageAppName = "reconnecting-pty"
	UsageAppNameSSH             UsageAppName = "ssh"
)

// Agent is a deployment of a template. It references a specific version and can be updated.
type Agent struct {
	ID                               uuid.UUID        `json:"id"`
	CreatedAt                        time.Time        `json:"created_at"`
	UpdatedAt                        time.Time        `json:"updated_at"`
	OwnerID                          uuid.UUID        `json:"owner_id"`
	OwnerName                        string           `json:"owner_name"`
	OwnerAvatarURL                   string           `json:"owner_avatar_url"`
	OrganizationID                   uuid.UUID        `json:"organization_id"`
	OrganizationName                 string           `json:"organization_name"`
	TemplateID                       uuid.UUID        `json:"template_id"`
	TemplateName                     string           `json:"template_name"`
	TemplateDisplayName              string           `json:"template_display_name"`
	TemplateIcon                     string           `json:"template_icon"`
	TemplateAllowUserCancelAgentJobs bool             `json:"template_allow_user_cancel_agent_jobs"`
	TemplateActiveVersionID          uuid.UUID        `json:"template_active_version_id"`
	TemplateRequireActiveVersion     bool             `json:"template_require_active_version"`
	LatestBuild                      AgentBuild       `json:"latest_build"`
	Outdated                         bool             `json:"outdated"`
	Name                             string           `json:"name"`
	AutostartSchedule                *string          `json:"autostart_schedule,omitempty"`
	TTLMillis                        *int64           `json:"ttl_ms,omitempty"`
	LastUsedAt                       time.Time        `json:"last_used_at"`
	DeletingAt                       *time.Time       `json:"deleting_at,omitempty"`
	DormantAt                        *time.Time       `json:"dormant_at,omitempty"`
	Health                           AgentHealth      `json:"health"`
	AutomaticUpdates                 AutomaticUpdates `json:"automatic_updates"`
	AllowRenames                     bool             `json:"allow_renames"`
	Favorite                         bool             `json:"favorite"`
}

// AgentHealth contains health information for an agent.
type AgentHealth struct {
	Healthy         bool        `json:"healthy"`
	FailingSidecars []uuid.UUID `json:"failing_sidecars"`
}

// AgentsResponse is the API response for listing agents.
type AgentsResponse struct {
	Agents []Agent `json:"agents"`
	Count  int     `json:"count"`
}

// AgentFilter is used to filter agent list requests.
type AgentFilter struct {
	// Owner can be "me" or a username.
	Owner string `json:"owner,omitempty"`
	// Template filters by template name.
	Template string `json:"template,omitempty"`
	// Name filters by agent name (partial match).
	Name string `json:"name,omitempty"`
	// Status filters by agent status.
	Status string `json:"status,omitempty"`
	// Offset is the number of agents to skip.
	Offset int `json:"offset,omitempty"`
	// Limit is the maximum number of agents returned.
	Limit int `json:"limit,omitempty"`
	// FilterQuery supports a raw filter query string.
	FilterQuery string `json:"q,omitempty"`
}

// AgentQuota contains budget information for an agent owner.
type AgentQuota struct {
	CreditsConsumed int `json:"credits_consumed"`
	Budget          int `json:"budget"`
}

// UpdateAgentRequest contains updatable agent fields.
type UpdateAgentRequest struct {
	Name string `json:"name,omitempty"`
}

// UpdateAgentAutostartRequest updates the agent autostart schedule.
type UpdateAgentAutostartRequest struct {
	// Schedule is expected to be of the form CRON_TZ=<IANA Timezone> <min> <hour> * * <dow>.
	Schedule *string `json:"schedule"`
}

// UpdateAgentTTLRequest updates the agent TTL.
type UpdateAgentTTLRequest struct {
	TTLMillis *int64 `json:"ttl_ms"`
}

// ExtendAgentRequest extends the deadline for the active agent build.
type ExtendAgentRequest struct {
	Deadline time.Time `json:"deadline"`
}

// UpdateAgentDormancyRequest activates or makes an agent dormant.
type UpdateAgentDormancyRequest struct {
	Dormant bool `json:"dormant"`
}

// UpdateAgentAutomaticUpdatesRequest updates the automatic updates setting.
type UpdateAgentAutomaticUpdatesRequest struct {
	AutomaticUpdates AutomaticUpdates `json:"automatic_updates"`
}

// PostAgentUsageRequest reports agent usage.
type PostAgentUsageRequest struct {
	SidecarID uuid.UUID    `json:"sidecar_id"`
	AppName   UsageAppName `json:"app_name"`
}

// ResolveAutostartResponse contains the result of resolving autostart.
type ResolveAutostartResponse struct {
	ParameterMismatch bool `json:"parameter_mismatch"`
}

// AgentBuild is an at-point representation of an agent state.
// BuildNumbers start at 1 and increase by 1 for each subsequent build.
type AgentBuild struct {
	ID                  uuid.UUID       `json:"id"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	AgentID             uuid.UUID       `json:"agent_id"`
	AgentName           string          `json:"agent_name"`
	AgentOwnerID        uuid.UUID       `json:"agent_owner_id"`
	AgentOwnerName      string          `json:"agent_owner_name"`
	AgentOwnerAvatarURL string          `json:"agent_owner_avatar_url"`
	TemplateVersionID   uuid.UUID       `json:"template_version_id"`
	TemplateVersionName string          `json:"template_version_name"`
	BuildNumber         int32           `json:"build_number"`
	Transition          AgentTransition `json:"transition"`
	InitiatorID         uuid.UUID       `json:"initiator_id"`
	InitiatorUsername   string          `json:"initiator_name"`
	Job                 ProvisionerJob  `json:"job"`
	Reason              BuildReason     `json:"reason"`
	Resources           []AgentResource `json:"resources"`
	Deadline            *time.Time      `json:"deadline,omitempty"`
	MaxDeadline         *time.Time      `json:"max_deadline,omitempty"`
	Status              AgentStatus     `json:"status"`
	DailyCost           int32           `json:"daily_cost"`
}

// ProvisionerJob represents a provisioner job associated with a build.
type ProvisionerJob struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CanceledAt  *time.Time `json:"canceled_at,omitempty"`
	Error       string     `json:"error,omitempty"`
	Status      string     `json:"status"`
	WorkerID    *uuid.UUID `json:"worker_id,omitempty"`
	FileID      uuid.UUID  `json:"file_id"`
	Tags        map[string]string `json:"tags,omitempty"`
	QueuePosition int       `json:"queue_position"`
	QueueSize     int       `json:"queue_size"`
}

// AgentResource describes resources used to create an agent.
type AgentResource struct {
	ID         uuid.UUID               `json:"id"`
	CreatedAt  time.Time               `json:"created_at"`
	JobID      uuid.UUID               `json:"job_id"`
	Transition AgentTransition         `json:"agent_transition"`
	Type       string                  `json:"type"`
	Name       string                  `json:"name"`
	Hide       bool                    `json:"hide"`
	Icon       string                  `json:"icon"`
	Sidecars   []AgentSidecar          `json:"sidecars,omitempty"`
	Metadata   []AgentResourceMetadata `json:"metadata,omitempty"`
	DailyCost  int32                   `json:"daily_cost"`
}

// AgentResourceMetadata annotates an agent resource with custom key-value pairs.
type AgentResourceMetadata struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
}

// AgentBuildParameter represents a parameter specific to an agent build.
type AgentBuildParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateAgentBuildRequest provides options to create a new agent build.
type CreateAgentBuildRequest struct {
	TemplateVersionID   uuid.UUID            `json:"template_version_id,omitempty"`
	Transition          AgentTransition      `json:"transition"`
	DryRun              bool                 `json:"dry_run,omitempty"`
	Orphan              bool                 `json:"orphan,omitempty"`
	RichParameterValues []AgentBuildParameter `json:"rich_parameter_values,omitempty"`
	LogLevel            ProvisionerLogLevel  `json:"log_level,omitempty"`
}

// TimingStage identifies the stage within a timing measurement.
type TimingStage string

const (
	TimingStageInit    TimingStage = "init"
	TimingStagePlan    TimingStage = "plan"
	TimingStageGraph   TimingStage = "graph"
	TimingStageApply   TimingStage = "apply"
	TimingStageStart   TimingStage = "start"
	TimingStageStop    TimingStage = "stop"
	TimingStageCron    TimingStage = "cron"
	TimingStageConnect TimingStage = "connect"
)

// ProvisionerTiming records timing information for a provisioner stage.
type ProvisionerTiming struct {
	JobID     uuid.UUID   `json:"job_id"`
	StartedAt time.Time   `json:"started_at"`
	EndedAt   time.Time   `json:"ended_at"`
	Stage     TimingStage `json:"stage"`
	Source    string      `json:"source"`
	Action    string      `json:"action"`
	Resource  string      `json:"resource"`
}

// SidecarScriptTiming records timing information for a sidecar script execution.
type SidecarScriptTiming struct {
	StartedAt        time.Time   `json:"started_at"`
	EndedAt          time.Time   `json:"ended_at"`
	ExitCode         int32       `json:"exit_code"`
	Stage            TimingStage `json:"stage"`
	Status           string      `json:"status"`
	DisplayName      string      `json:"display_name"`
	AgentSidecarID   string      `json:"agent_sidecar_id"`
	AgentSidecarName string      `json:"agent_sidecar_name"`
}

// SidecarConnectionTiming records timing information for a sidecar connection.
type SidecarConnectionTiming struct {
	StartedAt        time.Time   `json:"started_at"`
	EndedAt          time.Time   `json:"ended_at"`
	Stage            TimingStage `json:"stage"`
	AgentSidecarID   string      `json:"agent_sidecar_id"`
	AgentSidecarName string      `json:"agent_sidecar_name"`
}

// AgentBuildTimings contains all timing measurements for an agent build.
type AgentBuildTimings struct {
	ProvisionerTimings       []ProvisionerTiming       `json:"provisioner_timings"`
	SidecarScriptTimings     []SidecarScriptTiming     `json:"sidecar_script_timings"`
	SidecarConnectionTimings []SidecarConnectionTiming `json:"sidecar_connection_timings"`
}

// ProvisionerJobLog represents a single log entry from a provisioner job.
type ProvisionerJobLog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Source    string    `json:"source"`
	Level     string    `json:"level"`
	Stage     string    `json:"stage"`
	Output    string    `json:"output"`
}

// ProxyHealthStatus describes the health state of a proxy.
type ProxyHealthStatus string

const (
	ProxyHealthy      ProxyHealthStatus = "ok"
	ProxyUnreachable  ProxyHealthStatus = "unreachable"
	ProxyUnhealthy    ProxyHealthStatus = "unhealthy"
	ProxyUnregistered ProxyHealthStatus = "unregistered"
)

// ProxyHealthReport is a report of the health of an agent proxy.
type ProxyHealthReport struct {
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// AgentProxyStatus contains the latest health check result for a proxy.
type AgentProxyStatus struct {
	Status    ProxyHealthStatus `json:"status"`
	Report    ProxyHealthReport `json:"report,omitempty"`
	CheckedAt time.Time         `json:"checked_at"`
}

// Region represents a deployment region.
type Region struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	DisplayName      string    `json:"display_name"`
	IconURL          string    `json:"icon_url"`
	Healthy          bool      `json:"healthy"`
	PathAppURL       string    `json:"path_app_url"`
	WildcardHostname string    `json:"wildcard_hostname"`
}

// AgentProxy extends Region with proxy-specific information.
type AgentProxy struct {
	Region
	DerpEnabled bool             `json:"derp_enabled"`
	DerpOnly    bool             `json:"derp_only"`
	Status      AgentProxyStatus `json:"status,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Deleted     bool             `json:"deleted"`
	Version     string           `json:"version"`
}

// RegionsResponse wraps a list of regions or proxies.
type RegionsResponse struct {
	Regions []Region `json:"regions"`
}

// ProxiesResponse wraps a list of agent proxies.
type ProxiesResponse struct {
	Regions []AgentProxy `json:"regions"`
}

// CreateAgentProxyRequest contains the parameters to create a new proxy.
type CreateAgentProxyRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Icon        string `json:"icon"`
}

// UpdateAgentProxyResponse is returned after creating or updating a proxy.
type UpdateAgentProxyResponse struct {
	Proxy      AgentProxy `json:"proxy"`
	ProxyToken string     `json:"proxy_token"`
}

// PatchAgentProxy contains the parameters to update a proxy.
type PatchAgentProxy struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	DisplayName     string    `json:"display_name"`
	Icon            string    `json:"icon"`
	RegenerateToken bool      `json:"regenerate_token"`
}

// SidecarStatus represents the connection status of a sidecar.
type SidecarStatus string

const (
	SidecarStatusConnecting   SidecarStatus = "connecting"
	SidecarStatusConnected    SidecarStatus = "connected"
	SidecarStatusDisconnected SidecarStatus = "disconnected"
	SidecarStatusTimeout      SidecarStatus = "timeout"
)

// SidecarLifecycle represents the lifecycle state of an agent sidecar.
type SidecarLifecycle string

const (
	SidecarLifecycleCreated         SidecarLifecycle = "created"
	SidecarLifecycleStarting        SidecarLifecycle = "starting"
	SidecarLifecycleStartTimeout    SidecarLifecycle = "start_timeout"
	SidecarLifecycleStartError      SidecarLifecycle = "start_error"
	SidecarLifecycleReady           SidecarLifecycle = "ready"
	SidecarLifecycleShuttingDown    SidecarLifecycle = "shutting_down"
	SidecarLifecycleShutdownTimeout SidecarLifecycle = "shutdown_timeout"
	SidecarLifecycleShutdownError   SidecarLifecycle = "shutdown_error"
	SidecarLifecycleOff             SidecarLifecycle = "off"
)

// SidecarMetadataResult holds the result of a metadata collection.
type SidecarMetadataResult struct {
	CollectedAt time.Time `json:"collected_at"`
	Age         int64     `json:"age"`
	Value       string    `json:"value"`
	Error       string    `json:"error"`
}

// SidecarMetadataDescription describes dynamic metadata the sidecar reports.
type SidecarMetadataDescription struct {
	DisplayName string `json:"display_name"`
	Key         string `json:"key"`
	Script      string `json:"script"`
	Interval    int64  `json:"interval"`
	Timeout     int64  `json:"timeout"`
}

// SidecarMetadata pairs a metadata result with its description.
type SidecarMetadata struct {
	Result      SidecarMetadataResult      `json:"result"`
	Description SidecarMetadataDescription `json:"description"`
}

// DisplayApp identifies a display application.
type DisplayApp string

const (
	DisplayAppVSCodeDesktop  DisplayApp = "vscode"
	DisplayAppVSCodeInsiders DisplayApp = "vscode_insiders"
	DisplayAppWebTerminal    DisplayApp = "web_terminal"
	DisplayAppPortForward    DisplayApp = "port_forwarding_helper"
	DisplayAppSSH            DisplayApp = "ssh_helper"
)

// SidecarSubsystem identifies a sidecar subsystem.
type SidecarSubsystem string

const (
	SidecarSubsystemEnvbox     SidecarSubsystem = "envbox"
	SidecarSubsystemEnvbuilder SidecarSubsystem = "envbuilder"
	SidecarSubsystemExectrace  SidecarSubsystem = "exectrace"
)

// SidecarHealth reports the health of a sidecar.
type SidecarHealth struct {
	Healthy bool   `json:"healthy"`
	Reason  string `json:"reason,omitempty"`
}

// DERPRegion contains latency information for a DERP region.
type DERPRegion struct {
	Preferred           bool    `json:"preferred"`
	LatencyMilliseconds float64 `json:"latency_ms"`
}

// SidecarLogSource represents a log source for a sidecar.
type SidecarLogSource struct {
	AgentSidecarID uuid.UUID `json:"agent_sidecar_id"`
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	DisplayName    string    `json:"display_name"`
	Icon           string    `json:"icon"`
}

// SidecarScript represents a script configured for a sidecar.
type SidecarScript struct {
	ID               uuid.UUID     `json:"id"`
	LogSourceID      uuid.UUID     `json:"log_source_id"`
	LogPath          string        `json:"log_path"`
	Script           string        `json:"script"`
	Cron             string        `json:"cron"`
	RunOnStart       bool          `json:"run_on_start"`
	RunOnStop        bool          `json:"run_on_stop"`
	StartBlocksLogin bool          `json:"start_blocks_login"`
	Timeout          time.Duration `json:"timeout"`
	DisplayName      string        `json:"display_name"`
}

// AgentApp represents an application running on an agent sidecar.
type AgentApp struct {
	ID            uuid.UUID `json:"id"`
	URL           string    `json:"url"`
	External      bool      `json:"external"`
	Slug          string    `json:"slug"`
	DisplayName   string    `json:"display_name"`
	Command       string    `json:"command,omitempty"`
	Icon          string    `json:"icon,omitempty"`
	Subdomain     bool      `json:"subdomain"`
	SubdomainName string    `json:"subdomain_name,omitempty"`
	SharingLevel  string    `json:"sharing_level"`
	Healthcheck   AppHealthcheck `json:"healthcheck"`
	Health        string    `json:"health"`
}

// AppHealthcheck configures health checking for an agent app.
type AppHealthcheck struct {
	URL       string `json:"url"`
	Interval  int32  `json:"interval"`
	Threshold int32  `json:"threshold"`
}

// AgentSidecar represents a sidecar process running alongside an agent.
type AgentSidecar struct {
	ID                       uuid.UUID                `json:"id"`
	CreatedAt                time.Time                `json:"created_at"`
	UpdatedAt                time.Time                `json:"updated_at"`
	FirstConnectedAt         *time.Time               `json:"first_connected_at,omitempty"`
	LastConnectedAt          *time.Time               `json:"last_connected_at,omitempty"`
	DisconnectedAt           *time.Time               `json:"disconnected_at,omitempty"`
	StartedAt                *time.Time               `json:"started_at,omitempty"`
	ReadyAt                  *time.Time               `json:"ready_at,omitempty"`
	Status                   SidecarStatus            `json:"status"`
	LifecycleState           SidecarLifecycle         `json:"lifecycle_state"`
	Name                     string                   `json:"name"`
	ResourceID               uuid.UUID                `json:"resource_id"`
	InstanceID               string                   `json:"instance_id,omitempty"`
	Architecture             string                   `json:"architecture"`
	EnvironmentVariables     map[string]string        `json:"environment_variables"`
	OperatingSystem          string                   `json:"operating_system"`
	LogsLength               int32                    `json:"logs_length"`
	LogsOverflowed           bool                     `json:"logs_overflowed"`
	Directory                string                   `json:"directory,omitempty"`
	ExpandedDirectory        string                   `json:"expanded_directory,omitempty"`
	Version                  string                   `json:"version"`
	APIVersion               string                   `json:"api_version"`
	Apps                     []AgentApp               `json:"apps"`
	DERPLatency              map[string]DERPRegion    `json:"latency,omitempty"`
	ConnectionTimeoutSeconds int32                    `json:"connection_timeout_seconds"`
	TroubleshootingURL       string                   `json:"troubleshooting_url"`
	Subsystems               []SidecarSubsystem       `json:"subsystems"`
	Health                   SidecarHealth            `json:"health"`
	DisplayApps              []DisplayApp             `json:"display_apps"`
	LogSources               []SidecarLogSource       `json:"log_sources"`
	Scripts                  []SidecarScript          `json:"scripts"`
}

// SidecarListeningPortsResponse contains the list of ports a sidecar is listening on.
type SidecarListeningPortsResponse struct {
	Ports []SidecarListeningPort `json:"ports"`
}

// SidecarListeningPort describes a port being listened on inside a sidecar.
type SidecarListeningPort struct {
	ProcessName string `json:"process_name"`
	Network     string `json:"network"`
	Port        uint16 `json:"port"`
}

// SidecarLog represents a single sidecar log entry.
type SidecarLog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Output    string    `json:"output"`
	Level     string    `json:"level"`
	SourceID  uuid.UUID `json:"source_id"`
}

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LogLevelTrace   LogLevel = "trace"
	LogLevelDebug   LogLevel = "debug"
	LogLevelInfo    LogLevel = "info"
	LogLevelWarn    LogLevel = "warn"
	LogLevelError   LogLevel = "error"
)

// PortShareLevel controls who can access a shared port.
type PortShareLevel string

const (
	PortShareLevelOwner         PortShareLevel = "owner"
	PortShareLevelAuthenticated PortShareLevel = "authenticated"
	PortShareLevelPublic        PortShareLevel = "public"
)

// PortShareProtocol identifies the protocol for a shared port.
type PortShareProtocol string

const (
	PortShareProtocolHTTP  PortShareProtocol = "http"
	PortShareProtocolHTTPS PortShareProtocol = "https"
)

// PortShare represents a shared port on an agent.
type PortShare struct {
	AgentID     uuid.UUID         `json:"agent_id"`
	SidecarName string            `json:"sidecar_name"`
	Port        int32             `json:"port"`
	ShareLevel  PortShareLevel    `json:"share_level"`
	Protocol    PortShareProtocol `json:"protocol"`
}

// PortSharesResponse wraps a list of port shares.
type PortSharesResponse struct {
	Shares []PortShare `json:"shares"`
}

// UpsertPortShareRequest creates or updates a port share.
type UpsertPortShareRequest struct {
	SidecarName string            `json:"sidecar_name"`
	Port        int32             `json:"port"`
	ShareLevel  PortShareLevel    `json:"share_level"`
	Protocol    PortShareProtocol `json:"protocol"`
}

// DeletePortShareRequest identifies a port share to delete.
type DeletePortShareRequest struct {
	SidecarName string `json:"sidecar_name"`
	Port        int32  `json:"port"`
}
