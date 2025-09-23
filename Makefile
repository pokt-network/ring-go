# Ring-Go Makefile
# Simplified targets for core functionality

.PHONY: help
help: ## Prints all the targets in all the Makefiles
	@echo "📋 Ring-Go Makefile Targets"
	@echo ""
	@echo "\033[1m=== 🔍 Information & Discovery ===\033[0m"
	@echo "  \033[36mhelp\033[0m                                                       Prints all the targets in all the Makefiles"
	@echo ""
	@echo "\033[1m=== 🧪 Testing ===\033[0m"
	@echo "  \033[36mtest_all\033[0m                                                   Run all tests"
	@echo ""
	@echo "\033[1m=== ⚡ Benchmarking ===\033[0m"
	@echo "  \033[36mbenchmark_all\033[0m                                              Run all benchmarks (tests both Decred and Ethereum backends)"
	@echo "  \033[36mbenchmark_report\033[0m                                           Compare crypto backends with formatted report for ring signatures"
	@echo ""
	@echo "\033[1m=== 🔨 Building ===\033[0m"
	@echo "  \033[36mbuild_fast\033[0m                                                 Build with Ethereum backend (50% faster signing, 80% faster verification, requires CGO)"
	@echo "  \033[36mbuild_portable\033[0m                                             Build with Decred backend (pure Go, maximum portability, no CGO dependencies)"
	@echo "  \033[36mclean_builds\033[0m                                               Remove all Ring-Go built binaries"
	@echo ""
	@echo "\033[1m=== 🧹 Code Quality ===\033[0m"
	@echo "  \033[36mgo_fmt_and_lint\033[0m                                            Format Go code with gofmt"

####################
### Testing 🧪 ###
####################

.PHONY: test_all
test_all: ## Run all tests
	@echo "🧪 Running all tests..."
	go test -v ./...

####################
### Benchmarking ###
####################

.PHONY: benchmark_all
benchmark_all: ## Run all benchmarks (tests both Decred and Ethereum backends)
	@echo "🔬 Running benchmarks with Decred backend (Pure Go)..."
	@echo "=================================================="
	go test -v -bench=. -benchmem -run=^$$ ./...
	@echo ""
	@echo "🔬 Running benchmarks with Ethereum backend (CGO + libsecp256k1)..."
	@echo "=================================================================="
	go test -tags=ethereum_secp256k1 -v -bench=. -benchmem -run=^$$ ./...

.PHONY: benchmark_report
benchmark_report: ## Compare crypto backends with formatted report for ring signatures
	@echo "🔬 Benchmarking ring signature crypto backends..."
	@echo "=================================================================="
	@( \
		echo "# Decred Backend Results:" && \
		timeout 45s go test . \
			-bench="BenchmarkSign(2|4|8|16|32)_Secp256k1|BenchmarkVerify(2|4|8|16|32)_Secp256k1" \
			-benchmem \
			-run=^$$ \
			-benchtime=2s \
			2>/dev/null | sed 's/_Secp256k1/_Decred/g' && \
		echo "# Ethereum Backend Results:" && \
		timeout 45s go test -tags=ethereum_secp256k1 . \
			-bench="BenchmarkSign(2|4|8|16|32)_Secp256k1|BenchmarkVerify(2|4|8|16|32)_Secp256k1" \
			-benchmem \
			-run=^$$ \
			-benchtime=2s \
			2>/dev/null | sed 's/_Secp256k1/_Ethereum/g' \
	) | python3 format_benchmark.py || ( \
		echo "⚠️  Benchmark timed out or failed. Trying simpler benchmark..." && \
		( \
			echo "# Decred Backend Results:" && \
			go test . \
				-bench="BenchmarkSign2_|BenchmarkSign32_|BenchmarkVerify2_|BenchmarkVerify32_" \
				-benchmem \
				-run=^$$ \
				-benchtime=1s \
				2>/dev/null | sed 's/_Secp256k1/_Decred/g' && \
			echo "# Ethereum Backend Results:" && \
			go test -tags=ethereum_secp256k1 . \
				-bench="BenchmarkSign2_|BenchmarkSign32_|BenchmarkVerify2_|BenchmarkVerify32_" \
				-benchmem \
				-run=^$$ \
				-benchtime=1s \
				2>/dev/null | sed 's/_Secp256k1/_Ethereum/g' \
		) | python3 format_benchmark.py \
	)
	@echo "=================================================================="

####################
### Building 🔨 ###
####################

.PHONY: build_fast
build_fast: ## Build with Ethereum backend (50% faster signing, 80% faster verification, requires CGO)
	@echo "🔨 Building with Ethereum backend (high performance)..."
	CGO_ENABLED=1 go build -tags=ethereum_secp256k1 -o ring-go-fast .
	@echo "✅ Built: ring-go-fast (Ethereum backend - requires CGO)"

.PHONY: build_portable
build_portable: ## Build with Decred backend (pure Go, maximum portability, no CGO dependencies)
	@echo "🔨 Building with Decred backend (maximum portability)..."
	CGO_ENABLED=0 go build -o ring-go-portable .
	@echo "✅ Built: ring-go-portable (Decred backend - pure Go)"

.PHONY: clean_builds
clean_builds: ## Remove all Ring-Go built binaries
	@echo "🧹 Cleaning built binaries..."
	rm -f ring-go-fast ring-go-portable
	@echo "✅ Cleaned all binaries"

####################
### Code Quality ###
####################

.PHONY: go_fmt_and_lint
go_fmt_and_lint: ## Format Go code with gofmt
	@echo "🧹 Formatting Go code..."
	gofmt -l -w .
	@echo "✅ Go code formatted"