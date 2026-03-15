package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// CreateTemplateVersion creates a new template version within an organization.
func (s *Service) CreateTemplateVersion(ctx context.Context, orgID uuid.UUID, req CreateTemplateVersionRequest) (TemplateVersion, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/templateversions", orgID), req)
	if err != nil {
		return TemplateVersion{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return TemplateVersion{}, client.ReadBodyAsError(resp)
	}

	var version TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return TemplateVersion{}, fmt.Errorf("decode response: %w", err)
	}
	return version, nil
}

// GetTemplateVersion returns a template version by ID.
func (s *Service) GetTemplateVersion(ctx context.Context, id uuid.UUID) (TemplateVersion, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s", id), nil)
	if err != nil {
		return TemplateVersion{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemplateVersion{}, client.ReadBodyAsError(resp)
	}

	var version TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return TemplateVersion{}, fmt.Errorf("decode response: %w", err)
	}
	return version, nil
}

// GetTemplateVersionByName returns a template version by its friendly name within a template.
func (s *Service) GetTemplateVersionByName(ctx context.Context, templateID uuid.UUID, name string) (TemplateVersion, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/templates/%s/versions/%s", templateID, name), nil)
	if err != nil {
		return TemplateVersion{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemplateVersion{}, client.ReadBodyAsError(resp)
	}

	var version TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return TemplateVersion{}, fmt.Errorf("decode response: %w", err)
	}
	return version, nil
}

// GetTemplateVersionByOrgAndName returns a template version by organization, template name, and version name.
func (s *Service) GetTemplateVersionByOrgAndName(ctx context.Context, orgID uuid.UUID, templateName, versionName string) (TemplateVersion, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/templates/%s/versions/%s", orgID, templateName, versionName), nil)
	if err != nil {
		return TemplateVersion{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemplateVersion{}, client.ReadBodyAsError(resp)
	}

	var version TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return TemplateVersion{}, fmt.Errorf("decode response: %w", err)
	}
	return version, nil
}

// UpdateTemplateVersion patches a template version's name or message.
func (s *Service) UpdateTemplateVersion(ctx context.Context, id uuid.UUID, req PatchTemplateVersionRequest) (TemplateVersion, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/templateversions/%s", id), req)
	if err != nil {
		return TemplateVersion{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemplateVersion{}, client.ReadBodyAsError(resp)
	}

	var version TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return TemplateVersion{}, fmt.Errorf("decode response: %w", err)
	}
	return version, nil
}

