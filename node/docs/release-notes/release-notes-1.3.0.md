# vgld v1.3.0

This release of vgld contains significant performance enhancements for startup
speed, validation, and network operations that directly benefit lightweight
clients, such as SPV (Simplified Payment Verification) wallets, a policy change
to reduce the default minimum transaction fee rate, a new public test network
version, removal of bloom filter support, infrastructure improvements, and other
quality assurance changes.

**It is highly recommended that everyone upgrade to this latest release as it
contains many important scalability improvements and is required to be able to
use the new public test network.**

## Downgrade Warning

The database format in v1.3.0 is not compatible with previous versions of the
software.  This only affects downgrades as users upgrading from previous
versions will see a one time database migration.

Once this migration has been completed, it will no longer be possible to
downgrade to a previous version of the software without having to delete the
database and redownload the chain.

## Notable Changes

### Reduction of Default Minimum Transaction Fee Rate Policy

The default setting for the policy which specifies the minimum transaction fee
rate that will be accepted and relayed to the rest of the network has been
reduced to 0.0001 VGL/kB (10,000 atoms/kB) from the previous value of 0.001
VGL/kB (100,000 atoms/kB).

Transactions should not attempt to use the reduced fee rate until the majority
of the network has upgraded to this release as otherwise the transactions will
likely have issues relaying through the network since old nodes that have not
updated their policy will reject them due to not paying a high enough fee.

### Several Speed Optimizations

This release contains several enhancements to improve speed for startup,
the initial sync process, validation, and network operations.

In order to achieve these speedups, there is a one time database migration, as
previously mentioned, that typically only takes a few seconds to complete on
most hardware.

#### Further Improved Startup Speed

The startup time has been improved by roughly 2x on both slower hard disk drives
(HDDs) and solid state drives (SSDs) as compared to v1.2.0.

#### Significantly Faster Network Operations

The ability to serve information to other peers on the network has received
several optimizations which, in addition to generally improving the overall
scalability and throughput of the network, also directly benefits SPV
(Simplified Payment Verification) clients by delivering the block headers they
require roughly 3x to 4x faster.

#### Signature Hash Calculation Optimization

Part of validating that transactions are only spending coins that the owner has
authorized involves ensuring the validity of cryptographic signatures.  This
release provides a speedup of about 75% to a key portion of that validation
which results in a roughly 20% faster initial sync process.

### Bloom Filters Removal

Bloom filters were deprecated as of the last release in favor of the more recent
privacy-preserving GCS committed filters.  Consequently, this release removes
support for bloom filters completely.  There are no known clients which use
bloom filters, however, if there are any unknown clients which use them, those
clients will need to be updated to use the GCS committed filters accordingly.

### Public Test Network Version 3

The public test network has been reset and bumped to version 3.  All of the new
consensus rules voted in by version 2 of the public test network have been
retained and are therefore active on the new version 3 test network without
having to vote them in again.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/release-v1.2.0...release-v1.3.0).

### Protocol and network:

