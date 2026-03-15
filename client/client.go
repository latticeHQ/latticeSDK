// Package client provides the HTTP client for communicating with Lattice Runtime.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const (
	// SessionTokenHeader is the custom header used for authentication.
	SessionTokenHeader = "Lattice-Session-Token"

	// APIKeyHeader is the header used for API key authentication.
	APIKeyHeader = "Lattice-API-Key"
)

// Client is an HTTP client for the Lattice Runtime API.
type Client struct {
	mu           sync.RWMutex
	httpClient   *http.Client
	baseURL      *url.URL
	apiKey       string
	sessionToken string
	headers      http.Header
}

// New creates a new Client with the given base URL and options.
func New(baseURL string, opts ...Option) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base URL: %w", err)
	}

	c := &Client{
		httpClient: http.DefaultClient,
		baseURL:    u,
		headers:    make(http.Header),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// SessionToken returns the current session token.
func (c *Client) SessionToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sessionToken
}

// SetSessionToken sets the session token for authentication.
func (c *Client) SetSessionToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sessionToken = token
}

// Request performs an HTTP request to the Runtime API.
// The path is relative to the base URL. If body is not nil and not an io.Reader,
// it will be JSON-encoded. The response is returned directly — callers are
// responsible for closing the response body.
func (c *Client) Request(ctx context.Context, method, path string, body interface{}, opts ...RequestOption) (*http.Response, error) {
	reqURL, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse path %q: %w", path, err)
	}

	var bodyReader io.Reader
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			bodyReader = v
		case []byte:
			bodyReader = bytes.NewReader(v)
		default:
			data, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(data)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Set authentication headers.
	c.mu.RLock()
	if c.apiKey != "" {
		req.Header.Set(APIKeyHeader, c.apiKey)
	}
	if c.sessionToken != "" {
		req.Header.Set(SessionTokenHeader, c.sessionToken)
	}
	// Copy custom headers.
	for k, vs := range c.headers {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	c.mu.RUnlock()

	// Apply request options.
	for _, opt := range opts {
		opt(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	return resp, nil
}

// RequestOption modifies an outgoing HTTP request.
type RequestOption func(*http.Request)

// WithQueryParam returns a RequestOption that adds a query parameter.
func WithQueryParam(key, value string) RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		q.Set(key, value)
		r.URL.RawQuery = q.Encode()
	}
}
