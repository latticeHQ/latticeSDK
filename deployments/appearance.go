package deployments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// GetAppearance returns the appearance configuration that controls the visual
// display of the dashboard, including branding, banners, and navigation.
func (s *Service) GetAppearance(ctx context.Context) (AppearanceConfig, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/appearance", nil)
	if err != nil {
		return AppearanceConfig{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AppearanceConfig{}, client.ReadBodyAsError(resp)
	}

	var cfg AppearanceConfig
	if err := json.NewDecoder(resp.Body).Decode(&cfg); err != nil {
		return AppearanceConfig{}, fmt.Errorf("decode response: %w", err)
	}
	return cfg, nil
}

// UpdateAppearance updates the dashboard appearance settings including
// the application name, logo, banners, and navigation configuration.
func (s *Service) UpdateAppearance(ctx context.Context, req UpdateAppearanceConfig) error {
	resp, err := s.client.Request(ctx, http.MethodPut, "/api/v2/appearance", req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
