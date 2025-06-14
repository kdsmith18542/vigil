module github.com/vigilnetwork/vgl/test/genesis_test

go 1.21

require (
	github.com/vigilnetwork/vgl/chaincfg/chainhash v1.0.4
	github.com/vigilnetwork/vgl/chaincfg/v3 v3.2.1
	github.com/vigilnetwork/vgl/wire v1.6.0
	
)

replace (
	github.com/vigilnetwork/vgl => ../../
	github.com/vigilnetwork/vgl/chaincfg/chainhash => ../../chaincfg/chainhash
	github.com/vigilnetwork/vgl/chaincfg/v3 => ../../chaincfg/v3
	github.com/vigilnetwork/vgl/wire => ../../wire
	
)
