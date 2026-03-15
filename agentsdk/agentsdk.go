// Package agentsdk provides agent-to-sidecar connectivity operations.
// It wraps the public REST and WebSocket endpoints exposed by Lattice Runtime,
// allowing external code to connect to running agent sidecars.
package agentsdk

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides agent-to-sidecar connectivity operations.
// Use this when connecting to running agent sidecars from external code.
type Service struct {
	client *client.Client
}

// New creates a new agentsdk Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
