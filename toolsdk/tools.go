package toolsdk

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// allTools is the registry of all generic tools, populated in init.
var allTools []GenericTool

func init() {
	allTools = []GenericTool{
		ReportTask.Generic(),
		GetAgent.Generic(),
		CreateAgent.Generic(),
		ListAgents.Generic(),
		CreateAgentBuild.Generic(),
		ListTemplates.Generic(),
		ListTemplateVersionParameters.Generic(),
		CreateTemplateVersion.Generic(),
		UpdateTemplateActiveVersion.Generic(),
		CreateTemplate.Generic(),
		DeleteTemplate.Generic(),
		UploadTarFile.Generic(),
		GetAuthenticatedUser.Generic(),
		CreateTask.Generic(),
		DeleteTask.Generic(),
		ListTasks.Generic(),
		GetTaskStatus.Generic(),
		SendTaskInput.Generic(),
		GetTaskLogs.Generic(),
		AgentBash.Generic(),
		AgentLS.Generic(),
		AgentReadFile.Generic(),
		AgentWriteFile.Generic(),
		AgentEditFile.Generic(),
		AgentEditFiles.Generic(),
		AgentPortForward.Generic(),
		AgentListApps.Generic(),
		GetAgentSidecarLogs.Generic(),
		GetAgentBuildLogs.Generic(),
		GetTemplateVersionLogs.Generic(),
		Search.Generic(),
		Fetch.Generic(),
	}
}

// ---------------------------------------------------------------------------
// ReportTask
// ---------------------------------------------------------------------------

// ReportTask reports a status update for the current task.
var ReportTask = Tool[ReportTaskArgs, ReportTaskResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameReportTask,
		Description: "Report the status of the current task. Use this to update the task state (working, idle, complete, failure) and provide a human-readable summary of progress.",
		Parameters: map[string]Parameter{
			"link":    {Type: "string", Description: "A relevant URL associated with this status update.", Required: false},
			"state":   {Type: "string", Description: "The current task state: working, idle, complete, or failure.", Required: true},
			"summary": {Type: "string", Description: "A brief human-readable summary of the current task progress.", Required: true},
		},
	},
	UserClientOptional: true,
	Handler: func(ctx context.Context, deps Deps, args ReportTaskArgs) (ReportTaskResult, error) {
		if deps.ReportTask == nil {
			return ReportTaskResult{OK: false}, fmt.Errorf("no task reporter configured")
		}
		if err := deps.ReportTask(args); err != nil {
			return ReportTaskResult{OK: false}, err
		}
		return ReportTaskResult{OK: true}, nil
	},
}

// ---------------------------------------------------------------------------
// Agent management tools
// ---------------------------------------------------------------------------

