package insights

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// InsightsReportInterval is the interval of time over which to generate a
// smaller insights report within a time range.
type InsightsReportInterval string

const (
	// InsightsReportIntervalDay generates daily interval reports.
	InsightsReportIntervalDay InsightsReportInterval = "day"
	// InsightsReportIntervalWeek generates weekly interval reports.
	InsightsReportIntervalWeek InsightsReportInterval = "week"
)

// Days returns the duration of the interval in days.
func (interval InsightsReportInterval) Days() int32 {
	switch interval {
	case InsightsReportIntervalDay:
		return 1
	case InsightsReportIntervalWeek:
		return 7
	default:
		panic("developer error: unsupported report interval")
	}
}

// TemplateInsightsSection defines the section to be included in the template
// insights response.
type TemplateInsightsSection string

const (
	// TemplateInsightsSectionIntervalReports requests interval-level reports.
	TemplateInsightsSectionIntervalReports TemplateInsightsSection = "interval_reports"
	// TemplateInsightsSectionReport requests the aggregate report.
	TemplateInsightsSectionReport TemplateInsightsSection = "report"
)

// --- User Latency ---

// UserLatencyInsightsRequest contains the parameters for querying user latency
// insights.
type UserLatencyInsightsRequest struct {
	StartTime   time.Time   `json:"start_time" format:"date-time"`
	EndTime     time.Time   `json:"end_time" format:"date-time"`
	TemplateIDs []uuid.UUID `json:"template_ids" format:"uuid"`
}

// UserLatencyInsightsResponse is the response from the user latency insights
// endpoint.
type UserLatencyInsightsResponse struct {
	Report UserLatencyInsightsReport `json:"report"`
}

// UserLatencyInsightsReport is the report from the user latency insights
// endpoint.
type UserLatencyInsightsReport struct {
	StartTime   time.Time     `json:"start_time" format:"date-time"`
	EndTime     time.Time     `json:"end_time" format:"date-time"`
	TemplateIDs []uuid.UUID   `json:"template_ids" format:"uuid"`
	Users       []UserLatency `json:"users"`
}

// UserLatency shows the connection latency for a user.
type UserLatency struct {
	TemplateIDs []uuid.UUID       `json:"template_ids" format:"uuid"`
	UserID      uuid.UUID         `json:"user_id" format:"uuid"`
	Username    string            `json:"username"`
	AvatarURL   string            `json:"avatar_url" format:"uri"`
	LatencyMS   ConnectionLatency `json:"latency_ms"`
}

// ConnectionLatency shows latency percentiles for a connection.
type ConnectionLatency struct {
	P50 float64 `json:"p50"`
	P95 float64 `json:"p95"`
}

// --- User Activity ---

// UserActivityInsightsRequest contains the parameters for querying user
// activity insights.
type UserActivityInsightsRequest struct {
	StartTime   time.Time   `json:"start_time" format:"date-time"`
	EndTime     time.Time   `json:"end_time" format:"date-time"`
	TemplateIDs []uuid.UUID `json:"template_ids" format:"uuid"`
}

// UserActivityInsightsResponse is the response from the user activity insights
// endpoint.
type UserActivityInsightsResponse struct {
	Report UserActivityInsightsReport `json:"report"`
}

// UserActivityInsightsReport is the report from the user activity insights
// endpoint.
type UserActivityInsightsReport struct {
	StartTime   time.Time      `json:"start_time" format:"date-time"`
	EndTime     time.Time      `json:"end_time" format:"date-time"`
	TemplateIDs []uuid.UUID    `json:"template_ids" format:"uuid"`
	Users       []UserActivity `json:"users"`
}

// UserActivity shows the session time for a user.
type UserActivity struct {
	TemplateIDs []uuid.UUID `json:"template_ids" format:"uuid"`
	UserID      uuid.UUID   `json:"user_id" format:"uuid"`
	Username    string      `json:"username"`
	AvatarURL   string      `json:"avatar_url" format:"uri"`
	Seconds     int64       `json:"seconds"`
}

// --- Template Insights ---

// TemplateInsightsRequest contains the parameters for querying template
// insights.
type TemplateInsightsRequest struct {
	StartTime   time.Time                 `json:"start_time" format:"date-time"`
	EndTime     time.Time                 `json:"end_time" format:"date-time"`
	TemplateIDs []uuid.UUID               `json:"template_ids" format:"uuid"`
	Interval    InsightsReportInterval    `json:"interval"`
	Sections    []TemplateInsightsSection `json:"sections"`
}

// TemplateInsightsResponse is the response from the template insights endpoint.
type TemplateInsightsResponse struct {
	Report          *TemplateInsightsReport          `json:"report,omitempty"`
	IntervalReports []TemplateInsightsIntervalReport `json:"interval_reports,omitempty"`
}

// TemplateInsightsReport is the aggregate report from the template insights
// endpoint.
type TemplateInsightsReport struct {
	StartTime       time.Time                `json:"start_time" format:"date-time"`
	EndTime         time.Time                `json:"end_time" format:"date-time"`
	TemplateIDs     []uuid.UUID              `json:"template_ids" format:"uuid"`
	ActiveUsers     int64                    `json:"active_users"`
	AppsUsage       []TemplateAppUsage       `json:"apps_usage"`
	ParametersUsage []TemplateParameterUsage `json:"parameters_usage"`
}

