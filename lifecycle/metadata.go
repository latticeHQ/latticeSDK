package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Metadata contains key-value metadata for a stack.
type Metadata map[string]interface{}

// GetMetadata returns the metadata for a stack.
func (s *Service) GetMetadata(ctx context.Context, stackID uuid.UUID) (Metadata, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/stacks/%s/metadata", stackID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var meta Metadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return meta, nil
}

// SetMetadata sets the metadata for a stack.
func (s *Service) SetMetadata(ctx context.Context, stackID uuid.UUID, meta Metadata) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/stacks/%s/metadata", stackID), meta)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
