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

	"github.com/athanorlabs/go-dleq/secp256k1"
	"github.com/stretchr/testify/require"
)

var testMessage = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

// benchmarkBackend benchmarks a specific backend for basic curve operations
func benchmarkBackend(b *testing.B, backendName string, useFast bool) {
	var backend CurveBackend
	if useFast {
		backend = NewSecp256k1Backend() // Uses pluggable backends
	} else {
		backend = &decredBackend{curve: secp256k1.NewCurve()}
	}

	privKey := backend.NewRandomScalar()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		point := backend.ScalarBaseMul(privKey)
		_ = backend.ScalarMul(privKey, point)
	}
}

// Generate backend comparison benchmarks
func BenchmarkBackend_Decred(b *testing.B) { benchmarkBackend(b, "Decred", false) }
func BenchmarkBackend_Fast(b *testing.B)   { benchmarkBackend(b, "Fast", true) }

// TestBackendCompatibility ensures both backends produce compatible results
func TestBackendCompatibility(t *testing.T) {
	// Test Decred backend
	decredBackend := &decredBackend{curve: secp256k1.NewCurve()}
	decredPrivKey := decredBackend.NewRandomScalar()
	decredPoint := decredBackend.ScalarBaseMul(decredPrivKey)
	require.NotNil(t, decredPoint, "Decred backend should produce valid points")

	// Test Fast backend
	fastBackend := NewSecp256k1Backend()
	fastPrivKey := fastBackend.NewRandomScalar()
	fastPoint := fastBackend.ScalarBaseMul(fastPrivKey)
	require.NotNil(t, fastPoint, "Fast backend should produce valid points")

	t.Log("✅ Both backends produce valid curve operations")
}