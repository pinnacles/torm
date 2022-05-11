setup.staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: lint
lint: setup.staticcheck
	staticcheck .

.PHONY: test
test:
	go test -race -v ./...

.PHONY: build
build:
	@echo "not implemented"

.PHONY: post-process
post-process:
	@echo "not implemented"
