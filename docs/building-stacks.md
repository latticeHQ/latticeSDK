# Building Department Stacks

## What Is a Department Stack?

A Department Stack is a vertical AI application built on Lattice Runtime. It gets identity, authorization, audit, budget, and coordination for free — you focus on your domain logic.

## Architecture

```
┌─────────────────────────────────────────┐
│              Your Stack                  │
│  ┌──────────┐  ┌──────────┐  ┌────────┐│
│  │ Domain   │  │  HTTP    │  │ Worker ││
│  │ Logic    │  │  Server  │  │ Loops  ││
│  └────┬─────┘  └────┬─────┘  └───┬────┘│
│       │              │             │      │
│  ┌────┴──────────────┴─────────────┴────┐│
│  │          Lattice SDK                  ││
│  │  Identity│Authz│Audit│Budget│Coord   ││
│  └────────────────┬──────────────────────┘│
└───────────────────┼───────────────────────┘
                    │ HTTP/SSE
┌───────────────────┼───────────────────────┐
│          Lattice Runtime                   │
│  Identity │ RBAC │ Audit │ Budget │ Net   │
└────────────────────────────────────────────┘
```

## Step-by-Step Guide

### 1. Initialize

```go
s, err := stack.New(stack.Config{
    RuntimeURL:   "http://localhost:3000",
    APIKey:       os.Getenv("LATTICE_API_KEY"),
    StackName:    "hr-stack",
    StackVersion: "0.1.0",
})
```

### 2. Use Identity

```go
// Who am I?
me, _ := s.Identity.GetCurrentUser(ctx)

// List agents running my template
agents, _, _ := s.Identity.ListAgents(ctx, &types.AgentFilter{
    Template: "hr-agent",
})
```

### 3. Check Authorization

```go
// Can this user update the agent?
allowed, _ := s.Authz.Can(ctx, types.RBACActionUpdate, types.RBACResourceAgent, agentID)

// Batch check
result, _ := s.Authz.Check(ctx, types.AuthorizationRequest{
    Checks: map[string]types.AuthorizationCheck{
        "read":   {Object: obj, Action: types.RBACActionRead},
        "delete": {Object: obj, Action: types.RBACActionDelete},
    },
})
```

### 4. Audit Actions

```go
// Query audit history
logs, _, _ := s.Audit.Query(ctx, &types.AuditLogFilter{
    SearchQuery: "action:delete",
}, nil)

// Emit custom events
s.Audit.Emit(ctx, audit.EmitEventRequest{
    Action:       types.AuditActionWrite,
    ResourceType: "contract_review",
})
```

### 5. Check Budget

```go
hasQuota, _ := s.Budget.CheckBudget(ctx, orgID)
if !hasQuota {
    log.Println("budget exceeded")
}
```

### 6. Coordinate with Other Stacks

```go
// Publish events
s.Coordination.Publish(ctx, coordination.PublishRequest{
    Topic:   "hr.onboarding-complete",
    Payload: map[string]string{"employee_id": "123"},
})

// Subscribe to events
next, _ := s.Coordination.Subscribe(ctx, "finance.budget-alert")
for {
    event, err := next()
    if err != nil {
        break
    }
    // Handle event...
}

// Escalate to another stack
s.Coordination.Escalate(ctx, coordination.EscalationRequest{
    TargetStack: "legal-stack",
    Action:      "review-contract",
    Payload:     contractData,
})
```

### 7. Run the Stack

```go
// Handles registration, heartbeats, and graceful shutdown
if err := s.Run(ctx); err != nil {
    log.Fatal(err)
}
```

## Serving HTTP

If your stack exposes an HTTP API:

```go
mux := http.NewServeMux()
mux.HandleFunc("/api/hr/employees", handleEmployees)

// Add Lattice auth middleware
handler := s.AuthMiddleware(mux)

go http.ListenAndServe(":8080", handler)
s.Run(ctx)
```

## Testing

Use the `client` package directly for unit tests:

```go
c, _ := client.New("http://test-runtime:3000", client.WithAPIKey("test-key"))
svc := identity.New(c)
agents, _, _ := svc.ListAgents(ctx, nil)
```
