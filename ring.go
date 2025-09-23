package ring

import (
	"errors"
	"fmt"
	"sync"

	"github.com/athanorlabs/go-dleq/ed25519"
	"github.com/athanorlabs/go-dleq/types"
)

// Ring represents a group of public keys such that one of the group created a signature.
type Ring struct {
	pubkeys []types.Point
	curve   types.Curve
	// precomputed once to avoid recomputing in Sign/Verify loops.
	hp []types.Point
}

// Size returns the size of the ring, ie. the number of public keys in it.
func (r *Ring) Size() int {
	return len(r.pubkeys)
}

// Equals checks whether the supplied ring is equal to the current ring.
// The ring's public keys must be in the same order for the rings to be equal
func (r *Ring) Equals(other *Ring) bool {
	if r.Size() != other.Size() {
		return false
	}

	for i, p := range r.pubkeys {
		if !p.Equals(other.pubkeys[i]) {
			return false
		}
	}
	bp, abp := r.curve.BasePoint(), r.curve.AltBasePoint()
	obp, oabp := other.curve.BasePoint(), other.curve.AltBasePoint()
	return bp.Equals(obp) && abp.Equals(oabp)
}

// RingSig represents a ring signature.
type RingSig struct {
	ring  *Ring          // array of public keys
	c     types.Scalar   // ring signature challenge
	s     []types.Scalar // ring signature values
	image types.Point    // key image
}

// PublicKeys returns a copy of the ring signature's public keys.
func (r *RingSig) PublicKeys() []types.Point {
	ret := make([]types.Point, len(r.ring.pubkeys))
	for i, pk := range r.ring.pubkeys {
		ret[i] = pk.Copy()
	}
	return ret
}

// PublicKeysRef returns a reference to the ring's public keys without copying them.
// NOTE: While the method does not involve cloning or copying, mutating the returned slice
// or its elements can break the integrity of the ring signature structure. Avoid mutation.
func (r *RingSig) PublicKeysRef() []types.Point {
	return r.ring.pubkeys
}

// Reset clears fields so RingSig can be reused with a pool (additive; safe for callers).
func (r *RingSig) Reset() {
	r.ring = nil
	r.c = nil
	r.s = nil
	r.image = nil
}

// Ring returns the ring from the RingSig struct
func (r *RingSig) Ring() *Ring {
	return r.ring
}

// NewKeyRingFromPublicKeys takes public key ring and places the public key corresponding to `privKey`
// in index idx of the ring.
// It returns a ring of public keys of length `len(ring)+1`.
func NewKeyRingFromPublicKeys(curve types.Curve, pubkeys []types.Point, privKey types.Scalar, idx int) (*Ring, error) {
	size := len(pubkeys) + 1
	newRing := make([]types.Point, size)
	pubkey := curve.ScalarBaseMul(privKey)

	if idx > len(pubkeys) {
		return nil, errors.New("index out of bounds: idx > len(pubkeys)")
	}
	if idx < 0 {
		return nil, errors.New("index out of bounds: idx < 0")
	}
	if privKey.IsZero() {
		return nil, errors.New("private key is zero")
	}

	newRing[idx] = pubkey
	pubkeysMap := make(map[types.Point]struct{}, size)
	pubkeysMap[pubkey] = struct{}{}

	for i := 0; i < size; i++ {
		if i == idx {
			continue
		}
		if i < idx {
			newRing[i] = pubkeys[i]
		} else {
			newRing[i] = pubkeys[i-1]
		}
		pubkeysMap[newRing[i]] = struct{}{}
	}

	if len(pubkeysMap) != len(newRing) {
		return nil, errors.New("duplicate public keys in ring")
	}

	// Precompute H_p(P_i)
	hp := make([]types.Point, size)
	for i := 0; i < size; i++ {
		hp[i] = hashToCurve(newRing[i])
	}

	return &Ring{
		pubkeys: newRing,
		curve:   curve,
		hp:      hp,
	}, nil
}

