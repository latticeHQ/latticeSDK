package sidecarsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// PostLogSource registers a new log source with Runtime.
func (s *Service) PostLogSource(ctx context.Context, req PostLogSourceRequest) (LogSource, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/agentsidecars/me/log-source", req)
	if err != nil {
		return LogSource{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LogSource{}, client.ReadBodyAsError(resp)
	}

	var result LogSource
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return LogSource{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// PatchLogs sends log entries for a given log source to Runtime.
func (s *Service) PatchLogs(ctx context.Context, req PatchLogsRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, "/api/v2/agentsidecars/me/logs", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
