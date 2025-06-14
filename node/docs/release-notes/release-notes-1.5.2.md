# vgld v1.5.2

This is a patch release of vgld to address a potential denial of service vector.

## Changelog

This patch release consists of 5 commits from 2 contributors which total to 4 files changed, 114 additional lines of code, and 20 deleted lines of code.

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/release-v1.5.1...release-v1.5.2).

### Protocol and network:

- blockmanager: handle notfound messages from peers ([Vigil/vgld#2344](https://github.com/vigilnetwork/vgl/pull/2344))
- blockmanager: limit the requested maps ([Vigil/vgld#2344](https://github.com/vigilnetwork/vgl/pull/2344))
- server: increase ban score for notfound messages ([Vigil/vgld#2344](https://github.com/vigilnetwork/vgl/pull/2344))
- server: return whether addBanScore disconnected the peer ([Vigil/vgld#2344](https://github.com/vigilnetwork/vgl/pull/2344))

### Misc:

- release: Bump for 1.5.2([Vigil/vgld#2345](https://github.com/vigilnetwork/vgl/pull/2345))

### Code Contributors (alphabetical order):

- Dave Collins
- David Hill
