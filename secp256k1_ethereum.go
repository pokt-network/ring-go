//go:build cgo && ethereum_secp256k1
// +build cgo,ethereum_secp256k1

package ring

import (
	"math/big"

	"github.com/athanorlabs/go-dleq/secp256k1"
	dsecp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// newPointFromFieldVals creates a secp256k1 point from Decred FieldVal coordinates
// For Ethereum backend, convert FieldVal to *big.Int
func newPointFromFieldVals(fe, maybeY *dsecp256k1.FieldVal) *secp256k1.PointImpl {
	// Convert FieldVal to *big.Int for Ethereum backend API
	// Need to normalize first for proper byte representation
	fe.Normalize()
	maybeY.Normalize()
	x := new(big.Int).SetBytes(fe.Bytes()[:])
	y := new(big.Int).SetBytes(maybeY.Bytes()[:])
	return secp256k1.NewPointFromCoordinates(x, y)
}