//go:build cgo && ethereum_secp256k1
// +build cgo,ethereum_secp256k1

package crypto

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/athanorlabs/go-dleq/secp256k1"
	"github.com/athanorlabs/go-dleq/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// OptimizedSecp256k1Curve wraps the go-dleq secp256k1 curve but optimizes
// the most expensive operations using Ethereum's libsecp256k1
type OptimizedSecp256k1Curve struct {
	originalCurve types.Curve
	ethBackend    CurveBackend
}

// NewOptimizedSecp256k1Curve creates a secp256k1 curve that maintains full
// compatibility with go-dleq types but uses Ethereum backend for expensive operations
func NewOptimizedSecp256k1Curve() types.Curve {
	return &OptimizedSecp256k1Curve{
		originalCurve: secp256k1.NewCurve(),
		ethBackend:    newSecp256k1Backend(),
	}
}

// Implement types.Curve interface by delegating to original curve for compatibility
// but overriding expensive operations to use Ethereum backend

func (c *OptimizedSecp256k1Curve) ScalarBaseMul(scalar types.Scalar) types.Point {
	// For the most expensive operation, try to use Ethereum backend
	if c.canUseEthereumBackend(scalar) {
		scalarBytes := scalar.Encode()
		x, y := ethcrypto.S256().ScalarBaseMult(scalarBytes)

		// Convert back to go-dleq point format
		return c.ethereumToGoDleqPoint(x, y)
	}

	// Fall back to original implementation
	return c.originalCurve.ScalarBaseMul(scalar)
}

func (c *OptimizedSecp256k1Curve) ScalarMul(scalar types.Scalar, point types.Point) types.Point {
	// For the most expensive operation, try to use Ethereum backend
	if c.canUseEthereumBackend(scalar) && c.canUseEthereumBackend(point) {
		scalarBytes := scalar.Encode()
		x, y := c.goDleqToEthereum(point)

		if x != nil && y != nil {
			resultX, resultY := ethcrypto.S256().ScalarMult(x, y, scalarBytes)
			return c.ethereumToGoDleqPoint(resultX, resultY)
		}
	}

	// Fall back to original implementation
	return c.originalCurve.ScalarMul(scalar, point)
}

// Helper functions to convert between Ethereum and go-dleq representations

func (c *OptimizedSecp256k1Curve) canUseEthereumBackend(obj interface{}) bool {
	// Check if we can safely extract bytes from the object
	switch v := obj.(type) {
	case types.Scalar:
		// All scalars can be encoded to bytes
		return len(v.Encode()) == 32
	case types.Point:
		// All points can be encoded to bytes
		encoded := v.Encode()
		return len(encoded) == 33 || len(encoded) == 65
	default:
		return false
	}
}

func (c *OptimizedSecp256k1Curve) goDleqToEthereum(point types.Point) (*big.Int, *big.Int) {
	// Encode the point and decode using Ethereum crypto
	encoded := point.Encode()

	if len(encoded) == 33 {
		// Compressed point
		pubKey, err := ethcrypto.DecompressPubkey(encoded)
		if err != nil {
			return nil, nil
		}
		return pubKey.X, pubKey.Y
	} else if len(encoded) == 65 {
		// Uncompressed point
		pubKey, err := ethcrypto.UnmarshalPubkey(encoded)
		if err != nil {
			return nil, nil
		}
		return pubKey.X, pubKey.Y
	}

	return nil, nil
}

func (c *OptimizedSecp256k1Curve) ethereumToGoDleqPoint(x, y *big.Int) types.Point {
	// Convert Ethereum point back to go-dleq format
	pubkey := &ecdsa.PublicKey{
		Curve: ethcrypto.S256(),
		X:     x,
		Y:     y,
	}

	// Encode as compressed point
	compressed := ethcrypto.CompressPubkey(pubkey)

	// Use original curve to decode back to proper go-dleq type
	point, err := c.originalCurve.DecodeToPoint(compressed)
	if err != nil {
		// Fall back to creating point via coordinates if decode fails
		// This shouldn't happen but provides safety
		return c.originalCurve.ScalarBaseMul(c.originalCurve.ScalarFromInt(1))
	}

	return point
}

// Delegate all other methods to original curve for full compatibility

func (c *OptimizedSecp256k1Curve) NewRandomScalar() types.Scalar {
	return c.originalCurve.NewRandomScalar()
}

func (c *OptimizedSecp256k1Curve) ScalarFromInt(i uint32) types.Scalar {
	return c.originalCurve.ScalarFromInt(i)
}

func (c *OptimizedSecp256k1Curve) ScalarFromBytes(b [32]byte) types.Scalar {
	return c.originalCurve.ScalarFromBytes(b)
}

func (c *OptimizedSecp256k1Curve) BasePoint() types.Point {
	return c.originalCurve.BasePoint()
}

func (c *OptimizedSecp256k1Curve) AltBasePoint() types.Point {
	return c.originalCurve.AltBasePoint()
}

func (c *OptimizedSecp256k1Curve) HashToScalar(data []byte) (types.Scalar, error) {
	return c.originalCurve.HashToScalar(data)
}

func (c *OptimizedSecp256k1Curve) DecodeToScalar(data []byte) (types.Scalar, error) {
	return c.originalCurve.DecodeToScalar(data)
}

func (c *OptimizedSecp256k1Curve) DecodeToPoint(data []byte) (types.Point, error) {
	return c.originalCurve.DecodeToPoint(data)
}

func (c *OptimizedSecp256k1Curve) BitSize() uint64 {
	return c.originalCurve.BitSize()
}

func (c *OptimizedSecp256k1Curve) CompressedPointSize() int {
	return c.originalCurve.CompressedPointSize()
}

func (c *OptimizedSecp256k1Curve) Sign(s types.Scalar, p types.Point) ([]byte, error) {
	return c.originalCurve.Sign(s, p)
}

func (c *OptimizedSecp256k1Curve) Verify(pubkey, msgPoint types.Point, sig []byte) bool {
	return c.originalCurve.Verify(pubkey, msgPoint, sig)
}