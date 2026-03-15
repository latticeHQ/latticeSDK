// Package sessions provides session management services for the Lattice Runtime API.
// It covers sessions, session builds, sidecars, and real-time sessions.
package sessions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// AutomaticUpdates represents the automatic updates policy for a session.
type AutomaticUpdates string

const (
	// AutomaticUpdatesAlways enables automatic template updates.
	AutomaticUpdatesAlways AutomaticUpdates = "always"
	// AutomaticUpdatesNever disables automatic template updates.
	AutomaticUpdatesNever AutomaticUpdates = "never"
)

// SessionTransition represents the desired state change for a session build.
type SessionTransition string

const (
	// SessionTransitionStart starts the session.
	SessionTransitionStart SessionTransition = "start"
	// SessionTransitionStop stops the session.
	SessionTransitionStop SessionTransition = "stop"
	// SessionTransitionDelete deletes the session.
	SessionTransitionDelete SessionTransition = "delete"
)

// SessionStatus represents the status of a session, derived from the latest build.
type SessionStatus string

const (
	SessionStatusPending   SessionStatus = "pending"
	SessionStatusStarting  SessionStatus = "starting"
	SessionStatusRunning   SessionStatus = "running"
	SessionStatusStopping  SessionStatus = "stopping"
	SessionStatusStopped   SessionStatus = "stopped"
	SessionStatusFailed    SessionStatus = "failed"
	SessionStatusCanceling SessionStatus = "canceling"
	SessionStatusCanceled  SessionStatus = "canceled"
	SessionStatusDeleting  SessionStatus = "deleting"
	SessionStatusDeleted   SessionStatus = "deleted"
)

// BuildReason indicates what triggered a session build.
type BuildReason string

const (
	// BuildReasonInitiator indicates a user-triggered build.
	BuildReasonInitiator BuildReason = "initiator"
	// BuildReasonAutostart indicates an autostart-triggered build.
	BuildReasonAutostart BuildReason = "autostart"
	// BuildReasonAutostop indicates an autostop-triggered build.
	BuildReasonAutostop BuildReason = "autostop"
)

// ProvisionerLogLevel controls the logging verbosity of a provisioner.
type ProvisionerLogLevel string

const (
	// ProvisionerLogLevelDebug enables debug-level logging.
	ProvisionerLogLevelDebug ProvisionerLogLevel = "debug"
)

// LogLevel represents the severity level of a log entry.
type LogLevel string

const (
	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// LogSource identifies the origin of a provisioner log entry.
type LogSource string

const (
	LogSourceProvisionerDaemon LogSource = "provisioner_daemon"
	LogSourceProvisioner       LogSource = "provisioner"
)

// ProvisionerJobStatus represents the status of a provisioner job.
type ProvisionerJobStatus string

const (
	ProvisionerJobPending   ProvisionerJobStatus = "pending"
	ProvisionerJobRunning   ProvisionerJobStatus = "running"
	ProvisionerJobSucceeded ProvisionerJobStatus = "succeeded"
	ProvisionerJobCanceling ProvisionerJobStatus = "canceling"
	ProvisionerJobCanceled  ProvisionerJobStatus = "canceled"
	ProvisionerJobFailed    ProvisionerJobStatus = "failed"
)

// Active returns whether the provisioner job is still active.
func (p ProvisionerJobStatus) Active() bool {
	return p == ProvisionerJobPending ||
		p == ProvisionerJobRunning ||
		p == ProvisionerJobCanceling
}

// JobErrorCode represents an error code returned by a provisioner job.
type JobErrorCode string

const (
	// RequiredTemplateVariables indicates missing required template variables.
	RequiredTemplateVariables JobErrorCode = "REQUIRED_TEMPLATE_VARIABLES"
)

// ProvisionerJob represents the state of a provisioner job.
type ProvisionerJob struct {
	ID            uuid.UUID            `json:"id" format:"uuid"`
	CreatedAt     time.Time            `json:"created_at" format:"date-time"`
	StartedAt     *time.Time           `json:"started_at,omitempty" format:"date-time"`
	CompletedAt   *time.Time           `json:"completed_at,omitempty" format:"date-time"`
	CanceledAt    *time.Time           `json:"canceled_at,omitempty" format:"date-time"`
	Error         string               `json:"error,omitempty"`
	ErrorCode     JobErrorCode         `json:"error_code,omitempty"`
	Status        ProvisionerJobStatus `json:"status"`
	WorkerID      *uuid.UUID           `json:"worker_id,omitempty" format:"uuid"`
	FileID        uuid.UUID            `json:"file_id" format:"uuid"`
	Tags          map[string]string    `json:"tags"`
	QueuePosition int                  `json:"queue_position"`
	QueueSize     int                  `json:"queue_size"`
}

// ProvisionerJobLog is a single log entry from a provisioner job.
type ProvisionerJobLog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	Source    LogSource `json:"log_source"`
	Level     LogLevel  `json:"log_level"`
	Stage     string    `json:"stage"`
	Output    string    `json:"output"`
}

