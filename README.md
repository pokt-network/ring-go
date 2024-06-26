# ring-go <!-- omit in toc -->

Implementation of linkable ring signatures using elliptic curve crypto in pure Go.
It supports ring signatures over both ed25519 and secp256k1.

- [Requirements](#requirements)
- [Install](#install)
- [References](#references)
- [Usage](#usage)

## Requirements

go 1.22.2

## Install

`go get github.com/pokt-network/ring-go`

## References

This implementation is based off of [Ring Confidential Transactions](https://eprint.iacr.org/2015/1098.pdf), in particular section 2, which defines MLSAG (Multilayered Linkable Spontaneous Anonymous Group signatures).

## Usage

See `examples/main.go`.

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
