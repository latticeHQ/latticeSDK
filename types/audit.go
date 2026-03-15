package types

import (
	"encoding/json"
	"net/netip"
	"time"

	"github.com/google/uuid"
)

// ResourceType identifies the type of resource in an audit log.
type ResourceType string

const (
	ResourceTypeTemplate        ResourceType = "template"
	ResourceTypeTemplateVersion ResourceType = "template_version"
	ResourceTypeUser            ResourceType = "user"
	ResourceTypeAgent           ResourceType = "agent"
	ResourceTypeAgentBuild      ResourceType = "agent_build"
	ResourceTypeAPIKey          ResourceType = "api_key"
	ResourceTypeGroup           ResourceType = "group"
	ResourceTypeLicense         ResourceType = "license"
	ResourceTypeOrganization    ResourceType = "organization"
	ResourceTypeCustomRole      ResourceType = "custom_role"
)

// AuditAction identifies the action taken in an audit log.
type AuditAction string

const (
	AuditActionCreate   AuditAction = "create"
	AuditActionWrite    AuditAction = "write"
	AuditActionDelete   AuditAction = "delete"
	AuditActionStart    AuditAction = "start"
	AuditActionStop     AuditAction = "stop"
	AuditActionLogin    AuditAction = "login"
	AuditActionLogout   AuditAction = "logout"
	AuditActionRegister AuditAction = "register"
)

// AuditDiff maps field names to their old and new values.
type AuditDiff map[string]AuditDiffField

// AuditDiffField shows the old and new values for a changed field.
type AuditDiffField struct {
	Old    interface{} `json:"old,omitempty"`
	New    interface{} `json:"new,omitempty"`
	Secret bool        `json:"secret"`
}

// AuditLog represents an entry in the audit trail.
type AuditLog struct {
	ID               uuid.UUID            `json:"id"`
	RequestID        uuid.UUID            `json:"request_id"`
	Time             time.Time            `json:"time"`
	IP               netip.Addr           `json:"ip"`
	UserSidecar      string               `json:"user_sidecar"`
	ResourceType     ResourceType         `json:"resource_type"`
	ResourceID       uuid.UUID            `json:"resource_id"`
	ResourceTarget   string               `json:"resource_target"`
	ResourceIcon     string               `json:"resource_icon"`
	Action           AuditAction          `json:"action"`
	Diff             AuditDiff            `json:"diff"`
	StatusCode       int32                `json:"status_code"`
	AdditionalFields json.RawMessage      `json:"additional_fields"`
	Description      string               `json:"description"`
	ResourceLink     string               `json:"resource_link"`
	IsDeleted        bool                 `json:"is_deleted"`
	OrganizationID   uuid.UUID            `json:"organization_id"`
	Organization     *MinimalOrganization `json:"organization,omitempty"`
	User             *User                `json:"user"`
}

// AuditLogFilter is used to filter audit log queries.
type AuditLogFilter struct {
	SearchQuery string `json:"q,omitempty"`
}
