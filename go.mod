module github.com/pokt-network/ring-go

go 1.24.3

toolchain go1.24.5

require (
	filippo.io/edwards25519 v1.0.0
	github.com/athanorlabs/go-dleq v0.1.0
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0
	github.com/ethereum/go-ethereum v1.14.12
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.36.0
)

replace github.com/athanorlabs/go-dleq => /Users/olshansky/workspace/pocket/go-dleq

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/holiman/uint256 v1.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
