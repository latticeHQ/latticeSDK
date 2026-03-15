package toolsdk

import "github.com/google/uuid"

// ---- Report Task ----

// ReportTaskArgs contains task status update information.
type ReportTaskArgs struct {
	Link    string `json:"link"`
	State   string `json:"state"`   // working, idle, failure, complete
	Summary string `json:"summary"`
}

// ReportTaskResult is the result of reporting task status.
type ReportTaskResult struct {
	OK bool `json:"ok"`
}

// ---- Agent management ----

// GetAgentArgs contains the parameters for getting an agent.
type GetAgentArgs struct {
	AgentID string `json:"agent_id"`
}

// CreateAgentArgs contains the parameters for creating an agent from a template.
type CreateAgentArgs struct {
	User              string            `json:"user"`
	TemplateVersionID string            `json:"template_version_id"`
	Name              string            `json:"name"`
	Parameters        map[string]string `json:"parameters,omitempty"`
}

// ListAgentsArgs contains the parameters for listing agents.
type ListAgentsArgs struct {
	Owner    string `json:"owner,omitempty"`
	Template string `json:"template,omitempty"`
	Status   string `json:"status,omitempty"`
}

// CreateAgentBuildArgs contains the parameters for creating an agent build.
type CreateAgentBuildArgs struct {
	AgentID           string `json:"agent_id"`
	TemplateVersionID string `json:"template_version_id,omitempty"`
	Transition        string `json:"transition"` // start, stop, delete
}

// ---- Template tools ----

// ListTemplatesArgs contains the parameters for listing templates.
type ListTemplatesArgs struct {
	OrganizationID string `json:"organization_id,omitempty"`
}

// TemplateVersionParametersArgs contains the parameters for getting template version parameters.
type TemplateVersionParametersArgs struct {
	TemplateVersionID string `json:"template_version_id"`
}

// CreateTemplateVersionArgs contains the parameters for creating a template version.
type CreateTemplateVersionArgs struct {
	OrganizationID string            `json:"organization_id"`
	TemplateID     string            `json:"template_id,omitempty"`
	Name           string            `json:"name,omitempty"`
	Message        string            `json:"message,omitempty"`
	FileID         string            `json:"file_id"`
	Provisioner    string            `json:"provisioner"`
	Tags           map[string]string `json:"tags,omitempty"`
}

// UpdateTemplateActiveVersionArgs contains the parameters for updating the active version.
type UpdateTemplateActiveVersionArgs struct {
	TemplateID string `json:"template_id"`
	VersionID  string `json:"version_id"`
}

// CreateTemplateArgs contains the parameters for creating a template.
type CreateTemplateArgs struct {
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	DisplayName    string `json:"display_name,omitempty"`
	Description    string `json:"description,omitempty"`
	VersionID      string `json:"version_id"`
}

// DeleteTemplateArgs contains the parameters for deleting a template.
type DeleteTemplateArgs struct {
	TemplateID string `json:"template_id"`
}

// UploadTarFileArgs contains the parameters for uploading a tar file.
type UploadTarFileArgs struct {
	// FileBase64 is the base64-encoded tar file content.
	FileBase64 string `json:"file_base64"`
}

// UploadTarFileResult contains the result of uploading a tar file.
type UploadTarFileResult struct {
	ID uuid.UUID `json:"id"`
}

// ---- User tools ----

// GetAuthenticatedUserArgs is empty; no parameters are required.
type GetAuthenticatedUserArgs struct{}

// ---- Task tools ----

// CreateTaskArgs contains the parameters for creating a task.
type CreateTaskArgs struct {
	User              string `json:"user"`
	TemplateVersionID string `json:"template_version_id"`
	Input             string `json:"input"`
	Name              string `json:"name,omitempty"`
}

// DeleteTaskArgs contains the parameters for deleting a task.
type DeleteTaskArgs struct {
	User   string `json:"user"`
	TaskID string `json:"task_id"`
}

// ListTasksArgs contains the parameters for listing tasks.
type ListTasksArgs struct {
	Owner  string `json:"owner,omitempty"`
	Status string `json:"status,omitempty"`
}

// GetTaskStatusArgs contains the parameters for getting task status.
type GetTaskStatusArgs struct {
	TaskIdentifier string `json:"task_identifier"`
}

// SendTaskInputArgs contains the parameters for sending input to a task.
type SendTaskInputArgs struct {
	User   string `json:"user"`
	TaskID string `json:"task_id"`
	Input  string `json:"input"`
}

// GetTaskLogsArgs contains the parameters for getting task logs.
type GetTaskLogsArgs struct {
	User   string `json:"user"`
	TaskID string `json:"task_id"`
}

// ---- Agent operation tools ----

// AgentBashArgs contains the parameters for executing a bash command in an agent.
type AgentBashArgs struct {
	AgentID   string `json:"agent_id"`
	SidecarID string `json:"sidecar_id"`
	Command   string `json:"command"`
}

// AgentBashResult contains the result of a bash command execution.
type AgentBashResult struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

// AgentLSArgs contains the parameters for listing files in an agent directory.
type AgentLSArgs struct {
	AgentID   string `json:"agent_id"`
	SidecarID string `json:"sidecar_id"`
	Path      string `json:"path"`
}

