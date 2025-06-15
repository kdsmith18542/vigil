# vgld v2.0.3 Release Notes

This is a patch release of vgld which includes the following changes:

- Improved sender privacy for transactions and mix messages via randomized
  announcements
- Nodes now prefer to maintain at least three mixing-capable outbound connections
- Recent transactions and mix messages will now be available to serve for longer
- Reduced memory usage during periods of lower activity
- Mixing-related performance enhancements

## Changelog

This patch release consists of 26 commits from 2 contributors which total to 37
files changed, 4527 additional lines of code, and 499 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v2.0.2...release-v2.0.3).

### Protocol and network:

- [release-v2.0] peer: Randomize inv batching delays ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] server: Cache advertised txns ([Vigil/vgld#3392](https://github.com/vigilnetwork/vgl/pull/3392))
- [release-v2.0] server: Prefer 3 min mix capable outbound peers ([Vigil/vgld#3392](https://github.com/vigilnetwork/vgl/pull/3392))

### Mixing message relay (mix pool):

- [release-v2.0] mixpool: Cache recently removed msgs ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] mixclient: Introduce random message jitter ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] netsync: Remove spent PRs from tip block txns ([Vigil/vgld#3392](https://github.com/vigilnetwork/vgl/pull/3392))

### Documentation:

- [release-v2.0] docs: Update for container/lru module ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] rand: Add README.md ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))

### Developer-related package and module changes:

- [release-v2.0] container/lru: Implement type safe generic LRUs ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] peer: Use container/lru module ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] crypto/rand: Implement module ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] rand: Add BigInt ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] rand: Uppercase N in half-open-range funcs ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] rand: Add rand.N generic func ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] rand: Add ShuffleSlice ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] rand: Add benchmarks ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] peer: Use new crypto/rand module ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] mixing: Prealloc buffers ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] mixpool: Remove Receive expectedMessages argument ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] mixing: Use new crypto/rand module ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] mixpool: Remove run from conflicting msg err ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] mixpool: Remove more references to runs ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))
- [release-v2.0] mixing: Reduce slot reservation mix pads allocs ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))

### Developer-related module management:

- [release-v2.0] main: Use backported peer updates ([Vigil/vgld#3390](https://github.com/vigilnetwork/vgl/pull/3390))
- [release-v2.0] main: Use backported mixing updates ([Vigil/vgld#3391](https://github.com/vigilnetwork/vgl/pull/3391))

### Misc:

- [release-v2.0] release: Bump for 2.0.3 ([Vigil/vgld#3393](https://github.com/vigilnetwork/vgl/pull/3393))

### Code Contributors (alphabetical order):

- Dave Collins
- Josh Rickmar
