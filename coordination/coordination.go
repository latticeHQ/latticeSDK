// Package coordination provides cross-stack messaging and shared state services.
package coordination

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides coordination operations between stacks.
type Service struct {
	client *client.Client
}

// New creates a new coordination Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
