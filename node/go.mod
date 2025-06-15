module github.com/kdsmith18542/vigil

go 1.23.0

toolchain go1.24.4

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/btcsuite/btcutil v1.0.2
	github.com/kdsmith18542/vigil/addrmgr/v3 v3.0.0
	github.com/kdsmith18542/vigil/bech32 v1.1.4
	github.com/kdsmith18542/vigil/blockchain/stake/v5 v5.0.1
	github.com/kdsmith18542/vigil/blockchain/standalone/v2 v2.2.1
	github.com/kdsmith18542/vigil/blockchain/v5 v5.0.1
	github.com/kdsmith18542/vigil/certgen v1.1.3
	github.com/kdsmith18542/vigil/chaincfg/chainhash v1.0.4
	github.com/kdsmith18542/vigil/chaincfg/v3 v3.2.1
	github.com/kdsmith18542/vigil/connmgr/v3 v3.1.2
	github.com/kdsmith18542/vigil/container/apbf v1.0.1
	github.com/kdsmith18542/vigil/container/lru v1.0.0
	github.com/kdsmith18542/vigil/crypto/rand v1.0.1
	github.com/kdsmith18542/vigil/crypto/ripemd160 v1.0.2
	github.com/kdsmith18542/vigil/database/v3 v3.0.2
	github.com/kdsmith18542/vigil/dcrec v1.0.1
	github.com/kdsmith18542/vigil/dcrec/secp256k1/v4 v4.3.0
	github.com/kdsmith18542/vigil/dcrjson/v4 v4.1.0
	github.com/kdsmith18542/vigil/dcrutil/v4 v4.0.2
	github.com/kdsmith18542/vigil/gcs/v4 v4.1.0
	github.com/kdsmith18542/vigil/math/uint256 v1.0.2
	github.com/kdsmith18542/vigil/mixing v0.3.0
	github.com/kdsmith18542/vigil/peer/v3 v3.1.1
	github.com/kdsmith18542/vigil/rpc/jsonrpc/types/v4 v4.2.0
	github.com/kdsmith18542/vigil/rpcclient/v8 v8.0.1
	github.com/kdsmith18542/vigil/txscript/v4 v4.1.1
	github.com/kdsmith18542/vigil/wire v1.7.0
	github.com/kdsmith18542/vigil/dcrtest/vgldtest v1.0.1-0.20240404170936-a2529e936df1
	github.com/kdsmith18542/vigil/go-socks v1.1.0
	github.com/kdsmith18542/vigil/slog v1.2.0
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

require github.com/kdsmith18542/vigil/crypto/blake256 v1.1.0 // indirect

require (
	vigil.network/vgl/cspp/v2 v2.4.0 // indirect
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412 // indirect
	github.com/companyzero/sntrup4591761 v0.0.0-20220309191932-9e0f3af2f07a // indirect
	github.com/dchest/siphash v1.2.3 // indirect
	github.com/kdsmith18542/vigil/chaincfg v1.5.1 // indirect
	github.com/kdsmith18542/vigil/chaincfg/v2 v2.0.2 // indirect
	github.com/kdsmith18542/vigil/dcrec/edwards v1.0.0 // indirect
	github.com/kdsmith18542/vigil/dcrec/edwards/v2 v2.0.3 // indirect
	github.com/kdsmith18542/vigil/dcrec/secp256k1 v1.0.2 // indirect
	github.com/kdsmith18542/vigil/dcrutil v1.4.1 // indirect
	github.com/kdsmith18542/vigil/dcrutil/v2 v2.0.0 // indirect
	github.com/kdsmith18542/vigil/hdkeychain/v3 v3.1.2 // indirect
	github.com/kdsmith18542/vigil/kawpow v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/text v0.26.0 // indirect
)

replace (
	github.com/kdsmith18542/vigil/addrmgr/v3 => ./addrmgr
	github.com/kdsmith18542/vigil/bech32 => ./bech32
	github.com/kdsmith18542/vigil/blockchain/stake/v5 => ./blockchain/stake
	github.com/kdsmith18542/vigil/blockchain/standalone/v2 => ./blockchain/standalone
	github.com/kdsmith18542/vigil/blockchain/v5 => ./blockchain
	github.com/kdsmith18542/vigil/certgen => ./certgen
	github.com/kdsmith18542/vigil/chaincfg/chainhash => ./chaincfg/chainhash
	github.com/kdsmith18542/vigil/chaincfg/v3 => ./chaincfg
	github.com/kdsmith18542/vigil/connmgr/v3 => ./connmgr
	github.com/kdsmith18542/vigil/container/apbf => ./container/apbf
	github.com/kdsmith18542/vigil/container/lru => ./container/lru
	github.com/kdsmith18542/vigil/crypto/blake256 => ./crypto/blake256
	github.com/kdsmith18542/vigil/crypto/rand => ./crypto/rand
	github.com/kdsmith18542/vigil/crypto/ripemd160 => ./crypto/ripemd160
	github.com/kdsmith18542/vigil/database/v3 => ./database
	github.com/kdsmith18542/vigil/dcrec => ./dcrec
	github.com/kdsmith18542/vigil/dcrec/secp256k1/v4 => ./dcrec/secp256k1
	github.com/kdsmith18542/vigil/dcrjson/v4 => ./dcrjson
	github.com/kdsmith18542/vigil/dcrutil/v4 => ./dcrutil
	github.com/kdsmith18542/vigil/gcs/v4 => ./gcs
	github.com/kdsmith18542/vigil/hdkeychain/v3 => ./hdkeychain
	github.com/kdsmith18542/vigil/limits => ./limits
	github.com/kdsmith18542/vigil/math/uint256 => ./math/uint256
	github.com/kdsmith18542/vigil/mixing => ./mixing
	github.com/kdsmith18542/vigil/peer/v3 => ./peer
	github.com/kdsmith18542/vigil/rpc/jsonrpc/types/v4 => ./rpc/jsonrpc/types
	github.com/kdsmith18542/vigil/rpcclient/v8 => ./rpcclient
	github.com/kdsmith18542/vigil/txscript/v4 => ./txscript
	github.com/kdsmith18542/vigil/wire => ./wire

)

replace github.com/kdsmith18542/vigil/kawpow => ./kawpow
