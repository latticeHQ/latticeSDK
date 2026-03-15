package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/latticehq/latticesdk/client"
)

// Service provides task management operations against the Lattice Runtime API.
type Service struct {
	client *client.Client
}

// New creates a new tasks Service backed by the given client.
func New(c *client.Client) *Service {
	return &Service{client: c}
}

// tasksListResponse is the API response for listing tasks.
type tasksListResponse struct {
	Tasks []Task `json:"tasks"`
	Count int    `json:"count"`
}

// CreateTask creates a new task for the specified user.
func (s *Service) CreateTask(ctx context.Context, user string, req CreateTaskRequest) (Task, error) {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/tasks/%s", user), req)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Task{}, client.ReadBodyAsError(resp)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return Task{}, fmt.Errorf("decode response: %w", err)
	}
	return task, nil
}

// ListTasks returns tasks matching the given filter.
func (s *Service) ListTasks(ctx context.Context, filter *TasksFilter) ([]Task, error) {
	if filter == nil {
		filter = &TasksFilter{}
	}

	resp, err := s.client.Request(ctx, http.MethodGet, "/api/v2/tasks", nil, filterToRequestOption(filter))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, client.ReadBodyAsError(resp)
	}

	var result tasksListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result.Tasks, nil
}

// GetTaskByID returns a single task by its ID.
func (s *Service) GetTaskByID(ctx context.Context, id uuid.UUID) (Task, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/tasks/me/%s", id), nil)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Task{}, client.ReadBodyAsError(resp)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return Task{}, fmt.Errorf("decode response: %w", err)
	}
	return task, nil
}

// GetTaskByOwnerAndName returns a single task by its owner and name.
// If owner is empty it defaults to "me".
func (s *Service) GetTaskByOwnerAndName(ctx context.Context, owner, ident string) (Task, error) {
	if owner == "" {
		owner = "me"
	}
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/tasks/%s/%s", owner, ident), nil)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Task{}, client.ReadBodyAsError(resp)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return Task{}, fmt.Errorf("decode response: %w", err)
	}
	return task, nil
}

// GetTaskByIdentifier fetches a task by an identifier that may be a UUID,
// a bare name (owned by the current user), or an "owner/name" pair.
func (s *Service) GetTaskByIdentifier(ctx context.Context, identifier string) (Task, error) {
	identifier = strings.TrimSpace(identifier)

	// Try parsing as UUID first.
	if taskID, err := uuid.Parse(identifier); err == nil {
		return s.GetTaskByID(ctx, taskID)
	}

	owner, taskName, err := splitTaskIdentifier(identifier)
	if err != nil {
		return Task{}, err
	}
	return s.GetTaskByOwnerAndName(ctx, owner, taskName)
}

// DeleteTask deletes a task by its ID.
func (s *Service) DeleteTask(ctx context.Context, user string, id uuid.UUID) error {
	resp, err := s.client.Request(ctx, http.MethodDelete, fmt.Sprintf("/api/v2/tasks/%s/%s", user, id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// SendToTask submits input to the task's sidebar app.
func (s *Service) SendToTask(ctx context.Context, user string, id uuid.UUID, req TaskSendRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPost, fmt.Sprintf("/api/v2/tasks/%s/%s/send", user, id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// UpdateTaskInput updates the task's input.
func (s *Service) UpdateTaskInput(ctx context.Context, user string, id uuid.UUID, req UpdateTaskInputRequest) error {
	resp, err := s.client.Request(ctx, http.MethodPatch, fmt.Sprintf("/api/v2/tasks/%s/%s/input", user, id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return client.ReadBodyAsError(resp)
	}
	return nil
}

// GetTaskLogs retrieves logs from the task app.
func (s *Service) GetTaskLogs(ctx context.Context, user string, id uuid.UUID) (TaskLogsResponse, error) {
	resp, err := s.client.Request(ctx, http.MethodGet, fmt.Sprintf("/api/v2/tasks/%s/%s/logs", user, id), nil)
	if err != nil {
		return TaskLogsResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TaskLogsResponse{}, client.ReadBodyAsError(resp)
	}

	var logs TaskLogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return TaskLogsResponse{}, fmt.Errorf("decode response: %w", err)
	}
	return logs, nil
}

// filterToRequestOption converts a TasksFilter into a client.RequestOption
// that sets query parameters on the outgoing request.
func filterToRequestOption(f *TasksFilter) client.RequestOption {
	return func(r *http.Request) {
		var params []string
		if f.Owner != "" {
			params = append(params, fmt.Sprintf("owner:%q", f.Owner))
		}
		if f.Organization != "" {
			params = append(params, fmt.Sprintf("organization:%q", f.Organization))
		}
		if f.Status != "" {
			params = append(params, fmt.Sprintf("status:%q", string(f.Status)))
		}
		if f.FilterQuery != "" {
			params = append(params, f.FilterQuery)
		}
		if len(params) > 0 {
			q := r.URL.Query()
			q.Set("q", strings.Join(params, " "))
			r.URL.RawQuery = q.Encode()
		}
	}
}

// splitTaskIdentifier splits a task identifier into owner and task name.
// Accepted formats are "name" (defaults owner to "me") or "owner/name".
func splitTaskIdentifier(identifier string) (owner, taskName string, err error) {
	parts := strings.Split(identifier, "/")
	switch len(parts) {
	case 1:
		return "me", parts[0], nil
	case 2:
		return parts[0], parts[1], nil
	default:
		return "", "", fmt.Errorf("invalid task identifier: %q", identifier)
	}
}
