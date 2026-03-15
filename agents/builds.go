package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// GetAgentBuild returns a single agent build by its ID.
func (s *Service) GetAgentBuild(ctx context.Context, id uuid.UUID) (AgentBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentbuilds/%s", id), nil)
	if err != nil {
		return AgentBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentBuild{}, client.ReadBodyAsError(resp)
	}

	var build AgentBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return AgentBuild{}, fmt.Errorf("decode response: %w", err)
	}
	return build, nil
}

// GetAgentBuildByOwnerAndName returns an agent build by owner, agent name, and build number.
func (s *Service) GetAgentBuildByOwnerAndName(ctx context.Context, owner, agent, buildNumber string) (AgentBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s/agent/%s/builds/%s", owner, agent, buildNumber), nil)
	if err != nil {
		return AgentBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentBuild{}, client.ReadBodyAsError(resp)
	}

	var build AgentBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return AgentBuild{}, fmt.Errorf("decode response: %w", err)
	}
	return build, nil
}

// CancelAgentBuild marks an agent build job as canceled.
func (s *Service) CancelAgentBuild(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/agentbuilds/%s/cancel", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// CreateAgentBuild queues a new build for the specified agent.
func (s *Service) CreateAgentBuild(ctx context.Context, agentID uuid.UUID, req CreateAgentBuildRequest) (AgentBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/agents/%s/builds", agentID), req)
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

// GetAgentBuildParameters returns the parameters used for a specific build.
func (s *Service) GetAgentBuildParameters(ctx context.Context, buildID uuid.UUID) ([]AgentBuildParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentbuilds/%s/parameters", buildID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []AgentBuildParameter
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return params, nil
}

// GetAgentBuildState returns the provisioner state of the build as raw bytes.
func (s *Service) GetAgentBuildState(ctx context.Context, buildID uuid.UUID) ([]byte, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentbuilds/%s/state", buildID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}
	return io.ReadAll(resp.Body)
}

// GetAgentBuildTimings returns timing measurements for a specific build.
func (s *Service) GetAgentBuildTimings(ctx context.Context, buildID uuid.UUID) (AgentBuildTimings, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentbuilds/%s/timings", buildID), nil)
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

// WatchAgentBuildLogs opens an SSE stream for provisioner job logs.
// The after parameter specifies the log ID after which to begin streaming.
// The returned function yields events when called repeatedly.
func (s *Service) WatchAgentBuildLogs(ctx context.Context, buildID uuid.UUID, after int64) (func() (*client.ServerSentEvent, error), error) {
	path := fmt.Sprintf("/api/v2/agentbuilds/%s/logs", buildID)
	var opts []client.RequestOption
	if after > 0 {
		opts = append(opts, client.WithQueryParam("after", fmt.Sprintf("%d", after)))
	}
	opts = append(opts, client.WithQueryParam("follow", "true"))

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil, opts...)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	return client.ServerSentEventReader(resp.Body), nil
}
