package licenses

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/latticehq/latticesdk/client"
)

// Service provides license management operations.
type Service struct {
	client *client.Client
}

// New creates a new licenses Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// AddLicense adds a new license to the deployment.
func (s *Service) AddLicense(ctx context.Context, req AddLicenseRequest) (License, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, "/api/v2/licenses", req)
	if err != nil {
		return License{}, fmt.Errorf("request add license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return License{}, client.ReadBodyAsError(resp)
	}

	var l License
	d := json.NewDecoder(resp.Body)
	d.UseNumber()
	return l, d.Decode(&l)
}

// ListLicenses returns all licenses for the deployment.
func (s *Service) ListLicenses(ctx context.Context) ([]License, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/licenses", nil)
	if err != nil {
		return nil, fmt.Errorf("request list licenses: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var licenses []License
	d := json.NewDecoder(resp.Body)
	d.UseNumber()
	return licenses, d.Decode(&licenses)
}

// DeleteLicense deletes a license by its numeric ID.
func (s *Service) DeleteLicense(ctx context.Context, id int32) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/licenses/%d", id), nil)
	if err != nil {
		return fmt.Errorf("request delete license: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
