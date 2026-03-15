package insights

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// GetSessionAnalytics returns session analytics for the specified organization.
func (s *Service) GetSessionAnalytics(ctx context.Context, orgID uuid.UUID) (SessionAnalyticsResponse, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/analytics/sessions", orgID),
		nil,
	)
	if err != nil {
		return SessionAnalyticsResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return SessionAnalyticsResponse{}, client.ReadBodyAsError(res)
	}
	var resp SessionAnalyticsResponse
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetSessionCountsByTemplate returns session counts grouped by template for
// the specified organization.
func (s *Service) GetSessionCountsByTemplate(ctx context.Context, orgID uuid.UUID) ([]SessionCountsByTemplate, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/analytics/sessions/by-template", orgID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var resp []SessionCountsByTemplate
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetEvalAnalytics returns eval analytics for the specified organization.
func (s *Service) GetEvalAnalytics(ctx context.Context, orgID uuid.UUID) (EvalAnalyticsResponse, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/analytics/evals", orgID),
		nil,
	)
	if err != nil {
		return EvalAnalyticsResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return EvalAnalyticsResponse{}, client.ReadBodyAsError(res)
	}
	var resp EvalAnalyticsResponse
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetEvalScoresByTemplate returns eval counts grouped by template for the
// specified organization.
func (s *Service) GetEvalScoresByTemplate(ctx context.Context, orgID uuid.UUID) ([]EvalScoresByTemplate, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/analytics/evals/by-template", orgID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var resp []EvalScoresByTemplate
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}
