package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// StackRegistration contains the registration details for a Department Stack.
type StackRegistration struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Version  string    `json:"version"`
	Endpoint string    `json:"endpoint"`
	Status   string    `json:"status"`
}

// RegisterRequest contains the parameters for registering a stack.
type RegisterRequest struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Endpoint string `json:"endpoint,omitempty"`
}

// RegisterStack registers a Department Stack with the Runtime.
func (s *Service) RegisterStack(ctx context.Context, req RegisterRequest) (StackRegistration, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/stacks/register", req)
	if err != nil {
		return StackRegistration{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return StackRegistration{}, client.ReadBodyAsError(resp)
	}

	var reg StackRegistration
	if err := json.NewDecoder(resp.Body).Decode(&reg); err != nil {
		return StackRegistration{}, fmt.Errorf("decode response: %w", err)
	}
	return reg, nil
}

// Deregister removes a stack registration from the Runtime.
func (s *Service) Deregister(ctx context.Context, stackID uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/stacks/%s", stackID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
