package budget

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// CostSummary contains aggregated cost information.
type CostSummary struct {
	TotalCredits     int `json:"total_credits"`
	UsedCredits      int `json:"used_credits"`
	RemainingCredits int `json:"remaining_credits"`
}

// GetCostSummary returns the cost summary for the deployment.
func (s *Service) GetCostSummary(ctx context.Context) (CostSummary, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/deployment/cost", nil)
	if err != nil {
		return CostSummary{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CostSummary{}, client.ReadBodyAsError(resp)
	}

	var summary CostSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		return CostSummary{}, fmt.Errorf("decode response: %w", err)
	}
	return summary, nil
}
