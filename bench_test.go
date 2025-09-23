package ring

import (
	"testing"

	"github.com/athanorlabs/go-dleq/types"
)

const idx = 0

func benchmarkSign(b *testing.B, curve types.Curve, keyring *Ring, privKey types.Scalar, size, idx int) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := keyring.Sign(testMsg, privKey)
		if err != nil {
			panic(err)
		}
	}
}

func mustKeyRing(curve types.Curve, privKey types.Scalar, size, idx int) *Ring {
	keyring, err := NewKeyRing(curve, size, privKey, idx)
	if err != nil {
		panic(err)
	}
	return keyring
}

func BenchmarkSign2_Secp256k1(b *testing.B)   { const size = 2; runSignBench(b, Secp256k1(), size) }
func BenchmarkSign4_Secp256k1(b *testing.B)   { const size = 4; runSignBench(b, Secp256k1(), size) }
func BenchmarkSign8_Secp256k1(b *testing.B)   { const size = 8; runSignBench(b, Secp256k1(), size) }
func BenchmarkSign16_Secp256k1(b *testing.B)  { const size = 16; runSignBench(b, Secp256k1(), size) }
func BenchmarkSign32_Secp256k1(b *testing.B)  { const size = 32; runSignBench(b, Secp256k1(), size) }
func BenchmarkSign64_Secp256k1(b *testing.B)  { const size = 64; runSignBench(b, Secp256k1(), size) }
func BenchmarkSign128_Secp256k1(b *testing.B) { const size = 128; runSignBench(b, Secp256k1(), size) }

func BenchmarkSign2_Ed25519(b *testing.B)   { const size = 2; runSignBench(b, Ed25519(), size) }
func BenchmarkSign4_Ed25519(b *testing.B)   { const size = 4; runSignBench(b, Ed25519(), size) }
func BenchmarkSign8_Ed25519(b *testing.B)   { const size = 8; runSignBench(b, Ed25519(), size) }
func BenchmarkSign16_Ed25519(b *testing.B)  { const size = 16; runSignBench(b, Ed25519(), size) }
func BenchmarkSign32_Ed25519(b *testing.B)  { const size = 32; runSignBench(b, Ed25519(), size) }
func BenchmarkSign64_Ed25519(b *testing.B)  { const size = 64; runSignBench(b, Ed25519(), size) }
func BenchmarkSign128_Ed25519(b *testing.B) { const size = 128; runSignBench(b, Ed25519(), size) }

func runSignBench(b *testing.B, curve types.Curve, size int) {
	privKey := curve.NewRandomScalar()
	keyring := mustKeyRing(curve, privKey, size, idx)
	benchmarkSign(b, curve, keyring, privKey, size, idx)
}

func benchmarkVerify(b *testing.B, sig *RingSig) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ok := sig.Verify(testMsg)
		if !ok {
			panic("did not verify signature")
		}
	}
}

func mustSig(curve types.Curve, size int) *RingSig {
	privKey := curve.NewRandomScalar()
	keyring := mustKeyRing(curve, privKey, size, idx)

	sig, err := keyring.Sign(testMsg, privKey)
	if err != nil {
		panic(err)
	}

	return sig
}

func BenchmarkVerify2_Secp256k1(b *testing.B)  { const size = 2; runVerifyBench(b, Secp256k1(), size) }
func BenchmarkVerify4_Secp256k1(b *testing.B)  { const size = 4; runVerifyBench(b, Secp256k1(), size) }
func BenchmarkVerify8_Secp256k1(b *testing.B)  { const size = 8; runVerifyBench(b, Secp256k1(), size) }
func BenchmarkVerify16_Secp256k1(b *testing.B) { const size = 16; runVerifyBench(b, Secp256k1(), size) }
func BenchmarkVerify32_Secp256k1(b *testing.B) { const size = 32; runVerifyBench(b, Secp256k1(), size) }
func BenchmarkVerify64_Secp256k1(b *testing.B) { const size = 64; runVerifyBench(b, Secp256k1(), size) }
func BenchmarkVerify128_Secp256k1(b *testing.B) {
	const size = 128
	runVerifyBench(b, Secp256k1(), size)
}

