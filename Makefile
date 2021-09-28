.PHONY=test mocks

GO111MODULE?=on
GOPRIVATE?=github.com/sbs-discovery-sweden/*
CGO_ENABLED?=0
GOOS?=linux
GO?=go
MAIN_DIR?=/var/app
AWS?=aws
AWS_PAGER?=
SM_HOST?=http://localhost:14566
WAIT_TIME_MS?=2000

.EXPORT_ALL_VARIABLES:


test: mocks
	$(GO) test -v ./... -cover -coverprofile coverage.out 
	# this will be cached, just needed to the test.json
	$(GO) test -v ./... -cover -coverprofile coverage.out -json > test.json

secretsmanager:
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/password" --description "test password" --secret-string '1234'
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/key" --description "test key" --secret-string file://test.key
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/cert" --description "test cert" --secret-string file://test.cert
	- $(AWS) secretsmanager create-secret --endpoint-url $(SM_HOST) --name "test/bin" --description "test binary" --secret-binary fileb://test.bin
	- $(AWS) secretsmanager put-secret-value --endpoint-url $(SM_HOST) --secret-id "test/password" --secret-string '1234'
	- $(AWS) secretsmanager put-secret-value --endpoint-url $(SM_HOST) --secret-id "test/key" --secret-string file://test.key
	- $(AWS) secretsmanager put-secret-value --endpoint-url $(SM_HOST) --secret-id "test/cert" --secret-string file://test.cert
	- $(AWS) secretsmanager put-secret-value --endpoint-url $(SM_HOST) --secret-id "test/bin" --secret-binary fileb://test.bin

mocks:
	$(GO) get github.com/golang/mock/gomock
	$(GO) get github.com/golang/mock/mockgen
	# generate gomocks
	$(GO) generate ./...

