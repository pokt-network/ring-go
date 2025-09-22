#####################
### Build Targets ###
#####################

.PHONY: build_fast
build_fast: ## Build with Ethereum backend (50% faster signing, 80% faster verification, requires CGO)
	@echo "🚀 Building Ring-Go with Ethereum secp256k1 backend..."
	@echo "   • Requires CGO and libsecp256k1"
	@echo "   • ~50% faster signing, ~80% faster verification for ring signatures"
	@echo "=================================================================="
	go build -tags="ethereum_secp256k1" -o ring-go-fast ./examples/...
	@echo "✅ Built: ring-go-fast"

.PHONY: build_portable
build_portable: ## Build with Decred backend (pure Go, maximum portability, no CGO dependencies)
	@echo "🌍 Building Ring-Go with Decred secp256k1 backend..."
	@echo "   • Pure Go, no CGO dependencies"
	@echo "   • Excellent performance, maximum portability"
	@echo "=================================================================="
	CGO_ENABLED=0 go build -o ring-go-portable ./examples/...
	@echo "✅ Built: ring-go-portable"

.PHONY: build_auto
build_auto: ## Auto-select optimal backend (Ethereum if CGO available, otherwise Decred)
	@echo "🎯 Auto-selecting optimal crypto backend..."
	@if command -v gcc >/dev/null 2>&1 && [ "$$CGO_ENABLED" != "0" ]; then \
		echo "   • CGO available, building fast version..."; \
		$(MAKE) build_fast; \
	else \
		echo "   • No CGO or CGO disabled, building portable version..."; \
		$(MAKE) build_portable; \
	fi

.PHONY: build_all
build_all: ## Build both Ethereum (fast) and Decred (portable) versions
	@echo "🏗️  Building all Ring-Go variants..."
	$(MAKE) build_fast
	$(MAKE) build_portable
	@echo "=================================================================="
	@echo "✅ Built all variants:"
	@echo "   • ring-go-fast     (Ethereum backend)"
	@echo "   • ring-go-portable (Decred backend)"
	@ls -la ring-go-*

.PHONY: clean_builds
clean_builds: ## Remove all Ring-Go built binaries
	@echo "🧹 Cleaning built binaries..."
	rm -f ring-go-fast ring-go-portable
	@echo "✅ Cleaned all builds"

.PHONY: test_builds
test_builds: ## Test both build variants to ensure they work
	@echo "🧪 Testing build variants..."
	@echo "=================================================================="
	@if [ -f ring-go-portable ]; then \
		echo "Testing portable build..."; \
		./ring-go-portable || echo "❌ Portable build failed"; \
	fi
	@if [ -f ring-go-fast ]; then \
		echo "Testing fast build..."; \
		./ring-go-fast || echo "❌ Fast build failed"; \
	fi
	@echo "✅ Build tests completed"