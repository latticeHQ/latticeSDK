// Package oauth2 provides OAuth2 provider application management for the Lattice Runtime API.
package oauth2

import (
	"time"

	"github.com/google/uuid"
)

// OAuth2AppEndpoints contains the OAuth2 endpoint URLs for an application.
type OAuth2AppEndpoints struct {
	Authorization string `json:"authorization"`
	Token         string `json:"token"`
	// DeviceAuth is the optional device authorization endpoint.
	DeviceAuth string `json:"device_authorization"`
}

// OAuth2ProviderApp represents an application configured to authenticate
// using Lattice as an OAuth2 provider.
type OAuth2ProviderApp struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	CallbackURL string             `json:"callback_url"`
	Icon        string             `json:"icon"`
	Endpoints   OAuth2AppEndpoints `json:"endpoints"`
}

// OAuth2ProviderAppFilter filters the list of OAuth2 provider applications.
type OAuth2ProviderAppFilter struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

// PostOAuth2ProviderAppRequest is the request body for creating a new
// OAuth2 provider application.
type PostOAuth2ProviderAppRequest struct {
	Name        string `json:"name"`
	CallbackURL string `json:"callback_url"`
	Icon        string `json:"icon,omitempty"`
}

// PutOAuth2ProviderAppRequest is the request body for updating an
// OAuth2 provider application.
type PutOAuth2ProviderAppRequest struct {
	Name        string `json:"name"`
	CallbackURL string `json:"callback_url"`
	Icon        string `json:"icon,omitempty"`
}

// OAuth2ProviderAppSecret represents a truncated secret for an OAuth2 application.
type OAuth2ProviderAppSecret struct {
	ID                    uuid.UUID  `json:"id"`
	LastUsedAt            *time.Time `json:"last_used_at"`
	ClientSecretTruncated string     `json:"client_secret_truncated"`
}

// OAuth2ProviderAppSecretFull contains the full secret, only returned at
// creation time.
type OAuth2ProviderAppSecretFull struct {
	ID               uuid.UUID `json:"id"`
	ClientSecretFull string    `json:"client_secret_full"`
}

// OAuth2ProviderGrantType represents a supported OAuth2 grant type.
type OAuth2ProviderGrantType string

const (
	// OAuth2ProviderGrantTypeAuthorizationCode is the authorization code grant type.
	OAuth2ProviderGrantTypeAuthorizationCode OAuth2ProviderGrantType = "authorization_code"
	// OAuth2ProviderGrantTypeRefreshToken is the refresh token grant type.
	OAuth2ProviderGrantTypeRefreshToken OAuth2ProviderGrantType = "refresh_token"
)

// Valid reports whether the grant type is a known value.
func (g OAuth2ProviderGrantType) Valid() bool {
	switch g {
	case OAuth2ProviderGrantTypeAuthorizationCode, OAuth2ProviderGrantTypeRefreshToken:
		return true
	}
	return false
}

// OAuth2ProviderResponseType represents a supported OAuth2 response type.
type OAuth2ProviderResponseType string

const (
	// OAuth2ProviderResponseTypeCode is the authorization code response type.
	OAuth2ProviderResponseTypeCode OAuth2ProviderResponseType = "code"
)

// Valid reports whether the response type is a known value.
func (r OAuth2ProviderResponseType) Valid() bool {
	switch r {
	case OAuth2ProviderResponseTypeCode:
		return true
	}
	return false
}
