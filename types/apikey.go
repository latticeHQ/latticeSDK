package types

import (
	"time"

	"github.com/google/uuid"
)

// APIKeyScope defines the scope of an API key.
type APIKeyScope string

const (
	APIKeyScopeAll                APIKeyScope = "all"
	APIKeyScopeApplicationConnect APIKeyScope = "application_connect"
)

// APIKey represents an API key for authentication.
type APIKey struct {
	ID              string      `json:"id"`
	UserID          uuid.UUID   `json:"user_id"`
	LastUsed        time.Time   `json:"last_used"`
	ExpiresAt       time.Time   `json:"expires_at"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	LoginType       LoginType   `json:"login_type"`
	Scope           APIKeyScope `json:"scope"`
	TokenName       string      `json:"token_name"`
	LifetimeSeconds int64       `json:"lifetime_seconds"`
}

// CreateTokenRequest contains the parameters for creating a new token.
type CreateTokenRequest struct {
	Lifetime  time.Duration `json:"lifetime"`
	Scope     APIKeyScope   `json:"scope"`
	TokenName string        `json:"token_name"`
}

// GenerateAPIKeyResponse contains the generated API key.
type GenerateAPIKeyResponse struct {
	Key string `json:"key"`
}

// TokensFilter filters token list requests.
type TokensFilter struct {
	IncludeAll bool `json:"include_all"`
}

// TokenConfig contains token configuration limits.
type TokenConfig struct {
	MaxTokenLifetime time.Duration `json:"max_token_lifetime"`
}
