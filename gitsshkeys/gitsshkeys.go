package gitsshkeys

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Service provides Git SSH key operations.
type Service struct {
	client *client.Client
}

// New creates a new gitsshkeys Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// GetGitSSHKey returns the Git SSH public key for the given user.
// The user parameter can be a user ID or "me".
func (s *Service) GetGitSSHKey(ctx context.Context, user string) (GitSSHKey, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/users/%s/gitsshkey", user), nil)
	if err != nil {
		return GitSSHKey{}, fmt.Errorf("request get git ssh key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitSSHKey{}, client.ReadBodyAsError(resp)
	}

	var key GitSSHKey
	return key, json.NewDecoder(resp.Body).Decode(&key)
}

// RegenerateGitSSHKey creates a new SSH key pair for the given user and
// returns the new key. The user parameter can be a user ID or "me".
func (s *Service) RegenerateGitSSHKey(ctx context.Context, user string) (GitSSHKey, error) {
	resp, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/users/%s/gitsshkey", user), nil)
	if err != nil {
		return GitSSHKey{}, fmt.Errorf("request regenerate git ssh key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitSSHKey{}, client.ReadBodyAsError(resp)
	}

	var key GitSSHKey
	return key, json.NewDecoder(resp.Body).Decode(&key)
}
