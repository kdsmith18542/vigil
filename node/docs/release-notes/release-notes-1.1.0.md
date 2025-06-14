# vgld v1.1.0

This release of vgld primarily introduces a new consensus vote agenda which
allows the stakeholders to decide whether or not to activate the features needed
for providing full support for Lightning Network.  For those unfamiliar with the
voting process in Vigil, this means that all code in order to support these
features is already included in this release, however its enforcement will
remain dormant until the stakeholders vote to activate it.

The following Vigil Change Proposals (VGLPs) describe the proposed changes in detail:
- [VGLP0002](https://github.com/Vigil/VGLPs/blob/master/VGLP-0002/VGLP-0002.mediawiki)
- [VGLP0003](https://github.com/Vigil/VGLPs/blob/master/VGLP-0003/VGLP-0003.mediawiki)

It is important for everyone to upgrade their software to this latest release
even if you don't intend to vote in favor of the agenda.

## Notable Changes

### Lightning Network Features Vote

In order to fully support many of the benefits that the Lightning Network will
bring, there are some primitives that involve changes to the current consensus
that need to be enabled.  A new vote with the id `lnfeatures` is now available
as of this release.  After upgrading, stakeholders may set their preferences
through their wallet or stake pool's website.

### Transaction Finality Policy

The standard policy for transaction relay has been changed to use the median
time of the past several blocks instead of the current network adjusted time
when examining lock times to determine if a transaction is final.  This provides
a more deterministic check across all peers and prevents the possibility of
miners attempting to game the timestamps in order to include more transactions.

Consensus enforcement of this change relies on the result of the aforementioned
`lnfeatures` vote.

### Relative Time Locks Policy

The standard policy for transaction relay has been modified to enforce relative
lock times for version 2 transactions via their sequence numbers and a new
`OP_CHECKSEQUENCEVERIFY` opcode.

Consensus enforcement of this change relies on the result of the aforementioned
`lnfeatures` vote.

### OP_SHA256 Opcode

In order to better support cross-chain interoperability, a new opcode to compute
the SHA-256 hash is being proposed.  Since this opcode is implemented as a hard
fork, it will not be available for use in scripts unless the aforementioned
`lnfeatures` vote passes.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/v1.0.7...v1.1.0).

### Protocol and network:
- chaincfg: update checkpoints for 1.1.0 release [Vigil/vgld#850](https://github.com/vigilnetwork/vgl/pull/850)
- chaincfg: Introduce agenda for v5 lnfeatures vote [Vigil/vgld#848](https://github.com/vigilnetwork/vgl/pull/848)
- txscript: Introduce OP_SHA256 [Vigil/vgld#851](https://github.com/vigilnetwork/vgl/pull/851)
- wire: Decrease num allocs when decoding headers [Vigil/vgld#861](https://github.com/vigilnetwork/vgl/pull/861)
- blockchain: Implement enforced relative seq locks [Vigil/vgld#864](https://github.com/vigilnetwork/vgl/pull/864)
- txscript: Implement CheckSequenceVerify [Vigil/vgld#864](https://github.com/vigilnetwork/vgl/pull/864)
- multi: Enable vote for VGLP0002 and VGLP0003 [Vigil/vgld#855](https://github.com/vigilnetwork/vgl/pull/855)

### Transaction relay (memory pool):
- mempool: Use median time for tx finality checks [Vigil/vgld#860](https://github.com/vigilnetwork/vgl/pull/860)
- mempool: Enforce relative sequence locks [Vigil/vgld#864](https://github.com/vigilnetwork/vgl/pull/864)
- policy/mempool: Enforce CheckSequenceVerify opcode [Vigil/vgld#864](https://github.com/vigilnetwork/vgl/pull/864)

### RPC:
- rpcserver: check whether ticketUtx was found [Vigil/vgld#824](https://github.com/vigilnetwork/vgl/pull/824)
- rpcserver: return rule error on rejected raw tx [Vigil/vgld#808](https://github.com/vigilnetwork/vgl/pull/808)

### vgld command-line flags:
- config: Extend --profile cmd line option to allow interface to be specified [Vigil/vgld#838](https://github.com/vigilnetwork/vgl/pull/838)

### Documentation
- docs: rpcapi format update [Vigil/vgld#807](https://github.com/vigilnetwork/vgl/pull/807)
- config: export sampleconfig for use by VGLinstall [Vigil/vgld#834](https://github.com/vigilnetwork/vgl/pull/834)
- sampleconfig: Add package README and doc.go [Vigil/vgld#835](https://github.com/vigilnetwork/vgl/pull/835)
- docs: create entry for getstakeversions in rpcapi [Vigil/vgld#819](https://github.com/vigilnetwork/vgl/pull/819)
- docs: crosscheck and update all rpc doc entries [Vigil/vgld#847](https://github.com/vigilnetwork/vgl/pull/847)
- docs: update git commit messages section heading [Vigil/vgld#863](https://github.com/vigilnetwork/vgl/pull/863)

### Developer-related package changes:
- Fix and regenerate precomputed secp256k1 curve [Vigil/vgld#823](https://github.com/vigilnetwork/vgl/pull/823)
- VGLec: use hardcoded datasets in tests [Vigil/vgld#822](https://github.com/vigilnetwork/vgl/pull/822)
- Use dchest/blake256  [Vigil/vgld#827](https://github.com/vigilnetwork/vgl/pull/827)
- glide: use jessevdk/go-flags for consistency [Vigil/vgld#833](https://github.com/vigilnetwork/vgl/pull/833)
- multi: Error descriptions are in lower case [Vigil/vgld#842](https://github.com/vigilnetwork/vgl/pull/842)
- txscript: Rename OP_SHA256 to OP_BLAKE256 [Vigil/vgld#840](https://github.com/vigilnetwork/vgl/pull/840)
- multi: Abstract standard verification flags [Vigil/vgld#852](https://github.com/vigilnetwork/vgl/pull/852)
- chain: Remove memory block node pruning [Vigil/vgld#858](https://github.com/vigilnetwork/vgl/pull/858)
- txscript: Add API to parse atomic swap contracts [Vigil/vgld#862](https://github.com/vigilnetwork/vgl/pull/862)

### Testing and Quality Assurance:
- Test against go 1.9 [Vigil/vgld#836](https://github.com/vigilnetwork/vgl/pull/836)
- VGLec: remove testify dependency [Vigil/vgld#829](https://github.com/vigilnetwork/vgl/pull/829)
- mining_test: add edge conditions from btcd [Vigil/vgld#831](https://github.com/vigilnetwork/vgl/pull/831)
- stake: Modify ticket tests to use chaincfg params [Vigil/vgld#844](https://github.com/vigilnetwork/vgl/pull/844)
- blockchain: Modify tests to use chaincfg params [Vigil/vgld#845](https://github.com/vigilnetwork/vgl/pull/845)
- blockchain: Cleanup various tests [Vigil/vgld#843](https://github.com/vigilnetwork/vgl/pull/843)
- Ensure run_vgl_tests.sh local fails correctly when gometalinter errors [Vigil/vgld#846](https://github.com/vigilnetwork/vgl/pull/846)
- peer: fix logic race in peer connection test [Vigil/vgld#865](https://github.com/vigilnetwork/vgl/pull/865)

### Misc:
- glide: sync deps [Vigil/vgld#837](https://github.com/vigilnetwork/vgl/pull/837)
- Update Vigil deps for v1.1.0 [Vigil/vgld#868](https://github.com/vigilnetwork/vgl/pull/868)
- Bump for v1.1.0 [Vigil/vgld#867](https://github.com/vigilnetwork/vgl/pull/867)

### Code Contributors (alphabetical order):

- Alex Yocom-Piatt
- Dave Collins
- David Hill
- Donald Adu-Poku
- Jason Zavaglia
- Jean-Christophe Mincke
- Jolan Luff
- Josh Rickmar
