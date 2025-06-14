module github.com/vigilnetwork/vgl

go 1.23.0

toolchain go1.24.4

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/btcsuite/btcutil v1.0.2
	github.com/vigilnetwork/vgl/addrmgr/v3 v3.0.0
	github.com/vigilnetwork/vgl/bech32 v1.1.4
	github.com/vigilnetwork/vgl/blockchain/stake/v5 v5.0.1
	github.com/vigilnetwork/vgl/blockchain/standalone/v2 v2.2.1
	github.com/vigilnetwork/vgl/blockchain/v5 v5.0.1
	github.com/vigilnetwork/vgl/certgen v1.1.3
	github.com/vigilnetwork/vgl/chaincfg/chainhash v1.0.4
	github.com/vigilnetwork/vgl/chaincfg/v3 v3.2.1
	github.com/vigilnetwork/vgl/connmgr/v3 v3.1.2
	github.com/vigilnetwork/vgl/container/apbf v1.0.1
	github.com/vigilnetwork/vgl/container/lru v1.0.0
	github.com/vigilnetwork/vgl/crypto/rand v1.0.1
	github.com/vigilnetwork/vgl/crypto/ripemd160 v1.0.2
	github.com/vigilnetwork/vgl/database/v3 v3.0.2
	github.com/vigilnetwork/vgl/dcrec v1.0.1
	github.com/vigilnetwork/vgl/dcrec/secp256k1/v4 v4.3.0
	github.com/vigilnetwork/vgl/dcrjson/v4 v4.1.0
	github.com/vigilnetwork/vgl/dcrutil/v4 v4.0.2
	github.com/vigilnetwork/vgl/gcs/v4 v4.1.0
	github.com/vigilnetwork/vgl/math/uint256 v1.0.2
	github.com/vigilnetwork/vgl/mixing v0.3.0
	github.com/vigilnetwork/vgl/peer/v3 v3.1.1
	github.com/vigilnetwork/vgl/rpc/jsonrpc/types/v4 v4.2.0
	github.com/vigilnetwork/vgl/rpcclient/v8 v8.0.1
	github.com/vigilnetwork/vgl/txscript/v4 v4.1.1
	github.com/vigilnetwork/vgl/wire v1.7.0
	github.com/Vigil/dcrtest/vgldtest v1.0.1-0.20240404170936-a2529e936df1
	github.com/Vigil/go-socks v1.1.0
	github.com/vigilnetwork/vgl/slog v1.2.0
	github.com/gorilla/websocket v1.5.1
	github.com/jessevdk/go-flags v1.5.0
	github.com/jrick/bitset v1.0.0
	github.com/jrick/logrotate v1.0.0
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	golang.org/x/net v0.28.0
	golang.org/x/sys v0.33.0
	golang.org/x/term v0.32.0
	lukechampine.com/blake3 v1.3.0
)

require github.com/vigilnetwork/vgl/crypto/blake256 v1.1.0 // indirect

require (
	vigil.network/vgl/cspp/v2 v2.4.0 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/companyzero/sntrup4591761 v0.0.0-20220309191932-9e0f3af2f07a // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/vigilnetwork/vgl/chaincfg v1.5.1 // indirect
	github.com/vigilnetwork/vgl/chaincfg/v2 v2.0.2 // indirect
	github.com/vigilnetwork/vgl/dcrec/edwards v1.0.0 // indirect
	github.com/vigilnetwork/vgl/dcrec/edwards/v2 v2.0.3 // indirect
	github.com/vigilnetwork/vgl/dcrec/secp256k1 v1.0.2 // indirect
	github.com/vigilnetwork/vgl/dcrutil v1.4.1 // indirect
	github.com/vigilnetwork/vgl/dcrutil/v2 v2.0.0 // indirect
	github.com/vigilnetwork/vgl/hdkeychain/v3 v3.1.2 // indirect
	github.com/vigilnetwork/vgl/kawpow v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)

replace (
	github.com/vigilnetwork/vgl/addrmgr/v3 => ./addrmgr
	github.com/vigilnetwork/vgl/bech32 => ./bech32
	github.com/vigilnetwork/vgl/blockchain/stake/v5 => ./blockchain/stake
	github.com/vigilnetwork/vgl/blockchain/standalone/v2 => ./blockchain/standalone
	github.com/vigilnetwork/vgl/blockchain/v5 => ./blockchain
	github.com/vigilnetwork/vgl/certgen => ./certgen
	github.com/vigilnetwork/vgl/chaincfg/chainhash => ./chaincfg/chainhash
	github.com/vigilnetwork/vgl/chaincfg/v3 => ./chaincfg
	github.com/vigilnetwork/vgl/connmgr/v3 => ./connmgr
	github.com/vigilnetwork/vgl/container/apbf => ./container/apbf
	github.com/vigilnetwork/vgl/container/lru => ./container/lru
	github.com/vigilnetwork/vgl/crypto/blake256 => ./crypto/blake256
	github.com/vigilnetwork/vgl/crypto/rand => ./crypto/rand
	github.com/vigilnetwork/vgl/crypto/ripemd160 => ./crypto/ripemd160
	github.com/vigilnetwork/vgl/database/v3 => ./database
	github.com/vigilnetwork/vgl/dcrec => ./dcrec
	github.com/vigilnetwork/vgl/dcrec/secp256k1/v4 => ./dcrec/secp256k1
	github.com/vigilnetwork/vgl/dcrjson/v4 => ./dcrjson
	github.com/vigilnetwork/vgl/dcrutil/v4 => ./dcrutil
	github.com/vigilnetwork/vgl/gcs/v4 => ./gcs
	github.com/vigilnetwork/vgl/hdkeychain/v3 => ./hdkeychain
	github.com/vigilnetwork/vgl/limits => ./limits
	github.com/vigilnetwork/vgl/math/uint256 => ./math/uint256
	github.com/vigilnetwork/vgl/mixing => ./mixing
	github.com/vigilnetwork/vgl/peer/v3 => ./peer
	github.com/vigilnetwork/vgl/rpc/jsonrpc/types/v4 => ./rpc/jsonrpc/types
	github.com/vigilnetwork/vgl/rpcclient/v8 => ./rpcclient
	github.com/vigilnetwork/vgl/txscript/v4 => ./txscript
	github.com/vigilnetwork/vgl/wire => ./wire

)

replace github.com/vigilnetwork/vgl/kawpow => ./kawpow
