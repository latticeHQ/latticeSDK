<div align="center">

# Lattice SDK

### The complete Go SDK for Lattice Runtime.

**380 functions across 25 packages — full API parity with the platform.**
**Build Department Stacks without touching Runtime source code.**

---

</div>

## What is a Department Stack?

A **Department Stack** is a vertical AI application built on top of Lattice Runtime's horizontal coordination layer. Think of it like building a Shopify app on the Shopify platform — you get identity, authorization, audit trails, budgets, and cross-stack coordination for free.

Examples:
- **HR Stack** — AI agents for recruiting, onboarding, and employee support
- **Legal Stack** — Contract review, compliance monitoring, regulatory tracking
- **Finance Stack** — Expense approval, forecasting, audit automation

## Install

```bash
go get github.com/latticehq/latticesdk@latest
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/latticehq/latticesdk/stack"
)

func main() {
    ctx := context.Background()

    s, err := stack.New(stack.Config{
        RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
        APIKey:       os.Getenv("LATTICE_API_KEY"),
        StackName:    "hr-stack",
        StackVersion: "0.1.0",
    })
    if err != nil {
        panic(err)
    }

    // Full platform access through typed services
    agents, _, _ := s.Agents.ListAgents(ctx, nil)
    fmt.Printf("Found %d agents\n", len(agents))

    templates, _ := s.Templates.ListTemplates(ctx, nil)
    fmt.Printf("Found %d templates\n", len(templates))

    allowed, _ := s.Authz.Can(ctx, "update", "agent", "some-id")
    fmt.Printf("Allowed: %v\n", allowed)

    // Run with heartbeats and graceful shutdown
    s.Run(ctx)
}
```

## Packages

### Core Coordination
| Package | Funcs | Description |
|---------|-------|-------------|
| `client` | 22 | HTTP client — auth, retries, SSE, pagination, error handling |
| `types` | — | Shared domain types — agents, users, orgs, templates, roles |
| `stack` | 8 | High-level bootstrap — all 25 services pre-wired with `Run()` |

### Identity & Access
| Package | Funcs | Description |
|---------|-------|-------------|
| `identity` | 25 | Identity management — agents, users, organizations, tokens |
| `users` | 38 | Full user lifecycle — create, auth, roles, password, org membership |
| `authz` | 5 | Authorization — permission checks, RBAC, site/org roles |
| `groups` | 17 | Group management — CRUD, IDP sync settings |
| `oauth2` | 12 | OAuth2 provider apps — secrets, revocation |
| `externalauth` | 8 | External auth providers — device auth, link/unlink |

### Agents & Sessions
| Package | Funcs | Description |
|---------|-------|-------------|
| `agents` | 44 | Complete agent management — CRUD, builds, proxies, sidecars, port shares |
| `sessions` | 45 | Complete session management — CRUD, builds, sidecars, real-time |
| `templates` | 39 | Template management — CRUD, versions, ACL, parameters, dry runs |

### Platform Operations
| Package | Funcs | Description |
|---------|-------|-------------|
| `deployments` | 23 | Deployment config, stats, appearance, entitlements, experiments |
| `tasks` | 12 | AI task management — CRUD, send messages, logs |
| `evals` | 13 | Evaluation framework — runs, comparisons, passes, parameters |
| `insights` | 9 | Usage analytics — latency, activity, template insights |
| `notifications` | 8 | Notification settings, templates, user preferences |
| `provisioners` | 10 | Provisioner daemons and key management |
| `licenses` | 8 | License management — add, list, delete |

### Coordination & Lifecycle
| Package | Funcs | Description |
|---------|-------|-------------|
| `coordination` | 6 | Cross-stack messaging — pub/sub, shared state, escalation |
| `lifecycle` | 7 | Stack lifecycle — registration, health, heartbeats, metadata |
| `audit` | 4 | Audit trail — query logs, emit custom events |
| `budget` | 4 | Budget control — quotas, cost reporting |

