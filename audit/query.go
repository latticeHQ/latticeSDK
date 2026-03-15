package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/types"
)

// auditLogResponse is the API response for listing audit logs.
type auditLogResponse struct {
	AuditLogs []types.AuditLog `json:"audit_logs"`
	Count     int64            `json:"count"`
}

// Query returns audit logs matching the given filter.
// Pass nil for filter to list all audit logs.
func (s *Service) Query(ctx context.Context, filter *types.AuditLogFilter, p *client.Pagination) ([]types.AuditLog, int64, error) {
	var opts []client.RequestOption
	if filter != nil && filter.SearchQuery != "" {
		opts = append(opts, client.WithQueryParam("q", filter.SearchQuery))
	}
	if p != nil {
		opts = append(opts, p.AsRequestOption())
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/audit", nil, opts...)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, client.ReadBodyAsError(resp)
	}

	var result auditLogResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("decode response: %w", err)
	}
	return result.AuditLogs, result.Count, nil
}
