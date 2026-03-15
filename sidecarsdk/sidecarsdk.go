// Package sidecarsdk provides the sidecar-to-Runtime communication API.
// Use this package when building a custom sidecar or agent process that needs
// to report status, send logs, and authenticate with Runtime.
package sidecarsdk

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides the sidecar-to-Runtime communication API.
// Use this when building a custom sidecar or agent that needs to
// report status, send logs, and authenticate with Runtime.
type Service struct {
	client *client.Client
}

// New creates a new sidecar SDK Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
