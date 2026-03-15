package deployments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// GetConfig returns the deployment configuration for the Lattice server.
// The response includes both the current configuration values and the available
// configuration options.
func (s *Service) GetConfig(ctx context.Context) (DeploymentConfig, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/deployment/config", nil)
	if err != nil {
		return DeploymentConfig{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DeploymentConfig{}, client.ReadBodyAsError(resp)
	}

	var cfg DeploymentConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return DeploymentConfig{}, fmt.Errorf("decode response: %w", err)
	}
	return cfg, nil
}

// GetStats returns aggregate deployment statistics including agent counts,
// session counts, and connection latency.
func (s *Service) GetStats(ctx context.Context) (DeploymentStats, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/deployment/stats", nil)
	if err != nil {
		return DeploymentStats{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DeploymentStats{}, client.ReadBodyAsError(resp)
	}

	var stats DeploymentStats
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return DeploymentStats{}, fmt.Errorf("decode response: %w", err)
	}
	return stats, nil
}

// GetEntitlements returns the feature entitlements for the deployment.
func (s *Service) GetEntitlements(ctx context.Context) (Entitlements, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/entitlements", nil)
	if err != nil {
		return Entitlements{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Entitlements{}, client.ReadBodyAsError(resp)
	}

	var ent Entitlements
	if err := json.NewDecoder(resp.Body).Decode(&ent); err != nil {
		return Entitlements{}, fmt.Errorf("decode response: %w", err)
	}
	return ent, nil
}

// GetExperiments returns the list of experiments enabled on the deployment.
func (s *Service) GetExperiments(ctx context.Context) (Experiments, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/experiments", nil)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var exp Experiments
	if err := json.NewDecoder(resp.Body).Decode(&exp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return exp, nil
}

// GetSafeExperiments returns the list of experiments that are safe for users
// to opt in to.
func (s *Service) GetSafeExperiments(ctx context.Context) (AvailableExperiments, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/experiments/available", nil)
	if err != nil {
		return AvailableExperiments{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AvailableExperiments{}, client.ReadBodyAsError(resp)
	}

	var exp AvailableExperiments
	if err := json.NewDecoder(resp.Body).Decode(&exp); err != nil {
		return AvailableExperiments{}, fmt.Errorf("decode response: %w", err)
	}
	return exp, nil
}

// GetBuildInfo returns build information for the Lattice instance, including
// version, dashboard URL, and API version details.
func (s *Service) GetBuildInfo(ctx context.Context) (BuildInfoResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/buildinfo", nil)
	if err != nil {
		return BuildInfoResponse{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return BuildInfoResponse{}, client.ReadBodyAsError(resp)
	}

	var info BuildInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return BuildInfoResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return info, nil
}

// GetAppHost returns the site-wide application wildcard hostname, e.g.
// "*--apps.latticeruntime.com". If the app host is not configured, the
// response will contain an empty string.
func (s *Service) GetAppHost(ctx context.Context) (AppHostResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/applications/host", nil)
	if err != nil {
		return AppHostResponse{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppHostResponse{}, client.ReadBodyAsError(resp)
	}

	var host AppHostResponse
	if err := json.NewDecoder(resp.Body).Decode(&host); err != nil {
		return AppHostResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return host, nil
}

// GetSSHConfig returns SSH configuration for the Lattice instance, including
// the hostname prefix and any custom SSH config options.
func (s *Service) GetSSHConfig(ctx context.Context) (SSHConfigResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/deployment/ssh", nil)
	if err != nil {
		return SSHConfigResponse{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SSHConfigResponse{}, client.ReadBodyAsError(resp)
	}

	var cfg SSHConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return SSHConfigResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return cfg, nil
}

// UpdateCheck returns information about the latest release version of Lattice
// and whether the server is running the latest release.
func (s *Service) UpdateCheck(ctx context.Context) (UpdateCheckResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/updatecheck", nil)
	if err != nil {
		return UpdateCheckResponse{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateCheckResponse{}, client.ReadBodyAsError(resp)
	}

	var check UpdateCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&check); err != nil {
		return UpdateCheckResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return check, nil
}
