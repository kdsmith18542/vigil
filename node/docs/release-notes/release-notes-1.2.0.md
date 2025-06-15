# vgld v1.2.0

This release of vgld contains significant performance enhancements,
infrastructure improvements, improved access to chain-related information for
providing better SPV (Simplified Payment Verification) support, and other
quality assurance changes.

A significant amount of infrastructure work has also been done this release
cycle towards being able to support several planned scalability optimizations.

## Downgrade Warning

The database format in v1.2.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Significantly Faster Startup

The startup time has been improved by roughly 17x on slower hard disk drives
(HDDs) and 8x on solid state drives (SSDs).

In order to achieve these speedups, there is a one time database migration, as
previously mentioned, that will likely take a while to complete (typically
around 5 to 6 minutes on HDDs and 2 to 3 minutes on SSDs).

### Support For DNS Seed Filtering

In order to better support the forthcoming SPV wallets, support for finding
other peers based upon their enabled services has been added.  This is useful
for both SPV wallets and full nodes since SPV wallets will require access to
full nodes in order to retrieve the necessary proofs and full nodes are
generally not interested in making outgoing connections to SPV wallets.

### Committed Filters

With the intention of supporting light clients, such as SPV wallets, in a
privacy-preserving way while still minimizing the amount of data that needs to
be downloaded, this release adds support for committed filters.  A committed
filter is a combination of a probalistic data structure that is used to test
whether an element is a member of a set with a predetermined collision
probability along with a commitment by consensus-validating full nodes to that
data.

A committed filter is created for every block which allows light clients to
download the filters and match against them locally rather than uploading
personal data to other nodes.

A new service flag is also provided to allow clients to discover nodes that
provide access to filters.

There is a one time database update to build and store the filters for all
existing historical blocks which will likely take a while to complete (typically
around 2 to 3 minutes on HDDs and 1 to 1.5 minutes on SSDs).

### Updated Atomic Swap Contracts

The standard checks for atomic swap contracts have been updated to ensure the
contracts enforce the secret size for safer support between chains with
disparate script rules.

### RPC Server Changes

#### New `getchaintips` RPC

