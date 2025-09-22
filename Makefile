########################
### Makefile Helpers ###
########################

# Include modular makefiles
include makefiles/benchmark.mk
include makefiles/build.mk

.PHONY: prompt_user
prompt_user: ## Internal helper target - prompt the user before continuing
	@echo "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: help
.DEFAULT_GOAL := help
help: ## Prints all the targets in all the Makefiles
	@echo ""
	@echo "\033[1;34m📋 Ring-Go Makefile Targets\033[0m"
	@echo ""
	@echo "\033[1;34m=== 🔍 Information & Discovery ===\033[0m"
	@grep -h -E '^help:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-58s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;34m=== 🧪 Testing ===\033[0m"
	@grep -h -E '^test.*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-58s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;34m=== ⚡ Benchmarking ===\033[0m"
	@grep -h -E '^benchmark_.*:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-58s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;34m=== 🔨 Building ===\033[0m"
	@grep -h -E '^(build_.*|clean_builds):.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-58s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[1;34m=== 🧹 Code Quality ===\033[0m"
	@grep -h -E '^(fmt|lint|vet):.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-58s\033[0m %s\n", $$1, $$2}'
	@echo ""


################
### Testing  ###
################

.PHONY: test
test: ## Run all tests with verbose output
	go test -v -race ./...

.PHONY: test_all
test_all: ## Run all tests (legacy alias for compatibility)
	go test -v -p 1 ./... -mod=readonly -race

.PHONY: test_short
test_short: ## Run tests with short flag (skip long-running tests)
	go test -v -short ./...

.PHONY: test_coverage
test_coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: test_crypto
test_crypto: ## Run only crypto backend tests
	go test -v ./crypto/...

.PHONY: test_integration
test_integration: ## Run integration tests between backends
	@echo "🔄 Running cross-backend integration tests..."
	go test -v ./crypto/ -run="TestCompatibility"

################
### Linting  ###
################

.PHONY: fmt
fmt: ## Format Go code with gofmt
	@echo "🎨 Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatted"

.PHONY: lint
lint: ## Run golangci-lint (install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@echo "🔍 Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found. Install with:"; \
		echo "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...
	@echo "✅ Vet checks passed"

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "🧹 Tidying go modules..."
	go mod tidy
	@echo "✅ Modules tidied"

##################
### Validation ###
##################

.PHONY: validate
validate: fmt vet lint test ## Run all validation checks (format, vet, lint, test)
	@echo "✅ All validation checks passed!"

.PHONY: ci
ci: validate benchmark_compatibility ## Run CI pipeline (validation + compatibility tests)
	@echo "🎉 CI pipeline completed successfully!"

##################
### Quick Demo ###
##################

.PHONY: demo
demo: ## Quick demo showing ring signature with both backends
	@echo "🎭 Ring Signature Demo..."
	@echo "=================================================================="
	@echo "📱 Testing Decred backend (Pure Go):"
	go run ./examples/main.go
	@echo ""
	@echo "⚡ Testing Ethereum backend (libsecp256k1):"
	go run -tags=ethereum_secp256k1 ./examples/main.go 2>/dev/null || echo "⚠️  Ethereum backend requires CGO"
	@echo "=================================================================="

###########################
### Legacy Compatibility ###
###########################

.PHONY: tag_bug_fix
tag_bug_fix: ## Tag a new bug fix release (e.g., v1.0.1 -> v1.0.2)
	@$(eval LATEST_TAG=$(shell git tag --sort=-v:refname | head -n 1))
	@$(eval NEW_TAG=$(shell echo $(LATEST_TAG) | awk -F. -v OFS=. '{ $$NF = sprintf("%d", $$NF + 1); print }'))
	@git tag $(NEW_TAG)
	@echo "New bug fix version tagged: $(NEW_TAG)"
	@echo "Run the following commands to push the new tag:"
	@echo "  git push origin $(NEW_TAG)"
	@echo "And draft a new release at https://github.com/pokt-network/ring-go/releases/new"

.PHONY: tag_minor_release
tag_minor_release: ## Tag a new minor release (e.g. v1.0.0 -> v1.1.0)
	@$(eval LATEST_TAG=$(shell git tag --sort=-v:refname | head -n 1))
	@$(eval NEW_TAG=$(shell echo $(LATEST_TAG) | awk -F. '{$$2 += 1; $$3 = 0; print $$1 "." $$2 "." $$3}'))
	@git tag $(NEW_TAG)
	@echo "New minor release version tagged: $(NEW_TAG)"
	@echo "Run the following commands to push the new tag:"
	@echo "  git push origin $(NEW_TAG)"
	@echo "And draft a new release at https://github.com/pokt-network/ring-go/releases/new"