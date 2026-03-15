package deployments

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/latticehq/latticesdk/client"
)

// GetDAUs returns daily active user metrics for the deployment.
// The tzOffset parameter is the timezone hour offset (e.g., 0 for UTC, -7 for PDT).
// Use TimezoneOffsetHour to compute the offset from a *time.Location.
func (s *Service) GetDAUs(ctx context.Context, tzOffset int) (*DAUsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/insights/daus", nil,
		client.WithQueryParam("tz_offset", strconv.Itoa(tzOffset)),
	)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
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

// GetDAUsLocalTZ returns daily active user metrics using the local timezone offset.
func (s *Service) GetDAUsLocalTZ(ctx context.Context) (*DAUsResponse, error) {
	return s.GetDAUs(ctx, TimezoneOffsetHour(time.Local))
}

// TimezoneOffsetHour returns the timezone offset in hours for the given location,
// matching the JavaScript getTimezoneOffset() convention. For example, UTC-7 returns 7.
func TimezoneOffsetHour(loc *time.Location) int {
	return TimezoneOffsetHourWithTime(time.Now(), loc)
}

// TimezoneOffsetHourWithTime returns the timezone offset in hours for the given
// location at the specified time. This accounts for daylight saving time changes.
func TimezoneOffsetHourWithTime(now time.Time, loc *time.Location) int {
	if loc == nil {
		loc = time.UTC
	}
	_, offsetSec := now.In(loc).Zone()
	return -1 * offsetSec / 60 / 60
}