// CancelTemplateVersion marks a template version job as canceled.
func (s *Service) CancelTemplateVersion(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/templateversions/%s/cancel", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ListTemplateVersions lists versions associated with a template.
func (s *Service) ListTemplateVersions(ctx context.Context, templateID uuid.UUID, req TemplateVersionsByTemplateRequest) ([]TemplateVersion, error) {
	u := fmt.Sprintf("/api/v2/templates/%s/versions", templateID)
	if req.IncludeArchived {
		u += "?include_archived=true"
	}

	var opts []client.RequestOption
	if req.Pagination != nil {
		opts = append(opts, req.Pagination.AsRequestOption())
	}

	resp, err := s.client.Request(ctx, http.MethodGet, u, nil, opts...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var versions []TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return versions, nil
}

// GetTemplateVersionRichParameters returns the rich parameters for a template version.
func (s *Service) GetTemplateVersionRichParameters(ctx context.Context, id uuid.UUID) ([]TemplateVersionParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s/rich-parameters", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []TemplateVersionParameter
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return params, nil
}

// GetTemplateVersionResolvedParameters returns parameters with inheritance resolution and provenance tracking.
func (s *Service) GetTemplateVersionResolvedParameters(ctx context.Context, id uuid.UUID) ([]ResolvedParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s/resolved-parameters", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []ResolvedParameter
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return params, nil
}

// GetTemplateVersionExternalAuth returns external authentication providers for a template version.
func (s *Service) GetTemplateVersionExternalAuth(ctx context.Context, id uuid.UUID) ([]TemplateVersionExternalAuth, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s/external-auth", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var extAuth []TemplateVersionExternalAuth
	if err := json.NewDecoder(resp.Body).Decode(&extAuth); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return extAuth, nil
}

// GetTemplateVersionResources returns the resources declared by a template version.
func (s *Service) GetTemplateVersionResources(ctx context.Context, id uuid.UUID) ([]AgentResource, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s/resources", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var resources []AgentResource
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return resources, nil
}

// GetTemplateVersionVariables returns the variables for a template version.
func (s *Service) GetTemplateVersionVariables(ctx context.Context, id uuid.UUID) ([]TemplateVersionVariable, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s/variables", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var variables []TemplateVersionVariable
	if err := json.NewDecoder(resp.Body).Decode(&variables); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return variables, nil
}

// WatchTemplateVersionLogs streams provisioner job logs for a template version via SSE.
// Logs that occurred after the given log ID are returned. Pass 0 to start from the beginning.
// The returned io.Closer must be called to release the connection.
func (s *Service) WatchTemplateVersionLogs(ctx context.Context, id uuid.UUID, after int64) (<-chan ProvisionerJobLog, io.Closer, error) {
	u := fmt.Sprintf("/api/v2/templateversions/%s/logs?follow", id)
	if after != 0 {
		u += fmt.Sprintf("&after=%d", after)
	}

	resp, err := s.client.Request(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, nil, client.ReadBodyAsError(resp)
	}

	return streamProvisionerJobLogs(resp.Body)
}

// CreateDryRun begins a dry-run provisioner job against the given template version.
func (s *Service) CreateDryRun(ctx context.Context, versionID uuid.UUID, req CreateTemplateVersionDryRunRequest) (ProvisionerJob, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/templateversions/%s/dry-run", versionID), req)
	if err != nil {
		return ProvisionerJob{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return ProvisionerJob{}, client.ReadBodyAsError(resp)
	}

	var job ProvisionerJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return ProvisionerJob{}, fmt.Errorf("decode response: %w", err)
	}
	return job, nil
}

// GetDryRun returns the current state of a template version dry-run job.
func (s *Service) GetDryRun(ctx context.Context, versionID, jobID uuid.UUID) (ProvisionerJob, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/templateversions/%s/dry-run/%s", versionID, jobID), nil)
	if err != nil {
		return ProvisionerJob{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ProvisionerJob{}, client.ReadBodyAsError(resp)
	}

	var job ProvisionerJob
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return ProvisionerJob{}, fmt.Errorf("decode response: %w", err)
	}
	return job, nil
}

// GetDryRunResources returns the resources of a completed dry-run job.
func (s *Service) GetDryRunResources(ctx context.Context, versionID, jobID uuid.UUID) ([]AgentResource, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/templateversions/%s/dry-run/%s/resources", versionID, jobID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var resources []AgentResource
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return resources, nil
}

// WatchDryRunLogs streams provisioner job logs for a dry-run via SSE.
// Logs that occurred after the given log ID are returned. Pass 0 to start from the beginning.
// The returned io.Closer must be called to release the connection.
func (s *Service) WatchDryRunLogs(ctx context.Context, versionID, jobID uuid.UUID, after int64) (<-chan ProvisionerJobLog, io.Closer, error) {
	u := fmt.Sprintf("/api/v2/templateversions/%s/dry-run/%s/logs?follow", versionID, jobID)
	if after != 0 {
		u += fmt.Sprintf("&after=%d", after)
	}

	resp, err := s.client.Request(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, nil, client.ReadBodyAsError(resp)
	}

	return streamProvisionerJobLogs(resp.Body)
}

// CancelDryRun marks a template version dry-run job as canceled.
func (s *Service) CancelDryRun(ctx context.Context, versionID, jobID uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/templateversions/%s/dry-run/%s/cancel", versionID, jobID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// PreviousTemplateVersion returns the version that precedes the named version in a template.
func (s *Service) PreviousTemplateVersion(ctx context.Context, orgID uuid.UUID, templateName, versionName string) (TemplateVersion, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/templates/%s/versions/%s/previous", orgID, templateName, versionName), nil)
	if err != nil {
		return TemplateVersion{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemplateVersion{}, client.ReadBodyAsError(resp)
	}

	var version TemplateVersion
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return TemplateVersion{}, fmt.Errorf("decode response: %w", err)
	}
	return version, nil
}

// streamProvisionerJobLogs reads SSE events from the response body and sends
// decoded ProvisionerJobLog entries to a channel. The returned io.Closer
// closes the underlying response body.
func streamProvisionerJobLogs(body io.ReadCloser) (<-chan ProvisionerJobLog, io.Closer, error) {
	logCh := make(chan ProvisionerJobLog, 8)
	nextEvent := client.ServerSentEventReader(body)

	go func() {
		defer close(logCh)
		defer body.Close()
		for {
			sse, err := nextEvent()
			if err != nil {
				return
			}
			switch sse.Type {
			case client.ServerSentEventTypePing:
				continue
			case client.ServerSentEventTypeError:
				return
			case client.ServerSentEventTypeData:
				var log ProvisionerJobLog
				if err := json.Unmarshal(sse.Data, &log); err != nil {
					return
				}
				logCh <- log
			}
		}
	}()

	return logCh, body, nil
}

