GO111MODULE?=on
GOPRIVATE?=github.com/shipperizer/*
CGO_ENABLED?=0
GOOS?=linux
GO?=go
MAIN_DIR?=/var/app
AWS?=aws
AWS_PAGER?=
SM_HOST?=http://localhost:14566
WAIT_TIME_MS?=2000
KAFKA_CLI_PATH?=/usr/local/bin
KAFKA_HOST?=kafka:9092
TOPIC_NAME?=test

.EXPORT_ALL_VARIABLES:


mocks: vendor
	$(GO) install github.com/golang/mock/mockgen@v1.6.0
	# generate gomocks
	$(GO) generate ./...
.PHONY: mocks

test: mocks vet
	$(GO) test ./... -cover -coverprofile coverage_source.out
	# this will be cached, just needed to the test.json
	$(GO) test ./... -cover -coverprofile coverage_source.out -json > test_source.json
	cat coverage_source.out | grep -v "mock_*" | tee coverage.out
	cat test_source.json | grep -v "mock_*" | tee test.json
.PHONY: test

vet:
	$(GO) vet ./...
.PHONY: vet

vendor:
	$(GO) mod vendor
.PHONY: vendor

secretsmanager:
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/password" --description "test password" --secret-string '1234'
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/key" --description "test key" --secret-string file://test.key
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/cert" --description "test cert" --secret-string file://test.cert
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/p12" --description "test p12" --secret-binary fileb://test.p12
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/p12.sha256" --description "test p12 sha256" --secret-binary fileb://test.p12.sha256
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/bin" --description "test binary" --secret-binary fileb://test.bin
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "sasl/credentials" --description "sasl credentials" --secret-string file://sasl.json
.PHONY: secretsmanager

topics:
	$(KAFKA_CLI_PATH)/kafka-topics.sh --create --bootstrap-server $(KAFKA_HOST) --replication-factor 1 --partitions 100 --topic $(TOPIC_NAME)
.PHONY: topics
