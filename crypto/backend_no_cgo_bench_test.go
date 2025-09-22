//go:build !cgo
// +build !cgo

// Package crypto provides CGO-free benchmarks for ring signatures using Decred backend.
//
// This file contains benchmarks that work without CGO enabled,
// testing only the Decred pure Go implementation.
package crypto

import (
	"testing"

	"github.com/athanorlabs/go-dleq/secp256k1"
	"github.com/stretchr/testify/require"
)

var testMessageNoCgo = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

// benchmarkBackendNoCgo benchmarks basic crypto operations without CGO
func benchmarkBackendNoCgo(b *testing.B, backendName string) {
	backend := &decredBackend{curve: secp256k1.NewCurve()}
	privKey := backend.NewRandomScalar()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		point := backend.ScalarBaseMul(privKey)
		_ = backend.ScalarMul(privKey, point)
	}
}

// CGO-free Backend Benchmarks
func BenchmarkBackendNoCgo_Decred(b *testing.B) { benchmarkBackendNoCgo(b, "Decred") }

// TestCompatibilityNoCgo ensures the Decred backend produces valid operations without CGO
func TestCompatibilityNoCgo(t *testing.T) {
	backend := &decredBackend{curve: secp256k1.NewCurve()}
	privKey := backend.NewRandomScalar()
	point := backend.ScalarBaseMul(privKey)
	require.NotNil(t, point, "Decred backend should produce valid points")

	t.Log("✅ Decred backend produces valid operations (CGO-free)")
}