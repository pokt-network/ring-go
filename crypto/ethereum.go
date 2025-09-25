//go:build cgo && ethereum_secp256k1
// +build cgo,ethereum_secp256k1

package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/pokt-network/go-dleq/types"
)

var _ CurveBackend = (*ethereumBackend)(nil)

// ethereumBackend implements CurveBackend using Ethereum's libsecp256k1 wrapper.
// This provides better performance through CGO-optimized C library.
type ethereumBackend struct{}

// ethereumScalar implements types.Scalar using big.Int
type ethereumScalar struct {
	value *big.Int
}

// ethereumPoint implements types.Point using big.Int coordinates
type ethereumPoint struct {
	x, y *big.Int
}

// newSecp256k1Backend creates a new Ethereum-based secp256k1 backend.
func newSecp256k1Backend() CurveBackend {
	fmt.Println("RING-GO CRYPTO BACKEND: Using 'Ethereum' backend. CGO is enabled so this will be faster than 'Decred' backend.")
	return &ethereumBackend{}
}

// ScalarBaseMul multiplies the base point by the scalar
func (e *ethereumBackend) ScalarBaseMul(scalar types.Scalar) types.Point {
	es := scalar.(*ethereumScalar)
	// Use ethereum's crypto for scalar base multiplication
	x, y := ethcrypto.S256().ScalarBaseMult(es.value.Bytes())
	return &ethereumPoint{x: x, y: y}
}

// ScalarMul multiplies a point by a scalar
func (e *ethereumBackend) ScalarMul(scalar types.Scalar, point types.Point) types.Point {
	es := scalar.(*ethereumScalar)
	ep := point.(*ethereumPoint)
	// Use ethereum's crypto for scalar multiplication
	x, y := ethcrypto.S256().ScalarMult(ep.x, ep.y, es.value.Bytes())
	return &ethereumPoint{x: x, y: y}
}

// NewRandomScalar generates a new random scalar
func (e *ethereumBackend) NewRandomScalar() types.Scalar {
	// Generate random bytes
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(fmt.Sprintf("failed to generate random scalar: %v", err))
	}

	// Convert to big.Int and ensure it's within the curve order
	scalar := new(big.Int).SetBytes(bytes)
	scalar.Mod(scalar, ethcrypto.S256().Params().N)

	return &ethereumScalar{value: scalar}
}

// ScalarFromInt creates a scalar from a uint32
func (e *ethereumBackend) ScalarFromInt(i uint32) types.Scalar {
	return &ethereumScalar{value: big.NewInt(int64(i))}
}

// ScalarFromBytes creates a scalar from a byte array
func (e *ethereumBackend) ScalarFromBytes(bytes [32]byte) types.Scalar {
	scalar := new(big.Int).SetBytes(bytes[:])
	scalar.Mod(scalar, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: scalar}
}

// NewScalar creates a new scalar from a big.Int
func (e *ethereumBackend) NewScalar(value *big.Int) types.Scalar {
	v := new(big.Int).Set(value)
	v.Mod(v, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: v}
}

// PointFromCoordinates creates a point from x,y coordinates
func (e *ethereumBackend) PointFromCoordinates(x, y *big.Int) types.Point {
	// Verify the point is on the curve
	if !ethcrypto.S256().IsOnCurve(x, y) {
		panic("point is not on curve")
	}
	return &ethereumPoint{x: new(big.Int).Set(x), y: new(big.Int).Set(y)}
}

// PointFromBytes creates a point from compressed or uncompressed bytes
func (e *ethereumBackend) PointFromBytes(bytes []byte) (types.Point, error) {
	if len(bytes) == 33 {
		// Compressed point
		pubKey, err := ethcrypto.DecompressPubkey(bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress pubkey: %w", err)
		}
		// DecompressPubkey returns an ecdsa.PublicKey, extract x,y
		return &ethereumPoint{x: pubKey.X, y: pubKey.Y}, nil
	} else if len(bytes) == 65 {
		// Uncompressed point
		pubKey, err := ethcrypto.UnmarshalPubkey(bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal pubkey: %w", err)
		}
		return &ethereumPoint{x: pubKey.X, y: pubKey.Y}, nil
	}
	return nil, fmt.Errorf("invalid point bytes length: %d", len(bytes))
}

