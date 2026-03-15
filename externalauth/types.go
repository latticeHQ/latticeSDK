// Package externalauth provides external authentication provider management
// for the Lattice Runtime API.
package externalauth

import "time"

// EnhancedExternalAuthProvider represents enhanced support for a type of
// external authentication. Git providers are examples of enhanced providers
// because they support intercepting "git clone".
type EnhancedExternalAuthProvider string

// String returns the string representation of the provider.
func (e EnhancedExternalAuthProvider) String() string {
	return string(e)
}

// Git returns whether the provider is a Git provider.
func (e EnhancedExternalAuthProvider) Git() bool {
	switch e {
	case EnhancedExternalAuthProviderGitHub,
		EnhancedExternalAuthProviderGitLab,
		EnhancedExternalAuthProviderBitBucketCloud,
		EnhancedExternalAuthProviderBitBucketServer,
		EnhancedExternalAuthProviderAzureDevops,
		EnhancedExternalAuthProviderAzureDevopsEntra,
		EnhancedExternalAuthProviderGitea:
		return true
	default:
		return false
	}
}

const (
	EnhancedExternalAuthProviderAzureDevops      EnhancedExternalAuthProvider = "azure-devops"
	EnhancedExternalAuthProviderAzureDevopsEntra  EnhancedExternalAuthProvider = "azure-devops-entra"
	EnhancedExternalAuthProviderGitHub            EnhancedExternalAuthProvider = "github"
	EnhancedExternalAuthProviderGitLab            EnhancedExternalAuthProvider = "gitlab"
	EnhancedExternalAuthProviderBitBucketCloud    EnhancedExternalAuthProvider = "bitbucket-cloud"
	EnhancedExternalAuthProviderBitBucketServer   EnhancedExternalAuthProvider = "bitbucket-server"
	EnhancedExternalAuthProviderSlack             EnhancedExternalAuthProvider = "slack"
	EnhancedExternalAuthProviderJFrog             EnhancedExternalAuthProvider = "jfrog"
	EnhancedExternalAuthProviderGitea             EnhancedExternalAuthProvider = "gitea"
)

// ExternalAuth contains the authentication state for an external provider.
type ExternalAuth struct {
	Authenticated    bool                         `json:"authenticated"`
	Device           bool                         `json:"device"`
	DisplayName      string                       `json:"display_name"`
	User             *ExternalAuthUser             `json:"user"`
	AppInstallable   bool                         `json:"app_installable"`
	AppInstallations []ExternalAuthAppInstallation `json:"installations"`
	AppInstallURL    string                       `json:"app_install_url"`
}

// ExternalAuthUser represents the user who authenticated with an external provider.
type ExternalAuthUser struct {
	ID         int64  `json:"id"`
	Login      string `json:"login"`
	AvatarURL  string `json:"avatar_url"`
	ProfileURL string `json:"profile_url"`
	Name       string `json:"name"`
}

// ExternalAuthAppInstallation represents an installed external auth application.
type ExternalAuthAppInstallation struct {
	ID           int              `json:"id"`
	Account      ExternalAuthUser `json:"account"`
	ConfigureURL string           `json:"configure_url"`
}

// ExternalAuthDevice is the response from the device authorization endpoint.
// See: https://tools.ietf.org/html/rfc8628#section-3.2
type ExternalAuthDevice struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// ExternalAuthDeviceExchange is the request body to exchange a device code for a token.
type ExternalAuthDeviceExchange struct {
	DeviceCode string `json:"device_code"`
}

// ExternalAuthLink is a link between a user and an external auth provider.
type ExternalAuthLink struct {
	ProviderID      string    `json:"provider_id"`
	CreatedAt       time.Time `json:"created_at" format:"date-time"`
	UpdatedAt       time.Time `json:"updated_at" format:"date-time"`
	HasRefreshToken bool      `json:"has_refresh_token"`
	Expires         time.Time `json:"expires" format:"date-time"`
	Authenticated   bool      `json:"authenticated"`
	ValidateError   string    `json:"validate_error"`
}

// ExternalAuthLinkProvider contains the static details of an external auth provider.
type ExternalAuthLinkProvider struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Device        bool   `json:"device"`
	DisplayName   string `json:"display_name"`
	DisplayIcon   string `json:"display_icon"`
	AllowRefresh  bool   `json:"allow_refresh"`
	AllowValidate bool   `json:"allow_validate"`
}

// ListUserExternalAuthResponse is the response listing external auth providers
// and the user's authenticated links.
type ListUserExternalAuthResponse struct {
	Providers []ExternalAuthLinkProvider `json:"providers"`
	Links     []ExternalAuthLink        `json:"links"`
}
