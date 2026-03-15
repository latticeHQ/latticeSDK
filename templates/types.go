// Package templates provides template management services for the Lattice Runtime API.
// It covers template CRUD, versions, parameters, variables, ACLs, and dry runs.
package templates

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// Template is the full representation of a Lattice template.
type Template struct {
	ID                             uuid.UUID                              `json:"id"`
	CreatedAt                      time.Time                              `json:"created_at"`
	UpdatedAt                      time.Time                              `json:"updated_at"`
	OrganizationID                 uuid.UUID                              `json:"organization_id"`
	OrganizationName               string                                 `json:"organization_name"`
	OrganizationDisplayName        string                                 `json:"organization_display_name"`
	OrganizationIcon               string                                 `json:"organization_icon"`
	Name                           string                                 `json:"name"`
	DisplayName                    string                                 `json:"display_name"`
	Provisioner                    types.ProvisionerType                  `json:"provisioner"`
	ActiveVersionID                uuid.UUID                              `json:"active_version_id"`
	ActiveUserCount                int                                    `json:"active_user_count"`
	BuildTimeStats                 map[types.AgentTransition]TransitionStats `json:"build_time_stats"`
	Description                    string                                 `json:"description"`
	Deprecated                     bool                                   `json:"deprecated"`
	DeprecationMessage             string                                 `json:"deprecation_message"`
	Icon                           string                                 `json:"icon"`
	DefaultTTLMillis               int64                                  `json:"default_ttl_ms"`
	ActivityBumpMillis             int64                                  `json:"activity_bump_ms"`
	AutostopRequirement            AutostopRequirement                    `json:"autostop_requirement"`
	AutostartRequirement           AutostartRequirement                   `json:"autostart_requirement"`
	CreatedByID                    uuid.UUID                              `json:"created_by_id"`
	CreatedByName                  string                                 `json:"created_by_name"`
	AllowUserAutostart             bool                                   `json:"allow_user_autostart"`
	AllowUserAutostop              bool                                   `json:"allow_user_autostop"`
	AllowUserCancelAgentJobs       bool                                   `json:"allow_user_cancel_agent_jobs"`
	FailureTTLMillis               int64                                  `json:"failure_ttl_ms"`
	TimeTilDormantMillis           int64                                  `json:"time_til_dormant_ms"`
	TimeTilDormantAutoDeleteMillis int64                                  `json:"time_til_dormant_autodelete_ms"`
	RequireActiveVersion           bool                                   `json:"require_active_version"`
	MaxPortShareLevel              PortShareLevel                         `json:"max_port_share_level"`
	Category                       TemplateCategory                       `json:"category"`
	RequiresProvisioning           bool                                   `json:"requires_provisioning"`
	ParentTemplateID               *uuid.UUID                             `json:"parent_template_id,omitempty"`
}

// TemplateCategory classifies templates.
type TemplateCategory string

const (
	// TemplateCategoryAgent is for agent templates that provision infrastructure.
	TemplateCategoryAgent TemplateCategory = "agent"
	// TemplateCategoryPreset is for preset templates without provisioning.
	TemplateCategoryPreset TemplateCategory = "preset"
	// TemplateCategoryEval is for evaluation templates.
	TemplateCategoryEval TemplateCategory = "eval"
)

// TemplateRole defines a user's role on a template.
type TemplateRole string

const (
	// TemplateRoleAdmin grants full admin access to a template.
	TemplateRoleAdmin TemplateRole = "admin"
	// TemplateRoleUse grants usage access to a template.
	TemplateRoleUse TemplateRole = "use"
)

// PortShareLevel controls agent port sharing behavior.
type PortShareLevel string

const (
	// PortShareLevelOwner restricts port sharing to the agent owner.
	PortShareLevelOwner PortShareLevel = "owner"
	// PortShareLevelAuthenticated allows port sharing with authenticated users.
	PortShareLevelAuthenticated PortShareLevel = "authenticated"
	// PortShareLevelPublic allows public port sharing.
	PortShareLevelPublic PortShareLevel = "public"
)

// TransitionStats contains percentile build time statistics.
type TransitionStats struct {
	P50 *int64 `json:"p50"`
	P95 *int64 `json:"p95"`
}