// UsageAppName identifies an application type for session usage tracking.
type UsageAppName string

const (
	UsageAppNameVscode          UsageAppName = "vscode"
	UsageAppNameJetbrains       UsageAppName = "jetbrains"
	UsageAppNameReconnectingPty UsageAppName = "reconnecting-pty"
	UsageAppNameSSH             UsageAppName = "ssh"
)

// SessionHealth reports the health of a session.
type SessionHealth struct {
	Healthy         bool        `json:"healthy"`
	FailingSidecars []uuid.UUID `json:"failing_sidecars" format:"uuid"`
}

// Session is a deployment of a template. It references a specific
// version and can be updated.
type Session struct {
	ID                                 uuid.UUID        `json:"id" format:"uuid"`
	CreatedAt                          time.Time        `json:"created_at" format:"date-time"`
	UpdatedAt                          time.Time        `json:"updated_at" format:"date-time"`
	OwnerID                            uuid.UUID        `json:"owner_id" format:"uuid"`
	OwnerName                          string           `json:"owner_name"`
	OwnerAvatarURL                     string           `json:"owner_avatar_url"`
	OrganizationID                     uuid.UUID        `json:"organization_id" format:"uuid"`
	OrganizationName                   string           `json:"organization_name"`
	TemplateID                         uuid.UUID        `json:"template_id" format:"uuid"`
	TemplateName                       string           `json:"template_name"`
	TemplateDisplayName                string           `json:"template_display_name"`
	TemplateIcon                       string           `json:"template_icon"`
	TemplateAllowUserCancelSessionJobs bool             `json:"template_allow_user_cancel_session_jobs"`
	TemplateActiveVersionID            uuid.UUID        `json:"template_active_version_id" format:"uuid"`
	TemplateRequireActiveVersion       bool             `json:"template_require_active_version"`
	LatestBuild                        SessionBuild     `json:"latest_build"`
	Outdated                           bool             `json:"outdated"`
	Name                               string           `json:"name"`
	AutostartSchedule                  *string          `json:"autostart_schedule,omitempty"`
	TTLMillis                          *int64           `json:"ttl_ms,omitempty"`
	LastUsedAt                         time.Time        `json:"last_used_at" format:"date-time"`
	DeletingAt                         *time.Time       `json:"deleting_at" format:"date-time"`
	DormantAt                          *time.Time       `json:"dormant_at" format:"date-time"`
	Health                             SessionHealth    `json:"health"`
	AutomaticUpdates                   AutomaticUpdates `json:"automatic_updates"`
	AllowRenames                       bool             `json:"allow_renames"`
	Favorite                           bool             `json:"favorite"`
}

// FullName returns the session's fully qualified name as "owner/name".
func (s Session) FullName() string {
	return fmt.Sprintf("%s/%s", s.OwnerName, s.Name)
}

// SessionsResponse is the paginated response for listing sessions.
type SessionsResponse struct {
	Sessions []Session `json:"sessions"`
	Count    int       `json:"count"`
}

