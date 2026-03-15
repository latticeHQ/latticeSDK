package externalauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Service provides external authentication operations.
type Service struct {
	client *client.Client
}

// New creates a new external auth Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// GetDeviceByID returns the device authorization details for the given provider.
func (s *Service) GetDeviceByID(ctx context.Context, provider string) (ExternalAuthDevice, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/external-auth/%s/device", provider), nil)
	if err != nil {
		return ExternalAuthDevice{}, fmt.Errorf("request external auth device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ExternalAuthDevice{}, client.ReadBodyAsError(resp)
	}

	var device ExternalAuthDevice
	return device, json.NewDecoder(resp.Body).Decode(&device)
}

// ExchangeDevice exchanges a device code for an external auth token.
func (s *Service) ExchangeDevice(ctx context.Context, provider string, req ExternalAuthDeviceExchange) error {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/external-auth/%s/device", provider), req)
	if err != nil {
		return fmt.Errorf("request exchange device: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetByID returns the external auth state for the given provider.
func (s *Service) GetByID(ctx context.Context, provider string) (ExternalAuth, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/external-auth/%s", provider), nil)
	if err != nil {
		return ExternalAuth{}, fmt.Errorf("request external auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ExternalAuth{}, client.ReadBodyAsError(resp)
	}

	var auth ExternalAuth
	return auth, json.NewDecoder(resp.Body).Decode(&auth)
}

// UnlinkByID deletes the external auth link for the given provider.
// This does not revoke the token from the identity provider.
func (s *Service) UnlinkByID(ctx context.Context, provider string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/external-auth/%s", provider), nil)
	if err != nil {
		return fmt.Errorf("request unlink external auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ListExternalAuths returns the available external auth providers and the
// user's authenticated links if they exist.
func (s *Service) ListExternalAuths(ctx context.Context) (ListUserExternalAuthResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/external-auth", nil)
	if err != nil {
		return ListUserExternalAuthResponse{}, fmt.Errorf("request list external auths: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ListUserExternalAuthResponse{}, client.ReadBodyAsError(resp)
	}

	var result ListUserExternalAuthResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}
