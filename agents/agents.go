package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides agent lifecycle management operations.
type Service struct {
	client *client.Client
}

// New creates a new agents Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// ListAgents returns all agents the authenticated user has access to,
// filtered by the provided AgentFilter. Pass nil for no filtering.
func (s *Service) ListAgents(ctx context.Context, filter *AgentFilter) (AgentsResponse, error) {
	var opts []client.RequestOption
	if filter != nil {
		opts = append(opts, agentFilterOption(filter))
		if filter.Limit > 0 || filter.Offset > 0 {
			opts = append(opts, client.Pagination{
				Limit:  filter.Limit,
				Offset: filter.Offset,
			}.AsRequestOption())
		}
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/agents", nil, opts...)
	if err != nil {
		return AgentsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentsResponse{}, client.ReadBodyAsError(resp)
	}

	var result AgentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return AgentsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// agentFilterOption builds query parameters from an AgentFilter.
func agentFilterOption(f *AgentFilter) client.RequestOption {
	return func(r *http.Request) {
		var params []string
		if f.Owner != "" {
			params = append(params, fmt.Sprintf("owner:%q", f.Owner))
		}
		if f.Name != "" {
			params = append(params, fmt.Sprintf("name:%q", f.Name))
		}
		if f.Template != "" {
			params = append(params, fmt.Sprintf("template:%q", f.Template))
		}
		if f.Status != "" {
			params = append(params, fmt.Sprintf("status:%q", f.Status))
		}
		if f.FilterQuery != "" {
			params = append(params, f.FilterQuery)
		}
		if len(params) > 0 {
			q := r.URL.Query()
			q.Set("q", strings.Join(params, " "))
			r.URL.RawQuery = q.Encode()
		}
	}
}

// GetAgent returns a single agent by ID.
func (s *Service) GetAgent(ctx context.Context, id uuid.UUID) (Agent, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s", id), nil)
	if err != nil {
		return Agent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Agent{}, client.ReadBodyAsError(resp)
	}

	var agent Agent
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return Agent{}, fmt.Errorf("decode response: %w", err)
	}
	return agent, nil
}

// GetAgentByOwnerAndName returns an agent by the owner's username and agent name.
func (s *Service) GetAgentByOwnerAndName(ctx context.Context, owner, name string) (Agent, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s/agent/%s", owner, name), nil)
	if err != nil {
		return Agent{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Agent{}, client.ReadBodyAsError(resp)
	}

	var agent Agent
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return Agent{}, fmt.Errorf("decode response: %w", err)
	}
	return agent, nil
}

// CreateAgent creates a new agent build under the specified user.
// The request body is a CreateAgentBuildRequest which specifies the template
// version and transition.
func (s *Service) CreateAgent(ctx context.Context, user string, req CreateAgentBuildRequest) (AgentBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/users/%s/agents", user), req)
	if err != nil {
		return AgentBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return AgentBuild{}, client.ReadBodyAsError(resp)
	}

	var build AgentBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return AgentBuild{}, fmt.Errorf("decode response: %w", err)
	}
	return build, nil
}

// UpdateAgent updates an existing agent's properties.
func (s *Service) UpdateAgent(ctx context.Context, id uuid.UUID, req UpdateAgentRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/agents/%s", id), req)
	if err != nil {
		return fmt.Errorf("update agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// DeleteAgent triggers a delete transition build for an agent.
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

// UpdateAgentAutostart sets the autostart schedule for an agent.
// If the provided schedule is nil, autostart is disabled.
func (s *Service) UpdateAgentAutostart(ctx context.Context, id uuid.UUID, req UpdateAgentAutostartRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/agents/%s/autostart", id), req)
	if err != nil {
		return fmt.Errorf("update agent autostart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateAgentTTL sets the time-to-live for an agent.
// If the provided TTL is nil, autostop is disabled.
func (s *Service) UpdateAgentTTL(ctx context.Context, id uuid.UUID, req UpdateAgentTTLRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/agents/%s/ttl", id), req)
	if err != nil {
		return fmt.Errorf("update agent ttl: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ExtendAgent updates the deadline for resources of the latest agent build.
func (s *Service) ExtendAgent(ctx context.Context, id uuid.UUID, req ExtendAgentRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/agents/%s/extend", id), req)
	if err != nil {
		return fmt.Errorf("extend agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateAgentDormancy activates or makes an agent dormant.
// Setting dormant=true makes the agent dormant; dormant=false activates it.
func (s *Service) UpdateAgentDormancy(ctx context.Context, id uuid.UUID, req UpdateAgentDormancyRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/agents/%s/dormant", id), req)
	if err != nil {
		return fmt.Errorf("update agent dormancy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateAgentAutomaticUpdates sets the automatic updates setting for an agent.
func (s *Service) UpdateAgentAutomaticUpdates(ctx context.Context, id uuid.UUID, req UpdateAgentAutomaticUpdatesRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/agents/%s/autoupdates", id), req)
	if err != nil {
		return fmt.Errorf("update agent automatic updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// FavoriteAgent marks an agent as a favorite.
func (s *Service) FavoriteAgent(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/agents/%s/favorite", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UnfavoriteAgent removes an agent from favorites.
func (s *Service) UnfavoriteAgent(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/agents/%s/favorite", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// WatchAgent opens an SSE stream that delivers agent updates in real time.
// The returned function yields events when called repeatedly.
// Returns (nil, io.EOF) when the stream ends.
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

// GetAgentTimings returns timing measurements for the latest build of an agent.
func (s *Service) GetAgentTimings(ctx context.Context, id uuid.UUID) (AgentBuildTimings, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s/timings", id), nil)
	if err != nil {
		return AgentBuildTimings{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentBuildTimings{}, client.ReadBodyAsError(resp)
	}

	var timings AgentBuildTimings
	if err := json.NewDecoder(resp.Body).Decode(&timings); err != nil {
		return AgentBuildTimings{}, fmt.Errorf("decode response: %w", err)
	}
	return timings, nil
}

// ResolveAutostart checks whether the agent's autostart can proceed
// without parameter mismatches.
func (s *Service) ResolveAutostart(ctx context.Context, id string) (ResolveAutostartResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s/resolve-autostart", id), nil)
	if err != nil {
		return ResolveAutostartResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResolveAutostartResponse{}, client.ReadBodyAsError(resp)
	}

	var result ResolveAutostartResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ResolveAutostartResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// PostAgentUsage marks the agent as having been used recently.
func (s *Service) PostAgentUsage(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/agents/%s/usage", id), nil)
	if err != nil {
		return fmt.Errorf("post agent usage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// PostAgentUsageWithBody marks the agent as having been used recently
// and records an app stat with the given sidecar and app name.
func (s *Service) PostAgentUsageWithBody(ctx context.Context, id uuid.UUID, req PostAgentUsageRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/agents/%s/usage", id), req)
	if err != nil {
		return fmt.Errorf("post agent usage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetAgentQuota returns budget and quota information for a user in an organization.
func (s *Service) GetAgentQuota(ctx context.Context, orgID, userID string) (AgentQuota, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s/members/%s/agent-quota", orgID, userID), nil)
	if err != nil {
		return AgentQuota{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentQuota{}, client.ReadBodyAsError(resp)
	}

	var quota AgentQuota
	if err := json.NewDecoder(resp.Body).Decode(&quota); err != nil {
		return AgentQuota{}, fmt.Errorf("decode response: %w", err)
	}
	return quota, nil
}