// SessionFilter provides filtering options for listing sessions.
type SessionFilter struct {
	// Owner can be "me" or a username.
	Owner string `json:"owner,omitempty"`
	// Template filters by template name.
	Template string `json:"template,omitempty"`
	// Name filters sessions by partial name match.
	Name string `json:"name,omitempty"`
	// Status filters by session status (derived from latest build status).
	Status string `json:"status,omitempty"`
	// Offset is the number of sessions to skip.
	Offset int `json:"offset,omitempty"`
	// Limit is the maximum number of sessions to return.
	Limit int `json:"limit,omitempty"`
	// FilterQuery supports a raw filter query string.
	FilterQuery string `json:"q,omitempty"`
}

// asRequestOption returns a RequestOption that applies session filter query parameters.
func (f SessionFilter) asRequestOption() client.RequestOption {
	return func(r *http.Request) {
		var params []string
		if f.Owner != "" {
			params = append(params, fmt.Sprintf("owner:%q", f.Owner))
		}
		if f.Name != "" {
			params = append(params, fmt.Sprintf("name:%q", f.Name))
		}
		if f.Template != "" {
			params = append(params, fmt.Sprintf("template:%q", f.Template))
		}
		if f.Status != "" {
			params = append(params, fmt.Sprintf("status:%q", f.Status))
		}
		if f.FilterQuery != "" {
			params = append(params, f.FilterQuery)
		}
		q := r.URL.Query()
		q.Set("q", strings.Join(params, " "))
		r.URL.RawQuery = q.Encode()
	}
}

// SessionOptions provides additional options when retrieving a session.
type SessionOptions struct {
	// IncludeDeleted includes soft-deleted sessions in the response.
	IncludeDeleted bool `json:"include_deleted,omitempty"`
}

// asRequestOption returns a RequestOption that applies session options as query parameters.
func (o SessionOptions) asRequestOption() client.RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		if o.IncludeDeleted {
			q.Set("include_deleted", "true")
		}
		r.URL.RawQuery = q.Encode()
	}
}

// CreateSessionBuildRequest provides options to create a new session build.
type CreateSessionBuildRequest struct {
	// TemplateVersionID optionally pins the build to a specific template version.
	TemplateVersionID uuid.UUID `json:"template_version_id,omitempty" format:"uuid"`
	// Transition is the desired state change (start, stop, or delete).
	Transition SessionTransition `json:"transition" validate:"oneof=start stop delete,required"`
	// DryRun, when true, validates the build without executing it.
	DryRun bool `json:"dry_run,omitempty"`
	// ProvisionerState is optional raw provisioner state.
	ProvisionerState []byte `json:"state,omitempty"`
	// Orphan may be set for the destroy transition to leave resources intact.
	Orphan bool `json:"orphan,omitempty"`
	// RichParameterValues are optional build parameters that overwrite existing ones.
	RichParameterValues []SessionBuildParameter `json:"rich_parameter_values,omitempty"`
	// LogLevel changes the default logging verbosity of a provider ("info" if empty).
	LogLevel ProvisionerLogLevel `json:"log_level,omitempty" validate:"omitempty,oneof=debug"`
}

// UpdateSessionRequest provides options for renaming a session.
type UpdateSessionRequest struct {
	Name string `json:"name,omitempty" validate:"username"`
}

// UpdateSessionAutostartRequest is a request to update a session's autostart schedule.
type UpdateSessionAutostartRequest struct {
	// Schedule is expected to be of the form "CRON_TZ=<IANA Timezone> <min> <hour> * * <dow>".
	// Example: "CRON_TZ=US/Central 30 9 * * 1-5".
	// Defaults to UTC if CRON_TZ is not present. Set to nil to disable autostart.
	Schedule *string `json:"schedule"`
}

// UpdateSessionTTLRequest is a request to update a session's TTL (time-to-live).
type UpdateSessionTTLRequest struct {
	// TTLMillis is the time-to-live in milliseconds. Nil disables autostop.
	TTLMillis *int64 `json:"ttl_ms"`
}

// PutExtendSessionRequest is a request to extend the deadline of the active session build.
type PutExtendSessionRequest struct {
	Deadline time.Time `json:"deadline" validate:"required" format:"date-time"`
}

