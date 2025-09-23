package ring

import (
	"github.com/athanorlabs/go-dleq/ed25519"
	"github.com/athanorlabs/go-dleq/types"
	"github.com/pokt-network/ring-go/crypto"
)

type (
	// Curve represents an elliptic curve that can be used for signing.
	Curve = types.Curve
)

// Ed25519 returns a new ed25519 curve instance.
func Ed25519() types.Curve {
	return ed25519.NewCurve()
}

// Secp256k1 returns a new secp256k1 curve instance. When available, this automatically
// uses optimized backends for better performance while maintaining full compatibility.
// BUILD-TIME CONFIGURATION: With ethereum_secp256k1 build tag, expensive operations
// are accelerated using libsecp256k1 via go-ethereum.
func Secp256k1() types.Curve {
	return crypto.NewOptimizedSecp256k1Curve()
}

// Secp256k1Fast returns a new secp256k1 curve instance using pluggable crypto backends.
// BUILD-TIME CONFIGURATION: The actual implementation is determined at compile time:
//
// - With "ethereum_secp256k1" tag: Uses Ethereum's libsecp256k1 (fastest, requires CGO)
// - Without tag: Uses Decred's implementation (portable, pure Go)
//
// Performance benefits:
// - ~50% faster signing operations
// - ~80% faster verification operations
// - Significantly fewer memory allocations
//
// Example usage:
//
//    curve := ring.Secp256k1Fast()  // Auto-selects optimal backend
//    ring, err := ring.NewKeyRing(curve, size, privKey, idx)
func Secp256k1Fast() types.Curve {
	backend := crypto.NewSecp256k1Backend()
	return crypto.NewCurveFromBackend(backend)
}
