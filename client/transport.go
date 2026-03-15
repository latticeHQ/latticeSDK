package client

import "net/http"

// HeaderTransport is an http.RoundTripper that adds headers to all requests.
type HeaderTransport struct {
	Transport http.RoundTripper
	Header    http.Header
}

// RoundTrip implements http.RoundTripper.
func (t *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, vs := range t.Header {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	transport := t.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return transport.RoundTrip(req)
}

// CloseIdleConnections closes idle connections on the underlying transport.
func (t *HeaderTransport) CloseIdleConnections() {
	type closeIdler interface {
		CloseIdleConnections()
	}
	if tr, ok := t.Transport.(closeIdler); ok {
		tr.CloseIdleConnections()
	}
}
