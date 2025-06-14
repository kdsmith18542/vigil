module github.com/kdsmith18542/vigil/test/genesis_test

go 1.21

require (
	github.com/kdsmith18542/vigil/chaincfg/chainhash v1.0.4
	github.com/kdsmith18542/vigil/chaincfg/v3 v3.2.1
	github.com/kdsmith18542/vigil/wire v1.6.0
	
)

replace (
	github.com/kdsmith18542/vigil => ../../
	github.com/kdsmith18542/vigil/chaincfg/chainhash => ../../chaincfg/chainhash
	github.com/kdsmith18542/vigil/chaincfg/v3 => ../../chaincfg/v3
	github.com/kdsmith18542/vigil/wire => ../../wire
	
)
