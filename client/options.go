package client

import "net/http"

// Option configures a Client.
type Option func(*Client)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

// WithSessionToken sets the session token for authentication.
func WithSessionToken(token string) Option {
	return func(c *Client) {
		c.sessionToken = token
	}
}

// WithHTTPClient sets a custom http.Client for requests.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithHeader adds a custom header to all requests.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.headers.Set(key, value)
	}
}