// UpdateSessionDormancy is a request to activate or make a session dormant.
type UpdateSessionDormancy struct {
	// Dormant, when true, marks the session as dormant. When false, activates a dormant session.
	Dormant bool `json:"dormant"`
}

// UpdateSessionAutomaticUpdatesRequest is a request to update a session's automatic updates policy.
type UpdateSessionAutomaticUpdatesRequest struct {
	AutomaticUpdates AutomaticUpdates `json:"automatic_updates"`
}

// PostSessionUsageRequest marks a session as recently used and records app usage.
type PostSessionUsageRequest struct {
	SidecarID uuid.UUID    `json:"sidecar_id" format:"uuid"`
	AppName   UsageAppName `json:"app_name"`
}

// SessionQuota represents the session quota for a user within an organization.
type SessionQuota struct {
	CreditsConsumed int `json:"credits_consumed"`
	Budget          int `json:"budget"`
}

// ResolveAutostartResponse contains the result of resolving autostart compatibility.
type ResolveAutostartResponse struct {
	ParameterMismatch bool `json:"parameter_mismatch"`
}

// SessionBuild is a point-in-time representation of a session state.
// BuildNumbers start at 1 and increase by 1 for each subsequent build.
type SessionBuild struct {
	ID                    uuid.UUID         `json:"id" format:"uuid"`
	CreatedAt             time.Time         `json:"created_at" format:"date-time"`
	UpdatedAt             time.Time         `json:"updated_at" format:"date-time"`
	SessionID             uuid.UUID         `json:"session_id" format:"uuid"`
	SessionName           string            `json:"session_name"`
	SessionOwnerID        uuid.UUID         `json:"session_owner_id" format:"uuid"`
	SessionOwnerName      string            `json:"session_owner_name"`
	SessionOwnerAvatarURL string            `json:"session_owner_avatar_url"`
	TemplateVersionID     uuid.UUID         `json:"template_version_id" format:"uuid"`
	TemplateVersionName   string            `json:"template_version_name"`
	BuildNumber           int32             `json:"build_number"`
	Transition            SessionTransition `json:"transition"`
	InitiatorID           uuid.UUID         `json:"initiator_id" format:"uuid"`
	InitiatorUsername     string            `json:"initiator_name"`
	Job                   ProvisionerJob    `json:"job"`
	Reason                BuildReason       `json:"reason"`
	Resources             []SessionResource `json:"resources"`
	Deadline              *time.Time        `json:"deadline,omitempty" format:"date-time"`
	MaxDeadline           *time.Time        `json:"max_deadline,omitempty" format:"date-time"`
	Status                SessionStatus     `json:"status"`
	DailyCost             int32             `json:"daily_cost"`
}

// SessionResource describes resources used to create a session.
type SessionResource struct {
	ID         uuid.UUID                 `json:"id" format:"uuid"`
	CreatedAt  time.Time                 `json:"created_at" format:"date-time"`
	JobID      uuid.UUID                 `json:"job_id" format:"uuid"`
	Transition SessionTransition         `json:"session_transition"`
	Type       string                    `json:"type"`
	Name       string                    `json:"name"`
	Hide       bool                      `json:"hide"`
	Icon       string                    `json:"icon"`
	Sidecars   []SessionSidecar          `json:"sidecars,omitempty"`
	Metadata   []SessionResourceMetadata `json:"metadata,omitempty"`
	DailyCost  int32                     `json:"daily_cost"`
}

// SessionResourceMetadata annotates a session resource with custom key-value pairs.
type SessionResourceMetadata struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
}

// SessionBuildParameter represents a parameter specific to a session build.
type SessionBuildParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TimingStage identifies the stage of a provisioner or sidecar timing event.
type TimingStage string

const (
	TimingStageInit  TimingStage = "init"
	TimingStagePlan  TimingStage = "plan"
	TimingStageGraph TimingStage = "graph"
	TimingStageApply TimingStage = "apply"
	TimingStageStart TimingStage = "start"
	TimingStageStop  TimingStage = "stop"
)

