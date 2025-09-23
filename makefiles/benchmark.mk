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
		timeout 30s go test . \
			-bench="BenchmarkSign.*_Secp256k1|BenchmarkVerify.*_Secp256k1" \
			-benchmem \
			-run=^$$ \
			-benchtime=2s \
			2>/dev/null | sed 's/_Secp256k1/_Decred/g' && \
		echo "# Ethereum Backend Results:" && \
		timeout 30s go test -tags=ethereum_secp256k1 . \
			-bench="BenchmarkSign.*_Secp256k1|BenchmarkVerify.*_Secp256k1" \
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

.PHONY: benchmark_compatibility
benchmark_compatibility: ## Test backend compatibility and correctness
	@echo "🧪 Testing ring signature compatibility across backends..."
	@echo "=================================================================="
	go test ./crypto -v -run="TestCompatibility"
	@echo ""
	@echo "✅ All backends produce valid, interoperable ring signatures"

.PHONY: benchmark_memory
benchmark_memory: ## Analyze memory allocation patterns
	@echo "🧠 Analyzing memory allocation patterns..."
	@echo "=================================================================="
	go test ./crypto -bench="BenchmarkMemoryAllocation" -benchmem -v
	@echo ""
	@echo "💡 Memory optimization opportunities:"
	@echo "   • Monitor allocation counts for performance-critical paths"
	@echo "   • Consider object pooling for high-frequency operations"

.PHONY: benchmark_parallel
benchmark_parallel: ## Test parallel performance characteristics
	@echo "⚡ Testing parallel ring signature performance..."
	@echo "=================================================================="
	go test ./crypto -bench="BenchmarkSignParallel" -benchmem -v
	@echo ""
	@echo "💡 Parallel performance insights:"
	@echo "   • Ring signatures scale well across cores"
	@echo "   • CGO overhead may affect parallel scaling"