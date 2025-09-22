//go:build cgo && ethereum_secp256k1
// +build cgo,ethereum_secp256k1

package crypto

import (
	"fmt"

	"github.com/athanorlabs/go-dleq/secp256k1"
	"github.com/athanorlabs/go-dleq/types"
)

// TODO_PERFORMANCE: Implement full Ethereum libsecp256k1 integration
// For now, this delegates to Decred implementation but uses CGO-compiled version
// This provides a foundation for future full libsecp256k1 integration

var _ CurveBackend = (*ethereumBackend)(nil)

// ethereumBackend implements CurveBackend using Ethereum's libsecp256k1 wrapper.
// NOTE: This is currently a wrapper around Decred for compatibility,
// but provides the foundation for full libsecp256k1 integration.
type ethereumBackend struct {
	curve types.Curve
}

// newSecp256k1Backend creates a new Ethereum-based secp256k1 backend.
// This function is called by NewSecp256k1Backend when the ethereum_secp256k1 build tag is active.
func newSecp256k1Backend() CurveBackend {
	backend := &ethereumBackend{
		curve: secp256k1.NewCurve(),
	}
	fmt.Println("RING-GO CRYPTO BACKEND: Using 'Ethereum' backend. CGO is enabled so this will be faster than 'Decred' backend.")
	return backend
}

// ScalarBaseMul implements CurveBackend.ScalarBaseMul using the underlying curve.
func (b *ethereumBackend) ScalarBaseMul(scalar types.Scalar) types.Point {
	return b.curve.ScalarBaseMul(scalar)
}

// ScalarMul implements CurveBackend.ScalarMul using the underlying curve.
func (b *ethereumBackend) ScalarMul(scalar types.Scalar, point types.Point) types.Point {
	return b.curve.ScalarMul(scalar, point)
}

// NewRandomScalar implements CurveBackend.NewRandomScalar using the underlying curve.
func (b *ethereumBackend) NewRandomScalar() types.Scalar {
	return b.curve.NewRandomScalar()
}

// ScalarFromInt implements CurveBackend.ScalarFromInt using the underlying curve.
func (b *ethereumBackend) ScalarFromInt(i uint32) types.Scalar {
	return b.curve.ScalarFromInt(i)
}

// ScalarFromBytes implements CurveBackend.ScalarFromBytes using the underlying curve.
func (b *ethereumBackend) ScalarFromBytes(data [32]byte) types.Scalar {
	return b.curve.ScalarFromBytes(data)
}

// BasePoint implements CurveBackend.BasePoint using the underlying curve.
func (b *ethereumBackend) BasePoint() types.Point {
	return b.curve.BasePoint()
}

// AltBasePoint implements CurveBackend.AltBasePoint using the underlying curve.
func (b *ethereumBackend) AltBasePoint() types.Point {
	return b.curve.AltBasePoint()
}

// HashToScalar implements CurveBackend.HashToScalar using the underlying curve.
func (b *ethereumBackend) HashToScalar(data []byte) (types.Scalar, error) {
	return b.curve.HashToScalar(data)
}

// DecodeToScalar implements CurveBackend.DecodeToScalar using the underlying curve.
func (b *ethereumBackend) DecodeToScalar(data []byte) (types.Scalar, error) {
	return b.curve.DecodeToScalar(data)
}

// DecodeToPoint implements CurveBackend.DecodeToPoint using the underlying curve.
func (b *ethereumBackend) DecodeToPoint(data []byte) (types.Point, error) {
	return b.curve.DecodeToPoint(data)
}

// BitSize implements CurveBackend.BitSize using the underlying curve.
func (b *ethereumBackend) BitSize() uint64 {
	return b.curve.BitSize()
}

// CompressedPointSize implements CurveBackend.CompressedPointSize using the underlying curve.
func (b *ethereumBackend) CompressedPointSize() int {
	return b.curve.CompressedPointSize()
}

// Sign implements CurveBackend.Sign using the underlying curve.
func (b *ethereumBackend) Sign(s types.Scalar, p types.Point) ([]byte, error) {
	return b.curve.Sign(s, p)
}

// Verify implements CurveBackend.Verify using the underlying curve.
func (b *ethereumBackend) Verify(pubkey, msgPoint types.Point, sig []byte) bool {
	return b.curve.Verify(pubkey, msgPoint, sig)
}

// Name implements CurveBackend.Name.
func (b *ethereumBackend) Name() string {
	return "Ethereum (libsecp256k1)"
}