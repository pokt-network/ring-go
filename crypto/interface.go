package crypto

import (
	"github.com/athanorlabs/go-dleq/types"
)

// CurveBackend defines the interface for pluggable secp256k1 implementations.
//
// BUILD-TIME CONFIGURATION: Different implementations are selected at compile time
// based on build tags for optimal performance vs portability trade-offs.
//
// Available backends:
// - Ethereum (build tag: ethereum_secp256k1): Uses libsecp256k1 C library, fastest performance, requires CGO
// - Decred (default, no build tag): Pure Go implementation, excellent performance, maximum portability
type CurveBackend interface {
	// Core curve operations
	ScalarBaseMul(scalar types.Scalar) types.Point
	ScalarMul(scalar types.Scalar, point types.Point) types.Point
	NewRandomScalar() types.Scalar
	ScalarFromInt(i uint32) types.Scalar
	ScalarFromBytes(b [32]byte) types.Scalar

	// Point operations
	BasePoint() types.Point
	AltBasePoint() types.Point

	// Encoding/decoding
	HashToScalar(data []byte) (types.Scalar, error)
	DecodeToScalar(data []byte) (types.Scalar, error)
	DecodeToPoint(data []byte) (types.Point, error)

	// Curve properties
	BitSize() uint64
	CompressedPointSize() int
	Sign(s types.Scalar, p types.Point) ([]byte, error)
	Verify(pubkey, msgPoint types.Point, sig []byte) bool
	Name() string
}

// NewSecp256k1Backend creates a new secp256k1 backend instance.
// BUILD-TIME CONFIGURATION: The actual implementation is determined at compile time
// based on build tags:
//
// - With "ethereum_secp256k1" tag: Uses Ethereum's libsecp256k1 (fastest, requires CGO)
// - Without tag: Uses Decred's implementation (portable, pure Go)
//
// Example usage:
//
//    backend := crypto.NewSecp256k1Backend()
//    curve := NewCurveFromBackend(backend)
func NewSecp256k1Backend() CurveBackend {
	return newSecp256k1Backend()
}

// CurveWrapper wraps a CurveBackend to implement the types.Curve interface
type CurveWrapper struct {
	backend CurveBackend
}

// NewCurveFromBackend creates a types.Curve from a CurveBackend
func NewCurveFromBackend(backend CurveBackend) types.Curve {
	return &CurveWrapper{backend: backend}
}

// Implement types.Curve interface
func (c *CurveWrapper) ScalarBaseMul(scalar types.Scalar) types.Point {
	return c.backend.ScalarBaseMul(scalar)
}

func (c *CurveWrapper) ScalarMul(scalar types.Scalar, point types.Point) types.Point {
	return c.backend.ScalarMul(scalar, point)
}

func (c *CurveWrapper) NewRandomScalar() types.Scalar {
	return c.backend.NewRandomScalar()
}

func (c *CurveWrapper) ScalarFromInt(i uint32) types.Scalar {
	return c.backend.ScalarFromInt(i)
}

func (c *CurveWrapper) ScalarFromBytes(b [32]byte) types.Scalar {
	return c.backend.ScalarFromBytes(b)
}

func (c *CurveWrapper) BasePoint() types.Point {
	return c.backend.BasePoint()
}

func (c *CurveWrapper) AltBasePoint() types.Point {
	return c.backend.AltBasePoint()
}

func (c *CurveWrapper) HashToScalar(data []byte) (types.Scalar, error) {
	return c.backend.HashToScalar(data)
}

func (c *CurveWrapper) DecodeToScalar(data []byte) (types.Scalar, error) {
	return c.backend.DecodeToScalar(data)
}

func (c *CurveWrapper) DecodeToPoint(data []byte) (types.Point, error) {
	return c.backend.DecodeToPoint(data)
}

func (c *CurveWrapper) BitSize() uint64 {
	return c.backend.BitSize()
}

func (c *CurveWrapper) CompressedPointSize() int {
	return c.backend.CompressedPointSize()
}

func (c *CurveWrapper) Sign(s types.Scalar, p types.Point) ([]byte, error) {
	return c.backend.Sign(s, p)
}

func (c *CurveWrapper) Verify(pubkey, msgPoint types.Point, sig []byte) bool {
	return c.backend.Verify(pubkey, msgPoint, sig)
}

func (c *CurveWrapper) Name() string {
	return c.backend.Name()
}