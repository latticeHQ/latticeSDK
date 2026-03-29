<div align="center">

# Lattice SDK

### Go SDK for building on Lattice Runtime

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache_2.0-blue.svg?style=flat-square)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat-square&logo=go)](https://go.dev)

**Build applications that inherit crash-proof execution, cryptographic audit, and governance from Lattice Runtime.**

[Lattice Runtime](https://github.com/latticeHQ/latticeRuntime) · [Docs](https://docs.latticeruntime.com) · [Discussions](https://github.com/latticeHQ/latticeRuntime/discussions)

</div>

---

## What It Does

The Lattice SDK lets you build Go applications on top of Lattice Runtime. Your application automatically inherits:

- **Durable execution** — embedded Temporal, crash-proof workflows
- **Cryptographic audit** — every action hash-chained and immutable
- **Budget enforcement** — per-workload spending limits
- **Identity & authorization** — RBAC + ABAC via Rego
- **Zero-trust networking** — WireGuard mesh

You write domain logic. Lattice handles governance.

## Install

```bash
go get github.com/latticeHQ/latticeSDK
```

## Usage

```go
import "github.com/latticeHQ/latticeSDK/agentsdk"

client := agentsdk.New(agentsdk.Config{
    RuntimeURL: "https://your-lattice-instance.com",
    Token:      os.Getenv("LATTICE_TOKEN"),
})

// Create an agent — it inherits all runtime governance
agent, err := client.CreateAgent(ctx, agentsdk.CreateAgentRequest{
    Name:     "my-worker",
    Template: "default",
})
```

## Part of the Lattice Ecosystem

| Component | Role |
|-----------|------|
| [**Runtime**](https://github.com/latticeHQ/latticeRuntime) | Crash-proof runtime — identity, auth, audit, budget, mesh |
| [**Workbench**](https://github.com/latticeHQ/latticeWorkbench) | 316K-line multi-model agent workspace |
| **SDK** (this repo) | Go SDK for building on Lattice |

---

<div align="center">

**[latticeruntime.com](https://latticeruntime.com)** — Crash-proof governed runtime for AI agents.

</div>
