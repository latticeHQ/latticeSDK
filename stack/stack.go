// Package stack provides the high-level entry point for building Department Stacks.
// It wires together all SDK services and manages the stack lifecycle.
package stack

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/latticehq/latticesdk/agents"
	"github.com/latticehq/latticesdk/agentsdk"
	"github.com/latticehq/latticesdk/audit"
	"github.com/latticehq/latticesdk/authz"
	"github.com/latticehq/latticesdk/budget"
	"github.com/latticehq/latticesdk/client"
	"github.com/latticehq/latticesdk/coordination"
	"github.com/latticehq/latticesdk/deployments"
	"github.com/latticehq/latticesdk/evals"
	"github.com/latticehq/latticesdk/externalauth"
	"github.com/latticehq/latticesdk/files"
	"github.com/latticehq/latticesdk/gitsshkeys"
	"github.com/latticehq/latticesdk/groups"
	"github.com/latticehq/latticesdk/healthsdk"
	"github.com/latticehq/latticesdk/identity"
	"github.com/latticehq/latticesdk/insights"
	"github.com/latticehq/latticesdk/licenses"
	"github.com/latticehq/latticesdk/lifecycle"
	"github.com/latticehq/latticesdk/notifications"
	"github.com/latticehq/latticesdk/oauth2"
	"github.com/latticehq/latticesdk/provisioners"
	"github.com/latticehq/latticesdk/replicas"
	"github.com/latticehq/latticesdk/sessions"
	"github.com/latticehq/latticesdk/sidecarsdk"
	"github.com/latticehq/latticesdk/tasks"
	"github.com/latticehq/latticesdk/templates"
	"github.com/latticehq/latticesdk/users"
)

// Stack is the main entry point for a Department Stack. It provides access to
// all Lattice Runtime services and manages registration, heartbeats, and
// graceful shutdown.
type Stack struct {
	// --- Core coordination services ---

	// Identity provides agent, user, organization, and token management.
	Identity *identity.Service

	// Authz provides authorization checks and RBAC operations.
	Authz *authz.Service

	// Audit provides audit trail query and event emission.
	Audit *audit.Service

	// Budget provides quota and cost management.
	Budget *budget.Service

	// Coordination provides cross-stack messaging and shared state.
	Coordination *coordination.Service

	// Lifecycle provides stack registration and health reporting.
	Lifecycle *lifecycle.Service

	// --- Full platform services ---

	// Agents provides complete agent management — CRUD, builds, proxies, sidecars.
	Agents *agents.Service

	// Sessions provides session management — CRUD, builds, sidecars, real-time.
	Sessions *sessions.Service

	// Templates provides template management — CRUD, versions, ACL, parameters.
	Templates *templates.Service

	// Users provides user management — CRUD, auth, roles, org membership.
	Users *users.Service

	// Deployments provides deployment config, stats, appearance, and entitlements.
	Deployments *deployments.Service

	// Tasks provides AI task management — CRUD, send, logs.
	Tasks *tasks.Service

	// Evals provides evaluation management — runs, comparisons, passes.
	Evals *evals.Service

	// Groups provides group management and IDP sync settings.
	Groups *groups.Service

	// OAuth2 provides OAuth2 provider app management.
	OAuth2 *oauth2.Service

	// Notifications provides notification settings, templates, and preferences.
	Notifications *notifications.Service

	// Files provides file upload and download.
	Files *files.Service

	// Insights provides usage analytics — latency, activity, template insights.
	Insights *insights.Service

	// Provisioners provides provisioner daemon and key management.
	Provisioners *provisioners.Service

	// ExternalAuth provides external authentication provider management.
	ExternalAuth *externalauth.Service

	// Licenses provides license management.
	Licenses *licenses.Service

	// Replicas provides replica information.
	Replicas *replicas.Service

	// GitSSHKeys provides Git SSH key management.
	GitSSHKeys *gitsshkeys.Service

	// --- Agent-side services ---

	// AgentSDK provides agent-to-sidecar connectivity — connections, PTY, file ops, debug.
	AgentSDK *agentsdk.Service

	// SidecarSDK provides sidecar-to-Runtime communication — auth, lifecycle, logs, stats.
	SidecarSDK *sidecarsdk.Service

	// Health provides deployment health checks and diagnostics.
	Health *healthsdk.Service

	// Client is the underlying HTTP client. Use this for custom API calls.
	Client *client.Client

	config       Config
	registration *lifecycle.StackRegistration
}

