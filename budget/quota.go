package budget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// GetQuota returns the budget quota for a member within an organization.
func (s *Service) GetQuota(ctx context.Context, orgID, userID uuid.UUID) (types.AgentQuota, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s/members/%s/agent-quota", orgID, userID), nil)
	if err != nil {
		return types.AgentQuota{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.AgentQuota{}, client.ReadBodyAsError(resp)
	}

	var quota types.AgentQuota
	if err := json.NewDecoder(resp.Body).Decode(&quota); err != nil {
		return types.AgentQuota{}, fmt.Errorf("decode response: %w", err)
	}
	return quota, nil
}

// CheckBudget returns true if the member has remaining budget within the organization.
func (s *Service) CheckBudget(ctx context.Context, orgID, userID uuid.UUID) (bool, error) {
	quota, err := s.GetQuota(ctx, orgID, userID)
	if err != nil {
		return false, err
	}
	// Budget of 0 means unlimited.
	if quota.Budget == 0 {
		return true, nil
	}
	return quota.CreditsConsumed < quota.Budget, nil
}
