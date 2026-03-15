package stack

import (
	"fmt"
	"os"
)

// Config holds the configuration for a Department Stack.
type Config struct {
	// RuntimeURL is the base URL of the Lattice Runtime API.
	RuntimeURL string

	// APIKey is the API key for authenticating with Runtime.
	APIKey string

	// SessionToken is an alternative to APIKey for session-based auth.
	SessionToken string

	// StackName is the human-readable name of this stack.
	StackName string

	// StackVersion is the version string of this stack (e.g., "0.1.0").
	StackVersion string

	// Endpoint is the optional URL where this stack can be reached by other stacks.
	Endpoint string
}

// LoadConfigFromEnv creates a Config from environment variables.
//
//	LATTICE_RUNTIME_URL  — Runtime API base URL (required)
//	LATTICE_API_KEY      — API key (required unless LATTICE_SESSION_TOKEN is set)
//	LATTICE_SESSION_TOKEN — Session token (alternative to API key)
//	LATTICE_STACK_NAME   — Stack name (required)
//	LATTICE_STACK_VERSION — Stack version (optional, defaults to "0.0.0")
//	LATTICE_STACK_ENDPOINT — Stack endpoint URL (optional)
func LoadConfigFromEnv() (Config, error) {
	c := Config{
		RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
		APIKey:       os.Getenv("LATTICE_API_KEY"),
		SessionToken: os.Getenv("LATTICE_SESSION_TOKEN"),
		StackName:    os.Getenv("LATTICE_STACK_NAME"),
		StackVersion: os.Getenv("LATTICE_STACK_VERSION"),
		Endpoint:     os.Getenv("LATTICE_STACK_ENDPOINT"),
	}

	if c.RuntimeURL == "" {
		return Config{}, fmt.Errorf("LATTICE_RUNTIME_URL is required")
	}
	if c.APIKey == "" && c.SessionToken == "" {
		return Config{}, fmt.Errorf("LATTICE_API_KEY or LATTICE_SESSION_TOKEN is required")
	}
	if c.StackName == "" {
		return Config{}, fmt.Errorf("LATTICE_STACK_NAME is required")
	}
	if c.StackVersion == "" {
		c.StackVersion = "0.0.0"
	}

	return c, nil
}
