package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// GetOrganization returns an organization by ID.
func (s *Service) GetOrganization(ctx context.Context, id uuid.UUID) (types.Organization, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s", id), nil)
	if err != nil {
		return types.Organization{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Organization{}, client.ReadBodyAsError(resp)
	}

	var org types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return types.Organization{}, fmt.Errorf("decode response: %w", err)
	}
	return org, nil
}

// ListOrganizations returns all organizations visible to the caller.
func (s *Service) ListOrganizations(ctx context.Context) ([]types.Organization, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/organizations", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var orgs []types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return orgs, nil
}

// ListOrganizationMembers returns the members of an organization.
func (s *Service) ListOrganizationMembers(ctx context.Context, orgID uuid.UUID, p *client.Pagination) ([]types.OrganizationMember, error) {
	var opts []client.RequestOption
	if p != nil {
		opts = append(opts, p.AsRequestOption())
	}

	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s/members", orgID), nil, opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var members []types.OrganizationMember
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return members, nil
}

// CreateOrganization creates a new organization.
func (s *Service) CreateOrganization(ctx context.Context, req types.CreateOrganizationRequest) (types.Organization, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/organizations", req)
	if err != nil {
		return types.Organization{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return types.Organization{}, client.ReadBodyAsError(resp)
	}

	var org types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return types.Organization{}, fmt.Errorf("decode response: %w", err)
	}
	return org, nil
}

// UpdateOrganization updates an organization.
func (s *Service) UpdateOrganization(ctx context.Context, id uuid.UUID, req types.UpdateOrganizationRequest) (types.Organization, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/organizations/%s", id), req)
	if err != nil {
		return types.Organization{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Organization{}, client.ReadBodyAsError(resp)
	}

	var org types.Organization
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return types.Organization{}, fmt.Errorf("decode response: %w", err)
	}
	return org, nil
}

// DeleteOrganization deletes an organization.
func (s *Service) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/organizations/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