### Utilities
| Package | Funcs | Description |
|---------|-------|-------------|
| `files` | 4 | File upload and download |
| `gitsshkeys` | 3 | Git SSH key management |
| `replicas` | 2 | Replica information |

## Architecture

```
┌─────────────────────────────────────────────────┐
│                  Your Department Stack            │
│  (HR, Legal, Finance, Security, Engineering...)   │
├─────────────────────────────────────────────────┤
│                   latticeSDK                      │
│  ┌─────────┐ ┌──────────┐ ┌──────────────────┐  │
│  │ agents  │ │ sessions │ │   templates      │  │
│  │ users   │ │ tasks    │ │   deployments    │  │
│  │ authz   │ │ evals    │ │   coordination   │  │
│  │ groups  │ │ audit    │ │   lifecycle      │  │
│  │ oauth2  │ │ budget   │ │   provisioners   │  │
│  └─────────┘ └──────────┘ └──────────────────┘  │
│                    client                         │
│            (HTTP, SSE, pagination)                │
├─────────────────────────────────────────────────┤
│              Lattice Runtime API                  │
│         access.latticeruntime.com                 │
└─────────────────────────────────────────────────┘
```

**Zero dependency on Runtime source code.** Only external dependency: `github.com/google/uuid`.

## Using Individual Services

You don't have to use the `stack` bootstrap. Each service works independently:

```go
import (
    "github.com/latticehq/latticesdk/client"
    "github.com/latticehq/latticesdk/agents"
    "github.com/latticehq/latticesdk/templates"
)

c, _ := client.New("https://access.latticeruntime.com", client.WithAPIKey("..."))

agentSvc := agents.New(c)
tmplSvc := templates.New(c)

// Use directly
myAgents, _, _ := agentSvc.ListAgents(ctx, nil)
myTemplates, _ := tmplSvc.ListTemplates(ctx, nil)
```

## The Ecosystem

| Component | What it does | License |
|-----------|-------------|---------|
| [**Runtime**](https://github.com/latticeHQ/latticeRuntime) | Coordination layer — identity, authorization, audit, budget, networking | Apache 2.0 |
| [**SDK**](https://github.com/latticeHQ/latticeSDK) | Go SDK for building Department Stacks (380 functions, 25 packages) | Apache 2.0 |
| [**Workbench**](https://github.com/latticeHQ/latticeWorkbench) | Reference Engineering Stack — multi-model agent workspace | MIT |
| [**Inference**](https://github.com/latticeHQ/latticeInference) | Local AI serving — MLX on Apple Silicon, zero-config clustering | Apache 2.0 |
| [**Operator**](https://github.com/latticeHQ/latticeOperator) | Self-hosted deployment management for Lattice infrastructure | Apache 2.0 |
| [**Registry**](https://github.com/latticeHQ/latticeRegistry) | Community ecosystem — Terraform modules, templates, department stacks | Apache 2.0 |
| [**Terraform Provider**](https://github.com/latticeHQ/terraform-provider-lattice) | Infrastructure as code for Lattice deployments | MPL 2.0 |
| [**Toolbox**](https://github.com/latticeHQ/LatticeToolbox) | macOS app manager for Lattice products | MIT |
| [**Homebrew**](https://github.com/latticeHQ/latticeHomebrew) | One-line install on macOS and Linux | MIT |
| [**Enterprise**](https://github.com/latticeHQ/latticeEnterprise) | Enterprise administration and governance | Coming soon |

## Links

- **Website**: [latticeruntime.com](https://latticeruntime.com)
- **Runtime**: [github.com/latticeHQ/latticeRuntime](https://github.com/latticeHQ/latticeRuntime)
- **Security**: security@latticeruntime.com

---

<div align="center">

**Your agents. Your coordination. Your rules. Your infrastructure.**

</div>
