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

// GetQuota returns the budget quota for an owner (user or organization).
func (s *Service) GetQuota(ctx context.Context, ownerID uuid.UUID) (types.AgentQuota, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/quota/%s", ownerID), nil)
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

// CheckBudget returns true if the owner has remaining budget.
func (s *Service) CheckBudget(ctx context.Context, ownerID uuid.UUID) (bool, error) {
	quota, err := s.GetQuota(ctx, ownerID)
	if err != nil {
		return false, err
	}
	// Budget of 0 means unlimited.
	if quota.Budget == 0 {
		return true, nil
	}
	return quota.CreditsConsumed < quota.Budget, nil
}
