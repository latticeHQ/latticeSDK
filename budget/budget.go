// Package budget provides budget and spend control services for the Lattice Runtime API.
package budget

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides budget management operations.
type Service struct {
	client *client.Client
}

// New creates a new budget Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
