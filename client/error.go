package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Response is the standard API response envelope.
type Response struct {
	Message     string            `json:"message"`
	Detail      string            `json:"detail,omitempty"`
	Validations []ValidationError `json:"validations,omitempty"`
}

// ValidationError represents a field-level validation error.
type ValidationError struct {
	Field  string `json:"field"`
	Detail string `json:"detail"`
}

// Error returns a formatted validation error string.
func (v ValidationError) Error() string {
	return fmt.Sprintf("field: %s detail: %s", v.Field, v.Detail)
}

// Error represents an API error response.
type Error struct {
	Response
	statusCode int
	method     string
	url        string
	Helper     string
}

// StatusCode returns the HTTP status code.
func (e *Error) StatusCode() int { return e.statusCode }

// Method returns the HTTP method of the failed request.
func (e *Error) Method() string { return e.method }

// URL returns the URL of the failed request.
func (e *Error) URL() string { return e.url }

// Error implements the error interface.
func (e *Error) Error() string {
	var b strings.Builder
	_, _ = fmt.Fprintf(&b, "%s %s: status %d", e.method, e.url, e.statusCode)
	if e.Message != "" {
		_, _ = fmt.Fprintf(&b, "\nmessage: %s", e.Message)
	}
	if e.Helper != "" {
		_, _ = fmt.Fprintf(&b, "\nhelper: %s", e.Helper)
	}
	if e.Detail != "" {
		_, _ = fmt.Fprintf(&b, "\ndetail: %s", e.Detail)
	}
	for _, v := range e.Validations {
		_, _ = fmt.Fprintf(&b, "\nvalidation: %s", v.Error())
	}
	return b.String()
}

// Friendly returns a user-friendly error message.
func (e *Error) Friendly() string {
	if e.Message == "" {
		return fmt.Sprintf("unexpected status code %d", e.statusCode)
	}
	var b strings.Builder
	b.WriteString(e.Message)
	if e.Helper != "" {
		_, _ = fmt.Fprintf(&b, "\n%s", e.Helper)
	}
	for _, v := range e.Validations {
		_, _ = fmt.Fprintf(&b, "\n- %s: %s", v.Field, v.Detail)
	}
	return b.String()
}

// ReadBodyAsError reads an HTTP response body and returns it as an *Error.
func ReadBodyAsError(resp *http.Response) error {
	defer resp.Body.Close()

	apiErr := &Error{
		statusCode: resp.StatusCode,
		method:     resp.Request.Method,
		url:        resp.Request.URL.String(),
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		apiErr.Message = fmt.Sprintf("read response body: %v", err)
		return apiErr
	}

	if len(body) == 0 {
		apiErr.Message = "empty response body"
		return apiErr
	}

	if err := json.Unmarshal(body, &apiErr.Response); err != nil {
		apiErr.Message = string(body)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		apiErr.Helper = "Try re-authenticating with a new API key or session token."
	}

	return apiErr
}

// IsConnectionError returns true if the error is a network connection error.
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "connection refused") ||
		strings.Contains(s, "no such host") ||
		strings.Contains(s, "network is unreachable")
}

// AsError attempts to unwrap err as an *Error.
func AsError(err error) (*Error, bool) {
	var apiErr *Error
	if e, ok := err.(*Error); ok {
		return e, true
	}
	_ = apiErr
	return nil, false
}