// BasePoint returns the secp256k1 base point
func (e *ethereumBackend) BasePoint() types.Point {
	params := ethcrypto.S256().Params()
	return &ethereumPoint{x: new(big.Int).Set(params.Gx), y: new(big.Int).Set(params.Gy)}
}

// AltBasePoint returns an alternative base point (for H parameter in ring signatures)
func (e *ethereumBackend) AltBasePoint() types.Point {
	// Use a deterministic point derived from the generator
	// This is a common pattern in ring signature implementations
	h := ethcrypto.Keccak256([]byte("alternative-base-point"))
	scalar := new(big.Int).SetBytes(h)
	scalar.Mod(scalar, ethcrypto.S256().Params().N)
	x, y := ethcrypto.S256().ScalarBaseMult(scalar.Bytes())
	return &ethereumPoint{x: x, y: y}
}

// HashToScalar hashes data to a scalar using Keccak256
func (e *ethereumBackend) HashToScalar(data []byte) (types.Scalar, error) {
	hash := ethcrypto.Keccak256(data)
	scalar := new(big.Int).SetBytes(hash)
	scalar.Mod(scalar, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: scalar}, nil
}

// DecodeToScalar decodes bytes to a scalar
func (e *ethereumBackend) DecodeToScalar(data []byte) (types.Scalar, error) {
	if len(data) != 32 {
		return nil, fmt.Errorf("invalid scalar length: %d", len(data))
	}
	scalar := new(big.Int).SetBytes(data)
	scalar.Mod(scalar, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: scalar}, nil
}

// DecodeToPoint decodes bytes to a point
func (e *ethereumBackend) DecodeToPoint(data []byte) (types.Point, error) {
	return e.PointFromBytes(data)
}

// BitSize returns the bit size of scalars (256 for secp256k1)
func (e *ethereumBackend) BitSize() uint64 {
	return 256
}

// CompressedPointSize returns the size of compressed points (33 bytes for secp256k1)
func (e *ethereumBackend) CompressedPointSize() int {
	return 33
}

// Sign creates a signature using the scalar and point
func (e *ethereumBackend) Sign(s types.Scalar, p types.Point) ([]byte, error) {
	// For ring signatures, this is typically implemented at a higher level
	// For now, we'll create a basic ECDSA signature
	scalar := s.(*ethereumScalar)

	// Create a private key from the scalar
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: ethcrypto.S256(),
		},
		D: scalar.value,
	}

	// Derive public key
	privKey.PublicKey.X, privKey.PublicKey.Y = ethcrypto.S256().ScalarBaseMult(scalar.value.Bytes())

	// Sign the point coordinates as message
	point := p.(*ethereumPoint)
	message := append(point.x.Bytes(), point.y.Bytes()...)

	return ethcrypto.Sign(message, privKey)
}

// Verify verifies a signature
func (e *ethereumBackend) Verify(pubkey, msgPoint types.Point, sig []byte) bool {
	// Basic ECDSA verification
	pubPoint := pubkey.(*ethereumPoint)
	msgPt := msgPoint.(*ethereumPoint)

	// Reconstruct message
	message := append(msgPt.x.Bytes(), msgPt.y.Bytes()...)

	// Create public key
	pubKey := &ecdsa.PublicKey{
		Curve: ethcrypto.S256(),
		X:     pubPoint.x,
		Y:     pubPoint.y,
	}

	return ethcrypto.VerifySignature(ethcrypto.CompressPubkey(pubKey), message, sig[:len(sig)-1])
}

// Name returns the backend name
func (e *ethereumBackend) Name() string {
	return "Ethereum (libsecp256k1)"
}

// Implement types.Scalar interface for ethereumScalar

// Add adds two scalars modulo the curve order
func (es *ethereumScalar) Add(other types.Scalar) types.Scalar {
	otherScalar := other.(*ethereumScalar)
	result := new(big.Int).Add(es.value, otherScalar.value)
	result.Mod(result, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: result}
}