func BenchmarkVerify2_Ed25519(b *testing.B)   { const size = 2; runVerifyBench(b, Ed25519(), size) }
func BenchmarkVerify4_Ed25519(b *testing.B)   { const size = 4; runVerifyBench(b, Ed25519(), size) }
func BenchmarkVerify8_Ed25519(b *testing.B)   { const size = 8; runVerifyBench(b, Ed25519(), size) }
func BenchmarkVerify16_Ed25519(b *testing.B)  { const size = 16; runVerifyBench(b, Ed25519(), size) }
func BenchmarkVerify32_Ed25519(b *testing.B)  { const size = 32; runVerifyBench(b, Ed25519(), size) }
func BenchmarkVerify64_Ed25519(b *testing.B)  { const size = 64; runVerifyBench(b, Ed25519(), size) }
func BenchmarkVerify128_Ed25519(b *testing.B) { const size = 128; runVerifyBench(b, Ed25519(), size) }

func runVerifyBench(b *testing.B, curve types.Curve, size int) {
	sig := mustSig(curve, size)
	benchmarkVerify(b, sig)
}

/*** New micro-benchmarks to surface the optimizations ***/

// 1) Compare copy vs zero-copy pubkey access
func BenchmarkPublicKeysCopy_Secp256k1_Size16(b *testing.B) {
	const size = 16
	curve := Secp256k1()
	sig := mustSig(curve, size)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = sig.PublicKeys() // copies
	}
}

func BenchmarkPublicKeysRef_Secp256k1_Size16(b *testing.B) {
	const size = 16
	curve := Secp256k1()
	sig := mustSig(curve, size)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = sig.PublicKeysRef() // zero-copy
	}
}

// 2) Exercise the optimized challenge path (single allocation)
func BenchmarkChallenge_Secp256k1(b *testing.B) {
	const size = 8
	curve := Secp256k1()
	priv := curve.NewRandomScalar()
	r := mustKeyRing(curve, priv, size, idx)

	// Prepare inputs that resemble the Verify loop
	p := r.pubkeys[0]
	l := curve.ScalarBaseMul(priv)
	hp := r.hp[0]
	rp := curve.ScalarMul(priv, hp)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = challenge(curve, testMsg, p.Add(l), rp)
	}
}

// Ed25519 challenge too (apples-to-apples)
func BenchmarkChallenge_Ed25519(b *testing.B) {
	curve := Ed25519()
	priv := curve.NewRandomScalar()
	r, err := NewKeyRing(curve, 4, priv, 0)
	if err != nil {
		b.Fatal(err)
	}
	L := r.pubkeys[0]
	R := r.pubkeys[1]

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = challenge(curve, testMsg, L, R)
	}
}

// EncodeInto vs Encode (secp256k1)
func BenchmarkPointEncodeInto_Secp256k1(b *testing.B) {
	curve := Secp256k1()
	p := curve.ScalarBaseMul(curve.NewRandomScalar())
	dst := make([]byte, curve.CompressedPointSize())

	pi, ok := p.(types.PointEncodeInto)
	if !ok {
		b.Skip("secp256k1 point does not implement EncodeInto")
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pi.EncodeInto(dst)
	}
}

func BenchmarkPointEncode_Secp256k1(b *testing.B) {
	curve := Secp256k1()
	p := curve.ScalarBaseMul(curve.NewRandomScalar())
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = p.Encode()
	}
}

// EncodeInto vs Encode (ed25519)
func BenchmarkPointEncodeInto_Ed25519(b *testing.B) {
	curve := Ed25519()
	p := curve.ScalarBaseMul(curve.NewRandomScalar())
	dst := make([]byte, curve.CompressedPointSize())

	pi, ok := p.(types.PointEncodeInto)
	if !ok {
		b.Skip("ed25519 point does not implement EncodeInto")
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pi.EncodeInto(dst)
	}
}

