# renovate: datasource=github-releases depName=mvdan/gofumpt
GOFUMPT_PACKAGE_VERSION := v0.7.0
# renovate: datasource=github-releases depName=golangci/golangci-lint
GOLANGCI_LINT_PACKAGE_VERSION := v1.62.0

GO ?= go
PACKAGES ?= $(shell go list ./...)
SOURCES ?= $(shell find . -name "*.go" -type f)

GOFUMPT_PACKAGE ?= mvdan.cc/gofumpt@$(GOFUMPT_PACKAGE_VERSION)
GOLANGCI_LINT_PACKAGE ?= github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_PACKAGE_VERSION)
GOTESTSUM_PACKAGE ?= gotest.tools/gotestsum@latest

GENERATE ?=

.PHONY: fmt
fmt:
	$(GO) run $(GOFUMPT_PACKAGE) -extra -w $(SOURCES)

.PHONY: golangci-lint
golangci-lint:
	$(GO) run $(GOLANGCI_LINT_PACKAGE) run

.PHONY: lint
lint: golangci-lint

.PHONY: generate
generate:
	$(GO) generate $(GENERATE)

.PHONY: test
test:
	$(GO) run $(GOTESTSUM_PACKAGE) --no-color=false -- -coverprofile=coverage.out $(PACKAGES)

.PHONY: deps
deps:
	$(GO) mod download
	$(GO) install $(GOFUMPT_PACKAGE)
	$(GO) install $(GOLANGCI_LINT_PACKAGE)
	$(GO) install $(GOTESTSUM_PACKAGE)
