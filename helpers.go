package ring

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

// helper function, returns type of v
func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func padTo32Bytes(in []byte) (out []byte) {
	out = append(out, in...)
	for {
		if len(out) == 32 {
			return
		}
		out = append([]byte{0}, out...)
	}
}

// based off https://github.com/particl/particl-core/blob/master/src/secp256k1/src/modules/mlsag/main_impl.h#L139
func hashToCurve(pk *ecdsa.PublicKey) *ecdsa.PublicKey {
	const safety = 128
	compressedKey := crypto.CompressPubkey(pk)
	hash := sha3.Sum256(compressedKey)
	fe := &secp256k1.FieldVal{}
	fe.SetBytes(&hash)
	maybeY := &secp256k1.FieldVal{}

	for i := 0; i < safety; i++ {
		ok := secp256k1.DecompressY(fe, false, maybeY)
		if ok {
			return &ecdsa.PublicKey{
				Curve: secp256k1.S256(),
				X:     big.NewInt(0).SetBytes((fe.Bytes())[:]),
				Y:     big.NewInt(0).SetBytes((maybeY.Bytes())[:]),
			}
		}

		hash = sha3.Sum256(hash[:])
		fe.SetBytes(&hash)
	}

	return nil
}
