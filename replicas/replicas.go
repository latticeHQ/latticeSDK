package replicas

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Service provides replica query operations.
type Service struct {
	client *client.Client
}

// New creates a new replicas Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// ListReplicas returns all replicas in the deployment.
func (s *Service) ListReplicas(ctx context.Context) ([]Replica, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/replicas", nil)
	if err != nil {
		return nil, fmt.Errorf("request list replicas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var replicas []Replica
	return replicas, json.NewDecoder(resp.Body).Decode(&replicas)
}
