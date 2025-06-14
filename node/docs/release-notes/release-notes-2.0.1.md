# vgld v2.0.1 Release Notes

This is a patch release of vgld which includes the following key changes:

* Provides a new JSON-RPC API method named `getmixmessage` that can be used to
  query decentralized StakeShuffle mixing messages
* No longer relays mixing messages when transaction relay is disabled
* Transaction outputs with one confirmation may now be used as part of a mix
* Improves best network address candidate selection
* More consistent logging of banned peers along with the reason they were banned

## Changelog

This patch release consists of 19 commits from 3 contributors which total to 18
files changed, 388 additional lines of code, and 187 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v2.0.0...release-v2.0.1).

### Protocol and network:

- netsync: Request new headers on rejected tip ([Vigil/vgld#3315](https://github.com/vigilnetwork/vgl/pull/3315))
- vgld: Make DisableRelayTx also include mixing messages ([Vigil/vgld#3315](https://github.com/vigilnetwork/vgl/pull/3315))
- server: Add logs on why a peer is banned ([Vigil/vgld#3318](https://github.com/vigilnetwork/vgl/pull/3318))

### RPC:

- VGLjson,rpcserver,types: Add getmixmessage ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- rpcserver: Allow getmixmessage from limited users ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))

### Mixing message relay (mix pool):

- mixpool: Require 1 block confirmation ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))

### Documentation:

- docs: Add getmixmessage ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))

### Developer-related package and module changes:

- addrmgr: Give unattempted addresses a chance ([Vigil/vgld#3316](https://github.com/vigilnetwork/vgl/pull/3316))
- mixpool: Add missing mutex acquire ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- mixing: Use stdaddr.Hash160 instead of VGLutil ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- mixing: Reduce a couple of allocation cases ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- mixpool: Wrap non-bannable errors with RuleError ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- mixclient: Do not remove PRs from failed runs ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- mixclient: Improve error returning ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))
- server: Make peer banning synchronous ([Vigil/vgld#3318](https://github.com/vigilnetwork/vgl/pull/3318))
- server: Consolidate ban reason logging ([Vigil/vgld#3318](https://github.com/vigilnetwork/vgl/pull/3318))

### Developer-related module management:

- main: Use backported addrmgr updates ([Vigil/vgld#3316](https://github.com/vigilnetwork/vgl/pull/3316))
- main: Use backported module updates ([Vigil/vgld#3317](https://github.com/vigilnetwork/vgl/pull/3317))

### Misc:

- release: Bump for 2.0.1 ([Vigil/vgld#3319](https://github.com/vigilnetwork/vgl/pull/3319))

### Code Contributors (alphabetical order):

- Dave Collins
- David Hill
- Josh Rickmar
