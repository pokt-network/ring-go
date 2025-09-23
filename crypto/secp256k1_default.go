//go:build !ethereum_secp256k1
// +build !ethereum_secp256k1

package crypto

import (
	"github.com/athanorlabs/go-dleq/secp256k1"
	"github.com/athanorlabs/go-dleq/types"
)

// NewOptimizedSecp256k1Curve creates a secp256k1 curve using the default backend
// When CGO/Ethereum backend is not available, this just returns the standard go-dleq curve
func NewOptimizedSecp256k1Curve() types.Curve {
	return secp256k1.NewCurve()
}