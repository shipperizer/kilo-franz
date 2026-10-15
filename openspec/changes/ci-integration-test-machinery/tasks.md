# Tasks

## 1. Workflow

- [ ] 1.1 Create integration GitHub Actions workflow and verify it triggers on pull requests
- [ ] 1.2 Define service containers for Kafka Vault and LocalStack and verify health checks pass
- [ ] 1.3 Add integration job to CI and verify unit job still runs

## 2. Test Scaffolding

- [ ] 2.1 Create tests/integration directory with smoke test and verify build tag works
- [ ] 2.2 Add test helpers for Kafka bootstrap and verify connection succeeds
- [ ] 2.3 Add fixtures for TLS and SASL and verify loading works

## 3. Makefile and Compose

- [ ] 3.1 Add integration Makefile target and verify local run works
- [ ] 3.2 Update docker-compose for CI and verify services start
- [ ] 3.3 Document local integration run and verify README steps

## 4. Verification

- [ ] 4.1 Run go test -tags=integration on CI and verify smoke test passes
- [ ] 4.2 Confirm flaky retry handling and verify rerun works

## Documentation and rollout

- [ ] 5.1 Update CHANGELOG and verify CI changes noted

## Implementation Notes

(none yet)