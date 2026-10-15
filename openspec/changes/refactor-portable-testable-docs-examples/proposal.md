# Proposal

## Why

Kilo-franz ships Kafka consumer and producer helpers built on segmentio/kafka-go. The public API ties callers to concrete structs, hides configuration behind unexported fields, and lacks runnable examples. These limits make the library hard to embed in other Go services and hard to unit test in isolation.

## What Changes

- Replace concrete consumer and producer dependencies with small interface contracts.
- Export configuration builders so callers can assemble configs from their own inputs.
- Add deterministic fakes for Kafka readers and writers to support unit tests.
- Introduce table-driven unit tests for core packages (`config`, `core`, `encoder`, `publisher`, `subscriber`).
- Refresh package documentation and Go doc comments for all public APIs.
- Add an `examples/` directory with minimal, runnable consumer and producer programs.
- Restructure package layout to keep platform-specific code behind interfaces.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

(none)

## Impact

- Affected packages: `config`, `core`, `encoder`, `logging`, `monitoring`, `publisher`, `refresh`, `sasl`, `subscriber`, `tls`, `vault`.
- Public API changes are additive; existing callers keep compiling.
- New example binaries live under `examples/` and do not ship in the library import path.