func BenchmarkPointEncode_Ed25519(b *testing.B) {
	curve := Ed25519()
	p := curve.ScalarBaseMul(curve.NewRandomScalar())
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = p.Encode()
	}
}

// Parallel verify throughput
func benchmarkVerifyParallel(b *testing.B, curve types.Curve, size int) {
	priv := curve.NewRandomScalar()
	r := mustKeyRing(curve, priv, size, idx)
	sig, err := r.Sign(testMsg, priv)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if !sig.Verify(testMsg) {
				panic("verify failed")
			}
		}
	})
}

// secp256k1 parallel
func BenchmarkVerifyParallel2_Secp256k1(b *testing.B)  { benchmarkVerifyParallel(b, Secp256k1(), 2) }
func BenchmarkVerifyParallel4_Secp256k1(b *testing.B)  { benchmarkVerifyParallel(b, Secp256k1(), 4) }
func BenchmarkVerifyParallel8_Secp256k1(b *testing.B)  { benchmarkVerifyParallel(b, Secp256k1(), 8) }
func BenchmarkVerifyParallel16_Secp256k1(b *testing.B) { benchmarkVerifyParallel(b, Secp256k1(), 16) }
func BenchmarkVerifyParallel32_Secp256k1(b *testing.B) { benchmarkVerifyParallel(b, Secp256k1(), 32) }
func BenchmarkVerifyParallel64_Secp256k1(b *testing.B) { benchmarkVerifyParallel(b, Secp256k1(), 64) }
func BenchmarkVerifyParallel128_Secp256k1(b *testing.B) {
	benchmarkVerifyParallel(b, Secp256k1(), 128)
}

// ed25519 parallel (same sizes)
func BenchmarkVerifyParallel2_Ed25519(b *testing.B)   { benchmarkVerifyParallel(b, Ed25519(), 2) }
func BenchmarkVerifyParallel4_Ed25519(b *testing.B)   { benchmarkVerifyParallel(b, Ed25519(), 4) }
func BenchmarkVerifyParallel8_Ed25519(b *testing.B)   { benchmarkVerifyParallel(b, Ed25519(), 8) }
func BenchmarkVerifyParallel16_Ed25519(b *testing.B)  { benchmarkVerifyParallel(b, Ed25519(), 16) }
func BenchmarkVerifyParallel32_Ed25519(b *testing.B)  { benchmarkVerifyParallel(b, Ed25519(), 32) }
func BenchmarkVerifyParallel64_Ed25519(b *testing.B)  { benchmarkVerifyParallel(b, Ed25519(), 64) }
func BenchmarkVerifyParallel128_Ed25519(b *testing.B) { benchmarkVerifyParallel(b, Ed25519(), 128) }

// Verify loop attribution: HP vs no-HP (secp256k1 & ed25519, size=32)
func verifyLoopWithHP(curve types.Curve, pubkeys []types.Point, hp []types.Point, s []types.Scalar, img types.Point, m [32]byte, c0 types.Scalar) {
	n := len(pubkeys)
	c := make([]types.Scalar, n)
	c[0] = c0
	for i := 0; i < n; i++ {
		cP := curve.ScalarMul(c[i], pubkeys[i])
		sG := curve.ScalarBaseMul(s[i])
		l := cP.Add(sG)

		cI := curve.ScalarMul(c[i], img)
		sH := curve.ScalarMul(s[i], hp[i])
		r := cI.Add(sH)

		if i == n-1 {
			c[0] = challenge(curve, m, l, r)
		} else {
			c[i+1] = challenge(curve, m, l, r)
		}
	}
}

