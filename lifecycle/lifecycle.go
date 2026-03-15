// Package lifecycle provides stack lifecycle management services.
package lifecycle

import (
	"github.com/latticehq/latticesdk/client"
)

// Service provides lifecycle management operations.
type Service struct {
	client *client.Client
}

// New creates a new lifecycle Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}