// ProvisionerTiming records the duration of a provisioner stage.
type ProvisionerTiming struct {
	JobID     uuid.UUID   `json:"job_id" format:"uuid"`
	StartedAt time.Time   `json:"started_at" format:"date-time"`
	EndedAt   time.Time   `json:"ended_at" format:"date-time"`
	Stage     TimingStage `json:"stage"`
	Source    string      `json:"source"`
	Action    string      `json:"action"`
	Resource  string      `json:"resource"`
}

// SidecarScriptTiming records the duration of a sidecar script execution.
type SidecarScriptTiming struct {
	StartedAt              time.Time   `json:"started_at" format:"date-time"`
	EndedAt                time.Time   `json:"ended_at" format:"date-time"`
	ExitCode               int32       `json:"exit_code"`
	Stage                  TimingStage `json:"stage"`
	Status                 string      `json:"status"`
	DisplayName            string      `json:"display_name"`
	SessionSidecarID       string      `json:"agent_sidecar_id"`
	SessionSidecarName     string      `json:"agent_sidecar_name"`
}

// SidecarConnectionTiming records the time taken for a sidecar to connect.
type SidecarConnectionTiming struct {
	StartedAt          time.Time   `json:"started_at" format:"date-time"`
	EndedAt            time.Time   `json:"ended_at" format:"date-time"`
	Stage              TimingStage `json:"stage"`
	SessionSidecarID   string      `json:"agent_sidecar_id"`
	SessionSidecarName string      `json:"agent_sidecar_name"`
}

// SessionBuildTimings contains all timing information for a session build.
type SessionBuildTimings struct {
	ProvisionerTimings       []ProvisionerTiming       `json:"provisioner_timings"`
	SidecarScriptTimings     []SidecarScriptTiming     `json:"sidecar_script_timings"`
	SidecarConnectionTimings []SidecarConnectionTiming `json:"sidecar_connection_timings"`
}

// SessionSidecarStatus represents the connection status of a session sidecar.
type SessionSidecarStatus string

const (
	SessionSidecarConnecting   SessionSidecarStatus = "connecting"
	SessionSidecarConnected    SessionSidecarStatus = "connected"
	SessionSidecarDisconnected SessionSidecarStatus = "disconnected"
	SessionSidecarTimeout      SessionSidecarStatus = "timeout"
)

// SessionSidecarLifecycle represents the lifecycle state of a session sidecar.
type SessionSidecarLifecycle string

const (
	SessionSidecarLifecycleCreated         SessionSidecarLifecycle = "created"
	SessionSidecarLifecycleStarting        SessionSidecarLifecycle = "starting"
	SessionSidecarLifecycleStartTimeout    SessionSidecarLifecycle = "start_timeout"
	SessionSidecarLifecycleStartError      SessionSidecarLifecycle = "start_error"
	SessionSidecarLifecycleReady           SessionSidecarLifecycle = "ready"
	SessionSidecarLifecycleShuttingDown    SessionSidecarLifecycle = "shutting_down"
	SessionSidecarLifecycleShutdownTimeout SessionSidecarLifecycle = "shutdown_timeout"
	SessionSidecarLifecycleShutdownError   SessionSidecarLifecycle = "shutdown_error"
	SessionSidecarLifecycleOff             SessionSidecarLifecycle = "off"
)

// Starting returns true if the sidecar is in the process of starting.
func (l SessionSidecarLifecycle) Starting() bool {
	switch l {
	case SessionSidecarLifecycleCreated, SessionSidecarLifecycleStarting:
		return true
	default:
		return false
	}
}

// ShuttingDown returns true if the sidecar is shutting down or has shut down.
func (l SessionSidecarLifecycle) ShuttingDown() bool {
	switch l {
	case SessionSidecarLifecycleShuttingDown, SessionSidecarLifecycleShutdownTimeout, SessionSidecarLifecycleShutdownError, SessionSidecarLifecycleOff:
		return true
	default:
		return false
	}
}

