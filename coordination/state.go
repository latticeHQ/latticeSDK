package coordination

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// SharedState represents key-value state shared between stacks.
type SharedState struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// GetSharedState retrieves shared state by key.
func (s *Service) GetSharedState(ctx context.Context, key string) (SharedState, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/coordination/state/%s", key), nil)
	if err != nil {
		return SharedState{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SharedState{}, client.ReadBodyAsError(resp)
	}

	var state SharedState
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return SharedState{}, fmt.Errorf("decode response: %w", err)
	}
	return state, nil
}

// SetSharedState sets shared state for the given key.
func (s *Service) SetSharedState(ctx context.Context, key string, value interface{}) error {
	body := map[string]interface{}{
		"key":   key,
		"value": value,
	}
	resp, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/coordination/state/%s", key), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