// AutostopRequirement configures required restart behavior for agents.
type AutostopRequirement struct {
	// DaysOfWeek is a list of days of the week on which restarts are required.
	DaysOfWeek []string `json:"days_of_week"`
	// Weeks is the number of weeks between required restarts.
	Weeks int64 `json:"weeks"`
}

// AutostartRequirement configures autostart behavior for agents.
type AutostartRequirement struct {
	// DaysOfWeek is a list of days of the week in which autostart is allowed.
	DaysOfWeek []string `json:"days_of_week"`
}

// TemplateFilter is used to filter template list requests.
type TemplateFilter struct {
	OrganizationID uuid.UUID
	ExactName      string
	FuzzyName      string
	SearchQuery    string
}

// AsRequestOption returns a RequestOption that applies filter query parameters.
func (f TemplateFilter) AsRequestOption() client.RequestOption {
	return func(r *http.Request) {
		var params []string
		if f.OrganizationID != uuid.Nil {
			params = append(params, fmt.Sprintf("organization:%q", f.OrganizationID.String()))
		}
		if f.ExactName != "" {
			params = append(params, fmt.Sprintf("exact_name:%q", f.ExactName))
		}
		if f.FuzzyName != "" {
			params = append(params, fmt.Sprintf("name:%q", f.FuzzyName))
		}
		if f.SearchQuery != "" {
			params = append(params, f.SearchQuery)
		}
		q := r.URL.Query()
		q.Set("q", strings.Join(params, " "))
		r.URL.RawQuery = q.Encode()
	}
}

// CreateTemplateRequest provides options when creating a template.
type CreateTemplateRequest struct {
	// Name is the name of the template.
	Name string `json:"name"`
	// DisplayName is the displayed name of the template.
	DisplayName string `json:"display_name,omitempty"`
	// Description is a description of what the template contains.
	Description string `json:"description,omitempty"`
	// Icon is a relative path or external URL for the template icon.
	Icon string `json:"icon,omitempty"`
	// VersionID is an in-progress or completed job to use as the initial version.
	VersionID uuid.UUID `json:"template_version_id"`
	// DefaultTTLMillis optionally specifies the default TTL for agents.
	DefaultTTLMillis *int64 `json:"default_ttl_ms,omitempty"`
	// ActivityBumpMillis optionally specifies the activity bump duration.
	ActivityBumpMillis *int64 `json:"activity_bump_ms,omitempty"`
	// AutostopRequirement is an enterprise feature for autostop configuration.
	AutostopRequirement *AutostopRequirement `json:"autostop_requirement,omitempty"`
	// AutostartRequirement is an enterprise feature for autostart configuration.
	AutostartRequirement *AutostartRequirement `json:"autostart_requirement,omitempty"`
	// AllowUserCancelAgentJobs allows users to cancel in-progress agent jobs.
	AllowUserCancelAgentJobs *bool `json:"allow_user_cancel_agent_jobs"`
	// AllowUserAutostart allows users to set autostart schedules.
	AllowUserAutostart *bool `json:"allow_user_autostart"`
	// AllowUserAutostop allows users to set custom agent TTL.
	AllowUserAutostop *bool `json:"allow_user_autostop"`
	// FailureTTLMillis optionally specifies the max lifetime for failed agents.
	FailureTTLMillis *int64 `json:"failure_ttl_ms,omitempty"`
	// TimeTilDormantMillis optionally specifies the max lifetime before locking.
	TimeTilDormantMillis *int64 `json:"dormant_ttl_ms,omitempty"`
	// TimeTilDormantAutoDeleteMillis optionally specifies the max lifetime before deletion.
	TimeTilDormantAutoDeleteMillis *int64 `json:"delete_ttl_ms,omitempty"`
	// DisableEveryoneGroupAccess disables default everyone-group template access.
	DisableEveryoneGroupAccess bool `json:"disable_everyone_group_access"`
	// RequireActiveVersion mandates agents are built with the active version.
	RequireActiveVersion bool `json:"require_active_version"`
	// MaxPortShareLevel optionally specifies the maximum port share level.
	MaxPortShareLevel *PortShareLevel `json:"max_port_share_level"`
	// Category is the template category: agent, preset, eval.
	Category TemplateCategory `json:"category,omitempty"`
	// RequiresProvisioning indicates whether this template requires terraform provisioning.
	RequiresProvisioning *bool `json:"requires_provisioning,omitempty"`
	// ParentTemplateID is the parent template for inheritance.
	ParentTemplateID *uuid.UUID `json:"parent_template_id,omitempty"`
}