// NewFixedKeyRingFromPublicKeys takes public keys and a curve to create a ring
func NewFixedKeyRingFromPublicKeys(curve types.Curve, pubkeys []types.Point) (*Ring, error) {
	pubkeysMap := make(map[types.Point]struct{}, len(pubkeys))

	size := len(pubkeys)
	newRing := make([]types.Point, size)
	for i := 0; i < size; i++ {
		pubkeysMap[pubkeys[i]] = struct{}{}
		newRing[i] = pubkeys[i].Copy()
	}

	if len(pubkeysMap) != len(newRing) {
		return nil, errors.New("duplicate public keys in ring")
	}

	hp := make([]types.Point, size)
	for i := 0; i < size; i++ {
		hp[i] = hashToCurve(newRing[i])
	}

	return &Ring{
		pubkeys: newRing,
		curve:   curve,
		hp:      hp,
	}, nil
}

// NewKeyRing creates a ring with size specified by `size` and places the public key corresponding
// to `privKey` in index idx of the ring.
// It returns a ring of public keys of length `size`.
func NewKeyRing(curve types.Curve, size int, privKey types.Scalar, idx int) (*Ring, error) {
	if idx >= size {
		return nil, errors.New("index out of bounds")
	}
	if privKey.IsZero() {
		return nil, errors.New("private key is zero")
	}

	ring := make([]types.Point, size)
	pubkey := curve.ScalarBaseMul(privKey)
	ring[idx] = pubkey

	for i := 0; i < size; i++ {
		if i == idx {
			continue
		}
		priv := curve.NewRandomScalar()
		ring[i] = curve.ScalarBaseMul(priv)
	}

	hp := make([]types.Point, size)
	for i := 0; i < size; i++ {
		hp[i] = hashToCurve(ring[i])
	}

	return &Ring{
		pubkeys: ring,
		curve:   curve,
		hp:      hp,
	}, nil
}

// Sign creates a ring signature on the given message using the public key ring
// and a private key of one of the members of the ring.
func (r *Ring) Sign(m [32]byte, privKey types.Scalar) (*RingSig, error) {
	ourIdx := -1
	pubkey := r.curve.ScalarBaseMul(privKey)
	for i, pk := range r.pubkeys {
		if pk.Equals(pubkey) {
			ourIdx = i
			break
		}
	}
	if ourIdx == -1 {
		return nil, errors.New("failed to find given key in public key set")
	}
	return Sign(m, r, privKey, ourIdx)
}

// ensureHP computes H_p(P_i) for all pubkeys if missing or out of date.
// It normalizes points through the ring's curve before hashing-to-curve to
// avoid "unsupported point type" panics from mixed concrete types.
func (r *Ring) ensureHP() error {
	if r == nil || r.curve == nil {
		return fmt.Errorf("nil ring/curve")
	}
	if r.pubkeys == nil || len(r.pubkeys) == 0 {
		return fmt.Errorf("no pubkeys in ring")
	}
	if r.hp != nil && len(r.hp) == len(r.pubkeys) {
		return nil
	}

	hp := make([]types.Point, len(r.pubkeys))
	for i, pk := range r.pubkeys {
		if pk == nil {
			return fmt.Errorf("nil pubkey at index %d", i)
		}
		// Normalize the point to the concrete type produced by this curve.
		enc := pk.Encode()
		dec, err := r.curve.DecodeToPoint(enc)
		if err != nil {
			return fmt.Errorf("failed to decode pubkey[%d]: %w", i, err)
		}
		hp[i] = hashToCurve(dec) // uses your existing helper
	}
	r.hp = hp
	return nil
}