// Sub subtracts another scalar from this one
func (es *ethereumScalar) Sub(other types.Scalar) types.Scalar {
	otherScalar := other.(*ethereumScalar)
	result := new(big.Int).Sub(es.value, otherScalar.value)
	result.Mod(result, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: result}
}

// Mul multiplies two scalars modulo the curve order
func (es *ethereumScalar) Mul(other types.Scalar) types.Scalar {
	otherScalar := other.(*ethereumScalar)
	result := new(big.Int).Mul(es.value, otherScalar.value)
	result.Mod(result, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: result}
}

// Negate returns the negation of the scalar
func (es *ethereumScalar) Negate() types.Scalar {
	result := new(big.Int).Neg(es.value)
	result.Mod(result, ethcrypto.S256().Params().N)
	return &ethereumScalar{value: result}
}

// Inverse returns the multiplicative inverse of the scalar
func (es *ethereumScalar) Inverse() types.Scalar {
	result := new(big.Int).ModInverse(es.value, ethcrypto.S256().Params().N)
	if result == nil {
		// Return zero if inverse doesn't exist
		return &ethereumScalar{value: new(big.Int)}
	}
	return &ethereumScalar{value: result}
}

// Eq checks if two scalars are equal
func (es *ethereumScalar) Eq(other types.Scalar) bool {
	otherScalar, ok := other.(*ethereumScalar)
	if !ok {
		return false
	}
	return es.value.Cmp(otherScalar.value) == 0
}

// IsZero returns true if the scalar is zero
func (es *ethereumScalar) IsZero() bool {
	return es.value.Sign() == 0
}

// Encode returns the scalar as bytes
func (es *ethereumScalar) Encode() []byte {
	bytes := es.value.Bytes()
	// Pad to 32 bytes
	result := make([]byte, 32)
	copy(result[32-len(bytes):], bytes)
	return result
}

// Implement types.Point interface for ethereumPoint

// Add adds two points on the curve
func (ep *ethereumPoint) Add(other types.Point) types.Point {
	otherPoint := other.(*ethereumPoint)
	x, y := ethcrypto.S256().Add(ep.x, ep.y, otherPoint.x, otherPoint.y)
	return &ethereumPoint{x: x, y: y}
}

// Sub subtracts another point from this one
func (ep *ethereumPoint) Sub(other types.Point) types.Point {
	otherPoint := other.(*ethereumPoint)
	// Negate the other point and add
	negY := new(big.Int).Neg(otherPoint.y)
	negY.Mod(negY, ethcrypto.S256().Params().P)
	x, y := ethcrypto.S256().Add(ep.x, ep.y, otherPoint.x, negY)
	return &ethereumPoint{x: x, y: y}
}

// ScalarMul multiplies the point by a scalar
func (ep *ethereumPoint) ScalarMul(scalar types.Scalar) types.Point {
	es := scalar.(*ethereumScalar)
	x, y := ethcrypto.S256().ScalarMult(ep.x, ep.y, es.value.Bytes())
	return &ethereumPoint{x: x, y: y}
}

// IsZero returns true if the point is the identity element
func (ep *ethereumPoint) IsZero() bool {
	return ep.x == nil || ep.y == nil || (ep.x.Sign() == 0 && ep.y.Sign() == 0)
}

// Equals checks if two points are equal
func (ep *ethereumPoint) Equals(other types.Point) bool {
	otherPoint, ok := other.(*ethereumPoint)
	if !ok {
		return false
	}
	if ep.IsZero() && otherPoint.IsZero() {
		return true
	}
	if ep.IsZero() || otherPoint.IsZero() {
		return false
	}
	return ep.x.Cmp(otherPoint.x) == 0 && ep.y.Cmp(otherPoint.y) == 0
}

// Copy creates a copy of the point
func (ep *ethereumPoint) Copy() types.Point {
	return &ethereumPoint{
		x: new(big.Int).Set(ep.x),
		y: new(big.Int).Set(ep.y),
	}
}

// Encode returns the point in compressed format
func (ep *ethereumPoint) Encode() []byte {
	pubkey := &ecdsa.PublicKey{
		Curve: ethcrypto.S256(),
		X:     ep.x,
		Y:     ep.y,
	}
	return ethcrypto.CompressPubkey(pubkey)
}
