# Proposal

## Why

The repository has no automated integration coverage. CI only runs unit tests and linting. Kafka-dependent behavior, SASL, TLS, and AWS Secrets Manager flows stay untested until manual runs. This slows reviews and hides regressions.

## What Changes

- Add a GitHub Actions workflow that boots Kafka, Vault, and LocalStack services.
- Add a reusable Makefile target for integration tests with build tags.
- Add integration test scaffolding under `tests/integration/`.
- Provide docker-compose configuration tuned for CI runners.
- Add test fixtures and helpers for TLS material and SASL credentials.
- Gate integration jobs on pull requests and pushes to main.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none)

## Impact

- Affected files: `.github/workflows/ci.yaml`, `Makefile`, `docker-compose.yaml`, `tests/integration/`.
- No library behavior changes; CI-only change.