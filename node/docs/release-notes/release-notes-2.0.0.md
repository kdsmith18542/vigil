# vgld v2.0.0 Release Notes

This is a new major release of vgld.  Some of the key highlights are:

* Decentralized StakeShuffle mixing
* Higher network throughput
* Lightweight client sync time reduced by around 50%
* Improved initial peer discovery
* Reject protocol message removal
* Various updates to the RPC server:
  * Dynamic TLS certificate reload
  * Proof-of-Work hash in block information
  * New JSON-RPC API additions related to decentralized StakeShuffle mixing
* Quality assurance changes

## Upgrade Required To Continue Participating in StakeShuffle Mixing

Although upgrading to this latest release is not absolutely required for
continued operation of the core network, it is required for anyone who wishes to
continue participating in StakeShuffle mixing.

## Notable Changes

### Decentralized StakeShuffle Mixing

The StakeShuffle mixing process is now fully decentralized via the peer-to-peer
network as of this release.  All core software has been upgraded to make use of
the new decentralized coordination facilities.

This release introduces several new peer-to-peer protocol messages to provide
the decentralized coordination.  The following is a brief summary of the new
messages:

|Message      |Overall Purpose                                                |
|-------------|---------------------------------------------------------------|
|`mixpairreq` |Request to participate in a mix with relevant data and proofs. |
|`mixkeyxchg` |Publishes public keys and commitments for blame assignment.    |
|`mixcphrtxt` |Enables quantum resistant (PQ) blinded key exchange.           |
|`mixslotres` |Establishes slot reservations used in the blinding process.    |
|`mixfactpoly`|Encodes solution to the factored slot reservation polynomial.  |
|`mixdcnet`   |Untraceable multi-party broadcast (dining cryptographers).     |
|`mixconfirm` |Provides partial signatures to create the mix transaction.     |
|`mixsecrets` |Reveals secrets of an unsuccessful mix for blame assignment.   |

### Higher Network Throughput

This release now supports concurrent data requests (`getdata`) which allows for
higher network throughput, particularly when the communications channel is
experiencing high latency.

A couple of notable benefits are:

* Reduced vote times since it allows blocks and transactions to propagate more
  quickly throughout the network
* More responsive during traffic bursts and general network congestion

### Lightweight client sync time reduced by around 50%

Lightweight clients may now request version 2 compact filters in batches as
opposed to one at a time.  This has the effect of drastically reducing the
initial sync time for lightweight clients such as Simplified Payment
Verification (SPV) wallets.

This release introduces a new pair of peer-to-peer protocol messages named
`getcfsv2` and `cfiltersv2` which provide the aforementioned capability.

### Improved Initial Peer Discovery

Peers will now continue to query unreachable seeders in the background with an
increasing backoff timeout when they have not already discovered a sufficient
number of peers on the network to achieve the target connectivity.

This primarily improves the experience for peers joining the network for the
first time and those that have not been online for a long time since they do not
have a known list of good peers to use.
 
### Reject Protocol Message Removal (`reject`)

The previously deprecated `reject` peer-to-peer protocol message is no longer
available.

Consequently, the minimum required network protocol version to
participate on the network is now version 9.

Note that all nodes on older protocol versions are already not able to
participate in the network due to consensus changes that have passed.

Recall from previous release notes that this message is being removed because it
is a holdover from the original codebase where it was required, but it really is
not a useful message and has several downsides:

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

### RPC Server Changes

The RPC server version as of this release is 8.2.0.

#### Dynamic TLS Certificate Reload

The RPC server now checks for updates to the TLS certificate key pair
(`rpc.cert` and `rpc.key` by default) on new connections and dynamically reloads
them if needed.  Similarly, the authorized client certificates (`clients.pem` by
default) when running with the client certificate authorization type mode
(`--authtype=clientcert`).

Some key highlights of this change:

* Certificates can now be updated without needing to shutdown and restart the
  process which enables things such as:
  * Updating the certificates to change the allowed domain name and/or IP addresses
  * Dynamically adding or removing authorized clients
  * Changing the cryptographic primitives used to newer supported variants
* All existing connections will continue to use the certificates that were
  loaded at the time the connection was made
