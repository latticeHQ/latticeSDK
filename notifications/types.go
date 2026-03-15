package notifications

import (
	"time"

	"github.com/google/uuid"
)

// NotificationsSettings describes the global notification configuration.
type NotificationsSettings struct {
	// NotifierPaused indicates whether all notifications are paused from sending.
	NotifierPaused bool `json:"notifier_paused"`
}

// NotificationTemplate describes a system notification template.
type NotificationTemplate struct {
	ID            uuid.UUID `json:"id" format:"uuid"`
	Name          string    `json:"name"`
	TitleTemplate string    `json:"title_template"`
	BodyTemplate  string    `json:"body_template"`
	Actions       string    `json:"actions"`
	Group         string    `json:"group"`
	Method        string    `json:"method"`
	Kind          string    `json:"kind"`
}

// NotificationPreference represents a user's preference for a specific
// notification template.
type NotificationPreference struct {
	NotificationTemplateID uuid.UUID `json:"id" format:"uuid"`
	Disabled               bool      `json:"disabled"`
	UpdatedAt              time.Time `json:"updated_at" format:"date-time"`
}

// UpdateUserNotificationPreferences contains the map of template IDs to their
// disabled state for updating a user's notification preferences.
type UpdateUserNotificationPreferences struct {
	TemplateDisabledMap map[string]bool `json:"template_disabled_map"`
}

// NotificationMethodsResponse describes the available and default notification
// dispatch methods.
type NotificationMethodsResponse struct {
	AvailableNotificationMethods []string `json:"available"`
	DefaultNotificationMethod    string   `json:"default"`
}

// updateNotificationTemplateMethod is the request body for updating a
// template's dispatch method.
type updateNotificationTemplateMethod struct {
	Method string `json:"method,omitempty"`
}
