//go:build !ethereum_secp256k1
// +build !ethereum_secp256k1

package crypto

import (
	"fmt"

	"github.com/athanorlabs/go-dleq/secp256k1"
	"github.com/athanorlabs/go-dleq/types"
)

// TODO_OPTIMIZE: Consider caching computed points to avoid regenerating them on each call
// decredBackend implements CurveBackend using Decred's pure Go secp256k1 implementation.
// This provides excellent performance without requiring CGO, making it highly portable.

var _ CurveBackend = (*decredBackend)(nil)

// decredBackend implements CurveBackend using the existing go-dleq Decred implementation.
type decredBackend struct {
	curve types.Curve
}

// newSecp256k1Backend creates a new Decred-based secp256k1 backend.
// This function is called by NewSecp256k1Backend when the ethereum_secp256k1 build tag is NOT active.
func newSecp256k1Backend() CurveBackend {
	backend := &decredBackend{
		curve: secp256k1.NewCurve(),
	}
	fmt.Println("RING-GO CRYPTO BACKEND: Using 'Decred' backend. CGO is disabled so this will be slower than 'Ethereum' backend.")
	return backend
}

// ScalarBaseMul implements CurveBackend.ScalarBaseMul using Decred's secp256k1 implementation.
func (b *decredBackend) ScalarBaseMul(scalar types.Scalar) types.Point {
	return b.curve.ScalarBaseMul(scalar)
}

// ScalarMul implements CurveBackend.ScalarMul using Decred's secp256k1 implementation.
func (b *decredBackend) ScalarMul(scalar types.Scalar, point types.Point) types.Point {
	return b.curve.ScalarMul(scalar, point)
}

// NewRandomScalar implements CurveBackend.NewRandomScalar using Decred's secp256k1 implementation.
func (b *decredBackend) NewRandomScalar() types.Scalar {
	return b.curve.NewRandomScalar()
}

// ScalarFromInt implements CurveBackend.ScalarFromInt using Decred's secp256k1 implementation.
func (b *decredBackend) ScalarFromInt(i uint32) types.Scalar {
	return b.curve.ScalarFromInt(i)
}

// ScalarFromBytes implements CurveBackend.ScalarFromBytes using Decred's secp256k1 implementation.
func (b *decredBackend) ScalarFromBytes(data [32]byte) types.Scalar {
	return b.curve.ScalarFromBytes(data)
}

// BasePoint implements CurveBackend.BasePoint using Decred's secp256k1 implementation.
func (b *decredBackend) BasePoint() types.Point {
	return b.curve.BasePoint()
}

// AltBasePoint implements CurveBackend.AltBasePoint using Decred's secp256k1 implementation.
func (b *decredBackend) AltBasePoint() types.Point {
	return b.curve.AltBasePoint()
}

// HashToScalar implements CurveBackend.HashToScalar using Decred's secp256k1 implementation.
func (b *decredBackend) HashToScalar(data []byte) (types.Scalar, error) {
	return b.curve.HashToScalar(data)
}

// DecodeToScalar implements CurveBackend.DecodeToScalar using Decred's secp256k1 implementation.
func (b *decredBackend) DecodeToScalar(data []byte) (types.Scalar, error) {
	return b.curve.DecodeToScalar(data)
}

// DecodeToPoint implements CurveBackend.DecodeToPoint using Decred's secp256k1 implementation.
func (b *decredBackend) DecodeToPoint(data []byte) (types.Point, error) {
	return b.curve.DecodeToPoint(data)
}

// BitSize implements CurveBackend.BitSize using Decred's secp256k1 implementation.
func (b *decredBackend) BitSize() uint64 {
	return b.curve.BitSize()
}

// CompressedPointSize implements CurveBackend.CompressedPointSize using Decred's secp256k1 implementation.
func (b *decredBackend) CompressedPointSize() int {
	return b.curve.CompressedPointSize()
}

// Sign implements CurveBackend.Sign using Decred's secp256k1 implementation.
func (b *decredBackend) Sign(s types.Scalar, p types.Point) ([]byte, error) {
	return b.curve.Sign(s, p)
}

// Verify implements CurveBackend.Verify using Decred's secp256k1 implementation.
func (b *decredBackend) Verify(pubkey, msgPoint types.Point, sig []byte) bool {
	return b.curve.Verify(pubkey, msgPoint, sig)
}

// Name implements CurveBackend.Name.
func (b *decredBackend) Name() string {
	return "Decred (Pure Go)"
}
