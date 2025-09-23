package ring

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/athanorlabs/go-dleq/types"
)

// Serialize converts the signature to a byte array.
func (r *RingSig) Serialize() ([]byte, error) {
	sig := []byte{}
	size := len(r.ring.pubkeys)

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(size))
	sig = append(sig, b[:]...)
	sig = append(sig, r.c.Encode()...)
	sig = append(sig, r.image.Encode()...)

	for i := 0; i < size; i++ {
		sig = append(sig, r.s[i].Encode()...)
		sig = append(sig, r.ring.pubkeys[i].Encode()...)
	}

	return sig, nil
}

// Deserialize converts the byteified signature into a *RingSig.
func (sig *RingSig) Deserialize(curve types.Curve, in []byte) error {
	// WARN: this assumes the groups have an encoded scalar length of 32!
	// which is fine for ed25519 and secp256k1, but may need to be changed
	// if other curves are added.
	const scalarLen = 32

	if len(in) < 4 {
		return errors.New("input too short: missing size header")
	}

	// total size sanity check
	size := int(binary.BigEndian.Uint32(in[:4]))
	if size < 2 {
		return fmt.Errorf("invalid ring size: %d", size)
	}

	ps := curve.CompressedPointSize()
	// Minimum expected bytes: 4(size) + 32(c) + ps(image) + size*(32(s_i) + ps(P_i))
	m := 4 + scalarLen + ps + size*(scalarLen+ps)
	if len(in) < m {
		return fmt.Errorf("input too short: got %d, need at least %d", len(in), m)
	}

	// Reader over the remaining bytes
	r := bytes.NewReader(in[4:])

	// c
	cBytes := make([]byte, scalarLen)
	if _, err := io.ReadFull(r, cBytes); err != nil {
		return fmt.Errorf("read c: %w", err)
	}
	c, err := curve.DecodeToScalar(cBytes)
	if err != nil {
		return fmt.Errorf("decode c: %w", err)
	}

	// image
	imgBytes := make([]byte, ps)
	if _, err := io.ReadFull(r, imgBytes); err != nil {
		return fmt.Errorf("read image: %w", err)
	}
	img, err := curve.DecodeToPoint(imgBytes)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	// s[i] and P[i]
	s := make([]types.Scalar, size)
	pubkeys := make([]types.Point, size)

	for i := 0; i < size; i++ {
		sb := make([]byte, scalarLen)
		if _, err := io.ReadFull(r, sb); err != nil {
			return fmt.Errorf("read s[%d]: %w", i, err)
		}
		si, err := curve.DecodeToScalar(sb)
		if err != nil {
			return fmt.Errorf("decode s[%d]: %w", i, err)
		}
		s[i] = si

		pb := make([]byte, ps)
		if _, err := io.ReadFull(r, pb); err != nil {
			return fmt.Errorf("read pubkey[%d]: %w", i, err)
		}
		pi, err := curve.DecodeToPoint(pb)
		if err != nil {
			return fmt.Errorf("decode pubkey[%d]: %w", i, err)
		}
		pubkeys[i] = pi
	}

	// Build ring and precompute hp AFTER pubkeys exist
	ring := &Ring{
		pubkeys: pubkeys,
		curve:   curve,
		hp:      make([]types.Point, size),
	}
	for i := 0; i < size; i++ {
		if ring.pubkeys[i] == nil {
			return fmt.Errorf("nil pubkey at index %d", i)
		}
		ring.hp[i] = hashToCurve(ring.pubkeys[i])
	}

	sig.ring = ring
	sig.c = c
	sig.s = s
	sig.image = img
	return nil
}
