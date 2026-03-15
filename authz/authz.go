// Package authz provides authorization and RBAC services for the Lattice Runtime API.
package authz

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides authorization operations.
type Service struct {
	client *client.Client
}

// New creates a new authz Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
