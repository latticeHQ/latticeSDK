package evals

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides evaluation run and comparison operations against the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new evals Service backed by the given client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// CreateEvalRun creates a new eval run in the specified organization.
func (s *Service) CreateEvalRun(ctx context.Context, org string, req CreateEvalRunRequest) (EvalRun, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/evals", org), req)
	if err != nil {
		return EvalRun{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return EvalRun{}, client.ReadBodyAsError(resp)
	}

	var evalRun EvalRun
	if err := json.NewDecoder(resp.Body).Decode(&evalRun); err != nil {
		return EvalRun{}, fmt.Errorf("decode response: %w", err)
	}
	return evalRun, nil
}

// ListEvalRuns returns eval runs in the specified organization.
func (s *Service) ListEvalRuns(ctx context.Context, org string, filter EvalRunFilter) (EvalRunsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/evals", org), nil,
		filterToRequestOption(filter))
	if err != nil {
		return EvalRunsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EvalRunsResponse{}, client.ReadBodyAsError(resp)
	}

	var result EvalRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return EvalRunsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// GetEvalRun returns a specific eval run by ID.
func (s *Service) GetEvalRun(ctx context.Context, id uuid.UUID) (EvalRun, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/evals/%s", id), nil)
	if err != nil {
		return EvalRun{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EvalRun{}, client.ReadBodyAsError(resp)
	}

	var evalRun EvalRun
	if err := json.NewDecoder(resp.Body).Decode(&evalRun); err != nil {
		return EvalRun{}, fmt.Errorf("decode response: %w", err)
	}
	return evalRun, nil
}

// DeleteEvalRun deletes an eval run by ID.
func (s *Service) DeleteEvalRun(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/evals/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetEvalRunPasses returns the passes for a specific eval run.
func (s *Service) GetEvalRunPasses(ctx context.Context, id uuid.UUID) ([]EvalRunPass, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/evals/%s/passes", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var passes []EvalRunPass
	if err := json.NewDecoder(resp.Body).Decode(&passes); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return passes, nil
}

// GetEvalRunParameters returns the parameters for a specific eval run.
func (s *Service) GetEvalRunParameters(ctx context.Context, id uuid.UUID) ([]EvalRunParameter, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/evals/%s/parameters", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var params []EvalRunParameter
	if err := json.NewDecoder(resp.Body).Decode(&params); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return params, nil
}

// CreateEvalComparison creates a new eval comparison in the specified organization.
func (s *Service) CreateEvalComparison(ctx context.Context, org string, req CreateEvalComparisonRequest) (EvalComparison, error) {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/organizations/%s/eval-comparisons", org), req)
	if err != nil {
		return EvalComparison{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return EvalComparison{}, client.ReadBodyAsError(resp)
	}

	var comparison EvalComparison
	if err := json.NewDecoder(resp.Body).Decode(&comparison); err != nil {
		return EvalComparison{}, fmt.Errorf("decode response: %w", err)
	}
	return comparison, nil
}

// ListEvalComparisons returns eval comparisons in the specified organization.
func (s *Service) ListEvalComparisons(ctx context.Context, org string) (EvalComparisonsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/organizations/%s/eval-comparisons", org), nil)
	if err != nil {
		return EvalComparisonsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EvalComparisonsResponse{}, client.ReadBodyAsError(resp)
	}

	var result EvalComparisonsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return EvalComparisonsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

// GetEvalComparison returns a specific eval comparison with all its runs.
func (s *Service) GetEvalComparison(ctx context.Context, id uuid.UUID) (EvalComparison, error) {
	resp, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/eval-comparisons/%s", id), nil)
	if err != nil {
		return EvalComparison{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return EvalComparison{}, client.ReadBodyAsError(resp)
	}

	var comparison EvalComparison
	if err := json.NewDecoder(resp.Body).Decode(&comparison); err != nil {
		return EvalComparison{}, fmt.Errorf("decode response: %w", err)
	}
	return comparison, nil
}

// DeleteEvalComparison deletes an eval comparison by ID.
func (s *Service) DeleteEvalComparison(ctx context.Context, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete,
		fmt.Sprintf("/api/v2/eval-comparisons/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// AddEvalComparisonRun adds an eval run to a comparison.
func (s *Service) AddEvalComparisonRun(ctx context.Context, comparisonID uuid.UUID, req AddEvalComparisonRunRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost,
		fmt.Sprintf("/api/v2/eval-comparisons/%s/runs", comparisonID), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// filterToRequestOption converts an EvalRunFilter into a client.RequestOption
// that sets query parameters on the outgoing request.
func filterToRequestOption(f EvalRunFilter) client.RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		if f.EvalTemplateID != uuid.Nil {
			q.Set("eval_template_id", f.EvalTemplateID.String())
		}
		if f.TargetSessionID != uuid.Nil {
			q.Set("target_session_id", f.TargetSessionID.String())
		}
		if f.InitiatorID != uuid.Nil {
			q.Set("initiator_id", f.InitiatorID.String())
		}
		if f.Status != "" {
			q.Set("status", string(f.Status))
		}
		r.URL.RawQuery = q.Encode()
	}
}