// SessionSidecarHealth reports the health of a sidecar.
type SessionSidecarHealth struct {
	Healthy bool   `json:"healthy"`
	Reason  string `json:"reason,omitempty"`
}

// SessionSidecarMetadataResult contains the result of a metadata collection.
type SessionSidecarMetadataResult struct {
	CollectedAt time.Time `json:"collected_at" format:"date-time"`
	Age         int64     `json:"age"`
	Value       string    `json:"value"`
	Error       string    `json:"error"`
}

// SessionSidecarMetadataDescription describes dynamic metadata the sidecar reports.
type SessionSidecarMetadataDescription struct {
	DisplayName string `json:"display_name"`
	Key         string `json:"key"`
	Script      string `json:"script"`
	Interval    int64  `json:"interval"`
	Timeout     int64  `json:"timeout"`
}

// SessionSidecarMetadata pairs a metadata result with its description.
type SessionSidecarMetadata struct {
	Result      SessionSidecarMetadataResult      `json:"result"`
	Description SessionSidecarMetadataDescription `json:"description"`
}

// SessionSidecarLogSource identifies the origin of a sidecar log.
type SessionSidecarLogSource struct {
	SessionSidecarID uuid.UUID `json:"session_sidecar_id" format:"uuid"`
	ID               uuid.UUID `json:"id" format:"uuid"`
	CreatedAt        time.Time `json:"created_at" format:"date-time"`
	DisplayName      string    `json:"display_name"`
	Icon             string    `json:"icon"`
}

// SessionSidecarScript represents a script configured on a sidecar.
type SessionSidecarScript struct {
	ID               uuid.UUID     `json:"id" format:"uuid"`
	LogSourceID      uuid.UUID     `json:"log_source_id" format:"uuid"`
	LogPath          string        `json:"log_path"`
	Script           string        `json:"script"`
	Cron             string        `json:"cron"`
	RunOnStart       bool          `json:"run_on_start"`
	RunOnStop        bool          `json:"run_on_stop"`
	StartBlocksLogin bool          `json:"start_blocks_login"`
	Timeout          time.Duration `json:"timeout"`
	DisplayName      string        `json:"display_name"`
}

// DisplayApp identifies a display application for a sidecar.
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

// DERPRegion contains latency information for a DERP relay region.
type DERPRegion struct {
	Preferred           bool    `json:"preferred"`
	LatencyMilliseconds float64 `json:"latency_ms"`
}

// SessionAppSharingLevel defines the sharing level for a session app.
type SessionAppSharingLevel string

const (
	SessionAppSharingLevelOwner         SessionAppSharingLevel = "owner"
	SessionAppSharingLevelAuthenticated SessionAppSharingLevel = "authenticated"
	SessionAppSharingLevelPublic        SessionAppSharingLevel = "public"
)

// SessionAppHealth represents the health status of a session app.
type SessionAppHealth string

const (
	SessionAppHealthDisabled     SessionAppHealth = "disabled"
	SessionAppHealthInitializing SessionAppHealth = "initializing"
	SessionAppHealthHealthy      SessionAppHealth = "healthy"
	SessionAppHealthUnhealthy    SessionAppHealth = "unhealthy"
)

// Healthcheck defines the configuration for checking app health.
type Healthcheck struct {
	URL       string `json:"url"`
	Interval  int32  `json:"interval"`
	Threshold int32  `json:"threshold"`
}

// SessionApp represents an application exposed by a session sidecar.
type SessionApp struct {
	ID            uuid.UUID              `json:"id" format:"uuid"`
	URL           string                 `json:"url"`
	External      bool                   `json:"external"`
	Slug          string                 `json:"slug"`
	DisplayName   string                 `json:"display_name"`
	Command       string                 `json:"command,omitempty"`
	Icon          string                 `json:"icon,omitempty"`
	Subdomain     bool                   `json:"subdomain"`
	SubdomainName string                 `json:"subdomain_name,omitempty"`
	SharingLevel  SessionAppSharingLevel `json:"sharing_level"`
	Healthcheck   Healthcheck            `json:"healthcheck"`
	Health        SessionAppHealth       `json:"health"`
	Hidden        bool                   `json:"hidden"`
}

