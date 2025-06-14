# vgld v1.7.7

This is a patch release of vgld that includes the following changes:

- Use the latest network protocol version
- Reduce bandwidth usage in certain scenarios by avoiding requests for inventory that is already known
- Mitigate excessive CPU usage in some rare scenarios specific to the test network
- Improve best address candidate selection efficiency

## Changelog

This patch release consists of 19 commits from 3 contributors which total to 92
files changed, 1357 additional lines of code, and 1191 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v1.7.5...release-v1.7.7).

### Protocol and network:

- peer: Use latest pver by default ([Vigil/vgld#3083](https://github.com/vigilnetwork/vgl/pull/3083))
- peer: Correct known inventory check ([Vigil/vgld#3083](https://github.com/vigilnetwork/vgl/pull/3083))

### Documentation:

- peer: Go 1.19 doc comment formatting ([Vigil/vgld#3083](https://github.com/vigilnetwork/vgl/pull/3083))
- addrmgr: Go 1.19 doc comment formatting ([Vigil/vgld#3084](https://github.com/vigilnetwork/vgl/pull/3084))
- multi: Go 1.19 doc comment formatting ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))
- docs: Update README.md to required Go 1.19/1.20 ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))

### Developer-related package and module changes:

- peer: Support module graph prune and lazy load ([Vigil/vgld#3083](https://github.com/vigilnetwork/vgl/pull/3083))
- main: Use backported peer updates ([Vigil/vgld#3083](https://github.com/vigilnetwork/vgl/pull/3083))
- addmrgr: Use TempDir to create temp test dirs ([Vigil/vgld#3084](https://github.com/vigilnetwork/vgl/pull/3084))
- addrmgr: Support module graph prune and lazy load ([Vigil/vgld#3084](https://github.com/vigilnetwork/vgl/pull/3084))
- addrmgr: Break after selecting random address ([Vigil/vgld#3084](https://github.com/vigilnetwork/vgl/pull/3084))
- addrmgr: Set min value and optimize address chance ([Vigil/vgld#3084](https://github.com/vigilnetwork/vgl/pull/3084))
- main: Use backported addrmgr updates ([Vigil/vgld#3084](https://github.com/vigilnetwork/vgl/pull/3084))
- main: Update to use latest sys module ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))

### Testing and Quality Assurance:

- build: Enable run_vgl_tests.sh to work with go.work ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))
- build: Update to latest action versions ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))
- build: Update golangci-lint to v1.51.1 ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))
- build: Test against Go 1.20 ([Vigil/vgld#3087](https://github.com/vigilnetwork/vgl/pull/3087))

### Misc:

- release: Bump for 1.7.7 ([Vigil/vgld#3085](https://github.com/vigilnetwork/vgl/pull/3085))

### Code Contributors (alphabetical order):

- Dave Collins
- Eng Zer Jun
- Jonathan Chappelow
