package types

// RBACResource identifies a resource type for authorization checks.
type RBACResource string

const (
	RBACResourceAgent           RBACResource = "agent"
	RBACResourceTemplate        RBACResource = "template"
	RBACResourceUser            RBACResource = "user"
	RBACResourceOrganization    RBACResource = "organization"
	RBACResourceAuditLog        RBACResource = "audit_log"
	RBACResourceGroup           RBACResource = "group"
	RBACResourceAPIKey          RBACResource = "api_key"
	RBACResourceLicense         RBACResource = "license"
	RBACResourceCustomRole      RBACResource = "custom_role"
)

// RBACAction identifies an action for authorization checks.
type RBACAction string

const (
	RBACActionCreate RBACAction = "create"
	RBACActionRead   RBACAction = "read"
	RBACActionUpdate RBACAction = "update"
	RBACActionDelete RBACAction = "delete"
)

// Well-known role names.
const (
	RoleOwner                     = "owner"
	RoleMember                    = "member"
	RoleTemplateAdmin             = "template-admin"
	RoleUserAdmin                 = "user-admin"
	RoleAuditor                   = "auditor"
	RoleOrganizationAdmin         = "organization-admin"
	RoleOrganizationMember        = "organization-member"
	RoleOrganizationAuditor       = "organization-auditor"
	RoleOrganizationTemplateAdmin = "organization-template-admin"
	RoleOrganizationUserAdmin     = "organization-user-admin"
)

// SlimRole is a lightweight role reference.
type SlimRole struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// Role represents a full RBAC role definition.
type Role struct {
	Name                    string       `json:"name"`
	OrganizationID          string       `json:"organization_id,omitempty"`
	DisplayName             string       `json:"display_name"`
	SitePermissions         []Permission `json:"site_permissions"`
	OrganizationPermissions []Permission `json:"organization_permissions"`
	UserPermissions         []Permission `json:"user_permissions"`
}

// Permission defines a single RBAC permission.
type Permission struct {
	Negate       bool         `json:"negate"`
	ResourceType RBACResource `json:"resource_type"`
	Action       RBACAction   `json:"action"`
}

// AssignableRoles wraps a Role with assignment metadata.
type AssignableRoles struct {
	Role
	Assignable bool `json:"assignable"`
	BuiltIn    bool `json:"built_in"`
}

// AuthorizationCheck is a single authorization check.
type AuthorizationCheck struct {
	Object AuthorizationObject `json:"object"`
	Action RBACAction          `json:"action"`
}

// AuthorizationObject identifies the target of an authorization check.
type AuthorizationObject struct {
	ResourceType   RBACResource `json:"resource_type"`
	OwnerID        string       `json:"owner_id,omitempty"`
	OrganizationID string       `json:"organization_id,omitempty"`
	ResourceID     string       `json:"resource_id,omitempty"`
	AnyOrgOwner    bool         `json:"any_org,omitempty"`
}

// AuthorizationRequest is a batch authorization check request.
type AuthorizationRequest struct {
	Checks map[string]AuthorizationCheck `json:"checks"`
}

// AuthorizationResponse maps check keys to allow/deny results.
type AuthorizationResponse map[string]bool
