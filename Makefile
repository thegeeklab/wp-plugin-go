# renovate: datasource=github-releases depName=mvdan/gofumpt
GOFUMPT_PACKAGE_VERSION := v0.9.2
# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_PACKAGE_VERSION := v2.5.0

GO ?= go
PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)

GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@$(GOFUMPT_PACKAGE_VERSION)
GOTESTSUM_PACKAGE ?= gotest.tools/gotestsum@latest

GENERATE ?=

.PHONY: fmt
fmt:
	$(shell go env GOPATH)/bin/gofumpt -extra -w $(SOURCES)

.PHONY: golangci-lint
golangci-lint:
	$(shell go env GOPATH)/bin/golangci-lint run

.PHONY: lint
lint: golangci-lint

.PHONY: generate
generate:
	$(GO) generate $(GENERATE)

.PHONY: test
test:
	$(shell go env GOPATH)/bin/gotestsum --no-color=false -- -coverprofile=coverage.out $(PACKAGES)

.PHONY: deps
deps:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_PACKAGE_VERSION)
	$(GO) mod download
	$(GO) install $(GOFUMPT_PACKAGE)
	$(GO) install $(GOTESTSUM_PACKAGE)
