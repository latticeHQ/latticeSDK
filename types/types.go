// Package types defines the domain types for the Lattice SDK.
// These types are pure data — no client dependency, only uuid + stdlib.
package types

// AutomaticUpdates controls agent automatic update behavior.
type AutomaticUpdates string

const (
	AutomaticUpdatesAlways AutomaticUpdates = "always"
	AutomaticUpdatesNever  AutomaticUpdates = "never"
)

// ProvisionerLogLevel controls provisioner log verbosity.
type ProvisionerLogLevel string

const (
	ProvisionerLogLevelDebug ProvisionerLogLevel = "debug"
)

// ProvisionerType identifies the provisioner backend.
type ProvisionerType string

const (
	ProvisionerTypeTerraform ProvisionerType = "terraform"
	ProvisionerTypeEcho      ProvisionerType = "echo"
)

// AgentTransition represents agent lifecycle transitions.
type AgentTransition string

const (
	AgentTransitionStart  AgentTransition = "start"
	AgentTransitionStop   AgentTransition = "stop"
	AgentTransitionDelete AgentTransition = "delete"
)

// BuildReason explains why a build was triggered.
type BuildReason string

const (
	BuildReasonAutostart BuildReason = "autostart"
	BuildReasonAutostop  BuildReason = "autostop"
	BuildReasonInitiator BuildReason = "initiator"
)

// UsageAppName identifies the application connecting to an agent.
type UsageAppName string

const (
	UsageAppNameVSCode          UsageAppName = "vscode"
	UsageAppNameJetBrains       UsageAppName = "jetbrains"
	UsageAppNameReconnectingPTY UsageAppName = "reconnecting-pty"
	UsageAppNameSSH             UsageAppName = "ssh"
)
