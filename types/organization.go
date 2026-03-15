package types

import (
	"time"

	"github.com/google/uuid"
)

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

// MinimalOrganization contains the minimum organization identification fields.
type MinimalOrganization struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Icon        string    `json:"icon"`
}

// OrganizationMember represents a user's membership in an organization.
type OrganizationMember struct {
	UserID         uuid.UUID  `json:"user_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Roles          []SlimRole `json:"roles"`
}

// CreateOrganizationRequest contains the parameters for creating an organization.
type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name,omitempty"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// UpdateOrganizationRequest contains updatable organization fields.
type UpdateOrganizationRequest struct {
	Name        string  `json:"name,omitempty"`
	DisplayName string  `json:"display_name,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}
