# vgld v1.5.0

This release of vgld introduces a large number of updates.  Some of the key highlights are:

* A new consensus vote agenda which allows the stakeholders to decide whether or not to activate support for block header commitments
* More efficient block filters
* Significant improvements to the mining infrastructure including asynchronous work notifications
* Major performance enhancements for transaction script validation
* Automatic external IP address discovery
* Support for IPv6 over Tor
* Various updates to the RPC server such as:
  * A new method to query information about the network
  * A method to retrieve the new version 2 block filters
  * More calls available to limited access users
* Infrastructure improvements
* Quality assurance changes

For those unfamiliar with the voting process in Vigil, all code in order to support block header commitments is already included in this release, however its enforcement will remain dormant until the stakeholders vote to activate it.

For reference, block header commitments were originally proposed and approved for initial implementation via the following Vigiliteia proposal:
- [Block Header Commitments Consensus Change](https://proposals.vigil.network/proposals/0a1ff846ec271184ea4e3a921a3ccd8d478f69948b984445ee1852f272d54c58)


The following Vigil Change Proposal (VGLP) describes the proposed changes in detail and provides a full technical specification:
- [VGLP0005](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki)

**It is important for everyone to upgrade their software to this latest release even if you don't intend to vote in favor of the agenda.**

## Downgrade Warning

The database format in v1.5.0 is not compatible with previous versions of the software.  This only affects downgrades as users upgrading from previous versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to downgrade to a previous version of the software without having to delete the database and redownload the chain.

## Notable Changes

### Block Header Commitments Vote

A new vote with the id `headercommitments` is now available as of this release.  After upgrading, stakeholders may set their preferences through their wallet or Voting Service Provider's (VSP) website.

The primary goal of this change is to increase the security and efficiency of lightweight clients, such as Vigiliton in its lightweight mode and the VGLandroid/VGLios mobile wallets, as well as add infrastructure that paves the
way for several future scalability enhancements.

A high level overview aimed at a general audience including a cost benefit analysis can be found in the  [Vigiliteia proposal](https://proposals.vigil.network/proposals/0a1ff846ec271184ea4e3a921a3ccd8d478f69948b984445ee1852f272d54c58).

In addition, a much more in-depth treatment can be found in the [motivation section of VGLP0005](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#motivation).

### Version 2 Block Filters

The block filters used by lightweight clients, such as SPV (Simplified Payment Verification) wallets, have been updated to improve their efficiency, ergonomics, and include additional information such as the full ticket
commitment script.  The new block filters are version 2.  The older version 1 filters are now deprecated and scheduled to be removed in the next release, so consumers should update to the new filters as soon as possible.

An overview of block filters can be found in the [block filters section of VGLP0005](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#block-filters).

Also, the specific contents and technical specification of the new version 2 block filters is available in the
[version 2 block filters section of VGLP0005](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#version-2-block-filters).

Finally, there is a one time database update to build and store the new filters for all existing historical blocks which will likely take a while to complete (typically around 8 to 10 minutes on HDDs and 4 to 5 minutes on SSDs).

### Mining Infrastructure Overhaul

The mining infrastructure for building block templates and delivering the work to miners has been significantly overhauled to improve several aspects as follows:

* Support asynchronous background template generation with intelligent vote propagation handling
* Improved handling of chain reorganizations necessary when the current tip is unable to obtain enough votes
* Current state synchronization
* Near elimination of stale templates when new blocks and votes are received
* Subscriptions for streaming template updates

The standard [getwork RPC](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getwork) that PoW miners currently use to perform the mining process has been updated to make use of this new infrastructure, so existing PoW miners will seamlessly get the vast majority of benefits without requiring any updates.

However, in addition, a new [notifywork RPC](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#notifywork) is now available that allows miners to register for work to be delivered
asynchronously as it becomes available via a WebSockets [work notification](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#work).  These notifications include the same information that `getwork` provides along with an additional `reason` parameter which allows the miners to make better decisions about when they should instruct workers to discard the current template immediately or should be allowed to finish their current round before being provided with the new template.

Miners are highly encouraged to update their software to make use of the new asynchronous notification infrastructure since it is more robust, efficient, and faster than polling `getwork` to manually determine the aforementioned conditions.

The following is a non-exhaustive overview that highlights the major benefits of the changes for both cases:

- Requests for updated templates during the normal mining process in between tip   changes will now be nearly instant instead of potentially taking several seconds to build the new template on the spot
- When the chain tip changes, requesting a template will now attempt to wait until either all votes have been received or a timeout occurs prior to handing out a template which is beneficial for PoW miners, PoS miners, and the network as a whole
- PoW miners are much less likely to end up with template with less than the max number of votes which means they are less likely to receive a reduced subsidy
- PoW miners will be much less likely to receive stale templates during chain tip changes due to vote propagation
- PoS voters whose votes end up arriving to the miner slightly slower than the minimum number required are much less likely to have their votes excluded despite having voted simply due to propagation delay

PoW miners who choose to update their software, pool or otherwise, to make use of the asynchronous work notifications will receive additional benefits such as:

- Ability to start mining a new block sooner due to receiving updated work as soon as it becomes available
- Immediate notification with new work that includes any votes that arrive late
- Periodic notifications with new work that include new transactions only when there have actually been new transaction
- Simplified interface code due to removal of the need for polling and manually checking the work bytes for special cases such as the number of votes

**NOTE: Miners that are not rolling the timestamp field as they mine should ensure their software is upgraded to roll the timestamp to the latest timestamp each time they hand work out to a miner.  This helps ensure the block timestamps are as accurate as possible.**

### Transaction Script Validation Optimizations

Transaction script validation has been almost completely rewritten to significantly improve its speed and reduce the number of memory allocations. While this has many more benefits than enumerated here, probably the most
important ones for most stakeholders are:

- Votes can be cast more quickly which helps reduce the number of missed votes
- Blocks are able to propagate more quickly throughout the network, which in turn further improves votes times
- The initial sync process is around 20-25% faster

### Automatic External IP Address Discovery

In order for nodes to fully participate in the peer-to-peer network, they must be publicly accessible and made discoverable by advertising their external IP address.  This is typically made slightly more complicated since most users run their nodes on networks behind Network Address Translation (NAT).

Previously, in addition to configuring the network firewall and/or router to allow inbound connections to port 9108 and forwarding the port to the internal IP address running vgld, it was also required to manually set the public external IP address via the `--externalip` CLI option.

This release will now make use of other nodes on the network in a decentralized fashion to automatically discover the external IP address, so it is no longer necessary to manually set CLI option for the vast majority of users.

### Tor IPv6 Support

It is now possible to resolve and connect to IPv6 peers over Tor in addition to the existing IPv4 support.

### RPC Server Changes

#### New Version 2 Block Filter Query RPC (`getcfilterv2`)

A new RPC named `getcfilterv2` is now available which can be used to retrieve the version 2 [block filter](https://github.com/Vigil/VGLPs/blob/master/VGLP-0005/VGLP-0005.mediawiki#Block_Filters)
for a given block along with its associated inclusion proof.  See the [getcfilterv2 JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getcfilterv2)
for API details.

#### New Network Information Query RPC (`getnetworkinfo`)

A new RPC named `getnetworkinfo` is now available which can be used to query information related to the peer-to-peer network such as the protocol version, the local time offset, the number of current connections, the supported network protocols, the current transaction relay fee, and the external IP addresses for
the local interfaces.  See the [getnetworkinfo JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getnetworkinfo) for API details.

#### Updates to Chain State Query RPC (`getblockchaininfo`)

The `difficulty` field of the `getblockchaininfo` RPC is now deprecated in favor of a new field named `difficultyratio` which matches the result returned by the `getdifficulty` RPC.

See the [getblockchaininfo JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblockchaininfo) for API details.

#### New Optional Version Parameter on Script Decode RPC (`decodescript`)

The `decodescript` RPC now accepts an additional optional parameter to specify the script version.  The only currently supported script version in Vigil is version 0 which means decoding scripts with versions other than 0 will be seen as non standard.

#### Removal of Deprecated Block Template RPC (`getblocktemplate`)

The previously deprecated `getblocktemplate` RPC is no longer available.  All known miners are already using the preferred `getwork` RPC since Vigil's block header supports more than enough nonce space to keep mining hardware busy without needing to resort to building custom templates with less efficient extra nonce coinbase workarounds.

#### Additional RPCs Available To Limited Access Users

The following RPCs that were previously unavailable to the limited access RPC user are now available to it:

- `estimatefee`
- `estimatesmartfee`
- `estimatestakediff`
- `existsaddress`
- `existsaddresses`
- `existsexpiredtickets`
- `existsliveticket`
- `existslivetickets`
- `existsmempoltxs`
- `existsmissedtickets`
- `getblocksubsidy`
- `getcfilter`
- `getcoinsupply`
- `getheaders`
- `getstakedifficulty`
- `getstakeversioninfo`
- `getstakeversions`
- `getvoteinfo`
- `livetickets`
- `missedtickets`
- `rebroadcastmissed`
- `rebroadcastwinners`
- `ticketfeeinfo`
- `ticketsforaddress`
- `ticketvwap`
- `txfeeinfo`

### Single Mining State Request

The peer-to-peer protocol message to request the current mining state (`getminings`) is used when peers first connect to retrieve all known votes for the current tip block.  This is only useful when the peer first connects because all future votes will be relayed once the connection has been established.  Consequently, nodes will now only respond to a single mining state request.  Subsequent requests are ignored.

### Developer Go Modules

A full suite of versioned Go modules (essentially code libraries) are now available for use by applications written in Go that wish to create robust software with reproducible, verifiable, and verified builds.

These modules are used to build vgld itself and are therefore well maintained, tested, documented, and relatively efficient.

## Changelog

This release consists of 600 commits from 17 contributors which total to 537 files changed, 41494 additional lines of code, and 29215 deleted lines of code.

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/release-v1.4.0...release-v1.5.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.5.0 release ([Vigil/vgld#1924](https://github.com/vigilnetwork/vgl/pull/1924))
- chaincfg: Introduce agenda for header cmtmts vote ([Vigil/vgld#1904](https://github.com/vigilnetwork/vgl/pull/1904))
- multi: Implement combined merkle root and vote ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- blockchain: Implement v2 block filter storage ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- gcs/blockcf2: Implement v2 block filter creation ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- wire: Implement getcfilterv2/cfilterv2 messages ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- peer: Implement getcfilterv2/cfilterv2 listeners ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- server: Implement getcfilterv2 ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- multi: Implement header commitments and vote ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- server: Remove instead of disconnect node ([Vigil/vgld#1644](https://github.com/vigilnetwork/vgl/pull/1644))
- server: limit getminingstate requests ([Vigil/vgld#1678](https://github.com/vigilnetwork/vgl/pull/1678))
- peer: Prevent last block height going backwards ([Vigil/vgld#1747](https://github.com/vigilnetwork/vgl/pull/1747))
- connmgr: Add ability to remove pending connections ([Vigil/vgld#1724](https://github.com/vigilnetwork/vgl/pull/1724))
- connmgr: Add cancellation of pending requests ([Vigil/vgld#1724](https://github.com/vigilnetwork/vgl/pull/1724))
- connmgr: Check for canceled connection before connect ([Vigil/vgld#1724](https://github.com/vigilnetwork/vgl/pull/1724))
- multi: add automatic network address discovery ([Vigil/vgld#1522](https://github.com/vigilnetwork/vgl/pull/1522))
- connmgr: add TorLookupIPContext, deprecate TorLookupIP ([Vigil/vgld#1849](https://github.com/vigilnetwork/vgl/pull/1849))
- connmgr: support resolving ipv6 hosts over Tor ([Vigil/vgld#1908](https://github.com/vigilnetwork/vgl/pull/1908))

### Transaction relay (memory pool):

- mempool: Reject same block vote double spends ([Vigil/vgld#1597](https://github.com/vigilnetwork/vgl/pull/1597))
- mempool: Limit max vote double spends exactly ([Vigil/vgld#1596](https://github.com/vigilnetwork/vgl/pull/1596))
- mempool: Optimize pool double spend check ([Vigil/vgld#1561](https://github.com/vigilnetwork/vgl/pull/1561))
- txscript: Tighten standardness pubkey checks ([Vigil/vgld#1649](https://github.com/vigilnetwork/vgl/pull/1649))
- mempool: drop container/list for simple FIFO ([Vigil/vgld#1681](https://github.com/vigilnetwork/vgl/pull/1681))
- mempool: remove unused error return value ([Vigil/vgld#1785](https://github.com/vigilnetwork/vgl/pull/1785))
- mempool: Add ErrorCode to returned TxRuleErrors ([Vigil/vgld#1901](https://github.com/vigilnetwork/vgl/pull/1901))

### Mining:

- mining: Optimize get the block's votes tx ([Vigil/vgld#1563](https://github.com/vigilnetwork/vgl/pull/1563))
- multi: add BgBlkTmplGenerator ([Vigil/vgld#1424](https://github.com/vigilnetwork/vgl/pull/1424))
- mining: Remove unnecessary notify goroutine ([Vigil/vgld#1708](https://github.com/vigilnetwork/vgl/pull/1708))
- mining: Improve template key handling ([Vigil/vgld#1709](https://github.com/vigilnetwork/vgl/pull/1709))
- mining:  fix scheduled template regen ([Vigil/vgld#1717](https://github.com/vigilnetwork/vgl/pull/1717))
- miner: Improve background generator lifecycle ([Vigil/vgld#1715](https://github.com/vigilnetwork/vgl/pull/1715))
- cpuminer: No speed monitor on discrete mining ([Vigil/vgld#1716](https://github.com/vigilnetwork/vgl/pull/1716))
- mining: Run vote ntfn in a separate goroutine ([Vigil/vgld#1718](https://github.com/vigilnetwork/vgl/pull/1718))
- mining: Overhaul background template generator ([Vigil/vgld#1748](https://github.com/vigilnetwork/vgl/pull/1748))
- mining: Remove unused error return value ([Vigil/vgld#1859](https://github.com/vigilnetwork/vgl/pull/1859))
- cpuminer: Fix off-by-one issues in nonce handling ([Vigil/vgld#1865](https://github.com/vigilnetwork/vgl/pull/1865))
- mining: Remove dead code ([Vigil/vgld#1882](https://github.com/vigilnetwork/vgl/pull/1882))
- mining: Remove unused extra nonce update code ([Vigil/vgld#1883](https://github.com/vigilnetwork/vgl/pull/1883))
- mining: Minor cleanup of aggressive mining path ([Vigil/vgld#1888](https://github.com/vigilnetwork/vgl/pull/1888))
- mining: Remove unused error codes ([Vigil/vgld#1889](https://github.com/vigilnetwork/vgl/pull/1889))
- mining: fix data race ([Vigil/vgld#1894](https://github.com/vigilnetwork/vgl/pull/1894))
- mining: fix data race ([Vigil/vgld#1896](https://github.com/vigilnetwork/vgl/pull/1896))
- cpuminer: fix race ([Vigil/vgld#1899](https://github.com/vigilnetwork/vgl/pull/1899))
- cpuminer: Improve speed stat tracking ([Vigil/vgld#1921](https://github.com/vigilnetwork/vgl/pull/1921))
- rpcserver/mining: Use bg tpl generator for getwork ([Vigil/vgld#1922](https://github.com/vigilnetwork/vgl/pull/1922))
- mining: Export TemplateUpdateReason ([Vigil/vgld#1923](https://github.com/vigilnetwork/vgl/pull/1923))
- multi: Add tpl update reason to work ntfns ([Vigil/vgld#1923](https://github.com/vigilnetwork/vgl/pull/1923))
- mining: Store block templates given by notifywork ([Vigil/vgld#1949](https://github.com/vigilnetwork/vgl/pull/1949))

### RPC:

- VGLjson: add cointype to WalletInfoResult ([Vigil/vgld#1606](https://github.com/vigilnetwork/vgl/pull/1606))
- rpcclient: Introduce v2 module using wallet types ([Vigil/vgld#1608](https://github.com/vigilnetwork/vgl/pull/1608))
- rpcserver: Update for VGLjson/v2 ([Vigil/vgld#1612](https://github.com/vigilnetwork/vgl/pull/1612))
- rpcclient: Add EstimateSmartFee ([Vigil/vgld#1641](https://github.com/vigilnetwork/vgl/pull/1641))
- rpcserver: remove unused quit chan ([Vigil/vgld#1629](https://github.com/vigilnetwork/vgl/pull/1629))
- rpcserver: Undeprecate getwork ([Vigil/vgld#1635](https://github.com/vigilnetwork/vgl/pull/1635))
- rpcserver: Add difficultyratio to getblockchaininfo ([Vigil/vgld#1630](https://github.com/vigilnetwork/vgl/pull/1630))
- multi:  add version arg to decodescript rpc ([Vigil/vgld#1731](https://github.com/vigilnetwork/vgl/pull/1731))
- VGLjson: Remove API breaking change ([Vigil/vgld#1778](https://github.com/vigilnetwork/vgl/pull/1778))
- rpcclient: Add GetMasterPubkey ([Vigil/vgld#1777](https://github.com/vigilnetwork/vgl/pull/1777))
- multi: add getnetworkinfo rpc ([Vigil/vgld#1536](https://github.com/vigilnetwork/vgl/pull/1536))
- rpcserver: Better error message ([Vigil/vgld#1861](https://github.com/vigilnetwork/vgl/pull/1861))
- multi: update limited user rpcs ([Vigil/vgld#1870](https://github.com/vigilnetwork/vgl/pull/1870))
- multi: make rebroadcast winners & missed ws only ([Vigil/vgld#1872](https://github.com/vigilnetwork/vgl/pull/1872))
- multi: remove getblocktemplate ([Vigil/vgld#1736](https://github.com/vigilnetwork/vgl/pull/1736))
- rpcserver: Match tx filter on ticket commitments ([Vigil/vgld#1881](https://github.com/vigilnetwork/vgl/pull/1881))
- rpcserver: don't use activeNetParams ([Vigil/vgld#1733](https://github.com/vigilnetwork/vgl/pull/1733))
- rpcserver: update rpcAskWallet rpc set ([Vigil/vgld#1892](https://github.com/vigilnetwork/vgl/pull/1892))
- rpcclient: close the unused response body ([Vigil/vgld#1905](https://github.com/vigilnetwork/vgl/pull/1905))
- rpcclient: Support getcfilterv2 JSON-RPC ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))
- multi: add notifywork rpc ([Vigil/vgld#1410](https://github.com/vigilnetwork/vgl/pull/1410))
- rpcserver: Cleanup getvoteinfo RPC ([Vigil/vgld#2005](https://github.com/vigilnetwork/vgl/pull/2005))

### vgld command-line flags and configuration:

- config: Remove deprecated getworkkey option ([Vigil/vgld#1594](https://github.com/vigilnetwork/vgl/pull/1594))

### certgen utility changes:

- certgen: Support Ed25519 cert generation on Go 1.13 ([Vigil/vgld#1757](https://github.com/vigilnetwork/vgl/pull/1757))

### vglctl utility changes:

- vglctl: Make version string consistent ([Vigil/vgld#1598](https://github.com/vigilnetwork/vgl/pull/1598))
- vglctl: Update for VGLjson/v2 and wallet types ([Vigil/vgld#1609](https://github.com/vigilnetwork/vgl/pull/1609))
- sampleconfig: add export vglctl sample config ([Vigil/vgld#2006](https://github.com/vigilnetwork/vgl/pull/2006))

### promptsecret utility changes:

- promptsecret: Add -n flag to prompt multiple times ([Vigil/vgld#1705](https://github.com/vigilnetwork/vgl/pull/1705))

### Documentation:

- docs: Update for secp256k1 v2 module ([Vigil/vgld#1919](https://github.com/vigilnetwork/vgl/pull/1919))
- docs: document module breaking changes process ([Vigil/vgld#1891](https://github.com/vigilnetwork/vgl/pull/1891))
- docs: Link to btc whitepaper on vigil.network ([Vigil/vgld#1885](https://github.com/vigilnetwork/vgl/pull/1885))
- docs: Update for mempool v3 module ([Vigil/vgld#1835](https://github.com/vigilnetwork/vgl/pull/1835))
- docs: Update for peer v2 module ([Vigil/vgld#1834](https://github.com/vigilnetwork/vgl/pull/1834))
- docs: Update for connmgr v2 module ([Vigil/vgld#1833](https://github.com/vigilnetwork/vgl/pull/1833))
- docs: Update for mining v2 module ([Vigil/vgld#1831](https://github.com/vigilnetwork/vgl/pull/1831))
- docs: Update for blockchain v2 module ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- docs: Update for rpcclient v4 module ([Vigil/vgld#1807](https://github.com/vigilnetwork/vgl/pull/1807))
- docs: Update for blockchain/stake v2 module ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- docs: Update for database v2 module ([Vigil/vgld#1799](https://github.com/vigilnetwork/vgl/pull/1799))
- docs: Update for rpcclient v3 module ([Vigil/vgld#1793](https://github.com/vigilnetwork/vgl/pull/1793))
- docs: Update for VGLjson/v3 module ([Vigil/vgld#1792](https://github.com/vigilnetwork/vgl/pull/1792))
- docs: Update for txscript v2 module ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- docs: Update for VGLutil v2 module ([Vigil/vgld#1770](https://github.com/vigilnetwork/vgl/pull/1770))
- docs: Update for VGLec/edwards v2 module ([Vigil/vgld#1765](https://github.com/vigilnetwork/vgl/pull/1765))
- docs: Update for chaincfg v2 module ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- docs: Update for hdkeychain v2 module ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Correct docs key examples ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- docs: allowHighFees arg has been implemented ([Vigil/vgld#1695](https://github.com/vigilnetwork/vgl/pull/1695))
- docs: move json rpc docs to mediawiki ([Vigil/vgld#1687](https://github.com/vigilnetwork/vgl/pull/1687))
- docs: Update for lru module ([Vigil/vgld#1683](https://github.com/vigilnetwork/vgl/pull/1683))
- docs: fix formatting in json rpc doc ([Vigil/vgld#1633](https://github.com/vigilnetwork/vgl/pull/1633))
- docs: Update for mempool v2 module ([Vigil/vgld#1613](https://github.com/vigilnetwork/vgl/pull/1613))
- docs: Update for rpcclient v2 module ([Vigil/vgld#1608](https://github.com/vigilnetwork/vgl/pull/1608))
- docs: Update for VGLjson v2 module ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- jsonrpc/types: Add README.md and doc.go ([Vigil/vgld#1794](https://github.com/vigilnetwork/vgl/pull/1794))
- VGLjson: Update README.md ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLec/secp256k1: Update README.md broken link ([Vigil/vgld#1631](https://github.com/vigilnetwork/vgl/pull/1631))
- bech32: Correct README build badge reference ([Vigil/vgld#1689](https://github.com/vigilnetwork/vgl/pull/1689))
- hdkeychain: Update README.md ([Vigil/vgld#1686](https://github.com/vigilnetwork/vgl/pull/1686))
- bech32: Correct README links ([Vigil/vgld#1691](https://github.com/vigilnetwork/vgl/pull/1691))
- stake: Remove unnecessary language in comment ([Vigil/vgld#1752](https://github.com/vigilnetwork/vgl/pull/1752))
- multi: Use https links where available ([Vigil/vgld#1771](https://github.com/vigilnetwork/vgl/pull/1771))
- stake: Make doc.go formatting consistent ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- blockchain: Update doc.go to reflect reality ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- multi: update rpc documentation ([Vigil/vgld#1867](https://github.com/vigilnetwork/vgl/pull/1867))
- VGLec: fix examples links ([Vigil/vgld#1914](https://github.com/vigilnetwork/vgl/pull/1914))
- gcs: Improve package documentation ([Vigil/vgld#1915](https://github.com/vigilnetwork/vgl/pull/1915))

### Developer-related package and module changes:

- VGLutil: Return deep copied tx in NewTxDeepTxIns ([Vigil/vgld#1545](https://github.com/vigilnetwork/vgl/pull/1545))
- mining: Remove superfluous error check ([Vigil/vgld#1552](https://github.com/vigilnetwork/vgl/pull/1552))
- VGLutil: Block does not cache the header bytes ([Vigil/vgld#1571](https://github.com/vigilnetwork/vgl/pull/1571))
- blockchain: Remove superfluous GetVoteInfo check ([Vigil/vgld#1574](https://github.com/vigilnetwork/vgl/pull/1574))
- blockchain: Make consensus votes network agnostic ([Vigil/vgld#1590](https://github.com/vigilnetwork/vgl/pull/1590))
- blockchain: Optimize skip stakebase input ([Vigil/vgld#1565](https://github.com/vigilnetwork/vgl/pull/1565))
- txscript: code cleanup ([Vigil/vgld#1591](https://github.com/vigilnetwork/vgl/pull/1591))
- VGLjson: Move estimate fee test to matching file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLjson: Move raw stake tx cmds to correct file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLjson: Move best block result to correct file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLjson: Move winning tickets ntfn to correct file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLjson: Move spent tickets ntfn to correct file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLjson: Move stake diff ntfn to correct file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- VGLjson: Move new tickets ntfn to correct file ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- txscript: Rename p2sh indicator to isP2SH ([Vigil/vgld#1605](https://github.com/vigilnetwork/vgl/pull/1605))
- mempool: Remove deprecated min high prio constant ([Vigil/vgld#1613](https://github.com/vigilnetwork/vgl/pull/1613))
- mempool: Remove tight coupling with VGLjson ([Vigil/vgld#1613](https://github.com/vigilnetwork/vgl/pull/1613))
- blockmanager: only check if current once handling inv's ([Vigil/vgld#1621](https://github.com/vigilnetwork/vgl/pull/1621))
- connmngr: Add DialAddr config option ([Vigil/vgld#1642](https://github.com/vigilnetwork/vgl/pull/1642))
- txscript: Consistent checksigaltverify handling ([Vigil/vgld#1647](https://github.com/vigilnetwork/vgl/pull/1647))
- multi: preallocate memory ([Vigil/vgld#1646](https://github.com/vigilnetwork/vgl/pull/1646))
- wire: Fix maximum payload length of MsgAddr ([Vigil/vgld#1638](https://github.com/vigilnetwork/vgl/pull/1638))
- blockmanager: remove unused requestedEverTxns ([Vigil/vgld#1624](https://github.com/vigilnetwork/vgl/pull/1624))
- blockmanager: remove useless requestedEverBlocks ([Vigil/vgld#1624](https://github.com/vigilnetwork/vgl/pull/1624))
- txscript: Introduce constant for max CLTV bytes ([Vigil/vgld#1650](https://github.com/vigilnetwork/vgl/pull/1650))
- txscript: Introduce constant for max CSV bytes ([Vigil/vgld#1651](https://github.com/vigilnetwork/vgl/pull/1651))
- chaincfg: Remove unused definition ([Vigil/vgld#1661](https://github.com/vigilnetwork/vgl/pull/1661))
- chaincfg: Use expected regnet merkle root var ([Vigil/vgld#1662](https://github.com/vigilnetwork/vgl/pull/1662))
- blockchain: Deprecate BlockOneCoinbasePaysTokens ([Vigil/vgld#1657](https://github.com/vigilnetwork/vgl/pull/1657))
- blockchain: Explicit script ver in coinbase checks ([Vigil/vgld#1658](https://github.com/vigilnetwork/vgl/pull/1658))
- chaincfg: Explicit unique net addr prefix ([Vigil/vgld#1663](https://github.com/vigilnetwork/vgl/pull/1663))
- chaincfg: Introduce params lookup by addr prefix ([Vigil/vgld#1664](https://github.com/vigilnetwork/vgl/pull/1664))
- VGLutil: Lookup params by addr prefix in chaincfg ([Vigil/vgld#1665](https://github.com/vigilnetwork/vgl/pull/1665))
- peer: Deprecate dependency on chaincfg ([Vigil/vgld#1671](https://github.com/vigilnetwork/vgl/pull/1671))
- server: Update for deprecated peer chaincfg ([Vigil/vgld#1671](https://github.com/vigilnetwork/vgl/pull/1671))
- fees: drop unused chaincfg ([Vigil/vgld#1675](https://github.com/vigilnetwork/vgl/pull/1675))
- lru: Implement a new module with generic LRU cache ([Vigil/vgld#1683](https://github.com/vigilnetwork/vgl/pull/1683))
- peer: Use lru cache module for inventory ([Vigil/vgld#1683](https://github.com/vigilnetwork/vgl/pull/1683))
- peer: Use lru cache module for nonces ([Vigil/vgld#1683](https://github.com/vigilnetwork/vgl/pull/1683))
- server: Use lru cache module for addresses ([Vigil/vgld#1683](https://github.com/vigilnetwork/vgl/pull/1683))
- multi: drop init and just set default log ([Vigil/vgld#1676](https://github.com/vigilnetwork/vgl/pull/1676))
- multi: deprecate DisableLog ([Vigil/vgld#1676](https://github.com/vigilnetwork/vgl/pull/1676))
- blockchain: Remove unused params from block index ([Vigil/vgld#1674](https://github.com/vigilnetwork/vgl/pull/1674))
- bech32: Initial Version ([Vigil/vgld#1646](https://github.com/vigilnetwork/vgl/pull/1646))
- chaincfg: Add extended key accessor funcs ([Vigil/vgld#1694](https://github.com/vigilnetwork/vgl/pull/1694))
- chaincfg: Rename extended key accessor funcs ([Vigil/vgld#1699](https://github.com/vigilnetwork/vgl/pull/1699))
- wire: Accurate calculations of maximum length ([Vigil/vgld#1672](https://github.com/vigilnetwork/vgl/pull/1672))
- wire: Fix MsgCFTypes maximum payload length ([Vigil/vgld#1673](https://github.com/vigilnetwork/vgl/pull/1673))
- txscript: Deprecate HasP2SHScriptSigStakeOpCodes ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Deprecate IsStakeOutput ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Deprecate GetMultisigMandN ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Introduce zero-alloc script tokenizer ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize script disasm ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Introduce raw script sighash calc func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize CalcSignatureHash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make isSmallInt accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make asSmallInt accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make isStakeOpcode accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize IsPayToScriptHash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize IsMultisigScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize IsMultisigSigScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize GetSigOpCount ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize isAnyKindOfScriptHash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize IsPushOnlyScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize new engine push only script ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Check p2sh push before parsing scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize GetPreciseSigOpCount ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make typeOfScript accept raw script ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript pay-to-script-hash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isScriptHash function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript multisig ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isMultiSig function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript pay-to-pubkey ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isPubkey function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript pay-to-alt-pubkey ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript pay-to-pubkey-hash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isPubkeyHash function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript pay-to-alt-pk-hash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript nulldata detection ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isNullData function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript stakesub detection ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isStakeSubmission function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript stakegen detection ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isStakeGen function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript stakerev detection ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isStakeRevocation function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize typeOfScript stakechange detect ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isSStxChange function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ContainsStakeOpCodes ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractCoinbaseNullData ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Convert CalcScriptInfo ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isPushOnly function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused getSigOpCount function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize CalcMultiSigStats ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize multi sig redeem script func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Convert GetScriptHashFromP2SHScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize PushedData ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize IsUnspendable ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make canonicalPush accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractAtomicSwapDataPushes ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs scripthash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs pubkeyhash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs altpubkeyhash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs pubkey ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs altpubkey ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs multisig ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakesub ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakegen ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakerev ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs stakechange ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAddrs nulldata ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Optimize ExtractPkScriptAltSigType ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused extractOneBytePush func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isPubkeyAlt function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isPubkeyHashAlt function ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused isOneByteMaxDataPush func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: mergeMultiSig function def order cleanup ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Use raw scripts in RawTxInSignature ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Use raw scripts in RawTxInSignatureAlt ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Correct p2pkSignatureScriptAlt comment ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Use raw scripts in SignTxOutput ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Implement efficient opcode data removal ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make isDisabled accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make alwaysIllegal accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make isConditional accept raw opcode ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make min push accept raw opcode and data ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Convert to use non-parsed opcode disasm ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Refactor engine to use raw scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused removeOpcodeByData func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Rename removeOpcodeByDataRaw func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused calcSignatureHash func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Rename calcSignatureHashRaw func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused parseScript func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused unparseScript func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused parsedOpcode.bytes func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Remove unused parseScriptTemplate func ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make executeOpcode take opcode and data ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Make op callbacks take opcode and data ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- VGLutil: Fix NewTxDeepTxIns implementation ([Vigil/vgld#1685](https://github.com/vigilnetwork/vgl/pull/1685))
- stake: drop txscript.DefaultScriptVersion usage ([Vigil/vgld#1704](https://github.com/vigilnetwork/vgl/pull/1704))
- peer: invSendQueue is a FIFO ([Vigil/vgld#1680](https://github.com/vigilnetwork/vgl/pull/1680))
- peer: pendingMsgs is a FIFO ([Vigil/vgld#1680](https://github.com/vigilnetwork/vgl/pull/1680))
- blockchain: drop container/list ([Vigil/vgld#1682](https://github.com/vigilnetwork/vgl/pull/1682))
- blockmanager: use local var for the request queue ([Vigil/vgld#1622](https://github.com/vigilnetwork/vgl/pull/1622))
- server: return on outbound peer creation error ([Vigil/vgld#1637](https://github.com/vigilnetwork/vgl/pull/1637))
- hdkeychain: Remove Address method ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Remove SetNet method ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Require network on decode extended key ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Don't rely on global state ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Introduce NetworkParams interface ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- server: Remove unused ScheduleShutdown func ([Vigil/vgld#1711](https://github.com/vigilnetwork/vgl/pull/1711))
- server: Remove unused dynamicTickDuration func ([Vigil/vgld#1711](https://github.com/vigilnetwork/vgl/pull/1711))
- main: Convert signal handling to use context ([Vigil/vgld#1712](https://github.com/vigilnetwork/vgl/pull/1712))
- txscript: Remove checks for impossible conditions ([Vigil/vgld#1713](https://github.com/vigilnetwork/vgl/pull/1713))
- indexers: Remove unused func ([Vigil/vgld#1714](https://github.com/vigilnetwork/vgl/pull/1714))
- multi: fix onVoteReceivedHandler shutdown ([Vigil/vgld#1721](https://github.com/vigilnetwork/vgl/pull/1721))
- wire: Rename extended errors to malformed errors ([Vigil/vgld#1742](https://github.com/vigilnetwork/vgl/pull/1742))
- rpcwebsocket: convert from list to simple FIFO ([Vigil/vgld#1726](https://github.com/vigilnetwork/vgl/pull/1726))
- VGLec: implement GenerateKey ([Vigil/vgld#1652](https://github.com/vigilnetwork/vgl/pull/1652))
- txscript: Remove SigHashOptimization constant ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- txscript: Remove CheckForDuplicateHashes constant ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- txscript: Remove CPUMinerThreads constant ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Move DNSSeed stringer next to type def ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Remove all registration capabilities ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Move mainnet code to mainnet files ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Move testnet3 code to testnet files ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Move simnet code to testnet files ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Move regnet code to regnet files ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Concrete genesis hash in Params struct ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Use scripts in block one token payouts ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Convert global param defs to funcs ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- edwards: remove curve param ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- edwards: unexport EncodedBytesToBigIntPoint ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- edwards: unexport a slew of funcs ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- edwards: add signature IsEqual and Verify methods ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- edwards: add Sign method to PrivateKey ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- chaincfg: Add addr params accessor funcs ([Vigil/vgld#1766](https://github.com/vigilnetwork/vgl/pull/1766))
- schnorr: remove curve param ([Vigil/vgld#1764](https://github.com/vigilnetwork/vgl/pull/1764))
- schnorr: unexport functions ([Vigil/vgld#1764](https://github.com/vigilnetwork/vgl/pull/1764))
- schnorr: add signature IsEqual and Verify methods ([Vigil/vgld#1764](https://github.com/vigilnetwork/vgl/pull/1764))
- secp256k1: unexport NAF ([Vigil/vgld#1764](https://github.com/vigilnetwork/vgl/pull/1764))
- addrmgr: drop container/list ([Vigil/vgld#1679](https://github.com/vigilnetwork/vgl/pull/1679))
- VGLutil: Remove unused ErrAddressCollision ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- dcurtil: Remove unused ErrMissingDefaultNet ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Require network on address decode ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Remove IsForNet from Address interface ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Remove DSA from Address interface ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Remove Net from Address interface ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Rename EncodeAddress to Address ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Don't store net ref in addr impls ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Require network on WIF decode ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Accept magic bytes directly in NewWIF ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Introduce AddressParams interface ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- blockchain: Do coinbase nulldata check locally ([Vigil/vgld#1770](https://github.com/vigilnetwork/vgl/pull/1770))
- blockchain: update CalcBlockSubsidy ([Vigil/vgld#1750](https://github.com/vigilnetwork/vgl/pull/1750))
- txscript: Use const for sighashall optimization ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove DisableLog ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Unexport HasP2SHScriptSigStakeOpCodes ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove third GetPreciseSigOpCount param ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove IsMultisigScript err return ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Unexport IsStakeOutput ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove CalcScriptInfo ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove multisig redeem script err return ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove GetScriptHashFromP2SHScript ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove GetMultisigMandN ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Remove DefaultScriptVersion ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Use secp256k1 types in sig cache ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- multi: decouple BlockManager from server ([Vigil/vgld#1728](https://github.com/vigilnetwork/vgl/pull/1728))
- database: Introduce BlockSerializer interface ([Vigil/vgld#1799](https://github.com/vigilnetwork/vgl/pull/1799))
- hdkeychain: Add ChildNum and Depth methods ([Vigil/vgld#1800](https://github.com/vigilnetwork/vgl/pull/1800))
- chaincfg: Avoid block 1 subsidy codegen explosion ([Vigil/vgld#1801](https://github.com/vigilnetwork/vgl/pull/1801))
- chaincfg: Add stake params accessor funcs ([Vigil/vgld#1802](https://github.com/vigilnetwork/vgl/pull/1802))
- stake: Remove DisableLog ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Remove unused TxSSGenStakeOutputInfo ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Remove unused TxSSRtxStakeOutputInfo ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Remove unused SetTxTree ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Introduce StakeParams interface ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Accept AddressParams for ticket commit addr ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- gcs: Optimize AddSigScript ([Vigil/vgld#1804](https://github.com/vigilnetwork/vgl/pull/1804))
- chaincfg: Add subsidy params accessor funcs ([Vigil/vgld#1813](https://github.com/vigilnetwork/vgl/pull/1813))
- blockchain/standalone: Implement a new module ([Vigil/vgld#1808](https://github.com/vigilnetwork/vgl/pull/1808))
- blockchain/standalone: Add merkle root calc funcs ([Vigil/vgld#1809](https://github.com/vigilnetwork/vgl/pull/1809))
- blockchain/standalone: Add subsidy calc funcs ([Vigil/vgld#1812](https://github.com/vigilnetwork/vgl/pull/1812))
- blockchain/standalone: Add IsCoinBaseTx ([Vigil/vgld#1815](https://github.com/vigilnetwork/vgl/pull/1815))
- stake: Check minimum req outputs for votes earlier ([Vigil/vgld#1819](https://github.com/vigilnetwork/vgl/pull/1819))
- blockchain: Use standalone module for merkle calcs ([Vigil/vgld#1816](https://github.com/vigilnetwork/vgl/pull/1816))
- blockchain: Use standalone for coinbase checks ([Vigil/vgld#1816](https://github.com/vigilnetwork/vgl/pull/1816))
- blockchain: Use standalone module subsidy calcs ([Vigil/vgld#1816](https://github.com/vigilnetwork/vgl/pull/1816))
- blockchain: Use standalone module for work funcs ([Vigil/vgld#1816](https://github.com/vigilnetwork/vgl/pull/1816))
- blockchain: Remove deprecated code ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- blockchain: Accept subsidy cache in config ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- mining: Use lastest major version deps ([Vigil/vgld#1831](https://github.com/vigilnetwork/vgl/pull/1831))
- connmgr: Accept DNS seeds as string slice ([Vigil/vgld#1833](https://github.com/vigilnetwork/vgl/pull/1833))
- peer: Remove deprecated Config.ChainParams field ([Vigil/vgld#1834](https://github.com/vigilnetwork/vgl/pull/1834))
- peer: Accept hash slice for block locators ([Vigil/vgld#1834](https://github.com/vigilnetwork/vgl/pull/1834))
- peer: Use latest major version deps ([Vigil/vgld#1834](https://github.com/vigilnetwork/vgl/pull/1834))
- mempool: Use latest major version deps ([Vigil/vgld#1835](https://github.com/vigilnetwork/vgl/pull/1835))
- main: Update to use all new major module versions ([Vigil/vgld#1837](https://github.com/vigilnetwork/vgl/pull/1837))
- blockchain: Implement stricter bounds checking ([Vigil/vgld#1825](https://github.com/vigilnetwork/vgl/pull/1825))
- gcs: Start v2 module dev cycle ([Vigil/vgld#1843](https://github.com/vigilnetwork/vgl/pull/1843))
- gcs: Support empty filters ([Vigil/vgld#1844](https://github.com/vigilnetwork/vgl/pull/1844))
- gcs: Make error consistent with rest of codebase ([Vigil/vgld#1846](https://github.com/vigilnetwork/vgl/pull/1846))
- gcs: Add filter version support ([Vigil/vgld#1848](https://github.com/vigilnetwork/vgl/pull/1848))
- gcs: Correct zero hash filter matches ([Vigil/vgld#1857](https://github.com/vigilnetwork/vgl/pull/1857))
- gcs: Standardize serialization on a single format ([Vigil/vgld#1851](https://github.com/vigilnetwork/vgl/pull/1851))
- gcs: Optimize Hash ([Vigil/vgld#1853](https://github.com/vigilnetwork/vgl/pull/1853))
- gcs: Group V1 filter funcs after filter defs ([Vigil/vgld#1854](https://github.com/vigilnetwork/vgl/pull/1854))
- gcs: Support independent fp rate and bin size ([Vigil/vgld#1854](https://github.com/vigilnetwork/vgl/pull/1854))
- blockchain: Refactor best chain state init ([Vigil/vgld#1871](https://github.com/vigilnetwork/vgl/pull/1871))
- gcs: Implement version 2 filters ([Vigil/vgld#1856](https://github.com/vigilnetwork/vgl/pull/1856))
- blockchain: Cleanup subsidy cache init order ([Vigil/vgld#1873](https://github.com/vigilnetwork/vgl/pull/1873))
- multi: use chain ref. from blockmanager config ([Vigil/vgld#1879](https://github.com/vigilnetwork/vgl/pull/1879))
- multi: remove unused funcs and vars ([Vigil/vgld#1880](https://github.com/vigilnetwork/vgl/pull/1880))
- gcs: Prevent empty data elements in v2 filters ([Vigil/vgld#1911](https://github.com/vigilnetwork/vgl/pull/1911))
- crypto: import ripemd160 ([Vigil/vgld#1907](https://github.com/vigilnetwork/vgl/pull/1907))
- multi: Use secp256k1/v2 module ([Vigil/vgld#1919](https://github.com/vigilnetwork/vgl/pull/1919))
- multi: Use crypto/ripemd160 module ([Vigil/vgld#1918](https://github.com/vigilnetwork/vgl/pull/1918))
- multi: Use VGLec/edwards/v2 module ([Vigil/vgld#1920](https://github.com/vigilnetwork/vgl/pull/1920))
- gcs: Prevent empty data elements fp matches ([Vigil/vgld#1940](https://github.com/vigilnetwork/vgl/pull/1940))
- main: Update to use all new module versions ([Vigil/vgld#1946](https://github.com/vigilnetwork/vgl/pull/1946))
- blockchain/standalone: Add inclusion proof funcs ([Vigil/vgld#1906](https://github.com/vigilnetwork/vgl/pull/1906))

### Developer-related module management:

- build: Require VGLjson v1.2.0 ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- multi: Remove non-root module replacements ([Vigil/vgld#1599](https://github.com/vigilnetwork/vgl/pull/1599))
- VGLjson: Introduce v2 module without wallet types ([Vigil/vgld#1607](https://github.com/vigilnetwork/vgl/pull/1607))
- release: Freeze version 1 mempool module ([Vigil/vgld#1613](https://github.com/vigilnetwork/vgl/pull/1613))
- release: Introduce mempool v2 module ([Vigil/vgld#1613](https://github.com/vigilnetwork/vgl/pull/1613))
- main: Tidy module to latest ([Vigil/vgld#1613](https://github.com/vigilnetwork/vgl/pull/1613))
- main: Update for mempool/v2 ([Vigil/vgld#1616](https://github.com/vigilnetwork/vgl/pull/1616))
- multi: Add go 1.11 directive to all modules ([Vigil/vgld#1677](https://github.com/vigilnetwork/vgl/pull/1677))
- build: Tidy module sums (go mod tidy) ([Vigil/vgld#1692](https://github.com/vigilnetwork/vgl/pull/1692))
- release: Freeze version 1 hdkeychain module ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- release: Introduce hdkeychain v2 module ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- release: Freeze version 1 chaincfg module ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Introduce chaincfg v2 module ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- chaincfg: Use VGLec/edwards/v1.0.0 ([Vigil/vgld#1758](https://github.com/vigilnetwork/vgl/pull/1758))
- VGLutil: Prepare v1.3.0 ([Vigil/vgld#1761](https://github.com/vigilnetwork/vgl/pull/1761))
- release: freeze version 1 VGLec/edwards module ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- edwards: Introduce v2 module ([Vigil/vgld#1762](https://github.com/vigilnetwork/vgl/pull/1762))
- release: freeze version 1 VGLec/secp256k1 module ([Vigil/vgld#1764](https://github.com/vigilnetwork/vgl/pull/1764))
- secp256k1: Introduce v2 module ([Vigil/vgld#1764](https://github.com/vigilnetwork/vgl/pull/1764))
- multi: Update all modules for chaincfg v1.5.1 ([Vigil/vgld#1768](https://github.com/vigilnetwork/vgl/pull/1768))
- release: Freeze version 1 VGLutil module ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Update to use chaincfg/v2 module ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- release: Introduce VGLutil v2 module ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- database: Use chaincfg/v2 ([Vigil/vgld#1772](https://github.com/vigilnetwork/vgl/pull/1772))
- txscript: Prepare v1.1.0 ([Vigil/vgld#1773](https://github.com/vigilnetwork/vgl/pull/1773))
- stake: Prepare v1.2.0 ([Vigil/vgld#1775](https://github.com/vigilnetwork/vgl/pull/1775))
- release: Freeze version 1 txscript module ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- txscript: Use VGLutil/v2 ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- release: Introduce txscript v2 module ([Vigil/vgld#1774](https://github.com/vigilnetwork/vgl/pull/1774))
- main: Add requires for new version modules ([Vigil/vgld#1776](https://github.com/vigilnetwork/vgl/pull/1776))
- VGLjson: Introduce v3 and move types to module ([Vigil/vgld#1779](https://github.com/vigilnetwork/vgl/pull/1779))
- jsonrpc/types: Prepare 1.0.0 ([Vigil/vgld#1787](https://github.com/vigilnetwork/vgl/pull/1787))
- main: Use latest JSON-RPC types ([Vigil/vgld#1789](https://github.com/vigilnetwork/vgl/pull/1789))
- multi: Use Vigil fork of go-socks ([Vigil/vgld#1790](https://github.com/vigilnetwork/vgl/pull/1790))
- rpcclient: Prepare v2.1.0 ([Vigil/vgld#1791](https://github.com/vigilnetwork/vgl/pull/1791))
- release: Freeze version 2 rpcclient module ([Vigil/vgld#1793](https://github.com/vigilnetwork/vgl/pull/1793))
- rpcclient: Use VGLjson/v3 ([Vigil/vgld#1793](https://github.com/vigilnetwork/vgl/pull/1793))
- release: Introduce rpcclient v3 module ([Vigil/vgld#1793](https://github.com/vigilnetwork/vgl/pull/1793))
- main: Use rpcclient/v3 ([Vigil/vgld#1795](https://github.com/vigilnetwork/vgl/pull/1795))
- hdkeychain: Prepare v2.0.1 ([Vigil/vgld#1798](https://github.com/vigilnetwork/vgl/pull/1798))
- release: Freeze version 1 database module ([Vigil/vgld#1799](https://github.com/vigilnetwork/vgl/pull/1799))
- database: Use VGLutil/v2 ([Vigil/vgld#1799](https://github.com/vigilnetwork/vgl/pull/1799))
- release: Introduce database v2 module ([Vigil/vgld#1799](https://github.com/vigilnetwork/vgl/pull/1799))
- release: Freeze version 1 blockchain/stake module ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Use VGLutil/v2 and chaincfg/v2 ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- Use txscript/v2 ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- stake: Use database/v2 ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- release: Introduce blockchain/stake v2 module ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- gcs: Use txscript/v2 ([Vigil/vgld#1804](https://github.com/vigilnetwork/vgl/pull/1804))
- gcs: Prepare v1.1.0 ([Vigil/vgld#1804](https://github.com/vigilnetwork/vgl/pull/1804))
- release: Freeze version 3 rpcclient module ([Vigil/vgld#1807](https://github.com/vigilnetwork/vgl/pull/1807))
- rpcclient: Use VGLutil/v2 and chaincfg/v2 ([Vigil/vgld#1807](https://github.com/vigilnetwork/vgl/pull/1807))
- release: Introduce rpcclient v4 module ([Vigil/vgld#1807](https://github.com/vigilnetwork/vgl/pull/1807))
- blockchain/standalone: Prepare v1.0.0 ([Vigil/vgld#1817](https://github.com/vigilnetwork/vgl/pull/1817))
- main: Consume latest module minors and patches ([Vigil/vgld#1822](https://github.com/vigilnetwork/vgl/pull/1822))
- blockchain: Prepare v1.2.0 ([Vigil/vgld#1820](https://github.com/vigilnetwork/vgl/pull/1820))
- mining: Prepare v1.1.1 ([Vigil/vgld#1826](https://github.com/vigilnetwork/vgl/pull/1826))
- release: Freeze version 1 blockchain module use ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- blockchain: Use lastest major version deps ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- release: Introduce blockchain v2 module ([Vigil/vgld#1823](https://github.com/vigilnetwork/vgl/pull/1823))
- connmgr: Prepare v1.1.0 ([Vigil/vgld#1828](https://github.com/vigilnetwork/vgl/pull/1828))
- peer: Prepare v1.2.0 ([Vigil/vgld#1830](https://github.com/vigilnetwork/vgl/pull/1830))
- release: Freeze version 1 mining module use ([Vigil/vgld#1831](https://github.com/vigilnetwork/vgl/pull/1831))
- release: Introduce mining v2 module ([Vigil/vgld#1831](https://github.com/vigilnetwork/vgl/pull/1831))
- mempool: Prepare v2.1.0 ([Vigil/vgld#1832](https://github.com/vigilnetwork/vgl/pull/1832))
- release: Freeze version 1 connmgr module use ([Vigil/vgld#1833](https://github.com/vigilnetwork/vgl/pull/1833))
- release: Introduce connmgr v2 module ([Vigil/vgld#1833](https://github.com/vigilnetwork/vgl/pull/1833))
- release: Freeze version 1 peer module use ([Vigil/vgld#1834](https://github.com/vigilnetwork/vgl/pull/1834))
- release: Introduce peer v2 module ([Vigil/vgld#1834](https://github.com/vigilnetwork/vgl/pull/1834))
- blockchain: Prepare v2.0.1 ([Vigil/vgld#1836](https://github.com/vigilnetwork/vgl/pull/1836))
- release: Freeze version 2 mempool module use ([Vigil/vgld#1835](https://github.com/vigilnetwork/vgl/pull/1835))
- release: Introduce mempool v3 module ([Vigil/vgld#1835](https://github.com/vigilnetwork/vgl/pull/1835))
- go.mod: sync ([Vigil/vgld#1913](https://github.com/vigilnetwork/vgl/pull/1913))
- secp256k1: Prepare v2.0.0 ([Vigil/vgld#1916](https://github.com/vigilnetwork/vgl/pull/1916))
- wire: Prepare v1.3.0 ([Vigil/vgld#1925](https://github.com/vigilnetwork/vgl/pull/1925))
- chaincfg: Prepare v2.3.0 ([Vigil/vgld#1926](https://github.com/vigilnetwork/vgl/pull/1926))
- VGLjson: Prepare v3.0.1 ([Vigil/vgld#1927](https://github.com/vigilnetwork/vgl/pull/1927))
- rpc/jsonrpc/types: Prepare v2.0.0 ([Vigil/vgld#1928](https://github.com/vigilnetwork/vgl/pull/1928))
- VGLutil: Prepare v2.0.1 ([Vigil/vgld#1929](https://github.com/vigilnetwork/vgl/pull/1929))
- blockchain/standalone: Prepare v1.1.0 ([Vigil/vgld#1930](https://github.com/vigilnetwork/vgl/pull/1930))
- txscript: Prepare v2.1.0 ([Vigil/vgld#1931](https://github.com/vigilnetwork/vgl/pull/1931))
- database: Prepare v2.0.1 ([Vigil/vgld#1932](https://github.com/vigilnetwork/vgl/pull/1932))
- blockchain/stake: Prepare v2.0.2 ([Vigil/vgld#1933](https://github.com/vigilnetwork/vgl/pull/1933))
- gcs: Prepare v2.0.0 ([Vigil/vgld#1934](https://github.com/vigilnetwork/vgl/pull/1934))
- blockchain: Prepare v2.1.0 ([Vigil/vgld#1935](https://github.com/vigilnetwork/vgl/pull/1935))
- addrmgr: Prepare v1.1.0 ([Vigil/vgld#1936](https://github.com/vigilnetwork/vgl/pull/1936))
- connmgr: Prepare v2.1.0 ([Vigil/vgld#1937](https://github.com/vigilnetwork/vgl/pull/1937))
- hdkeychain: Prepare v2.1.0 ([Vigil/vgld#1938](https://github.com/vigilnetwork/vgl/pull/1938))
- peer: Prepare v2.1.0 ([Vigil/vgld#1939](https://github.com/vigilnetwork/vgl/pull/1939))
- fees: Prepare v2.0.0 ([Vigil/vgld#1941](https://github.com/vigilnetwork/vgl/pull/1941))
- rpcclient: Prepare v4.1.0 ([Vigil/vgld#1943](https://github.com/vigilnetwork/vgl/pull/1943))
- mining: Prepare v2.0.1 ([Vigil/vgld#1944](https://github.com/vigilnetwork/vgl/pull/1944))
- mempool: Prepare v3.1.0 ([Vigil/vgld#1945](https://github.com/vigilnetwork/vgl/pull/1945))

### Testing and Quality Assurance:

- mempool: Accept test mungers for vote tx ([Vigil/vgld#1595](https://github.com/vigilnetwork/vgl/pull/1595))
- build: Replace TravisCI with CI via Github actions ([Vigil/vgld#1903](https://github.com/vigilnetwork/vgl/pull/1903))
- build: Setup github actions for CI ([Vigil/vgld#1902](https://github.com/vigilnetwork/vgl/pull/1902))
- TravisCI: Recommended install for golangci-lint ([Vigil/vgld#1808](https://github.com/vigilnetwork/vgl/pull/1808))
- TravisCI: Use more portable module ver stripping ([Vigil/vgld#1784](https://github.com/vigilnetwork/vgl/pull/1784))
- TravisCI: Test and lint latest version modules ([Vigil/vgld#1776](https://github.com/vigilnetwork/vgl/pull/1776))
- TravisCI: Disable race detector ([Vigil/vgld#1749](https://github.com/vigilnetwork/vgl/pull/1749))
- TravisCI: Set ./run_vgl_tests.sh executable perms ([Vigil/vgld#1648](https://github.com/vigilnetwork/vgl/pull/1648))
- travis: bump golangci-lint to v1.18.0 ([Vigil/vgld#1890](https://github.com/vigilnetwork/vgl/pull/1890))
- travis: Test go1.13 and drop go1.11 ([Vigil/vgld#1875](https://github.com/vigilnetwork/vgl/pull/1875))
- travis: Allow staged builds with build cache ([Vigil/vgld#1797](https://github.com/vigilnetwork/vgl/pull/1797))
- travis: drop docker and test directly ([Vigil/vgld#1783](https://github.com/vigilnetwork/vgl/pull/1783))
- travis: test go1.12 ([Vigil/vgld#1627](https://github.com/vigilnetwork/vgl/pull/1627))
- travis: Add misspell linter ([Vigil/vgld#1618](https://github.com/vigilnetwork/vgl/pull/1618))
- travis: run linters in each module ([Vigil/vgld#1601](https://github.com/vigilnetwork/vgl/pull/1601))
- multi: switch to golangci-lint ([Vigil/vgld#1575](https://github.com/vigilnetwork/vgl/pull/1575))
- blockchain: Consistent legacy seq lock tests ([Vigil/vgld#1580](https://github.com/vigilnetwork/vgl/pull/1580))
- blockchain: Add test logic to find deployments ([Vigil/vgld#1581](https://github.com/vigilnetwork/vgl/pull/1581))
- blockchain: Introduce chaingen test harness ([Vigil/vgld#1583](https://github.com/vigilnetwork/vgl/pull/1583))
- blockchain: Use harness in force head reorg tests ([Vigil/vgld#1584](https://github.com/vigilnetwork/vgl/pull/1584))
- blockchain: Use harness in stake version tests ([Vigil/vgld#1585](https://github.com/vigilnetwork/vgl/pull/1585))
- blockchain: Use harness in checkblktemplate tests ([Vigil/vgld#1586](https://github.com/vigilnetwork/vgl/pull/1586))
- blockchain: Use harness in threshold state tests ([Vigil/vgld#1587](https://github.com/vigilnetwork/vgl/pull/1587))
- blockchain: Use harness in legacy seqlock tests ([Vigil/vgld#1588](https://github.com/vigilnetwork/vgl/pull/1588))
- blockchain: Use harness in fixed seqlock tests ([Vigil/vgld#1589](https://github.com/vigilnetwork/vgl/pull/1589))
- multi: cleanup linter warnings ([Vigil/vgld#1601](https://github.com/vigilnetwork/vgl/pull/1601))
- txscript: Add remove signature reference test ([Vigil/vgld#1604](https://github.com/vigilnetwork/vgl/pull/1604))
- rpctest: Update for rpccclient/v2 and VGLjson/v2 ([Vigil/vgld#1610](https://github.com/vigilnetwork/vgl/pull/1610))
- wire: Add tests for MsgCFTypes ([Vigil/vgld#1619](https://github.com/vigilnetwork/vgl/pull/1619))
- chaincfg: Move a test to chainhash package ([Vigil/vgld#1632](https://github.com/vigilnetwork/vgl/pull/1632))
- rpctest: Add RemoveNode ([Vigil/vgld#1643](https://github.com/vigilnetwork/vgl/pull/1643))
- rpctest: Add NodesConnected ([Vigil/vgld#1643](https://github.com/vigilnetwork/vgl/pull/1643))
- VGLutil: Reduce global refs in addr unit tests ([Vigil/vgld#1666](https://github.com/vigilnetwork/vgl/pull/1666))
- VGLutil: Consolidate tests into package ([Vigil/vgld#1669](https://github.com/vigilnetwork/vgl/pull/1669))
- peer: Consolidate tests into package ([Vigil/vgld#1670](https://github.com/vigilnetwork/vgl/pull/1670))
- wire: Add tests for BlockHeader (From)Bytes ([Vigil/vgld#1600](https://github.com/vigilnetwork/vgl/pull/1600))
- wire: Add tests for MsgGetCFilter ([Vigil/vgld#1628](https://github.com/vigilnetwork/vgl/pull/1628))
- VGLutil: Add tests for NewTxDeep ([Vigil/vgld#1684](https://github.com/vigilnetwork/vgl/pull/1684))
- rpctest: Introduce VotingWallet ([Vigil/vgld#1668](https://github.com/vigilnetwork/vgl/pull/1668))
- txscript: Add stake tx remove opcode tests ([Vigil/vgld#1210](https://github.com/vigilnetwork/vgl/pull/1210))
- txscript: Move init func in benchmarks to top ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for script parsing ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for DisasmString ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Convert sighash calc tests ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for IsPayToScriptHash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmarks for IsMutlsigScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmarks for IsMutlsigSigScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for GetSigOpCount ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add tests for stake-tagged script hash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for isAnyKindOfScriptHash ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for IsPushOnlyScript ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for GetPreciseSigOpCount ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for GetScriptClass ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for pay-to-pubkey scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for pay-to-alt-pubkey scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for pay-to-pubkey-hash scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for pay-to-alt-pubkey-hash scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for null scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for stake submission scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for stake generation scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for stake revocation scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add bench for stake change scripts ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for ContainsStakeOpCodes ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for ExtractCoinbaseNullData ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add CalcMultiSigStats benchmark ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add multisig redeem script extract bench ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for PushedData ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add benchmark for IsUnspendable ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add tests for atomic swap extraction ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add ExtractAtomicSwapDataPushes benches ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add ExtractPkScriptAddrs benchmarks ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- txscript: Add ExtractPkScriptAltSigType benchmark ([Vigil/vgld#1656](https://github.com/vigilnetwork/vgl/pull/1656))
- wire: Add tests for MsgGetCFTypes ([Vigil/vgld#1703](https://github.com/vigilnetwork/vgl/pull/1703))
- blockchain: Allow named blocks in chaingen harness ([Vigil/vgld#1701](https://github.com/vigilnetwork/vgl/pull/1701))
- txscript: Cleanup opcode removal by data tests ([Vigil/vgld#1702](https://github.com/vigilnetwork/vgl/pull/1702))
- hdkeychain: Correct benchmark extended key ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Consolidate tests into package ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Use locally-scoped netparams in tests ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- hdkeychain: Use mock net params in tests ([Vigil/vgld#1696](https://github.com/vigilnetwork/vgl/pull/1696))
- wire: Add tests for MsgGetCFHeaders ([Vigil/vgld#1720](https://github.com/vigilnetwork/vgl/pull/1720))
- wire: Add tests for MsgCFHeaders ([Vigil/vgld#1732](https://github.com/vigilnetwork/vgl/pull/1732))
- main/rpctest: Update for hdkeychain/v2 ([Vigil/vgld#1707](https://github.com/vigilnetwork/vgl/pull/1707))
- rpctest: Allow custom miner on voting wallet ([Vigil/vgld#1751](https://github.com/vigilnetwork/vgl/pull/1751))
- wire: Add tests for MsgCFilter ([Vigil/vgld#1741](https://github.com/vigilnetwork/vgl/pull/1741))
- chaincfg; Add tests for required unique fields ([Vigil/vgld#1698](https://github.com/vigilnetwork/vgl/pull/1698))
- fullblocktests: Add coinbase nulldata tests ([Vigil/vgld#1769](https://github.com/vigilnetwork/vgl/pull/1769))
- VGLutil: Make docs example testable and correct it ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- VGLutil: Use mock addr params in tests ([Vigil/vgld#1767](https://github.com/vigilnetwork/vgl/pull/1767))
- wire: assert MaxMessagePayload limit in tests ([Vigil/vgld#1755](https://github.com/vigilnetwork/vgl/pull/1755))
- docker: use go 1.12 ([Vigil/vgld#1782](https://github.com/vigilnetwork/vgl/pull/1782))
- docker: update alpine and include notes ([Vigil/vgld#1786](https://github.com/vigilnetwork/vgl/pull/1786))
- hdkeychain: Correct a few comment typos ([Vigil/vgld#1796](https://github.com/vigilnetwork/vgl/pull/1796))
- database: Use unique test db names for v2 module ([Vigil/vgld#1806](https://github.com/vigilnetwork/vgl/pull/1806))
- main: Add database/v2 override for tests ([Vigil/vgld#1806](https://github.com/vigilnetwork/vgl/pull/1806))
- gcs: Add benchmark for AddSigScript ([Vigil/vgld#1804](https://github.com/vigilnetwork/vgl/pull/1804))
- txscript: Fix typo in script test data ([Vigil/vgld#1821](https://github.com/vigilnetwork/vgl/pull/1821))
- database: Separate dbs for concurrent db tests ([Vigil/vgld#1824](https://github.com/vigilnetwork/vgl/pull/1824))
- gcs: Overhaul tests and benchmarks ([Vigil/vgld#1845](https://github.com/vigilnetwork/vgl/pull/1845))
- rpctest: Remove leftover debug print ([Vigil/vgld#1862](https://github.com/vigilnetwork/vgl/pull/1862))
- txscript: Fix duplicate test name ([Vigil/vgld#1863](https://github.com/vigilnetwork/vgl/pull/1863))
- gcs: Add benchmark for filter hashing ([Vigil/vgld#1853](https://github.com/vigilnetwork/vgl/pull/1853))
- gcs: Add tests for bit reader/writer ([Vigil/vgld#1855](https://github.com/vigilnetwork/vgl/pull/1855))
- peer: Ensure listener tests sync with messages ([Vigil/vgld#1874](https://github.com/vigilnetwork/vgl/pull/1874))
- rpctest: remove always-nil error ([Vigil/vgld#1913](https://github.com/vigilnetwork/vgl/pull/1913))
- rpctest: use errgroup to catch errors from go routines ([Vigil/vgld#1913](https://github.com/vigilnetwork/vgl/pull/1913))

### Misc:

- release: Bump for 1.5 release cycle ([Vigil/vgld#1546](https://github.com/vigilnetwork/vgl/pull/1546))
- mempool: Fix typo in fetchInputUtxos comment ([Vigil/vgld#1562](https://github.com/vigilnetwork/vgl/pull/1562))
- blockchain: Fix typos found by misspell ([Vigil/vgld#1617](https://github.com/vigilnetwork/vgl/pull/1617))
- VGLutil: Fix typos found by misspell ([Vigil/vgld#1617](https://github.com/vigilnetwork/vgl/pull/1617))
- main: Write memprofile on shutdown ([Vigil/vgld#1655](https://github.com/vigilnetwork/vgl/pull/1655))
- config: Parse network interfaces ([Vigil/vgld#1514](https://github.com/vigilnetwork/vgl/pull/1514))
- config: Cleanup and simplify network info parsing ([Vigil/vgld#1706](https://github.com/vigilnetwork/vgl/pull/1706))
- main: Rework windows service sod notification ([Vigil/vgld#1710](https://github.com/vigilnetwork/vgl/pull/1710))
- multi: fix recent govet findings ([Vigil/vgld#1727](https://github.com/vigilnetwork/vgl/pull/1727))
- rpcserver: Fix misspelling ([Vigil/vgld#1763](https://github.com/vigilnetwork/vgl/pull/1763))
- chaincfg: Run gofmt -s ([Vigil/vgld#1776](https://github.com/vigilnetwork/vgl/pull/1776))
- jsonrpc/types: Update copyright years ([Vigil/vgld#1794](https://github.com/vigilnetwork/vgl/pull/1794))
- stake: Correct comment typo on Hash256PRNG ([Vigil/vgld#1803](https://github.com/vigilnetwork/vgl/pull/1803))
- multi: Correct typos ([Vigil/vgld#1839](https://github.com/vigilnetwork/vgl/pull/1839))
- wire: Fix a few messageError string typos ([Vigil/vgld#1840](https://github.com/vigilnetwork/vgl/pull/1840))
- miningerror: Remove duplicate copyright ([Vigil/vgld#1860](https://github.com/vigilnetwork/vgl/pull/1860))
- multi: Correct typos ([Vigil/vgld#1864](https://github.com/vigilnetwork/vgl/pull/1864))

### Code Contributors (alphabetical order):

- Aaron Campbell
- Conner Fromknecht
- Dave Collins
- David Hill
- Donald Adu-Poku
- Hamid
- J Fixby
- Jamie Holdstock
- JoeGruffins
- Jonathan Chappelow
- Josh Rickmar
- Matheus Degiovani
- Nicola Larosa
- Olaoluwa Osuntokun
- Roei Erez
- Sarlor
- Victor Oliveira




