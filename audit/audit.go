// Package audit provides audit trail services for the Lattice Runtime API.
package audit

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides audit trail operations.
type Service struct {
	client *client.Client
}

// New creates a new audit Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