// UpdateTemplateMeta contains fields for updating template metadata.
type UpdateTemplateMeta struct {
	Name                           string               `json:"name,omitempty"`
	DisplayName                    string               `json:"display_name,omitempty"`
	Description                    string               `json:"description,omitempty"`
	Category                       TemplateCategory     `json:"category,omitempty"`
	Icon                           string               `json:"icon,omitempty"`
	DefaultTTLMillis               int64                `json:"default_ttl_ms,omitempty"`
	ActivityBumpMillis             int64                `json:"activity_bump_ms,omitempty"`
	AutostopRequirement            *AutostopRequirement `json:"autostop_requirement,omitempty"`
	AutostartRequirement           *AutostartRequirement `json:"autostart_requirement,omitempty"`
	AllowUserAutostart             bool                 `json:"allow_user_autostart,omitempty"`
	AllowUserAutostop              bool                 `json:"allow_user_autostop,omitempty"`
	AllowUserCancelAgentJobs       bool                 `json:"allow_user_cancel_agent_jobs,omitempty"`
	FailureTTLMillis               int64                `json:"failure_ttl_ms,omitempty"`
	TimeTilDormantMillis           int64                `json:"time_til_dormant_ms,omitempty"`
	TimeTilDormantAutoDeleteMillis int64                `json:"time_til_dormant_autodelete_ms,omitempty"`
	UpdateAgentLastUsedAt          bool                 `json:"update_agent_last_used_at"`
	UpdateAgentDormantAt           bool                 `json:"update_agent_dormant_at"`
	RequireActiveVersion           bool                 `json:"require_active_version,omitempty"`
	DeprecationMessage             *string              `json:"deprecation_message"`
	DisableEveryoneGroupAccess     bool                 `json:"disable_everyone_group_access"`
	MaxPortShareLevel              *PortShareLevel      `json:"max_port_share_level"`
}

// UpdateActiveTemplateVersion identifies the version to promote.
type UpdateActiveTemplateVersion struct {
	ID uuid.UUID `json:"id"`
}

// ArchiveTemplateVersionsRequest controls which versions are archived.
type ArchiveTemplateVersionsRequest struct {
	// All archives all unused versions when true; otherwise only failed versions.
	All bool `json:"all"`
}

// ArchiveTemplateVersionsResponse contains the result of an archive operation.
type ArchiveTemplateVersionsResponse struct {
	TemplateID  uuid.UUID   `json:"template_id"`
	ArchivedIDs []uuid.UUID `json:"archived_ids"`
}

// TemplateACL contains the users and groups with access to a template.
type TemplateACL struct {
	Users  []TemplateUser  `json:"users"`
	Groups []TemplateGroup `json:"group"`
}

// TemplateUser is a user with a template-specific role.
type TemplateUser struct {
	types.User
	Role TemplateRole `json:"role"`
}

// TemplateGroup is a group with a template-specific role.
type TemplateGroup struct {
	Group
	Role TemplateRole `json:"role"`
}

// Group represents a group in the system.
type Group struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	OrganizationID uuid.UUID `json:"organization_id"`
	AvatarURL      string    `json:"avatar_url"`
	Source         string    `json:"source"`
	QuotaAllowance int      `json:"quota_allowance"`
}

// UpdateTemplateACL maps user and group IDs to roles.
type UpdateTemplateACL struct {
	// UserPerms maps user UUIDs to roles.
	UserPerms map[string]TemplateRole `json:"user_perms,omitempty"`
	// GroupPerms maps group UUIDs to roles.
	GroupPerms map[string]TemplateRole `json:"group_perms,omitempty"`
}

// ACLAvailable lists users and groups that can be added to a template ACL.
type ACLAvailable struct {
	Users  []ReducedUser `json:"users"`
	Groups []Group       `json:"groups"`
}

