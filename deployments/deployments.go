package deployments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides deployment management operations against the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new deployments Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// ListDeployments returns all deployments the authenticated user has access to.
func (s *Service) ListDeployments(ctx context.Context) ([]Deployment, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/deployments", nil)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var deployments []Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return deployments, nil
}

// GetDeploymentBySlug fetches a deployment by its URL-safe slug.
func (s *Service) GetDeploymentBySlug(ctx context.Context, slug string) (Deployment, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/deployments/%s", slug), nil)
	if err != nil {
		return Deployment{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Deployment{}, client.ReadBodyAsError(resp)
	}

	var deployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return Deployment{}, fmt.Errorf("decode response: %w", err)
	}
	return deployment, nil
}

// GetDeploymentByID fetches a deployment by its UUID.
func (s *Service) GetDeploymentByID(ctx context.Context, id uuid.UUID) (Deployment, error) {
	return s.GetDeploymentBySlug(ctx, id.String())
}

// CreateDeployment creates a new deployment instance.
func (s *Service) CreateDeployment(ctx context.Context, req CreateDeploymentRequest) (Deployment, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/deployments", req)
	if err != nil {
		return Deployment{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Deployment{}, client.ReadBodyAsError(resp)
	}

	var deployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return Deployment{}, fmt.Errorf("decode response: %w", err)
	}
	return deployment, nil
}

// UpdateDeployment updates an existing deployment's settings.
// The slugOrID parameter accepts either a deployment slug or UUID string.
func (s *Service) UpdateDeployment(ctx context.Context, slugOrID string, req UpdateDeploymentRequest) (Deployment, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/deployments/%s", slugOrID), req)
	if err != nil {
		return Deployment{}, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Deployment{}, client.ReadBodyAsError(resp)
	}

	var deployment Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return Deployment{}, fmt.Errorf("decode response: %w", err)
	}
	return deployment, nil
}

// DeleteDeployment removes a deployment. Only non-default deployments can be deleted.
// The slugOrID parameter accepts either a deployment slug or UUID string.
func (s *Service) DeleteDeployment(ctx context.Context, slugOrID string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/deployments/%s", slugOrID), nil)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
