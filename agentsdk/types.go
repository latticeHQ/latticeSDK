package agentsdk

import (
	"encoding/json"

	"github.com/google/uuid"
)

// SidecarPort constants define well-known ports used by sidecars.
const (
	// SidecarSSHPort is the default SSH port exposed by sidecars.
	SidecarSSHPort = 22

	// SidecarHTTPPort is the default HTTP API port for the sidecar.
	SidecarHTTPPort = 4

	// SidecarMinPort is the minimum assignable port for sidecar services.
	SidecarMinPort = 1024

	// SidecarMaxPort is the maximum assignable port for sidecar services.
	SidecarMaxPort = 65535
)

// ConnectionInfo contains the information needed to establish a direct
// connection to a sidecar via DERP or direct networking.
type ConnectionInfo struct {
	DERPMap                    json.RawMessage `json:"derp_map"`
	DERPForceWebSockets        bool            `json:"derp_force_websockets"`
	DisableDirectConnections   bool            `json:"disable_direct_connections"`
}

// PTYOptions are the parameters for opening a reconnecting PTY session.
type PTYOptions struct {
	// SidecarID is the UUID of the sidecar to connect to.
	SidecarID uuid.UUID
	// Reconnect is a unique ID for reconnecting to an existing PTY session.
	// Use uuid.New() for a new session.
	Reconnect uuid.UUID
	// Width is the terminal width in columns.
	Width uint16
	// Height is the terminal height in rows.
	Height uint16
	// Command is an optional command to execute in the PTY.
	Command string
	// SignedToken is an optional pre-signed authentication token. When set,
	// it is sent as a query parameter instead of the session token header.
	SignedToken string
}

// PTYRequest is sent over the PTY connection to resize or send input.
type PTYRequest struct {
	Data   string `json:"data,omitempty"`
	Height uint16 `json:"height,omitempty"`
	Width  uint16 `json:"width,omitempty"`
}

// LSRelativity controls the base path used when listing directories.
type LSRelativity string

const (
	// LSRelativityRoot resolves paths relative to the filesystem root.
	LSRelativityRoot LSRelativity = "root"
	// LSRelativityHome resolves paths relative to the user's home directory.
	LSRelativityHome LSRelativity = "home"
)

// LSRequest is the request body for listing a directory.
type LSRequest struct {
	// Path is the path segments to the directory to list.
	Path []string `json:"path"`
	// Relativity controls how Path is resolved.
	Relativity LSRelativity `json:"relativity"`
}

// LSResponse is the response from a directory listing.
type LSResponse struct {
	// AbsolutePath is the absolute path as individual segments.
	AbsolutePath []string `json:"absolute_path"`
	// AbsolutePathString is the joined absolute path string.
	AbsolutePathString string `json:"absolute_path_string"`
	// Contents is the list of files and directories in the listed path.
	Contents []LSFile `json:"contents"`
}

// LSFile describes a single file or directory entry.
type LSFile struct {
	// Name is the base name of the file.
	Name string `json:"name"`
	// AbsolutePathString is the full absolute path to the file.
	AbsolutePathString string `json:"absolute_path_string"`
	// IsDir is true if the entry is a directory.
	IsDir bool `json:"is_dir"`
}

// FileEdit describes a single search-and-replace edit within a file.
type FileEdit struct {
	Search  string `json:"search"`
	Replace string `json:"replace"`
}

// FileEdits groups edits to be applied to a single file.
type FileEdits struct {
	Path  string     `json:"path"`
	Edits []FileEdit `json:"edits"`
}

// FileEditRequest is the request body for applying edits to files.
type FileEditRequest struct {
	Files []FileEdits `json:"files"`
}

// ListeningPortsResponse contains the ports a sidecar is listening on.
type ListeningPortsResponse struct {
	Ports []ListeningPort `json:"ports"`
}

// ListeningPort describes a single port being listened on inside a sidecar.
type ListeningPort struct {
	ProcessName string `json:"process_name"`
	Network     string `json:"network"`
	Port        uint16 `json:"port"`
}

// NetcheckReport contains the results of a network connectivity check.
// The raw JSON is preserved because the report structure depends on the
// sidecar version and may contain version-specific fields.
type NetcheckReport struct {
	Raw json.RawMessage
}

// MarshalJSON implements json.Marshaler.
func (n NetcheckReport) MarshalJSON() ([]byte, error) {
	return n.Raw, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *NetcheckReport) UnmarshalJSON(data []byte) error {
	n.Raw = append(n.Raw[:0], data...)
	return nil
}
