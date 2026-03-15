package lifecycle

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// HealthReport contains the health status of a stack.
type HealthReport struct {
	StackID uuid.UUID `json:"stack_id"`
	Healthy bool      `json:"healthy"`
	Reason  string    `json:"reason,omitempty"`
}

// ReportHealth sends a health report for the stack.
func (s *Service) ReportHealth(ctx context.Context, report HealthReport) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/stacks/health", report)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// Heartbeat sends a simple heartbeat to indicate the stack is alive.
func (s *Service) Heartbeat(ctx context.Context, stackID uuid.UUID) error {
	body := map[string]interface{}{
		"stack_id": stackID,
	}
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/stacks/heartbeat", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