// New creates a new Stack with all services pre-wired.
func New(cfg Config) (*Stack, error) {
	var opts []client.Option
	if cfg.APIKey != "" {
		opts = append(opts, client.WithAPIKey(cfg.APIKey))
	}
	if cfg.SessionToken != "" {
		opts = append(opts, client.WithSessionToken(cfg.SessionToken))
	}

	c, err := client.New(cfg.RuntimeURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	return &Stack{
		// Core coordination
		Identity:     identity.New(c),
		Authz:        authz.New(c),
		Audit:        audit.New(c),
		Budget:       budget.New(c),
		Coordination: coordination.New(c),
		Lifecycle:    lifecycle.New(c),
		// Full platform
		Agents:        agents.New(c),
		Sessions:      sessions.New(c),
		Templates:     templates.New(c),
		Users:         users.New(c),
		Deployments:   deployments.New(c),
		Tasks:         tasks.New(c),
		Evals:         evals.New(c),
		Groups:        groups.New(c),
		OAuth2:        oauth2.New(c),
		Notifications: notifications.New(c),
		Files:         files.New(c),
		Insights:      insights.New(c),
		Provisioners:  provisioners.New(c),
		ExternalAuth:  externalauth.New(c),
		Licenses:      licenses.New(c),
		Replicas:      replicas.New(c),
		GitSSHKeys:    gitsshkeys.New(c),
		// Agent-side
		AgentSDK:   agentsdk.New(c),
		SidecarSDK: sidecarsdk.New(c),
		Health:     healthsdk.New(c),
		// Internals
		Client: c,
		config: cfg,
	}, nil
}

// Run registers the stack, starts heartbeats, and blocks until the context
// is cancelled or a shutdown signal is received. It handles graceful shutdown
// including deregistration.
func (s *Stack) Run(ctx context.Context) error {
	// Register with Runtime.
	reg, err := s.Lifecycle.RegisterStack(ctx, lifecycle.RegisterRequest{
		Name:     s.config.StackName,
		Version:  s.config.StackVersion,
		Endpoint: s.config.Endpoint,
	})
	if err != nil {
		return fmt.Errorf("register stack: %w", err)
	}
	s.registration = &reg
	log.Printf("lattice: registered stack %q (id=%s)", s.config.StackName, reg.ID)

	// Create a context that cancels on OS signals.
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start heartbeat loop.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("lattice: shutting down stack %q", s.config.StackName)
			// Attempt graceful deregistration with a fresh context.
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := s.Lifecycle.Deregister(shutdownCtx, reg.ID); err != nil {
				log.Printf("lattice: deregister error: %v", err)
			} else {
				log.Printf("lattice: deregistered stack %q", s.config.StackName)
			}
			return ctx.Err()
		case <-ticker.C:
			if err := s.Lifecycle.Heartbeat(ctx, reg.ID); err != nil {
				log.Printf("lattice: heartbeat error: %v", err)
			}
		}
	}
}

// Registration returns the stack registration, or nil if Run hasn't been called.
func (s *Stack) Registration() *lifecycle.StackRegistration {
	return s.registration
}

// NewFromEnv creates a Stack from environment variables.
// See LoadConfigFromEnv for the required variables.
func NewFromEnv() (*Stack, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return New(cfg)
}

// init suppresses the log timestamp prefix for cleaner output.
func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}