// GetAgent retrieves an agent by its ID.
var GetAgent = Tool[GetAgentArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetAgent,
		Description: "Get details about a Lattice agent by its ID.",
		Parameters: map[string]Parameter{
			"agent_id": {Type: "string", Description: "The UUID of the agent to retrieve.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args GetAgentArgs) (json.RawMessage, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return nil, fmt.Errorf("parse agent_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/agents/%s", agentID), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// CreateAgent creates a new agent from a template.
var CreateAgent = Tool[CreateAgentArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameCreateAgent,
		Description: "Create a new Lattice agent from a template version. The agent is created under the specified user.",
		Parameters: map[string]Parameter{
			"user":                {Type: "string", Description: "The username to create the agent under.", Required: true},
			"template_version_id": {Type: "string", Description: "The UUID of the template version to use.", Required: true},
			"name":                {Type: "string", Description: "The name for the new agent.", Required: true},
			"parameters":          {Type: "object", Description: "Key-value map of build parameters.", Required: false},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args CreateAgentArgs) (json.RawMessage, error) {
		tvID, err := uuid.Parse(args.TemplateVersionID)
		if err != nil {
			return nil, fmt.Errorf("parse template_version_id: %w", err)
		}
		type buildParam struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		}
		type createReq struct {
			TemplateVersionID   uuid.UUID    `json:"template_version_id"`
			Transition          string       `json:"transition"`
			RichParameterValues []buildParam `json:"rich_parameter_values,omitempty"`
		}
		req := createReq{
			TemplateVersionID: tvID,
			Transition:        "start",
		}
		for k, v := range args.Parameters {
			req.RichParameterValues = append(req.RichParameterValues, buildParam{Name: k, Value: v})
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/users/%s/agents", args.User), req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// ListAgents lists agents the authenticated user has access to.
var ListAgents = Tool[ListAgentsArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameListAgents,
		Description: "List Lattice agents. Optionally filter by owner, template, or status.",
		Parameters: map[string]Parameter{
			"owner":    {Type: "string", Description: "Filter by owner username (or 'me').", Required: false},
			"template": {Type: "string", Description: "Filter by template name.", Required: false},
			"status":   {Type: "string", Description: "Filter by agent status (running, stopped, etc.).", Required: false},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args ListAgentsArgs) (json.RawMessage, error) {
		var opts []client.RequestOption
		var params []string
		if args.Owner != "" {
			params = append(params, fmt.Sprintf("owner:%q", args.Owner))
		}
		if args.Template != "" {
			params = append(params, fmt.Sprintf("template:%q", args.Template))
		}
		if args.Status != "" {
			params = append(params, fmt.Sprintf("status:%q", args.Status))
		}
		if len(params) > 0 {
			opts = append(opts, func(r *http.Request) {
				q := r.URL.Query()
				q.Set("q", strings.Join(params, " "))
				r.URL.RawQuery = q.Encode()
			})
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, "/api/v2/agents", nil, opts...)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// CreateAgentBuild creates a new build for an agent (start, stop, or delete).
var CreateAgentBuild = Tool[CreateAgentBuildArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameCreateAgentBuild,
		Description: "Create a new build for a Lattice agent, triggering a start, stop, or delete transition.",
		Parameters: map[string]Parameter{
			"agent_id":            {Type: "string", Description: "The UUID of the agent.", Required: true},
			"template_version_id": {Type: "string", Description: "The template version ID to use. If empty, uses the current version.", Required: false},
			"transition":          {Type: "string", Description: "The transition to perform: start, stop, or delete.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args CreateAgentBuildArgs) (json.RawMessage, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return nil, fmt.Errorf("parse agent_id: %w", err)
		}
		type buildReq struct {
			TemplateVersionID uuid.UUID `json:"template_version_id,omitempty"`
			Transition        string    `json:"transition"`
		}
		req := buildReq{Transition: args.Transition}
		if args.TemplateVersionID != "" {
			tvID, err := uuid.Parse(args.TemplateVersionID)
			if err != nil {
				return nil, fmt.Errorf("parse template_version_id: %w", err)
			}
			req.TemplateVersionID = tvID
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/agents/%s/builds", agentID), req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// ---------------------------------------------------------------------------
// Template tools
// ---------------------------------------------------------------------------

// ListTemplates lists available templates.
var ListTemplates = Tool[ListTemplatesArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameListTemplates,
		Description: "List available Lattice templates. Optionally filter by organization.",
		Parameters: map[string]Parameter{
			"organization_id": {Type: "string", Description: "Filter by organization UUID.", Required: false},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args ListTemplatesArgs) (json.RawMessage, error) {
		path := "/api/v2/templates"
		if args.OrganizationID != "" {
			orgID, err := uuid.Parse(args.OrganizationID)
			if err != nil {
				return nil, fmt.Errorf("parse organization_id: %w", err)
			}
			path = fmt.Sprintf("/api/v2/organizations/%s/templates", orgID)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// ListTemplateVersionParameters returns the rich parameters for a template version.
var ListTemplateVersionParameters = Tool[TemplateVersionParametersArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameTemplateVersionParameters,
		Description: "Get the configurable parameters for a template version. Useful for understanding what inputs are needed to create an agent.",
		Parameters: map[string]Parameter{
			"template_version_id": {Type: "string", Description: "The UUID of the template version.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args TemplateVersionParametersArgs) (json.RawMessage, error) {
		tvID, err := uuid.Parse(args.TemplateVersionID)
		if err != nil {
			return nil, fmt.Errorf("parse template_version_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/templateversions/%s/rich-parameters", tvID), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// CreateTemplateVersion creates a new template version.
var CreateTemplateVersion = Tool[CreateTemplateVersionArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameCreateTemplateVersion,
		Description: "Create a new version of a Lattice template. Requires a previously uploaded file and a provisioner type.",
		Parameters: map[string]Parameter{
			"organization_id": {Type: "string", Description: "The UUID of the organization.", Required: true},
			"template_id":     {Type: "string", Description: "The UUID of the template to attach this version to.", Required: false},
			"name":            {Type: "string", Description: "Name for the new version.", Required: false},
			"message":         {Type: "string", Description: "Description for the new version.", Required: false},
			"file_id":         {Type: "string", Description: "The UUID of the uploaded file containing template source.", Required: true},
			"provisioner":     {Type: "string", Description: "The provisioner type: terraform or echo.", Required: true},
			"tags":            {Type: "object", Description: "Key-value provisioner tags.", Required: false},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args CreateTemplateVersionArgs) (json.RawMessage, error) {
		orgID, err := uuid.Parse(args.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("parse organization_id: %w", err)
		}
		fileID, err := uuid.Parse(args.FileID)
		if err != nil {
			return nil, fmt.Errorf("parse file_id: %w", err)
		}
		type createReq struct {
			Name            string            `json:"name,omitempty"`
			Message         string            `json:"message,omitempty"`
			TemplateID      uuid.UUID         `json:"template_id,omitempty"`
			StorageMethod   string            `json:"storage_method"`
			FileID          uuid.UUID         `json:"file_id"`
			Provisioner     string            `json:"provisioner"`
			ProvisionerTags map[string]string `json:"tags,omitempty"`
		}
		req := createReq{
			Name:          args.Name,
			Message:       args.Message,
			StorageMethod: "file",
			FileID:        fileID,
			Provisioner:   args.Provisioner,
		}
		if args.TemplateID != "" {
			tmplID, err := uuid.Parse(args.TemplateID)
			if err != nil {
				return nil, fmt.Errorf("parse template_id: %w", err)
			}
			req.TemplateID = tmplID
		}
		if len(args.Tags) > 0 {
			req.ProvisionerTags = args.Tags
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/organizations/%s/templateversions", orgID), req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// UpdateTemplateActiveVersion promotes a template version to be the active version.
var UpdateTemplateActiveVersion = Tool[UpdateTemplateActiveVersionArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameUpdateTemplateActiveVersion,
		Description: "Update the active version of a Lattice template. This promotes a specific version to be used by default.",
		Parameters: map[string]Parameter{
			"template_id": {Type: "string", Description: "The UUID of the template.", Required: true},
			"version_id":  {Type: "string", Description: "The UUID of the version to promote.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args UpdateTemplateActiveVersionArgs) (json.RawMessage, error) {
		tmplID, err := uuid.Parse(args.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("parse template_id: %w", err)
		}
		versionID, err := uuid.Parse(args.VersionID)
		if err != nil {
			return nil, fmt.Errorf("parse version_id: %w", err)
		}
		type updateReq struct {
			ID uuid.UUID `json:"id"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/templates/%s/versions", tmplID), updateReq{ID: versionID})
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		result := map[string]bool{"ok": true}
		data, _ := json.Marshal(result)
		return json.RawMessage(data), nil
	},
}

// CreateTemplate creates a new template.
var CreateTemplate = Tool[CreateTemplateArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameCreateTemplate,
		Description: "Create a new Lattice template in an organization.",
		Parameters: map[string]Parameter{
			"organization_id": {Type: "string", Description: "The UUID of the organization.", Required: true},
			"name":            {Type: "string", Description: "The name for the new template.", Required: true},
			"display_name":    {Type: "string", Description: "Display name for the template.", Required: false},
			"description":     {Type: "string", Description: "Description of the template.", Required: false},
			"version_id":      {Type: "string", Description: "The UUID of the initial template version.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args CreateTemplateArgs) (json.RawMessage, error) {
		orgID, err := uuid.Parse(args.OrganizationID)
		if err != nil {
			return nil, fmt.Errorf("parse organization_id: %w", err)
		}
		versionID, err := uuid.Parse(args.VersionID)
		if err != nil {
			return nil, fmt.Errorf("parse version_id: %w", err)
		}
		type createReq struct {
			Name              string    `json:"name"`
			DisplayName       string    `json:"display_name,omitempty"`
			Description       string    `json:"description,omitempty"`
			TemplateVersionID uuid.UUID `json:"template_version_id"`
		}
		req := createReq{
			Name:              args.Name,
			DisplayName:       args.DisplayName,
			Description:       args.Description,
			TemplateVersionID: versionID,
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/organizations/%s/templates", orgID), req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// DeleteTemplate deletes a template by ID.
var DeleteTemplate = Tool[DeleteTemplateArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameDeleteTemplate,
		Description: "Delete a Lattice template by its ID. This is a destructive operation.",
		Parameters: map[string]Parameter{
			"template_id": {Type: "string", Description: "The UUID of the template to delete.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args DeleteTemplateArgs) (json.RawMessage, error) {
		tmplID, err := uuid.Parse(args.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("parse template_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/templates/%s", tmplID), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		result := map[string]bool{"ok": true}
		data, _ := json.Marshal(result)
		return json.RawMessage(data), nil
	},
}

// UploadTarFile uploads a base64-encoded tar file for use with template versions.
var UploadTarFile = Tool[UploadTarFileArgs, UploadTarFileResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameUploadTarFile,
		Description: "Upload a tar file (base64-encoded) to Lattice. Returns a file ID that can be used when creating template versions.",
		Parameters: map[string]Parameter{
			"file_base64": {Type: "string", Description: "The base64-encoded content of the tar file.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args UploadTarFileArgs) (UploadTarFileResult, error) {
		decoded, err := base64.StdEncoding.DecodeString(args.FileBase64)
		if err != nil {
			return UploadTarFileResult{}, fmt.Errorf("decode base64: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, "/api/v2/files", bytes.NewReader(decoded), func(r *http.Request) {
			r.Header.Set("Content-Type", "application/x-tar")
		})
		if err != nil {
			return UploadTarFileResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
			return UploadTarFileResult{}, client.ReadBodyAsError(resp)
		}
		var result struct {
			ID uuid.UUID `json:"hash"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return UploadTarFileResult{}, fmt.Errorf("decode response: %w", err)
		}
		return UploadTarFileResult{ID: result.ID}, nil
	},
}

// ---------------------------------------------------------------------------
// User tools
// ---------------------------------------------------------------------------

// GetAuthenticatedUser returns the currently authenticated user.
var GetAuthenticatedUser = Tool[GetAuthenticatedUserArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetAuthenticatedUser,
		Description: "Get the currently authenticated Lattice user.",
	},
	Handler: func(ctx context.Context, deps Deps, _ GetAuthenticatedUserArgs) (json.RawMessage, error) {
		resp, err := deps.Client.Request(ctx, http.MethodGet, "/api/v2/users/me", nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// ---------------------------------------------------------------------------
// Task tools
// ---------------------------------------------------------------------------

// CreateTask creates a new task.
var CreateTask = Tool[CreateTaskArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameCreateTask,
		Description: "Create a new Lattice task for the specified user.",
		Parameters: map[string]Parameter{
			"user":                {Type: "string", Description: "The username to create the task under.", Required: true},
			"template_version_id": {Type: "string", Description: "The UUID of the template version to use.", Required: true},
			"input":               {Type: "string", Description: "The initial prompt or input for the task.", Required: true},
			"name":                {Type: "string", Description: "Optional name for the task.", Required: false},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args CreateTaskArgs) (json.RawMessage, error) {
		tvID, err := uuid.Parse(args.TemplateVersionID)
		if err != nil {
			return nil, fmt.Errorf("parse template_version_id: %w", err)
		}
		type createReq struct {
			TemplateVersionID uuid.UUID `json:"template_version_id"`
			Input             string    `json:"input"`
			Name              string    `json:"name,omitempty"`
		}
		req := createReq{
			TemplateVersionID: tvID,
			Input:             args.Input,
			Name:              args.Name,
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/tasks/%s", args.User), req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusCreated {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// DeleteTask deletes a task by its ID.
var DeleteTask = Tool[DeleteTaskArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameDeleteTask,
		Description: "Delete a Lattice task. This is a destructive operation.",
		Parameters: map[string]Parameter{
			"user":    {Type: "string", Description: "The username that owns the task.", Required: true},
			"task_id": {Type: "string", Description: "The UUID of the task to delete.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args DeleteTaskArgs) (json.RawMessage, error) {
		taskID, err := uuid.Parse(args.TaskID)
		if err != nil {
			return nil, fmt.Errorf("parse task_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/tasks/%s/%s", args.User, taskID), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusAccepted {
			return nil, client.ReadBodyAsError(resp)
		}
		result := map[string]bool{"ok": true}
		data, _ := json.Marshal(result)
		return json.RawMessage(data), nil
	},
}

// ListTasks lists tasks matching the given filters.
var ListTasks = Tool[ListTasksArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameListTasks,
		Description: "List Lattice tasks. Optionally filter by owner or status.",
		Parameters: map[string]Parameter{
			"owner":  {Type: "string", Description: "Filter by owner username (or 'me').", Required: false},
			"status": {Type: "string", Description: "Filter by task status (pending, active, paused, error).", Required: false},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args ListTasksArgs) (json.RawMessage, error) {
		var opts []client.RequestOption
		var params []string
		if args.Owner != "" {
			params = append(params, fmt.Sprintf("owner:%q", args.Owner))
		}
		if args.Status != "" {
			params = append(params, fmt.Sprintf("status:%q", args.Status))
		}
		if len(params) > 0 {
			opts = append(opts, func(r *http.Request) {
				q := r.URL.Query()
				q.Set("q", strings.Join(params, " "))
				r.URL.RawQuery = q.Encode()
			})
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, "/api/v2/tasks", nil, opts...)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// GetTaskStatus gets the current status of a task.
var GetTaskStatus = Tool[GetTaskStatusArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetTaskStatus,
		Description: "Get the current status of a Lattice task by its identifier (UUID, name, or owner/name).",
		Parameters: map[string]Parameter{
			"task_identifier": {Type: "string", Description: "The task identifier: a UUID, bare name, or owner/name pair.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args GetTaskStatusArgs) (json.RawMessage, error) {
		identifier := strings.TrimSpace(args.TaskIdentifier)
		// Try UUID first.
		var path string
		if taskID, err := uuid.Parse(identifier); err == nil {
			path = fmt.Sprintf("/api/v2/tasks/me/%s", taskID)
		} else {
			parts := strings.SplitN(identifier, "/", 2)
			if len(parts) == 2 {
				path = fmt.Sprintf("/api/v2/tasks/%s/%s", parts[0], parts[1])
			} else {
				path = fmt.Sprintf("/api/v2/tasks/me/%s", identifier)
			}
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// SendTaskInput sends input to a running task.
var SendTaskInput = Tool[SendTaskInputArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameSendTaskInput,
		Description: "Send input to a running Lattice task's sidebar app.",
		Parameters: map[string]Parameter{
			"user":    {Type: "string", Description: "The username that owns the task.", Required: true},
			"task_id": {Type: "string", Description: "The UUID of the task.", Required: true},
			"input":   {Type: "string", Description: "The input text to send to the task.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args SendTaskInputArgs) (json.RawMessage, error) {
		taskID, err := uuid.Parse(args.TaskID)
		if err != nil {
			return nil, fmt.Errorf("parse task_id: %w", err)
		}
		type sendReq struct {
			Input string `json:"input"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/tasks/%s/%s/send", args.User, taskID), sendReq{Input: args.Input})
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusNoContent {
			return nil, client.ReadBodyAsError(resp)
		}
		result := map[string]bool{"ok": true}
		data, _ := json.Marshal(result)
		return json.RawMessage(data), nil
	},
}

// GetTaskLogs retrieves logs from a task.
var GetTaskLogs = Tool[GetTaskLogsArgs, json.RawMessage]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetTaskLogs,
		Description: "Get logs from a Lattice task.",
		Parameters: map[string]Parameter{
			"user":    {Type: "string", Description: "The username that owns the task.", Required: true},
			"task_id": {Type: "string", Description: "The UUID of the task.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args GetTaskLogsArgs) (json.RawMessage, error) {
		taskID, err := uuid.Parse(args.TaskID)
		if err != nil {
			return nil, fmt.Errorf("parse task_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/tasks/%s/%s/logs", args.User, taskID), nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}
		return json.RawMessage(data), nil
	},
}

// ---------------------------------------------------------------------------
// Agent operation tools
// ---------------------------------------------------------------------------

// AgentBash executes a bash command in an agent's sidecar.
var AgentBash = Tool[AgentBashArgs, AgentBashResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentBash,
		Description: "Execute a bash command inside a Lattice agent's sidecar environment.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar to execute in.", Required: true},
			"command":    {Type: "string", Description: "The bash command to execute.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentBashArgs) (AgentBashResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentBashResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentBashResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		type execReq struct {
			Command string `json:"command"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/exec", agentID, sidecarID),
			execReq{Command: args.Command})
		if err != nil {
			return AgentBashResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return AgentBashResult{}, client.ReadBodyAsError(resp)
		}
		var result AgentBashResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return AgentBashResult{}, fmt.Errorf("decode response: %w", err)
		}
		return result, nil
	},
}

// AgentLS lists directory contents in an agent's sidecar.
var AgentLS = Tool[AgentLSArgs, AgentLSResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentLS,
		Description: "List the contents of a directory inside a Lattice agent's sidecar.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
			"path":       {Type: "string", Description: "The directory path to list.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentLSArgs) (AgentLSResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentLSResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentLSResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/ls", agentID, sidecarID),
			nil,
			client.WithQueryParam("path", args.Path))
		if err != nil {
			return AgentLSResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return AgentLSResult{}, client.ReadBodyAsError(resp)
		}
		var result AgentLSResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return AgentLSResult{}, fmt.Errorf("decode response: %w", err)
		}
		return result, nil
	},
}

// AgentReadFile reads a file from an agent's sidecar.
var AgentReadFile = Tool[AgentReadFileArgs, AgentReadFileResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentReadFile,
		Description: "Read the contents of a file inside a Lattice agent's sidecar.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
			"path":       {Type: "string", Description: "The file path to read.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentReadFileArgs) (AgentReadFileResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentReadFileResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentReadFileResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/file", agentID, sidecarID),
			nil,
			client.WithQueryParam("path", args.Path))
		if err != nil {
			return AgentReadFileResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return AgentReadFileResult{}, client.ReadBodyAsError(resp)
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return AgentReadFileResult{}, fmt.Errorf("read response: %w", err)
		}
		return AgentReadFileResult{Content: string(content)}, nil
	},
}

// AgentWriteFile writes content to a file in an agent's sidecar.
var AgentWriteFile = Tool[AgentWriteFileArgs, AgentWriteFileResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentWriteFile,
		Description: "Write content to a file inside a Lattice agent's sidecar.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
			"path":       {Type: "string", Description: "The file path to write to.", Required: true},
			"content":    {Type: "string", Description: "The content to write.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentWriteFileArgs) (AgentWriteFileResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentWriteFileResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentWriteFileResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		type writeReq struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPut,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/file", agentID, sidecarID),
			writeReq{Path: args.Path, Content: args.Content})
		if err != nil {
			return AgentWriteFileResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return AgentWriteFileResult{}, client.ReadBodyAsError(resp)
		}
		return AgentWriteFileResult{OK: true}, nil
	},
}

// AgentEditFile applies edits to a file in an agent's sidecar.
var AgentEditFile = Tool[AgentEditFileArgs, AgentEditFileResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentEditFile,
		Description: "Apply text edits to a file inside a Lattice agent's sidecar. Each edit replaces old_text with new_text.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
			"path":       {Type: "string", Description: "The file path to edit.", Required: true},
			"edits":      {Type: "array", Description: "List of edits, each with old_text and new_text.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentEditFileArgs) (AgentEditFileResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentEditFileResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentEditFileResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		type editReq struct {
			Path  string     `json:"path"`
			Edits []FileEdit `json:"edits"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPatch,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/file", agentID, sidecarID),
			editReq{Path: args.Path, Edits: args.Edits})
		if err != nil {
			return AgentEditFileResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return AgentEditFileResult{}, client.ReadBodyAsError(resp)
		}
		return AgentEditFileResult{OK: true}, nil
	},
}

// AgentEditFiles applies edits to multiple files in an agent's sidecar.
var AgentEditFiles = Tool[AgentEditFilesArgs, AgentEditFilesResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentEditFiles,
		Description: "Apply text edits to multiple files inside a Lattice agent's sidecar in a single operation.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
			"files":      {Type: "array", Description: "List of file edit groups, each with a path and edits array.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentEditFilesArgs) (AgentEditFilesResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentEditFilesResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentEditFilesResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		type editFilesReq struct {
			Files []FileEdits `json:"files"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPatch,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/files", agentID, sidecarID),
			editFilesReq{Files: args.Files})
		if err != nil {
			return AgentEditFilesResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return AgentEditFilesResult{}, client.ReadBodyAsError(resp)
		}
		return AgentEditFilesResult{OK: true}, nil
	},
}

// AgentPortForward creates a port forward to an agent's sidecar.
var AgentPortForward = Tool[AgentPortForwardArgs, AgentPortForwardResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentPortForward,
		Description: "Set up port forwarding to a port on a Lattice agent's sidecar. Returns a URL that can be used to access the port.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
			"port":       {Type: "number", Description: "The port number to forward.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentPortForwardArgs) (AgentPortForwardResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentPortForwardResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentPortForwardResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		type portFwdReq struct {
			Port int `json:"port"`
		}
		resp, err := deps.Client.Request(ctx, http.MethodPost,
			fmt.Sprintf("/api/v2/agents/%s/sidecars/%s/port-forward", agentID, sidecarID),
			portFwdReq{Port: args.Port})
		if err != nil {
			return AgentPortForwardResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return AgentPortForwardResult{}, client.ReadBodyAsError(resp)
		}
		var result AgentPortForwardResult
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return AgentPortForwardResult{}, fmt.Errorf("decode response: %w", err)
		}
		return result, nil
	},
}

// AgentListApps lists the applications running on an agent's sidecar.
var AgentListApps = Tool[AgentListAppsArgs, AgentListAppsResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameAgentListApps,
		Description: "List the applications running on a Lattice agent's sidecar.",
		Parameters: map[string]Parameter{
			"agent_id":   {Type: "string", Description: "The UUID of the agent.", Required: true},
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args AgentListAppsArgs) (AgentListAppsResult, error) {
		agentID, err := uuid.Parse(args.AgentID)
		if err != nil {
			return AgentListAppsResult{}, fmt.Errorf("parse agent_id: %w", err)
		}
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return AgentListAppsResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/agentsidecars/%s", sidecarID), nil)
		if err != nil {
			return AgentListAppsResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return AgentListAppsResult{}, client.ReadBodyAsError(resp)
		}
		// Parse the sidecar response to extract apps.
		var sidecar struct {
			Apps []struct {
				Slug        string `json:"slug"`
				DisplayName string `json:"display_name"`
				URL         string `json:"url"`
				Health      string `json:"health"`
			} `json:"apps"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&sidecar); err != nil {
			return AgentListAppsResult{}, fmt.Errorf("decode response: %w", err)
		}
		var apps []AppInfo
		for _, a := range sidecar.Apps {
			apps = append(apps, AppInfo{
				Slug:        a.Slug,
				DisplayName: a.DisplayName,
				URL:         a.URL,
				Health:      a.Health,
			})
		}
		// Suppress unused variable warning for agentID.
		_ = agentID
		return AgentListAppsResult{Apps: apps}, nil
	},
}

// ---------------------------------------------------------------------------
// Log tools
// ---------------------------------------------------------------------------

// GetAgentSidecarLogs retrieves logs from an agent sidecar.
var GetAgentSidecarLogs = Tool[GetAgentSidecarLogsArgs, GetAgentSidecarLogsResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetAgentSidecarLogs,
		Description: "Get logs from a Lattice agent's sidecar process.",
		Parameters: map[string]Parameter{
			"sidecar_id": {Type: "string", Description: "The UUID of the sidecar.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args GetAgentSidecarLogsArgs) (GetAgentSidecarLogsResult, error) {
		sidecarID, err := uuid.Parse(args.SidecarID)
		if err != nil {
			return GetAgentSidecarLogsResult{}, fmt.Errorf("parse sidecar_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/agentsidecars/%s/logs", sidecarID), nil)
		if err != nil {
			return GetAgentSidecarLogsResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return GetAgentSidecarLogsResult{}, client.ReadBodyAsError(resp)
		}
		var logs []struct {
			Output    string `json:"output"`
			Level     string `json:"level"`
			CreatedAt string `json:"created_at"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
			return GetAgentSidecarLogsResult{}, fmt.Errorf("decode response: %w", err)
		}
		var entries []LogEntry
		for _, l := range logs {
			entries = append(entries, LogEntry{
				Output:    l.Output,
				Level:     l.Level,
				CreatedAt: l.CreatedAt,
			})
		}
		return GetAgentSidecarLogsResult{Logs: entries}, nil
	},
}

// GetAgentBuildLogs retrieves logs from an agent build.
var GetAgentBuildLogs = Tool[GetAgentBuildLogsArgs, GetAgentBuildLogsResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetAgentBuildLogs,
		Description: "Get provisioner logs from a Lattice agent build.",
		Parameters: map[string]Parameter{
			"build_id": {Type: "string", Description: "The UUID of the agent build.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args GetAgentBuildLogsArgs) (GetAgentBuildLogsResult, error) {
		buildID, err := uuid.Parse(args.BuildID)
		if err != nil {
			return GetAgentBuildLogsResult{}, fmt.Errorf("parse build_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/agentbuilds/%s/logs", buildID), nil)
		if err != nil {
			return GetAgentBuildLogsResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return GetAgentBuildLogsResult{}, client.ReadBodyAsError(resp)
		}
		var logs []struct {
			Output    string `json:"output"`
			Level     string `json:"log_level"`
			Stage     string `json:"stage"`
			CreatedAt string `json:"created_at"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
			return GetAgentBuildLogsResult{}, fmt.Errorf("decode response: %w", err)
		}
		var entries []BuildLogEntry
		for _, l := range logs {
			entries = append(entries, BuildLogEntry{
				Output:    l.Output,
				Level:     l.Level,
				Stage:     l.Stage,
				CreatedAt: l.CreatedAt,
			})
		}
		return GetAgentBuildLogsResult{Logs: entries}, nil
	},
}

// GetTemplateVersionLogs retrieves provisioner logs from a template version.
var GetTemplateVersionLogs = Tool[GetTemplateVersionLogsArgs, GetTemplateVersionLogsResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameGetTemplateVersionLogs,
		Description: "Get provisioner logs from a Lattice template version build.",
		Parameters: map[string]Parameter{
			"template_version_id": {Type: "string", Description: "The UUID of the template version.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args GetTemplateVersionLogsArgs) (GetTemplateVersionLogsResult, error) {
		tvID, err := uuid.Parse(args.TemplateVersionID)
		if err != nil {
			return GetTemplateVersionLogsResult{}, fmt.Errorf("parse template_version_id: %w", err)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet,
			fmt.Sprintf("/api/v2/templateversions/%s/logs", tvID), nil)
		if err != nil {
			return GetTemplateVersionLogsResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return GetTemplateVersionLogsResult{}, client.ReadBodyAsError(resp)
		}
		var logs []struct {
			Output    string `json:"output"`
			Level     string `json:"log_level"`
			Stage     string `json:"stage"`
			CreatedAt string `json:"created_at"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
			return GetTemplateVersionLogsResult{}, fmt.Errorf("decode response: %w", err)
		}
		var entries []BuildLogEntry
		for _, l := range logs {
			entries = append(entries, BuildLogEntry{
				Output:    l.Output,
				Level:     l.Level,
				Stage:     l.Stage,
				CreatedAt: l.CreatedAt,
			})
		}
		return GetTemplateVersionLogsResult{Logs: entries}, nil
	},
}

// ---------------------------------------------------------------------------
// Search tools
// ---------------------------------------------------------------------------

// Search performs a ChatGPT-style search across templates and agents.
var Search = Tool[SearchArgs, SearchResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameSearch,
		Description: "Search for Lattice templates and agents by name or description. Returns a list of matching items.",
		Parameters: map[string]Parameter{
			"query": {Type: "string", Description: "The search query string.", Required: true},
		},
	},
	UserClientOptional: false,
	Handler: func(ctx context.Context, deps Deps, args SearchArgs) (SearchResult, error) {
		var items []SearchResultItem

		// Search templates.
		templatesResp, err := deps.Client.Request(ctx, http.MethodGet, "/api/v2/templates", nil, func(r *http.Request) {
			q := r.URL.Query()
			q.Set("q", args.Query)
			r.URL.RawQuery = q.Encode()
		})
		if err == nil {
			defer templatesResp.Body.Close()
			if templatesResp.StatusCode == http.StatusOK {
				var templates []struct {
					ID          uuid.UUID `json:"id"`
					Name        string    `json:"name"`
					DisplayName string    `json:"display_name"`
					Description string    `json:"description"`
					Icon        string    `json:"icon"`
				}
				if err := json.NewDecoder(templatesResp.Body).Decode(&templates); err == nil {
					for _, t := range templates {
						items = append(items, SearchResultItem{
							ID:          t.ID.String(),
							Name:        t.Name,
							DisplayName: t.DisplayName,
							Type:        "template",
							Description: t.Description,
							Icon:        t.Icon,
						})
					}
				}
			}
		}

		// Search agents.
		agentsResp, err := deps.Client.Request(ctx, http.MethodGet, "/api/v2/agents", nil, func(r *http.Request) {
			q := r.URL.Query()
			q.Set("q", fmt.Sprintf("name:%q", args.Query))
			r.URL.RawQuery = q.Encode()
		})
		if err == nil {
			defer agentsResp.Body.Close()
			if agentsResp.StatusCode == http.StatusOK {
				var agentsResult struct {
					Agents []struct {
						ID           uuid.UUID `json:"id"`
						Name         string    `json:"name"`
						TemplateName string    `json:"template_name"`
						OwnerName    string    `json:"owner_name"`
					} `json:"agents"`
				}
				if err := json.NewDecoder(agentsResp.Body).Decode(&agentsResult); err == nil {
					for _, a := range agentsResult.Agents {
						items = append(items, SearchResultItem{
							ID:          a.ID.String(),
							Name:        a.Name,
							DisplayName: a.Name,
							Type:        "agent",
							Description: fmt.Sprintf("Template: %s, Owner: %s", a.TemplateName, a.OwnerName),
						})
					}
				}
			}
		}

		return SearchResult{Items: items}, nil
	},
}

// Fetch retrieves detailed information about a template or agent.
var Fetch = Tool[FetchArgs, FetchResult]{
	ToolDef: ToolDefinition{
		Name:        ToolNameFetch,
		Description: "Fetch detailed information about a specific Lattice template or agent by ID.",
		Parameters: map[string]Parameter{
			"id":   {Type: "string", Description: "The UUID of the template or agent.", Required: true},
			"type": {Type: "string", Description: "The type of resource: 'template' or 'agent'.", Required: true},
		},
	},
	Handler: func(ctx context.Context, deps Deps, args FetchArgs) (FetchResult, error) {
		id, err := uuid.Parse(args.ID)
		if err != nil {
			return FetchResult{}, fmt.Errorf("parse id: %w", err)
		}
		var path string
		switch args.Type {
		case "template":
			path = fmt.Sprintf("/api/v2/templates/%s", id)
		case "agent":
			path = fmt.Sprintf("/api/v2/agents/%s", id)
		default:
			return FetchResult{}, fmt.Errorf("unsupported type: %s (must be 'template' or 'agent')", args.Type)
		}
		resp, err := deps.Client.Request(ctx, http.MethodGet, path, nil)
		if err != nil {
			return FetchResult{}, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return FetchResult{}, client.ReadBodyAsError(resp)
		}
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return FetchResult{}, fmt.Errorf("read response: %w", err)
		}
		// Parse common fields.
		var common struct {
			ID          uuid.UUID `json:"id"`
			Name        string    `json:"name"`
			DisplayName string    `json:"display_name"`
			Description string    `json:"description"`
			Icon        string    `json:"icon"`
			Status      string    `json:"status"`
		}
		_ = json.Unmarshal(data, &common)
		displayName := common.DisplayName
		if displayName == "" {
			displayName = common.Name
		}
		return FetchResult{
			ID:          common.ID.String(),
			Name:        common.Name,
			DisplayName: displayName,
			Type:        args.Type,
			Description: common.Description,
			Icon:        common.Icon,
			Status:      common.Status,
			RawJSON:     string(data),
		}, nil
	},
}