// Sign creates a ring signature on the given message using the provided private key
// and ring of public keys.
func Sign(m [32]byte, ring *Ring, privKey types.Scalar, ourIdx int) (*RingSig, error) {
	size := len(ring.pubkeys)
	if size < 2 {
		return nil, errors.New("size of ring less than two")
	}
	if ourIdx >= size {
		return nil, errors.New("secret index out of range of ring size")
	}
	if privKey.IsZero() {
		return nil, errors.New("private key is zero")
	}

	// check that key at index s is indeed the signer
	pubkey := ring.curve.ScalarBaseMul(privKey)
	if !ring.pubkeys[ourIdx].Equals(pubkey) {
		return nil, errors.New("secret index in ring is not signer")
	}

	// setup
	curve := ring.curve
	h := hashToCurve(pubkey)
	sig := &RingSig{
		ring: ring,
		// calculate key image I = x * H_p(P) where H_p is a hash-to-curve function
		image: curve.ScalarMul(privKey, h),
	}

	// start at c[j]; pooled scratch
	c := getScalarScratch(size)
	defer putScalarScratch(c)

	// s IS RETAINED not needs to pool it
	s := make([]types.Scalar, size)

	// pick random scalar u, calculate L[j] = u*G
	u := curve.NewRandomScalar()
	l := curve.ScalarBaseMul(u)

	// compute R[j] = u*H_p(P[j])
	r := curve.ScalarMul(u, h)

	// calculate challenge c[j+1] = H(m, L_j, R_j)
	idx := (ourIdx + 1) % size
	c[idx] = challenge(ring.curve, m, l, r)

	// start loop at j+1
	for i := 1; i < size; i++ {
		idx := (ourIdx + i) % size
		if ring.pubkeys[idx] == nil {
			return nil, fmt.Errorf("no public key at index %d", idx)
		}

		// pick random scalar s_i
		s[idx] = curve.NewRandomScalar()

		// calculate L_i = s_i*G + c_i*P_i
		cP := curve.ScalarMul(c[idx], ring.pubkeys[idx])
		sG := curve.ScalarBaseMul(s[idx])
		l := cP.Add(sG)

		// calculate R_i = s_i*H_p(P_i) + c_i*I (use precomputed ring.hp[idx])
		cI := curve.ScalarMul(c[idx], sig.image)
		sH := curve.ScalarMul(s[idx], ring.hp[idx])
		r := cI.Add(sH)

		// calculate c[i+1] = H(m, L_i, R_i)
		c[(idx+1)%size] = challenge(curve, m, l, r)
	}

	// close the ring by finding s[j] = u - c[j]*x
	cx := c[ourIdx].Mul(privKey)
	s[ourIdx] = u.Sub(cx)

	// check that u*G = s[j]*G + c[j]*P[j]
	cP := curve.ScalarMul(c[ourIdx], pubkey)
	sG := curve.ScalarBaseMul(s[ourIdx])
	lNew := cP.Add(sG)
	if !lNew.Equals(l) {
		return nil, errors.New("failed to close ring: uG != sG + cP")
	}

	cI := curve.ScalarMul(c[ourIdx], sig.image)
	sH := curve.ScalarMul(s[ourIdx], h)
	rNew := cI.Add(sH)
	if !rNew.Equals(r) {
		return nil, errors.New("failed to close ring: uH(P) != sH(P) + cI")
	}

	cCheck := challenge(ring.curve, m, l, r)
	if !cCheck.Eq(c[(ourIdx+1)%size]) {
		return nil, errors.New("challenge check failed")
	}

	sig.s = s
	sig.c = c[0]
	return sig, nil
}

// Verify verifies the ring signature for the given message.
// It returns true if a valid signature, false otherwise.
func (sig *RingSig) Verify(m [32]byte) bool {
	if sig == nil || sig.ring == nil {
		return false
	}
	// setup
	ring := sig.ring
	size := len(ring.pubkeys)
	if size < 2 || len(sig.s) != size || sig.c == nil || sig.image == nil || ring.curve == nil {
		return false
	}
	if err := ring.ensureHP(); err != nil {
		return false
	}

	// pooled scratch for c[]
	c := getScalarScratch(size)
	defer putScalarScratch(c)

	c[0] = sig.c
	curve := ring.curve

	// calculate c[i+1] = H(m, s[i]*G + c[i]*P[i])
	// and c[0] = H)(m, s[n-1]*G + c[n-1]*P[n-1]) where n is the ring size
	for i := 0; i < size; i++ {
		// calculate L_i = s_i*G + c_i*P_i
		cP := curve.ScalarMul(c[i], ring.pubkeys[i])
		sG := curve.ScalarBaseMul(sig.s[i])
		l := cP.Add(sG)

		// calculate R_i = s_i*H_p(P_i) + c_i*I
		cI := curve.ScalarMul(c[i], sig.image)
		// use precomputed H_p(P_i)
		sH := curve.ScalarMul(sig.s[i], ring.hp[i])
		r := cI.Add(sH)

		// calculate c[i+1] = H(m, L_i, R_i)
		if i == size-1 {
			c[0] = challenge(curve, m, l, r)
		} else {
			c[i+1] = challenge(curve, m, l, r)
		}
	}

	return sig.c.Eq(c[0])
}

