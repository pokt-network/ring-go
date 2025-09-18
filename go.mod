module github.com/pokt-network/ring-go

go 1.24.0

toolchain go1.24.3

replace github.com/athanorlabs/go-dleq => github.com/jorgecuesta/go-dleq v0.0.0-20250918223310-7a1fc288336f

require (
	filippo.io/edwards25519 v1.1.0
	github.com/athanorlabs/go-dleq v0.1.0
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.4.0
	github.com/stretchr/testify v1.11.1
	golang.org/x/crypto v0.42.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
