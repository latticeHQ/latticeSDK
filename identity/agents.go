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

// agentsResponse is the API response for listing agents.
type agentsResponse struct {
	Agents []types.Agent `json:"agents"`
	Count  int           `json:"count"`
}

// ListAgents returns agents matching the given filter.
// Pass nil for filter to list all agents.
func (s *Service) ListAgents(ctx context.Context, filter *types.AgentFilter) ([]types.Agent, int, error) {
	var opts []client.RequestOption
	if filter != nil {
		if filter.Owner != "" {
			opts = append(opts, client.WithQueryParam("owner", filter.Owner))
		}
		if filter.Template != "" {
			opts = append(opts, client.WithQueryParam("template", filter.Template))
		}
		if filter.Name != "" {
			opts = append(opts, client.WithQueryParam("name", filter.Name))
		}
		if filter.Status != "" {
			opts = append(opts, client.WithQueryParam("status", filter.Status))
		}
		if filter.FilterQuery != "" {
			opts = append(opts, client.WithQueryParam("q", filter.FilterQuery))
		}
		if filter.Limit > 0 {
			opts = append(opts, client.WithQueryParam("limit", fmt.Sprintf("%d", filter.Limit)))
		}
		if filter.Offset > 0 {
			opts = append(opts, client.WithQueryParam("offset", fmt.Sprintf("%d", filter.Offset)))
		}
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/agents", nil, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, client.ReadBodyAsError(resp)
	}

	var result agentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}
	return result.Agents, result.Count, nil
}

// GetAgent returns a single agent by ID.
func (s *Service) GetAgent(ctx context.Context, id uuid.UUID) (types.Agent, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s", id), nil)
	if err != nil {
		return types.Agent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Agent{}, client.ReadBodyAsError(resp)
	}

	var agent types.Agent
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return types.Agent{}, fmt.Errorf("decode response: %w", err)
	}
	return agent, nil
}

// CreateAgent creates a new agent.
func (s *Service) CreateAgent(ctx context.Context, orgID uuid.UUID, req types.CreateAgentRequest) (types.Agent, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/organizations/%s/agents", orgID), req)
	if err != nil {
		return types.Agent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return types.Agent{}, client.ReadBodyAsError(resp)
	}

	var agent types.Agent
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return types.Agent{}, fmt.Errorf("decode response: %w", err)
	}
	return agent, nil
}

// UpdateAgent updates an existing agent.
func (s *Service) UpdateAgent(ctx context.Context, id uuid.UUID, req types.UpdateAgentRequest) (types.Agent, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/agents/%s", id), req)
	if err != nil {
		return types.Agent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Agent{}, client.ReadBodyAsError(resp)
	}

	var agent types.Agent
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return types.Agent{}, fmt.Errorf("decode response: %w", err)
	}
	return agent, nil
}

// DeleteAgent permanently deletes an agent.
func (s *Service) DeleteAgent(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/agents/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// CreateAgentBuild triggers a new build for an agent.
func (s *Service) CreateAgentBuild(ctx context.Context, agentID uuid.UUID, req types.CreateAgentBuildRequest) (types.AgentBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/agents/%s/builds", agentID), req)
	if err != nil {
		return types.AgentBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return types.AgentBuild{}, client.ReadBodyAsError(resp)
	}

	var build types.AgentBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return types.AgentBuild{}, fmt.Errorf("decode response: %w", err)
	}
	return build, nil
}

// GetAgentQuota returns budget/quota information for a member within an organization.
func (s *Service) GetAgentQuota(ctx context.Context, orgID, userID uuid.UUID) (types.AgentQuota, error) {
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

// WatchAgent opens an SSE stream for agent updates.
// The returned function yields events. Call it in a loop until it returns an error.
func (s *Service) WatchAgent(ctx context.Context, id uuid.UUID) (func() (*client.ServerSentEvent, error), error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s/watch", id), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	return client.ServerSentEventReader(resp.Body), nil
}
