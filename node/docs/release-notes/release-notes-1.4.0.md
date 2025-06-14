# vgld v1.4.0

This release of vgld introduces a new consensus vote agenda which allows the
stakeholders to decide whether or not to activate changes needed to modify the
sequence lock handling which is required for providing full support for the
Lightning Network.  For those unfamiliar with the voting process in Vigil, this
means that all code in order to make the necessary changes is already included
in this release, however its enforcement will remain dormant until the
stakeholders vote to activate it.

It also contains smart fee estimation, performance enhancements for block relay
and processing, a major internal restructuring of how unspent transaction
outputs are handled, support for whitelisting inbound peers to ensure service
for your own SPV (Simplified Payment Verification) wallets, various updates to
the RPC server such as a new method to query the state of the chain and more
easily supporting external RPC connections over TLS, infrastructure
improvements, and other quality assurance changes.

The following Vigil Change Proposals (VGLP) describes the proposed changes in detail:
- [VGLP0004](https://github.com/Vigil/VGLPs/blob/master/VGLP-0004/VGLP-0004.mediawiki)

**It is important for everyone to upgrade their software to this latest release
even if you don't intend to vote in favor of the agenda.**

## Downgrade Warning

The database format in v1.4.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a lengthy one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Fix Lightning Network Sequence Locks Vote

In order to fully support the Lightning Network, the current sequence lock
consensus rules need to be modified.  A new vote with the id `fixlnseqlocks` is
now available as of this release.  After upgrading, stakeholders may set their
preferences through their wallet or Voting Service Provider's (VSP) website.

### Smart Fee Estimation (`estimatesmartfee`)

A new RPC named `estimatesmartfee` is now available which returns a suitable
fee rate for transactions to use in order to have a high probability of them
being mined within a specified number of confirmations.  The estimation is based
on actual network usage and thus varies according to supply and demand.

This is important in the context of the Lightning Network (LN) and, more
generally, it provides services and users with a mechanism to choose how to
handle network congestion.  For example, payments that are high priority might
be willing to pay a higher fee to help ensure the transaction is mined more
quickly, while lower priority payments might be willing to wait longer in
exchange for paying a lower fee.  This estimation capability provides a way to
obtain a fee that will achieve the desired result with a high probability.

### Support for Whitelisting Inbound Peers

When peers are whitelisted via the `--whitelist` option, they will now be
allowed to connect even when they would otherwise exceed the maximum number of
peers.  This is highly useful in cases where users have configured their wallet
to use SPV mode and only connect to vgld instances that they control for
increased privacy and guaranteed service.

### Several Speed Optimizations

Similar to previous releases, this release also contains several enhancements to
improve speed for the initial sync process, validation, and network operations.

In order to achieve these speedups, there is a lengthy one time database
migration, as previously mentioned, that typically takes anywhere from 30
minutes to an hour to complete depending on hardware.

#### Faster Tip Block Relay

Blocks that extend the current best chain are now relayed to the network
immediately after they pass the initial sanity and contextual checks, most
notably valid proof of work.  This allows blocks to propagate more quickly
throughout the network, which in turn improves vote times.

#### UTXO Set Restructuring

The way the unspent transaction outputs are handled internally has been
overhauled to significantly decrease the time it takes to validate blocks and
transactions.  While this has many benefits, probably the most important one
for most stakeholders is that votes can be cast more quickly which helps reduce
the number of missed votes.

### RPC Server Changes

#### New Chain State Query RPC (`getblockchaininfo`)

A new RPC named `getblockchaininfo` is now available which can be used to query
the state of the chain including details such as its overall verification
progress during initial sync, the maximum supported block size, and that status
of consensus changes (deployments) which require stakeholder votes.  See the
[JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblockchaininfo)
for API details.

#### Removal of Vote Creation RPC (`createrawssgen`)

The deprecated `createrawssgen`, which was previously used to allow creating a
vote via RPC is no longer available.  Votes are time sensitive and thus it does
not make sense to create them offline.

#### Updates to Block and Transaction RPCs

The `getblock`, `getblockheader`, `getrawtransaction`, and
`searchrawtransactions` RPCs now contain additional information such as the
`extradata` field in the header, the `expiry` field in transactions, and the
`blockheight` and `blockindex` of  the block that contains a transaction if it
has been mined.  See the [JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.md)
for API details.

#### Built-in Support for Enabling External TLS RPC Connections

A new command line parameter (`--altdnsnames`) and environment variable
(`vgld_ALT_DNSNAMES`) can now be used before the first launch of drcd to specify
additional external IP addresses and DNS names to add during the certificate
creation that are permitted to connect to the RPC server via TLS.  Previously,
a separate tool was required to accomplish this configuration.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/release-v1.3.0...release-v1.4.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.4.0 release ([Vigil/vgld#1547](https://github.com/vigilnetwork/vgl/pull/1547))
- chaincfg: Introduce agenda for fixlnseqlocks vote ([Vigil/vgld#1578](https://github.com/vigilnetwork/vgl/pull/1578))
- multi: Enable vote for VGLP0004 ([Vigil/vgld#1579](https://github.com/vigilnetwork/vgl/pull/1579))
- peer: Add support for specifying ua comments ([Vigil/vgld#1413](https://github.com/vigilnetwork/vgl/pull/1413))
- blockmanager: Fast relay checked tip blocks ([Vigil/vgld#1443](https://github.com/vigilnetwork/vgl/pull/1443))
- multi: Latest consensus active from simnet genesis ([Vigil/vgld#1482](https://github.com/vigilnetwork/vgl/pull/1482))
- server: Always allow whitelisted inbound peers ([Vigil/vgld#1516](https://github.com/vigilnetwork/vgl/pull/1516))

### Transaction relay (memory pool):

- blockmanager: handle txs in invalid blocks ([Vigil/vgld#1430](https://github.com/vigilnetwork/vgl/pull/1430))
- mempool: Remove potential negative locktime check ([Vigil/vgld#1455](https://github.com/vigilnetwork/vgl/pull/1455))
- mempool: Stake-related readability improvements ([Vigil/vgld#1456](https://github.com/vigilnetwork/vgl/pull/1456))

### RPC:

- multi: Include additional fields on RPC tx results ([Vigil/vgld#1441](https://github.com/vigilnetwork/vgl/pull/1441))
- rpcserver: Allow scripthash addrs in createrawsstx ([Vigil/vgld#1444](https://github.com/vigilnetwork/vgl/pull/1444))
- rpcserver: Remove createrawssgen RPC ([Vigil/vgld#1448](https://github.com/vigilnetwork/vgl/pull/1448))
- rpcclient: support getchaintips RPC ([Vigil/vgld#1469](https://github.com/vigilnetwork/vgl/pull/1469))
- multi: Add getblockchaininfo rpc ([Vigil/vgld#1479](https://github.com/vigilnetwork/vgl/pull/1479))
- rpcserver: Adds ability to allow alternative dns names for TLS ([Vigil/vgld#1476](https://github.com/vigilnetwork/vgl/pull/1476))
- multi: Cleanup recent alt DNS names additions ([Vigil/vgld#1493](https://github.com/vigilnetwork/vgl/pull/1493))
- multi: Cleanup getblock and getblockheader RPCs ([Vigil/vgld#1497](https://github.com/vigilnetwork/vgl/pull/1497))
- multi: Return total chain work in RPC results ([Vigil/vgld#1498](https://github.com/vigilnetwork/vgl/pull/1498))
- rpcserver: Improve GenerateNBlocks error message ([Vigil/vgld#1507](https://github.com/vigilnetwork/vgl/pull/1507))
- rpcserver: Fix verify progress calculation ([Vigil/vgld#1508](https://github.com/vigilnetwork/vgl/pull/1508))
- rpcserver: Fix sendrawtransaction error code ([Vigil/vgld#1512](https://github.com/vigilnetwork/vgl/pull/1512))
- blockchain: Notify stake states after connected block ([Vigil/vgld#1515](https://github.com/vigilnetwork/vgl/pull/1515))
- rpcserver: bump version to 5.0. ([Vigil/vgld#1531](https://github.com/vigilnetwork/vgl/pull/1531))
- rpcclient: support getblockchaininfo RPC ([Vigil/vgld#1539](https://github.com/vigilnetwork/vgl/pull/1539))
- rpcserver: update block template reconstruction ([Vigil/vgld#1567](https://github.com/vigilnetwork/vgl/pull/1567))

### vgld command-line flags and configuration:

- config: add --maxsameip to limit # of conns to same IP ([Vigil/vgld#1517](https://github.com/vigilnetwork/vgl/pull/1517))

### Documentation:

- docs: Update docs for versioned modules ([Vigil/vgld#1391](https://github.com/vigilnetwork/vgl/pull/1391))
- docs: Update for fees package ([Vigil/vgld#1540](https://github.com/vigilnetwork/vgl/pull/1540))
- docs: Revamp main README.md and update docs ([Vigil/vgld#1447](https://github.com/vigilnetwork/vgl/pull/1447))
- docs: Use relative versions in contrib checklist ([Vigil/vgld#1451](https://github.com/vigilnetwork/vgl/pull/1451))
- docs: Use the correct binary name ([Vigil/vgld#1461](https://github.com/vigilnetwork/vgl/pull/1461))
- docs: Add github pull request template ([Vigil/vgld#1474](https://github.com/vigilnetwork/vgl/pull/1474))
- docs: Use unix line ending in mod hierarchy gv ([Vigil/vgld#1487](https://github.com/vigilnetwork/vgl/pull/1487))
- docs: Add README badge and link for goreportcard ([Vigil/vgld#1492](https://github.com/vigilnetwork/vgl/pull/1492))
- sampleconfig: Fix proxy typo ([Vigil/vgld#1513](https://github.com/vigilnetwork/vgl/pull/1513))

### Developer-related package and module changes:

- release: Bump module versions and deps ([Vigil/vgld#1541](https://github.com/vigilnetwork/vgl/pull/1541))
- build: Tidy module sums with go mod tidy ([Vigil/vgld#1408](https://github.com/vigilnetwork/vgl/pull/1408))
- blockchain: update BestState ([Vigil/vgld#1416](https://github.com/vigilnetwork/vgl/pull/1416))
- mempool: tweak trace logs ([Vigil/vgld#1429](https://github.com/vigilnetwork/vgl/pull/1429))
- blockchain: Correct best pool size on disconnect ([Vigil/vgld#1431](https://github.com/vigilnetwork/vgl/pull/1431))
- multi: Make use of new internal version package ([Vigil/vgld#1435](https://github.com/vigilnetwork/vgl/pull/1435))
- peer: Protect handlePongMsg with p.statsMtx ([Vigil/vgld#1438](https://github.com/vigilnetwork/vgl/pull/1438))
- limits: Make limits package internal ([Vigil/vgld#1436](https://github.com/vigilnetwork/vgl/pull/1436))
- indexers: Remove unneeded existsaddrindex iface ([Vigil/vgld#1439](https://github.com/vigilnetwork/vgl/pull/1439))
- blockchain: Reduce block availability assumptions ([Vigil/vgld#1442](https://github.com/vigilnetwork/vgl/pull/1442))
- peer: Provide immediate queue inventory func ([Vigil/vgld#1443](https://github.com/vigilnetwork/vgl/pull/1443))
- server: Add infrastruct for immediate inv relay ([Vigil/vgld#1443](https://github.com/vigilnetwork/vgl/pull/1443))
- blockchain: Add new tip block checked notification ([Vigil/vgld#1443](https://github.com/vigilnetwork/vgl/pull/1443))
- multi: remove chainState dependency in rpcserver ([Vigil/vgld#1417](https://github.com/vigilnetwork/vgl/pull/1417))
- mining: remove chainState dependency ([Vigil/vgld#1418](https://github.com/vigilnetwork/vgl/pull/1418))
- multi: remove chainState deps in server & cpuminer ([Vigil/vgld#1419](https://github.com/vigilnetwork/vgl/pull/1419))
- blockmanager: remove block manager chain state ([Vigil/vgld#1420](https://github.com/vigilnetwork/vgl/pull/1420))
- multi: move MinHighPriority into mining package ([Vigil/vgld#1421](https://github.com/vigilnetwork/vgl/pull/1421))
- multi: add BlkTmplGenerator ([Vigil/vgld#1422](https://github.com/vigilnetwork/vgl/pull/1422))
- multi: add cpuminerConfig ([Vigil/vgld#1423](https://github.com/vigilnetwork/vgl/pull/1423))
- multi: Move update blk time to blk templ generator ([Vigil/vgld#1454](https://github.com/vigilnetwork/vgl/pull/1454))
- multi: No stake height checks in check tx inputs ([Vigil/vgld#1457](https://github.com/vigilnetwork/vgl/pull/1457))
- blockchain: Separate tx input stake checks ([Vigil/vgld#1452](https://github.com/vigilnetwork/vgl/pull/1452))
- blockchain: Ensure no stake opcodes in tx sanity ([Vigil/vgld#1453](https://github.com/vigilnetwork/vgl/pull/1453))
- blockchain: Move finalized tx func to validation ([Vigil/vgld#1465](https://github.com/vigilnetwork/vgl/pull/1465))
- blockchain: Move unique coinbase func to validate ([Vigil/vgld#1466](https://github.com/vigilnetwork/vgl/pull/1466))
- blockchain: Store interrupt channel with state ([Vigil/vgld#1467](https://github.com/vigilnetwork/vgl/pull/1467))
- multi: Cleanup and optimize tx input check code ([Vigil/vgld#1468](https://github.com/vigilnetwork/vgl/pull/1468))
- blockmanager: Avoid duplicate header announcements ([Vigil/vgld#1473](https://github.com/vigilnetwork/vgl/pull/1473))
- VGLjson: additions for pay to contract hash ([Vigil/vgld#1260](https://github.com/vigilnetwork/vgl/pull/1260))
- multi: Break blockchain dependency on VGLjson ([Vigil/vgld#1488](https://github.com/vigilnetwork/vgl/pull/1488))
- chaincfg: Unexport internal errors ([Vigil/vgld#1489](https://github.com/vigilnetwork/vgl/pull/1489))
- multi: Cleanup the unsupported vglwallet commands ([Vigil/vgld#1478](https://github.com/vigilnetwork/vgl/pull/1478))
- multi: Rename ThresholdState to NextThresholdState ([Vigil/vgld#1494](https://github.com/vigilnetwork/vgl/pull/1494))
- VGLjson: Add listtickets command ([Vigil/vgld#1267](https://github.com/vigilnetwork/vgl/pull/1267))
- multi: Add started and done reorg notifications ([Vigil/vgld#1495](https://github.com/vigilnetwork/vgl/pull/1495))
- blockchain: Remove unused CheckWorklessBlockSanity ([Vigil/vgld#1496](https://github.com/vigilnetwork/vgl/pull/1496))
- blockchain: Simplify block template checking ([Vigil/vgld#1499](https://github.com/vigilnetwork/vgl/pull/1499))
- blockchain: Only mark nodes modified when modified ([Vigil/vgld#1503](https://github.com/vigilnetwork/vgl/pull/1503))
- blockchain: Cleanup and optimize stake node logic ([Vigil/vgld#1504](https://github.com/vigilnetwork/vgl/pull/1504))
- blockchain: Separate full data context checks ([Vigil/vgld#1509](https://github.com/vigilnetwork/vgl/pull/1509))
- blockchain: Reverse utxo set semantics ([Vigil/vgld#1471](https://github.com/vigilnetwork/vgl/pull/1471))
- blockchain: Convert to direct single-step reorgs ([Vigil/vgld#1500](https://github.com/vigilnetwork/vgl/pull/1500))
- multi: Migration for utxo set semantics reversal ([Vigil/vgld#1520](https://github.com/vigilnetwork/vgl/pull/1520))
- blockchain: Make version 5 update atomic ([Vigil/vgld#1529](https://github.com/vigilnetwork/vgl/pull/1529))
- blockchain: Simplify force head reorgs ([Vigil/vgld#1526](https://github.com/vigilnetwork/vgl/pull/1526))
- secp256k1: Correct edge case in deterministic sign ([Vigil/vgld#1533](https://github.com/vigilnetwork/vgl/pull/1533))
- VGLjson: Add gettransaction txtype/ticketstatus ([Vigil/vgld#1276](https://github.com/vigilnetwork/vgl/pull/1276))
- txscript: Use ScriptBuilder more ([Vigil/vgld#1519](https://github.com/vigilnetwork/vgl/pull/1519))
- fees: Add estimator package ([Vigil/vgld#1434](https://github.com/vigilnetwork/vgl/pull/1434))
- multi: Integrate fee estimation ([Vigil/vgld#1434](https://github.com/vigilnetwork/vgl/pull/1434))

### Testing and Quality Assurance:

- multi: Use temp directories for database tests ([Vigil/vgld#1404](https://github.com/vigilnetwork/vgl/pull/1404))
- multi: Only use module-scoped data in tests ([Vigil/vgld#1405](https://github.com/vigilnetwork/vgl/pull/1405))
- blockchain: Use temp dirs for fullblocks test ([Vigil/vgld#1406](https://github.com/vigilnetwork/vgl/pull/1406))
- database: Use module-scoped data in iface tests ([Vigil/vgld#1407](https://github.com/vigilnetwork/vgl/pull/1407))
- travis: Update for Go1.11 and module builds ([Vigil/vgld#1415](https://github.com/vigilnetwork/vgl/pull/1415))
- indexers: Use testable bucket for existsaddrindex ([Vigil/vgld#1440](https://github.com/vigilnetwork/vgl/pull/1440))
- txscript: group numeric encoding tests with their opcodes ([Vigil/vgld#1382](https://github.com/vigilnetwork/vgl/pull/1382))
- txscript: add p2sh opcode tests ([Vigil/vgld#1381](https://github.com/vigilnetwork/vgl/pull/1381))
- txscript: add stake opcode tests ([Vigil/vgld#1383](https://github.com/vigilnetwork/vgl/pull/1383))
- main: add address encoding magic constants test ([Vigil/vgld#1458](https://github.com/vigilnetwork/vgl/pull/1458))
- chaingen: Only revoke missed tickets once ([Vigil/vgld#1484](https://github.com/vigilnetwork/vgl/pull/1484))
- chaingen/fullblocktests: Add disapproval tests ([Vigil/vgld#1485](https://github.com/vigilnetwork/vgl/pull/1485))
- multi: Resurrect regression network ([Vigil/vgld#1480](https://github.com/vigilnetwork/vgl/pull/1480))
- multi: Use regression test network in unit tests ([Vigil/vgld#1481](https://github.com/vigilnetwork/vgl/pull/1481))
- main: move cert tests to a separated file ([Vigil/vgld#1502](https://github.com/vigilnetwork/vgl/pull/1502))
- mempool: Accept test mungers for create signed tx ([Vigil/vgld#1576](https://github.com/vigilnetwork/vgl/pull/1576))
- mempool: Implement test harness seq lock calc ([Vigil/vgld#1577](https://github.com/vigilnetwork/vgl/pull/1577))

### Misc:

- release: Bump for 1.4 release cycle ([Vigil/vgld#1414](https://github.com/vigilnetwork/vgl/pull/1414))
- multi: Make changes suggested by Go 1.11 gofmt -s ([Vigil/vgld#1415](https://github.com/vigilnetwork/vgl/pull/1415))
- build: Remove dep toml and lock file ([Vigil/vgld#1460](https://github.com/vigilnetwork/vgl/pull/1460))
- docker: Update to go 1.11 ([Vigil/vgld#1463](https://github.com/vigilnetwork/vgl/pull/1463))
- build: Support MacOS sed for obtaining module list ([Vigil/vgld#1483](https://github.com/vigilnetwork/vgl/pull/1483))
- multi: Correct a few typos found by misspell ([Vigil/vgld#1490](https://github.com/vigilnetwork/vgl/pull/1490))
- multi: Address some golint complaints ([Vigil/vgld#1491](https://github.com/vigilnetwork/vgl/pull/1491))
- multi: Remove unused code ([Vigil/vgld#1505](https://github.com/vigilnetwork/vgl/pull/1505))
- release: Bump siphash version to v1.2.1 ([Vigil/vgld#1538](https://github.com/vigilnetwork/vgl/pull/1538))
- release: Bump module versions and deps ([Vigil/vgld#1541](https://github.com/vigilnetwork/vgl/pull/1541))
- Fix required version of stake module ([Vigil/vgld#1549](https://github.com/vigilnetwork/vgl/pull/1549))
- release: Tidy module files with published versions ([Vigil/vgld#1543](https://github.com/vigilnetwork/vgl/pull/1543))
- mempool: Fix required version of mining module ([Vigil/vgld#1551](https://github.com/vigilnetwork/vgl/pull/1551))

### Code Contributors (alphabetical order):

- Corey Osman
- Dave Collins
- David Hill
- Dmitry Fedorov
- Donald Adu-Poku
- ggoranov
- githubsands
- J Fixby
- Jonathan Chappelow
- Josh Rickmar
- Matheus Degiovani
- Sarlor
- zhizhongzhiwai
