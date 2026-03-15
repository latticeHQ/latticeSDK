# Lattice SDK API Reference

## Package Overview

| Package | Import | Description |
|---------|--------|-------------|
| `stack` | `github.com/latticehq/latticesdk/stack` | High-level entry point |
| `client` | `github.com/latticehq/latticesdk/client` | HTTP client |
| `types` | `github.com/latticehq/latticesdk/types` | Domain types |
| `identity` | `github.com/latticehq/latticesdk/identity` | Identity management |
| `authz` | `github.com/latticehq/latticesdk/authz` | Authorization |
| `audit` | `github.com/latticehq/latticesdk/audit` | Audit trail |
| `budget` | `github.com/latticehq/latticesdk/budget` | Budget control |
| `coordination` | `github.com/latticehq/latticesdk/coordination` | Cross-stack messaging |
| `lifecycle` | `github.com/latticehq/latticesdk/lifecycle` | Stack lifecycle |

## stack

### `stack.New(cfg Config) (*Stack, error)`
Creates a new Stack with all services pre-wired.

### `stack.NewFromEnv() (*Stack, error)`
Creates a Stack from environment variables.

### `(*Stack).Run(ctx context.Context) error`
Registers, heartbeats, and blocks until shutdown.

## identity

### `(*Service).ListAgents(ctx, filter) ([]Agent, int, error)`
List agents matching the filter.

### `(*Service).GetAgent(ctx, id) (Agent, error)`
Get a single agent by ID.

### `(*Service).CreateAgent(ctx, orgID, req) (Agent, error)`
Create a new agent.

### `(*Service).GetCurrentUser(ctx) (User, error)`
Get the authenticated user.

### `(*Service).ListUsers(ctx, search, status, pagination) ([]User, int, error)`
List users with optional filters.

### `(*Service).ListOrganizations(ctx) ([]Organization, error)`
List all visible organizations.

### `(*Service).CreateToken(ctx, req) (GenerateAPIKeyResponse, error)`
Create a new API token.

### `(*Service).ListTokens(ctx) ([]APIKey, error)`
List all tokens for the current user.

### `(*Service).RevokeToken(ctx, keyID) error`
Delete a token.

## authz

### `(*Service).Can(ctx, action, resourceType, resourceID) (bool, error)`
Simple authorization check.

### `(*Service).Check(ctx, req) (AuthorizationResponse, error)`
Batch authorization check.

### `(*Service).ListSiteRoles(ctx) ([]AssignableRoles, error)`
List site-wide roles.

## audit

### `(*Service).Query(ctx, filter, pagination) ([]AuditLog, int64, error)`
Query audit logs.

### `(*Service).Emit(ctx, req) error`
Emit a custom audit event.

## budget

### `(*Service).GetQuota(ctx, ownerID) (AgentQuota, error)`
Get budget quota for an owner.

### `(*Service).CheckBudget(ctx, ownerID) (bool, error)`
Check if owner has remaining budget.

## coordination

### `(*Service).Publish(ctx, req) error`
Publish an event to a topic.

### `(*Service).Subscribe(ctx, topic) (func() (*SSE, error), error)`
Subscribe to events on a topic.

### `(*Service).GetSharedState(ctx, key) (SharedState, error)`
Get shared state by key.

### `(*Service).SetSharedState(ctx, key, value) error`
Set shared state.

### `(*Service).Escalate(ctx, req) error`
Escalate work to another stack.

## lifecycle

### `(*Service).RegisterStack(ctx, req) (StackRegistration, error)`
Register a stack with Runtime.

### `(*Service).Deregister(ctx, stackID) error`
Deregister a stack.

### `(*Service).Heartbeat(ctx, stackID) error`
Send a heartbeat.

### `(*Service).ReportHealth(ctx, report) error`
Report stack health status.

### `(*Service).GetMetadata(ctx, stackID) (Metadata, error)`
Get stack metadata.

### `(*Service).SetMetadata(ctx, stackID, meta) error`
Set stack metadata.
