package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// AddOrganizationMember adds a user to an organization.
func (s *Service) AddOrganizationMember(ctx context.Context, orgID uuid.UUID, user string) (OrganizationMember, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/members/%s", orgID, user), nil)
	if err != nil {
		return OrganizationMember{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OrganizationMember{}, client.ReadBodyAsError(resp)
	}

	var member OrganizationMember
	return member, json.NewDecoder(resp.Body).Decode(&member)
}

// RemoveOrganizationMember removes a user from an organization.
func (s *Service) RemoveOrganizationMember(ctx context.Context, orgID uuid.UUID, user string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/organizations/%s/members/%s", orgID, user), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ListOrganizationMembers returns all members of an organization.
func (s *Service) ListOrganizationMembers(ctx context.Context, orgID uuid.UUID) ([]OrganizationMemberWithUserData, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/members/", orgID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var members []OrganizationMemberWithUserData
	return members, json.NewDecoder(resp.Body).Decode(&members)
}

// UpdateOrganizationMemberRoles replaces a user's roles within an organization.
// Include ALL roles the user should have in the organization.
func (s *Service) UpdateOrganizationMemberRoles(ctx context.Context, orgID uuid.UUID, user string, req UpdateRoles) (OrganizationMember, error) {
	resp, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/organizations/%s/members/%s/roles", orgID, user), req)
	if err != nil {
		return OrganizationMember{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OrganizationMember{}, client.ReadBodyAsError(resp)
	}

	var member OrganizationMember
	return member, json.NewDecoder(resp.Body).Decode(&member)
}
