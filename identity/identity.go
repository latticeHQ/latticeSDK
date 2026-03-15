// Package identity provides identity management services for the Lattice Runtime API.
// It covers agents, users, organizations, and API tokens.
package identity

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides identity management operations.
type Service struct {
	client *client.Client
}

// New creates a new identity Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
