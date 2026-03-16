package coordination

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Event represents a coordination event between stacks.
type Event struct {
	Topic   string          `json:"topic"`
	Payload json.RawMessage `json:"payload"`
	Source  string          `json:"source"`
}

// PublishRequest publishes an event to a topic.
type PublishRequest struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

// Publish sends an event to a coordination topic.
func (s *Service) Publish(ctx context.Context, req PublishRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/coordination/events", req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// Subscribe opens an SSE stream for events on a topic.
// Returns an iterator function that yields events.
func (s *Service) Subscribe(ctx context.Context, topic string) (func() (*client.ServerSentEvent, error), error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/coordination/events/subscribe?topic=%s", topic), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	return client.ServerSentEventReader(resp.Body), nil
}