// ReducedUser is a lightweight user representation for ACL listings.
type ReducedUser struct {
	types.MinimalUser
	Name   string          `json:"name"`
	Email  string          `json:"email"`
	Status types.UserStatus `json:"status"`
}

// TemplateExample represents a starter template.
type TemplateExample struct {
	ID          string   `json:"id"`
	URL         string   `json:"url"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Tags        []string `json:"tags"`
	Markdown    string   `json:"markdown"`
}

// DAUsResponse contains daily active user statistics for a template.
type DAUsResponse struct {
	Entries      []DAUEntry `json:"entries"`
	TZHourOffset int        `json:"tz_hour_offset"`
}

// DAUEntry is a single daily active user data point.
type DAUEntry struct {
	// Date is formatted as 2024-01-31 (no timezone info).
	Date   string `json:"date"`
	Amount int    `json:"amount"`
}

// DAURequest configures the timezone offset for DAU queries.
type DAURequest struct {
	TZHourOffset int
}

// AsRequestOption returns a RequestOption that sets the tz_offset query parameter.
func (d DAURequest) AsRequestOption() client.RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set("tz_offset", strconv.Itoa(d.TZHourOffset))
		r.URL.RawQuery = q.Encode()
	}
}

// TemplateVersion represents a single version of a template.
type TemplateVersion struct {
	ID                  uuid.UUID                `json:"id"`
	TemplateID          *uuid.UUID               `json:"template_id,omitempty"`
	OrganizationID      uuid.UUID                `json:"organization_id,omitempty"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	Name                string                   `json:"name"`
	Message             string                   `json:"message"`
	Job                 ProvisionerJob           `json:"job"`
	Readme              string                   `json:"readme"`
	CreatedBy           types.MinimalUser        `json:"created_by"`
	Archived            bool                     `json:"archived"`
	Warnings            []TemplateVersionWarning `json:"warnings,omitempty"`
	MatchedProvisioners MatchedProvisioners      `json:"matched_provisioners,omitempty"`
}

// TemplateVersionWarning identifies warnings on a template version.
type TemplateVersionWarning string

const (
	// TemplateVersionWarningUnsupportedAgents warns about unsupported agent configurations.
	TemplateVersionWarningUnsupportedAgents TemplateVersionWarning = "UNSUPPORTED_AGENTS"
)

// ProvisionerStorageMethod identifies how provisioner files are stored.
type ProvisionerStorageMethod string

const (
	// ProvisionerStorageMethodFile stores provisioner files on disk.
	ProvisionerStorageMethodFile ProvisionerStorageMethod = "file"
)

// CreateTemplateVersionRequest provides options when creating a template version.
type CreateTemplateVersionRequest struct {
	// Name is the version name.
	Name string `json:"name,omitempty"`
	// Message is an optional description of the version.
	Message string `json:"message,omitempty"`
	// TemplateID optionally associates the version with a template.
	TemplateID uuid.UUID `json:"template_id,omitempty"`
	// StorageMethod is the provisioner storage method (currently only "file").
	StorageMethod ProvisionerStorageMethod `json:"storage_method"`
	// FileID references the uploaded file containing the template source.
	FileID uuid.UUID `json:"file_id,omitempty"`
	// ExampleID references a starter template to use as the source.
	ExampleID string `json:"example_id,omitempty"`
	// Provisioner is the provisioner type (terraform, echo).
	Provisioner types.ProvisionerType `json:"provisioner"`
	// ProvisionerTags are key-value tags for provisioner matching.
	ProvisionerTags map[string]string `json:"tags"`
	// UserVariableValues are user-provided variable values.
	UserVariableValues []VariableValue `json:"user_variable_values,omitempty"`
}

// PatchTemplateVersionRequest contains fields for updating a template version.
type PatchTemplateVersionRequest struct {
	// Name updates the version name.
	Name string `json:"name,omitempty"`
	// Message updates the version message.
	Message *string `json:"message,omitempty"`
}

