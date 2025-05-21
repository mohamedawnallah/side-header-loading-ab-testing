module example

go 1.24.2

replace github.com/lightninglabs/neutrino => /home/dev/neutrino

replace github.com/btcsuite/btcd/btcec => github.com/btcsuite/btcd/btcec/v2 v2.2.0

require (
	github.com/btcsuite/btcd v0.24.3-0.20250318170759-4f4ea81776d6
	github.com/btcsuite/btcd/btcutil v1.1.5
	github.com/btcsuite/btclog v1.0.0
	github.com/btcsuite/btcwallet/walletdb v1.5.1
	github.com/lightninglabs/neutrino v0.0.0-00010101000000-000000000000
)

require (
	github.com/aead/siphash v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.3.4 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0 // indirect
	github.com/btcsuite/btcwallet/wtxmgr v1.5.0 // indirect
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd // indirect
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.0.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/decred/dcrd/lru v1.1.2 // indirect
	github.com/kkdai/bstream v1.0.0 // indirect
	github.com/lightninglabs/neutrino/cache v1.1.2 // indirect
	github.com/lightningnetwork/lnd/clock v1.0.1 // indirect
	github.com/lightningnetwork/lnd/queue v1.0.1 // indirect
	github.com/lightningnetwork/lnd/ticker v1.0.0 // indirect
	go.etcd.io/bbolt v1.3.11 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
