# vgld v1.7.0

This is a new major release of vgld.  Some of the key highlights are:

* Four new consensus vote agendas which allow stakeholders to decide whether or
  not to activate support for the following:
  * Reverting the Treasury maximum expenditure policy
  * Enforcing explicit version upgrades
  * Support for automatic ticket revocations for missed votes
  * Changing the Proof-of-Work and Proof-of-Stake subsidy split from 60%/30% to 10%/80%
* Substantially reduced initial sync time
* Major performance enhancements to unspent transaction output handling
* Faster cryptographic signature validation
* Significant improvements to network synchronization
* Support for a configurable assumed valid block
* Block index memory usage reduction
* Asynchronous indexing
* Version 1 block filters removal
* Various updates to the RPC server:
  * Additional per-connection read limits
  * A more strict cross origin request policy
  * A new alternative client authentication mechanism based on TLS certificates
  * Availability of the scripting language version for transaction outputs
  * Several other notable updates, additions, and removals related to the JSON-RPC API
* New developer modules:
  * Age-Partitioned Bloom Filters
  * Fixed-Precision Unsigned 256-bit Integers
  * Standard Scripts
  * Standard Addresses
* Infrastructure improvements
* Quality assurance changes

For those unfamiliar with the
[voting process](https://docs.vigil.network/governance/consensus-rule-voting/overview/)
in Vigil, all code needed in order to support each of the aforementioned
consensus changes is already included in this release, however it will remain
dormant until the stakeholders vote to activate it.

For reference, the consensus change work for each of the four changes was
originally proposed and approved for initial implementation via the following
Vigiliteia proposals:
- [Decentralized Treasury Spending](https://proposals-archive.vigil.network/proposals/c96290a)
- [Explicit Version Upgrades Consensus Change](https://proposals.vigil.network/record/3a98861)
- [Automatic Ticket Revocations Consensus Change](https://proposals.vigil.network/record/e2d7b7d)
- [Change PoW/PoS Subsidy Split From 60/30 to 10/80](https://proposals.vigil.network/record/427e1d4)

The following Vigil Change Proposals (VGLPs) describe the proposed changes in
detail and provide full technical specifications:
- [VGLP0007](https://github.com/Vigil/VGLPs/blob/master/VGLP-0007/VGLP-0007.mediawiki)
- [VGLP0008](https://github.com/Vigil/VGLPs/blob/master/VGLP-0008/VGLP-0008.mediawiki)
- [VGLP0009](https://github.com/Vigil/VGLPs/blob/master/VGLP-0009/VGLP-0009.mediawiki)
- [VGLP0010](https://github.com/Vigil/VGLPs/blob/master/VGLP-0010/VGLP-0010.mediawiki)

## Upgrade Required

**It is extremely important for everyone to upgrade their software to this
latest release even if you don't intend to vote in favor of the agenda.  This
particularly applies to PoW miners as failure to upgrade will result in lost
rewards after block height 635775.  That is estimated to be around Feb 21st,
2022.**

## Downgrade Warning

The database format in v1.7.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

The database migration typically takes around 40-50 minutes on HDDs and 20-30
minutes on SSDs.

## Notable Changes

### Four New Consensus Change Votes

Four new consensus change votes are now available as of this release.  After
upgrading, stakeholders may set their preferences through their wallet.

#### Revert Treasury Maximum Expenditure Policy Vote

The first new vote available as of this release has the id `reverttreasurypolicy`.

The primary goal of this change is to revert the currently active maximum
expenditure policy of the decentralized Treasury to the one specified in the
[original Vigiliteia proposal](https://proposals-archive.vigil.network/proposals/c96290a).

See [VGLP0007](https://github.com/Vigil/VGLPs/blob/master/VGLP-0007/VGLP-0007.mediawiki) for
the full technical specification.

#### Explicit Version Upgrades Vote

The second new vote available as of this release has the id `explicitverupgrades`.

The primary goals of this change are to:

* Provide an easy, reliable, and efficient method for software and hardware to
  determine exactly which rules should be applied to transaction and script
  versions
* Further embrace the increased security and other desirable properties that
  hard forks provide over soft forks

See the following for more details:

* [Vigiliteia proposal](https://proposals.vigil.network/record/3a98861)
* [VGLP0008](https://github.com/Vigil/VGLPs/blob/master/VGLP-0008/VGLP-0008.mediawiki)

#### Automatic Ticket Revocations Vote

The third new vote available as of this release has the id `autorevocations`.

The primary goals of this change are to:

* Improve the Vigil stakeholder user experience by removing the requirement for
  stakeholders to manually revoke missed and expired tickets
* Enable the recovery of funds for users who lost their redeem script for the
  legacy VSP system (before the release of vspd, which removed the need for the
  redeem script)

See the following for more details:

* [Vigiliteia proposal](https://proposals.vigil.network/record/e2d7b7d)
* [VGLP0009](https://github.com/Vigil/VGLPs/blob/master/VGLP-0009/VGLP-0009.mediawiki)

#### Change PoW/PoS Subsidy Split to 10/80 Vote

The fourth new vote available as of this release has the id `changesubsidysplit`.

The proposed modification to the subsidy split is intended to substantially
diminish the ability to attack Vigil's markets with mined coins and improve
decentralization of the issuance process.

See the following for more details:

* [Vigiliteia proposal](https://proposals.vigil.network/record/427e1d4)
* [VGLP0010](https://github.com/Vigil/VGLPs/blob/master/VGLP-0010/VGLP-0010.mediawiki)

### Substantially Reduced Initial Sync Time

The amount of time it takes to complete the initial chain synchronization
process has been substantially reduced.  With default settings, it is around 48%
faster versus the previous release.

### Unspent Transaction Output Overhaul

The way unspent transaction outputs (UTXOs) are handled has been significantly
reworked to provide major performance enhancements to both steady-state
operation as well as the initial chain sync process as follows:

* Each UTXO is now tracked independently on a per-output basis
* The UTXOs now reside in a dedicated database
* All UTXO reads and writes now make use of a cache

#### Unspent Transaction Output Cache

All reads and writes of unspent transaction outputs (utxos) now go through a
cache that sits on top of the utxo set database which drastically reduces the
amount of reading and writing to disk, especially during the initial sync
process when a very large number of blocks are being processed in quick
succession.

This utxo cache provides significant runtime performance benefits at the cost of
some additional memory usage.  The maximum size of the cache can be configured
with the new `--utxocachemaxsize` command-line configuration option.  The
default value is 150 MiB, the minimum value is 25 MiB, and the maximum value is
32768 MiB (32 GiB).

Some key properties of the cache are as follows:

* For reads, the UTXO cache acts as a read-through cache
  * All UTXO reads go through the cache
  * Cache misses load the missing data from the disk and cache it for future lookups
* For writes, the UTXO cache acts as a write-back cache
  * Writes to the cache are acknowledged by the cache immediately, but are only
    periodically flushed to disk
* Allows intermediate steps to effectively be skipped thereby avoiding the need
  to write millions of entries to disk
* On average, recent UTXOs are much more likely to be spent in upcoming blocks
  than older UTXOs, so only the oldest UTXOs are evicted as needed in order to
  maximize the hit ratio of the cache
* The cache is periodically flushed with conditional eviction:
  * When the cache is NOT full, nothing is evicted, but the changes are still
    written to the disk set to allow for a quicker reconciliation in the case of
    an unclean shutdown
  * When the cache is full, 15% of the oldest UTXOs are evicted

### Faster Cryptographic Signature Validation

Some aspects of the underlying crypto code has been updated to further improve
its execution speed and reduce the number of memory allocations resulting in
about a 1% reduction to signature verification time.

The primary benefits are:

* Improved vote times since blocks and transactions propagate more quickly
  throughout the network
* Approximately a 2% reduction to the duration of the initial sync process

### Significant Improvements to Network Synchronization

The method used to obtain blocks from other peers on the network is now guided
entirely by block headers.  This provides a wide variety of benefits, but the
most notable ones for most users are:

* Faster initial synchronization
* Reduced bandwidth usage
* Enhanced protection against attempted DoS attacks
* Percentage-based progress reporting
* Improved steady state logging

### Support for Configurable Assumed Valid Block

This release introduces a new model for deciding when several historical
validation checks may be skipped for blocks that are an ancestor of a known good
block.

Specifically, a new `AssumeValid` parameter is now used to specify the
aforementioned known good block.  The default value of the parameter is updated
with each release to a recent block that is part of the main chain.

The default value of the parameter can be overridden with the `--assumevalid`
command-line option by setting it as follows:

* `--assumevalid=0`: Disable the feature resulting in no skipped validation checks
* `--assumevalid=[blockhash]`:  Set `AssumeValid` to the specified block hash

Specifying a block hash closer to the current best chain tip allows for faster
syncing.  This is useful since the validation requirements increase the longer a
particular release build is out as the default known good block becomes deeper
in the chain.

### Block Index Memory Usage Reduction

The block index that keeps track of block status and connectivity now occupies
around 30MiB less memory and scales better as more blocks are added to the
chain.

### Asynchronous Indexing

The various optional indexes are now created asynchronously versus when
blocks are processed as was previously the case.

This permits blocks to be validated more quickly when the indexes are enabled
since the validation no longer needs to wait for the indexing operations to
complete.

In order to help keep consistent behavior for RPC clients, RPCs that involve
interacting with the indexes will not return results until the associated
indexing operation completes when the indexing tip is close to the current best
chain tip.

One side effect of this change that RPC clients should be aware of is that it is
now possible to receive sync timeout errors on RPCs that involve interacting
with the indexes if the associated indexing tip gets so far behind it would end
up delaying results for too long.  In practice, errors of this type are rare and
should only ever be observed during the initial sync process before the
associated indexes are current.  However, callers should be aware of the
possibility and handle the error accordingly.

The following RPCs are affected:

* `existsaddress`
* `existsaddresses`
* `getrawtransaction`
* `searchrawtransactions`

### Version 1 Block Filters Removal

The previously deprecated version 1 block filters are no longer available on the
peer-to-peer network.  Use
[version 2 block filters](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#version-2-block-filters)
with their associated
[block header commitment](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#block-header-commitments)
and [inclusion proof](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#verifying-commitment-root-inclusion-proofs)
instead.

### RPC Server Changes

The RPC server version as of this release is 7.0.0.

#### Max Request Limits

The RPC server now imposes the following additional per-connection read limits
to help further harden it against potential abuse in non-standard configurations
on poorly-configured networks:

* 0 B / 8 MiB for pre and post auth HTTP connections
* 4 KiB / 16 MiB for pre and post auth WebSocket connections

In practice, these changes will not have any noticeable effect for the vast
majority of nodes since the RPC server is not publicly accessible by default and
also requires authentication.

Nevertheless, it can still be useful for scenarios such as authenticated fuzz
testing and improperly-configured networks that have disabled all other security
measures.

#### More Strict Cross Origin Request (CORS) Policy

The CORS policy for WebSocket clients is now more strict and rejects requests
from other domains.

In practice, CORS requests will be rejected before ever reaching that point due
to the use of a self-signed TLS certificate and the requirement for
authentication to issue any commands.  However, additional protection mechanisms
make it that much more difficult to attack by providing defense in depth.

#### Alternative Client Authentication Method Based on TLS Certificates

A new alternative method for TLS clients to authenticate to the RPC server by
presenting a client certificate in the TLS handshake is now available.

Under this authentication method, the certificate authority for a client
certificate must be added to the RPC server as a trusted root in order for it to
trust the client.  Once activated, clients will no longer be required to provide
HTTP Basic authentication nor use the `authenticate` RPC in the case of
WebSocket clients.

Note that while TLS client authentication has the potential to ultimately allow
more fine grained access controls on a per-client basis, it currently only
supports clients with full administrative privileges.  In other words, it is not
currently compatible with the `--rpclimituser` and `--rpclimitpass` mechanism,
so users depending on the limited user settings should avoid the new
authentication method for now.

The new authentication type can be activated with the `--authtype=clientcert`
configuration option.

By default, the trusted roots are loaded from the `clients.pem` file in vgld's
application data directory, however, that location can be modified via the
`--clientcafile` option if desired.

#### Updates to Transaction Output Query RPC (`gettxout`)

The `gettxout` RPC has the following modifications:

* An additional `tree` parameter is now required in order to explicitly identify
  the exact transaction output being requested
* The transaction `version` field is no longer available in the primary JSON
  object of the results
* The child `scriptPubKey` JSON object in the results now includes a new
  `version` field that identifies the scripting language version

See the
[gettxout JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#gettxout)
for API details.

#### Removal of Stake Difficulty Notification RPCs (`notifystakedifficulty` and `stakedifficulty`)

The deprecated `notifystakedifficulty` and `stakedifficulty` WebSocket-only RPCs
are no longer available.  This notification is unnecessary since the difficulty
change interval is well defined.  Callers may obtain the difficulty via
`getstakedifficulty` at the appropriate difficulty change intervals instead.

See the
[getstakedifficulty JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getstakedifficulty)
for API details.

#### Removal of Version 1 Filter RPCs (`getcfilter` and `getcfilterheader`)

The deprecated `getcfilter` and `getcfilterheader` RPCs, which were previously
used to obtain version 1 block filters via RPC are no longer available. Use
`getcfilterv2` instead.

See the
[getcfilterv2 JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getcfilterv2)
for API details.

#### New Median Time Field on Block Query RPCs (`getblock` and `getblockheader`)

The verbose results of the `getblock` and `getblockheader` RPCs now include a
`mediantime` field that specifies the median block time associated with the
block.

See the following for API details:

* [getblock JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblock)
* [getblockheader JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblockheader)

#### New Scripting Language Version Field on Raw Transaction RPCs (`getrawtransaction`, `decoderawtransaction`, `searchrawtransactions`, and `getblock`)

The verbose results of the `getrawtransaction`, `decoderawtransaction`,
`searchrawtransactions`, and `getblock` RPCs now include a `version` field in
the child `scriptPubKey` JSON object that identifies the scripting language
version.

See the following for API details:

* [getrawtransaction JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getrawtransaction)
* [decoderawtransaction JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#decoderawtransaction)
* [searchrawtransactions JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#searchrawtransactions)
* [getblock JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblock)

#### New Treasury Add Transaction Filter on Mempool Query RPC (`getrawmempool`)

The transaction type parameter of the `getrawmempool` RPC now accepts `tadd` to
only include treasury add transactions in the results.

See the
[getrawmempool JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getrawmempool)
for API details.

#### New Manual Block Invalidation and Reconsideration RPCs (`invalidateblock` and `reconsiderblock`)

A new pair of RPCs named `invalidateblock` and `reconsiderblock` are now
available.  These RPCs can be used to manually invalidate a block as if it had
violated consensus rules and reconsider a block for validation and best chain
selection by removing any invalid status from it and its ancestors, respectively.

This capability is provided for development, testing, and debugging.  It can be
particularly useful when developing services that build on top of Vigil to more
easily ensure edge conditions associated with invalid blocks and chain
reorganization are being handled properly.

These RPCs do not apply to regular users and can safely be ignored outside of
development.

See the following for API details:

* [invalidateblock JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#invalidateblock)
* [reconsiderblock JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#reconsiderblock)

### Reject Protocol Message Deprecated (`reject`)

The `reject` peer-to-peer protocol message is now deprecated and is scheduled to
be removed in a future release.

This message is a holdover from the original codebase where it was required, but
it really is not a useful message and has several downsides:

* Nodes on the network must be trustless, which means anything relying on such a
  message is setting itself up for failure because nodes are not obligated to
  send it at all nor be truthful as to the reason
* It can be harmful to privacy as it allows additional node fingerprinting
* It can lead to security issues for implementations that don't handle it with
  proper sanitization practices
* It can easily give software implementations the fully incorrect impression
  that it can be relied on for determining if transactions and blocks are valid
* The only way it is actually used currently is to show a debug log message,
  however, all of that information is already available via the peer and/or wire
  logging anyway
* It carries a non-trivial amount of development overhead to continue to support
  it when nothing actually uses it

### No DNS Seeds Command-Line Option Deprecated (`--nodnsseed`)

The `--nodnsseed` command-line configuration option is now deprecated and will
be removed in a future release.  Use `--noseeders` instead.

DNS seeding has not been used since the previous release.

## Notable New Developer Modules

### Age-Partitioned Bloom Filters

A new `github.com/vigilnetwork/vgl/container/apbf` module is now available that
provides Age-Partitioned Bloom Filters (APBFs).

An APBF is a probabilistic lookup device that can quickly determine if it
contains an element.  It permits tracking large amounts of data while using very
little memory at the cost of a controlled rate of false positives.  Unlike
classic Bloom filters, it is able to handle an unbounded amount of data by aging
and discarding old items.

For a concrete example of actual savings achieved in Vigil by making use of an
APBF, the memory to track addresses known by 125 peers was reduced from ~200 MiB
to ~5 MiB.

See the
[apbf module documentation](https://pkg.go.dev/github.com/vigilnetwork/vgl/container/apbf)
for full details on usage, accuracy under workloads, expected memory usage, and
performance benchmarks.

### Fixed-Precision Unsigned 256-bit Integers

A new `github.com/vigilnetwork/vgl/math/uint256` module is now available that provides
highly optimized allocation free fixed precision unsigned 256-bit integer
arithmetic.

The package has a strong focus on performance and correctness and features
arithmetic, boolean comparison, bitwise logic, bitwise shifts, conversion
to/from relevant types, and full formatting support - all served with an
ergonomic API, full test coverage, and benchmarks.

Every operation is faster than the standard library `big.Int` equivalent and the
primary math operations provide reductions of over 90% in the calculation time.
Most other operations are also significantly faster.

See the
[uint256 module documentation](https://pkg.go.dev/github.com/vigilnetwork/vgl/math/uint256)
for full details on usage, including a categorized summary, and performance
benchmarks.

### Standard Scripts

A new `github.com/vigilnetwork/vgl/txscript/v4/stdscript` package is now available
that provides facilities for identifying and extracting data from transaction
scripts that are considered standard by the default policy of most nodes.

The package is part of the `github.com/vigilnetwork/vgl/txscript/v4` module.

See the
[stdscript package documentation](https://pkg.go.dev/github.com/vigilnetwork/vgl/txscript/v4/stdscript)
for full details on usage and a list of the recognized standard scripts.

### Standard Addresses

A new `github.com/vigilnetwork/vgl/txscript/v4/stdaddr` package is now available that
provides facilities for working with human-readable Vigil payment addresses.

The package is part of the `github.com/vigilnetwork/vgl/txscript/v4` module.

See the
[stdaddr package documentation](https://pkg.go.dev/github.com/vigilnetwork/vgl/txscript/v4/stdaddr)
for full details on usage and a list of the supported addresses.

## Changelog

This release consists of 877 commits from 16 contributors which total to 492
files changed, 77937 additional lines of code, and 30961 deleted lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v1.6.0...release-v1.7.0).

### Protocol and network:

- chaincfg: Add extra seeders ([Vigil/vgld#2532](https://github.com/vigilnetwork/vgl/pull/2532))
- server: Stop serving v1 cfilters over p2p ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- blockchain: Decouple processing and download logic ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- blockchain: Improve current detection ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- netsync: Rework inventory announcement handling ([Vigil/vgld#2548](https://github.com/vigilnetwork/vgl/pull/2548))
- peer: Add inv type summary to debug message ([Vigil/vgld#2556](https://github.com/vigilnetwork/vgl/pull/2556))
- netsync: Remove unused submit block flags param ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- netsync: Remove submit/processblock orphan flag ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- netsync: Remove orphan block handling ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- netsync: Rework sync model to use hdr annoucements ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- progresslog: Add support for header sync progress ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- netsync: Add header sync progress log ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- multi: Add chain verify progress percentage ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- peer: Remove getheaders response deadline ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- chaincfg: Update seed URL ([Vigil/vgld#2564](https://github.com/vigilnetwork/vgl/pull/2564))
- upnp: Don't return loopback IPs in getOurIP ([Vigil/vgld#2566](https://github.com/vigilnetwork/vgl/pull/2566))
- server: Prevent duplicate pending conns ([Vigil/vgld#2563](https://github.com/vigilnetwork/vgl/pull/2563))
- multi: Use an APBF for recently confirmed txns ([Vigil/vgld#2580](https://github.com/vigilnetwork/vgl/pull/2580))
- multi: Use an APBF for per peer known addrs ([Vigil/vgld#2583](https://github.com/vigilnetwork/vgl/pull/2583))
- peer: Stop sending and logging reject messages ([Vigil/vgld#2586](https://github.com/vigilnetwork/vgl/pull/2586))
- netsync: Stop sending reject messages ([Vigil/vgld#2586](https://github.com/vigilnetwork/vgl/pull/2586))
- server: Stop sending reject messages ([Vigil/vgld#2586](https://github.com/vigilnetwork/vgl/pull/2586))
- peer: Remove deprecated onversion reject return ([Vigil/vgld#2586](https://github.com/vigilnetwork/vgl/pull/2586))
- peer: Remove unneeded PushRejectMsg ([Vigil/vgld#2586](https://github.com/vigilnetwork/vgl/pull/2586))
- wire: Deprecate reject message ([Vigil/vgld#2586](https://github.com/vigilnetwork/vgl/pull/2586))
- server: Respond to getheaders when same chain tip ([Vigil/vgld#2587](https://github.com/vigilnetwork/vgl/pull/2587))
- netsync: Use an APBF for recently rejected txns ([Vigil/vgld#2590](https://github.com/vigilnetwork/vgl/pull/2590))
- server: Only send fast block anns to full nodes ([Vigil/vgld#2606](https://github.com/vigilnetwork/vgl/pull/2606))
- upnp: More accurate getOurIP ([Vigil/vgld#2571](https://github.com/vigilnetwork/vgl/pull/2571))
- server: Correct tx not found ban reason ([Vigil/vgld#2677](https://github.com/vigilnetwork/vgl/pull/2677))
- chaincfg: Add VGLP0007 deployment ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- chaincfg: Introduce explicit ver upgrades agenda ([Vigil/vgld#2713](https://github.com/vigilnetwork/vgl/pull/2713))
- blockchain: Implement reject new tx vers vote ([Vigil/vgld#2716](https://github.com/vigilnetwork/vgl/pull/2716))
- blockchain: Implement reject new script vers vote ([Vigil/vgld#2716](https://github.com/vigilnetwork/vgl/pull/2716))
- chaincfg: Add agenda for auto ticket revocations ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- multi: VGLP0009 Auto revocations consensus change ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- chaincfg: Use single latest checkpoint ([Vigil/vgld#2762](https://github.com/vigilnetwork/vgl/pull/2762))
- peer: Offset ping interval from idle timeout ([Vigil/vgld#2796](https://github.com/vigilnetwork/vgl/pull/2796))
- chaincfg: Update checkpoint for upcoming release ([Vigil/vgld#2794](https://github.com/vigilnetwork/vgl/pull/2794))
- chaincfg: Update min known chain work for release ([Vigil/vgld#2795](https://github.com/vigilnetwork/vgl/pull/2795))
- netsync: Request init state immediately upon sync ([Vigil/vgld#2812](https://github.com/vigilnetwork/vgl/pull/2812))
- blockchain: Reject old block vers for HFV ([Vigil/vgld#2752](https://github.com/vigilnetwork/vgl/pull/2752))
- netsync: Rework next block download logic ([Vigil/vgld#2828](https://github.com/vigilnetwork/vgl/pull/2828))
- chaincfg: Add AssumeValid param ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- chaincfg: Introduce subsidy split change agenda ([Vigil/vgld#2847](https://github.com/vigilnetwork/vgl/pull/2847))
- multi: Implement VGLP0010 subsidy consensus vote ([Vigil/vgld#2848](https://github.com/vigilnetwork/vgl/pull/2848))
- server: Force PoW upgrade to v9 ([Vigil/vgld#2875](https://github.com/vigilnetwork/vgl/pull/2875))

### Transaction relay (memory pool):

- mempool: Limit ancestor tracking in mempool ([Vigil/vgld#2458](https://github.com/vigilnetwork/vgl/pull/2458))
- mempool: Remove old fix sequence lock rejection ([Vigil/vgld#2496](https://github.com/vigilnetwork/vgl/pull/2496))
- mempool: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- mempool: Enforce explicit versions ([Vigil/vgld#2716](https://github.com/vigilnetwork/vgl/pull/2716))
- mempool: Remove unneeded max tx ver std checks ([Vigil/vgld#2716](https://github.com/vigilnetwork/vgl/pull/2716))
- mempool: Update fraud proof data ([Vigil/vgld#2804](https://github.com/vigilnetwork/vgl/pull/2804))
- mempool: CheckTransactionInputs check fraud proof ([Vigil/vgld#2804](https://github.com/vigilnetwork/vgl/pull/2804))

### Mining:

- mining: Move txPriorityQueue to a separate file ([Vigil/vgld#2431](https://github.com/vigilnetwork/vgl/pull/2431))
- mining: Move interfaces to mining/interface.go ([Vigil/vgld#2431](https://github.com/vigilnetwork/vgl/pull/2431))
- mining: Add method comments to blockManagerFacade ([Vigil/vgld#2431](https://github.com/vigilnetwork/vgl/pull/2431))
- mining: Move BgBlkTmplGenerator to separate file ([Vigil/vgld#2431](https://github.com/vigilnetwork/vgl/pull/2431))
- mining: Prevent panic in child prio item handling ([Vigil/vgld#2434](https://github.com/vigilnetwork/vgl/pull/2434))
- mining: Add Config struct to house mining params ([Vigil/vgld#2436](https://github.com/vigilnetwork/vgl/pull/2436))
- mining: Move block chain functions to Config ([Vigil/vgld#2436](https://github.com/vigilnetwork/vgl/pull/2436))
- mining: Move txMiningView from mempool package ([Vigil/vgld#2467](https://github.com/vigilnetwork/vgl/pull/2467))
- mining: Switch to custom waitGroup impl ([Vigil/vgld#2477](https://github.com/vigilnetwork/vgl/pull/2477))
- mining: Remove leftover block manager facade iface ([Vigil/vgld#2510](https://github.com/vigilnetwork/vgl/pull/2510))
- mining: No error log on expected head reorg errors ([Vigil/vgld#2560](https://github.com/vigilnetwork/vgl/pull/2560))
- mining: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- mining: Add error kinds for auto revocations ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- mining: Add auto revocation priority to tx queue ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- mining: Add HeaderByHash to Config ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- mining: Prevent unnecessary reorg with equal votes ([Vigil/vgld#2840](https://github.com/vigilnetwork/vgl/pull/2840))
- mining: Update to latest block vers for HFV ([Vigil/vgld#2753](https://github.com/vigilnetwork/vgl/pull/2753))

### RPC:

- rpcserver: Upgrade is deprecated; switch to Upgrader ([Vigil/vgld#2409](https://github.com/vigilnetwork/vgl/pull/2409))
- multi: Add TAdd support to getrawmempool ([Vigil/vgld#2448](https://github.com/vigilnetwork/vgl/pull/2448))
- rpcserver: Update getrawmempool txtype help ([Vigil/vgld#2452](https://github.com/vigilnetwork/vgl/pull/2452))
- rpcserver: Hash auth using random-keyed MAC ([Vigil/vgld#2486](https://github.com/vigilnetwork/vgl/pull/2486))
- rpcserver: Use next stake diff from snapshot ([Vigil/vgld#2493](https://github.com/vigilnetwork/vgl/pull/2493))
- rpcserver: Make authenticate match header auth ([Vigil/vgld#2502](https://github.com/vigilnetwork/vgl/pull/2502))
- rpcserver: Check unauthorized access in const time ([Vigil/vgld#2509](https://github.com/vigilnetwork/vgl/pull/2509))
- multi: Subscribe for work ntfns in rpcserver ([Vigil/vgld#2501](https://github.com/vigilnetwork/vgl/pull/2501))
- rpcserver: Prune block templates in websocket path ([Vigil/vgld#2503](https://github.com/vigilnetwork/vgl/pull/2503))
- rpcserver: Remove version from gettxout result ([Vigil/vgld#2517](https://github.com/vigilnetwork/vgl/pull/2517))
- rpcserver: Add tree param to gettxout ([Vigil/vgld#2517](https://github.com/vigilnetwork/vgl/pull/2517))
- rpcserver/netsync: Remove notifystakedifficulty ([Vigil/vgld#2519](https://github.com/vigilnetwork/vgl/pull/2519))
- rpcserver: Remove v1 getcfilter{,header} ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- rpcserver: Remove unused Filterer interface ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- rpcserver: Update getblockchaininfo best header ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- rpcserver: Remove unused LocateBlocks iface method ([Vigil/vgld#2538](https://github.com/vigilnetwork/vgl/pull/2538))
- rpcserver: Allow TLS client cert authentication ([Vigil/vgld#2482](https://github.com/vigilnetwork/vgl/pull/2482))
- rpcserver: Add invalidate/reconsiderblock support ([Vigil/vgld#2536](https://github.com/vigilnetwork/vgl/pull/2536))
- rpcserver: Support getblockchaininfo genesis block ([Vigil/vgld#2550](https://github.com/vigilnetwork/vgl/pull/2550))
- rpcserver: Calc verify progress based on best hdr ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- rpcserver: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- rpcserver: Allow gettreasurybalance empty blk str ([Vigil/vgld#2640](https://github.com/vigilnetwork/vgl/pull/2640))
- rpcserver: Add median time to verbose results ([Vigil/vgld#2638](https://github.com/vigilnetwork/vgl/pull/2638))
- rpcserver: Allow interface names for dial addresses ([Vigil/vgld#2623](https://github.com/vigilnetwork/vgl/pull/2623))
- rpcserver: Add script version to gettxout ([Vigil/vgld#2650](https://github.com/vigilnetwork/vgl/pull/2650))
- rpcserver: Remove unused help entry ([Vigil/vgld#2648](https://github.com/vigilnetwork/vgl/pull/2648))
- rpcserver: Set script version in raw tx results ([Vigil/vgld#2663](https://github.com/vigilnetwork/vgl/pull/2663))
- rpcserver: Impose additional read limits ([Vigil/vgld#2675](https://github.com/vigilnetwork/vgl/pull/2675))
- rpcserver: Add more strict request origin check ([Vigil/vgld#2676](https://github.com/vigilnetwork/vgl/pull/2676))
- rpcserver: Use duplicate tx error for recently mined transactions ([Vigil/vgld#2705](https://github.com/vigilnetwork/vgl/pull/2705))
- rpcserver: Wait for sync on rpc request ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- rpcserver: Update websocket ping timeout handling ([Vigil/vgld#2866](https://github.com/vigilnetwork/vgl/pull/2866))

### vgld command-line flags and configuration:

- multi: Rename BMGR subsystem to SYNC ([Vigil/vgld#2500](https://github.com/vigilnetwork/vgl/pull/2500))
- server/indexers: Remove v1 cfilter indexing support ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- config: Add utxocachemaxsize ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- main: Update slog for LOGFLAGS=nodatetime support ([Vigil/vgld#2608](https://github.com/vigilnetwork/vgl/pull/2608))
- config: Allow interface names for listener addresses ([Vigil/vgld#2623](https://github.com/vigilnetwork/vgl/pull/2623))
- config: Correct dir create failure error message ([Vigil/vgld#2682](https://github.com/vigilnetwork/vgl/pull/2682))
- config: Add logsize config option ([Vigil/vgld#2711](https://github.com/vigilnetwork/vgl/pull/2711))
- config: conditionally generate rpc credentials ([Vigil/vgld#2779](https://github.com/vigilnetwork/vgl/pull/2779))
- multi: Add assumevalid config option ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))

### gencerts utility changes:

- gencerts: Add certificate authority capabilities ([Vigil/vgld#2478](https://github.com/vigilnetwork/vgl/pull/2478))
- gencerts: Add RSA support (4096 bit keys only) ([Vigil/vgld#2551](https://github.com/vigilnetwork/vgl/pull/2551))

### addblock utility changes:

- cmd/addblock: update block importer ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- addblock: Run index subscriber as a goroutine ([Vigil/vgld#2760](https://github.com/vigilnetwork/vgl/pull/2760))
- addblock: Fix blockchain initialization ([Vigil/vgld#2760](https://github.com/vigilnetwork/vgl/pull/2760))
- addblock: Use chain bulk import mode ([Vigil/vgld#2782](https://github.com/vigilnetwork/vgl/pull/2782))

### findcheckpoint utility changes:

- findcheckpoint: Fix blockchain initialization ([Vigil/vgld#2759](https://github.com/vigilnetwork/vgl/pull/2759))

### Documentation:

- docs: Fix JSON-RPC API gettxoutsetinfo description ([Vigil/vgld#2443](https://github.com/vigilnetwork/vgl/pull/2443))
- docs: Add JSON-RPC API getpeerinfo missing fields ([Vigil/vgld#2443](https://github.com/vigilnetwork/vgl/pull/2443))
- docs: Fix JSON-RPC API gettreasurybalance fmt ([Vigil/vgld#2443](https://github.com/vigilnetwork/vgl/pull/2443))
- docs: Fix JSON-RPC API gettreasuryspendvotes fmt ([Vigil/vgld#2443](https://github.com/vigilnetwork/vgl/pull/2443))
- docs: Add JSON-RPC API searchrawtxns req limit ([Vigil/vgld#2443](https://github.com/vigilnetwork/vgl/pull/2443))
- docs: Update JSON-RPC API getrawmempool ([Vigil/vgld#2453](https://github.com/vigilnetwork/vgl/pull/2453))
- progresslog: Add package documentation ([Vigil/vgld#2499](https://github.com/vigilnetwork/vgl/pull/2499))
- netsync: Add package documentation ([Vigil/vgld#2500](https://github.com/vigilnetwork/vgl/pull/2500))
- multi: update error code related documentation ([Vigil/vgld#2515](https://github.com/vigilnetwork/vgl/pull/2515))
- docs: Update JSON-RPC API getwork to match reality ([Vigil/vgld#2526](https://github.com/vigilnetwork/vgl/pull/2526))
- docs: Remove notifystakedifficulty JSON-RPC API ([Vigil/vgld#2519](https://github.com/vigilnetwork/vgl/pull/2519))
- docs: Remove v1 getcfilter{,header} JSON-RPC API ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- chaincfg: Update doc.go ([Vigil/vgld#2528](https://github.com/vigilnetwork/vgl/pull/2528))
- blockchain: Update README.md and doc.go ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- docs: Add invalidate/reconsiderblock JSON-RPC API ([Vigil/vgld#2536](https://github.com/vigilnetwork/vgl/pull/2536))
- docs: Add release notes for v1.6.0 ([Vigil/vgld#2451](https://github.com/vigilnetwork/vgl/pull/2451))
- multi: Update README.md files for go modules ([Vigil/vgld#2559](https://github.com/vigilnetwork/vgl/pull/2559))
- apbf: Add README.md ([Vigil/vgld#2579](https://github.com/vigilnetwork/vgl/pull/2579))
- docs: Add release notes for v1.6.1 ([Vigil/vgld#2601](https://github.com/vigilnetwork/vgl/pull/2601))
- docs: Update min recommended specs in README.md ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- stdaddr: Add README.md ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add serialized pubkey info to README.md ([Vigil/vgld#2619](https://github.com/vigilnetwork/vgl/pull/2619))
- docs: Add release notes for v1.6.2 ([Vigil/vgld#2630](https://github.com/vigilnetwork/vgl/pull/2630))
- docs: Add scriptpubkey json returns ([Vigil/vgld#2650](https://github.com/vigilnetwork/vgl/pull/2650))
- stdscript: Add README.md ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stake: Comment on max SSGen outputs with treasury ([Vigil/vgld#2664](https://github.com/vigilnetwork/vgl/pull/2664))
- docs: Update JSON-RPC API for script version ([Vigil/vgld#2663](https://github.com/vigilnetwork/vgl/pull/2663))
- docs: Update JSON-RPC API for max request limits ([Vigil/vgld#2675](https://github.com/vigilnetwork/vgl/pull/2675))
- docs: Add SECURITY.md file ([Vigil/vgld#2717](https://github.com/vigilnetwork/vgl/pull/2717))
- sampleconfig: Add missing log options ([Vigil/vgld#2723](https://github.com/vigilnetwork/vgl/pull/2723))
- docs: Update go versions in README.md ([Vigil/vgld#2722](https://github.com/vigilnetwork/vgl/pull/2722))
- docs: Correct generate description ([Vigil/vgld#2724](https://github.com/vigilnetwork/vgl/pull/2724))
- database: Correct README rpcclient link ([Vigil/vgld#2725](https://github.com/vigilnetwork/vgl/pull/2725))
- docs: Add accuracy and reliability to README.md ([Vigil/vgld#2726](https://github.com/vigilnetwork/vgl/pull/2726))
- sampleconfig: Update for deprecated nodnsseed ([Vigil/vgld#2728](https://github.com/vigilnetwork/vgl/pull/2728))
- docs: Update for secp256k1 v4 module ([Vigil/vgld#2732](https://github.com/vigilnetwork/vgl/pull/2732))
- docs: Update for new modules ([Vigil/vgld#2744](https://github.com/vigilnetwork/vgl/pull/2744))
- sampleconfig: update rpc credentials documentation ([Vigil/vgld#2779](https://github.com/vigilnetwork/vgl/pull/2779))
- docs: Update for addrmgr v2 module ([Vigil/vgld#2797](https://github.com/vigilnetwork/vgl/pull/2797))
- docs: Update for rpc/jsonrpc/types v3 module ([Vigil/vgld#2801](https://github.com/vigilnetwork/vgl/pull/2801))
- stdscript: Update README.md for provably pruneable ([Vigil/vgld#2803](https://github.com/vigilnetwork/vgl/pull/2803))
- docs: Update for txscript v3 module ([Vigil/vgld#2815](https://github.com/vigilnetwork/vgl/pull/2815))
- docs: Update for VGLutil v4 module ([Vigil/vgld#2818](https://github.com/vigilnetwork/vgl/pull/2818))
- uint256: Add README.md ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- docs: Update for peer v3 module ([Vigil/vgld#2820](https://github.com/vigilnetwork/vgl/pull/2820))
- docs: Update for database v3 module ([Vigil/vgld#2822](https://github.com/vigilnetwork/vgl/pull/2822))
- docs: Update for blockchain/stake v4 module ([Vigil/vgld#2824](https://github.com/vigilnetwork/vgl/pull/2824))
- docs: Update for gcs v3 module ([Vigil/vgld#2830](https://github.com/vigilnetwork/vgl/pull/2830))
- docs: Fix typos and trailing whitespace ([Vigil/vgld#2843](https://github.com/vigilnetwork/vgl/pull/2843))
- docs: Add max line length and wrapping guidelines ([Vigil/vgld#2843](https://github.com/vigilnetwork/vgl/pull/2843))
- docs: Update for math/uint256 module ([Vigil/vgld#2842](https://github.com/vigilnetwork/vgl/pull/2842))
- docs: Update simnet env docs for subsidy split ([Vigil/vgld#2848](https://github.com/vigilnetwork/vgl/pull/2848))
- docs: Update for blockchain v4 module ([Vigil/vgld#2831](https://github.com/vigilnetwork/vgl/pull/2831))
- docs: Update for rpcclient v7 module ([Vigil/vgld#2851](https://github.com/vigilnetwork/vgl/pull/2851))
- primitives: Add skeleton README.md ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))

### Contrib changes:

- contrib: Update OpenBSD rc script for 6.9 features ([Vigil/vgld#2646](https://github.com/vigilnetwork/vgl/pull/2646))
- contrib: Bump Dockerfile.alpine to alpine:3.14.0 ([Vigil/vgld#2681](https://github.com/vigilnetwork/vgl/pull/2681))
- build: Use go 1.17 in Dockerfiles ([Vigil/vgld#2722](https://github.com/vigilnetwork/vgl/pull/2722))
- build: Pin docker images with SHA instead of tag ([Vigil/vgld#2735](https://github.com/vigilnetwork/vgl/pull/2735))
- build/contrib: Improve docker support ([Vigil/vgld#2740](https://github.com/vigilnetwork/vgl/pull/2740))

### Developer-related package and module changes:

- VGLjson: Reject dup method type registrations ([Vigil/vgld#2417](https://github.com/vigilnetwork/vgl/pull/2417))
- peer: various cleanups ([Vigil/vgld#2396](https://github.com/vigilnetwork/vgl/pull/2396))
- blockchain: Create treasury buckets during upgrade ([Vigil/vgld#2441](https://github.com/vigilnetwork/vgl/pull/2441))
- blockchain: Fix stxosToScriptSource ([Vigil/vgld#2444](https://github.com/vigilnetwork/vgl/pull/2444))
- rpcserver: add NtfnManager interface ([Vigil/vgld#2410](https://github.com/vigilnetwork/vgl/pull/2410))
- lru: Fix lookup race on small caches ([Vigil/vgld#2464](https://github.com/vigilnetwork/vgl/pull/2464))
- gcs: update error types ([Vigil/vgld#2262](https://github.com/vigilnetwork/vgl/pull/2262))
- main: Switch windows service dependency ([Vigil/vgld#2479](https://github.com/vigilnetwork/vgl/pull/2479))
- blockchain: Simplify upgrade single run stage code ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- blockchain: Simplify upgrade batching logic ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- blockchain: Use new batching logic for filter init ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- blockchain: Use new batch logic for blkidx upgrade ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- blockchain: Use new batch logic for utxos upgrade ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- blockchain: Use new batch logic for spends upgrade ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- blockchain: Use new batch logic for clr failed ([Vigil/vgld#2457](https://github.com/vigilnetwork/vgl/pull/2457))
- windows: Switch to os.Executable ([Vigil/vgld#2485](https://github.com/vigilnetwork/vgl/pull/2485))
- blockchain: Revert fast add reversal ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Less order dependent full blocks tests ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Move context free tx sanity checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Move context free block sanity checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Rework contextual tx checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Move {coin,trsy}base contextual checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Move staketx-related contextual checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Move sigop-related contextual checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Make CheckBlockSanity context free ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Context free CheckTransactionSanity ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- blockchain: Move contextual treasury spend checks ([Vigil/vgld#2481](https://github.com/vigilnetwork/vgl/pull/2481))
- mempool: Comment and stylistic updates ([Vigil/vgld#2480](https://github.com/vigilnetwork/vgl/pull/2480))
- mining: Rename TxMiningView Remove method ([Vigil/vgld#2490](https://github.com/vigilnetwork/vgl/pull/2490))
- mining: Unexport TxMiningView methods ([Vigil/vgld#2490](https://github.com/vigilnetwork/vgl/pull/2490))
- mining: Update mergeUtxoView comment ([Vigil/vgld#2490](https://github.com/vigilnetwork/vgl/pull/2490))
- blockchain: Consolidate deployment errors ([Vigil/vgld#2487](https://github.com/vigilnetwork/vgl/pull/2487))
- blockchain: Consolidate unknown block errors ([Vigil/vgld#2487](https://github.com/vigilnetwork/vgl/pull/2487))
- blockchain: Consolidate no filter errors ([Vigil/vgld#2487](https://github.com/vigilnetwork/vgl/pull/2487))
- blockchain: Consolidate no treasury bal errors ([Vigil/vgld#2487](https://github.com/vigilnetwork/vgl/pull/2487))
- blockchain: Convert to LRU block cache ([Vigil/vgld#2488](https://github.com/vigilnetwork/vgl/pull/2488))
- blockchain: Remove unused error returns ([Vigil/vgld#2489](https://github.com/vigilnetwork/vgl/pull/2489))
- blockmanager: Remove unused stakediff infra ([Vigil/vgld#2493](https://github.com/vigilnetwork/vgl/pull/2493))
- server: Use next stake diff from snapshot ([Vigil/vgld#2493](https://github.com/vigilnetwork/vgl/pull/2493))
- blockchain: Explicit hash in next stake diff calcs ([Vigil/vgld#2494](https://github.com/vigilnetwork/vgl/pull/2494))
- blockchain: Explicit hash in LN agenda active func ([Vigil/vgld#2495](https://github.com/vigilnetwork/vgl/pull/2495))
- blockmanager: Remove unused config field ([Vigil/vgld#2497](https://github.com/vigilnetwork/vgl/pull/2497))
- blockmanager: Decouple block database code ([Vigil/vgld#2497](https://github.com/vigilnetwork/vgl/pull/2497))
- blockmanager: Decouple from global config var ([Vigil/vgld#2497](https://github.com/vigilnetwork/vgl/pull/2497))
- blockchain: Explicit hash in max block size func ([Vigil/vgld#2507](https://github.com/vigilnetwork/vgl/pull/2507))
- progresslog: Make block progress log internal ([Vigil/vgld#2499](https://github.com/vigilnetwork/vgl/pull/2499))
- server: Do not use unexported block manager cfg ([Vigil/vgld#2498](https://github.com/vigilnetwork/vgl/pull/2498))
- blockmanager: Rework chain current logic ([Vigil/vgld#2498](https://github.com/vigilnetwork/vgl/pull/2498))
- multi: Handle chain ntfn callback in server ([Vigil/vgld#2498](https://github.com/vigilnetwork/vgl/pull/2498))
- server: Rename blockManager field to syncManager ([Vigil/vgld#2500](https://github.com/vigilnetwork/vgl/pull/2500))
- server: Add temp sync manager interface ([Vigil/vgld#2500](https://github.com/vigilnetwork/vgl/pull/2500))
- netsync: Split blockmanager into separate package ([Vigil/vgld#2500](https://github.com/vigilnetwork/vgl/pull/2500))
- netsync: Rename blockManager to SyncManager ([Vigil/vgld#2500](https://github.com/vigilnetwork/vgl/pull/2500))
- internal/ticketdb: update error types ([Vigil/vgld#2279](https://github.com/vigilnetwork/vgl/pull/2279))
- secp256k1/ecdsa: update error types ([Vigil/vgld#2281](https://github.com/vigilnetwork/vgl/pull/2281))
- secp256k1/schnorr: update error types ([Vigil/vgld#2282](https://github.com/vigilnetwork/vgl/pull/2282))
- VGLjson: update error types ([Vigil/vgld#2271](https://github.com/vigilnetwork/vgl/pull/2271))
- VGLec/secp256k1: update error types ([Vigil/vgld#2265](https://github.com/vigilnetwork/vgl/pull/2265))
- blockchain/stake: update error types ([Vigil/vgld#2264](https://github.com/vigilnetwork/vgl/pull/2264))
- multi: update database error types ([Vigil/vgld#2261](https://github.com/vigilnetwork/vgl/pull/2261))
- blockchain: Remove unused treasury active func ([Vigil/vgld#2514](https://github.com/vigilnetwork/vgl/pull/2514))
- stake: update ticket lottery errors ([Vigil/vgld#2433](https://github.com/vigilnetwork/vgl/pull/2433))
- netsync: Improve is current detection ([Vigil/vgld#2513](https://github.com/vigilnetwork/vgl/pull/2513))
- internal/mining: update mining error types ([Vigil/vgld#2515](https://github.com/vigilnetwork/vgl/pull/2515))
- multi: sprinkle on more errors.As/Is ([Vigil/vgld#2522](https://github.com/vigilnetwork/vgl/pull/2522))
- mining: Correct fee calculations during reorgs ([Vigil/vgld#2530](https://github.com/vigilnetwork/vgl/pull/2530))
- fees: Remove deprecated DisableLog ([Vigil/vgld#2529](https://github.com/vigilnetwork/vgl/pull/2529))
- rpcclient: Remove deprecated DisableLog ([Vigil/vgld#2527](https://github.com/vigilnetwork/vgl/pull/2527))
- rpcclient: Remove notifystakedifficulty ([Vigil/vgld#2519](https://github.com/vigilnetwork/vgl/pull/2519))
- rpc/jsonrpc/types: Remove notifystakedifficulty ([Vigil/vgld#2519](https://github.com/vigilnetwork/vgl/pull/2519))
- netsync: Remove unneeded ForceReorganization ([Vigil/vgld#2520](https://github.com/vigilnetwork/vgl/pull/2520))
- mining: Remove duplicate method ([Vigil/vgld#2520](https://github.com/vigilnetwork/vgl/pull/2520))
- multi: use EstimateSmartFeeResult ([Vigil/vgld#2283](https://github.com/vigilnetwork/vgl/pull/2283))
- rpcclient: Remove v1 getcfilter{,header} ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- rpc/jsonrpc/types: Remove v1 getcfilter{,header} ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- gcs: Remove unused v1 blockcf package ([Vigil/vgld#2525](https://github.com/vigilnetwork/vgl/pull/2525))
- blockchain: Remove legacy sequence lock view ([Vigil/vgld#2534](https://github.com/vigilnetwork/vgl/pull/2534))
- blockchain: Remove IsFixSeqLocksAgendaActive ([Vigil/vgld#2534](https://github.com/vigilnetwork/vgl/pull/2534))
- blockchain: Explicit hash in estimate stake diff ([Vigil/vgld#2524](https://github.com/vigilnetwork/vgl/pull/2524))
- netsync: Remove unneeded TipGeneration ([Vigil/vgld#2537](https://github.com/vigilnetwork/vgl/pull/2537))
- netsync: Remove unused TicketPoolValue ([Vigil/vgld#2544](https://github.com/vigilnetwork/vgl/pull/2544))
- netsync: Embed peers vs separate peer states ([Vigil/vgld#2541](https://github.com/vigilnetwork/vgl/pull/2541))
- netsync/server: Update peer heights directly ([Vigil/vgld#2542](https://github.com/vigilnetwork/vgl/pull/2542))
- netsync: Move proactive sigcache evict to server ([Vigil/vgld#2543](https://github.com/vigilnetwork/vgl/pull/2543))
- blockchain: Add invalidate/reconsider infrastruct ([Vigil/vgld#2536](https://github.com/vigilnetwork/vgl/pull/2536))
- rpc/jsonrpc/types: Add invalidate/reconsiderblock ([Vigil/vgld#2536](https://github.com/vigilnetwork/vgl/pull/2536))
- netsync: Convert lifecycle to context ([Vigil/vgld#2545](https://github.com/vigilnetwork/vgl/pull/2545))
- multi: Rework utxoset/view to use outpoints ([Vigil/vgld#2540](https://github.com/vigilnetwork/vgl/pull/2540))
- blockchain: Remove compression version param ([Vigil/vgld#2547](https://github.com/vigilnetwork/vgl/pull/2547))
- blockchain: Remove error from LatestBlockLocator ([Vigil/vgld#2548](https://github.com/vigilnetwork/vgl/pull/2548))
- blockchain: Fix incorrect decompressScript calls ([Vigil/vgld#2552](https://github.com/vigilnetwork/vgl/pull/2552))
- blockchain: Fix V3 spend journal migration ([Vigil/vgld#2552](https://github.com/vigilnetwork/vgl/pull/2552))
- multi: Remove blockChain field from UtxoViewpoint ([Vigil/vgld#2553](https://github.com/vigilnetwork/vgl/pull/2553))
- blockchain: Move UtxoEntry to a separate file ([Vigil/vgld#2553](https://github.com/vigilnetwork/vgl/pull/2553))
- blockchain: Update UtxoEntry Clone method comment ([Vigil/vgld#2553](https://github.com/vigilnetwork/vgl/pull/2553))
- progresslog: Make logger more generic ([Vigil/vgld#2555](https://github.com/vigilnetwork/vgl/pull/2555))
- server: Remove several unused funcs ([Vigil/vgld#2561](https://github.com/vigilnetwork/vgl/pull/2561))
- mempool: Store staged transactions as TxDesc ([Vigil/vgld#2319](https://github.com/vigilnetwork/vgl/pull/2319))
- connmgr: Add func to iterate conn reqs ([Vigil/vgld#2562](https://github.com/vigilnetwork/vgl/pull/2562))
- netsync: Correct check for needTx ([Vigil/vgld#2568](https://github.com/vigilnetwork/vgl/pull/2568))
- rpcclient: Update EstimateSmartFee return type ([Vigil/vgld#2255](https://github.com/vigilnetwork/vgl/pull/2255))
- server: Notify sync mgr later and track ntfn ([Vigil/vgld#2582](https://github.com/vigilnetwork/vgl/pull/2582))
- apbf: Introduce Age-Partitioned Bloom Filters ([Vigil/vgld#2579](https://github.com/vigilnetwork/vgl/pull/2579))
- apbf: Add basic usage example ([Vigil/vgld#2579](https://github.com/vigilnetwork/vgl/pull/2579))
- apbf: Add support to go generate a KL table ([Vigil/vgld#2579](https://github.com/vigilnetwork/vgl/pull/2579))
- apbf: Switch to fast reduce method ([Vigil/vgld#2584](https://github.com/vigilnetwork/vgl/pull/2584))
- server: Remove unneeded child context ([Vigil/vgld#2593](https://github.com/vigilnetwork/vgl/pull/2593))
- blockchain: Separate utxo state from tx flags ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Add utxoStateFresh to UtxoEntry ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Add size method to UtxoEntry ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Deep copy view entry script from tx ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Add utxoSetState to the database ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- multi: Add UtxoCache ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Make InitUtxoCache a UtxoCache method ([Vigil/vgld#2599](https://github.com/vigilnetwork/vgl/pull/2599))
- blockchain: Add UtxoCacher interface ([Vigil/vgld#2599](https://github.com/vigilnetwork/vgl/pull/2599))
- VGLutil: Correct ed25519 address constructor ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Introduce package infra for std addrs ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add infrastructure for v0 decoding ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2pk-ecdsa-secp256k1 support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2pk-ed25519 support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2pk-schnorr-secp256k1 support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2pkh-ecdsa-secp256k1 support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2pkh-ed25519 support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2pkh-schnorr-secp256k1 support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add v0 p2sh support ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- stdaddr: Add decode address example ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- txscript: Rename script bldr add data to unchecked ([Vigil/vgld#2611](https://github.com/vigilnetwork/vgl/pull/2611))
- txscript: Add script bldr unchecked op add ([Vigil/vgld#2611](https://github.com/vigilnetwork/vgl/pull/2611))
- rpcserver: Remove uncompressed pubkeys fast path ([Vigil/vgld#2617](https://github.com/vigilnetwork/vgl/pull/2617))
- blockchain: Allow alternate tips for current check ([Vigil/vgld#2612](https://github.com/vigilnetwork/vgl/pull/2612))
- txscript: Accept raw public keys in MultiSigScript ([Vigil/vgld#2615](https://github.com/vigilnetwork/vgl/pull/2615))
- cpuminer: Remove unused MiningAddrs from Config ([Vigil/vgld#2616](https://github.com/vigilnetwork/vgl/pull/2616))
- stdaddr: Add ability to obtain raw public key ([Vigil/vgld#2619](https://github.com/vigilnetwork/vgl/pull/2619))
- stdaddr: Move from internal/staging to txscript ([Vigil/vgld#2620](https://github.com/vigilnetwork/vgl/pull/2620))
- stdaddr: Accept vote and revoke limits separately ([Vigil/vgld#2624](https://github.com/vigilnetwork/vgl/pull/2624))
- stake: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- indexers: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- blockchain: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- rpcclient: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- hdkeychain: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSStx ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSStxChange ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSSGen ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSSGenSHDirect ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSSGenPKHDirect ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSSRtx ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSSRtxPKHDirect ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToSSRtxSHDirect ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToAddrScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused PayToScriptHashScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToSchnorrPubKeyScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToEdwardsPubKeyScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToPubKeyScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToScriptHashScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToPubKeyHashSchnorrScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToPubKeyHashEdwardsScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused payToPubKeyHashScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused GenerateSStxAddrPush ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Remove unused ErrUnsupportedAddress ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- txscript: Break VGLutil dependency ([Vigil/vgld#2626](https://github.com/vigilnetwork/vgl/pull/2626))
- stdaddr: Replace Address method with String ([Vigil/vgld#2633](https://github.com/vigilnetwork/vgl/pull/2633))
- VGLutil: Convert to use new stdaddr package ([Vigil/vgld#2628](https://github.com/vigilnetwork/vgl/pull/2628))
- VGLutil: Remove all code related to Address ([Vigil/vgld#2628](https://github.com/vigilnetwork/vgl/pull/2628))
- blockchain: Trsy always inactive for genesis blk ([Vigil/vgld#2636](https://github.com/vigilnetwork/vgl/pull/2636))
- blockchain: Use agenda flags for tx check context ([Vigil/vgld#2639](https://github.com/vigilnetwork/vgl/pull/2639))
- blockchain: Move UTXO DB methods to separate file ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- blockchain: Move UTXO DB tests to separate file ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- ipc: Fix lifetimeEvent comments ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- blockchain: Add utxoDatabaseInfo ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- multi: Introduce UTXO database ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- blockchain: Decouple stxo and utxo migrations ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- multi: Migrate to UTXO database ([Vigil/vgld#2632](https://github.com/vigilnetwork/vgl/pull/2632))
- main: Handle SIGHUP with clean shutdown ([Vigil/vgld#2645](https://github.com/vigilnetwork/vgl/pull/2645))
- txscript: Split signing code to sign subpackage ([Vigil/vgld#2642](https://github.com/vigilnetwork/vgl/pull/2642))
- database: Add Flush to DB interface ([Vigil/vgld#2649](https://github.com/vigilnetwork/vgl/pull/2649))
- multi: Flush block DB before UTXO DB ([Vigil/vgld#2649](https://github.com/vigilnetwork/vgl/pull/2649))
- blockchain: Flush UTXO DB after init utxoSetState ([Vigil/vgld#2649](https://github.com/vigilnetwork/vgl/pull/2649))
- blockchain: Force flush in separateUtxoDatabase ([Vigil/vgld#2649](https://github.com/vigilnetwork/vgl/pull/2649))
- version: Rework to support single version override ([Vigil/vgld#2651](https://github.com/vigilnetwork/vgl/pull/2651))
- blockchain: Remove UtxoCacher DB Tx dependency ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add UtxoBackend interface ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Export UtxoSetState ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add FetchEntry to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add PutUtxos to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add FetchState to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add FetchStats to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add FetchInfo to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Move LoadUtxoDB to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Add Upgrade to UtxoBackend ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- multi: Remove UTXO db in BlockChain and UtxoCache ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- blockchain: Export ViewFilteredSet ([Vigil/vgld#2652](https://github.com/vigilnetwork/vgl/pull/2652))
- stake: Return StakeAddress from cmtmt conversion ([Vigil/vgld#2655](https://github.com/vigilnetwork/vgl/pull/2655))
- stdscript: Introduce pkg infra for std scripts ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pk-ecdsa-secp256k1 support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pk-ed25519 support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pk-schnorr-secp256k1 support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pkh-ecdsa-secp256k1 support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pkh-ed25519 support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pkh-schnorr-secp256k1 support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2sh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 ecdsa multisig support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 ecdsa multisig redeem support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 nulldata support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake sub p2pkh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake sub p2sh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake gen p2pkh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake gen p2sh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake revoke p2pkh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake revoke p2sh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake change p2pkh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake change p2sh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 treasury add support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 treasury gen p2pkh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 treasury gen p2sh support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add ecdsa multisig creation script ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 atomic swap redeem support ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add example for determining script type ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add example for p2pkh extract ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add example of script hash extract ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- blockchain: Use scripts in tickets address query ([Vigil/vgld#2657](https://github.com/vigilnetwork/vgl/pull/2657))
- stake: Do not use standardness code in consensus ([Vigil/vgld#2658](https://github.com/vigilnetwork/vgl/pull/2658))
- blockchain: Remove unneeded OP_TADD maturity check ([Vigil/vgld#2659](https://github.com/vigilnetwork/vgl/pull/2659))
- stake: Add is treasury gen script ([Vigil/vgld#2660](https://github.com/vigilnetwork/vgl/pull/2660))
- blockchain: No standardness code in consensus ([Vigil/vgld#2661](https://github.com/vigilnetwork/vgl/pull/2661))
- gcs: No standardness code in consensus ([Vigil/vgld#2662](https://github.com/vigilnetwork/vgl/pull/2662))
- stake: Remove stale TODOs from CheckSSGenVotes ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- stake: Remove stale TODOs from CheckSSRtx ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- txscript: Move contains stake opcode to consensus ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Move stake blockref script to consensus ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Move stake votebits script to consensus ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Remove unused IsPubKeyHashScript ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Remove unused IsStakeChangeScript ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Remove unused PushedData ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- blockchain: Flush UtxoCache when latch to current ([Vigil/vgld#2671](https://github.com/vigilnetwork/vgl/pull/2671))
- VGLjson: Minor jsonerr.go update ([Vigil/vgld#2672](https://github.com/vigilnetwork/vgl/pull/2672))
- rpcclient: Cancel client context on shutdown ([Vigil/vgld#2678](https://github.com/vigilnetwork/vgl/pull/2678))
- blockchain: Remove serializeUtxoEntry error ([Vigil/vgld#2683](https://github.com/vigilnetwork/vgl/pull/2683))
- blockchain: Add IsTreasuryEnabled to AgendaFlags ([Vigil/vgld#2686](https://github.com/vigilnetwork/vgl/pull/2686))
- multi: Update block ntfns to contain AgendaFlags ([Vigil/vgld#2686](https://github.com/vigilnetwork/vgl/pull/2686))
- multi: Update ProcessOrphans to use AgendaFlags ([Vigil/vgld#2686](https://github.com/vigilnetwork/vgl/pull/2686))
- mempool: Add maybeAcceptTransaction AgendaFlags ([Vigil/vgld#2686](https://github.com/vigilnetwork/vgl/pull/2686))
- secp256k1: Allow code generation to compile again ([Vigil/vgld#2687](https://github.com/vigilnetwork/vgl/pull/2687))
- jsonrpc/types: Add missing Method type to vars ([Vigil/vgld#2688](https://github.com/vigilnetwork/vgl/pull/2688))
- blockchain: Add UTXO backend error kinds ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Add helper to convert leveldb errors ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Add UtxoBackendIterator interface ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Add UtxoBackendTx interface ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Add levelDbUtxoBackendTx type ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- multi: Update UtxoBackend to use leveldb directly ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- multi: Move UTXO database ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Unexport levelDbUtxoBackend ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Always use node lookup methods ([Vigil/vgld#2685](https://github.com/vigilnetwork/vgl/pull/2685))
- blockchain: Use short keys for block index ([Vigil/vgld#2685](https://github.com/vigilnetwork/vgl/pull/2685))
- rpcclient: Shutdown breaks reconnect sleep ([Vigil/vgld#2696](https://github.com/vigilnetwork/vgl/pull/2696))
- secp256k1: No deps on adaptor code for precomps ([Vigil/vgld#2690](https://github.com/vigilnetwork/vgl/pull/2690))
- secp256k1: Always initialize adaptor instance ([Vigil/vgld#2690](https://github.com/vigilnetwork/vgl/pull/2690))
- secp256k1: Optimize precomp values to use affine ([Vigil/vgld#2690](https://github.com/vigilnetwork/vgl/pull/2690))
- rpcserver: Handle getwork nil err during reorg ([Vigil/vgld#2700](https://github.com/vigilnetwork/vgl/pull/2700))
- secp256k1: Improve scalar mult readability ([Vigil/vgld#2695](https://github.com/vigilnetwork/vgl/pull/2695))
- secp256k1: Optimize NAF conversion ([Vigil/vgld#2695](https://github.com/vigilnetwork/vgl/pull/2695))
- blockchain: Verify state of VGLP0007 voting ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- blockchain: Rename max expenditure funcs ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- stake: Add ExpiringNextBlock method to Node ([Vigil/vgld#2701](https://github.com/vigilnetwork/vgl/pull/2701))
- rpcclient: Add GetNetworkInfo call ([Vigil/vgld#2703](https://github.com/vigilnetwork/vgl/pull/2703))
- stake: Pre-allocate lottery ticket index slice ([Vigil/vgld#2710](https://github.com/vigilnetwork/vgl/pull/2710))
- blockchain: Switch to treasuryValueType.IsDebit ([Vigil/vgld#2680](https://github.com/vigilnetwork/vgl/pull/2680))
- blockchain: Sum amounts added to treasury ([Vigil/vgld#2680](https://github.com/vigilnetwork/vgl/pull/2680))
- blockchain: Add maxTreasuryExpenditureVGLP0007 ([Vigil/vgld#2680](https://github.com/vigilnetwork/vgl/pull/2680))
- blockchain: Use new expenditure policy if activated ([Vigil/vgld#2680](https://github.com/vigilnetwork/vgl/pull/2680))
- blockchain: Add checkTicketRedeemers ([Vigil/vgld#2702](https://github.com/vigilnetwork/vgl/pull/2702))
- blockchain: Add NextExpiringTickets to BestState ([Vigil/vgld#2708](https://github.com/vigilnetwork/vgl/pull/2708))
- multi: Add FetchUtxoEntry to mining Config ([Vigil/vgld#2709](https://github.com/vigilnetwork/vgl/pull/2709))
- stake: Add func to create revocation from ticket ([Vigil/vgld#2707](https://github.com/vigilnetwork/vgl/pull/2707))
- rpcserver: Use CreateRevocationFromTicket ([Vigil/vgld#2707](https://github.com/vigilnetwork/vgl/pull/2707))
- multi: Don't use deprecated ioutil package ([Vigil/vgld#2722](https://github.com/vigilnetwork/vgl/pull/2722))
- blockchain: Consolidate tx check flag construction ([Vigil/vgld#2716](https://github.com/vigilnetwork/vgl/pull/2716))
- stake: Export Hash256PRNG UniformRandom ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- blockchain: Check auto revocations agenda state ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- multi: Add mempool IsAutoRevocationsAgendaActive ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- multi: Add auto revocations to agenda flags ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- multi: Check tx inputs auto revocations flag ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- blockchain: Add auto revocation error kinds ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- stake: Add auto revocation error kinds ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- blockchain: Move revocation checks block context ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- multi: Add isAutoRevocationsEnabled to CheckSSRtx ([Vigil/vgld#2719](https://github.com/vigilnetwork/vgl/pull/2719))
- addrmgr: Remove deprecated code ([Vigil/vgld#2729](https://github.com/vigilnetwork/vgl/pull/2729))
- peer: Remove deprecated DisableLog ([Vigil/vgld#2730](https://github.com/vigilnetwork/vgl/pull/2730))
- database: Remove deprecated DisableLog ([Vigil/vgld#2731](https://github.com/vigilnetwork/vgl/pull/2731))
- addrmgr: Decouple IP network checks from wire ([Vigil/vgld#2596](https://github.com/vigilnetwork/vgl/pull/2596))
- addrmgr: Rename network address type ([Vigil/vgld#2596](https://github.com/vigilnetwork/vgl/pull/2596))
- addrmgr: Decouple addrmgr from wire NetAddress ([Vigil/vgld#2596](https://github.com/vigilnetwork/vgl/pull/2596))
- multi: add spend pruner ([Vigil/vgld#2641](https://github.com/vigilnetwork/vgl/pull/2641))
- multi: synchronize spend prunes and notifications ([Vigil/vgld#2641](https://github.com/vigilnetwork/vgl/pull/2641))
- blockchain: workSorterLess -> betterCandidate ([Vigil/vgld#2747](https://github.com/vigilnetwork/vgl/pull/2747))
- mempool: Add HeaderByHash to Config ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- rpctest: Remove unused BlockVersion const ([Vigil/vgld#2754](https://github.com/vigilnetwork/vgl/pull/2754))
- blockchain: Handle genesis auto revocation agenda ([Vigil/vgld#2755](https://github.com/vigilnetwork/vgl/pull/2755))
- indexers: remove index manager ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- indexers: add index subscriber ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- indexers: refactor interfaces ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- indexers: async transaction index ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- indexers: update address index ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- indexers: async exists address index ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- multi: integrate index subscriber ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- multi: avoid using subscriber lifecycle in catchup ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- multi: remove spend deps on index disc. notifs ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- multi: copy snapshot pkScript ([Vigil/vgld#2219](https://github.com/vigilnetwork/vgl/pull/2219))
- blockchain: Conditionally log difficulty retarget ([Vigil/vgld#2761](https://github.com/vigilnetwork/vgl/pull/2761))
- multi: Use single latest checkpoint ([Vigil/vgld#2763](https://github.com/vigilnetwork/vgl/pull/2763))
- blockchain: Move diff retarget log to connect ([Vigil/vgld#2765](https://github.com/vigilnetwork/vgl/pull/2765))
- multi: source index notif. from block notif ([Vigil/vgld#2256](https://github.com/vigilnetwork/vgl/pull/2256))
- server: fix wireToAddrmgrNetAddress data race ([Vigil/vgld#2758](https://github.com/vigilnetwork/vgl/pull/2758))
- multi: Flush cache before fetching UTXO stats ([Vigil/vgld#2767](https://github.com/vigilnetwork/vgl/pull/2767))
- blockchain: Don't use deprecated ioutil package ([Vigil/vgld#2769](https://github.com/vigilnetwork/vgl/pull/2769))
- blockchain: Fix ticket db disconnect revocations ([Vigil/vgld#2768](https://github.com/vigilnetwork/vgl/pull/2768))
- blockchain: Add convenience ancestor of func ([Vigil/vgld#2771](https://github.com/vigilnetwork/vgl/pull/2771))
- blockchain: Use new ancestor of convenience func ([Vigil/vgld#2771](https://github.com/vigilnetwork/vgl/pull/2771))
- blockchain: Remove unused latest blk locator func ([Vigil/vgld#2772](https://github.com/vigilnetwork/vgl/pull/2772))
- blockchain: Remove unused next lottery data func ([Vigil/vgld#2773](https://github.com/vigilnetwork/vgl/pull/2773))
- secp256k1: Correct 96-bit accum double overflow ([Vigil/vgld#2778](https://github.com/vigilnetwork/vgl/pull/2778))
- blockchain: Further decouple upgrade code ([Vigil/vgld#2776](https://github.com/vigilnetwork/vgl/pull/2776))
- blockchain: Add bulk import mode ([Vigil/vgld#2782](https://github.com/vigilnetwork/vgl/pull/2782))
- multi: Remove flags from SyncManager ProcessBlock ([Vigil/vgld#2783](https://github.com/vigilnetwork/vgl/pull/2783))
- netsync: Remove flags from processBlockMsg ([Vigil/vgld#2783](https://github.com/vigilnetwork/vgl/pull/2783))
- multi: Remove flags from blockchain ProcessBlock ([Vigil/vgld#2783](https://github.com/vigilnetwork/vgl/pull/2783))
- multi: Remove flags from ProcessBlockHeader ([Vigil/vgld#2785](https://github.com/vigilnetwork/vgl/pull/2785))
- blockchain: Remove flags maybeAcceptBlockHeader ([Vigil/vgld#2785](https://github.com/vigilnetwork/vgl/pull/2785))
- version: Use uint32 for major/minor/patch ([Vigil/vgld#2789](https://github.com/vigilnetwork/vgl/pull/2789))
- wire: Write message header directly ([Vigil/vgld#2790](https://github.com/vigilnetwork/vgl/pull/2790))
- stake: Correct treasury enabled vote discovery ([Vigil/vgld#2780](https://github.com/vigilnetwork/vgl/pull/2780))
- blockchain: Correct treasury spend vote data ([Vigil/vgld#2780](https://github.com/vigilnetwork/vgl/pull/2780))
- blockchain: UTXO database migration fix ([Vigil/vgld#2798](https://github.com/vigilnetwork/vgl/pull/2798))
- blockchain: Handle zero-length UTXO backend state ([Vigil/vgld#2798](https://github.com/vigilnetwork/vgl/pull/2798))
- mining: Remove unnecessary tx copy ([Vigil/vgld#2792](https://github.com/vigilnetwork/vgl/pull/2792))
- multi: Use VGLutil Tx in NewTxDeepTxIns ([Vigil/vgld#2802](https://github.com/vigilnetwork/vgl/pull/2802))
- indexers: synchronize index subscriber ntfn sends/receives ([Vigil/vgld#2806](https://github.com/vigilnetwork/vgl/pull/2806))
- stdscript: Add exported MaxDataCarrierSizeV0 ([Vigil/vgld#2803](https://github.com/vigilnetwork/vgl/pull/2803))
- stdscript: Add ProvablyPruneableScriptV0 ([Vigil/vgld#2803](https://github.com/vigilnetwork/vgl/pull/2803))
- stdscript: Add num required sigs support ([Vigil/vgld#2805](https://github.com/vigilnetwork/vgl/pull/2805))
- netsync: Remove unused RpcServer ([Vigil/vgld#2811](https://github.com/vigilnetwork/vgl/pull/2811))
- netsync: Consolidate initial sync handling ([Vigil/vgld#2812](https://github.com/vigilnetwork/vgl/pull/2812))
- stdscript: Add v0 p2pk-ed25519 extract ([Vigil/vgld#2807](https://github.com/vigilnetwork/vgl/pull/2807))
- stdscript: Add v0 p2pk-schnorr-secp256k1 extract ([Vigil/vgld#2807](https://github.com/vigilnetwork/vgl/pull/2807))
- stdscript: Add v0 p2pkh-ed25519 extract ([Vigil/vgld#2807](https://github.com/vigilnetwork/vgl/pull/2807))
- stdscript: Add v0 p2pkh-schnorr-secp256k1 extract ([Vigil/vgld#2807](https://github.com/vigilnetwork/vgl/pull/2807))
- stdscript: Add script to address conversion ([Vigil/vgld#2807](https://github.com/vigilnetwork/vgl/pull/2807))
- stdscript: Move from internal/staging to txscript ([Vigil/vgld#2810](https://github.com/vigilnetwork/vgl/pull/2810))
- mining: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- mempool: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- chaingen: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- blockchain: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- indexers: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- indexers: Remove unused trsy enabled params ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript/sign: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript/sign: Remove unused trsy enabled params ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- rpcserver: Convert to use stdscript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove deprecated ExtractAtomicSwapDataPushes ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused GenerateProvablyPruneableOut ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused MultiSigScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused MultisigRedeemScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused CalcMultiSigStats ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused IsMultisigScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused IsMultisigSigScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused ExtractPkScriptAltSigType ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused GetScriptClass ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused GetStakeOutSubclass ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused typeOfScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isTreasurySpendScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isMultisigScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused ExtractPkScriptAddrs ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused scriptHashToAddrs ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused pubKeyHashToAddrs ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isTreasuryAddScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractMultisigScriptDetails ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isStakeChangeScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isPubKeyHashScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isStakeRevocationScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isStakeGenScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isStakeSubmissionScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractStakeScriptHash ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractStakePubKeyHash ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isNullDataScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractPubKeyHash ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isPubKeyAltScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractPubKeyAltDetails ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isPubKeyScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractPubKey ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractUncompressedPubKey ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractCompressedPubKey ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isPubKeyHashAltScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused extractPubKeyHashAltDetails ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused isStandardAltSignatureType ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused MaxDataCarrierSize ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused ScriptClass ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused ErrNotMultisigScript ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused ErrTooManyRequiredSigs ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- txscript: Remove unused ErrTooMuchNullData ([Vigil/vgld#2808](https://github.com/vigilnetwork/vgl/pull/2808))
- stdaddr: Use txscript for opcode definitions ([Vigil/vgld#2809](https://github.com/vigilnetwork/vgl/pull/2809))
- stdscript: Add v0 stake-tagged p2pkh extract ([Vigil/vgld#2816](https://github.com/vigilnetwork/vgl/pull/2816))
- stdscript: Add v0 stake-tagged p2sh extract ([Vigil/vgld#2816](https://github.com/vigilnetwork/vgl/pull/2816))
- server: sync rebroadcast inv sends/receives ([Vigil/vgld#2814](https://github.com/vigilnetwork/vgl/pull/2814))
- multi: Move last ann block from peer to netsync ([Vigil/vgld#2821](https://github.com/vigilnetwork/vgl/pull/2821))
- uint256: Introduce package infrastructure ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add set from big endian bytes ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add set from little endian bytes ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add get big endian bytes ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add get little endian bytes ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add zero support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add uint32 casting support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add uint64 casting support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add equality comparison support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add less than comparison support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add less or equals comparison support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add greater than comparison support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add greater or equals comparison support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add general comparison support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add addition support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add subtraction support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add multiplication support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add squaring support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add division support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add negation support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add is odd support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise left shift support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise right shift support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise not support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise or support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise and support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise xor support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bit length support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add text formatting support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add conversion to stdlib big int support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add conversion from stdlib big int support ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add basic usage example ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- stake: Rename func to identify stake cmtmnt output ([Vigil/vgld#2824](https://github.com/vigilnetwork/vgl/pull/2824))
- progresslog: Make header logging concurrent safe ([Vigil/vgld#2833](https://github.com/vigilnetwork/vgl/pull/2833))
- netsync: Contiguous hashes for initial state reqs ([Vigil/vgld#2825](https://github.com/vigilnetwork/vgl/pull/2825))
- multi: Allow discrete mining with invalidated tip ([Vigil/vgld#2838](https://github.com/vigilnetwork/vgl/pull/2838))
- primitives: Add difficulty bits <-> uint256 ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add work calc from diff bits ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add hash to uint256 conversion ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add check proof of work ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add core merkle tree root calcs ([Vigil/vgld#2826](https://github.com/vigilnetwork/vgl/pull/2826))
- primitives: Add inclusion proof funcs ([Vigil/vgld#2827](https://github.com/vigilnetwork/vgl/pull/2827))
- indexers: update indexer error types ([Vigil/vgld#2770](https://github.com/vigilnetwork/vgl/pull/2770))
- rpcserver: Submit transactions directly ([Vigil/vgld#2835](https://github.com/vigilnetwork/vgl/pull/2835))
- netsync: Remove unused tx submission processing ([Vigil/vgld#2835](https://github.com/vigilnetwork/vgl/pull/2835))
- internal/staging: add ban manager ([Vigil/vgld#2554](https://github.com/vigilnetwork/vgl/pull/2554))
- uint256: Correct base 10 output formatting ([Vigil/vgld#2844](https://github.com/vigilnetwork/vgl/pull/2844))
- multi: Add assumeValid to BlockChain ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- blockchain: Track assumed valid node ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- blockchain: Set BFFastAdd based on assume valid ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- blockchain: Assume valid skip script validation ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- blockchain: Bulk import skip script validation ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- hdkeychain: Add a strict BIP32 child derivation method ([Vigil/vgld#2845](https://github.com/vigilnetwork/vgl/pull/2845))
- mempool: Consolidate tx check flag construction ([Vigil/vgld#2846](https://github.com/vigilnetwork/vgl/pull/2846))
- standalone: Add modified subsidy split calcs ([Vigil/vgld#2848](https://github.com/vigilnetwork/vgl/pull/2848))

### Developer-related module management:

- rpcclient: Prepare v6.0.1 ([Vigil/vgld#2455](https://github.com/vigilnetwork/vgl/pull/2455))
- multi: Start blockchain v4 module dev cycle ([Vigil/vgld#2463](https://github.com/vigilnetwork/vgl/pull/2463))
- multi: Start rpcclient v7 module dev cycle ([Vigil/vgld#2463](https://github.com/vigilnetwork/vgl/pull/2463))
- multi: Start gcs v3 module dev cycle ([Vigil/vgld#2463](https://github.com/vigilnetwork/vgl/pull/2463))
- multi: Start blockchain/stake v4 module dev cycle ([Vigil/vgld#2511](https://github.com/vigilnetwork/vgl/pull/2511))
- multi: Start txscript v4 module dev cycle ([Vigil/vgld#2511](https://github.com/vigilnetwork/vgl/pull/2511))
- multi: Start VGLutil v4 module dev cycle ([Vigil/vgld#2511](https://github.com/vigilnetwork/vgl/pull/2511))
- multi: Start VGLec/secp256k1 v4 module dev cycle ([Vigil/vgld#2511](https://github.com/vigilnetwork/vgl/pull/2511))
- rpc/jsonrpc/types: Start v3 module dev cycle ([Vigil/vgld#2517](https://github.com/vigilnetwork/vgl/pull/2517))
- multi: Round 1 prerel module release ver updates ([Vigil/vgld#2569](https://github.com/vigilnetwork/vgl/pull/2569))
- multi: Round 2 prerel module release ver updates ([Vigil/vgld#2570](https://github.com/vigilnetwork/vgl/pull/2570))
- multi: Round 3 prerel module release ver updates ([Vigil/vgld#2572](https://github.com/vigilnetwork/vgl/pull/2572))
- multi: Round 4 prerel module release ver updates ([Vigil/vgld#2573](https://github.com/vigilnetwork/vgl/pull/2573))
- multi: Round 5 prerel module release ver updates ([Vigil/vgld#2574](https://github.com/vigilnetwork/vgl/pull/2574))
- multi: Round 6 prerel module release ver updates ([Vigil/vgld#2575](https://github.com/vigilnetwork/vgl/pull/2575))
- multi: Update to siphash v1.2.2 ([Vigil/vgld#2577](https://github.com/vigilnetwork/vgl/pull/2577))
- peer: Start v3 module dev cycle ([Vigil/vgld#2585](https://github.com/vigilnetwork/vgl/pull/2585))
- addrmgr: Start v2 module dev cycle ([Vigil/vgld#2592](https://github.com/vigilnetwork/vgl/pull/2592))
- blockchain: Prerel module release ver updates ([Vigil/vgld#2634](https://github.com/vigilnetwork/vgl/pull/2634))
- blockchain: Bump database module minor version ([Vigil/vgld#2654](https://github.com/vigilnetwork/vgl/pull/2654))
- multi: Require last database/v2.0.3-x version ([Vigil/vgld#2689](https://github.com/vigilnetwork/vgl/pull/2689))
- multi: Introduce database/v3 module ([Vigil/vgld#2689](https://github.com/vigilnetwork/vgl/pull/2689))
- multi: Use database/v3 module ([Vigil/vgld#2693](https://github.com/vigilnetwork/vgl/pull/2693))
- main: Use pseudo-versions in bumped mods ([Vigil/vgld#2698](https://github.com/vigilnetwork/vgl/pull/2698))
- blockchain: Add replace to chaincfg dependency ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- VGLjson: Introduce v4 module ([Vigil/vgld#2733](https://github.com/vigilnetwork/vgl/pull/2733))
- secp256k1: Prepare v4.0.0 ([Vigil/vgld#2732](https://github.com/vigilnetwork/vgl/pull/2732))
- docs: Update for VGLjson v4 module ([Vigil/vgld#2734](https://github.com/vigilnetwork/vgl/pull/2734))
- VGLjson: Prepare v4.0.0 ([Vigil/vgld#2734](https://github.com/vigilnetwork/vgl/pull/2734))
- blockchain: Prerel module release ver updates ([Vigil/vgld#2748](https://github.com/vigilnetwork/vgl/pull/2748))
- gcs: Prerel module release ver updates ([Vigil/vgld#2749](https://github.com/vigilnetwork/vgl/pull/2749))
- multi: Update gcs prerel version ([Vigil/vgld#2750](https://github.com/vigilnetwork/vgl/pull/2750))
- multi: update build tags to pref. go1.17 syntax ([Vigil/vgld#2764](https://github.com/vigilnetwork/vgl/pull/2764))
- chaincfg: Prepare v3.1.0 ([Vigil/vgld#2799](https://github.com/vigilnetwork/vgl/pull/2799))
- addrmgr: Prepare v2.0.0 ([Vigil/vgld#2797](https://github.com/vigilnetwork/vgl/pull/2797))
- rpc/jsonrpc/types: Prepare v3.0.0 ([Vigil/vgld#2801](https://github.com/vigilnetwork/vgl/pull/2801))
- txscript: Prepare v4.0.0 ([Vigil/vgld#2815](https://github.com/vigilnetwork/vgl/pull/2815))
- hdkeychain: Prepare v3.0.1 ([Vigil/vgld#2817](https://github.com/vigilnetwork/vgl/pull/2817))
- VGLutil: Prepare v4.0.0 ([Vigil/vgld#2818](https://github.com/vigilnetwork/vgl/pull/2818))
- connmgr: Prepare v3.1.0 ([Vigil/vgld#2819](https://github.com/vigilnetwork/vgl/pull/2819))
- peer: Prepare v3.0.0 ([Vigil/vgld#2820](https://github.com/vigilnetwork/vgl/pull/2820))
- database: Prepare v3.0.0 ([Vigil/vgld#2822](https://github.com/vigilnetwork/vgl/pull/2822))
- blockchain/stake: Prepare v4.0.0 ([Vigil/vgld#2824](https://github.com/vigilnetwork/vgl/pull/2824))
- gcs: Prepare v3.0.0 ([Vigil/vgld#2830](https://github.com/vigilnetwork/vgl/pull/2830))
- math/uint256: Prepare v1.0.0 ([Vigil/vgld#2842](https://github.com/vigilnetwork/vgl/pull/2842))
- blockchain: Prepare v4.0.0 ([Vigil/vgld#2831](https://github.com/vigilnetwork/vgl/pull/2831))
- rpcclient: Prepare v7.0.0 ([Vigil/vgld#2851](https://github.com/vigilnetwork/vgl/pull/2851))
- version: Include VCS build info in version string ([Vigil/vgld#2841](https://github.com/vigilnetwork/vgl/pull/2841))
- main: Update to use all new module versions ([Vigil/vgld#2853](https://github.com/vigilnetwork/vgl/pull/2853))
- main: Remove module replacements ([Vigil/vgld#2855](https://github.com/vigilnetwork/vgl/pull/2855))

### Testing and Quality Assurance:

- rpcserver: Add handleGetTreasuryBalance tests ([Vigil/vgld#2390](https://github.com/vigilnetwork/vgl/pull/2390))
- rpcserver: Add handleGet{Generate,HashesPerSec} tests ([Vigil/vgld#2365](https://github.com/vigilnetwork/vgl/pull/2365))
- mining: Cleanup txPriorityQueue tests ([Vigil/vgld#2431](https://github.com/vigilnetwork/vgl/pull/2431))
- blockchain: fix errorlint warnings ([Vigil/vgld#2411](https://github.com/vigilnetwork/vgl/pull/2411))
- rpcserver: Add handleGetHeaders test ([Vigil/vgld#2366](https://github.com/vigilnetwork/vgl/pull/2366))
- rpcserver: add ticketsforaddress tests ([Vigil/vgld#2405](https://github.com/vigilnetwork/vgl/pull/2405))
- rpcserver: add ticketvwap tests ([Vigil/vgld#2406](https://github.com/vigilnetwork/vgl/pull/2406))
- rpcserver: add handleTxFeeInfo tests ([Vigil/vgld#2407](https://github.com/vigilnetwork/vgl/pull/2407))
- rpcserver: add handleTicketFeeInfo tests ([Vigil/vgld#2408](https://github.com/vigilnetwork/vgl/pull/2408))
- rpcserver: add handleVerifyMessage tests ([Vigil/vgld#2413](https://github.com/vigilnetwork/vgl/pull/2413))
- rpcserver: add handleSendRawTransaction tests ([Vigil/vgld#2410](https://github.com/vigilnetwork/vgl/pull/2410))
- rpcserver: add handleGetVoteInfo tests ([Vigil/vgld#2432](https://github.com/vigilnetwork/vgl/pull/2432))
- database: Fix errorlint warnings ([Vigil/vgld#2484](https://github.com/vigilnetwork/vgl/pull/2484))
- mining: Add mining test harness ([Vigil/vgld#2480](https://github.com/vigilnetwork/vgl/pull/2480))
- mining: Add NewBlockTemplate tests ([Vigil/vgld#2480](https://github.com/vigilnetwork/vgl/pull/2480))
- mining: Move TxMiningView tests to mining ([Vigil/vgld#2480](https://github.com/vigilnetwork/vgl/pull/2480))
- rpcserver: add handleGetRawTransaction tests ([Vigil/vgld#2483](https://github.com/vigilnetwork/vgl/pull/2483))
- blockchain: Improve synthetic treasury vote tests ([Vigil/vgld#2488](https://github.com/vigilnetwork/vgl/pull/2488))
- rpcserver: Add handleGetMempoolInfo test ([Vigil/vgld#2492](https://github.com/vigilnetwork/vgl/pull/2492))
- connmgr: Increase test timeouts ([Vigil/vgld#2505](https://github.com/vigilnetwork/vgl/pull/2505))
- run_vgl_tests.sh: Avoid command substitution ([Vigil/vgld#2506](https://github.com/vigilnetwork/vgl/pull/2506))
- mempool: Make sequence lock tests more consistent ([Vigil/vgld#2496](https://github.com/vigilnetwork/vgl/pull/2496))
- mempool: Rework sequence lock acceptance tests ([Vigil/vgld#2496](https://github.com/vigilnetwork/vgl/pull/2496))
- rpcserver: Add handleGetTxOut tests ([Vigil/vgld#2516](https://github.com/vigilnetwork/vgl/pull/2516))
- rpcserver: Add handleGetNetworkHashPS test ([Vigil/vgld#2512](https://github.com/vigilnetwork/vgl/pull/2512))
- rpcserver: Add handleGetMiningInfo test ([Vigil/vgld#2512](https://github.com/vigilnetwork/vgl/pull/2512))
- blockchain: Simplify TestFixedSequenceLocks ([Vigil/vgld#2534](https://github.com/vigilnetwork/vgl/pull/2534))
- chaingen: Support querying block test name by hash ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- blockchain: Improve test harness logging ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- blockchain: Support separate test block generation ([Vigil/vgld#2518](https://github.com/vigilnetwork/vgl/pull/2518))
- rpcserver: add handleVersion, handleHelp rpc tests ([Vigil/vgld#2549](https://github.com/vigilnetwork/vgl/pull/2549))
- blockchain: Use ReplaceVoteBits in utxoview tests ([Vigil/vgld#2553](https://github.com/vigilnetwork/vgl/pull/2553))
- blockchain: Add unit test coverage for UtxoEntry ([Vigil/vgld#2553](https://github.com/vigilnetwork/vgl/pull/2553))
- rpctest: Don't use installed node ([Vigil/vgld#2523](https://github.com/vigilnetwork/vgl/pull/2523))
- apbf: Add comprehensive tests ([Vigil/vgld#2579](https://github.com/vigilnetwork/vgl/pull/2579))
- apbf: Add benchmarks ([Vigil/vgld#2579](https://github.com/vigilnetwork/vgl/pull/2579))
- rpcserver: Add handleGetRawMempool test ([Vigil/vgld#2589](https://github.com/vigilnetwork/vgl/pull/2589))
- build: Test against go 1.16 ([Vigil/vgld#2598](https://github.com/vigilnetwork/vgl/pull/2598))
- blockchain: Add test name to TestUtxoEntry errors ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Add UtxoCache test coverage ([Vigil/vgld#2591](https://github.com/vigilnetwork/vgl/pull/2591))
- blockchain: Use new style for chainio test errors ([Vigil/vgld#2595](https://github.com/vigilnetwork/vgl/pull/2595))
- rpcserver: Add handleInvalidateBlock test ([Vigil/vgld#2604](https://github.com/vigilnetwork/vgl/pull/2604))
- blockchain: Mock time.Now for utxo cache tests ([Vigil/vgld#2605](https://github.com/vigilnetwork/vgl/pull/2605))
- blockchain: Add UtxoCache Initialize tests ([Vigil/vgld#2599](https://github.com/vigilnetwork/vgl/pull/2599))
- blockchain: Add TestShutdownUtxoCache tests ([Vigil/vgld#2599](https://github.com/vigilnetwork/vgl/pull/2599))
- rpcserver: Add handleReconsiderBlock test ([Vigil/vgld#2613](https://github.com/vigilnetwork/vgl/pull/2613))
- stdaddr: Add benchmarks ([Vigil/vgld#2610](https://github.com/vigilnetwork/vgl/pull/2610))
- rpctest: Make tests work properly with latest code ([Vigil/vgld#2614](https://github.com/vigilnetwork/vgl/pull/2614))
- mempool: Remove unused field from test struct ([Vigil/vgld#2618](https://github.com/vigilnetwork/vgl/pull/2618))
- mempool: Remove unused func from tests ([Vigil/vgld#2621](https://github.com/vigilnetwork/vgl/pull/2621))
- rpctest: Don't use Fatalf in goroutines ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- chaingen: Remove unused PurchaseCommitmentScript ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- rpctest: Convert to use new stdaddr package ([Vigil/vgld#2625](https://github.com/vigilnetwork/vgl/pull/2625))
- VGLutil: Move address params iface and mock impls ([Vigil/vgld#2628](https://github.com/vigilnetwork/vgl/pull/2628))
- stdscript: Add v0 p2pk-ecdsa-secp256k1 benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pk-ed25519 benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pk-schnorr-secp256k1 benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pkh-ecdsa-secp256k1 benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pkh-ed25519 benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2pkh-schnorr-secp256k1 benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 p2sh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 ecdsa multisig benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 ecdsa multisig redeem benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add extract v0 multisig redeem benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 nulldata benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake sub p2pkh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake sub p2sh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake gen p2pkh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake gen p2sh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake revoke p2pkh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake revoke p2sh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake change p2pkh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 stake change p2sh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 treasury add benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 treasury gen p2pkh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 treasury gen p2sh benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add determine script type benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- stdscript: Add v0 atomic swap redeem benchmark ([Vigil/vgld#2656](https://github.com/vigilnetwork/vgl/pull/2656))
- txscript: Separate short form script parsing ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Explicit consensus p2sh tests ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- txscript: Explicit consensus any kind p2sh tests ([Vigil/vgld#2666](https://github.com/vigilnetwork/vgl/pull/2666))
- stake: No standardness code in tests ([Vigil/vgld#2667](https://github.com/vigilnetwork/vgl/pull/2667))
- blockchain: Add outpointKey tests ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- blockchain: Add block index key collision tests ([Vigil/vgld#2685](https://github.com/vigilnetwork/vgl/pull/2685))
- secp256k1: Rework NAF tests ([Vigil/vgld#2695](https://github.com/vigilnetwork/vgl/pull/2695))
- secp256k1: Cleanup NAF benchmark ([Vigil/vgld#2695](https://github.com/vigilnetwork/vgl/pull/2695))
- rpctest: Add P2PAddress() function ([Vigil/vgld#2704](https://github.com/vigilnetwork/vgl/pull/2704))
- tests: Remove hardcoded CC=gcc from run_vgl_tests.sh ([Vigil/vgld#2706](https://github.com/vigilnetwork/vgl/pull/2706))
- build: Test against Go 1.17 ([Vigil/vgld#2712](https://github.com/vigilnetwork/vgl/pull/2712))
- blockchain: Support voting multiple agendas in test ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- blockchain: Single out treasury policy test ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- blockchain: Correct test harness err msg ([Vigil/vgld#2714](https://github.com/vigilnetwork/vgl/pull/2714))
- blockchain: Test new max expenditure policy ([Vigil/vgld#2680](https://github.com/vigilnetwork/vgl/pull/2680))
- chaingen: Add spendable coinbase out snapshots ([Vigil/vgld#2715](https://github.com/vigilnetwork/vgl/pull/2715))
- mempool: Accept test mungers for create tickets ([Vigil/vgld#2721](https://github.com/vigilnetwork/vgl/pull/2721))
- build: Don't set GO111MODULE unnecessarily ([Vigil/vgld#2722](https://github.com/vigilnetwork/vgl/pull/2722))
- build: Don't manually test changing go.{mod,sum} ([Vigil/vgld#2722](https://github.com/vigilnetwork/vgl/pull/2722))
- stake: Add CalculateRewards tests ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- stake: Add CheckSSRtx tests ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- blockchain: Test auto revocations deployment ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- chaingen: Add revocation mungers ([Vigil/vgld#2718](https://github.com/vigilnetwork/vgl/pull/2718))
- addrmgr: Improve test coverage ([Vigil/vgld#2596](https://github.com/vigilnetwork/vgl/pull/2596))
- addrmgr: Remove unnecessary test cases ([Vigil/vgld#2596](https://github.com/vigilnetwork/vgl/pull/2596))
- rpcserver: Tune large tspend test amount ([Vigil/vgld#2679](https://github.com/vigilnetwork/vgl/pull/2679))
- build: Pin GitHub Actions to SHA ([Vigil/vgld#2736](https://github.com/vigilnetwork/vgl/pull/2736))
- blockchain: Add calcTicketReturnAmounts tests ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- blockchain: Add checkTicketRedeemers tests ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- blockchain: Add auto revocation validation tests ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- mining: Add auto revocation block template tests ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- mempool: Add tests with auto revocations enabled ([Vigil/vgld#2720](https://github.com/vigilnetwork/vgl/pull/2720))
- txscript: Add versioned short form parsing ([Vigil/vgld#2756](https://github.com/vigilnetwork/vgl/pull/2756))
- txscript: Test consistency and cleanup ([Vigil/vgld#2757](https://github.com/vigilnetwork/vgl/pull/2757))
- mempool: Add blockHeight to AddFakeUTXO for tests ([Vigil/vgld#2804](https://github.com/vigilnetwork/vgl/pull/2804))
- mempool: Test fraud proof handling ([Vigil/vgld#2804](https://github.com/vigilnetwork/vgl/pull/2804))
- stdscript: Add extract v0 stake-tagged p2pkh bench ([Vigil/vgld#2816](https://github.com/vigilnetwork/vgl/pull/2816))
- stdscript: Add extract v0 stake-tagged p2sh bench ([Vigil/vgld#2816](https://github.com/vigilnetwork/vgl/pull/2816))
- mempool: Update test to check hash value ([Vigil/vgld#2804](https://github.com/vigilnetwork/vgl/pull/2804))
- stdscript: Add num required sigs benchmark ([Vigil/vgld#2805](https://github.com/vigilnetwork/vgl/pull/2805))
- uint256: Add big endian set benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add little endian set benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add big endian get benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add little endian get benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add zero benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add equality comparison benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add less than comparison benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add greater than comparison benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add general comparison benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add addition benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add subtraction benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add multiplication benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add squaring benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add division benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add negation benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add is odd benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise left shift benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise right shift benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise not benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise or benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise and benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bitwise xor benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add bit length benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add text formatting benchmarks ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add conversion to stdlib big int benchmark ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- uint256: Add conversion from stdlib big int benchmark ([Vigil/vgld#2787](https://github.com/vigilnetwork/vgl/pull/2787))
- primitives: Add diff bits conversion benchmarks ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add work calc benchmark ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add hash to uint256 benchmark ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add check proof of work benchmark ([Vigil/vgld#2788](https://github.com/vigilnetwork/vgl/pull/2788))
- primitives: Add merkle root benchmarks ([Vigil/vgld#2826](https://github.com/vigilnetwork/vgl/pull/2826))
- primitives: Add inclusion proof benchmarks ([Vigil/vgld#2827](https://github.com/vigilnetwork/vgl/pull/2827))
- blockchain: Add AssumeValid tests ([Vigil/vgld#2839](https://github.com/vigilnetwork/vgl/pull/2839))
- chaingen: Add vote subsidy munger ([Vigil/vgld#2848](https://github.com/vigilnetwork/vgl/pull/2848))

### Misc:

- release: Bump for 1.7 release cycle ([Vigil/vgld#2429](https://github.com/vigilnetwork/vgl/pull/2429))
- secp256k1: Correct const name for doc comment ([Vigil/vgld#2445](https://github.com/vigilnetwork/vgl/pull/2445))
- multi: Fix various typos ([Vigil/vgld#2607](https://github.com/vigilnetwork/vgl/pull/2607))
- rpcserver: Fix createrawssrtx comments ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- blockchain: Fix comment formatting in generator ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- stake: Fix MaxOutputsPerSSRtx comment ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- stake: Fix CheckSSGenVotes function comment ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- stake: Fix CheckSSRtx function comment ([Vigil/vgld#2665](https://github.com/vigilnetwork/vgl/pull/2665))
- database: Add comment on os.MkdirAll behavior ([Vigil/vgld#2670](https://github.com/vigilnetwork/vgl/pull/2670))
- multi: Address some linter complaints ([Vigil/vgld#2684](https://github.com/vigilnetwork/vgl/pull/2684))
- txscript: Fix a couple of a comment typos ([Vigil/vgld#2692](https://github.com/vigilnetwork/vgl/pull/2692))
- blockchain: Remove inapplicable comment ([Vigil/vgld#2742](https://github.com/vigilnetwork/vgl/pull/2742))
- mining: Fix error in comment ([Vigil/vgld#2743](https://github.com/vigilnetwork/vgl/pull/2743))
- blockchain: Fix several typos ([Vigil/vgld#2745](https://github.com/vigilnetwork/vgl/pull/2745))
- blockchain: Update a few BFFastAdd comments ([Vigil/vgld#2781](https://github.com/vigilnetwork/vgl/pull/2781))
- multi: Address some linter complaints ([Vigil/vgld#2791](https://github.com/vigilnetwork/vgl/pull/2791))
- netsync: Correct typo ([Vigil/vgld#2813](https://github.com/vigilnetwork/vgl/pull/2813))
- netsync: Fix misc typos ([Vigil/vgld#2834](https://github.com/vigilnetwork/vgl/pull/2834))
- mining: Fix typo ([Vigil/vgld#2834](https://github.com/vigilnetwork/vgl/pull/2834))
- blockchain: Correct comment typos for find fork ([Vigil/vgld#2828](https://github.com/vigilnetwork/vgl/pull/2828))
- rpcserver: Rename var to make linter happy ([Vigil/vgld#2835](https://github.com/vigilnetwork/vgl/pull/2835))
- blockchain: Wrap at max line length ([Vigil/vgld#2843](https://github.com/vigilnetwork/vgl/pull/2843))
- release: Bump for 1.7.0 ([Vigil/vgld#2856](https://github.com/vigilnetwork/vgl/pull/2856))

### Code Contributors (alphabetical order):

- briancolecoinmetrics
- Dave Collins
- David Hill
- degeri
- Donald Adu-Poku
- J Fixby
- Jamie Holdstock
- Joe Gruffins
- Jonathan Chappelow
- Josh Rickmar
- lolandhold
- Matheus Degiovani
- Naveen
- Ryan Staudt
- Youssef Boukenken
- Wisdom Arerosuoghene




