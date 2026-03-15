// Package groups provides group management for the Lattice Runtime API.
package groups

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides group management operations.
type Service struct {
	client *client.Client
}

// New creates a new groups Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// CreateGroup creates a new group within the specified organization.
func (s *Service) CreateGroup(ctx context.Context, orgID uuid.UUID, req CreateGroupRequest) (Group, error) {
	res, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/groups", orgID),
		req,
	)
	if err != nil {
		return Group{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return Group{}, client.ReadBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// ListGroupsByOrg returns all groups belonging to the given organization.
//
// Deprecated: Use ListGroups with GroupArguments instead.
func (s *Service) ListGroupsByOrg(ctx context.Context, orgID uuid.UUID) ([]Group, error) {
	return s.ListGroups(ctx, GroupArguments{Organization: orgID.String()})
}

// ListGroups returns groups matching the supplied filter arguments.
func (s *Service) ListGroups(ctx context.Context, args GroupArguments) ([]Group, error) {
	qp := url.Values{}
	if args.Organization != "" {
		qp.Set("organization", args.Organization)
	}
	if args.HasMember != "" {
		qp.Set("has_member", args.HasMember)
	}
	if len(args.GroupIDs) > 0 {
		idStrs := make([]string, 0, len(args.GroupIDs))
		for _, id := range args.GroupIDs {
			idStrs = append(idStrs, id.String())
		}
		qp.Set("group_ids", strings.Join(idStrs, ","))
	}

	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/groups?%s", qp.Encode()),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var groups []Group
	return groups, json.NewDecoder(res.Body).Decode(&groups)
}

// GetGroupByOrgAndName retrieves a group by its organization ID and name.
func (s *Service) GetGroupByOrgAndName(ctx context.Context, orgID uuid.UUID, name string) (Group, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/groups/%s", orgID, name),
		nil,
	)
	if err != nil {
		return Group{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Group{}, client.ReadBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetGroup retrieves a single group by its ID.
func (s *Service) GetGroup(ctx context.Context, id uuid.UUID) (Group, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/groups/%s", id),
		nil,
	)
	if err != nil {
		return Group{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Group{}, client.ReadBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// UpdateGroup patches an existing group.
func (s *Service) UpdateGroup(ctx context.Context, id uuid.UUID, req PatchGroupRequest) (Group, error) {
	res, err := s.client.Request(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/groups/%s", id),
		req,
	)
	if err != nil {
		return Group{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Group{}, client.ReadBodyAsError(res)
	}
	var resp Group
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// DeleteGroup deletes a group by its ID.
func (s *Service) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	res, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/groups/%s", id),
		nil,
	)
	if err != nil {
		return fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(res)
	}
	return nil
}