* The existing working certs are retained if any errors are encountered when
  loading the new ones in order to avoid breaking a working config

#### New Proof-of-Work Hash Field in Block Info RPCs (`getblock` and `getblockheader`)

The verbose results of the `getblock` and `getblockheader` RPCs now include a
`powhash` field in the JSON object that contains the Proof-of-Work hash for the
block.  The new field will be exactly the same as the `hash` (block hash) field
for all blocks prior to the activation of the stakeholder vote to change the PoW
hashing algorithm
([VGLP0011](https://github.com/Vigil/VGLPs/blob/master/VGLP-0011/VGLP-0011.mediawiki)).

See the following for API details:

* [getblock JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblock)
* [getblockheader JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#getblockheader)

#### New StakeShuffle Mixing Pool (mixpool) Message Send RPC (`sendrawmixmessage`)

A new RPC named `sendrawmixmessage` is now available.  This RPC can be used to
manually submit all mixing messages to the mixpool and broadcast them to the
network.

See the following for API details:

* [sendrawmixmessage JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#sendrawmixmessage)

#### New StakeShuffle Mixing Pool (mixpool) Message WebSocket Notification RPCs (`notifymixmessages`)

WebSocket notifications for mixing messages accepted to the mixpool are now
available.

Use `notifymixmessages` to request `mixmessage` notifications and
`stopnotifymixmessages` to stop receiving notifications.

See the following for API details:

* [notifymixmessages JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#notifymixmessages)
* [stopnotifymixmessages JSON-RPC API Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#stopnotifymixmessages)
* [mixmessage JSON-RPC API Notification Documentation](https://github.com/vigilnetwork/vgl/blob/master/docs/json_rpc_api.mediawiki#mixmessage)

## Changelog

This release consists of 168 commits from 11 contributors which total
to 265 files changed, 18292 additional lines of code, and 2978 deleted
lines of code.

All commits since the last release may be viewed on GitHub
[here](https://github.com/vigilnetwork/vgl/compare/release-v1.8.1...release-v2.0.0).

### Protocol and network:

- main: Add read header timeout to profile server ([Vigil/vgld#3186](https://github.com/vigilnetwork/vgl/pull/3186))
- server: Support concurrent getdata requests ([Vigil/vgld#3203](https://github.com/vigilnetwork/vgl/pull/3203))
- server: Update required minimum protocol version ([Vigil/vgld#3221](https://github.com/vigilnetwork/vgl/pull/3221))
- wire: add p2p mixing messages ([Vigil/vgld#3066](https://github.com/vigilnetwork/vgl/pull/3066))
- server: Retry seeder connections ([Vigil/vgld#3228](https://github.com/vigilnetwork/vgl/pull/3228))
- wire: Add epoch field to mix key exchange message ([Vigil/vgld#3235](https://github.com/vigilnetwork/vgl/pull/3235))
- wire: Remove PR hashes from (get)initstate ([Vigil/vgld#3244](https://github.com/vigilnetwork/vgl/pull/3244))
- wire: add previous revealed secrets hashes to RS message ([Vigil/vgld#3239](https://github.com/vigilnetwork/vgl/pull/3239))
- wire: Add Opcode field to MixPairReqUTXO ([Vigil/vgld#3243](https://github.com/vigilnetwork/vgl/pull/3243))
- wire: Zero out secrets signature for commitment hash ([Vigil/vgld#3248](https://github.com/vigilnetwork/vgl/pull/3248))
- wire: add flag bits to PR message ([Vigil/vgld#3246](https://github.com/vigilnetwork/vgl/pull/3246))
- wire: Add MsgMixFactoredPoly ([Vigil/vgld#3247](https://github.com/vigilnetwork/vgl/pull/3247))
- server: Require protocol v9 (removed reject msg) ([Vigil/vgld#3254](https://github.com/vigilnetwork/vgl/pull/3254))
- wire: Remove deprecated reject message support ([Vigil/vgld#3254](https://github.com/vigilnetwork/vgl/pull/3254))
- wire: Include unmixed session position in KE ([Vigil/vgld#3257](https://github.com/vigilnetwork/vgl/pull/3257))
- wire: Add msgs to get batch of cfilters ([Vigil/vgld#3211](https://github.com/vigilnetwork/vgl/pull/3211))
- multi: Respond to getcfsv2 message ([Vigil/vgld#3211](https://github.com/vigilnetwork/vgl/pull/3211))
- chaincfg: Update assume valid for release ([Vigil/vgld#3263](https://github.com/vigilnetwork/vgl/pull/3263))
- chaincfg: Update min known chain work for release ([Vigil/vgld#3264](https://github.com/vigilnetwork/vgl/pull/3264))
- multi: Integrate mixpool and propagate p2p mixing messages ([Vigil/vgld#3208](https://github.com/vigilnetwork/vgl/pull/3208))

### Mining:

- mining: Update blk templ diff for too few voters ([Vigil/vgld#3241](https://github.com/vigilnetwork/vgl/pull/3241))
- cpuminer: Rework discrete mining vote wait logic ([Vigil/vgld#3242](https://github.com/vigilnetwork/vgl/pull/3242))
- cpuminer: Remove unused IsKnownInvalidBlock method ([Vigil/vgld#3242](https://github.com/vigilnetwork/vgl/pull/3242))

### RPC:

- rpc: Add PoWHash to getblock/getblockheader (verbose) results ([Vigil/vgld#3154](https://github.com/vigilnetwork/vgl/pull/3154))
- rpcserver: Support dynamic cert reload ([Vigil/vgld#3153](https://github.com/vigilnetwork/vgl/pull/3153))
- rpcserver: Modify getnetworkhashps -1 blocks logic ([Vigil/vgld#3181](https://github.com/vigilnetwork/vgl/pull/3181))
- rpcserver: Remove unneeded AddedNodeInfo method ([Vigil/vgld#3236](https://github.com/vigilnetwork/vgl/pull/3236))

### vgld command-line flags and configuration:

- config: Support vgld_APPDATA env variable ([Vigil/vgld#3152](https://github.com/vigilnetwork/vgl/pull/3152))
- config: Show usage message on invalid cli flags ([Vigil/vgld#3282](https://github.com/vigilnetwork/vgl/pull/3282))

### Documentation:

- docs: Add release notes for v1.8.0 ([Vigil/vgld#3144](https://github.com/vigilnetwork/vgl/pull/3144))
- docs: Add release notes templates ([Vigil/vgld#3148](https://github.com/vigilnetwork/vgl/pull/3148))
- docs: Update for blockchain v5 module ([Vigil/vgld#3149](https://github.com/vigilnetwork/vgl/pull/3149))
- docs: Update JSON-RPC API for powhash ([Vigil/vgld#3154](https://github.com/vigilnetwork/vgl/pull/3154))
- docs: Update README.md to required Go 1.20/1.21 ([Vigil/vgld#3172](https://github.com/vigilnetwork/vgl/pull/3172))
- docs: Add release notes for v1.8.1 ([Vigil/vgld#3195](https://github.com/vigilnetwork/vgl/pull/3195))
- docs: Update README.md to required Go 1.21/1.22 ([Vigil/vgld#3220](https://github.com/vigilnetwork/vgl/pull/3220))
- docs: Add mixing cmds and ntfn JSON-RPC API ([Vigil/vgld#3275](https://github.com/vigilnetwork/vgl/pull/3275))

### Contrib changes:

- docker: Update image to golang:1.20.5-alpine3.18 ([Vigil/vgld#3146](https://github.com/vigilnetwork/vgl/pull/3146))
- docker: Update image to golang:1.20.6-alpine3.18 ([Vigil/vgld#3158](https://github.com/vigilnetwork/vgl/pull/3158))
- docker: Update image to golang:1.20.7-alpine3.18 ([Vigil/vgld#3170](https://github.com/vigilnetwork/vgl/pull/3170))
- docker: Update image to golang:1.21.0-alpine3.18 ([Vigil/vgld#3171](https://github.com/vigilnetwork/vgl/pull/3171))
- docker: Update image to golang:1.21.1-alpine3.18 ([Vigil/vgld#3183](https://github.com/vigilnetwork/vgl/pull/3183))
- docker: Update image to golang:1.21.2-alpine3.18 ([Vigil/vgld#3197](https://github.com/vigilnetwork/vgl/pull/3197))
- docker: Update image to golang:1.21.3-alpine3.18 ([Vigil/vgld#3198](https://github.com/vigilnetwork/vgl/pull/3198))
- docker: Update image to golang:1.21.4-alpine3.18 ([Vigil/vgld#3210](https://github.com/vigilnetwork/vgl/pull/3210))
- docker: Update image to golang:1.21.5-alpine3.18 ([Vigil/vgld#3214](https://github.com/vigilnetwork/vgl/pull/3214))
- docker: Update image to golang:1.21.6-alpine3.19 ([Vigil/vgld#3215](https://github.com/vigilnetwork/vgl/pull/3215))
- docker: Update image to golang:1.22.0-alpine3.19 ([Vigil/vgld#3219](https://github.com/vigilnetwork/vgl/pull/3219))
- docker: Update image to golang:1.22.1-alpine3.19 ([Vigil/vgld#3222](https://github.com/vigilnetwork/vgl/pull/3222))
- docker: Update image to golang:1.22.2-alpine3.19 ([Vigil/vgld#3231](https://github.com/vigilnetwork/vgl/pull/3231))
- docker: Update image to golang:1.22.3-alpine3.19 ([Vigil/vgld#3249](https://github.com/vigilnetwork/vgl/pull/3249))
- contrib: Add mixing to go workspace setup script ([Vigil/vgld#3265](https://github.com/vigilnetwork/vgl/pull/3265))

### Developer-related package and module changes:

- jsonrpc/types: Add powhash to verbose block output ([Vigil/vgld#3154](https://github.com/vigilnetwork/vgl/pull/3154))
- chaingen: More precise ASERT comments ([Vigil/vgld#3156](https://github.com/vigilnetwork/vgl/pull/3156))
- rpcclient: Explicitly require TLS >= 1.2 for HTTP ([Vigil/vgld#3169](https://github.com/vigilnetwork/vgl/pull/3169))
- multi: Avoid range capture for Go 1.22 changes ([Vigil/vgld#3165](https://github.com/vigilnetwork/vgl/pull/3165))
- blockchain: Remove unused progress logger param ([Vigil/vgld#3177](https://github.com/vigilnetwork/vgl/pull/3177))
- blockchain: Remove unused trsy enabled param ([Vigil/vgld#3177](https://github.com/vigilnetwork/vgl/pull/3177))
- multi: Wrap errors for better errors.Is/As support ([Vigil/vgld#3178](https://github.com/vigilnetwork/vgl/pull/3178))
- rpcserver: Improve internal error handling ([Vigil/vgld#3182](https://github.com/vigilnetwork/vgl/pull/3182))
- sampleconfig: Use embed with external files ([Vigil/vgld#3185](https://github.com/vigilnetwork/vgl/pull/3185))
- secp256k1/ecdsa: Expose r and s value of signature ([Vigil/vgld#3188](https://github.com/vigilnetwork/vgl/pull/3188))
- secp256k1/ecdsa: Remove some unnecessary subslices ([Vigil/vgld#3189](https://github.com/vigilnetwork/vgl/pull/3189))
- secp256k1/ecdsa: Add tests for new R and S methods ([Vigil/vgld#3190](https://github.com/vigilnetwork/vgl/pull/3190))
- secp256k1/ecdsa: Add test for order wraparound ([Vigil/vgld#3191](https://github.com/vigilnetwork/vgl/pull/3191))
- VGLutil: Use os.UserHomeDir in appDataDir ([Vigil/vgld#3196](https://github.com/vigilnetwork/vgl/pull/3196))
- multi: Reduce done goroutines ([Vigil/vgld#3199](https://github.com/vigilnetwork/vgl/pull/3199))
- multi: Consolidate waitgroup logic ([Vigil/vgld#3200](https://github.com/vigilnetwork/vgl/pull/3200))
- netsync: Rename NewPeer to PeerConnected ([Vigil/vgld#3202](https://github.com/vigilnetwork/vgl/pull/3202))
- netsync: Rename DonePeer to PeerDisconnected ([Vigil/vgld#3202](https://github.com/vigilnetwork/vgl/pull/3202))
- netsync: Export opaque peer and require it in API ([Vigil/vgld#3202](https://github.com/vigilnetwork/vgl/pull/3202))
- server: Use server peer in log statements ([Vigil/vgld#3202](https://github.com/vigilnetwork/vgl/pull/3202))
- server: Don't wait or try to send notfound data ([Vigil/vgld#3204](https://github.com/vigilnetwork/vgl/pull/3204))
- peer: select on p.quit during stall control chan writes ([Vigil/vgld#3209](https://github.com/vigilnetwork/vgl/pull/3209))
- peer: provide better debug for queued nil messages ([Vigil/vgld#3213](https://github.com/vigilnetwork/vgl/pull/3213))
- wire: Mark legacy message types as deprecated ([Vigil/vgld#3205](https://github.com/vigilnetwork/vgl/pull/3205))
- secp256k1: Add TinyGo support ([Vigil/vgld#3223](https://github.com/vigilnetwork/vgl/pull/3223))
- wire: Fix typo in comment ([Vigil/vgld#3226](https://github.com/vigilnetwork/vgl/pull/3226))
- secp256k1: No allocs in slow scalar base mult path ([Vigil/vgld#3225](https://github.com/vigilnetwork/vgl/pull/3225))
- VGLutil: Remove Getenv("HOME") fallback for appdata dir ([Vigil/vgld#3230](https://github.com/vigilnetwork/vgl/pull/3230))
- server: Do not update addrmgr on simnet/regnet ([Vigil/vgld#3240](https://github.com/vigilnetwork/vgl/pull/3240))
- connmgr: Only mark persistent peer reconn pending ([Vigil/vgld#3238](https://github.com/vigilnetwork/vgl/pull/3238))
- server: Use atomic types for some svr peer fields ([Vigil/vgld#3237](https://github.com/vigilnetwork/vgl/pull/3237))
- peer: Remove deprecated reject message support ([Vigil/vgld#3254](https://github.com/vigilnetwork/vgl/pull/3254))
- peer: Close mock connections in tests ([Vigil/vgld#3254](https://github.com/vigilnetwork/vgl/pull/3254))
- peer: Require protocol v9 (removed reject msg) ([Vigil/vgld#3254](https://github.com/vigilnetwork/vgl/pull/3254))
- gcs: Add func to determine max cfilter size ([Vigil/vgld#3211](https://github.com/vigilnetwork/vgl/pull/3211))
- blockchain: Add function to locate multiple cfilters ([Vigil/vgld#3211](https://github.com/vigilnetwork/vgl/pull/3211))
- server: Use sync.Mutex since the read lock isn't used ([Vigil/vgld#3270](https://github.com/vigilnetwork/vgl/pull/3270))
- mixing: Only validate compressed 33-byte pubkeys ([Vigil/vgld#3271](https://github.com/vigilnetwork/vgl/pull/3271))
- mixing: Add mixpool package ([Vigil/vgld#3082](https://github.com/vigilnetwork/vgl/pull/3082))
- mixing: Add mixclient package ([Vigil/vgld#3256](https://github.com/vigilnetwork/vgl/pull/3256))
- server: Implement separate mutex for peer state ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make server peer conn req concurrent safe ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Use iterator for connected ip count ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make add peer logic synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make done peer logic synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make conn count query synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make outbound group query synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make manual connect code synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make pending conn cancel code synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make persistent peer removal synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make persistent node query synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Make manual peer disconnect synchronous ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))
- server: Remove unused query chan and handler ([Vigil/vgld#3251](https://github.com/vigilnetwork/vgl/pull/3251))

### Developer-related module management:

- secp256k1: Prepare v4.3.0 ([Vigil/vgld#3227](https://github.com/vigilnetwork/vgl/pull/3227))
- main: update vgldtest module to master ([Vigil/vgld#3232](https://github.com/vigilnetwork/vgl/pull/3232))
- main: update vgldtest module to master ([Vigil/vgld#3233](https://github.com/vigilnetwork/vgl/pull/3233))
- wire: go mod tidy ([Vigil/vgld#3250](https://github.com/vigilnetwork/vgl/pull/3250))
- multi: Deduplicate external dependencies ([Vigil/vgld#3255](https://github.com/vigilnetwork/vgl/pull/3255))
- wire: Prepare v1.7.0 ([Vigil/vgld#3258](https://github.com/vigilnetwork/vgl/pull/3258))
- blockchain/standalone: Prepare v2.2.1 ([Vigil/vgld#3259](https://github.com/vigilnetwork/vgl/pull/3259))
- addrmgr: Prepare v2.0.3 ([Vigil/vgld#3260](https://github.com/vigilnetwork/vgl/pull/3260))
- mixing: Use latest crypto deps ([Vigil/vgld#3261](https://github.com/vigilnetwork/vgl/pull/3261))
- connmgr: Prepare v3.1.2 ([Vigil/vgld#3262](https://github.com/vigilnetwork/vgl/pull/3262))
- chaincfg: Prepare v3.2.1 ([Vigil/vgld#3266](https://github.com/vigilnetwork/vgl/pull/3266))
- txscript: Prepare v4.1.1 ([Vigil/vgld#3267](https://github.com/vigilnetwork/vgl/pull/3267))
- hdkeychain: Prepare v3.1.2 ([Vigil/vgld#3268](https://github.com/vigilnetwork/vgl/pull/3268))
- rpc/jsonrpc/types: Prepare v4.2.0 ([Vigil/vgld#3276](https://github.com/vigilnetwork/vgl/pull/3276))
- peer: Prepare v3.1.0 ([Vigil/vgld#3277](https://github.com/vigilnetwork/vgl/pull/3277))
- VGLutil: Prepare v4.0.2 ([Vigil/vgld#3278](https://github.com/vigilnetwork/vgl/pull/3278))
- database: Prepare v3.0.2 ([Vigil/vgld#3279](https://github.com/vigilnetwork/vgl/pull/3279))
- mixing: Prepare v0.1.0 ([Vigil/vgld#3280](https://github.com/vigilnetwork/vgl/pull/3280))
- blockchain/stake: Prepare v5.0.1 ([Vigil/vgld#3281](https://github.com/vigilnetwork/vgl/pull/3281))
- gcs: Prepare v4.1.0 ([Vigil/vgld#3283](https://github.com/vigilnetwork/vgl/pull/3283))
- blockchain: Prepare v5.0.1 ([Vigil/vgld#3284](https://github.com/vigilnetwork/vgl/pull/3284))
- rpcclient: Prepare v8.0.1 ([Vigil/vgld#3285](https://github.com/vigilnetwork/vgl/pull/3285))
- main: Update to use all new module versions ([Vigil/vgld#3286](https://github.com/vigilnetwork/vgl/pull/3286))
- main: Remove module replacements ([Vigil/vgld#3288](https://github.com/vigilnetwork/vgl/pull/3288))
- mixing: Introduce module ([Vigil/vgld#3207](https://github.com/vigilnetwork/vgl/pull/3207))

### Testing and Quality Assurance:

- main: Use release version of VGLtest framework ([Vigil/vgld#3142](https://github.com/vigilnetwork/vgl/pull/3142))
- database: Mark test helpers ([Vigil/vgld#3147](https://github.com/vigilnetwork/vgl/pull/3147))
- database: Use TempDir to create temp test dirs ([Vigil/vgld#3147](https://github.com/vigilnetwork/vgl/pull/3147))
- build: Add CI support for test and module cache ([Vigil/vgld#3145](https://github.com/vigilnetwork/vgl/pull/3145))
- main: improve test flag handling ([Vigil/vgld#3151](https://github.com/vigilnetwork/vgl/pull/3151))
- build: Add nilerr linter ([Vigil/vgld#3157](https://github.com/vigilnetwork/vgl/pull/3157))
- build: Update to latest action versions ([Vigil/vgld#3159](https://github.com/vigilnetwork/vgl/pull/3159))
- build: Move lint logic to its own script ([Vigil/vgld#3161](https://github.com/vigilnetwork/vgl/pull/3161))
- build: Use go install for linter and add cache ([Vigil/vgld#3162](https://github.com/vigilnetwork/vgl/pull/3162))
- build: Update golangci-lint to v1.53.3 ([Vigil/vgld#3163](https://github.com/vigilnetwork/vgl/pull/3163))
- build: Correct missing shebang in lint script ([Vigil/vgld#3164](https://github.com/vigilnetwork/vgl/pull/3164))
- build: Checkout source before Go setup ([Vigil/vgld#3166](https://github.com/vigilnetwork/vgl/pull/3166))
- build: Use setup-go action cache ([Vigil/vgld#3168](https://github.com/vigilnetwork/vgl/pull/3168))
- build: Update to latest action versions ([Vigil/vgld#3172](https://github.com/vigilnetwork/vgl/pull/3172))
- build: Update golangci-lint to v1.53.1 ([Vigil/vgld#3172](https://github.com/vigilnetwork/vgl/pull/3172))
- build: Test against Go 1.21 ([Vigil/vgld#3172](https://github.com/vigilnetwork/vgl/pull/3172))
- build: Test against Go 1.21 ([Vigil/vgld#3172](https://github.com/vigilnetwork/vgl/pull/3172))
- standalone: Add decreasing timestamps ASERT test ([Vigil/vgld#3173](https://github.com/vigilnetwork/vgl/pull/3173))
- build: Add dupword linter ([Vigil/vgld#3175](https://github.com/vigilnetwork/vgl/pull/3175))
- build: Add errorlint linter ([Vigil/vgld#3179](https://github.com/vigilnetwork/vgl/pull/3179))
- build: Update to latest action versions ([Vigil/vgld#3216](https://github.com/vigilnetwork/vgl/pull/3216))
- build: Update golangci-lint to v1.55.2 ([Vigil/vgld#3216](https://github.com/vigilnetwork/vgl/pull/3216))
- build: Update golangci-lint to v1.56.0 ([Vigil/vgld#3220](https://github.com/vigilnetwork/vgl/pull/3220))
- build: Test against Go 1.22 ([Vigil/vgld#3220](https://github.com/vigilnetwork/vgl/pull/3220))
- secp256k1: Add scalar base mult variant benchmarks ([Vigil/vgld#3224](https://github.com/vigilnetwork/vgl/pull/3224))
- run_vgl_tests.sh: allow passing of additional test arguments ([Vigil/vgld#3229](https://github.com/vigilnetwork/vgl/pull/3229))
- VGLutil: Fix message in test error ([Vigil/vgld#3230](https://github.com/vigilnetwork/vgl/pull/3230))
- VGLutil: Use os.UserHomedir for base home directory in tests ([Vigil/vgld#3230](https://github.com/vigilnetwork/vgl/pull/3230))
- rpctests: Pass test loggers to vgldtest package ([Vigil/vgld#3232](https://github.com/vigilnetwork/vgl/pull/3232))
- multi: Fix function names in some comments ([Vigil/vgld#3245](https://github.com/vigilnetwork/vgl/pull/3245))
- rpc/jsonrpc/types: Add tests for new mix types ([Vigil/vgld#3274](https://github.com/vigilnetwork/vgl/pull/3274))

### Misc:

- release: Bump for 1.9 release cycle ([Vigil/vgld#3141](https://github.com/vigilnetwork/vgl/pull/3141))
- main: Don't include requires in build ([Vigil/vgld#3143](https://github.com/vigilnetwork/vgl/pull/3143))
- multi: Address some linter complaints ([Vigil/vgld#3155](https://github.com/vigilnetwork/vgl/pull/3155))
- multi: Remove a bunch of dup words in comments ([Vigil/vgld#3174](https://github.com/vigilnetwork/vgl/pull/3174))
- multi: Cleanup superfluous trailing newlines ([Vigil/vgld#3176](https://github.com/vigilnetwork/vgl/pull/3176))
- main: Update license to 2024 ([Vigil/vgld#3217](https://github.com/vigilnetwork/vgl/pull/3217))
- release: Bump for 2.0 release cycle ([Vigil/vgld#3269](https://github.com/vigilnetwork/vgl/pull/3269))
- release: Bump for 2.0.0 ([Vigil/vgld#3289](https://github.com/vigilnetwork/vgl/pull/3289))

### Code Contributors (alphabetical order):

- Billy Zelani Malik
- Dave Collins
- David Hill
- Matheus Degiovani
- Nicola Larosa
- peicuiping
- Peter Zen
- Jamie Holdstock
- Jonathan Chappelow
- Josh Rickmark
- SeedHammer
