package agentsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// ListeningPorts returns the ports that the sidecar is currently listening on.
func (s *Service) ListeningPorts(ctx context.Context, sidecarID uuid.UUID) (ListeningPortsResponse, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/api/v0/listening-ports", sidecarID)

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListeningPortsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ListeningPortsResponse{}, client.ReadBodyAsError(resp)
	}

	var result ListeningPortsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ListeningPortsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// Netcheck runs a network connectivity check on the sidecar and returns
// the raw report. The report structure depends on the sidecar version.
func (s *Service) Netcheck(ctx context.Context, sidecarID uuid.UUID) (NetcheckReport, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/api/v0/netcheck", sidecarID)

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return NetcheckReport{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NetcheckReport{}, client.ReadBodyAsError(resp)
	}

	var report NetcheckReport
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return NetcheckReport{}, fmt.Errorf("decode response: %w", err)
	}
	return report, nil
}

// DebugMagicsock returns raw magicsock debug information from the sidecar.
// The returned bytes are the unprocessed response body.
func (s *Service) DebugMagicsock(ctx context.Context, sidecarID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/debug/magicsock", sidecarID)
	return s.debugEndpoint(ctx, path)
}

// DebugManifest returns the sidecar's debug manifest. The returned bytes
// are the unprocessed response body.
func (s *Service) DebugManifest(ctx context.Context, sidecarID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/debug/manifest", sidecarID)
	return s.debugEndpoint(ctx, path)
}

// DebugLogs returns the sidecar's debug logs. The returned bytes are the
// unprocessed response body.
func (s *Service) DebugLogs(ctx context.Context, sidecarID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/debug/logs", sidecarID)
	return s.debugEndpoint(ctx, path)
}

// PrometheusMetrics returns the sidecar's Prometheus metrics in the standard
// exposition format. The returned bytes are the unprocessed response body.
func (s *Service) PrometheusMetrics(ctx context.Context, sidecarID uuid.UUID) ([]byte, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/debug/prometheus", sidecarID)
	return s.debugEndpoint(ctx, path)
}

// debugEndpoint is a shared helper that performs a GET request and returns
// the raw response body bytes.
func (s *Service) debugEndpoint(ctx context.Context, path string) ([]byte, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return data, nil
}