// TemplateInsightsIntervalReport is the report from the template insights
// endpoint for a specific interval.
type TemplateInsightsIntervalReport struct {
	StartTime   time.Time              `json:"start_time" format:"date-time"`
	EndTime     time.Time              `json:"end_time" format:"date-time"`
	TemplateIDs []uuid.UUID            `json:"template_ids" format:"uuid"`
	Interval    InsightsReportInterval `json:"interval"`
	ActiveUsers int64                  `json:"active_users"`
}

// TemplateAppsType defines the type of app reported.
type TemplateAppsType string

const (
	// TemplateAppsTypeBuiltin represents a built-in application.
	TemplateAppsTypeBuiltin TemplateAppsType = "builtin"
	// TemplateAppsTypeApp represents a user-defined application.
	TemplateAppsTypeApp TemplateAppsType = "app"
)

// TemplateAppUsage shows the usage of an app for one or more templates.
type TemplateAppUsage struct {
	TemplateIDs []uuid.UUID      `json:"template_ids" format:"uuid"`
	Type        TemplateAppsType `json:"type"`
	DisplayName string           `json:"display_name"`
	Slug        string           `json:"slug"`
	Icon        string           `json:"icon"`
	Seconds     int64            `json:"seconds"`
	TimesUsed   int64            `json:"times_used"`
}

// TemplateParameterUsage shows the usage of a parameter for one or more
// templates.
type TemplateParameterUsage struct {
	TemplateIDs []uuid.UUID                    `json:"template_ids" format:"uuid"`
	DisplayName string                         `json:"display_name"`
	Name        string                         `json:"name"`
	Type        string                         `json:"type"`
	Description string                         `json:"description"`
	Options     []TemplateVersionParameterOption `json:"options,omitempty"`
	Values      []TemplateParameterValue       `json:"values"`
}

// TemplateVersionParameterOption represents an option for a template version
// parameter.
type TemplateVersionParameterOption struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	Icon        string `json:"icon"`
}

// TemplateParameterValue shows the usage count of a parameter value.
type TemplateParameterValue struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// --- Session Analytics ---

// SessionAnalytics represents a session row from the analytics view,
// enriched with template and owner metadata plus aggregated parameters.
type SessionAnalytics struct {
	ID                uuid.UUID       `json:"id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	OrganizationID    uuid.UUID       `json:"organization_id"`
	TemplateID        uuid.UUID       `json:"template_id"`
	TemplateVersionID uuid.UUID       `json:"template_version_id"`
	OwnerID           uuid.UUID       `json:"owner_id"`
	Name              string          `json:"name"`
	DisplayName       string          `json:"display_name"`
	Status            string          `json:"status"`
	StartedAt         *time.Time      `json:"started_at,omitempty"`
	EndedAt           *time.Time      `json:"ended_at,omitempty"`
	DurationSeconds   *int32          `json:"duration_seconds,omitempty"`
	TemplateName      string          `json:"template_name"`
	TemplateCategory  string          `json:"template_category"`
	ParentTemplateID  *uuid.UUID      `json:"parent_template_id,omitempty"`
	OwnerUsername     string          `json:"owner_username"`
	Parameters        json.RawMessage `json:"parameters"`
}

// SessionCountsByTemplate holds aggregated session counts per template.
type SessionCountsByTemplate struct {
	TemplateID         uuid.UUID `json:"template_id"`
	TemplateName       string    `json:"template_name"`
	TemplateCategory   string    `json:"template_category"`
	SessionCount       int32     `json:"session_count"`
	ActiveCount        int32     `json:"active_count"`
	CompletedCount     int32     `json:"completed_count"`
	AvgDurationSeconds int32     `json:"avg_duration_seconds"`
}

// SessionAnalyticsResponse wraps the session analytics results.
type SessionAnalyticsResponse struct {
	Sessions []SessionAnalytics `json:"sessions"`
	Count    int                `json:"count"`
}

// --- Eval Analytics ---

// EvalAnalytics represents an eval run from the analytics view,
// enriched with template metadata, session info, and aggregated parameters.
type EvalAnalytics struct {
	ID                    uuid.UUID       `json:"id"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	OrganizationID        uuid.UUID       `json:"organization_id"`
	EvalTemplateID        uuid.UUID       `json:"eval_template_id"`
	EvalTemplateVersionID uuid.UUID       `json:"eval_template_version_id"`
	TargetSessionID       *uuid.UUID      `json:"target_session_id,omitempty"`
	Status                string          `json:"status"`
	InitiatorID           uuid.UUID       `json:"initiator_id"`
	Results               json.RawMessage `json:"results"`
	Scores                json.RawMessage `json:"scores"`
	EvalTemplateName      string          `json:"eval_template_name"`
	SessionTemplateID     *uuid.UUID      `json:"session_template_id,omitempty"`
	SessionTemplateName   string          `json:"session_template_name,omitempty"`
	InitiatorUsername     string          `json:"initiator_username"`
	SessionParameters     json.RawMessage `json:"session_parameters,omitempty"`
	EvalParameters        json.RawMessage `json:"eval_parameters,omitempty"`
}

// EvalScoresByTemplate holds aggregated eval counts per template.
type EvalScoresByTemplate struct {
	EvalTemplateID   uuid.UUID `json:"eval_template_id"`
	EvalTemplateName string    `json:"eval_template_name"`
	RunCount         int32     `json:"run_count"`
	CompletedCount   int32     `json:"completed_count"`
	FailedCount      int32     `json:"failed_count"`
}

// EvalAnalyticsResponse wraps the eval analytics results.
type EvalAnalyticsResponse struct {
	EvalRuns []EvalAnalytics `json:"eval_runs"`
	Count    int             `json:"count"`
}
