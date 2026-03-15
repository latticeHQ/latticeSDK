package healthsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Service provides deployment health check and diagnostics.
type Service struct {
	client *client.Client
}

// New creates a new health check Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// DebugHealth returns the full health report for the deployment.
func (s *Service) DebugHealth(ctx context.Context) (HealthcheckReport, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/debug/health", nil)
	if err != nil {
		return HealthcheckReport{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthcheckReport{}, client.ReadBodyAsError(resp)
	}

	var report HealthcheckReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return HealthcheckReport{}, fmt.Errorf("decode response: %w", err)
	}
	return report, nil
}

// HealthSettings returns the current health check settings.
func (s *Service) HealthSettings(ctx context.Context) (HealthSettings, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/debug/health/settings", nil)
	if err != nil {
		return HealthSettings{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthSettings{}, client.ReadBodyAsError(resp)
	}

	var settings HealthSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return HealthSettings{}, fmt.Errorf("decode response: %w", err)
	}
	return settings, nil
}

// PutHealthSettings updates the health check settings.
func (s *Service) PutHealthSettings(ctx context.Context, settings UpdateHealthSettings) error {
	resp, err := s.client.Request(ctx, http.MethodPut, "/api/v2/debug/health/settings", settings)
	if err != nil {
		return fmt.Errorf("update health settings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
