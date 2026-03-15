package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// ListProxies returns all agent proxies.
func (s *Service) ListProxies(ctx context.Context) ([]AgentProxy, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/agentproxies", nil)
	if err != nil {
		return nil, fmt.Errorf("list proxies: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var result ProxiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Regions, nil
}

// GetProxyByID returns a single agent proxy by its UUID.
func (s *Service) GetProxyByID(ctx context.Context, id uuid.UUID) (AgentProxy, error) {
	return s.GetProxyByName(ctx, id.String())
}

// GetProxyByName returns a single agent proxy by its name or ID string.
func (s *Service) GetProxyByName(ctx context.Context, name string) (AgentProxy, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentproxies/%s", name), nil)
	if err != nil {
		return AgentProxy{}, fmt.Errorf("get proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentProxy{}, client.ReadBodyAsError(resp)
	}

	var proxy AgentProxy
	if err := json.NewDecoder(resp.Body).Decode(&proxy); err != nil {
		return AgentProxy{}, fmt.Errorf("decode response: %w", err)
	}
	return proxy, nil
}

// CreateProxy creates a new agent proxy and returns it along with the proxy token.
func (s *Service) CreateProxy(ctx context.Context, req CreateAgentProxyRequest) (UpdateAgentProxyResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentproxies", req)
	if err != nil {
		return UpdateAgentProxyResponse{}, fmt.Errorf("create proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return UpdateAgentProxyResponse{}, client.ReadBodyAsError(resp)
	}

	var result UpdateAgentProxyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return UpdateAgentProxyResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// UpdateProxy updates an existing agent proxy.
func (s *Service) UpdateProxy(ctx context.Context, req PatchAgentProxy) (UpdateAgentProxyResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/agentproxies/%s", req.ID), req)
	if err != nil {
		return UpdateAgentProxyResponse{}, fmt.Errorf("update proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateAgentProxyResponse{}, client.ReadBodyAsError(resp)
	}

	var result UpdateAgentProxyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return UpdateAgentProxyResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// DeleteProxyByID deletes an agent proxy by its UUID.
func (s *Service) DeleteProxyByID(ctx context.Context, id uuid.UUID) error {
	return s.DeleteProxyByName(ctx, id.String())
}

// DeleteProxyByName deletes an agent proxy by its name or ID string.
func (s *Service) DeleteProxyByName(ctx context.Context, name string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/agentproxies/%s", name), nil)
	if err != nil {
		return fmt.Errorf("delete proxy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ListRegions returns all available deployment regions.
func (s *Service) ListRegions(ctx context.Context) ([]Region, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/regions", nil)
	if err != nil {
		return nil, fmt.Errorf("list regions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var result RegionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Regions, nil
}
