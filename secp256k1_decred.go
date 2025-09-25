//go:build !ethereum_secp256k1
// +build !ethereum_secp256k1

package ring

import (
	dsecp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/athanorlabs/go-dleq/secp256k1"
)

// newPointFromFieldVals creates a secp256k1 point from Decred FieldVal coordinates
// For Decred backend, use FieldVal directly
func newPointFromFieldVals(fe, maybeY *dsecp256k1.FieldVal) *secp256k1.PointImpl {
	return secp256k1.NewPointFromCoordinates(*fe, *maybeY)
}
