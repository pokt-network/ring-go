package ring

import (
	"github.com/pokt-network/go-dleq/ed25519"
	"github.com/pokt-network/go-dleq/secp256k1"
	"github.com/pokt-network/go-dleq/types"
)

type (
	// Curve represents an elliptic curve that can be used for signing.
	Curve = types.Curve
)

// Ed25519 returns a new ed25519 curve instance.
func Ed25519() types.Curve {
	return ed25519.NewCurve()
}

// Secp256k1 returns a new secp256k1 curve instance.
// BUILD-TIME CONFIGURATION: With ethereum_secp256k1 build tag,
// expensive operations are accelerated using libsecp256k1 via go-ethereum.
func Secp256k1() types.Curve {
	return secp256k1.NewCurve()
}
