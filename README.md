<div align="center">

# Lattice SDK

### Build Department Stacks on the Lattice coordination layer.

**The official Go SDK for [Lattice Runtime](https://github.com/latticeHQ/latticeRuntime) — identity, authorization, audit, budget, and coordination for institutional AI.**

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

    // List agents using your template
    agents, err := s.Identity.ListAgents(ctx, nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d agents\n", len(agents))

    // Check authorization
    allowed, err := s.Authz.Can(ctx, "update", "agent", agentID)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Allowed: %v\n", allowed)

    // Run with heartbeats and graceful shutdown
    s.Run(ctx)
}
```

## Packages

| Package | Description |
|---------|-------------|
| `client` | HTTP client with authentication, retries, and error handling |
| `types` | Domain types — agents, users, organizations, templates, roles |
| `identity` | Identity management — agents, users, organizations, tokens |
| `authz` | Authorization — permission checks, RBAC |
| `audit` | Audit trail — query and emit events |
| `budget` | Budget control — quotas, cost reporting |
| `coordination` | Cross-stack messaging — pub/sub, shared state |
| `lifecycle` | Stack lifecycle — registration, health, metadata |
| `stack` | High-level bootstrap — all services pre-wired with `Run()` |

## The Ecosystem

| Component | What it does | License |
|-----------|-------------|---------|
| [**Runtime**](https://github.com/latticeHQ/latticeRuntime) | Coordination layer — identity, authorization, audit, budget, networking | Apache 2.0 |
| [**SDK**](https://github.com/latticeHQ/latticeSDK) | Go SDK for building Department Stacks | Apache 2.0 |
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