- chaincfg: Add checkpoints for 1.3.0 release ([Vigil/vgld#1385](https://github.com/vigilnetwork/vgl/pull/1385))
- multi: Remove everything to do about bloom filters ([Vigil/vgld#1162](https://github.com/vigilnetwork/vgl/pull/1162))
- wire: Remove TxSerializeWitnessSigning ([Vigil/vgld#1180](https://github.com/vigilnetwork/vgl/pull/1180))
- addrmgr: Skip low quality addresses for getaddr ([Vigil/vgld#1135](https://github.com/vigilnetwork/vgl/pull/1135))
- addrmgr: Fix race in save peers ([Vigil/vgld#1259](https://github.com/vigilnetwork/vgl/pull/1259))
- server: Only respond to getaddr once per conn ([Vigil/vgld#1257](https://github.com/vigilnetwork/vgl/pull/1257))
- peer: Rework version negotiation ([Vigil/vgld#1250](https://github.com/vigilnetwork/vgl/pull/1250))
- peer: Allow OnVersion callback to reject peer ([Vigil/vgld#1251](https://github.com/vigilnetwork/vgl/pull/1251))
- server: Reject outbound conns to non-full nodes ([Vigil/vgld#1252](https://github.com/vigilnetwork/vgl/pull/1252))
- peer: Improve net address service adverts ([Vigil/vgld#1253](https://github.com/vigilnetwork/vgl/pull/1253))
- addrmgr: Expose method to update services ([Vigil/vgld#1254](https://github.com/vigilnetwork/vgl/pull/1254))
- server: Update addrmgr services on outbound conns ([Vigil/vgld#1254](https://github.com/vigilnetwork/vgl/pull/1254))
- server: Use local inbound var in version handler ([Vigil/vgld#1255](https://github.com/vigilnetwork/vgl/pull/1255))
- server: Only advertise local addr when current ([Vigil/vgld#1256](https://github.com/vigilnetwork/vgl/pull/1256))
- server: Use local addr var in version handler ([Vigil/vgld#1258](https://github.com/vigilnetwork/vgl/pull/1258))
- chaincfg: split params into per-network files ([Vigil/vgld#1265](https://github.com/vigilnetwork/vgl/pull/1265))
- server: Always reply to getheaders with headers ([Vigil/vgld#1295](https://github.com/vigilnetwork/vgl/pull/1295))
- addrmgr: skip never-successful addresses ([Vigil/vgld#1313](https://github.com/vigilnetwork/vgl/pull/1313))
- multi: Introduce default coin type for SLIP0044 ([Vigil/vgld#1293](https://github.com/vigilnetwork/vgl/pull/1293))
- blockchain: Modify diff redux logic for testnet ([Vigil/vgld#1387](https://github.com/vigilnetwork/vgl/pull/1387))
- multi: Reset testnet and bump to version 3 ([Vigil/vgld#1387](https://github.com/vigilnetwork/vgl/pull/1387))
- multi: Remove testnet version 2 defs and refs ([Vigil/vgld#1387](https://github.com/vigilnetwork/vgl/pull/1387))

### Transaction relay (memory pool):

- policy: Lower default relay fee to 0.0001/kB ([Vigil/vgld#1202](https://github.com/vigilnetwork/vgl/pull/1202))
- mempool: Use blockchain for tx expiry check ([Vigil/vgld#1199](https://github.com/vigilnetwork/vgl/pull/1199))
- mempool: use secp256k1 functions directly ([Vigil/vgld#1213](https://github.com/vigilnetwork/vgl/pull/1213))
- mempool: Make expiry pruning self contained ([Vigil/vgld#1378](https://github.com/vigilnetwork/vgl/pull/1378))
- mempool: Stricter orphan evaluation and eviction ([Vigil/vgld#1207](https://github.com/vigilnetwork/vgl/pull/1207))
- mempool: use secp256k1 functions directly ([Vigil/vgld#1213](https://github.com/vigilnetwork/vgl/pull/1213))
- multi: add specialized rebroadcast handling for stake txs ([Vigil/vgld#979](https://github.com/vigilnetwork/vgl/pull/979))
- mempool: Make expiry pruning self contained ([Vigil/vgld#1378](https://github.com/vigilnetwork/vgl/pull/1378))

### RPC:

- rpcserver: Improve JSON-RPC compatibility ([Vigil/vgld#1150](https://github.com/vigilnetwork/vgl/pull/1150))
- rpcserver: Correct rebroadcastwinners handler ([Vigil/vgld#1234](https://github.com/vigilnetwork/vgl/pull/1234))
- VGLjson: Add Expiry field to CreateRawTransactionCmd ([Vigil/vgld#1149](https://github.com/vigilnetwork/vgl/pull/1149))
- VGLjson: add estimatesmartfee ([Vigil/vgld#1201](https://github.com/vigilnetwork/vgl/pull/1201))
- rpc: Use upstream gorilla/websocket ([Vigil/vgld#1218](https://github.com/vigilnetwork/vgl/pull/1218))
- VGLjson: add createvotingaccount and dropvotingaccount rpc methods ([Vigil/vgld#1217](https://github.com/vigilnetwork/vgl/pull/1217))
- multi: Change NoSplitTransaction param to SplitTx ([Vigil/vgld#1231](https://github.com/vigilnetwork/vgl/pull/1231))
- rpcclient: pass default value for NewPurchaseTicketCmd's comment param ([Vigil/vgld#1232](https://github.com/vigilnetwork/vgl/pull/1232))
- multi: No winning ticket ntfns for big reorg depth ([Vigil/vgld#1235](https://github.com/vigilnetwork/vgl/pull/1235))
- multi: modify PurchaseTicketCmd ([Vigil/vgld#1241](https://github.com/vigilnetwork/vgl/pull/1241))
- multi: move extension commands into associated normal command files ([Vigil/vgld#1238](https://github.com/vigilnetwork/vgl/pull/1238))
- VGLjson: Fix NewCreateRawTransactionCmd comment ([Vigil/vgld#1262](https://github.com/vigilnetwork/vgl/pull/1262))
- multi: revert TicketChange addition to PurchaseTicketCmd ([Vigil/vgld#1278](https://github.com/vigilnetwork/vgl/pull/1278))
- rpcclient: Implement fmt.Stringer for Client ([Vigil/vgld#1298](https://github.com/vigilnetwork/vgl/pull/1298))
- multi: add amount field to TransactionInput ([Vigil/vgld#1297](https://github.com/vigilnetwork/vgl/pull/1297))
- VGLjson: Ready GetStakeInfoResult for SPV wallets ([Vigil/vgld#1333](https://github.com/vigilnetwork/vgl/pull/1333))
- VGLjson: add fundrawtransaction command ([Vigil/vgld#1316](https://github.com/vigilnetwork/vgl/pull/1316))
- VGLjson: Make linter happy by renaming Id to ID ([Vigil/vgld#1374](https://github.com/vigilnetwork/vgl/pull/1374))
- VGLjson: Remove unused vote bit concat codec funcs ([Vigil/vgld#1384](https://github.com/vigilnetwork/vgl/pull/1384))
- rpcserver: Cleanup cfilter handling ([Vigil/vgld#1398](https://github.com/vigilnetwork/vgl/pull/1398))

### vgld command-line flags and configuration:

- multi: Correct clean and expand path handling ([Vigil/vgld#1186](https://github.com/vigilnetwork/vgl/pull/1186))

### vglctl utility changes:

- vglctl: Fix --skipverify failing if rpc.cert not found ([Vigil/vgld#1163](https://github.com/vigilnetwork/vgl/pull/1163))

### Documentation:

- hdkeychain: Correct hash algorithm in comment ([Vigil/vgld#1171](https://github.com/vigilnetwork/vgl/pull/1171))
- Fix typo in blockchain ([Vigil/vgld#1185](https://github.com/vigilnetwork/vgl/pull/1185))
- docs: Update node.js example for v8.11.1 LTS ([Vigil/vgld#1209](https://github.com/vigilnetwork/vgl/pull/1209))
- docs: Update txaccepted method in json_rpc_api.md ([Vigil/vgld#1242](https://github.com/vigilnetwork/vgl/pull/1242))
- docs: Correct blockmaxsize and blockprioritysize ([Vigil/vgld#1339](https://github.com/vigilnetwork/vgl/pull/1339))
- server: Correct comment in getblocks handler ([Vigil/vgld#1269](https://github.com/vigilnetwork/vgl/pull/1269))
- config: Fix typo ([Vigil/vgld#1274](https://github.com/vigilnetwork/vgl/pull/1274))
- multi: Fix badges in README ([Vigil/vgld#1277](https://github.com/vigilnetwork/vgl/pull/1277))
- stake: Correct comment in connectNode ([Vigil/vgld#1325](https://github.com/vigilnetwork/vgl/pull/1325))
- txscript: Update comments for removal of flags ([Vigil/vgld#1336](https://github.com/vigilnetwork/vgl/pull/1336))
- docs: Update docs for versioned modules ([Vigil/vgld#1391](https://github.com/vigilnetwork/vgl/pull/1391))
- mempool: Correct min relay tx fee comment to VGL ([Vigil/vgld#1396](https://github.com/vigilnetwork/vgl/pull/1396))

### Developer-related package and module changes:

- blockchain: CheckConnectBlockTemplate with tests ([Vigil/vgld#1086](https://github.com/vigilnetwork/vgl/pull/1086))
- addrmgr: Simplify package API ([Vigil/vgld#1136](https://github.com/vigilnetwork/vgl/pull/1136))
- txscript: Remove unused strict multisig flag ([Vigil/vgld#1203](https://github.com/vigilnetwork/vgl/pull/1203))
- txscript: Move sig hash logic to separate file ([Vigil/vgld#1174](https://github.com/vigilnetwork/vgl/pull/1174))
- txscript: Remove SigHashAllValue ([Vigil/vgld#1175](https://github.com/vigilnetwork/vgl/pull/1175))
- txscript: Decouple and optimize sighash calc ([Vigil/vgld#1179](https://github.com/vigilnetwork/vgl/pull/1179))
- wire: Remove TxSerializeWitnessValueSigning ([Vigil/vgld#1176](https://github.com/vigilnetwork/vgl/pull/1176))
- hdkeychain: Satisfy fmt.Stringer interface ([Vigil/vgld#1168](https://github.com/vigilnetwork/vgl/pull/1168))
- blockchain: Validate tx expiry in block context ([Vigil/vgld#1187](https://github.com/vigilnetwork/vgl/pull/1187))
- blockchain: rename ErrRegTxSpendStakeOut to ErrRegTxCreateStakeOut ([Vigil/vgld#1195](https://github.com/vigilnetwork/vgl/pull/1195))
- multi: Break coinbase dep on standardness rules ([Vigil/vgld#1196](https://github.com/vigilnetwork/vgl/pull/1196))
- txscript: Cleanup code for the substr opcode ([Vigil/vgld#1206](https://github.com/vigilnetwork/vgl/pull/1206))
- multi: use secp256k1 types and fields directly ([Vigil/vgld#1211](https://github.com/vigilnetwork/vgl/pull/1211))
- VGLec: add Pubkey func to secp256k1 and edwards elliptic curves ([Vigil/vgld#1214](https://github.com/vigilnetwork/vgl/pull/1214))
- blockchain: use secp256k1 functions directly ([Vigil/vgld#1212](https://github.com/vigilnetwork/vgl/pull/1212))
- multi: Replace btclog with slog ([Vigil/vgld#1216](https://github.com/vigilnetwork/vgl/pull/1216))
- multi: Define vgo modules ([Vigil/vgld#1223](https://github.com/vigilnetwork/vgl/pull/1223))
- chainhash: Define vgo module ([Vigil/vgld#1224](https://github.com/vigilnetwork/vgl/pull/1224))
- wire: Refine vgo deps ([Vigil/vgld#1225](https://github.com/vigilnetwork/vgl/pull/1225))
- addrmrg: Refine vgo deps ([Vigil/vgld#1226](https://github.com/vigilnetwork/vgl/pull/1226))
- chaincfg: Refine vgo deps ([Vigil/vgld#1227](https://github.com/vigilnetwork/vgl/pull/1227))
- multi: Return fork len from ProcessBlock ([Vigil/vgld#1233](https://github.com/vigilnetwork/vgl/pull/1233))
- blockchain: Panic on fatal assertions ([Vigil/vgld#1243](https://github.com/vigilnetwork/vgl/pull/1243))
- blockchain: Convert to full block index in mem ([Vigil/vgld#1229](https://github.com/vigilnetwork/vgl/pull/1229))
- blockchain: Optimize checkpoint handling ([Vigil/vgld#1230](https://github.com/vigilnetwork/vgl/pull/1230))
- blockchain: Optimize block locator generation ([Vigil/vgld#1237](https://github.com/vigilnetwork/vgl/pull/1237))
- multi: Refactor and optimize inv discovery ([Vigil/vgld#1239](https://github.com/vigilnetwork/vgl/pull/1239))
- peer: Minor function definition order cleanup ([Vigil/vgld#1247](https://github.com/vigilnetwork/vgl/pull/1247))
- peer: Remove superfluous dup version check ([Vigil/vgld#1248](https://github.com/vigilnetwork/vgl/pull/1248))
- txscript: export canonicalDataSize ([Vigil/vgld#1266](https://github.com/vigilnetwork/vgl/pull/1266))
- blockchain: Add BuildMerkleTreeStore alternative for MsgTx ([Vigil/vgld#1268](https://github.com/vigilnetwork/vgl/pull/1268))
- blockchain: Optimize exported header access ([Vigil/vgld#1273](https://github.com/vigilnetwork/vgl/pull/1273))
- txscript: Cleanup P2SH and stake opcode handling ([Vigil/vgld#1318](https://github.com/vigilnetwork/vgl/pull/1318))
- txscript: Significantly improve errors ([Vigil/vgld#1319](https://github.com/vigilnetwork/vgl/pull/1319))
- txscript: Remove pay-to-script-hash flag ([Vigil/vgld#1321](https://github.com/vigilnetwork/vgl/pull/1321))
- txscript: Remove DER signature verification flag ([Vigil/vgld#1323](https://github.com/vigilnetwork/vgl/pull/1323))
- txscript: Remove verify minimal data flag ([Vigil/vgld#1326](https://github.com/vigilnetwork/vgl/pull/1326))
- txscript: Remove script num require minimal flag ([Vigil/vgld#1328](https://github.com/vigilnetwork/vgl/pull/1328))
- txscript: Make PeekInt consistent with PopInt ([Vigil/vgld#1329](https://github.com/vigilnetwork/vgl/pull/1329))
- build: Add experimental support for vgo ([Vigil/vgld#1215](https://github.com/vigilnetwork/vgl/pull/1215))
- build: Update some vgo dependencies to use tags ([Vigil/vgld#1219](https://github.com/vigilnetwork/vgl/pull/1219))
- stake: add ExpiredByBlock to stake.Node ([Vigil/vgld#1221](https://github.com/vigilnetwork/vgl/pull/1221))
- server: Minor function definition order cleanup ([Vigil/vgld#1271](https://github.com/vigilnetwork/vgl/pull/1271))
- server: Convert CF code to use new inv discovery ([Vigil/vgld#1272](https://github.com/vigilnetwork/vgl/pull/1272))
- multi: add valueIn parameter to wire.NewTxIn ([Vigil/vgld#1287](https://github.com/vigilnetwork/vgl/pull/1287))
- txscript: Remove low S verification flag ([Vigil/vgld#1308](https://github.com/vigilnetwork/vgl/pull/1308))
- txscript: Remove unused old sig hash type ([Vigil/vgld#1309](https://github.com/vigilnetwork/vgl/pull/1309))
- txscript: Remove strict encoding verification flag ([Vigil/vgld#1310](https://github.com/vigilnetwork/vgl/pull/1310))
- blockchain: Combine block by hash functions ([Vigil/vgld#1330](https://github.com/vigilnetwork/vgl/pull/1330))
- multi: Continue conversion from chainec to VGLec ([Vigil/vgld#1304](https://github.com/vigilnetwork/vgl/pull/1304))
- multi: Remove unused secp256k1 sig parse parameter ([Vigil/vgld#1335](https://github.com/vigilnetwork/vgl/pull/1335))
- blockchain: Refactor db main chain idx to blk idx ([Vigil/vgld#1332](https://github.com/vigilnetwork/vgl/pull/1332))
- blockchain: Remove main chain index from db ([Vigil/vgld#1334](https://github.com/vigilnetwork/vgl/pull/1334))
- blockchain: Implement new chain view ([Vigil/vgld#1337](https://github.com/vigilnetwork/vgl/pull/1337))
- blockmanager: remove unused Pause() API ([Vigil/vgld#1340](https://github.com/vigilnetwork/vgl/pull/1340))
- chainhash: Remove dup code from hash funcs ([Vigil/vgld#1342](https://github.com/vigilnetwork/vgl/pull/1342))
- connmgr: Fix the ConnReq print out causing panic ([Vigil/vgld#1345](https://github.com/vigilnetwork/vgl/pull/1345))
- gcs: Pool MatchAny data allocations ([Vigil/vgld#1348](https://github.com/vigilnetwork/vgl/pull/1348))
- blockchain: Faster chain view block locator ([Vigil/vgld#1338](https://github.com/vigilnetwork/vgl/pull/1338))
- blockchain: Refactor to use new chain view ([Vigil/vgld#1344](https://github.com/vigilnetwork/vgl/pull/1344))
- blockchain: Remove unnecessary genesis block check ([Vigil/vgld#1368](https://github.com/vigilnetwork/vgl/pull/1368))
- chainhash: Update go build module support ([Vigil/vgld#1358](https://github.com/vigilnetwork/vgl/pull/1358))
- wire: Update go build module support ([Vigil/vgld#1359](https://github.com/vigilnetwork/vgl/pull/1359))
- addrmgr: Update go build module support ([Vigil/vgld#1360](https://github.com/vigilnetwork/vgl/pull/1360))
- chaincfg: Update go build module support ([Vigil/vgld#1361](https://github.com/vigilnetwork/vgl/pull/1361))
- connmgr: Refine go build module support ([Vigil/vgld#1363](https://github.com/vigilnetwork/vgl/pull/1363))
- secp256k1: Refine go build module support ([Vigil/vgld#1362](https://github.com/vigilnetwork/vgl/pull/1362))
- VGLec: Refine go build module support ([Vigil/vgld#1364](https://github.com/vigilnetwork/vgl/pull/1364))
- certgen: Update go build module support ([Vigil/vgld#1365](https://github.com/vigilnetwork/vgl/pull/1365))
- VGLutil: Refine go build module support ([Vigil/vgld#1366](https://github.com/vigilnetwork/vgl/pull/1366))
- hdkeychain: Refine go build module support ([Vigil/vgld#1369](https://github.com/vigilnetwork/vgl/pull/1369))
- txscript: Refine go build module support ([Vigil/vgld#1370](https://github.com/vigilnetwork/vgl/pull/1370))
- multi: Remove go modules that do not build ([Vigil/vgld#1371](https://github.com/vigilnetwork/vgl/pull/1371))
- database: Refine go build module support ([Vigil/vgld#1372](https://github.com/vigilnetwork/vgl/pull/1372))
- build: Refine build module support ([Vigil/vgld#1384](https://github.com/vigilnetwork/vgl/pull/1384))
- blockmanager: make pruning transactions consistent ([Vigil/vgld#1376](https://github.com/vigilnetwork/vgl/pull/1376))
- blockchain: Optimize reorg to use known status ([Vigil/vgld#1367](https://github.com/vigilnetwork/vgl/pull/1367))
- blockchain: Make block index flushable ([Vigil/vgld#1375](https://github.com/vigilnetwork/vgl/pull/1375))
- blockchain: Mark fastadd block valid ([Vigil/vgld#1392](https://github.com/vigilnetwork/vgl/pull/1392))
- release: Bump module versions and deps ([Vigil/vgld#1390](https://github.com/vigilnetwork/vgl/pull/1390))
- blockchain: Mark fastadd block valid ([Vigil/vgld#1392](https://github.com/vigilnetwork/vgl/pull/1392))
- gcs: use dchest/siphash ([Vigil/vgld#1395](https://github.com/vigilnetwork/vgl/pull/1395))
- VGLec: Make function defs more consistent ([Vigil/vgld#1432](https://github.com/vigilnetwork/vgl/pull/1432))

### Testing and Quality Assurance:

- addrmgr: Simplify tests for KnownAddress ([Vigil/vgld#1133](https://github.com/vigilnetwork/vgl/pull/1133))
- blockchain: move block validation rule tests into fullblocktests ([Vigil/vgld#1141](https://github.com/vigilnetwork/vgl/pull/1141))
- addrmgr: Test timestamp update during AddAddress ([Vigil/vgld#1137](https://github.com/vigilnetwork/vgl/pull/1137))
- txscript: Consolidate tests into txscript package ([Vigil/vgld#1177](https://github.com/vigilnetwork/vgl/pull/1177))
- txscript: Add JSON-based signature hash tests ([Vigil/vgld#1178](https://github.com/vigilnetwork/vgl/pull/1178))
- txscript: Correct JSON-based signature hash tests ([Vigil/vgld#1181](https://github.com/vigilnetwork/vgl/pull/1181))
- txscript: Add benchmark for sighash calculation ([Vigil/vgld#1179](https://github.com/vigilnetwork/vgl/pull/1179))
- mempool: Refactor pool membership test logic ([Vigil/vgld#1188](https://github.com/vigilnetwork/vgl/pull/1188))
- blockchain: utilize CalcNextReqStakeDifficulty in fullblocktests ([Vigil/vgld#1189](https://github.com/vigilnetwork/vgl/pull/1189))
- fullblocktests: add additional premine and malformed tests ([Vigil/vgld#1190](https://github.com/vigilnetwork/vgl/pull/1190))
- txscript: Improve substr opcode test coverage ([Vigil/vgld#1205](https://github.com/vigilnetwork/vgl/pull/1205))
- txscript: Convert reference tests to new format ([Vigil/vgld#1320](https://github.com/vigilnetwork/vgl/pull/1320))
- txscript: Remove P2SH flag from test data ([Vigil/vgld#1322](https://github.com/vigilnetwork/vgl/pull/1322))
- txscript: Remove DERSIG flag from test data ([Vigil/vgld#1324](https://github.com/vigilnetwork/vgl/pull/1324))
- txscript: Remove MINIMALDATA flag from test data ([Vigil/vgld#1327](https://github.com/vigilnetwork/vgl/pull/1327))
- fullblocktests: Add expired stake tx test ([Vigil/vgld#1184](https://github.com/vigilnetwork/vgl/pull/1184))
- travis: simplify Docker files ([Vigil/vgld#1275](https://github.com/vigilnetwork/vgl/pull/1275))
- docker: Add dockerfiles for running vgld nodes ([Vigil/vgld#1317](https://github.com/vigilnetwork/vgl/pull/1317))
- blockchain: Improve spend journal tests ([Vigil/vgld#1246](https://github.com/vigilnetwork/vgl/pull/1246))
- txscript: Cleanup and add tests for left opcode ([Vigil/vgld#1281](https://github.com/vigilnetwork/vgl/pull/1281))
- txscript: Cleanup and add tests for right opcode ([Vigil/vgld#1282](https://github.com/vigilnetwork/vgl/pull/1282))
- txscript: Cleanup and add tests for the cat opcode ([Vigil/vgld#1283](https://github.com/vigilnetwork/vgl/pull/1283))
- txscript: Cleanup and add tests for rotr opcode ([Vigil/vgld#1285](https://github.com/vigilnetwork/vgl/pull/1285))
- txscript: Cleanup and add tests for rotl opcode ([Vigil/vgld#1286](https://github.com/vigilnetwork/vgl/pull/1286))
- txscript: Cleanup and add tests for lshift opcode ([Vigil/vgld#1288](https://github.com/vigilnetwork/vgl/pull/1288))
- txscript: Cleanup and add tests for rshift opcode ([Vigil/vgld#1289](https://github.com/vigilnetwork/vgl/pull/1289))
- txscript: Cleanup and add tests for div opcode ([Vigil/vgld#1290](https://github.com/vigilnetwork/vgl/pull/1290))
- txscript: Cleanup and add tests for mod opcode ([Vigil/vgld#1291](https://github.com/vigilnetwork/vgl/pull/1291))
- txscript: Update CSV to match tests in VGLP0003 ([Vigil/vgld#1292](https://github.com/vigilnetwork/vgl/pull/1292))
- txscript: Introduce repeated syntax to test data ([Vigil/vgld#1299](https://github.com/vigilnetwork/vgl/pull/1299))
- txscript: Allow multi opcode test data repeat ([Vigil/vgld#1300](https://github.com/vigilnetwork/vgl/pull/1300))
- txscript: Improve and correct some script tests ([Vigil/vgld#1303](https://github.com/vigilnetwork/vgl/pull/1303))
- main: verify network pow limits ([Vigil/vgld#1302](https://github.com/vigilnetwork/vgl/pull/1302))
- txscript: Remove STRICTENC flag from test data ([Vigil/vgld#1311](https://github.com/vigilnetwork/vgl/pull/1311))
- txscript: Cleanup plus tests for checksig opcodes ([Vigil/vgld#1315](https://github.com/vigilnetwork/vgl/pull/1315))
- blockchain: Add negative tests for forced reorg ([Vigil/vgld#1341](https://github.com/vigilnetwork/vgl/pull/1341))
- VGLjson: Consolidate tests into VGLjson package ([Vigil/vgld#1373](https://github.com/vigilnetwork/vgl/pull/1373))
- txscript: add additional data push op code tests ([Vigil/vgld#1346](https://github.com/vigilnetwork/vgl/pull/1346))
- txscript: add/group control op code tests ([Vigil/vgld#1349](https://github.com/vigilnetwork/vgl/pull/1349))
- txscript: add/group stack op code tests ([Vigil/vgld#1350](https://github.com/vigilnetwork/vgl/pull/1350))
- txscript: group splice opcode tests ([Vigil/vgld#1351](https://github.com/vigilnetwork/vgl/pull/1351))
- txscript: add/group bitwise logic, comparison & rotation op code tests ([Vigil/vgld#1352](https://github.com/vigilnetwork/vgl/pull/1352))
- txscript: add/group numeric related opcode tests ([Vigil/vgld#1353](https://github.com/vigilnetwork/vgl/pull/1353))
- txscript: group reserved op code tests ([Vigil/vgld#1355](https://github.com/vigilnetwork/vgl/pull/1355))
- txscript: add/group crypto related op code tests ([Vigil/vgld#1354](https://github.com/vigilnetwork/vgl/pull/1354))
- multi: Reduce testnet2 refs in unit tests ([Vigil/vgld#1387](https://github.com/vigilnetwork/vgl/pull/1387))
- blockchain: Avoid deployment expiration in tests ([Vigil/vgld#1450](https://github.com/vigilnetwork/vgl/pull/1450))

### Misc:

- release: Bump for v1.3.0 ([Vigil/vgld#1388](https://github.com/vigilnetwork/vgl/pull/1388))
- multi: Correct typos found by misspell ([Vigil/vgld#1197](https://github.com/vigilnetwork/vgl/pull/1197))
- main: Correct mem profile error message ([Vigil/vgld#1183](https://github.com/vigilnetwork/vgl/pull/1183))
- multi: Use saner permissions saving certs ([Vigil/vgld#1263](https://github.com/vigilnetwork/vgl/pull/1263))
- server: only call time.Now() once ([Vigil/vgld#1313](https://github.com/vigilnetwork/vgl/pull/1313))
- multi: linter cleanup ([Vigil/vgld#1305](https://github.com/vigilnetwork/vgl/pull/1305))
- multi: Remove unnecessary network name funcs ([Vigil/vgld#1387](https://github.com/vigilnetwork/vgl/pull/1387))
- config: Warn if testnet2 database exists ([Vigil/vgld#1389](https://github.com/vigilnetwork/vgl/pull/1389))

### Code Contributors (alphabetical order):

- Dave Collins
- David Hill
- Dmitry Fedorov
- Donald Adu-Poku
- harzo
- hypernoob
- J Fixby
- Jonathan Chappelow
- Josh Rickmar
- Markus Richter
- matadormel
- Matheus Degiovani
- Michael Eze
- Orthomind
- Shuai Qi
- Tibor Bősze
- Victor Oliveira




