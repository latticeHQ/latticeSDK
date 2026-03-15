package provisioners

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides provisioner daemon and key management operations.
type Service struct {
	client *client.Client
}

// New creates a new provisioners Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// ListDaemons returns all provisioner daemons available to the caller.
func (s *Service) ListDaemons(ctx context.Context) ([]ProvisionerDaemon, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		"/api/v2/organizations/default/provisionerdaemons", nil)
	if err != nil {
		return nil, fmt.Errorf("request provisioner daemons: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var daemons []ProvisionerDaemon
	return daemons, json.NewDecoder(resp.Body).Decode(&daemons)
}

// ListOrgDaemons returns provisioner daemons for the given organization,
// optionally filtered by tags. Pass nil for tags to list all daemons.
func (s *Service) ListOrgDaemons(ctx context.Context, orgID uuid.UUID, tags map[string]string) ([]ProvisionerDaemon, error) {
	path := fmt.Sprintf("/api/v2/organizations/%s/provisionerdaemons", orgID.String())

	var opts []client.RequestOption
	if len(tags) > 0 {
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			return nil, fmt.Errorf("marshal tags: %w", err)
		}
		opts = append(opts, client.WithQueryParam("tags", string(tagsJSON)))
	}

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil, opts...)
	if err != nil {
		return nil, fmt.Errorf("request organization provisioner daemons: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var daemons []ProvisionerDaemon
	return daemons, json.NewDecoder(resp.Body).Decode(&daemons)
}

// CreateProvisionerKey creates a new provisioner key for an organization.
func (s *Service) CreateProvisionerKey(ctx context.Context, orgID uuid.UUID, req CreateProvisionerKeyRequest) (CreateProvisionerKeyResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/provisionerkeys", orgID.String()), req)
	if err != nil {
		return CreateProvisionerKeyResponse{}, fmt.Errorf("request create provisioner key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return CreateProvisionerKeyResponse{}, client.ReadBodyAsError(resp)
	}

	var result CreateProvisionerKeyResponse
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// ListProvisionerKeys lists all provisioner keys for an organization.
func (s *Service) ListProvisionerKeys(ctx context.Context, orgID uuid.UUID) ([]ProvisionerKey, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/provisionerkeys", orgID.String()), nil)
	if err != nil {
		return nil, fmt.Errorf("request list provisioner keys: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var keys []ProvisionerKey
	return keys, json.NewDecoder(resp.Body).Decode(&keys)
}

// GetProvisionerKey returns the provisioner key identified by the given key string.
func (s *Service) GetProvisionerKey(ctx context.Context, pk string) (ProvisionerKey, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/provisionerkeys/%s", pk), nil,
		func(req *http.Request) {
			req.Header.Add("Lattice-Provisioner-Daemon-Key", pk)
		},
	)
	if err != nil {
		return ProvisionerKey{}, fmt.Errorf("request get provisioner key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ProvisionerKey{}, client.ReadBodyAsError(resp)
	}

	var key ProvisionerKey
	return key, json.NewDecoder(resp.Body).Decode(&key)
}

// ListProvisionerKeyDaemons lists all provisioner keys with their associated
// daemons for an organization.
func (s *Service) ListProvisionerKeyDaemons(ctx context.Context, orgID uuid.UUID) ([]ProvisionerKeyDaemons, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/provisionerkeys/daemons", orgID.String()), nil)
	if err != nil {
		return nil, fmt.Errorf("request list provisioner key daemons: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var result []ProvisionerKeyDaemons
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

// DeleteProvisionerKey deletes a provisioner key by name within an organization.
func (s *Service) DeleteProvisionerKey(ctx context.Context, orgID uuid.UUID, name string) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/organizations/%s/provisionerkeys/%s", orgID.String(), name), nil)
	if err != nil {
		return fmt.Errorf("request delete provisioner key: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}