A new RPC named `getchaintips` has been added which allows callers to query
information about the status of known side chains and their branch lengths.
It currently only provides support for side chains that have been seen while the
current instance of the process is running.  This will be further improved in
future releases.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/v1.1.2...v1.2.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.2.0 release ([Vigil/vgld#1139](https://github.com/vigilnetwork/vgl/pull/1139))
- chaincfg: Introduce new type DNSSeed ([Vigil/vgld#961](https://github.com/vigilnetwork/vgl/pull/961))
- blockmanager: sync with the most updated peer ([Vigil/vgld#984](https://github.com/vigilnetwork/vgl/pull/984))
- multi: remove MsgAlert ([Vigil/vgld#1161](https://github.com/vigilnetwork/vgl/pull/1161))
- multi: Add initial committed filter (CF) support ([Vigil/vgld#1151](https://github.com/vigilnetwork/vgl/pull/1151))

### Transaction relay (memory pool):

- txscript: Correct nulldata standardness check ([Vigil/vgld#935](https://github.com/vigilnetwork/vgl/pull/935))
- mempool: Optimize orphan map limiting ([Vigil/vgld#1117](https://github.com/vigilnetwork/vgl/pull/1117))
- mining: Fix duplicate txns in the prio heap ([Vigil/vgld#1108](https://github.com/vigilnetwork/vgl/pull/1108))
- mining: Stop transactions losing their dependants ([Vigil/vgld#1109](https://github.com/vigilnetwork/vgl/pull/1109))

### RPC:

- rpcserver: skip cert create when RPC is disabled ([Vigil/vgld#949](https://github.com/vigilnetwork/vgl/pull/949))
- rpcserver: remove redundant checks in blockTemplateResult ([Vigil/vgld#826](https://github.com/vigilnetwork/vgl/pull/826))
- rpcserver: assert network for validateaddress rpc ([Vigil/vgld#963](https://github.com/vigilnetwork/vgl/pull/963))
- rpcserver: Do not rebroadcast stake transactions ([Vigil/vgld#973](https://github.com/vigilnetwork/vgl/pull/973))
- VGLjson: add ticket fee field to PurchaseTicketCmd ([Vigil/vgld#902](https://github.com/vigilnetwork/vgl/pull/902))
- vglwalletextcmds: remove getseed ([Vigil/vgld#985](https://github.com/vigilnetwork/vgl/pull/985))
- VGLjson: Add SweepAccountCmd & SweepAccountResult ([Vigil/vgld#1027](https://github.com/vigilnetwork/vgl/pull/1027))
- rpcserver: add sweepaccount to the wallet list of commands ([Vigil/vgld#1028](https://github.com/vigilnetwork/vgl/pull/1028))
- rpcserver: add batched request support (json 2.0) ([Vigil/vgld#841](https://github.com/vigilnetwork/vgl/pull/841))
- VGLjson: include summary totals in GetBalanceResult ([Vigil/vgld#1062](https://github.com/vigilnetwork/vgl/pull/1062))
- multi: Implement getchaintips JSON-RPC ([Vigil/vgld#1098](https://github.com/vigilnetwork/vgl/pull/1098))
- rpcserver: Add vgld version info to getversion RPC ([Vigil/vgld#1097](https://github.com/vigilnetwork/vgl/pull/1097))
- rpcserver: Correct getblockheader result text ([Vigil/vgld#1104](https://github.com/vigilnetwork/vgl/pull/1104))
- VGLjson: add StartAutoBuyerCmd & StopAutoBuyerCmd ([Vigil/vgld#903](https://github.com/vigilnetwork/vgl/pull/903))
- VGLjson: fix typo for StartAutoBuyerCmd ([Vigil/vgld#1146](https://github.com/vigilnetwork/vgl/pull/1146))
- VGLjson: require passphrase for StartAutoBuyerCmd ([Vigil/vgld#1147](https://github.com/vigilnetwork/vgl/pull/1147))
- VGLjson: fix StopAutoBuyerCmd registration bug ([Vigil/vgld#1148](https://github.com/vigilnetwork/vgl/pull/1148))
- blockchain: Support testnet stake diff estimation ([Vigil/vgld#1115](https://github.com/vigilnetwork/vgl/pull/1115))
- rpcserver: fix jsonRPCRead data race ([Vigil/vgld#1157](https://github.com/vigilnetwork/vgl/pull/1157))
- VGLjson: Add VerifySeedCmd ([Vigil/vgld#1160](https://github.com/vigilnetwork/vgl/pull/1160))

### vgld command-line flags and configuration:

- mempool: Rename RelayNonStd config option ([Vigil/vgld#1024](https://github.com/vigilnetwork/vgl/pull/1024))
- sampleconfig: Update min relay fee ([Vigil/vgld#959](https://github.com/vigilnetwork/vgl/pull/959))
- sampleconfig: Correct comment ([Vigil/vgld#1063](https://github.com/vigilnetwork/vgl/pull/1063))
- multi: Expand ~ to correct home directory on all OSes ([Vigil/vgld#1041](https://github.com/vigilnetwork/vgl/pull/1041))

### checkdevpremine utility changes:

- checkdevpremine: Remove --skipverify option ([Vigil/vgld#969](https://github.com/vigilnetwork/vgl/pull/969))
- checkdevpremine: Implement --notls option ([Vigil/vgld#969](https://github.com/vigilnetwork/vgl/pull/969))
- checkdevpremine: Make file naming consistent ([Vigil/vgld#969](https://github.com/vigilnetwork/vgl/pull/969))
- checkdevpremine: Fix comment ([Vigil/vgld#969](https://github.com/vigilnetwork/vgl/pull/969))
- checkdevpremine: Remove utility ([Vigil/vgld#1068](https://github.com/vigilnetwork/vgl/pull/1068))

### Documentation:

- fullblocktests: Add missing doc.go file ([Vigil/vgld#956](https://github.com/vigilnetwork/vgl/pull/956))
- docs: Add fullblocktests entry and make consistent ([Vigil/vgld#956](https://github.com/vigilnetwork/vgl/pull/956))
- docs: Add mempool entry to developer tools section ([Vigil/vgld#1058](https://github.com/vigilnetwork/vgl/pull/1058))
- mempool: Add docs.go and flesh out README.md ([Vigil/vgld#1058](https://github.com/vigilnetwork/vgl/pull/1058))
- docs: document packages and fix typo  ([Vigil/vgld#965](https://github.com/vigilnetwork/vgl/pull/965))
- docs: rpcclient is now part of the main vgld repo ([Vigil/vgld#970](https://github.com/vigilnetwork/vgl/pull/970))
- VGLjson: Update README.md ([Vigil/vgld#982](https://github.com/vigilnetwork/vgl/pull/982))
- docs: Remove carriage return ([Vigil/vgld#1106](https://github.com/vigilnetwork/vgl/pull/1106))
- Adjust README.md for new Go version ([Vigil/vgld#1105](https://github.com/vigilnetwork/vgl/pull/1105))
- docs: document how to use go test -coverprofile ([Vigil/vgld#1107](https://github.com/vigilnetwork/vgl/pull/1107))
- addrmgr: Improve documentation ([Vigil/vgld#1125](https://github.com/vigilnetwork/vgl/pull/1125))
- docs: Fix links for internal packages ([Vigil/vgld#1144](https://github.com/vigilnetwork/vgl/pull/1144))

### Developer-related package changes:

- chaingen: Add revocation generation infrastructure ([Vigil/vgld#1120](https://github.com/vigilnetwork/vgl/pull/1120))
- txscript: Add null data script creator ([Vigil/vgld#943](https://github.com/vigilnetwork/vgl/pull/943))
- txscript: Cleanup and improve NullDataScript tests ([Vigil/vgld#943](https://github.com/vigilnetwork/vgl/pull/943))
- txscript: Allow external signature hash calc ([Vigil/vgld#951](https://github.com/vigilnetwork/vgl/pull/951))
- secp256k1: update func signatures ([Vigil/vgld#934](https://github.com/vigilnetwork/vgl/pull/934))
- txscript: enforce MaxDataCarrierSize for GenerateProvablyPruneableOut ([Vigil/vgld#953](https://github.com/vigilnetwork/vgl/pull/953))
- txscript: Remove OP_SMALLDATA ([Vigil/vgld#954](https://github.com/vigilnetwork/vgl/pull/954))
- blockchain: Accept header in CheckProofOfWork ([Vigil/vgld#977](https://github.com/vigilnetwork/vgl/pull/977))
- blockchain: Make func definition style consistent ([Vigil/vgld#983](https://github.com/vigilnetwork/vgl/pull/983))
- blockchain: only fetch the parent block in BFFastAdd ([Vigil/vgld#972](https://github.com/vigilnetwork/vgl/pull/972))
- blockchain: Switch to FindSpentTicketsInBlock ([Vigil/vgld#915](https://github.com/vigilnetwork/vgl/pull/915))
- stake: Add Hash256PRNG init vector support ([Vigil/vgld#986](https://github.com/vigilnetwork/vgl/pull/986))
- blockchain/stake: Use Hash256PRNG init vector ([Vigil/vgld#987](https://github.com/vigilnetwork/vgl/pull/987))
- blockchain: Don't store full header in block node ([Vigil/vgld#988](https://github.com/vigilnetwork/vgl/pull/988))
- blockchain: Reconstruct headers from block nodes ([Vigil/vgld#989](https://github.com/vigilnetwork/vgl/pull/989))
- stake/multi: Don't return errors for IsX functions ([Vigil/vgld#995](https://github.com/vigilnetwork/vgl/pull/995))
- blockchain: Rename block index to main chain index ([Vigil/vgld#996](https://github.com/vigilnetwork/vgl/pull/996))
- blockchain: Refactor main block index logic ([Vigil/vgld#990](https://github.com/vigilnetwork/vgl/pull/990))
- blockchain: Use hash values in structs ([Vigil/vgld#992](https://github.com/vigilnetwork/vgl/pull/992))
- blockchain: Remove unused dump function ([Vigil/vgld#1001](https://github.com/vigilnetwork/vgl/pull/1001))
- blockchain: Generalize and optimize chain reorg ([Vigil/vgld#997](https://github.com/vigilnetwork/vgl/pull/997))
- blockchain: Pass parent block in connection code ([Vigil/vgld#998](https://github.com/vigilnetwork/vgl/pull/998))
- blockchain: Explicit block fetch semanticss ([Vigil/vgld#999](https://github.com/vigilnetwork/vgl/pull/999))
- blockchain: Use next detach block in reorg chain ([Vigil/vgld#1002](https://github.com/vigilnetwork/vgl/pull/1002))
- blockchain: Limit header sanity check to header ([Vigil/vgld#1003](https://github.com/vigilnetwork/vgl/pull/1003))
- blockchain: Validate num votes in header sanity ([Vigil/vgld#1005](https://github.com/vigilnetwork/vgl/pull/1005))
- blockchain: Validate max votes in header sanity ([Vigil/vgld#1006](https://github.com/vigilnetwork/vgl/pull/1006))
- blockchain: Validate stake diff in header context ([Vigil/vgld#1004](https://github.com/vigilnetwork/vgl/pull/1004))
- blockchain: No votes/revocations in header sanity ([Vigil/vgld#1007](https://github.com/vigilnetwork/vgl/pull/1007))
- blockchain: Validate max purchases in header sanity ([Vigil/vgld#1008](https://github.com/vigilnetwork/vgl/pull/1008))
- blockchain: Validate early votebits in header sanity ([Vigil/vgld#1009](https://github.com/vigilnetwork/vgl/pull/1009))
- blockchain: Validate block height in header context ([Vigil/vgld#1010](https://github.com/vigilnetwork/vgl/pull/1010))
- blockchain: Move check block context func ([Vigil/vgld#1011](https://github.com/vigilnetwork/vgl/pull/1011))
- blockchain: Block sanity cleanup and consistency ([Vigil/vgld#1012](https://github.com/vigilnetwork/vgl/pull/1012))
- blockchain: Remove dup ticket purchase value check ([Vigil/vgld#1013](https://github.com/vigilnetwork/vgl/pull/1013))
- blockchain: Only tickets before SVH in block sanity ([Vigil/vgld#1014](https://github.com/vigilnetwork/vgl/pull/1014))
- blockchain: Remove unused vote bits function ([Vigil/vgld#1015](https://github.com/vigilnetwork/vgl/pull/1015))
- blockchain: Move upgrade-only code to upgrade.go ([Vigil/vgld#1016](https://github.com/vigilnetwork/vgl/pull/1016))
- stake: Static assert of vote commitment ([Vigil/vgld#1020](https://github.com/vigilnetwork/vgl/pull/1020))
- blockchain: Remove unused error code ([Vigil/vgld#1021](https://github.com/vigilnetwork/vgl/pull/1021))
- blockchain: Improve readability of parent approval ([Vigil/vgld#1022](https://github.com/vigilnetwork/vgl/pull/1022))
- peer: rename mruinvmap, mrunoncemap to lruinvmap, lrunoncemap ([Vigil/vgld#976](https://github.com/vigilnetwork/vgl/pull/976))
- peer: rename noncemap to noncecache ([Vigil/vgld#976](https://github.com/vigilnetwork/vgl/pull/976))
- peer: rename inventorymap to inventorycache ([Vigil/vgld#976](https://github.com/vigilnetwork/vgl/pull/976))
- connmgr: convert state to atomic ([Vigil/vgld#1025](https://github.com/vigilnetwork/vgl/pull/1025))
- blockchain/mining: Full checks in CCB ([Vigil/vgld#1017](https://github.com/vigilnetwork/vgl/pull/1017))
- blockchain: Validate pool size in header context ([Vigil/vgld#1018](https://github.com/vigilnetwork/vgl/pull/1018))
- blockchain: Vote commitments in block sanity ([Vigil/vgld#1023](https://github.com/vigilnetwork/vgl/pull/1023))
- blockchain: Validate early final state is zero ([Vigil/vgld#1031](https://github.com/vigilnetwork/vgl/pull/1031))
- blockchain: Validate final state in header context ([Vigil/vgld#1034](https://github.com/vigilnetwork/vgl/pull/1033))
- blockchain: Max revocations in block sanity ([Vigil/vgld#1034](https://github.com/vigilnetwork/vgl/pull/1034))
- blockchain: Allowed stake txns in block sanity ([Vigil/vgld#1035](https://github.com/vigilnetwork/vgl/pull/1035))
- blockchain: Validate allowed votes in block context ([Vigil/vgld#1036](https://github.com/vigilnetwork/vgl/pull/1036))
- blockchain: Validate allowed revokes in blk contxt ([Vigil/vgld#1037](https://github.com/vigilnetwork/vgl/pull/1037))
- blockchain/stake: Rename tix spent to tix voted ([Vigil/vgld#1038](https://github.com/vigilnetwork/vgl/pull/1038))
- txscript: Require atomic swap contracts to specify the secret size ([Vigil/vgld#1039](https://github.com/vigilnetwork/vgl/pull/1039))
- blockchain: Remove unused struct ([Vigil/vgld#1043](https://github.com/vigilnetwork/vgl/pull/1043))
- blockchain: Store side chain blocks in database ([Vigil/vgld#1000](https://github.com/vigilnetwork/vgl/pull/1000))
- blockchain: Simplify initial chain state ([Vigil/vgld#1045](https://github.com/vigilnetwork/vgl/pull/1045))
- blockchain: Rework database versioning ([Vigil/vgld#1047](https://github.com/vigilnetwork/vgl/pull/1047))
- blockchain: Don't require chain for db upgrades ([Vigil/vgld#1051](https://github.com/vigilnetwork/vgl/pull/1051))
- blockchain/indexers: Allow interrupts ([Vigil/vgld#1052](https://github.com/vigilnetwork/vgl/pull/1052))
- blockchain: Remove old version information ([Vigil/vgld#1055](https://github.com/vigilnetwork/vgl/pull/1055))
- stake: optimize FindSpentTicketsInBlock slightly ([Vigil/vgld#1049](https://github.com/vigilnetwork/vgl/pull/1049))
- blockchain: Do not accept orphans/genesis block ([Vigil/vgld#1057](https://github.com/vigilnetwork/vgl/pull/1057))
- blockchain: Separate node ticket info population ([Vigil/vgld#1056](https://github.com/vigilnetwork/vgl/pull/1056))
- blockchain: Accept parent in blockNode constructor ([Vigil/vgld#1056](https://github.com/vigilnetwork/vgl/pull/1056))
- blockchain: Combine ErrDoubleSpend & ErrMissingTx ([Vigil/vgld#1064](https://github.com/vigilnetwork/vgl/pull/1064))
- blockchain: Calculate the lottery IV on demand ([Vigil/vgld#1065](https://github.com/vigilnetwork/vgl/pull/1065))
- blockchain: Simplify add/remove node logic ([Vigil/vgld#1067](https://github.com/vigilnetwork/vgl/pull/1067))
- blockchain: Infrastructure to manage block index ([Vigil/vgld#1044](https://github.com/vigilnetwork/vgl/pull/1044))
- blockchain: Add block validation status to index ([Vigil/vgld#1044](https://github.com/vigilnetwork/vgl/pull/1044))
- blockchain: Migrate to new block indexuse it ([Vigil/vgld#1044](https://github.com/vigilnetwork/vgl/pull/1044))
- blockchain: Lookup child in force head reorg ([Vigil/vgld#1070](https://github.com/vigilnetwork/vgl/pull/1070))
- blockchain: Refactor block idx entry serialization ([Vigil/vgld#1069](https://github.com/vigilnetwork/vgl/pull/1069))
- blockchain: Limit GetStakeVersions count ([Vigil/vgld#1071](https://github.com/vigilnetwork/vgl/pull/1071))
- blockchain: Remove dry run flag ([Vigil/vgld#1073](https://github.com/vigilnetwork/vgl/pull/1073))
- blockchain: Remove redundant stake ver calc func ([Vigil/vgld#1087](https://github.com/vigilnetwork/vgl/pull/1087))
- blockchain: Reduce GetGeneration to TipGeneration ([Vigil/vgld#1083](https://github.com/vigilnetwork/vgl/pull/1083))
- blockchain: Add chain tip tracking ([Vigil/vgld#1084](https://github.com/vigilnetwork/vgl/pull/1084))
- blockchain: Switch tip generation to chain tips ([Vigil/vgld#1085](https://github.com/vigilnetwork/vgl/pull/1085))
- blockchain: Simplify voter version calculation ([Vigil/vgld#1088](https://github.com/vigilnetwork/vgl/pull/1088))
- blockchain: Remove unused threshold serialization ([Vigil/vgld#1092](https://github.com/vigilnetwork/vgl/pull/1092))
- blockchain: Simplify chain tip tracking ([Vigil/vgld#1092](https://github.com/vigilnetwork/vgl/pull/1092))
- blockchain: Cache tip and parent at init ([Vigil/vgld#1100](https://github.com/vigilnetwork/vgl/pull/1100))
- mining: Obtain block by hash instead of top block ([Vigil/vgld#1094](https://github.com/vigilnetwork/vgl/pull/1094))
- blockchain: Remove unused GetTopBlock function ([Vigil/vgld#1094](https://github.com/vigilnetwork/vgl/pull/1094))
- multi: Rename BIP0111Version to NodeBloomVersion ([Vigil/vgld#1112](https://github.com/vigilnetwork/vgl/pull/1112))
- mining/mempool: Move priority code to mining pkg ([Vigil/vgld#1110](https://github.com/vigilnetwork/vgl/pull/1110))
- mining: Use single uint64 coinbase extra nonce ([Vigil/vgld#1116](https://github.com/vigilnetwork/vgl/pull/1116))
- mempool/mining: Clarify tree validity semantics ([Vigil/vgld#1118](https://github.com/vigilnetwork/vgl/pull/1118))
- mempool/mining: TxSource separation ([Vigil/vgld#1119](https://github.com/vigilnetwork/vgl/pull/1119))
- connmgr: Use same Dial func signature as net.Dial ([Vigil/vgld#1113](https://github.com/vigilnetwork/vgl/pull/1113))
- addrmgr: Declutter package API ([Vigil/vgld#1124](https://github.com/vigilnetwork/vgl/pull/1124))
- mining: Correct initial template generation ([Vigil/vgld#1122](https://github.com/vigilnetwork/vgl/pull/1122))
- cpuminer: Use header for extra nonce ([Vigil/vgld#1123](https://github.com/vigilnetwork/vgl/pull/1123))
- addrmgr: Make writing of peers file safer ([Vigil/vgld#1126](https://github.com/vigilnetwork/vgl/pull/1126))
- addrmgr: Save peers file only if necessary ([Vigil/vgld#1127](https://github.com/vigilnetwork/vgl/pull/1127))
- addrmgr: Factor out common code ([Vigil/vgld#1138](https://github.com/vigilnetwork/vgl/pull/1138))
- addrmgr: Improve isBad() performance ([Vigil/vgld#1134](https://github.com/vigilnetwork/vgl/pull/1134))
- VGLutil: Disallow creation of hybrid P2PK addrs ([Vigil/vgld#1154](https://github.com/vigilnetwork/vgl/pull/1154))
- chainec/VGLec: Remove hybrid pubkey support ([Vigil/vgld#1155](https://github.com/vigilnetwork/vgl/pull/1155))
- blockchain: Only fetch inputs once in connect txns ([Vigil/vgld#1152](https://github.com/vigilnetwork/vgl/pull/1152))
- indexers: Provide interface for index removal ([Vigil/vgld#1158](https://github.com/vigilnetwork/vgl/pull/1158))

### Testing and Quality Assurance:

- travis: set GOVERSION environment properly ([Vigil/vgld#958](https://github.com/vigilnetwork/vgl/pull/958))
- stake: Override false positive vet error ([Vigil/vgld#960](https://github.com/vigilnetwork/vgl/pull/960))
- docs: make example code compile ([Vigil/vgld#970](https://github.com/vigilnetwork/vgl/pull/970))
- blockchain: Add median time tests ([Vigil/vgld#991](https://github.com/vigilnetwork/vgl/pull/991))
- chaingen: Update vote commitments on hdr updates ([Vigil/vgld#1023](https://github.com/vigilnetwork/vgl/pull/1023))
- fullblocktests: Add tests for early final state ([Vigil/vgld#1031](https://github.com/vigilnetwork/vgl/pull/1031))
- travis: test in docker container ([Vigil/vgld#1053](https://github.com/vigilnetwork/vgl/pull/1053))
- blockchain: Correct error stringer tests ([Vigil/vgld#1066](https://github.com/vigilnetwork/vgl/pull/1066))
- blockchain: Remove superfluous reorg tests ([Vigil/vgld#1072](https://github.com/vigilnetwork/vgl/pull/1072))
- blockchain: Use chaingen for forced reorg tests ([Vigil/vgld#1074](https://github.com/vigilnetwork/vgl/pull/1074))
- blockchain: Remove superfluous test checks ([Vigil/vgld#1075](https://github.com/vigilnetwork/vgl/pull/1075))
- blockchain: move block validation rule tests into fullblocktests ([Vigil/vgld#1060](https://github.com/vigilnetwork/vgl/pull/1060))
- fullblocktests: Cleanup after refactor ([Vigil/vgld#1080](https://github.com/vigilnetwork/vgl/pull/1080))
- chaingen: Prevent dup block names in NextBlock ([Vigil/vgld#1079](https://github.com/vigilnetwork/vgl/pull/1079))
- blockchain: Remove duplicate val tests ([Vigil/vgld#1082](https://github.com/vigilnetwork/vgl/pull/1082))
- chaingen: Break dependency on blockchain ([Vigil/vgld#1076](https://github.com/vigilnetwork/vgl/pull/1076))
- blockchain: Consolidate tests into the main package ([Vigil/vgld#1077](https://github.com/vigilnetwork/vgl/pull/1077))
- chaingen: Export vote commitment script function ([Vigil/vgld#1081](https://github.com/vigilnetwork/vgl/pull/1081))
- fullblocktests: Improve vote on wrong block tests ([Vigil/vgld#1081](https://github.com/vigilnetwork/vgl/pull/1081))
- chaingen: Export func to check if block is solved ([Vigil/vgld#1089](https://github.com/vigilnetwork/vgl/pull/1089))
- fullblocktests: Use new exported IsSolved func ([Vigil/vgld#1089](https://github.com/vigilnetwork/vgl/pull/1089))
- chaingen: Accept mungers for create premine block ([Vigil/vgld#1090](https://github.com/vigilnetwork/vgl/pull/1090))
- blockchain: Add tests for chain tip tracking ([Vigil/vgld#1096](https://github.com/vigilnetwork/vgl/pull/1096))
- blockchain: move block validation rule tests into fullblocktests (2/x) ([Vigil/vgld#1095](https://github.com/vigilnetwork/vgl/pull/1095))
- addrmgr: Remove obsolete coverage script ([Vigil/vgld#1103](https://github.com/vigilnetwork/vgl/pull/1103))
- chaingen: Track expected blk heights separately ([Vigil/vgld#1101](https://github.com/vigilnetwork/vgl/pull/1101))
- addrmgr: Improve test coverage ([Vigil/vgld#1111](https://github.com/vigilnetwork/vgl/pull/1111))
- chaingen: Add revocation generation infrastructure ([Vigil/vgld#1120](https://github.com/vigilnetwork/vgl/pull/1120))
- fullblocktests: Add some basic revocation tests ([Vigil/vgld#1121](https://github.com/vigilnetwork/vgl/pull/1121))
- addrmgr: Test removal of corrupt peers file ([Vigil/vgld#1129](https://github.com/vigilnetwork/vgl/pull/1129))

### Misc:

- release: Bump for v1.2.0 ([Vigil/vgld#1140](https://github.com/vigilnetwork/vgl/pull/1140))
- goimports -w . ([Vigil/vgld#968](https://github.com/vigilnetwork/vgl/pull/968))
- dep: sync ([Vigil/vgld#980](https://github.com/vigilnetwork/vgl/pull/980))
- multi: Simplify code per gosimple linter ([Vigil/vgld#993](https://github.com/vigilnetwork/vgl/pull/993))
- multi: various cleanups ([Vigil/vgld#1019](https://github.com/vigilnetwork/vgl/pull/1019))
- multi: release the mutex earlier ([Vigil/vgld#1026](https://github.com/vigilnetwork/vgl/pull/1026))
- multi: fix some maligned linter warnings ([Vigil/vgld#1025](https://github.com/vigilnetwork/vgl/pull/1025))
- blockchain: Correct a few log statements ([Vigil/vgld#1042](https://github.com/vigilnetwork/vgl/pull/1042))
- mempool: cleaner ([Vigil/vgld#1050](https://github.com/vigilnetwork/vgl/pull/1050))
- multi: fix misspell linter warnings ([Vigil/vgld#1054](https://github.com/vigilnetwork/vgl/pull/1054))
- dep: sync ([Vigil/vgld#1091](https://github.com/vigilnetwork/vgl/pull/1091))
- multi: Properly capitalize Vigil ([Vigil/vgld#1102](https://github.com/vigilnetwork/vgl/pull/1102))
- build: Correct semver build handling ([Vigil/vgld#1097](https://github.com/vigilnetwork/vgl/pull/1097))
- main: Make func definition style consistent ([Vigil/vgld#1114](https://github.com/vigilnetwork/vgl/pull/1114))
- main: Allow semver prerel via linker flags ([Vigil/vgld#1128](https://github.com/vigilnetwork/vgl/pull/1128))

### Code Contributors (alphabetical order):

- Andrew Chiw
- Daniel Krawsiz
- Dave Collins
- David Hill
- Donald Adu-Poku
- Javed Khan
- Jolan Luff
- Jon Gillham
- Josh Rickmar
- Markus Richter
- Matheus Degiovani
- Ryan Vacek




