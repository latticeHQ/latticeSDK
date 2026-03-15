package types

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusDormant   UserStatus = "dormant"
	UserStatusSuspended UserStatus = "suspended"
)

// LoginType identifies the authentication method.
type LoginType string

const (
	LoginTypePassword LoginType = "password"
	LoginTypeGithub   LoginType = "github"
	LoginTypeOIDC     LoginType = "oidc"
	LoginTypeToken    LoginType = "token"
	LoginTypeNone     LoginType = "none"
)

// MinimalUser contains the minimum user identification fields.
type MinimalUser struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
}

// User represents a full user profile.
type User struct {
	ID              uuid.UUID   `json:"id"`
	Username        string      `json:"username"`
	AvatarURL       string      `json:"avatar_url"`
	Name            string      `json:"name"`
	Email           string      `json:"email"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	LastSeenAt      time.Time   `json:"last_seen_at"`
	Status          UserStatus  `json:"status"`
	LoginType       LoginType   `json:"login_type"`
	ThemePreference string      `json:"theme_preference"`
	OrganizationIDs []uuid.UUID `json:"organization_ids"`
	Roles           []SlimRole  `json:"roles"`
}

// CreateUserRequest contains the parameters for creating a new user.
type CreateUserRequest struct {
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Name           string    `json:"name"`
	Password       string    `json:"password"`
	UserLoginType  LoginType `json:"login_type"`
	DisableLogin   bool      `json:"disable_login"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
}

// UpdateUserProfileRequest contains updatable user profile fields.
type UpdateUserProfileRequest struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

// UpdateRoles sets a user's roles.
type UpdateRoles struct {
	Roles []string `json:"roles"`
}

// LoginWithPasswordRequest authenticates with email and password.
type LoginWithPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginWithPasswordResponse contains the session token from a successful login.
type LoginWithPasswordResponse struct {
	SessionToken string `json:"session_token"`
}

// AuthMethods describes available authentication methods.
type AuthMethods struct {
	TermsOfServiceURL string         `json:"terms_of_service_url,omitempty"`
	Password          AuthMethod     `json:"password"`
	Github            AuthMethod     `json:"github"`
	OIDC              OIDCAuthMethod `json:"oidc"`
}

// AuthMethod indicates whether an auth method is enabled.
type AuthMethod struct {
	Enabled bool `json:"enabled"`
}

// OIDCAuthMethod extends AuthMethod with OIDC-specific display configuration.
type OIDCAuthMethod struct {
	AuthMethod
	SignInText string `json:"signInText"`
	IconURL    string `json:"iconUrl"`
}
