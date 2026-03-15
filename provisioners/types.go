// Package provisioners provides provisioner daemon and key management
// for the Lattice Runtime API.
package provisioners

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LogSource identifies the origin of a provisioner log entry.
type LogSource string

// LogLevel represents the severity of a provisioner log entry.
type LogLevel string

const (
	LogSourceProvisionerDaemon LogSource = "provisioner_daemon"
	LogSourceProvisioner       LogSource = "provisioner"

	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// ProvisionerDaemon represents a running provisioner daemon instance.
type ProvisionerDaemon struct {
	ID             uuid.UUID         `json:"id" format:"uuid"`
	OrganizationID uuid.UUID         `json:"organization_id" format:"uuid"`
	KeyID          uuid.UUID         `json:"key_id" format:"uuid"`
	CreatedAt      time.Time         `json:"created_at" format:"date-time"`
	LastSeenAt     *time.Time        `json:"last_seen_at,omitempty" format:"date-time"`
	Name           string            `json:"name"`
	Version        string            `json:"version"`
	APIVersion     string            `json:"api_version"`
	Provisioners   []string          `json:"provisioners"`
	Tags           map[string]string `json:"tags"`
}

// ProvisionerJobStatus represents the at-time state of a provisioner job.
type ProvisionerJobStatus string

// Active returns whether the job is still active.
// It returns true if canceling as well, since the job is not yet
// in an entirely inactive state.
func (p ProvisionerJobStatus) Active() bool {
	return p == ProvisionerJobPending ||
		p == ProvisionerJobRunning ||
		p == ProvisionerJobCanceling
}

const (
	ProvisionerJobPending   ProvisionerJobStatus = "pending"
	ProvisionerJobRunning   ProvisionerJobStatus = "running"
	ProvisionerJobSucceeded ProvisionerJobStatus = "succeeded"
	ProvisionerJobCanceling ProvisionerJobStatus = "canceling"
	ProvisionerJobCanceled  ProvisionerJobStatus = "canceled"
	ProvisionerJobFailed    ProvisionerJobStatus = "failed"
	ProvisionerJobUnknown   ProvisionerJobStatus = "unknown"
)

// JobErrorCode defines the error code returned by a provisioner job runner.
type JobErrorCode string

const (
	// RequiredTemplateVariables indicates required template variables were not supplied.
	RequiredTemplateVariables JobErrorCode = "REQUIRED_TEMPLATE_VARIABLES"
)

// ProvisionerJobLog represents a single log entry from a provisioner job.
type ProvisionerJobLog struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	Source    LogSource `json:"log_source"`
	Level     LogLevel  `json:"log_level"`
	Stage     string    `json:"stage"`
	Output    string    `json:"output"`
}

// ProvisionerKeyTags is a set of key-value tags associated with a provisioner key.
type ProvisionerKeyTags map[string]string

// String returns a sorted, space-separated representation of the tags (key=value pairs).
func (p ProvisionerKeyTags) String() string {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	tags := make([]string, 0, len(keys))
	for _, key := range keys {
		tags = append(tags, fmt.Sprintf("%s=%s", key, p[key]))
	}
	return strings.Join(tags, " ")
}

// ProvisionerKey represents an authentication key for provisioner daemons.
type ProvisionerKey struct {
	ID             uuid.UUID          `json:"id" format:"uuid"`
	CreatedAt      time.Time          `json:"created_at" format:"date-time"`
	OrganizationID uuid.UUID          `json:"organization" format:"uuid"`
	Name           string             `json:"name"`
	Tags           ProvisionerKeyTags `json:"tags"`
}

// ProvisionerKeyDaemons pairs a provisioner key with its associated daemons.
type ProvisionerKeyDaemons struct {
	Key     ProvisionerKey      `json:"key"`
	Daemons []ProvisionerDaemon `json:"daemons"`
}

// CreateProvisionerKeyRequest is the request body for creating a new provisioner key.
type CreateProvisionerKeyRequest struct {
	Name string            `json:"name"`
	Tags map[string]string `json:"tags"`
}

// CreateProvisionerKeyResponse is the response from creating a provisioner key.
type CreateProvisionerKeyResponse struct {
	Key string `json:"key"`
}
