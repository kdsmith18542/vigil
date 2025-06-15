module github.com/kdsmith18542/vigil/wallet

go 1.23.0

require (
	github.com/kdsmith18542/vigil/addrmgr/v3 v3.0.0
	github.com/kdsmith18542/vigil/blockchain/v5 v5.0.0
	github.com/kdsmith18542/vigil/chaincfg/v3 v3.0.0
	github.com/kdsmith18542/vigil/dcrec/secp256k1/v4 v4.0.0
	github.com/kdsmith18542/vigil/dcrjson/v4 v4.0.0
	github.com/kdsmith18542/vigil/dcrutil/v4 v4.0.0
	github.com/kdsmith18542/vigil/rpcclient/v8 v8.0.0
	github.com/kdsmith18542/vigil/txscript/v4 v4.0.0
	github.com/kdsmith18542/vigil/wire v1.0.0
)

replace (
	github.com/kdsmith18542/vigil/addrmgr/v3 => ../node/addrmgr
	github.com/kdsmith18542/vigil/blockchain/v5 => ../node/blockchain
	github.com/kdsmith18542/vigil/chaincfg/v3 => ../node/chaincfg
	github.com/kdsmith18542/vigil/dcrec/secp256k1/v4 => ../node/dcrec/secp256k1
	github.com/kdsmith18542/vigil/dcrjson/v4 => ../node/dcrjson
	github.com/kdsmith18542/vigil/dcrutil/v4 => ../node/dcrutil
	github.com/kdsmith18542/vigil/rpcclient/v8 => ../node/rpcclient
	github.com/kdsmith18542/vigil/txscript/v4 => ../node/txscript
	github.com/kdsmith18542/vigil/wire => ../node/wire
)