package types

import (
	"time"

	"github.com/google/uuid"
)

// Agent represents a deployed agent instance.
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

// AgentBuild represents a build of an agent.
type AgentBuild struct {
	ID                uuid.UUID       `json:"id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	AgentID           uuid.UUID       `json:"agent_id"`
	AgentName         string          `json:"agent_name"`
	TemplateVersionID uuid.UUID       `json:"template_version_id"`
	BuildNumber       int32           `json:"build_number"`
	Transition        AgentTransition `json:"transition"`
	InitiatorID       uuid.UUID       `json:"initiator_id"`
	InitiatorName     string          `json:"initiator_username"`
	JobID             uuid.UUID       `json:"job_id"`
	Reason            BuildReason     `json:"reason"`
	Status            string          `json:"status"`
}

// AgentHealth contains health information for an agent.
type AgentHealth struct {
	Healthy         bool        `json:"healthy"`
	FailingSidecars []uuid.UUID `json:"failing_sidecars"`
}

// AgentFilter is used to filter agent list requests.
type AgentFilter struct {
	Owner       string `json:"owner,omitempty"`
	Template    string `json:"template,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
	Offset      int    `json:"offset,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	FilterQuery string `json:"q,omitempty"`
}

// AgentQuota contains budget information for an agent.
type AgentQuota struct {
	CreditsConsumed int `json:"credits_consumed"`
	Budget          int `json:"budget"`
}

// CreateAgentRequest contains the parameters for creating a new agent.
type CreateAgentRequest struct {
	TemplateID          uuid.UUID            `json:"template_id,omitempty"`
	TemplateVersionID   uuid.UUID            `json:"template_version_id,omitempty"`
	Name                string               `json:"name"`
	AutostartSchedule   *string              `json:"autostart_schedule,omitempty"`
	TTLMillis           *int64               `json:"ttl_ms,omitempty"`
	RichParameterValues []AgentBuildParameter `json:"rich_parameter_values,omitempty"`
	AutomaticUpdates    AutomaticUpdates     `json:"automatic_updates,omitempty"`
}

// AgentBuildParameter represents a parameter value for a build.
type AgentBuildParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateAgentBuildRequest contains the parameters for creating a new agent build.
type CreateAgentBuildRequest struct {
	TemplateVersionID   uuid.UUID            `json:"template_version_id,omitempty"`
	Transition          AgentTransition      `json:"transition"`
	DryRun              bool                 `json:"dry_run,omitempty"`
	Orphan              bool                 `json:"orphan,omitempty"`
	RichParameterValues []AgentBuildParameter `json:"rich_parameter_values,omitempty"`
	LogLevel            ProvisionerLogLevel  `json:"log_level,omitempty"`
}

// UpdateAgentRequest contains updatable agent fields.
type UpdateAgentRequest struct {
	Name string `json:"name,omitempty"`
}

// UpdateAgentAutostartRequest updates the agent autostart schedule.
type UpdateAgentAutostartRequest struct {
	Schedule *string `json:"schedule"`
}

// UpdateAgentTTLRequest updates the agent TTL.
type UpdateAgentTTLRequest struct {
	TTLMillis *int64 `json:"ttl_ms"`
}

// PostAgentUsageRequest reports agent usage.
type PostAgentUsageRequest struct {
	SidecarID uuid.UUID    `json:"sidecar_id"`
	AppName   UsageAppName `json:"app_name"`
}
