# Contributing to Lattice SDK

We welcome contributions to the Lattice SDK.

## Getting Started

1. Fork the repository
2. Create a feature branch from `develop`
3. Make your changes
4. Run checks: `make check`
5. Submit a pull request to `develop`

## Development

```bash
# Build
make build

# Test
make test

# Lint
make lint

# All checks
make check
```

## Guidelines

- Follow existing code patterns
- Add tests for new functionality
- Keep dependencies minimal (stdlib + uuid only)
- Use functional options for configuration
- Propagate context through all public APIs

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
