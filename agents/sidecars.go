package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// GetSidecar returns a single agent sidecar by its ID.
func (s *Service) GetSidecar(ctx context.Context, id uuid.UUID) (AgentSidecar, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentsidecars/%s", id), nil)
	if err != nil {
		return AgentSidecar{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AgentSidecar{}, client.ReadBodyAsError(resp)
	}

	var sidecar AgentSidecar
	if err := json.NewDecoder(resp.Body).Decode(&sidecar); err != nil {
		return AgentSidecar{}, fmt.Errorf("decode response: %w", err)
	}
	return sidecar, nil
}

// WatchSidecarMetadata opens an SSE stream for sidecar metadata updates.
// Returns a channel of metadata slices and an error channel. The error
// channel receives exactly one error when the stream ends. The metadata
// channel is never closed by the sender; it becomes unreadable when the
// error channel fires.
func (s *Service) WatchSidecarMetadata(ctx context.Context, id uuid.UUID) (<-chan []SidecarMetadata, <-chan error) {
	metadataChan := make(chan []SidecarMetadata, 256)
	errorChan := make(chan error, 1)
	ready := make(chan struct{})

	go func() {
		defer close(errorChan)

		resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentsidecars/%s/watch-metadata", id), nil)
		if err != nil {
			close(ready)
			errorChan <- err
			return
		}

		if resp.StatusCode != http.StatusOK {
			close(ready)
			errorChan <- client.ReadBodyAsError(resp)
			return
		}

		nextEvent := client.ServerSentEventReader(resp.Body)
		defer resp.Body.Close()

		firstEvent := true
		for {
			select {
			case <-ctx.Done():
				if firstEvent {
					close(ready)
				}
				errorChan <- ctx.Err()
				return
			default:
			}

			sse, err := nextEvent()
			if err != nil {
				if firstEvent {
					close(ready)
				}
				errorChan <- err
				return
			}

			if firstEvent {
				close(ready)
				firstEvent = false
			}

			if sse.Type == client.ServerSentEventTypePing {
				continue
			}

			if sse.Type == client.ServerSentEventTypeError {
				errorChan <- fmt.Errorf("sse error: %s", string(sse.Data))
				return
			}

			if sse.Type == client.ServerSentEventTypeData {
				var met []SidecarMetadata
				if err := json.Unmarshal(sse.Data, &met); err != nil {
					errorChan <- fmt.Errorf("unmarshal metadata: %w", err)
					return
				}
				select {
				case <-ctx.Done():
					errorChan <- ctx.Err()
					return
				case metadataChan <- met:
				}
			}
		}
	}()

	// Wait until the first event is received so the subscription is established.
	<-ready

	return metadataChan, errorChan
}

// GetSidecarListeningPorts returns the list of ports currently being listened
// on inside the sidecar's network namespace.
func (s *Service) GetSidecarListeningPorts(ctx context.Context, sidecarID uuid.UUID) (SidecarListeningPortsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentsidecars/%s/listening-ports", sidecarID), nil)
	if err != nil {
		return SidecarListeningPortsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SidecarListeningPortsResponse{}, client.ReadBodyAsError(resp)
	}

	var ports SidecarListeningPortsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ports); err != nil {
		return SidecarListeningPortsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return ports, nil
}

// WatchSidecarLogs streams sidecar logs. When follow is false, it fetches
// existing logs and returns them on the channel, then closes it. When follow
// is true, it opens an SSE stream and sends log batches as they arrive.
// The after parameter specifies the log ID after which to begin.
func (s *Service) WatchSidecarLogs(ctx context.Context, sidecarID uuid.UUID, after int64, follow bool) (<-chan []SidecarLog, <-chan error) {
	logChan := make(chan []SidecarLog, 256)
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)
		defer close(logChan)

		var opts []client.RequestOption
		if after > 0 {
			opts = append(opts, client.WithQueryParam("after", fmt.Sprintf("%d", after)))
		}
		if follow {
			opts = append(opts, client.WithQueryParam("follow", "true"))
		}

		resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentsidecars/%s/logs", sidecarID), nil, opts...)
		if err != nil {
			errorChan <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errorChan <- client.ReadBodyAsError(resp)
			return
		}

		if !follow {
			// Non-follow: decode all logs at once.
			var logs []SidecarLog
			if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
				errorChan <- fmt.Errorf("decode logs: %w", err)
				return
			}
			logChan <- logs
			return
		}

		// Follow mode: read SSE events.
		nextEvent := client.ServerSentEventReader(resp.Body)
		for {
			select {
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			default:
			}

			sse, err := nextEvent()
			if err != nil {
				errorChan <- err
				return
			}

			if sse.Type == client.ServerSentEventTypePing {
				continue
			}

			if sse.Type == client.ServerSentEventTypeError {
				errorChan <- fmt.Errorf("sse error: %s", string(sse.Data))
				return
			}

			if sse.Type == client.ServerSentEventTypeData {
				var logs []SidecarLog
				if err := json.Unmarshal(sse.Data, &logs); err != nil {
					errorChan <- fmt.Errorf("unmarshal logs: %w", err)
					return
				}
				select {
				case <-ctx.Done():
					errorChan <- ctx.Err()
					return
				case logChan <- logs:
				}
			}
		}
	}()

	return logChan, errorChan
}
