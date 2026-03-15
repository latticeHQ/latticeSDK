package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// GetSidecar returns a session sidecar by ID.
func (s *Service) GetSidecar(ctx context.Context, id uuid.UUID) (SessionSidecar, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionsidecars/%s", id), nil)
	if err != nil {
		return SessionSidecar{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionSidecar{}, client.ReadBodyAsError(resp)
	}

	var sidecar SessionSidecar
	if err := json.NewDecoder(resp.Body).Decode(&sidecar); err != nil {
		return SessionSidecar{}, fmt.Errorf("decode session sidecar: %w", err)
	}
	return sidecar, nil
}

// WatchSidecarMetadata opens an SSE stream that emits sidecar metadata updates.
// It returns a channel of metadata slices and an error channel. The metadata
// channel receives updates whenever metadata values change. Exactly one error
// will be sent on the error channel when the stream ends. The caller should
// cancel the context to stop the stream.
func (s *Service) WatchSidecarMetadata(ctx context.Context, id uuid.UUID) (<-chan []SessionSidecarMetadata, <-chan error) {
	metadataChan := make(chan []SessionSidecarMetadata, 256)
	errorChan := make(chan error, 1)

	ready := make(chan struct{})
	watch := func() error {
		resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionsidecars/%s/watch-metadata", id), nil)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return client.ReadBodyAsError(resp)
		}

		nextEvent := client.ServerSentEventReader(resp.Body)
		defer resp.Body.Close()

		firstEvent := true
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			sse, err := nextEvent()
			if err != nil {
				return err
			}

			if firstEvent {
				close(ready)
				firstEvent = false
			}

			if sse.Type == client.ServerSentEventTypePing {
				continue
			}

			switch sse.Type {
			case client.ServerSentEventTypeData:
				var met []SessionSidecarMetadata
				if err := json.Unmarshal(sse.Data, &met); err != nil {
					return fmt.Errorf("unmarshal metadata: %w", err)
				}
				metadataChan <- met
			case client.ServerSentEventTypeError:
				var r client.Response
				if err := json.Unmarshal(sse.Data, &r); err != nil {
					return fmt.Errorf("unmarshal error: %w", err)
				}
				return fmt.Errorf("%s", r.Message)
			default:
				return fmt.Errorf("unexpected event type: %s", sse.Type)
			}
		}
	}

	go func() {
		defer close(errorChan)
		err := watch()
		select {
		case <-ready:
		default:
			close(ready)
		}
		errorChan <- err
	}()

	// Wait until first event is received and the subscription is registered.
	<-ready

	return metadataChan, errorChan
}

// GetSidecarListeningPorts returns a list of ports currently being listened on
// inside the sidecar's network namespace.
func (s *Service) GetSidecarListeningPorts(ctx context.Context, sidecarID uuid.UUID) (SessionSidecarListeningPortsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionsidecars/%s/listening-ports", sidecarID), nil)
	if err != nil {
		return SessionSidecarListeningPortsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionSidecarListeningPortsResponse{}, client.ReadBodyAsError(resp)
	}

	var ports SessionSidecarListeningPortsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ports); err != nil {
		return SessionSidecarListeningPortsResponse{}, fmt.Errorf("decode listening ports: %w", err)
	}
	return ports, nil
}

// WatchSidecarLogs streams log entries from a session sidecar.
//
// The after parameter specifies the log ID after which to start streaming.
// When follow is false, the method returns all available logs after the given
// ID in a single batch and closes the channel. When follow is true, the method
// opens a streaming connection and continues to emit logs as they arrive.
//
// The returned function yields the next server-sent event on each call.
// Returns (nil, io.EOF) when the stream ends. If follow is false, only a single
// batch of logs is returned. Callers should cancel the context to stop a
// followed stream.
func (s *Service) WatchSidecarLogs(ctx context.Context, sidecarID uuid.UUID, after int64, follow bool) (func() (*client.ServerSentEvent, error), error) {
	path := fmt.Sprintf("/api/v2/sessionsidecars/%s/logs", sidecarID)

	var queryParts []string
	if after != 0 {
		queryParts = append(queryParts, fmt.Sprintf("after=%d", after))
	}
	if follow {
		queryParts = append(queryParts, "follow")
	}
	if len(queryParts) > 0 {
		path += "?"
		for i, part := range queryParts {
			if i > 0 {
				path += "&"
			}
			path += part
		}
	}

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	nextEvent := client.ServerSentEventReader(resp.Body)
	return nextEvent, nil
}
