package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// CreateRealTimeSession creates a new real-time session in the specified organization.
func (s *Service) CreateRealTimeSession(ctx context.Context, org string, req CreateRealTimeSessionRequest) (RealTimeSession, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/realtime-sessions", org),
		req,
	)
	if err != nil {
		return RealTimeSession{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return RealTimeSession{}, client.ReadBodyAsError(resp)
	}

	var session RealTimeSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return RealTimeSession{}, fmt.Errorf("decode real-time session: %w", err)
	}
	return session, nil
}

// ListRealTimeSessions returns all real-time sessions in the specified organization,
// filtered by the given criteria.
func (s *Service) ListRealTimeSessions(ctx context.Context, org string, filter RealTimeSessionFilter) (RealTimeSessionsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/realtime-sessions", org),
		nil,
		filter.asRequestOption(),
	)
	if err != nil {
		return RealTimeSessionsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RealTimeSessionsResponse{}, client.ReadBodyAsError(resp)
	}

	var result RealTimeSessionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RealTimeSessionsResponse{}, fmt.Errorf("decode real-time sessions: %w", err)
	}
	return result, nil
}

// GetRealTimeSession returns a specific real-time session by ID.
func (s *Service) GetRealTimeSession(ctx context.Context, id uuid.UUID) (RealTimeSession, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/realtime-sessions/%s", id),
		nil,
	)
	if err != nil {
		return RealTimeSession{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RealTimeSession{}, client.ReadBodyAsError(resp)
	}

	var session RealTimeSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return RealTimeSession{}, fmt.Errorf("decode real-time session: %w", err)
	}
	return session, nil
}

// UpdateRealTimeSession updates an existing real-time session.
func (s *Service) UpdateRealTimeSession(ctx context.Context, id uuid.UUID, req UpdateRealTimeSessionRequest) (RealTimeSession, error) {
	resp, err := s.client.Request(ctx, http.MethodPatch,
		fmt.Sprintf("/api/v2/realtime-sessions/%s", id),
		req,
	)
	if err != nil {
		return RealTimeSession{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RealTimeSession{}, client.ReadBodyAsError(resp)
	}

	var session RealTimeSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return RealTimeSession{}, fmt.Errorf("decode real-time session: %w", err)
	}
	return session, nil
}

// DeleteRealTimeSession deletes a real-time session by ID.
func (s *Service) DeleteRealTimeSession(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/realtime-sessions/%s", id),
		nil,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetRealTimeSessionParameters returns the resolved parameters for a specific real-time session.
func (s *Service) GetRealTimeSessionParameters(ctx context.Context, id uuid.UUID) ([]RealTimeSessionParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/realtime-sessions/%s/parameters", id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []RealTimeSessionParameter
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		return nil, fmt.Errorf("decode real-time session parameters: %w", err)
	}
	return params, nil
}

// ListMyRealTimeSessions returns all real-time sessions owned by the authenticated user.
func (s *Service) ListMyRealTimeSessions(ctx context.Context) (RealTimeSessionsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/users/me/realtime-sessions", nil)
	if err != nil {
		return RealTimeSessionsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RealTimeSessionsResponse{}, client.ReadBodyAsError(resp)
	}

	var result RealTimeSessionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RealTimeSessionsResponse{}, fmt.Errorf("decode real-time sessions: %w", err)
	}
	return result, nil
}
