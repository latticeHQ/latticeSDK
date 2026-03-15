package sidecarsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// PostStartup reports the sidecar's startup information to Runtime.
func (s *Service) PostStartup(ctx context.Context, req PostStartupRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/me/startup", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// PostLifecycle reports a lifecycle state change to Runtime.
func (s *Service) PostLifecycle(ctx context.Context, req PostLifecycleRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/me/lifecycle", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// PostMetadata sends metadata key-value pairs to Runtime.
func (s *Service) PostMetadata(ctx context.Context, req PostMetadataRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/me/metadata", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// PostStats sends connection and session statistics to Runtime and receives the
// recommended report interval.
func (s *Service) PostStats(ctx context.Context, req Stats) (StatsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/me/stats", req)
	if err != nil {
		return StatsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return StatsResponse{}, client.ReadBodyAsError(resp)
	}

	var result StatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return StatsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
