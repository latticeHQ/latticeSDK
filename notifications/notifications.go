// Package notifications provides notification management for the Lattice Runtime API.
package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides notification management operations.
type Service struct {
	client *client.Client
}

// New creates a new notifications Service.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// GetSettings retrieves the global notifications settings, which currently
// describes whether all notifications are paused from sending.
func (s *Service) GetSettings(ctx context.Context) (NotificationsSettings, error) {
	res, err := s.client.Request(ctx, http.MethodGet, "/api/v2/notifications/settings", nil)
	if err != nil {
		return NotificationsSettings{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return NotificationsSettings{}, client.ReadBodyAsError(res)
	}
	var settings NotificationsSettings
	return settings, json.NewDecoder(res.Body).Decode(&settings)
}

// UpdateSettings modifies the global notifications settings.
func (s *Service) UpdateSettings(ctx context.Context, settings NotificationsSettings) error {
	res, err := s.client.Request(ctx, http.MethodPut, "/api/v2/notifications/settings", settings)
	if err != nil {
		return fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(res)
	}
	return nil
}

// UpdateTemplateMethod modifies a notification template to use a specific
// dispatch method, overriding the method set in the deployment configuration.
func (s *Service) UpdateTemplateMethod(ctx context.Context, templateID uuid.UUID, method string) error {
	res, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/notifications/templates/%s/method", templateID),
		updateNotificationTemplateMethod{Method: method},
	)
	if err != nil {
		return fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return client.ReadBodyAsError(res)
	}
	return nil
}

// ListSystemTemplates retrieves all notification templates pertaining to
// internal system events.
func (s *Service) ListSystemTemplates(ctx context.Context) ([]NotificationTemplate, error) {
	res, err := s.client.Request(ctx, http.MethodGet, "/api/v2/notifications/templates/system", nil)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var templates []NotificationTemplate
	return templates, json.NewDecoder(res.Body).Decode(&templates)
}

// GetUserPreferences retrieves notification preferences for the specified user.
func (s *Service) GetUserPreferences(ctx context.Context, userID uuid.UUID) ([]NotificationPreference, error) {
	res, err := s.client.Request(ctx, http.MethodGet,
		fmt.Sprintf("/api/v2/users/%s/notifications/preferences", userID),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var prefs []NotificationPreference
	return prefs, json.NewDecoder(res.Body).Decode(&prefs)
}

// UpdateUserPreferences updates notification preferences for the specified user.
func (s *Service) UpdateUserPreferences(ctx context.Context, userID uuid.UUID, req UpdateUserNotificationPreferences) ([]NotificationPreference, error) {
	res, err := s.client.Request(ctx, http.MethodPut,
		fmt.Sprintf("/api/v2/users/%s/notifications/preferences", userID),
		req,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(res)
	}
	var prefs []NotificationPreference
	return prefs, json.NewDecoder(res.Body).Decode(&prefs)
}

// GetDispatchMethods returns the available and default notification dispatch
// methods for the deployment.
func (s *Service) GetDispatchMethods(ctx context.Context) (NotificationMethodsResponse, error) {
	res, err := s.client.Request(ctx, http.MethodGet, "/api/v2/notifications/dispatch-methods", nil)
	if err != nil {
		return NotificationMethodsResponse{}, fmt.Errorf("make request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return NotificationMethodsResponse{}, client.ReadBodyAsError(res)
	}
	var resp NotificationMethodsResponse
	return resp, json.NewDecoder(res.Body).Decode(&resp)
}