// SessionSidecar represents a sidecar process running within a session.
type SessionSidecar struct {
	ID                       uuid.UUID                   `json:"id" format:"uuid"`
	CreatedAt                time.Time                   `json:"created_at" format:"date-time"`
	UpdatedAt                time.Time                   `json:"updated_at" format:"date-time"`
	FirstConnectedAt         *time.Time                  `json:"first_connected_at,omitempty" format:"date-time"`
	LastConnectedAt          *time.Time                  `json:"last_connected_at,omitempty" format:"date-time"`
	DisconnectedAt           *time.Time                  `json:"disconnected_at,omitempty" format:"date-time"`
	StartedAt                *time.Time                  `json:"started_at,omitempty" format:"date-time"`
	ReadyAt                  *time.Time                  `json:"ready_at,omitempty" format:"date-time"`
	Status                   SessionSidecarStatus        `json:"status"`
	LifecycleState           SessionSidecarLifecycle     `json:"lifecycle_state"`
	Name                     string                      `json:"name"`
	ResourceID               uuid.UUID                   `json:"resource_id" format:"uuid"`
	InstanceID               string                      `json:"instance_id,omitempty"`
	Architecture             string                      `json:"architecture"`
	EnvironmentVariables     map[string]string           `json:"environment_variables"`
	OperatingSystem          string                      `json:"operating_system"`
	LogsLength               int32                       `json:"logs_length"`
	LogsOverflowed           bool                        `json:"logs_overflowed"`
	Directory                string                      `json:"directory,omitempty"`
	ExpandedDirectory        string                      `json:"expanded_directory,omitempty"`
	Version                  string                      `json:"version"`
	APIVersion               string                      `json:"api_version"`
	Apps                     []SessionApp                `json:"apps"`
	DERPLatency              map[string]DERPRegion       `json:"latency,omitempty"`
	ConnectionTimeoutSeconds int32                       `json:"connection_timeout_seconds"`
	TroubleshootingURL       string                      `json:"troubleshooting_url"`
	Subsystems               []SidecarSubsystem          `json:"subsystems"`
	Health                   SessionSidecarHealth        `json:"health"`
	DisplayApps              []DisplayApp                `json:"display_apps"`
	LogSources               []SessionSidecarLogSource   `json:"log_sources"`
	Scripts                  []SessionSidecarScript      `json:"scripts"`
}

// SessionSidecarLog is a single log entry from a session sidecar.
type SessionSidecarLog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	Output    string    `json:"output"`
	Level     LogLevel  `json:"level"`
	SourceID  uuid.UUID `json:"source_id" format:"uuid"`
}

// SessionSidecarListeningPortsResponse contains listening ports for a sidecar.
type SessionSidecarListeningPortsResponse struct {
	Ports []SessionSidecarListeningPort `json:"ports"`
}

// SessionSidecarListeningPort describes a port being listened on inside a sidecar.
type SessionSidecarListeningPort struct {
	ProcessName string `json:"process_name"`
	Network     string `json:"network"`
	Port        uint16 `json:"port"`
}

// RealTimeSessionStatus represents the status of a real-time session.
type RealTimeSessionStatus string

const (
	RealTimeSessionStatusPending   RealTimeSessionStatus = "pending"
	RealTimeSessionStatusActive    RealTimeSessionStatus = "active"
	RealTimeSessionStatusCompleted RealTimeSessionStatus = "completed"
	RealTimeSessionStatusFailed    RealTimeSessionStatus = "failed"
	RealTimeSessionStatusCanceled  RealTimeSessionStatus = "canceled"
)

