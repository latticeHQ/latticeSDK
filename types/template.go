package types

import (
	"time"

	"github.com/google/uuid"
)

// TemplateCategory classifies templates.
type TemplateCategory string

const (
	TemplateCategoryAgent  TemplateCategory = "agent"
	TemplateCategoryPreset TemplateCategory = "preset"
	TemplateCategoryEval   TemplateCategory = "eval"
)

// TemplateRole defines a user's role on a template.
type TemplateRole string

const (
	TemplateRoleAdmin TemplateRole = "admin"
	TemplateRoleUse   TemplateRole = "use"
)

// Template represents an agent template.
type Template struct {
	ID                             uuid.UUID                   `json:"id"`
	CreatedAt                      time.Time                   `json:"created_at"`
	UpdatedAt                      time.Time                   `json:"updated_at"`
	OrganizationID                 uuid.UUID                   `json:"organization_id"`
	OrganizationName               string                      `json:"organization_name"`
	OrganizationDisplayName        string                      `json:"organization_display_name"`
	OrganizationIcon               string                      `json:"organization_icon"`
	Name                           string                      `json:"name"`
	DisplayName                    string                      `json:"display_name"`
	Provisioner                    ProvisionerType             `json:"provisioner"`
	ActiveVersionID                uuid.UUID                   `json:"active_version_id"`
	ActiveUserCount                int                         `json:"active_user_count"`
	BuildTimeStats                 map[AgentTransition]TransitionStats `json:"build_time_stats"`
	Description                    string                      `json:"description"`
	Deprecated                     bool                        `json:"deprecated"`
	DeprecationMessage             string                      `json:"deprecation_message"`
	Icon                           string                      `json:"icon"`
	DefaultTTLMillis               int64                       `json:"default_ttl_ms"`
	ActivityBumpMillis             int64                       `json:"activity_bump_ms"`
	AutostopRequirement            TemplateAutostopRequirement `json:"autostop_requirement"`
	AutostartRequirement           TemplateAutostartRequirement `json:"autostart_requirement"`
	CreatedByID                    uuid.UUID                   `json:"created_by_id"`
	CreatedByName                  string                      `json:"created_by_name"`
	AllowUserAutostart             bool                        `json:"allow_user_autostart"`
	AllowUserAutostop              bool                        `json:"allow_user_autostop"`
	AllowUserCancelAgentJobs       bool                        `json:"allow_user_cancel_agent_jobs"`
	FailureTTLMillis               int64                       `json:"failure_ttl_ms"`
	TimeTilDormantMillis           int64                       `json:"time_til_dormant_ms"`
	TimeTilDormantAutoDeleteMillis int64                       `json:"time_til_dormant_autodelete_ms"`
	RequireActiveVersion           bool                        `json:"require_active_version"`
	Category                       TemplateCategory            `json:"category"`
	RequiresProvisioning           bool                        `json:"requires_provisioning"`
	ParentTemplateID               *uuid.UUID                  `json:"parent_template_id,omitempty"`
}

// TransitionStats contains percentile build time statistics.
type TransitionStats struct {
	P50 *int64 `json:"p50"`
	P95 *int64 `json:"p95"`
}

// TemplateAutostopRequirement configures autostop behavior.
type TemplateAutostopRequirement struct {
	DaysOfWeek []string `json:"days_of_week"`
	Weeks      int64    `json:"weeks"`
}

// TemplateAutostartRequirement configures autostart behavior.
type TemplateAutostartRequirement struct {
	DaysOfWeek []string `json:"days_of_week"`
}

// TemplateVersion represents a version of a template.
type TemplateVersion struct {
	ID             uuid.UUID `json:"id"`
	TemplateID     uuid.UUID `json:"template_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name"`
	Message        string    `json:"message"`
	CreatedByID    uuid.UUID `json:"created_by_id"`
	CreatedByName  string    `json:"created_by_name"`
	Archived       bool      `json:"archived"`
}

// TemplateFilter is used to filter template list requests.
type TemplateFilter struct {
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	ExactName      string    `json:"exact_name,omitempty"`
	FuzzyName      string    `json:"fuzzy_name,omitempty"`
	SearchQuery    string    `json:"q,omitempty"`
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
