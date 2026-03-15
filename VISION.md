# Lattice SDK — Vision

## The Problem

Every organization deploying AI agents needs the same infrastructure: identity, authorization, audit trails, budget controls, and cross-agent coordination. But building this from scratch for each vertical application is wasteful and error-prone.

Lattice Runtime provides this horizontal coordination layer. The SDK makes it accessible.

## The Vision

**Make building institutional AI applications as easy as building a web app with Stripe.**

Developers should be able to:
1. `go get` the SDK
2. Point it at a Lattice Runtime instance
3. Build their domain logic — not infrastructure

The SDK is the interface between **Department Stacks** (vertical AI applications) and the **Runtime** (horizontal coordination layer).

## Design Principles

### 1. Zero Infrastructure Knowledge Required
Stack developers shouldn't need to understand Runtime internals. The SDK abstracts away HTTP, authentication, pagination, and error handling.

### 2. One Import, All Services
The `stack` package provides a single entry point with all services pre-wired. No dependency injection, no service locators — just `s.Identity`, `s.Authz`, `s.Budget`.

### 3. Idiomatic Go
Standard patterns: functional options, context propagation, structured errors, iterator functions. No magic, no code generation.

### 4. Minimal Dependencies
Only `github.com/google/uuid` beyond the standard library. The SDK should never force dependency choices on stack developers.

### 5. Runtime Version Independence
The SDK targets the stable Runtime API (`/api/v2/`). It should work across Runtime versions without requiring lock-step upgrades.

## What Is a Department Stack?

A Department Stack is a vertical AI application that runs on Lattice Runtime. Examples:

- **HR Stack**: AI agents for recruiting, onboarding, benefits, employee support
- **Legal Stack**: Contract review, compliance monitoring, regulatory tracking, IP management
- **Finance Stack**: Expense approval, forecasting, audit automation, tax preparation
- **Security Stack**: Threat detection, incident response, vulnerability assessment
- **Engineering Stack**: Code review, deployment automation, on-call management

Each stack:
- Registers with Runtime on startup
- Authenticates via API key or session token
- Uses Runtime for identity, authorization, and audit
- Can communicate with other stacks via coordination APIs
- Reports health and cost metrics back to Runtime

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
