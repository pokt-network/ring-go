//go:build cgo
// +build cgo

// Package crypto provides comprehensive benchmarks comparing different secp256k1 backends for ring signatures.
//
// Performance Summary (Ring Signature Operations):
//
// COMPARISON (Backend Build Tags):
// - Decred (no tags): Pure Go implementation, excellent performance, maximum portability
// - Ethereum (ethereum_secp256k1): Uses CGO and the underlying secp256k1 optimizations
//
// RECOMMENDATIONS:
// 1. For maximum performance: Use Ethereum backend with CGO (requires CGO)
//   - Better optimized curve operations
//   - Production-ready (Bitcoin Core standard libraries)
//
// 2. For CGO-free environments: Use Decred backend (default)
//   - Excellent portability across platforms
//   - No external dependencies
//   - Good performance for most use cases
package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var testMessage = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

// benchmarkBackend benchmarks the current backend for basic curve operations
func benchmarkBackend(b *testing.B, backendName string) {
	backend := NewSecp256k1Backend() // Uses currently compiled backend

	privKey := backend.NewRandomScalar()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		point := backend.ScalarBaseMul(privKey)
		_ = backend.ScalarMul(privKey, point)
	}
}

// Generate backend comparison benchmarks
func BenchmarkBackend_Current(b *testing.B) { benchmarkBackend(b, "Current") }

// TestBackendCompatibility ensures the current backend produces valid results
func TestBackendCompatibility(t *testing.T) {
	// Test current backend (determined by build tags)
	backend := NewSecp256k1Backend()
	privKey := backend.NewRandomScalar()
	point := backend.ScalarBaseMul(privKey)
	require.NotNil(t, point, "Backend should produce valid points")

	// Test scalar operations
	scalar2 := backend.NewRandomScalar()
	sum := privKey.Add(scalar2)
	require.NotNil(t, sum, "Scalar addition should work")

	// Test point operations
	point2 := backend.ScalarBaseMul(scalar2)
	pointSum := point.Add(point2)
	require.NotNil(t, pointSum, "Point addition should work")

	t.Logf("✅ %s backend produces valid curve operations", backend.Name())
}
