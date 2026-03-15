package coordination

import (
	"context"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// EscalationRequest sends work to another stack for processing.
type EscalationRequest struct {
	TargetStack string      `json:"target_stack"`
	Action      string      `json:"action"`
	Payload     interface{} `json:"payload"`
	Priority    string      `json:"priority,omitempty"`
}

// Escalate sends an escalation request to another Department Stack.
func (s *Service) Escalate(ctx context.Context, req EscalationRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/coordination/escalate", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
