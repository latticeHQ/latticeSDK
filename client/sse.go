package client

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ServerSentEvent represents a single server-sent event.
type ServerSentEvent struct {
	Type ServerSentEventType `json:"type"`
	Data []byte              `json:"data,omitempty"`
}

// ServerSentEventType identifies the type of SSE event.
type ServerSentEventType string

const (
	ServerSentEventTypePing  ServerSentEventType = "ping"
	ServerSentEventTypeData  ServerSentEventType = "data"
	ServerSentEventTypeError ServerSentEventType = "error"
)

// ServerSentEventReader returns an iterator function that reads server-sent
// events from the given ReadCloser. Call the returned function repeatedly
// to get the next event. Returns (nil, io.EOF) when the stream ends.
func ServerSentEventReader(rc io.ReadCloser) func() (*ServerSentEvent, error) {
	reader := bufio.NewReader(rc)

	nextLineValue := func(prefix string) ([]byte, error) {
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("reading line: %w", err)
			}
			if strings.TrimSpace(line) != "" {
				if !strings.HasPrefix(line, prefix+": ") {
					return nil, fmt.Errorf("expected %q prefix, got: %s", prefix, line)
				}
				s := strings.TrimPrefix(line, prefix+": ")
				s = strings.TrimSpace(s)
				return []byte(s), nil
			}
		}
	}

	return func() (*ServerSentEvent, error) {
		t, err := nextLineValue("event")
		if err != nil {
			return nil, fmt.Errorf("reading event type: %w", err)
		}

		switch ServerSentEventType(t) {
		case ServerSentEventTypePing:
			return &ServerSentEvent{Type: ServerSentEventTypePing}, nil
		case ServerSentEventTypeData:
			d, err := nextLineValue("data")
			if err != nil {
				return nil, fmt.Errorf("reading data: %w", err)
			}
			return &ServerSentEvent{Type: ServerSentEventTypeData, Data: d}, nil
		case ServerSentEventTypeError:
			d, err := nextLineValue("data")
			if err != nil {
				return nil, fmt.Errorf("reading error data: %w", err)
			}
			return &ServerSentEvent{Type: ServerSentEventTypeError, Data: d}, nil
		default:
			return nil, fmt.Errorf("unknown event type: %s", t)
		}
	}
}
