package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides template management operations against the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new templates Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// GetTemplate returns a single template by ID.
func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (Template, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templates/%s", id), nil)
	if err != nil {
		return Template{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Template{}, client.ReadBodyAsError(resp)
	}

	var tmpl Template
	if err := json.NewDecoder(resp.Body).Decode(&tmpl); err != nil {
		return Template{}, fmt.Errorf("decode response: %w", err)
	}
	return tmpl, nil
}

// ListTemplates returns all templates matching the given filter.
func (s *Service) ListTemplates(ctx context.Context, filter TemplateFilter) ([]Template, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/templates", nil, filter.AsRequestOption())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var templates []Template
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return templates, nil
}

// ListTemplatesByOrg returns all templates belonging to the given organization.
func (s *Service) ListTemplatesByOrg(ctx context.Context, orgID uuid.UUID) ([]Template, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s/templates", orgID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var templates []Template
	if err := json.NewDecoder(resp.Body).Decode(&templates); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return templates, nil
}

// GetTemplateByName returns a template by name within an organization (case-insensitive).
func (s *Service) GetTemplateByName(ctx context.Context, orgID uuid.UUID, name string) (Template, error) {
	if name == "" {
		return Template{}, fmt.Errorf("template name cannot be empty")
	}
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/organizations/%s/templates/%s", orgID, name), nil)
	if err != nil {
		return Template{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Template{}, client.ReadBodyAsError(resp)
	}

	var tmpl Template
	if err := json.NewDecoder(resp.Body).Decode(&tmpl); err != nil {
		return Template{}, fmt.Errorf("decode response: %w", err)
	}
	return tmpl, nil
}

// CreateTemplate creates a new template within an organization.
func (s *Service) CreateTemplate(ctx context.Context, orgID uuid.UUID, req CreateTemplateRequest) (Template, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/organizations/%s/templates", orgID), req)
	if err != nil {
		return Template{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Template{}, client.ReadBodyAsError(resp)
	}

	var tmpl Template
	if err := json.NewDecoder(resp.Body).Decode(&tmpl); err != nil {
		return Template{}, fmt.Errorf("decode response: %w", err)
	}
	return tmpl, nil
}

// UpdateTemplateMeta updates the metadata for a template.
func (s *Service) UpdateTemplateMeta(ctx context.Context, id uuid.UUID, req UpdateTemplateMeta) (Template, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/templates/%s", id), req)
	if err != nil {
		return Template{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		return Template{}, fmt.Errorf("template metadata not modified")
	}
	if resp.StatusCode != http.StatusOK {
		return Template{}, client.ReadBodyAsError(resp)
	}

	var tmpl Template
	if err := json.NewDecoder(resp.Body).Decode(&tmpl); err != nil {
		return Template{}, fmt.Errorf("decode response: %w", err)
	}
	return tmpl, nil
}

// DeleteTemplate deletes a template by ID.
func (s *Service) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/templates/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetTemplateACL returns the access control list for a template.
func (s *Service) GetTemplateACL(ctx context.Context, id uuid.UUID) (TemplateACL, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templates/%s/acl", id), nil)
	if err != nil {
		return TemplateACL{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemplateACL{}, client.ReadBodyAsError(resp)
	}

	var acl TemplateACL
	if err := json.NewDecoder(resp.Body).Decode(&acl); err != nil {
		return TemplateACL{}, fmt.Errorf("decode response: %w", err)
	}
	return acl, nil
}

// UpdateTemplateACL updates the access control list for a template.
func (s *Service) UpdateTemplateACL(ctx context.Context, id uuid.UUID, req UpdateTemplateACL) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/templates/%s/acl", id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetTemplateACLAvailable returns users and groups that can be added to a template ACL.
func (s *Service) GetTemplateACLAvailable(ctx context.Context, id uuid.UUID) (ACLAvailable, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templates/%s/acl/available", id), nil)
	if err != nil {
		return ACLAvailable{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ACLAvailable{}, client.ReadBodyAsError(resp)
	}

	var avail ACLAvailable
	if err := json.NewDecoder(resp.Body).Decode(&avail); err != nil {
		return ACLAvailable{}, fmt.Errorf("decode response: %w", err)
	}
	return avail, nil
}

// UpdateActiveTemplateVersion promotes a version to be the active version of a template.
// The version must be attached to the template.
func (s *Service) UpdateActiveTemplateVersion(ctx context.Context, id uuid.UUID, req UpdateActiveTemplateVersion) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/templates/%s/versions", id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ArchiveTemplateVersions archives unused versions of a template.
// When all is true, all unused versions are archived regardless of job status.
// When all is false, only failed versions are archived.
func (s *Service) ArchiveTemplateVersions(ctx context.Context, id uuid.UUID, all bool) (ArchiveTemplateVersionsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/templates/%s/versions/archive", id),
		ArchiveTemplateVersionsRequest{All: all},
	)
	if err != nil {
		return ArchiveTemplateVersionsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ArchiveTemplateVersionsResponse{}, client.ReadBodyAsError(resp)
	}

	var result ArchiveTemplateVersionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ArchiveTemplateVersionsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// SetArchiveTemplateVersion sets or clears the archived status of a single template version.
func (s *Service) SetArchiveTemplateVersion(ctx context.Context, versionID uuid.UUID, archive bool) error {
	action := "/unarchive"
	if archive {
		action = "/archive"
	}
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/templateversions/%s%s", versionID, action), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetTemplateDAUs returns daily active user statistics for a template.
// Use tzOffset=0 for UTC or time-based offsets for local time zones.
func (s *Service) GetTemplateDAUs(ctx context.Context, id uuid.UUID, tzOffset int) (*DAUsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/templates/%s/daus", id),
		nil,
		DAURequest{TZHourOffset: tzOffset}.AsRequestOption(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var daus DAUsResponse
	if err := json.NewDecoder(resp.Body).Decode(&daus); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &daus, nil
}

// ListTemplateExamples returns starter templates for an organization.
//
// Deprecated: Use StarterTemplates instead.
func (s *Service) ListTemplateExamples(ctx context.Context, _ uuid.UUID) ([]TemplateExample, error) {
	return s.StarterTemplates(ctx)
}

// StarterTemplates returns the list of example/starter templates available in Lattice.
func (s *Service) StarterTemplates(ctx context.Context) ([]TemplateExample, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/templates/examples", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var examples []TemplateExample
	if err := json.NewDecoder(resp.Body).Decode(&examples); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return examples, nil
}