func verifyLoopNoHP(curve types.Curve, pubkeys []types.Point, s []types.Scalar, img types.Point, m [32]byte, c0 types.Scalar) {
	n := len(pubkeys)
	c := make([]types.Scalar, n)
	c[0] = c0
	for i := 0; i < n; i++ {
		cP := curve.ScalarMul(c[i], pubkeys[i])
		sG := curve.ScalarBaseMul(s[i])
		l := cP.Add(sG)

		cI := curve.ScalarMul(c[i], img)
		h := hashToCurve(pubkeys[i]) // recompute each time
		sH := curve.ScalarMul(s[i], h)
		r := cI.Add(sH)

		if i == n-1 {
			c[0] = challenge(curve, m, l, r)
		} else {
			c[i+1] = challenge(curve, m, l, r)
		}
	}
}

func prepSig(b *testing.B, curve types.Curve, size int) (*RingSig, *Ring) {
	priv := curve.NewRandomScalar()
	r := mustKeyRing(curve, priv, size, idx)
	sig, err := r.Sign(testMsg, priv)
	if err != nil {
		b.Fatal(err)
	}
	return sig, r
}

func BenchmarkVerifyLoop_HP_vs_NoHP_Secp256k1_32(b *testing.B) {
	const size = 32
	sig, r := prepSig(b, Secp256k1(), size)

	b.Run("with_hp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			verifyLoopWithHP(r.curve, r.pubkeys, r.hp, sig.s, sig.image, testMsg, sig.c)
		}
	})
	b.Run("no_hp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			verifyLoopNoHP(r.curve, r.pubkeys, sig.s, sig.image, testMsg, sig.c)
		}
	})
}

func BenchmarkVerifyLoop_HP_vs_NoHP_Ed25519_32(b *testing.B) {
	const size = 32
	sig, r := prepSig(b, Ed25519(), size)

	b.Run("with_hp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			verifyLoopWithHP(r.curve, r.pubkeys, r.hp, sig.s, sig.image, testMsg, sig.c)
		}
	})
	b.Run("no_hp", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			verifyLoopNoHP(r.curve, r.pubkeys, sig.s, sig.image, testMsg, sig.c)
		}
	})
}

// Batch/stream verify of many sigs over same ring (common in handlers)
func benchmarkVerifyBatchSameRing(b *testing.B, curve types.Curve, size, batch int) {
	priv := curve.NewRandomScalar()
	r := mustKeyRing(curve, priv, size, idx)

	sigs := make([]*RingSig, batch)
	for i := 0; i < batch; i++ {
		s, err := r.Sign(testMsg, priv)
		if err != nil {
			b.Fatal(err)
		}
		sigs[i] = s
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < batch; j++ {
			if !sigs[j].Verify(testMsg) {
				b.Fatal("verify failed")
			}
		}
	}
}

// secp256k1 batch benches (sizes match your main benches; batch=64)
func BenchmarkVerifyBatchSameRing_2x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 2, 64)
}
func BenchmarkVerifyBatchSameRing_4x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 4, 64)
}
func BenchmarkVerifyBatchSameRing_8x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 8, 64)
}
func BenchmarkVerifyBatchSameRing_16x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 16, 64)
}
func BenchmarkVerifyBatchSameRing_32x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 32, 64)
}
func BenchmarkVerifyBatchSameRing_64x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 64, 64)
}
func BenchmarkVerifyBatchSameRing_128x64_Secp256k1(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Secp256k1(), 128, 64)
}

// ed25519 batch benches (same sizes, batch=64)
func BenchmarkVerifyBatchSameRing_2x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 2, 64)
}
func BenchmarkVerifyBatchSameRing_4x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 4, 64)
}
func BenchmarkVerifyBatchSameRing_8x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 8, 64)
}
func BenchmarkVerifyBatchSameRing_16x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 16, 64)
}
func BenchmarkVerifyBatchSameRing_32x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 32, 64)
}
func BenchmarkVerifyBatchSameRing_64x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 64, 64)
}
func BenchmarkVerifyBatchSameRing_128x64_Ed25519(b *testing.B) {
	benchmarkVerifyBatchSameRing(b, Ed25519(), 128, 64)
}
