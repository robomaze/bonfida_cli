BUFFER := $(shell mktemp)
DIST_DIR=./dist
REPORT_DIR=$(DIST_DIR)/report
IGNORE_FILES=

default: lint test

deps:
	@echo "Installing dependencies"
	go get ./...

lint:
	@echo "Checking code style"
	test -f $(GOPATH)/bin/gofumpt || go get mvdan.cc/gofumpt
	gofumpt -l . | tee $(BUFFER)
	@! test -s $(BUFFER)
	go vet ./...
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.1
	golangci-lint run

lint-fix:
	test -f $(GOPATH)/bin/gofumpt || go get mvdan.cc/gofumpt
	gofumpt -l -w .

test:
	@echo "Running unit tests"
	mkdir -p $(REPORT_DIR)
	go test ./...

test-race:
	@echo "Running unit tests and checking for race conditions"
	mkdir -p $(REPORT_DIR)
	go test ./... -race
