package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides session management operations for the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new sessions Service backed by the given client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// ListSessions returns all sessions the authenticated user has access to,
// filtered by the given criteria.
func (s *Service) ListSessions(ctx context.Context, filter SessionFilter) (SessionsResponse, error) {
	page := client.Pagination{
		Offset: filter.Offset,
		Limit:  filter.Limit,
	}
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/sessions", nil, filter.asRequestOption(), page.AsRequestOption())
	if err != nil {
		return SessionsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionsResponse{}, client.ReadBodyAsError(resp)
	}

	var result SessionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SessionsResponse{}, fmt.Errorf("decode sessions response: %w", err)
	}
	return result, nil
}

// GetSession returns a single session by ID.
func (s *Service) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	return s.getSession(ctx, id)
}

// GetDeletedSession returns a single session that was deleted.
// It includes soft-deleted sessions in the lookup.
func (s *Service) GetDeletedSession(ctx context.Context, id uuid.UUID) (Session, error) {
	opts := SessionOptions{IncludeDeleted: true}
	return s.getSession(ctx, id, opts.asRequestOption())
}

func (s *Service) getSession(ctx context.Context, id uuid.UUID, opts ...client.RequestOption) (Session, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessions/%s", id), nil, opts...)
	if err != nil {
		return Session{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Session{}, client.ReadBodyAsError(resp)
	}

	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return Session{}, fmt.Errorf("decode session: %w", err)
	}
	return session, nil
}

// GetSessionByOwnerAndName returns a session by owner username and session name.
func (s *Service) GetSessionByOwnerAndName(ctx context.Context, owner, name string, opts SessionOptions) (Session, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/users/%s/session/%s", owner, name),
		nil,
		func(r *http.Request) {
			q := r.URL.Query()
			q.Set("include_deleted", fmt.Sprintf("%t", opts.IncludeDeleted))
			r.URL.RawQuery = q.Encode()
		},
	)
	if err != nil {
		return Session{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Session{}, client.ReadBodyAsError(resp)
	}

	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return Session{}, fmt.Errorf("decode session: %w", err)
	}
	return session, nil
}

// CreateSessionBuild queues a new build to occur for a session.
// The build transitions the session to the state specified in the request.
func (s *Service) CreateSessionBuild(ctx context.Context, sessionID uuid.UUID, req CreateSessionBuildRequest) (SessionBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/sessions/%s/builds", sessionID), req)
	if err != nil {
		return SessionBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return SessionBuild{}, client.ReadBodyAsError(resp)
	}

	var build SessionBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return SessionBuild{}, fmt.Errorf("decode session build: %w", err)
	}
	return build, nil
}

// WatchSession opens an SSE stream that emits Session updates whenever the
// session's state changes. The returned function yields the next server-sent
// event on each call. Returns (nil, io.EOF) when the stream ends. Callers
// should cancel the context to stop the stream.
func (s *Service) WatchSession(ctx context.Context, id uuid.UUID) (func() (*client.ServerSentEvent, error), error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessions/%s/watch", id), nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	nextEvent := client.ServerSentEventReader(resp.Body)
	return nextEvent, nil
}

// UpdateSession updates a session's mutable fields (currently only name).
func (s *Service) UpdateSession(ctx context.Context, id uuid.UUID, req UpdateSessionRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/sessions/%s", id), req)
	if err != nil {
		return fmt.Errorf("update session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateSessionAutostart sets the autostart schedule for a session.
// If the provided schedule is nil, autostart is disabled.
func (s *Service) UpdateSessionAutostart(ctx context.Context, id uuid.UUID, req UpdateSessionAutostartRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/sessions/%s/autostart", id), req)
	if err != nil {
		return fmt.Errorf("update session autostart: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateSessionTTL sets the time-to-live for a session.
// If the provided TTL is nil, autostop is disabled.
func (s *Service) UpdateSessionTTL(ctx context.Context, id uuid.UUID, req UpdateSessionTTLRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/sessions/%s/ttl", id), req)
	if err != nil {
		return fmt.Errorf("update session TTL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// ExtendSession updates the deadline for resources of the latest session build.
func (s *Service) ExtendSession(ctx context.Context, id uuid.UUID, req PutExtendSessionRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/sessions/%s/extend", id), req)
	if err != nil {
		return fmt.Errorf("extend session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateSessionDormancy sets a session as dormant or activates a dormant session.
func (s *Service) UpdateSessionDormancy(ctx context.Context, id uuid.UUID, req UpdateSessionDormancy) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/sessions/%s/dormant", id), req)
	if err != nil {
		return fmt.Errorf("update session dormancy: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateSessionAutomaticUpdates sets the automatic updates policy for a session.
func (s *Service) UpdateSessionAutomaticUpdates(ctx context.Context, id uuid.UUID, req UpdateSessionAutomaticUpdatesRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/sessions/%s/autoupdates", id), req)
	if err != nil {
		return fmt.Errorf("update session automatic updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// FavoriteSession marks a session as a favorite for the authenticated user.
func (s *Service) FavoriteSession(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPut, fmt.Sprintf("/api/v2/sessions/%s/favorite", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UnfavoriteSession removes a session from the authenticated user's favorites.
func (s *Service) UnfavoriteSession(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/sessions/%s/favorite", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// PostSessionUsage marks the session as having been used recently.
func (s *Service) PostSessionUsage(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/sessions/%s/usage", id), nil)
	if err != nil {
		return fmt.Errorf("post session usage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetSessionQuota returns the session quota for a user in the given organization.
func (s *Service) GetSessionQuota(ctx context.Context, organizationID, userID string) (SessionQuota, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/members/%s/session-quota", organizationID, userID),
		nil,
	)
	if err != nil {
		return SessionQuota{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionQuota{}, client.ReadBodyAsError(resp)
	}

	var quota SessionQuota
	if err := json.NewDecoder(resp.Body).Decode(&quota); err != nil {
		return SessionQuota{}, fmt.Errorf("decode session quota: %w", err)
	}
	return quota, nil
}

// ResolveSessionAutostart checks whether the session's autostart is compatible
// with the current template version parameters.
func (s *Service) ResolveSessionAutostart(ctx context.Context, id uuid.UUID) (ResolveAutostartResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessions/%s/resolve-autostart", id), nil)
	if err != nil {
		return ResolveAutostartResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ResolveAutostartResponse{}, client.ReadBodyAsError(resp)
	}

	var result ResolveAutostartResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ResolveAutostartResponse{}, fmt.Errorf("decode resolve autostart: %w", err)
	}
	return result, nil
}

// GetSessionTimings returns timing information for the latest build of the given session.
func (s *Service) GetSessionTimings(ctx context.Context, id uuid.UUID) (SessionBuildTimings, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessions/%s/timings", id), nil)
	if err != nil {
		return SessionBuildTimings{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionBuildTimings{}, client.ReadBodyAsError(resp)
	}

	var timings SessionBuildTimings
	if err := json.NewDecoder(resp.Body).Decode(&timings); err != nil {
		return SessionBuildTimings{}, fmt.Errorf("decode session timings: %w", err)
	}
	return timings, nil
}
