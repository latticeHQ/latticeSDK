package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// CreateToken creates a new API token for the authenticated user.
func (s *Service) CreateToken(ctx context.Context, req types.CreateTokenRequest) (types.GenerateAPIKeyResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/users/me/keys/tokens", req)
	if err != nil {
		return types.GenerateAPIKeyResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return types.GenerateAPIKeyResponse{}, client.ReadBodyAsError(resp)
	}

	var result types.GenerateAPIKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return types.GenerateAPIKeyResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// ListTokens returns all tokens for the authenticated user.
func (s *Service) ListTokens(ctx context.Context) ([]types.APIKey, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/me/keys/tokens", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var tokens []types.APIKey
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return tokens, nil
}

// RevokeToken deletes a token by its ID.
func (s *Service) RevokeToken(ctx context.Context, keyID string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/users/me/keys/%s", keyID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetTokenConfig returns token configuration (e.g., max lifetime).
func (s *Service) GetTokenConfig(ctx context.Context) (types.TokenConfig, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/me/keys/tokens/tokenconfig", nil)
	if err != nil {
		return types.TokenConfig{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.TokenConfig{}, client.ReadBodyAsError(resp)
	}

	var config types.TokenConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return types.TokenConfig{}, fmt.Errorf("decode response: %w", err)
	}
	return config, nil
}
