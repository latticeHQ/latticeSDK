// Package users provides comprehensive user management for the Lattice Runtime API.
// It covers user CRUD, authentication, roles, organization membership, quiet hours,
// autofill parameters, and password management.
package users

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account.
type UserStatus string

const (
	// UserStatusActive indicates the user is active.
	UserStatusActive UserStatus = "active"
	// UserStatusDormant indicates the user has not yet logged in.
	UserStatusDormant UserStatus = "dormant"
	// UserStatusSuspended indicates the user has been suspended.
	UserStatusSuspended UserStatus = "suspended"
)

// LoginType identifies the authentication method for a user.
type LoginType string

const (
	// LoginTypePassword authenticates with email and password.
	LoginTypePassword LoginType = "password"
	// LoginTypeGithub authenticates via GitHub OAuth.
	LoginTypeGithub LoginType = "github"
	// LoginTypeOIDC authenticates via OpenID Connect.
	LoginTypeOIDC LoginType = "oidc"
	// LoginTypeToken authenticates via API token.
	LoginTypeToken LoginType = "token"
	// LoginTypeNone disables all login methods for the user.
	LoginTypeNone LoginType = "none"
)

// MinimalUser contains the minimum fields needed to identify a user.
type MinimalUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
}

// ReducedUser omits role and organization information from a user.
// Roles are derived from site and organization memberships, which require
// additional database lookups.
type ReducedUser struct {
	MinimalUser
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastSeenAt      time.Time  `json:"last_seen_at"`
	Status          UserStatus `json:"status"`
	LoginType       LoginType  `json:"login_type"`
	ThemePreference string     `json:"theme_preference"`
}

// User represents a full user profile in Lattice.
type User struct {
	ReducedUser
	OrganizationIDs []uuid.UUID `json:"organization_ids"`
	Roles           []SlimRole  `json:"roles"`
}

// UsersRequest specifies filters and pagination for listing users.
type UsersRequest struct {
	// Search filters users by a free-text search term.
	Search string `json:"search,omitempty"`
	// Status filters users by account status.
	Status UserStatus `json:"status,omitempty"`
	// Role filters users that have the given role.
	Role string `json:"role,omitempty"`
	// SearchQuery is passed as the "q" query parameter.
	SearchQuery string `json:"q,omitempty"`
	// Limit sets the maximum number of results per page.
	Limit int `json:"limit,omitempty"`
	// Offset is the number of results to skip.
	Offset int `json:"offset,omitempty"`
	// AfterID returns results after the given UUID (cursor-based pagination).
	AfterID uuid.UUID `json:"after_id,omitempty"`
}

// GetUsersResponse is the API response for listing users.
type GetUsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

// CreateFirstUserRequest contains the parameters for creating the initial
// superadmin user on a new Lattice deployment.
type CreateFirstUserRequest struct {
	Email     string                   `json:"email"`
	Username  string                   `json:"username"`
	Name      string                   `json:"name"`
	Password  string                   `json:"password"`
	Trial     bool                     `json:"trial"`
	TrialInfo CreateFirstUserTrialInfo `json:"trial_info"`
}

// CreateFirstUserTrialInfo contains optional trial registration details.
type CreateFirstUserTrialInfo struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	JobTitle    string `json:"job_title"`
	CompanyName string `json:"company_name"`
	Country     string `json:"country"`
	Developers  string `json:"developers"`
}

// CreateFirstUserResponse contains the IDs for a newly created first user.
type CreateFirstUserResponse struct {
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
}

