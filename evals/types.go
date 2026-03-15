// Package evals provides evaluation run and comparison services for the Lattice Runtime API.
package evals

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EvalRunStatus represents the status of an eval run.
type EvalRunStatus string

const (
	// EvalRunStatusPending indicates the eval run has been created but not started.
	EvalRunStatusPending EvalRunStatus = "pending"
	// EvalRunStatusRunning indicates the eval run is in progress.
	EvalRunStatusRunning EvalRunStatus = "running"
	// EvalRunStatusCompleted indicates the eval run finished successfully.
	EvalRunStatusCompleted EvalRunStatus = "completed"
	// EvalRunStatusFailed indicates the eval run encountered an error.
	EvalRunStatusFailed EvalRunStatus = "failed"
	// EvalRunStatusCanceled indicates the eval run was canceled.
	EvalRunStatusCanceled EvalRunStatus = "canceled"
)

// EvalComparisonStatus represents the status of an eval comparison.
type EvalComparisonStatus string

const (
	// EvalComparisonStatusPending indicates the comparison is awaiting execution.
	EvalComparisonStatusPending EvalComparisonStatus = "pending"
	// EvalComparisonStatusRunning indicates the comparison is in progress.
	EvalComparisonStatusRunning EvalComparisonStatus = "running"
	// EvalComparisonStatusCompleted indicates the comparison finished successfully.
	EvalComparisonStatusCompleted EvalComparisonStatus = "completed"
	// EvalComparisonStatusFailed indicates the comparison encountered an error.
	EvalComparisonStatusFailed EvalComparisonStatus = "failed"
)

// EvalComparisonType represents the type of comparison.
type EvalComparisonType string

const (
	// EvalComparisonTypeModel compares across different models.
	EvalComparisonTypeModel EvalComparisonType = "model"
	// EvalComparisonTypePreset compares across different presets.
	EvalComparisonTypePreset EvalComparisonType = "preset"
	// EvalComparisonTypeUser compares across different users.
	EvalComparisonTypeUser EvalComparisonType = "user"
	// EvalComparisonTypeCustom is a user-defined comparison type.
	EvalComparisonTypeCustom EvalComparisonType = "custom"
)

// EvalRun represents a single evaluation of a session or transcript.
type EvalRun struct {
	ID                    uuid.UUID       `json:"id"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	OrganizationID        uuid.UUID       `json:"organization_id"`
	EvalTemplateID        uuid.UUID       `json:"eval_template_id"`
	EvalTemplateVersionID uuid.UUID       `json:"eval_template_version_id"`
	TargetSessionID       *uuid.UUID      `json:"target_session_id,omitempty"`
	Transcript            string          `json:"transcript,omitempty"`
	SessionIDExternal     string          `json:"session_id_external,omitempty"`
	Status                EvalRunStatus   `json:"status"`
	InitiatorID           uuid.UUID       `json:"initiator_id"`
	Results               json.RawMessage `json:"results,omitempty"`
	Scores                json.RawMessage `json:"scores,omitempty"`
	Metadata              json.RawMessage `json:"metadata,omitempty"`

	// Populated via joins.
	EvalTemplateName    string `json:"eval_template_name,omitempty"`
	InitiatorUsername   string `json:"initiator_username,omitempty"`
	SessionTemplateName string `json:"session_template_name,omitempty"`
}

// EvalRunParameter stores a resolved parameter value for an eval run.
type EvalRunParameter struct {
	EvalRunID uuid.UUID `json:"eval_run_id"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
}

// EvalRunPass represents a single analysis pass within an eval run.
type EvalRunPass struct {
	ID          uuid.UUID       `json:"id"`
	EvalRunID   uuid.UUID       `json:"eval_run_id"`
	PassNumber  int32           `json:"pass_number"`
	PassName    string          `json:"pass_name"`
	Fields      string          `json:"fields"`
	Instruction string          `json:"instruction,omitempty"`
	Status      EvalRunStatus   `json:"status"`
	Results     json.RawMessage `json:"results,omitempty"`
	TokensUsed  int32           `json:"tokens_used"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

// EvalRunFilter allows filtering eval runs.
type EvalRunFilter struct {
	EvalTemplateID  uuid.UUID     `json:"eval_template_id,omitempty"`
	TargetSessionID uuid.UUID     `json:"target_session_id,omitempty"`
	InitiatorID     uuid.UUID     `json:"initiator_id,omitempty"`
	Status          EvalRunStatus `json:"status,omitempty"`
}

// EvalRunsResponse contains a list of eval runs with a total count.
type EvalRunsResponse struct {
	EvalRuns []EvalRun `json:"eval_runs"`
	Count    int       `json:"count"`
}

// EvalRunParameterInput is a parameter name/value pair for eval run creation.
type EvalRunParameterInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateEvalRunRequest is the request body for creating a new eval run.
type CreateEvalRunRequest struct {
	EvalTemplateID        uuid.UUID             `json:"eval_template_id"`
	EvalTemplateVersionID uuid.UUID             `json:"eval_template_version_id,omitempty"`
	TargetSessionID       *uuid.UUID            `json:"target_session_id,omitempty"`
	Transcript            string                `json:"transcript,omitempty"`
	SessionIDExternal     string                `json:"session_id_external,omitempty"`
	Parameters            []EvalRunParameterInput `json:"parameters,omitempty"`
	Metadata              json.RawMessage       `json:"metadata,omitempty"`
}

// EvalComparison represents a comparison between multiple eval runs.
type EvalComparison struct {
	ID             uuid.UUID            `json:"id"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
	OrganizationID uuid.UUID            `json:"organization_id"`
	EvalTemplateID uuid.UUID            `json:"eval_template_id"`
	Name           string               `json:"name"`
	ComparisonType EvalComparisonType   `json:"comparison_type"`
	Status         EvalComparisonStatus `json:"status"`
	Results        json.RawMessage      `json:"results,omitempty"`
	Metadata       json.RawMessage      `json:"metadata,omitempty"`

	// Runs is populated when fetching a comparison with its runs.
	Runs []EvalComparisonRunEntry `json:"runs,omitempty"`
}

// EvalComparisonRunEntry is an eval run with its comparison label.
type EvalComparisonRunEntry struct {
	Label   string  `json:"label"`
	EvalRun EvalRun `json:"eval_run"`
}

// CreateEvalComparisonRequest is the request body for creating an eval comparison.
type CreateEvalComparisonRequest struct {
	EvalTemplateID uuid.UUID          `json:"eval_template_id"`
	Name           string             `json:"name"`
	ComparisonType EvalComparisonType `json:"comparison_type,omitempty"`
	Metadata       json.RawMessage    `json:"metadata,omitempty"`
}

// AddEvalComparisonRunRequest adds an eval run to a comparison.
type AddEvalComparisonRunRequest struct {
	EvalRunID uuid.UUID `json:"eval_run_id"`
	Label     string    `json:"label,omitempty"`
}

// EvalComparisonsResponse contains a list of eval comparisons with a total count.
type EvalComparisonsResponse struct {
	Comparisons []EvalComparison `json:"comparisons"`
	Count       int              `json:"count"`
}
