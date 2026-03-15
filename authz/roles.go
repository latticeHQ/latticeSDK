package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// ListSiteRoles returns all site-wide roles.
func (s *Service) ListSiteRoles(ctx context.Context) ([]types.AssignableRoles, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/roles", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var roles []types.AssignableRoles
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return roles, nil
}

// ListOrganizationRoles returns roles for an organization.
func (s *Service) ListOrganizationRoles(ctx context.Context, orgID uuid.UUID) ([]types.AssignableRoles, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s/members/roles", orgID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var roles []types.AssignableRoles
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return roles, nil
}
