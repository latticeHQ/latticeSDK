package agentsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// ListDirectory lists the contents of a directory on the sidecar's filesystem.
// The request specifies the path and relativity (root or home) for resolution.
func (s *Service) ListDirectory(ctx context.Context, sidecarID uuid.UUID, req LSRequest) (LSResponse, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/api/v0/list-directory", sidecarID)
	resp, err := s.client.Request(ctx, http.MethodPost, path, req)
	if err != nil {
		return LSResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return LSResponse{}, client.ReadBodyAsError(resp)
	}

	var result LSResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return LSResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// ReadFile reads a file from the sidecar's filesystem. It returns the file
// contents as an io.ReadCloser, the content type, and any error. The caller
// is responsible for closing the returned reader.
//
// The offset and limit parameters control which byte range is read. Pass 0
// for both to read the entire file.
func (s *Service) ReadFile(ctx context.Context, sidecarID uuid.UUID, filePath string, offset, limit int64) (io.ReadCloser, string, error) {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/api/v0/read-file", sidecarID)

	var opts []client.RequestOption
	opts = append(opts, client.WithQueryParam("path", filePath))
	if offset > 0 {
		opts = append(opts, client.WithQueryParam("offset", strconv.FormatInt(offset, 10)))
	}
	if limit > 0 {
		opts = append(opts, client.WithQueryParam("limit", strconv.FormatInt(limit, 10)))
	}

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil, opts...)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", client.ReadBodyAsError(resp)
	}

	contentType := resp.Header.Get("Content-Type")
	return resp.Body, contentType, nil
}

// WriteFile writes content to a file on the sidecar's filesystem.
// The content is sent as the request body. The path is specified as a query
// parameter.
func (s *Service) WriteFile(ctx context.Context, sidecarID uuid.UUID, filePath string, content io.Reader) error {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/api/v0/write-file", sidecarID)

	resp, err := s.client.Request(ctx, http.MethodPost, path, content,
		client.WithQueryParam("path", filePath),
	)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// EditFiles applies search-and-replace edits to one or more files on the
// sidecar's filesystem. Each FileEdits entry targets a single file and
// contains one or more find/replace operations.
func (s *Service) EditFiles(ctx context.Context, sidecarID uuid.UUID, req FileEditRequest) error {
	path := fmt.Sprintf("/api/v2/agentsidecars/%s/api/v0/edit-files", sidecarID)

	resp, err := s.client.Request(ctx, http.MethodPost, path, req)
	if err != nil {
		return fmt.Errorf("edit files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
