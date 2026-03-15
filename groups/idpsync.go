package groups

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// GetGroupIDPSyncSettings retrieves the IDP sync settings for groups within
// the specified organization.
func (s *Service) GetGroupIDPSyncSettings(ctx context.Context, orgID string) (GroupSyncSettings, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/settings/idpsync/groups", orgID),
		nil,
	)
	if err != nil {
		return GroupSyncSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return GroupSyncSettings{}, client.ReadBodyAsError(res)
	}
	var resp GroupSyncSettings
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// UpdateGroupIDPSyncSettings updates the IDP sync settings for groups within
// the specified organization.
func (s *Service) UpdateGroupIDPSyncSettings(ctx context.Context, orgID string, req GroupSyncSettings) (GroupSyncSettings, error) {
	res, err := s.client.Request(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/organizations/%s/settings/idpsync/groups", orgID),
		req,
	)
	if err != nil {
		return GroupSyncSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return GroupSyncSettings{}, client.ReadBodyAsError(res)
	}
	var resp GroupSyncSettings
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetRoleIDPSyncSettings retrieves the IDP sync settings for roles within
// the specified organization.
func (s *Service) GetRoleIDPSyncSettings(ctx context.Context, orgID string) (RoleSyncSettings, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/settings/idpsync/roles", orgID),
		nil,
	)
	if err != nil {
		return RoleSyncSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return RoleSyncSettings{}, client.ReadBodyAsError(res)
	}
	var resp RoleSyncSettings
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// UpdateRoleIDPSyncSettings updates the IDP sync settings for roles within
// the specified organization.
func (s *Service) UpdateRoleIDPSyncSettings(ctx context.Context, orgID string, req RoleSyncSettings) (RoleSyncSettings, error) {
	res, err := s.client.Request(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/organizations/%s/settings/idpsync/roles", orgID),
		req,
	)
	if err != nil {
		return RoleSyncSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return RoleSyncSettings{}, client.ReadBodyAsError(res)
	}
	var resp RoleSyncSettings
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetOrganizationIDPSyncSettings retrieves the deployment-wide IDP sync
// settings for organization membership.
func (s *Service) GetOrganizationIDPSyncSettings(ctx context.Context) (OrganizationSyncSettings, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		"/api/v2/settings/idpsync/organization",
		nil,
	)
	if err != nil {
		return OrganizationSyncSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return OrganizationSyncSettings{}, client.ReadBodyAsError(res)
	}
	var resp OrganizationSyncSettings
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// UpdateOrganizationIDPSyncSettings updates the deployment-wide IDP sync
// settings for organization membership.
func (s *Service) UpdateOrganizationIDPSyncSettings(ctx context.Context, req OrganizationSyncSettings) (OrganizationSyncSettings, error) {
	res, err := s.client.Request(ctx, http.MethodPatch,
		"/api/v2/settings/idpsync/organization",
		req,
	)
	if err != nil {
		return OrganizationSyncSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return OrganizationSyncSettings{}, client.ReadBodyAsError(res)
	}
	var resp OrganizationSyncSettings
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetAvailableIDPSyncFields returns the claim fields available for IDP sync
// at the deployment level.
func (s *Service) GetAvailableIDPSyncFields(ctx context.Context) ([]string, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		"/api/v2/settings/idpsync/available-fields",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var resp []string
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetOrgAvailableIDPSyncFields returns the claim fields available for IDP sync
// within the specified organization.
func (s *Service) GetOrgAvailableIDPSyncFields(ctx context.Context, orgID string) ([]string, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/settings/idpsync/available-fields", orgID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var resp []string
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}