// VariableValue is a name-value pair for template variables.
type VariableValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TemplateVersionsByTemplateRequest configures listing versions for a template.
type TemplateVersionsByTemplateRequest struct {
	TemplateID      uuid.UUID `json:"template_id"`
	IncludeArchived bool      `json:"include_archived"`
	Pagination      *client.Pagination
}

// TemplateVersionParameter represents a parameter for a template version.
type TemplateVersionParameter struct {
	Name                 string                           `json:"name"`
	DisplayName          string                           `json:"display_name,omitempty"`
	Description          string                           `json:"description"`
	DescriptionPlaintext string                           `json:"description_plaintext"`
	Type                 string                           `json:"type"`
	Mutable              bool                             `json:"mutable"`
	DefaultValue         string                           `json:"default_value"`
	Icon                 string                           `json:"icon"`
	Options              []TemplateVersionParameterOption `json:"options"`
	ValidationError      string                           `json:"validation_error,omitempty"`
	ValidationRegex      string                           `json:"validation_regex,omitempty"`
	ValidationMin        *int32                           `json:"validation_min,omitempty"`
	ValidationMax        *int32                           `json:"validation_max,omitempty"`
	ValidationMonotonic  ValidationMonotonicOrder         `json:"validation_monotonic,omitempty"`
	Required             bool                             `json:"required"`
	Ephemeral            bool                             `json:"ephemeral"`
}

// TemplateVersionParameterOption represents a selectable option for a template parameter.
type TemplateVersionParameterOption struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Icon        string `json:"icon"`
}

// ValidationMonotonicOrder controls monotonic validation direction.
type ValidationMonotonicOrder string

const (
	// MonotonicOrderIncreasing requires values to increase.
	MonotonicOrderIncreasing ValidationMonotonicOrder = "increasing"
	// MonotonicOrderDecreasing requires values to decrease.
	MonotonicOrderDecreasing ValidationMonotonicOrder = "decreasing"
)

// TemplateVersionVariable represents a managed template variable.
type TemplateVersionVariable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	Value        string `json:"value"`
	DefaultValue string `json:"default_value"`
	Required     bool   `json:"required"`
	Sensitive    bool   `json:"sensitive"`
}

// TemplateVersionExternalAuth describes an external authentication provider for a template version.
type TemplateVersionExternalAuth struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	DisplayName     string `json:"display_name"`
	DisplayIcon     string `json:"display_icon"`
	AuthenticateURL string `json:"authenticate_url"`
	Authenticated   bool   `json:"authenticated"`
	Optional        bool   `json:"optional,omitempty"`
}

// ProvisionerJobStatus represents the at-time state of a provisioner job.
type ProvisionerJobStatus string

const (
	// ProvisionerJobPending means the job is waiting to be picked up.
	ProvisionerJobPending ProvisionerJobStatus = "pending"
	// ProvisionerJobRunning means the job is currently executing.
	ProvisionerJobRunning ProvisionerJobStatus = "running"
	// ProvisionerJobSucceeded means the job completed successfully.
	ProvisionerJobSucceeded ProvisionerJobStatus = "succeeded"
	// ProvisionerJobCanceling means the job is being canceled.
	ProvisionerJobCanceling ProvisionerJobStatus = "canceling"
	// ProvisionerJobCanceled means the job was canceled.
	ProvisionerJobCanceled ProvisionerJobStatus = "canceled"
	// ProvisionerJobFailed means the job failed.
	ProvisionerJobFailed ProvisionerJobStatus = "failed"
)

// JobErrorCode identifies specific provisioner job error conditions.
type JobErrorCode string

const (
	// JobErrorCodeRequiredTemplateVariables indicates required variables were missing.
	JobErrorCodeRequiredTemplateVariables JobErrorCode = "REQUIRED_TEMPLATE_VARIABLES"
)

