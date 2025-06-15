# vgld v1.7.4

This is a patch release of vgld to support modifications to version 3 of the test
network as well as provide some minor improvements related to mining.

## Changelog

This patch release consists of 10 commits from 2 contributors which total to 17
files changed, 225 additional lines of code, and 57 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v1.7.2...release-v1.7.4) and
[here](https://github.com/vigilnetwork/vgl/compare/blockchain/v4.0.1...blockchain/v4.0.2).

### Protocol and network:

- blockchain: Enforce testnet difficulty throttling ([Vigil/vgld#2979](https://github.com/vigilnetwork/vgl/pull/2979))
- netsync: Improve sync height tracking ([Vigil/vgld#2985](https://github.com/vigilnetwork/vgl/pull/2985))

### Mining:

- mining: Copy regular txns for alternate templates ([Vigil/vgld#2985](https://github.com/vigilnetwork/vgl/pull/2985))
- server: Send winning tickets when unsynced mining ([Vigil/vgld#2985](https://github.com/vigilnetwork/vgl/pull/2985))

### RPC:

- rpcserver: Return template errors from getwork RPC ([Vigil/vgld#2985](https://github.com/vigilnetwork/vgl/pull/2985))

### Developer-related package and module changes:

- blockchain: Consistency pass for next req dif calc ([Vigil/vgld#2979](https://github.com/vigilnetwork/vgl/pull/2979))
- main: Use backported blockchain updates ([Vigil/vgld#2984](https://github.com/vigilnetwork/vgl/pull/2984))

### Testing and Quality Assurance:

- build: Update to latest actions and linter ([Vigil/vgld#2979](https://github.com/vigilnetwork/vgl/pull/2979))
- build: Use recommended golangci-lint installer  ([Vigil/vgld#2984](https://github.com/vigilnetwork/vgl/pull/2984))

### Misc:

- release: Bump for 1.7.4 ([Vigil/vgld#2986](https://github.com/vigilnetwork/vgl/pull/2986))

### Code Contributors (alphabetical order):

- Dave Collins
- Jamie Holdstock
