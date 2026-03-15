// Package tasks provides task management services for the Lattice Runtime API.
package tasks

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	// TaskStatusPending indicates the task has been created but no agent
	// has been provisioned yet, or the agent build job status is unknown.
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusInitializing indicates the agent build is pending/running,
	// the sidecar is connecting, or apps are initializing.
	TaskStatusInitializing TaskStatus = "initializing"
	// TaskStatusActive indicates the task's agent is running with a
	// successful start transition, the sidecar is connected, and all agent
	// apps are either healthy or disabled.
	TaskStatusActive TaskStatus = "active"
	// TaskStatusPaused indicates the task's agent has been stopped or
	// deleted (stop/delete transition with successful job status).
	TaskStatusPaused TaskStatus = "paused"
	// TaskStatusUnknown indicates the task's status cannot be determined
	// based on the agent build, sidecar lifecycle, or app health states.
	TaskStatusUnknown TaskStatus = "unknown"
	// TaskStatusError indicates the task's agent build job has failed,
	// or the agent apps are reporting unhealthy status.
	TaskStatusError TaskStatus = "error"
)

// TaskState represents the high-level lifecycle of a task.
type TaskState string

const (
	// TaskStateWorking indicates the AI sidecar is actively processing work.
	TaskStateWorking TaskState = "working"
	// TaskStateIdle indicates the AI sidecar's screen is stable and no work
	// is being performed.
	TaskStateIdle TaskState = "idle"
	// TaskStateComplete indicates the AI sidecar has successfully completed
	// the task.
	TaskStateComplete TaskState = "complete"
	// TaskStateFailed indicates the AI sidecar reported a failure state.
	TaskStateFailed TaskState = "failed"
)

// TaskStateEntry represents a single entry in the task's state history.
type TaskStateEntry struct {
	Timestamp time.Time `json:"timestamp"`
	State     TaskState `json:"state"`
	Message   string    `json:"message"`
	URI       string    `json:"uri"`
}

// Task represents a task in the Lattice Runtime.
type Task struct {
	ID                  uuid.UUID       `json:"id"`
	OrganizationID      uuid.UUID       `json:"organization_id"`
	OwnerID             uuid.UUID       `json:"owner_id"`
	OwnerName           string          `json:"owner_name"`
	OwnerAvatarURL      string          `json:"owner_avatar_url,omitempty"`
	Name                string          `json:"name"`
	DisplayName         string          `json:"display_name"`
	TemplateID          uuid.UUID       `json:"template_id"`
	TemplateVersionID   uuid.UUID       `json:"template_version_id"`
	TemplateName        string          `json:"template_name"`
	TemplateDisplayName string          `json:"template_display_name"`
	TemplateIcon        string          `json:"template_icon"`
	AgentID             *uuid.UUID      `json:"agent_id,omitempty"`
	AgentName           string          `json:"agent_name"`
	AgentStatus         string          `json:"agent_status,omitempty"`
	AgentBuildNumber    int32           `json:"agent_build_number,omitempty"`
	AgentSidecarID      *uuid.UUID      `json:"agent_sidecar_id,omitempty"`
	InitialPrompt       string          `json:"initial_prompt"`
	Status              TaskStatus      `json:"status"`
	CurrentState        *TaskStateEntry `json:"current_state,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

// TasksFilter filters the list of tasks.
type TasksFilter struct {
	// Owner can be a username, UUID, or "me".
	Owner string `json:"owner,omitempty"`
	// Organization can be an organization name or UUID.
	Organization string `json:"organization,omitempty"`
	// Status filters the tasks by their task status.
	Status TaskStatus `json:"status,omitempty"`
	// FilterQuery allows specifying a raw filter query.
	FilterQuery string `json:"filter_query,omitempty"`
}

// CreateTaskRequest represents the request to create a new task.
type CreateTaskRequest struct {
	TemplateVersionID       uuid.UUID `json:"template_version_id"`
	TemplateVersionPresetID uuid.UUID `json:"template_version_preset_id,omitempty"`
	Input                   string    `json:"input"`
	Name                    string    `json:"name,omitempty"`
	DisplayName             string    `json:"display_name,omitempty"`
}

// TaskSendRequest is used to send task input to the tasks sidebar app.
type TaskSendRequest struct {
	Input string `json:"input"`
}

// UpdateTaskInputRequest is used to update a task's input.
type UpdateTaskInputRequest struct {
	Input string `json:"input"`
}

// TaskLogType indicates the source of a task log entry.
type TaskLogType string

const (
	// TaskLogTypeInput represents a user input log entry.
	TaskLogTypeInput TaskLogType = "input"
	// TaskLogTypeOutput represents a task output log entry.
	TaskLogTypeOutput TaskLogType = "output"
)

// TaskLogEntry represents a single log entry for a task.
type TaskLogEntry struct {
	ID      int         `json:"id"`
	Content string      `json:"content"`
	Type    TaskLogType `json:"type"`
	Time    time.Time   `json:"time"`
}

// TaskLogsResponse contains the logs for a task.
type TaskLogsResponse struct {
	Logs []TaskLogEntry `json:"logs"`
}
