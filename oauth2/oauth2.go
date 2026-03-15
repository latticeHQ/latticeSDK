package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides OAuth2 provider application management against the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new oauth2 Service backed by the given client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// ListProviderApps returns the OAuth2 provider applications, optionally filtered.
func (s *Service) ListProviderApps(ctx context.Context, filter OAuth2ProviderAppFilter) ([]OAuth2ProviderApp, error) {
	var opts []client.RequestOption
	if filter.UserID != uuid.Nil {
		opts = append(opts, client.WithQueryParam("user_id", filter.UserID.String()))
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/oauth2-provider/apps", nil, opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var apps []OAuth2ProviderApp
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return apps, nil
}

// GetProviderApp returns a single OAuth2 provider application by ID.
func (s *Service) GetProviderApp(ctx context.Context, id uuid.UUID) (OAuth2ProviderApp, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/oauth2-provider/apps/%s", id), nil)
	if err != nil {
		return OAuth2ProviderApp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OAuth2ProviderApp{}, client.ReadBodyAsError(resp)
	}

	var app OAuth2ProviderApp
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return OAuth2ProviderApp{}, fmt.Errorf("decode response: %w", err)
	}
	return app, nil
}

// CreateProviderApp creates a new OAuth2 provider application.
func (s *Service) CreateProviderApp(ctx context.Context, req PostOAuth2ProviderAppRequest) (OAuth2ProviderApp, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/oauth2-provider/apps", req)
	if err != nil {
		return OAuth2ProviderApp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return OAuth2ProviderApp{}, client.ReadBodyAsError(resp)
	}

	var app OAuth2ProviderApp
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return OAuth2ProviderApp{}, fmt.Errorf("decode response: %w", err)
	}
	return app, nil
}

// UpdateProviderApp updates an existing OAuth2 provider application.
func (s *Service) UpdateProviderApp(ctx context.Context, id uuid.UUID, req PutOAuth2ProviderAppRequest) (OAuth2ProviderApp, error) {
	resp, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/oauth2-provider/apps/%s", id), req)
	if err != nil {
		return OAuth2ProviderApp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return OAuth2ProviderApp{}, client.ReadBodyAsError(resp)
	}

	var app OAuth2ProviderApp
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return OAuth2ProviderApp{}, fmt.Errorf("decode response: %w", err)
	}
	return app, nil
}

// DeleteProviderApp deletes an OAuth2 provider application, also invalidating
// any tokens that were generated from it.
func (s *Service) DeleteProviderApp(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/oauth2-provider/apps/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ListProviderAppSecrets returns the truncated secrets for an OAuth2 application.
func (s *Service) ListProviderAppSecrets(ctx context.Context, appID uuid.UUID) ([]OAuth2ProviderAppSecret, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/oauth2-provider/apps/%s/secrets", appID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var secrets []OAuth2ProviderAppSecret
	if err := json.NewDecoder(resp.Body).Decode(&secrets); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return secrets, nil
}

// CreateProviderAppSecret creates a new secret for an OAuth2 application.
// This is the only time the full secret will be revealed.
func (s *Service) CreateProviderAppSecret(ctx context.Context, appID uuid.UUID) (OAuth2ProviderAppSecretFull, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/oauth2-provider/apps/%s/secrets", appID), nil)
	if err != nil {
		return OAuth2ProviderAppSecretFull{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return OAuth2ProviderAppSecretFull{}, client.ReadBodyAsError(resp)
	}

	var secret OAuth2ProviderAppSecretFull
	if err := json.NewDecoder(resp.Body).Decode(&secret); err != nil {
		return OAuth2ProviderAppSecretFull{}, fmt.Errorf("decode response: %w", err)
	}
	return secret, nil
}

// DeleteProviderAppSecret deletes a secret from an OAuth2 application,
// also invalidating any tokens that were generated from it.
func (s *Service) DeleteProviderAppSecret(ctx context.Context, appID, secretID uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/oauth2-provider/apps/%s/secrets/%s", appID, secretID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// RevokeProviderApp completely revokes an app's access for the authenticated user.
func (s *Service) RevokeProviderApp(ctx context.Context, appID uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, "/oauth2/tokens", nil,
		client.WithQueryParam("client_id", appID.String()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
