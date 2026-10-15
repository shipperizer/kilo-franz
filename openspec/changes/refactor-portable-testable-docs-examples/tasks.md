# Tasks

## 1. Interfaces

- [ ] 1.1 Define ReaderInterface in core and verify unit tests compile against it
- [ ] 1.2 Define WriterInterface in core and verify unit tests compile against it
- [ ] 1.3 Define MonitorInterface and verify publisher compiles with it

## 2. Fakes

- [ ] 2.1 Add FakeReader in new testkit package and verify it satisfies ReaderInterface
- [ ] 2.2 Add FakeWriter in testkit and verify it satisfies WriterInterface
- [ ] 2.3 Add FakeMonitor in testkit and verify log output captured

## 3. Configuration

- [ ] 3.1 Export config constructor helpers and verify go vet passes
- [ ] 3.2 Add validation errors to NewConfig and verify unit tests cover them
- [ ] 3.3 Document config fields and verify godoc renders without warnings

## 4. Package Refactor

- [ ] 4.1 Update subscriber to depend on ReaderInterface and verify tests pass
- [ ] 4.2 Update publisher to depend on WriterInterface and verify tests pass
- [ ] 4.3 Remove unused concrete dependencies and verify build passes

## 5. Tests

- [ ] 5.1 Add table tests for config package and verify coverage rises
- [ ] 5.2 Add table tests for encoder package and verify round trip cases
- [ ] 5.3 Add tests for sasl builders and verify error cases

## 6. Documentation and Examples

- [ ] 6.1 Refresh README with usage and verify examples compile
- [ ] 6.2 Add examples/consumer and verify make example runs
- [ ] 6.3 Add examples/producer and verify make example runs

## Verification suite

- [ ] 7.1 Run go test ./... and verify all packages pass
- [ ] 7.2 Run go vet ./... and verify no findings

## Documentation and rollout

- [ ] 8.1 Update CHANGELOG and verify release notes list new APIs
- [ ] 8.2 Tag pre-release and verify examples build against it

## Implementation Notes

(none yet)