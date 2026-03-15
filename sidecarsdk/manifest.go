package sidecarsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Manifest retrieves the full configuration manifest for this sidecar from Runtime.
func (s *Service) Manifest(ctx context.Context) (Manifest, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/agentsidecars/me/manifest", nil)
	if err != nil {
		return Manifest{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Manifest{}, client.ReadBodyAsError(resp)
	}

	var result Manifest
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Manifest{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}
