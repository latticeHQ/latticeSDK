package sidecarsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/latticehq/latticesdk/client"
)

// GitSSHKey retrieves the SSH key pair used for Git operations.
func (s *Service) GitSSHKey(ctx context.Context) (GitSSHKey, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/agentsidecars/me/gitsshkey", nil)
	if err != nil {
		return GitSSHKey{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitSSHKey{}, client.ReadBodyAsError(resp)
	}

	var result GitSSHKey
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GitSSHKey{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// ExternalAuth requests external authentication credentials from Runtime.
func (s *Service) ExternalAuth(ctx context.Context, req ExternalAuthRequest) (ExternalAuthResponse, error) {
	q := url.Values{}
	if req.ID != "" {
		q.Set("id", req.ID)
	}
	if req.Match != "" {
		q.Set("match", req.Match)
	}
	if req.Listen {
		q.Set("listen", strconv.FormatBool(req.Listen))
	}

	path := "/api/v2/agentsidecars/me/external-auth"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ExternalAuthResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ExternalAuthResponse{}, client.ReadBodyAsError(resp)
	}

	var result ExternalAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ExternalAuthResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
