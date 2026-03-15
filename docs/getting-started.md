# Getting Started with the Lattice SDK

## Prerequisites

- Go 1.22 or later
- A running Lattice Runtime instance
- An API key or session token

## Installation

```bash
go get github.com/latticehq/latticesdk@latest
```

## Quick Start

### 1. Create a new stack

```bash
go run github.com/latticehq/latticesdk/cmd/lattice-stack@latest init my-stack
cd my-stack
```

### 2. Configure your environment

```bash
cp .env.example .env
# Edit .env with your Runtime URL and API key
```

### 3. Run your stack

```bash
go run .
```

## Manual Setup

If you prefer to set up manually:

```go
package main

import (
    "context"
    "log"
    "os"

    "github.com/latticehq/latticesdk/stack"
)

func main() {
    ctx := context.Background()

    s, err := stack.New(stack.Config{
        RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
        APIKey:       os.Getenv("LATTICE_API_KEY"),
        StackName:    "my-stack",
        StackVersion: "0.1.0",
    })
    if err != nil {
        log.Fatal(err)
    }

    // Use services: s.Identity, s.Authz, s.Audit, s.Budget, s.Coordination
    _ = s

    // Run with heartbeats and graceful shutdown
    if err := s.Run(ctx); err != nil {
        log.Fatal(err)
    }
}
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `LATTICE_RUNTIME_URL` | Yes | Base URL of the Lattice Runtime API |
| `LATTICE_API_KEY` | Yes* | API key for authentication |
| `LATTICE_SESSION_TOKEN` | Yes* | Session token (alternative to API key) |
| `LATTICE_STACK_NAME` | Yes | Human-readable name for this stack |
| `LATTICE_STACK_VERSION` | No | Version string (defaults to "0.0.0") |
| `LATTICE_STACK_ENDPOINT` | No | URL where other stacks can reach this one |

*Either `LATTICE_API_KEY` or `LATTICE_SESSION_TOKEN` is required.

## Next Steps

- [API Reference](api-reference.md) — Complete SDK API documentation
- [Building Stacks](building-stacks.md) — Guide to building Department Stacks
- [Examples](../examples/) — Working example stacks
