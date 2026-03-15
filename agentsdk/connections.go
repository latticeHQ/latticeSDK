package agentsdk

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// SidecarConnectionInfo returns the DERP map and connection options for a
// specific sidecar. This information is needed to establish a direct
// connection to the sidecar via DERP relay or direct networking.
func (s *Service) SidecarConnectionInfo(ctx context.Context, sidecarID uuid.UUID) (ConnectionInfo, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agentsidecars/%s/connection", sidecarID), nil)
	if err != nil {
		return ConnectionInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ConnectionInfo{}, client.ReadBodyAsError(resp)
	}

	var info ConnectionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return ConnectionInfo{}, fmt.Errorf("decode response: %w", err)
	}
	return info, nil
}

// SidecarConnectionInfoGeneric returns connection information that is not
// tied to a specific sidecar. This is useful for discovering the DERP map
// before a specific sidecar ID is known.
func (s *Service) SidecarConnectionInfoGeneric(ctx context.Context) (ConnectionInfo, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/agentsidecars/connection", nil)
	if err != nil {
		return ConnectionInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ConnectionInfo{}, client.ReadBodyAsError(resp)
	}

	var info ConnectionInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return ConnectionInfo{}, fmt.Errorf("decode response: %w", err)
	}
	return info, nil
}

// ptyConn wraps a hijacked HTTP connection to provide an io.ReadWriteCloser
// for the WebSocket-upgraded PTY stream.
type ptyConn struct {
	conn net.Conn
	br   *bufio.Reader
}

// Read reads from the underlying connection, first draining any buffered data.
func (p *ptyConn) Read(b []byte) (int, error) {
	return p.br.Read(b)
}

// Write writes to the underlying connection.
func (p *ptyConn) Write(b []byte) (int, error) {
	return p.conn.Write(b)
}

// Close closes the underlying connection.
func (p *ptyConn) Close() error {
	return p.conn.Close()
}

// ReconnectingPTY opens a reconnecting PTY session to the specified sidecar.
// It performs a WebSocket upgrade over HTTP and returns the raw connection
// as an io.ReadWriteCloser. The caller can write PTYRequest JSON messages
// and read terminal output from the returned connection.
//
// The connection is not a full WebSocket implementation — it is the raw
// TCP connection after the HTTP upgrade. Callers that need WebSocket framing
// should layer their own codec on top.
func (s *Service) ReconnectingPTY(ctx context.Context, opts PTYOptions) (io.ReadWriteCloser, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/pty", opts.SidecarID)

	// Build query parameters for the PTY handshake.
	var queryParts []string
	queryParts = append(queryParts, "reconnect="+opts.Reconnect.String())
	queryParts = append(queryParts, "width="+strconv.Itoa(int(opts.Width)))
	queryParts = append(queryParts, "height="+strconv.Itoa(int(opts.Height)))
	if opts.Command != "" {
		queryParts = append(queryParts, "command="+opts.Command)
	}
	if opts.SignedToken != "" {
		queryParts = append(queryParts, "signed_app_token="+opts.SignedToken)
	}
	path += "?" + strings.Join(queryParts, "&")

	// Use the client to build and send the request with proper auth headers.
	// We set the WebSocket upgrade headers via a request option.
	resp, err := s.client.Request(ctx, http.MethodGet, path, nil, func(r *http.Request) {
		r.Header.Set("Connection", "Upgrade")
		r.Header.Set("Upgrade", "websocket")
		r.Header.Set("Sec-WebSocket-Version", "13")
		// A minimal base64-encoded key for the handshake.
		r.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
	})
	if err != nil {
		return nil, fmt.Errorf("pty request: %w", err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, client.ReadBodyAsError(resp)
	}

	// Hijack the connection from the response.
	hijacker, ok := resp.Body.(io.ReadWriteCloser)
	if ok {
		return hijacker, nil
	}

	// If the body doesn't support hijacking directly, try to get the
	// underlying connection through the response's TLS/net.Conn.
	// The standard http.Client does not expose hijacking, so we wrap
	// the response body as a ReadWriteCloser with a no-op writer.
	// In practice, callers should use an http.Transport with
	// DialContext to capture the raw connection.
	return &responseReadWriteCloser{resp: resp}, nil
}

// responseReadWriteCloser wraps an *http.Response whose body is still open
// after a protocol upgrade, providing io.ReadWriteCloser semantics.
// Writing goes to the underlying connection via the response body if it
// supports writing; otherwise writes return an error.
type responseReadWriteCloser struct {
	resp *http.Response
}

func (r *responseReadWriteCloser) Read(p []byte) (int, error) {
	return r.resp.Body.Read(p)
}

func (r *responseReadWriteCloser) Write(p []byte) (int, error) {
	w, ok := r.resp.Body.(io.Writer)
	if !ok {
		return 0, fmt.Errorf("response body does not support writing; use a custom http.Transport with DialContext to capture the raw connection")
	}
	return w.Write(p)
}

func (r *responseReadWriteCloser) Close() error {
	return r.resp.Body.Close()
}