// AgentLSResult contains the result of listing files in an agent.
type AgentLSResult struct {
	Files []string `json:"files"`
}

// AgentReadFileArgs contains the parameters for reading a file from an agent.
type AgentReadFileArgs struct {
	AgentID   string `json:"agent_id"`
	SidecarID string `json:"sidecar_id"`
	Path      string `json:"path"`
}

// AgentReadFileResult contains the content of a file read from an agent.
type AgentReadFileResult struct {
	Content string `json:"content"`
}

// AgentWriteFileArgs contains the parameters for writing a file to an agent.
type AgentWriteFileArgs struct {
	AgentID   string `json:"agent_id"`
	SidecarID string `json:"sidecar_id"`
	Path      string `json:"path"`
	Content   string `json:"content"`
}

// AgentWriteFileResult contains the result of writing a file.
type AgentWriteFileResult struct {
	OK bool `json:"ok"`
}

// FileEdit describes a single edit operation on a file.
type FileEdit struct {
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
}

// AgentEditFileArgs contains the parameters for editing a file in an agent.
type AgentEditFileArgs struct {
	AgentID   string     `json:"agent_id"`
	SidecarID string     `json:"sidecar_id"`
	Path      string     `json:"path"`
	Edits     []FileEdit `json:"edits"`
}

// AgentEditFileResult contains the result of editing a file.
type AgentEditFileResult struct {
	OK bool `json:"ok"`
}

// FileEdits groups a file path with its associated edits.
type FileEdits struct {
	Path  string     `json:"path"`
	Edits []FileEdit `json:"edits"`
}

// AgentEditFilesArgs contains the parameters for editing multiple files in an agent.
type AgentEditFilesArgs struct {
	AgentID   string      `json:"agent_id"`
	SidecarID string      `json:"sidecar_id"`
	Files     []FileEdits `json:"files"`
}

// AgentEditFilesResult contains the result of editing multiple files.
type AgentEditFilesResult struct {
	OK bool `json:"ok"`
}

// AgentPortForwardArgs contains the parameters for port forwarding to an agent.
type AgentPortForwardArgs struct {
	AgentID   string `json:"agent_id"`
	SidecarID string `json:"sidecar_id"`
	Port      int    `json:"port"`
}

// AgentPortForwardResult contains the URL for an active port forward.
type AgentPortForwardResult struct {
	URL string `json:"url"`
}

// AgentListAppsArgs contains the parameters for listing apps on an agent.
type AgentListAppsArgs struct {
	AgentID   string `json:"agent_id"`
	SidecarID string `json:"sidecar_id"`
}

// AppInfo describes an application running on an agent sidecar.
type AppInfo struct {
	Slug        string `json:"slug"`
	DisplayName string `json:"display_name"`
	URL         string `json:"url"`
	Health      string `json:"health"`
}

// AgentListAppsResult contains the list of apps on an agent.
type AgentListAppsResult struct {
	Apps []AppInfo `json:"apps"`
}

// ---- Log tools ----

// GetAgentSidecarLogsArgs contains the parameters for getting agent sidecar logs.
type GetAgentSidecarLogsArgs struct {
	SidecarID string `json:"sidecar_id"`
}

// LogEntry represents a single log entry.
type LogEntry struct {
	Output    string `json:"output"`
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
}

// GetAgentSidecarLogsResult contains sidecar log entries.
type GetAgentSidecarLogsResult struct {
	Logs []LogEntry `json:"logs"`
}

// GetAgentBuildLogsArgs contains the parameters for getting agent build logs.
type GetAgentBuildLogsArgs struct {
	BuildID string `json:"build_id"`
}

// BuildLogEntry represents a single build log entry.
type BuildLogEntry struct {
	Output    string `json:"output"`
	Level     string `json:"level"`
	Stage     string `json:"stage"`
	CreatedAt string `json:"created_at"`
}

// GetAgentBuildLogsResult contains agent build log entries.
type GetAgentBuildLogsResult struct {
	Logs []BuildLogEntry `json:"logs"`
}

// GetTemplateVersionLogsArgs contains the parameters for getting template version logs.
type GetTemplateVersionLogsArgs struct {
	TemplateVersionID string `json:"template_version_id"`
}

// GetTemplateVersionLogsResult contains template version log entries.
type GetTemplateVersionLogsResult struct {
	Logs []BuildLogEntry `json:"logs"`
}

// ---- Search tools ----

// SearchArgs contains the search query.
type SearchArgs struct {
	Query string `json:"query"`
}

// SearchResultItem represents a single search result.
type SearchResultItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // "template" or "agent"
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
}

// SearchResult contains search results.
type SearchResult struct {
	Items []SearchResultItem `json:"items"`
}

// FetchArgs contains the parameters for fetching details.
type FetchArgs struct {
	ID   string `json:"id"`
	Type string `json:"type"` // "template" or "agent"
}

// FetchResult contains details about a template or agent.
type FetchResult struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Status      string `json:"status,omitempty"`
	// RawJSON holds the full API response for advanced consumers.
	RawJSON string `json:"raw_json,omitempty"`
}
