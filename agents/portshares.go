package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// GetPortShares returns all port shares configured for an agent.
func (s *Service) GetPortShares(ctx context.Context, agentID uuid.UUID) (PortSharesResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s/port-share", agentID), nil)
	if err != nil {
		return PortSharesResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PortSharesResponse{}, client.ReadBodyAsError(resp)
	}

	var shares PortSharesResponse
	if err := json.NewDecoder(resp.Body).Decode(&shares); err != nil {
		return PortSharesResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return shares, nil
}

// UpsertPortShare creates or updates a port share for an agent.
func (s *Service) UpsertPortShare(ctx context.Context, agentID uuid.UUID, req UpsertPortShareRequest) (PortShare, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/agents/%s/port-share", agentID), req)
	if err != nil {
		return PortShare{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PortShare{}, client.ReadBodyAsError(resp)
	}

	var share PortShare
	if err := json.NewDecoder(resp.Body).Decode(&share); err != nil {
		return PortShare{}, fmt.Errorf("decode response: %w", err)
	}
	return share, nil
}

// DeletePortShare removes a port share from an agent.
func (s *Service) DeletePortShare(ctx context.Context, agentID uuid.UUID, req DeletePortShareRequest) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/agents/%s/port-share", agentID), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
