files := $(shell find . -name vendor -prune -o -name '*.go' -print)
TESTARGS ?= -cover

test: vendor vet
	@go test ${TESTARGS} ./...

fmt:
	@gofmt -l -s -w $(files)

vendor:
	@dep ensure

vet:
	@go vet ./...

.PHONY: test fmt vendor vet
