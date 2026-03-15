package authz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// Can checks whether the authenticated user is allowed to perform an action
// on a resource. This is the simplest authorization check.
func (s *Service) Can(ctx context.Context, action types.RBACAction, resourceType types.RBACResource, resourceID string) (bool, error) {
	checks := types.AuthorizationRequest{
		Checks: map[string]types.AuthorizationCheck{
			"check": {
				Object: types.AuthorizationObject{
					ResourceType: resourceType,
					ResourceID:   resourceID,
				},
				Action: action,
			},
		},
	}

	result, err := s.Check(ctx, checks)
	if err != nil {
		return false, err
	}
	return result["check"], nil
}

// Check performs a batch authorization check. Multiple checks can be performed
// in a single request. Returns a map of check names to allow/deny results.
func (s *Service) Check(ctx context.Context, req types.AuthorizationRequest) (types.AuthorizationResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/authcheck", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var result types.AuthorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