// CreateUserRequest contains the parameters for creating a user.
//
// Deprecated: Use CreateUserRequestWithOrgs instead.
type CreateUserRequest struct {
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	Password       string    `json:"password"`
	UserLoginType  LoginType `json:"login_type"`
	DisableLogin   bool      `json:"disable_login"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
}

// CreateUserRequestWithOrgs creates a user with membership in multiple organizations.
type CreateUserRequestWithOrgs struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
	// UserLoginType defaults to LoginTypePassword.
	UserLoginType LoginType `json:"login_type"`
	// UserStatus defaults to UserStatusDormant.
	UserStatus *UserStatus `json:"user_status,omitempty"`
	// OrganizationIDs lists organizations the user should be a member of.
	OrganizationIDs []uuid.UUID `json:"organization_ids,omitempty"`
}

// UpdateUserProfileRequest contains updatable profile fields.
type UpdateUserProfileRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

// UpdateUserPasswordRequest changes a user's password.
type UpdateUserPasswordRequest struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}

// ValidateUserPasswordRequest checks password complexity.
type ValidateUserPasswordRequest struct {
	Password string `json:"password"`
}

// ValidateUserPasswordResponse indicates whether a password meets complexity requirements.
type ValidateUserPasswordResponse struct {
	Valid   bool   `json:"valid"`
	Details string `json:"details"`
}

// UpdateUserAppearanceSettingsRequest updates a user's theme preference.
type UpdateUserAppearanceSettingsRequest struct {
	ThemePreference string `json:"theme_preference"`
}

// UpdateRoles sets the complete list of role names for a user.
type UpdateRoles struct {
	Roles []string `json:"roles"`
}

// UserRoles contains all roles assigned to a user, both site-wide and per-organization.
type UserRoles struct {
	Roles             []string               `json:"roles"`
	OrganizationRoles map[uuid.UUID][]string `json:"organization_roles"`
}

// LoginWithPasswordRequest authenticates a user with email and password.
type LoginWithPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginWithPasswordResponse contains the session token from a successful password login.
type LoginWithPasswordResponse struct {
	SessionToken string `json:"session_token"`
}

// RequestOneTimePasscodeRequest triggers a one-time passcode email for password reset.
type RequestOneTimePasscodeRequest struct {
	Email string `json:"email"`
}

// ChangePasswordWithOneTimePasscodeRequest resets a password using a one-time passcode.
type ChangePasswordWithOneTimePasscodeRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	OneTimePasscode string `json:"one_time_passcode"`
}

// ConvertLoginRequest converts a user's authentication method (e.g. password to OAuth).
type ConvertLoginRequest struct {
	// ToType is the target login type.
	ToType   LoginType `json:"to_type"`
	Password string    `json:"password"`
}

// OAuthConversionResponse contains the OAuth state needed to complete a login type conversion.
type OAuthConversionResponse struct {
	StateString string    `json:"state_string"`
	ExpiresAt   time.Time `json:"expires_at"`
	ToType      LoginType `json:"to_type"`
	UserID      uuid.UUID `json:"user_id"`
}

// AuthMethods describes available authentication methods on the deployment.
type AuthMethods struct {
	TermsOfServiceURL string         `json:"terms_of_service_url,omitempty"`
	Password          AuthMethod     `json:"password"`
	Github            AuthMethod     `json:"github"`
	OIDC              OIDCAuthMethod `json:"oidc"`
}

// AuthMethod indicates whether an authentication method is enabled.
type AuthMethod struct {
	Enabled bool `json:"enabled"`
}

// OIDCAuthMethod extends AuthMethod with OIDC-specific display configuration.
type OIDCAuthMethod struct {
	AuthMethod
	SignInText string `json:"signInText"`
	IconURL    string `json:"iconUrl"`
}

// UserQuietHoursScheduleResponse contains quiet hours settings for a user.
type UserQuietHoursScheduleResponse struct {
	// RawSchedule is the raw cron expression.
	RawSchedule string `json:"raw_schedule"`
	// UserSet is true if the user has configured a custom quiet hours schedule.
	UserSet bool `json:"user_set"`
	// UserCanSet is true if the user is allowed to configure their own quiet hours.
	UserCanSet bool `json:"user_can_set"`
	// Time is the start time of the quiet hours window (HH:mm, 24-hour).
	Time string `json:"time"`
	// Timezone is the timezone for the schedule (UTC if unspecified).
	Timezone string `json:"timezone"`
	// Next is the next time the quiet hours window will start.
	Next time.Time `json:"next"`
}

// UpdateUserQuietHoursScheduleRequest updates a user's quiet hours schedule.
type UpdateUserQuietHoursScheduleRequest struct {
	// Schedule is a cron expression defining the start of the quiet hours window.
	// Must be daily with a single time. Use a CRON_TZ prefix for timezone
	// (e.g. "CRON_TZ=America/New_York 0 2 * * *"). An empty value resets to the
	// deployment default.
	Schedule string `json:"schedule"`
}

// UserParameter represents an autofill parameter for template builds.
type UserParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// OrganizationMember represents a user's membership in an organization.
type OrganizationMember struct {
	UserID         uuid.UUID  `json:"user_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Roles          []SlimRole `json:"roles"`
}

