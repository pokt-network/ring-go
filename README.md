# ring-go <!-- omit in toc -->

Implementation of linkable ring signatures using elliptic curve crypto in Go.
It supports ring signatures over both ed25519 and secp256k1 with **pluggable crypto backends** for optimal performance.

- [Requirements](#requirements)
- [Install](#install)
- [Crypto Backends](#crypto-backends)
  - [Comparison](#comparison)
  - [Build Commands](#build-commands)
    - [Portable (Pure Go)](#portable-pure-go)
    - [High Performance (CGO)](#high-performance-cgo)
  - [Benchmarking Crypto Backends](#benchmarking-crypto-backends)
- [References](#references)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [High-Performance Usage](#high-performance-usage)

## Requirements

- Go 1.23.0 or later
- **Optional**: CGO and libsecp256k1 for high-performance backend

## Install

```bash
go get github.com/pokt-network/ring-go
```

## Crypto Backends

⚠️ **The crypto backend is a BUILD-TIME, not a RUN-TIME configuration** ⚠️

Ring-go supports multiple secp256k1 crypto backends for optimal performance vs portability tradeoffs.

### Comparison

| Backend                | Build Configuration                           | Dependencies                      | Performance                                   | Portability                |
| ---------------------- | -------------------------------------------- | --------------------------------- | --------------------------------------------- | -------------------------- |
| **Decred (Pure Go)**   | Default (no tags)                           | None                              | Excellent (pure Go secp256k1)                | Runs anywhere              |
| **Ethereum secp256k1** | `CGO_ENABLED=1` + `-tags=ethereum_secp256k1` | `gcc`, `libsecp256k1`             | ~50% faster signing, ~80% faster verification | Requires CGO + system libs |

### Build Commands

#### Portable (Pure Go)

**Default behavior** - no special configuration needed:

```bash
go build ./examples/...
# OR using Makefile
make build_portable
```

#### High Performance (CGO)

Build with Ethereum's libsecp256k1 backend:

```bash
go build -tags=ethereum_secp256k1 ./examples/...
# OR using Makefile
make build_fast
```

Auto-select optimal backend:
```bash
make build_auto  # Chooses fast if CGO available, portable otherwise
```

### Benchmarking Crypto Backends

Run performance comparison:

```bash
make benchmark_report
```

Example output:
```
🔍 SIGN PERFORMANCE (Ring Signatures):
Ring Size  Backend         Time/op      Memory/op    Allocs/op    Performance
--------   -------         --------     ---------    ---------    -----------
2          Ethereum        780.5 μs     3.2 KB       45           🥇
2          Decred          1.55 ms      5.0 KB       84           🥈
32         Ethereum        8.4 ms       57.7 KB      926          🥇
32         Decred          16.8 ms      114.4 KB     1828         🥈
```

Run all benchmarks:
```bash
make benchmark_all
```

## References

This implementation is based off of [Ring Confidential Transactions](https://eprint.iacr.org/2015/1098.pdf), in particular section 2, which defines MLSAG (Multilayered Linkable Spontaneous Anonymous Group signatures).

## Usage

### Basic Usage

See `examples/main.go` for complete examples.

**Standard ring signatures using default backend:**

```go
package main

import (
    "fmt"

    ring "github.com/pokt-network/ring-go"
    "golang.org/x/crypto/sha3"
)

func signAndVerify(curve ring.Curve) {
    privKey := curve.NewRandomScalar()
    msgHash := sha3.Sum256([]byte("helloworld"))

    // size of the public key ring (anonymity set)
    const size = 16

    // our key's secret index within the set
    const idx = 7

    keyring, err := ring.NewKeyRing(curve, size, privKey, idx)
    if err != nil {
        panic(err)
    }

    sig, err := keyring.Sign(msgHash, privKey)
    if err != nil {
        panic(err)
    }

    ok := sig.Verify(msgHash)
    if !ok {
        fmt.Println("failed to verify :(")
        return
    }

    fmt.Println("verified signature!")
}

func main() {
    fmt.Println("using secp256k1...")
    signAndVerify(ring.Secp256k1())
    fmt.Println("using ed25519...")
    signAndVerify(ring.Ed25519())
}
```

### High-Performance Usage

**Using the fast backend for performance-critical applications:**

```go
package main

import (
    "fmt"

    ring "github.com/pokt-network/ring-go"
    "golang.org/x/crypto/sha3"
)

func highPerformanceSignAndVerify() {
    // Use the fast backend (auto-selects optimal implementation)
    curve := ring.Secp256k1Fast()

    privKey := curve.NewRandomScalar()
    msgHash := sha3.Sum256([]byte("high-performance ring signature"))

    // Larger ring size for stronger anonymity
    const size = 64
    const idx = 32

    keyring, err := ring.NewKeyRing(curve, size, privKey, idx)
    if err != nil {
        panic(err)
    }

    // ~50% faster signing with Ethereum backend
    sig, err := keyring.Sign(msgHash, privKey)
    if err != nil {
        panic(err)
    }

    // ~80% faster verification with Ethereum backend
    ok := sig.Verify(msgHash)
    if !ok {
        fmt.Println("failed to verify :(")
        return
    }

    fmt.Printf("✅ Verified ring signature with %d-member anonymity set!\n", size)
}

func main() {
    highPerformanceSignAndVerify()
}
```

**Quick demo:**

```bash
make demo
```