// ProvisionerJob represents a provisioner job.
type ProvisionerJob struct {
	ID          uuid.UUID            `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	StartedAt   *time.Time           `json:"started_at,omitempty"`
	CompletedAt *time.Time           `json:"completed_at,omitempty"`
	CanceledAt  *time.Time           `json:"canceled_at,omitempty"`
	Error       string               `json:"error,omitempty"`
	ErrorCode   JobErrorCode         `json:"error_code,omitempty"`
	Status      ProvisionerJobStatus `json:"status"`
	WorkerID    *uuid.UUID           `json:"worker_id,omitempty"`
	FileID      uuid.UUID            `json:"file_id"`
	Tags        map[string]string    `json:"tags"`
	QueuePosition int               `json:"queue_position"`
	QueueSize     int               `json:"queue_size"`
}

// LogSource identifies the source of a provisioner log entry.
type LogSource string

const (
	// LogSourceProvisionerDaemon logs from the provisioner daemon process.
	LogSourceProvisionerDaemon LogSource = "provisioner_daemon"
	// LogSourceProvisioner logs from the provisioner itself.
	LogSourceProvisioner LogSource = "provisioner"
)

// LogLevel represents log verbosity.
type LogLevel string

const (
	// LogLevelTrace is the most verbose log level.
	LogLevelTrace LogLevel = "trace"
	// LogLevelDebug is for debug-level messages.
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo is for informational messages.
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn is for warning messages.
	LogLevelWarn LogLevel = "warn"
	// LogLevelError is for error messages.
	LogLevelError LogLevel = "error"
)

// ProvisionerJobLog represents a provisioner log entry.
type ProvisionerJobLog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Source    LogSource `json:"log_source"`
	Level     LogLevel  `json:"log_level"`
	Stage     string    `json:"stage"`
	Output    string    `json:"output"`
}

// MatchedProvisioners contains information about provisioners that matched a job's tags.
type MatchedProvisioners struct {
	// Count is the total number of matching provisioner daemons.
	Count int `json:"count"`
	// Available is the number of provisioners available to take jobs.
	Available int `json:"available"`
	// MostRecentlySeen is the last seen time of the matched provisioners.
	MostRecentlySeen *time.Time `json:"most_recently_seen,omitempty"`
}

// AgentResource represents a provisioned resource associated with an agent.
type AgentResource struct {
	ID         uuid.UUID               `json:"id"`
	CreatedAt  time.Time               `json:"created_at"`
	JobID      uuid.UUID               `json:"job_id"`
	Transition types.AgentTransition   `json:"agent_transition"`
	Type       string                  `json:"type"`
	Name       string                  `json:"name"`
	Hide       bool                    `json:"hide"`
	Icon       string                  `json:"icon"`
	Metadata   []AgentResourceMetadata `json:"metadata,omitempty"`
	DailyCost  int32                   `json:"daily_cost"`
}

// AgentResourceMetadata annotates an agent resource with custom key-value pairs.
type AgentResourceMetadata struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
}

// CreateTemplateVersionDryRunRequest defines the request parameters for a dry run.
type CreateTemplateVersionDryRunRequest struct {
	// AgentName is the name of the agent to simulate.
	AgentName string `json:"agent_name"`
	// RichParameterValues are parameter values for the dry run.
	RichParameterValues []types.AgentBuildParameter `json:"rich_parameter_values"`
	// UserVariableValues are variable values for the dry run.
	UserVariableValues []VariableValue `json:"user_variable_values,omitempty"`
}

// ParameterSource indicates where a resolved parameter value originated.
type ParameterSource string

const (
	// ParameterSourceDefault means the value came from the parameter's own default.
	ParameterSourceDefault ParameterSource = "default"
	// ParameterSourceParentTemplate means the value was inherited from a parent template.
	ParameterSourceParentTemplate ParameterSource = "parent_template"
	// ParameterSourceChildTemplate means the value was set by the child template override.
	ParameterSourceChildTemplate ParameterSource = "child_template"
	// ParameterSourceInstanceOverride means the value was provided at instance creation time.
	ParameterSourceInstanceOverride ParameterSource = "instance_override"
)

// ResolvedParameter represents a parameter with its resolved value and provenance.
type ResolvedParameter struct {
	// Parameter is the full parameter definition.
	Parameter TemplateVersionParameter `json:"parameter"`
	// Value is the resolved value after inheritance.
	Value string `json:"value"`
	// Source indicates where this value originated.
	Source ParameterSource `json:"source"`
	// SourceTemplateName is the name of the template that provided this value.
	SourceTemplateName string `json:"source_template_name,omitempty"`
}
