package sessions

import (
	"encoding/json"
	"fmt"
	"context"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// GetSessionBuild returns a single session build by its ID.
func (s *Service) GetSessionBuild(ctx context.Context, id uuid.UUID) (SessionBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionbuilds/%s", id), nil)
	if err != nil {
		return SessionBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionBuild{}, client.ReadBodyAsError(resp)
	}

	var build SessionBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return SessionBuild{}, fmt.Errorf("decode session build: %w", err)
	}
	return build, nil
}

// GetSessionBuildByOwnerAndName returns a session build by the owner's username,
// session name, and build number.
func (s *Service) GetSessionBuildByOwnerAndName(ctx context.Context, owner, session, buildNumber string) (SessionBuild, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/users/%s/session/%s/builds/%s", owner, session, buildNumber),
		nil,
	)
	if err != nil {
		return SessionBuild{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionBuild{}, client.ReadBodyAsError(resp)
	}

	var build SessionBuild
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return SessionBuild{}, fmt.Errorf("decode session build: %w", err)
	}
	return build, nil
}

// CancelSessionBuild marks a session build job as canceled.
func (s *Service) CancelSessionBuild(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/sessionbuilds/%s/cancel", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetSessionBuildParameters returns the build parameters for a session build.
func (s *Service) GetSessionBuildParameters(ctx context.Context, buildID uuid.UUID) ([]SessionBuildParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionbuilds/%s/parameters", buildID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []SessionBuildParameter
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		return nil, fmt.Errorf("decode session build parameters: %w", err)
	}
	return params, nil
}

// GetSessionBuildState returns the raw provisioner state of a build.
func (s *Service) GetSessionBuildState(ctx context.Context, buildID uuid.UUID) ([]byte, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionbuilds/%s/state", buildID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	return io.ReadAll(resp.Body)
}

// GetSessionBuildTimings returns timing information for a specific session build.
func (s *Service) GetSessionBuildTimings(ctx context.Context, buildID uuid.UUID) (SessionBuildTimings, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/sessionbuilds/%s/timings", buildID), nil)
	if err != nil {
		return SessionBuildTimings{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SessionBuildTimings{}, client.ReadBodyAsError(resp)
	}

	var timings SessionBuildTimings
	if err := json.NewDecoder(resp.Body).Decode(&timings); err != nil {
		return SessionBuildTimings{}, fmt.Errorf("decode session build timings: %w", err)
	}
	return timings, nil
}

// WatchSessionBuildLogs opens an SSE stream for provisioner job logs that
// occurred after the given log ID. The returned function yields the next
// server-sent event on each call. Returns (nil, io.EOF) when the stream ends.
// Callers should cancel the context to stop the stream.
func (s *Service) WatchSessionBuildLogs(ctx context.Context, buildID uuid.UUID, after int64) (func() (*client.ServerSentEvent, error), error) {
	path := fmt.Sprintf("/api/v2/sessionbuilds/%s/logs", buildID)
	if after > 0 {
		path = fmt.Sprintf("%s?after=%d&follow", path, after)
	} else {
		path = path + "?follow"
	}

	resp, err := s.client.Request(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	nextEvent := client.ServerSentEventReader(resp.Body)
	return nextEvent, nil
}