// Link returns true if the two signatures were created by the same signer,
// false otherwise.
func Link(sigA, sigB *RingSig) bool {
	switch sigA.Ring().curve.(type) {
	case *ed25519.CurveImpl:
		cofactor := Ed25519().ScalarFromInt(8)
		imageA := sigA.image.ScalarMul(cofactor)
		imageB := sigB.image.ScalarMul(cofactor)
		return imageA.Equals(imageB)
	default:
		return sigA.image.Equals(sigB.image)
	}
}

// Optimized challenge with pooled buffer and EncodeInto

// encodePointInto tries types.PointEncodeInto; falls back to Encode(). Returns bytes written.
func encodePointInto(p types.Point, dst []byte) int {
	if ei, ok := p.(types.PointEncodeInto); ok {
		return ei.EncodeInto(dst)
	}
	b := p.Encode()
	return copy(dst, b)
}

var challengeBufPool = sync.Pool{
	New: func() any {
		// default cap fits secp256k1: 32 (msg) + 33 + 33 = 98
		return make([]byte, 0, 98)
	},
}

func challenge(curve types.Curve, m [32]byte, l, r types.Point) types.Scalar {
	ps := curve.CompressedPointSize()
	need := 32 + 2*ps

	buf := challengeBufPool.Get().([]byte)
	if cap(buf) < need {
		buf = make([]byte, need)
	}
	buf = buf[:need]

	// [0:32) = message
	copy(buf[:32], m[:])

	// [32:32+ps) = L (compressed)
	off := 32
	_ = encodePointInto(l, buf[off:off+ps])

	// [32+ps:32+2*ps) = R (compressed)
	off += ps
	_ = encodePointInto(r, buf[off:off+ps])

	c, err := curve.HashToScalar(buf)

	// return buffer to pool
	challengeBufPool.Put(buf[:0])

	if err != nil {
		// this should not happen
		panic(err)
	}
	return c
}

// Pooled scratch for c[] in Verify/Sign

// bucketed pools to avoid resizing and reduce GC churn
var scalarPools = [...]struct {
	cap  int
	pool sync.Pool
}{
	{16, sync.Pool{New: func() any { return make([]types.Scalar, 0, 16) }}},
	{32, sync.Pool{New: func() any { return make([]types.Scalar, 0, 32) }}},
	{64, sync.Pool{New: func() any { return make([]types.Scalar, 0, 64) }}},
	{128, sync.Pool{New: func() any { return make([]types.Scalar, 0, 128) }}},
	{256, sync.Pool{New: func() any { return make([]types.Scalar, 0, 256) }}},
}

func getScalarScratch(n int) []types.Scalar {
	for i := range scalarPools {
		if n <= scalarPools[i].cap {
			b := scalarPools[i].pool.Get().([]types.Scalar)
			if cap(b) < n {
				// extremely unlikely, but be safe
				return make([]types.Scalar, n)
			}
			return b[:n]
		}
	}
	// fallback for huge rings (rare)
	return make([]types.Scalar, n)
}

func putScalarScratch(b []types.Scalar) {
	capb := cap(b)
	// clear interface elements to avoid retaining pointers
	for i := range b {
		b[i] = nil
	}
	b = b[:0]
	for i := range scalarPools {
		if capb == scalarPools[i].cap {
			scalarPools[i].pool.Put(b)
			return
		}
	}
	// non-pooled capacity: drop
}
