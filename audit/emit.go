package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// EmitEventRequest contains the parameters for emitting a custom audit event.
type EmitEventRequest struct {
	Action           types.AuditAction   `json:"action"`
	ResourceType     types.ResourceType  `json:"resource_type"`
	ResourceID       uuid.UUID           `json:"resource_id,omitempty"`
	AdditionalFields json.RawMessage     `json:"additional_fields,omitempty"`
	Time             time.Time           `json:"time,omitempty"`
	OrganizationID   uuid.UUID           `json:"organization_id,omitempty"`
}

// Emit creates a custom audit log entry.
func (s *Service) Emit(ctx context.Context, req EmitEventRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/audit/testgenerate", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetAuditLog returns a single audit log entry by ID.
func (s *Service) GetAuditLog(ctx context.Context, id uuid.UUID) (types.AuditLog, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/audit/%s", id), nil)
	if err != nil {
		return types.AuditLog{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.AuditLog{}, client.ReadBodyAsError(resp)
	}

	var log types.AuditLog
	if err := json.NewDecoder(resp.Body).Decode(&log); err != nil {
		return types.AuditLog{}, fmt.Errorf("decode response: %w", err)
	}
	return log, nil
}
