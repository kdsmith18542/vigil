# vgld v1.5.1

This is a patch release of vgld to address a minor memory leak with authenticated RPC websocket clients on intermittent connections.   It also updates the `vglctl` utility to include the new `auditreuse` vglwallet command.

## Changelog

This patch release consists of 4 commits from 3 contributors which total to 4 files changed, 27 additional lines of code, and 6 deleted lines of code.

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/release-v1.5.0...release-v1.5.1).

### RPC:

- rpcwebsocket: Remove client from missed maps ([Vigil/vgld#2049](https://github.com/vigilnetwork/vgl/pull/2049))
- rpcwebsocket: Use nonblocking messages and ntfns ([Vigil/vgld#2050](https://github.com/vigilnetwork/vgl/pull/2050))

### vglctl utility changes:

- vglctl: Update vglwallet RPC types package ([Vigil/vgld#2051](https://github.com/vigilnetwork/vgl/pull/2051))

### Misc:

- release: Bump for 1.5.1([Vigil/vgld#2052](https://github.com/vigilnetwork/vgl/pull/2052))

### Code Contributors (alphabetical order):

- Dave Collins
- Josh Rickmar
- Matheus Degiovani




