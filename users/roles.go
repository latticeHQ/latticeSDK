package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// customRoleRequest is the wire format for creating and updating custom roles.
type customRoleRequest struct {
	Name                    string       `json:"name"`
	DisplayName             string       `json:"display_name"`
	SitePermissions         []Permission `json:"site_permissions"`
	OrganizationPermissions []Permission `json:"organization_permissions"`
	UserPermissions         []Permission `json:"user_permissions"`
}

// ListSiteRoles returns all assignable site-wide roles.
func (s *Service) ListSiteRoles(ctx context.Context) ([]AssignableRoles, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/roles", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var roles []AssignableRoles
	return roles, json.NewDecoder(resp.Body).Decode(&roles)
}

// ListOrganizationRoles returns all assignable roles for an organization.
func (s *Service) ListOrganizationRoles(ctx context.Context, orgID uuid.UUID) ([]AssignableRoles, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/members/roles", orgID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var roles []AssignableRoles
	return roles, json.NewDecoder(resp.Body).Decode(&roles)
}

// CreateOrganizationRole creates a custom role within an organization.
// The role's OrganizationID field determines which organization the role belongs to.
func (s *Service) CreateOrganizationRole(ctx context.Context, role Role) (Role, error) {
	req := customRoleRequest{
		Name:                    role.Name,
		DisplayName:             role.DisplayName,
		SitePermissions:         role.SitePermissions,
		OrganizationPermissions: role.OrganizationPermissions,
		UserPermissions:         role.UserPermissions,
	}

	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/members/roles", role.OrganizationID), req)
	if err != nil {
		return Role{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Role{}, client.ReadBodyAsError(resp)
	}

	var result Role
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// UpdateOrganizationRole updates an existing custom role within an organization.
// The role's OrganizationID field determines which organization the role belongs to.
func (s *Service) UpdateOrganizationRole(ctx context.Context, role Role) (Role, error) {
	req := customRoleRequest{
		Name:                    role.Name,
		DisplayName:             role.DisplayName,
		SitePermissions:         role.SitePermissions,
		OrganizationPermissions: role.OrganizationPermissions,
		UserPermissions:         role.UserPermissions,
	}

	resp, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/organizations/%s/members/roles", role.OrganizationID), req)
	if err != nil {
		return Role{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Role{}, client.ReadBodyAsError(resp)
	}

	var result Role
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// DeleteOrganizationRole deletes a custom role from an organization.
func (s *Service) DeleteOrganizationRole(ctx context.Context, orgID uuid.UUID, roleName string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/organizations/%s/members/roles/%s", orgID, roleName), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