// RealTimeSession represents a real-time LiveKit session created from a preset template.
// Unlike the agent-based Session type, this is a lightweight entity for real-time
// audio/video sessions that does not involve provisioning or builds.
type RealTimeSession struct {
	ID                uuid.UUID             `json:"id" format:"uuid"`
	CreatedAt         time.Time             `json:"created_at" format:"date-time"`
	UpdatedAt         time.Time             `json:"updated_at" format:"date-time"`
	OrganizationID    uuid.UUID             `json:"organization_id" format:"uuid"`
	TemplateID        uuid.UUID             `json:"template_id" format:"uuid"`
	TemplateVersionID uuid.UUID             `json:"template_version_id" format:"uuid"`
	OwnerID           uuid.UUID             `json:"owner_id" format:"uuid"`
	Name              string                `json:"name"`
	DisplayName       string                `json:"display_name"`
	RoomName          string                `json:"room_name"`
	RoomID            string                `json:"room_id"`
	Status            RealTimeSessionStatus `json:"status"`
	StartedAt         *time.Time            `json:"started_at,omitempty" format:"date-time"`
	EndedAt           *time.Time            `json:"ended_at,omitempty" format:"date-time"`
	DurationSeconds   *int32                `json:"duration_seconds,omitempty"`
	Transcript        string                `json:"transcript,omitempty"`
	TranscriptURL     string                `json:"transcript_url,omitempty"`
	RecordingURL      string                `json:"recording_url,omitempty"`
	Metadata          json.RawMessage       `json:"metadata,omitempty"`
	OwnerName         string                `json:"owner_name,omitempty"`
	TemplateName      string                `json:"template_name,omitempty"`
	OrganizationName  string                `json:"organization_name,omitempty"`
}

// RealTimeSessionParameter stores a resolved parameter value for a real-time session.
type RealTimeSessionParameter struct {
	SessionID uuid.UUID `json:"session_id" format:"uuid"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
}

// RealTimeSessionParameterInput is a parameter name/value pair for session creation.
type RealTimeSessionParameterInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateRealTimeSessionRequest is the request body for creating a new real-time session.
type CreateRealTimeSessionRequest struct {
	TemplateID        uuid.UUID                       `json:"template_id" validate:"required" format:"uuid"`
	TemplateVersionID uuid.UUID                       `json:"template_version_id,omitempty" format:"uuid"`
	Name              string                          `json:"name" validate:"required"`
	DisplayName       string                          `json:"display_name,omitempty"`
	RoomName          string                          `json:"room_name,omitempty"`
	Metadata          json.RawMessage                 `json:"metadata,omitempty"`
	Parameters        []RealTimeSessionParameterInput `json:"parameters,omitempty"`
}

// UpdateRealTimeSessionRequest is the request body for updating a real-time session.
type UpdateRealTimeSessionRequest struct {
	Status        *RealTimeSessionStatus `json:"status,omitempty"`
	Transcript    *string                `json:"transcript,omitempty"`
	TranscriptURL *string                `json:"transcript_url,omitempty"`
	RecordingURL  *string                `json:"recording_url,omitempty"`
	RoomID        *string                `json:"room_id,omitempty"`
}

// RealTimeSessionFilter allows filtering real-time sessions by various criteria.
type RealTimeSessionFilter struct {
	OrganizationID uuid.UUID             `json:"organization_id,omitempty"`
	OwnerID        uuid.UUID             `json:"owner_id,omitempty"`
	TemplateID     uuid.UUID             `json:"template_id,omitempty"`
	Status         RealTimeSessionStatus `json:"status,omitempty"`
	Search         string                `json:"q,omitempty"`
}

// asRequestOption returns a RequestOption that applies the filter as query parameters.
func (f RealTimeSessionFilter) asRequestOption() client.RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		if f.TemplateID != uuid.Nil {
			q.Set("template_id", f.TemplateID.String())
		}
		if f.OwnerID != uuid.Nil {
			q.Set("owner_id", f.OwnerID.String())
		}
		if f.Status != "" {
			q.Set("status", string(f.Status))
		}
		if f.Search != "" {
			q.Set("q", f.Search)
		}
		r.URL.RawQuery = q.Encode()
	}
}

// RealTimeSessionsResponse contains a list of real-time sessions.
type RealTimeSessionsResponse struct {
	Sessions []RealTimeSession `json:"sessions"`
	Count    int               `json:"count"`
}
