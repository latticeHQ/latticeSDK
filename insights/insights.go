// Package insights provides usage insights and analytics for the Lattice Runtime API.
package insights

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/latticehq/latticesdk/client"
)

// insightsTimeLayout is the time format used by the insights API.
const insightsTimeLayout = time.RFC3339

// Service provides insights and analytics operations.
type Service struct {
	client *client.Client
}

// New creates a new insights Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// GetUserLatencyInsights returns connection latency insights for users over
// the specified time range and optional template filter.
func (s *Service) GetUserLatencyInsights(ctx context.Context, req UserLatencyInsightsRequest) (UserLatencyInsightsResponse, error) {
	qp := url.Values{}
	qp.Add("start_time", req.StartTime.Format(insightsTimeLayout))
	qp.Add("end_time", req.EndTime.Format(insightsTimeLayout))
	if len(req.TemplateIDs) > 0 {
		ids := make([]string, 0, len(req.TemplateIDs))
		for _, id := range req.TemplateIDs {
			ids = append(ids, id.String())
		}
		qp.Add("template_ids", strings.Join(ids, ","))
	}

	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/insights/user-latency?%s", qp.Encode()),
		nil,
	)
	if err != nil {
		return UserLatencyInsightsResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return UserLatencyInsightsResponse{}, client.ReadBodyAsError(res)
	}
	var result UserLatencyInsightsResponse
	return result, json.NewDecoder(res.Body).Decode(&result)
}

// GetUserActivityInsights returns session activity insights for users over
// the specified time range and optional template filter.
func (s *Service) GetUserActivityInsights(ctx context.Context, req UserActivityInsightsRequest) (UserActivityInsightsResponse, error) {
	qp := url.Values{}
	qp.Add("start_time", req.StartTime.Format(insightsTimeLayout))
	qp.Add("end_time", req.EndTime.Format(insightsTimeLayout))
	if len(req.TemplateIDs) > 0 {
		ids := make([]string, 0, len(req.TemplateIDs))
		for _, id := range req.TemplateIDs {
			ids = append(ids, id.String())
		}
		qp.Add("template_ids", strings.Join(ids, ","))
	}

	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/insights/user-activity?%s", qp.Encode()),
		nil,
	)
	if err != nil {
		return UserActivityInsightsResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return UserActivityInsightsResponse{}, client.ReadBodyAsError(res)
	}
	var result UserActivityInsightsResponse
	return result, json.NewDecoder(res.Body).Decode(&result)
}

// GetTemplateInsights returns template usage insights including app usage,
// parameter usage, and optional interval reports.
func (s *Service) GetTemplateInsights(ctx context.Context, req TemplateInsightsRequest) (TemplateInsightsResponse, error) {
	qp := url.Values{}
	qp.Add("start_time", req.StartTime.Format(insightsTimeLayout))
	qp.Add("end_time", req.EndTime.Format(insightsTimeLayout))
	if len(req.TemplateIDs) > 0 {
		ids := make([]string, 0, len(req.TemplateIDs))
		for _, id := range req.TemplateIDs {
			ids = append(ids, id.String())
		}
		qp.Add("template_ids", strings.Join(ids, ","))
	}
	if req.Interval != "" {
		qp.Add("interval", string(req.Interval))
	}
	if len(req.Sections) > 0 {
		sections := make([]string, 0, len(req.Sections))
		for _, sec := range req.Sections {
			sections = append(sections, string(sec))
		}
		qp.Add("sections", strings.Join(sections, ","))
	}

	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/insights/templates?%s", qp.Encode()),
		nil,
	)
	if err != nil {
		return TemplateInsightsResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return TemplateInsightsResponse{}, client.ReadBodyAsError(res)
	}
	var result TemplateInsightsResponse
	return result, json.NewDecoder(res.Body).Decode(&result)
}
