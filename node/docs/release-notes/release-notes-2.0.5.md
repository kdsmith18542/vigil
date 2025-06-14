# vgld v2.0.5 Release Notes

This is a patch release of vgld which includes the following changes:

- Improved StakeShuffle mixing robustness against misbehaving peers
- Peers are no longer intermittently disconnected when serving factored
  polynomial data

## Upgrade Required To Participate in StakeShuffle Mixing

Although upgrading to this latest release is not required for continued
operation of the core network, it is required for anyone who wishes to
participate in StakeShuffle mixing with the highest anonymity set guarantees and
fastest matching.

## Changelog

This patch release consists of 8 commits from 2 contributors which total to 8
files changed, 83 additional lines of code, and 56 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v2.0.4...release-v2.0.5).

### Protocol and network:

- [release-v2.0] peer: Add factored poly to stall detector ([Vigil/vgld#3457](https://github.com/vigilnetwork/vgl/pull/3457))

### Mixing message relay (mix pool):

- [release-v2.0] mixpool: Do not return early for revealed secrets ([Vigil/vgld#3458](https://github.com/vigilnetwork/vgl/pull/3458))

### Developer-related package and module changes:

- [release-v2.0] mixclient: Avoid jitter calculation panic ([Vigil/vgld#3458](https://github.com/vigilnetwork/vgl/pull/3458))
- [release-v2.0] mixclient: Detect exited csppsolver processes ([Vigil/vgld#3458](https://github.com/vigilnetwork/vgl/pull/3458))
- [release-v2.0] mixclient: Sort roots for slot assignment ([Vigil/vgld#3458](https://github.com/vigilnetwork/vgl/pull/3458))

### Developer-related module management:

- [release-v2.0] main: Use backported peer updates ([Vigil/vgld#3457](https://github.com/vigilnetwork/vgl/pull/3457))
- [release-v2.0] main: Use backported mixing updates ([Vigil/vgld#3458](https://github.com/vigilnetwork/vgl/pull/3458))

### Misc:

- [release-v2.0] release: Bump for 2.0.5 ([Vigil/vgld#3459](https://github.com/vigilnetwork/vgl/pull/3459))

### Code Contributors (alphabetical order):

- Dave Collins
- Josh Rickmar
