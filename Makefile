files := $(shell find . -name vendor -prune -o -name '*.go' -print)

test: vendor
	@go test -cover ./...

fmt:
	@gofmt -l -s -w $(files)

vendor:
	@dep ensure

.PHONY: test fmt vendor