// OrganizationMemberWithUserData extends OrganizationMember with user profile fields.
type OrganizationMemberWithUserData struct {
	Username    string     `json:"username"`
	Name        string     `json:"name"`
	AvatarURL   string     `json:"avatar_url"`
	Email       string     `json:"email"`
	GlobalRoles []SlimRole `json:"global_roles"`
	OrganizationMember
}

// Organization represents an organization in Lattice.
type Organization struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	DisplayName  string    `json:"display_name"`
	Icon         string    `json:"icon"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsDefault    bool      `json:"is_default"`
	DeploymentID uuid.UUID `json:"deployment_id"`
}

// --- Role types ---

// SlimRole is a lightweight role reference without permission details.
type SlimRole struct {
	Name           string `json:"name"`
	DisplayName    string `json:"display_name"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// Role represents a full RBAC role definition with permissions.
type Role struct {
	Name                    string       `json:"name"`
	OrganizationID          string       `json:"organization_id,omitempty"`
	DisplayName             string       `json:"display_name"`
	SitePermissions         []Permission `json:"site_permissions"`
	OrganizationPermissions []Permission `json:"organization_permissions"`
	UserPermissions         []Permission `json:"user_permissions"`
}

// AssignableRoles wraps a Role with metadata about whether it can be assigned.
type AssignableRoles struct {
	Role
	// Assignable is true if the current user can assign this role.
	Assignable bool `json:"assignable"`
	// BuiltIn is true if this is a system-defined immutable role.
	BuiltIn bool `json:"built_in"`
}

// Permission defines a single RBAC permission entry.
type Permission struct {
	// Negate makes this a negative (deny) permission.
	Negate       bool         `json:"negate"`
	ResourceType RBACResource `json:"resource_type"`
	Action       RBACAction   `json:"action"`
}

// RBACResource identifies a resource type for authorization checks.
type RBACResource string

const (
	RBACResourceAgent        RBACResource = "agent"
	RBACResourceTemplate     RBACResource = "template"
	RBACResourceUser         RBACResource = "user"
	RBACResourceOrganization RBACResource = "organization"
	RBACResourceAuditLog     RBACResource = "audit_log"
	RBACResourceGroup        RBACResource = "group"
	RBACResourceAPIKey       RBACResource = "api_key"
	RBACResourceLicense      RBACResource = "license"
	RBACResourceCustomRole   RBACResource = "custom_role"
)

// RBACAction identifies an action for authorization checks.
type RBACAction string

const (
	RBACActionCreate RBACAction = "create"
	RBACActionRead   RBACAction = "read"
	RBACActionUpdate RBACAction = "update"
	RBACActionDelete RBACAction = "delete"
)

// CreatePermissions is a helper that builds a Permission slice from a resource-to-actions mapping.
func CreatePermissions(mapping map[RBACResource][]RBACAction) []Permission {
	perms := make([]Permission, 0)
	for t, actions := range mapping {
		for _, action := range actions {
			perms = append(perms, Permission{
				ResourceType: t,
				Action:       action,
			})
		}
	}
	return perms
}
