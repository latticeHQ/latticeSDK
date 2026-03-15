// Package files provides file upload and download for the Lattice Runtime API.
package files

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides file upload and download operations.
type Service struct {
	client *client.Client
}

// New creates a new files Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// Upload uploads an arbitrary file with the specified content type.
// The reader provides the raw file contents. On success, the returned
// UploadResponse contains the file hash that can be used for subsequent
// downloads.
func (s *Service) Upload(ctx context.Context, contentType string, rd io.Reader) (UploadResponse, error) {
	res, err := s.client.Request(ctx, http.MethodPost, "/api/v2/files", rd, func(r *http.Request) {
		r.Header.Set("Content-Type", contentType)
	})
	if err != nil {
		return UploadResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return UploadResponse{}, client.ReadBodyAsError(res)
	}
	var resp UploadResponse
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// Download fetches a file by its uploaded hash. It returns the raw bytes,
// the Content-Type header value, and any error.
func (s *Service) Download(ctx context.Context, id uuid.UUID) ([]byte, string, error) {
	return s.DownloadWithFormat(ctx, id, "")
}

// DownloadWithFormat fetches a file by its uploaded hash, optionally forcing
// format conversion (e.g., FormatZip). It returns the raw bytes, the
// Content-Type header value, and any error.
func (s *Service) DownloadWithFormat(ctx context.Context, id uuid.UUID, format string) ([]byte, string, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/files/%s?format=%s", id, format),
		nil,
	)
	if err != nil {
		return nil, "", fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, "", client.ReadBodyAsError(res)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", fmt.Errorf("read response body: %w", err)
	}
	return data, res.Header.Get("Content-Type"), nil
}
