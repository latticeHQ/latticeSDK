package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// usersResponse is the API response for listing users.
type usersResponse struct {
	Users []types.User `json:"users"`
	Count int          `json:"count"`
}

// GetCurrentUser returns the authenticated user.
func (s *Service) GetCurrentUser(ctx context.Context) (types.User, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/me", nil)
	if err != nil {
		return types.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.User{}, client.ReadBodyAsError(resp)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return types.User{}, fmt.Errorf("decode response: %w", err)
	}
	return user, nil
}

// GetUser returns a user by ID.
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (types.User, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/users/%s", id), nil)
	if err != nil {
		return types.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.User{}, client.ReadBodyAsError(resp)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return types.User{}, fmt.Errorf("decode response: %w", err)
	}
	return user, nil
}

// ListUsers returns users matching the given parameters.
func (s *Service) ListUsers(ctx context.Context, search string, status types.UserStatus, p *client.Pagination) ([]types.User, int, error) {
	var opts []client.RequestOption
	if search != "" {
		opts = append(opts, client.WithQueryParam("q", search))
	}
	if status != "" {
		opts = append(opts, client.WithQueryParam("status", string(status)))
	}
	if p != nil {
		opts = append(opts, p.AsRequestOption())
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users", nil, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, client.ReadBodyAsError(resp)
	}

	var result usersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}
	return result.Users, result.Count, nil
}

// UpdateUserProfile updates the current user's profile.
func (s *Service) UpdateUserProfile(ctx context.Context, id uuid.UUID, req types.UpdateUserProfileRequest) (types.User, error) {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/profile", id), req)
	if err != nil {
		return types.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.User{}, client.ReadBodyAsError(resp)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return types.User{}, fmt.Errorf("decode response: %w", err)
	}
	return user, nil
}

// UpdateUserRoles updates a user's site-wide roles.
func (s *Service) UpdateUserRoles(ctx context.Context, id uuid.UUID, req types.UpdateRoles) (types.User, error) {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/users/%s/roles", id), req)
	if err != nil {
		return types.User{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.User{}, client.ReadBodyAsError(resp)
	}

	var user types.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return types.User{}, fmt.Errorf("decode response: %w", err)
	}
	return user, nil
}

// GetAuthMethods returns the available authentication methods.
func (s *Service) GetAuthMethods(ctx context.Context) (types.AuthMethods, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/authmethods", nil)
	if err != nil {
		return types.AuthMethods{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.AuthMethods{}, client.ReadBodyAsError(resp)
	}

	var methods types.AuthMethods
	if err := json.NewDecoder(resp.Body).Decode(&methods); err != nil {
		return types.AuthMethods{}, fmt.Errorf("decode response: %w", err)
	}
	return methods, nil
}
