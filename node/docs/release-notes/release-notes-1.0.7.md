## vgld v1.0.7

This release of vgld primarily contains improvements to the infrastructure and
other quality assurance changes that are bringing us closer to providing full
support for Lightning Network.

A lot of work required for Lightning Network support went into getting the
required code merged into the upstream project, btcd, which now fully supports
it.  These changes also must be synced and integrated with vgld as well and
therefore many of the changes in this release are related to that process.

## Notable Changes

### Dust check removed from stake transactions

The standard policy for regular transactions is to reject any transactions that
have outputs so small that they cost more to the network than their value.  This
behavior is desirable for regular transactions, however it was also being
applied to vote and revocation transactions which could lead to a situation
where stake pools with low fees could result in votes and revocations having
difficulty being mined.

This check has been changed to only apply to regular transactions now in order
to prevent any issues.  Stake transactions have several other checks that make
this one unnecessary for them.

### New `feefilter` peer-to-peer message

A new optional peer-to-peer message named `feefilter` has been added that allows
peers to inform others about the minimum transaction fee rate they are willing
to accept.  This will enable peers to avoid notifying others about transactions
they will not accept anyways and therefore can result in a significant bandwidth
savings.

### Bloom filter service bit enforcement

Peers that are configured to disable bloom filter support will now disconnect
remote peers that send bloom filter related commands rather than simply ignoring
them.  This allows any light clients that do not observe the service bit to
potentially find another peer that provides the service.  Additionally, remote
peers that have negotiated a high enough protocol version to observe the service
bit and still send bloom filter related commands anyways will now be banned.


## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/v1.0.5...v1.0.7).

### Protocol and network:
- Allow reorg of block one [Vigil/vgld#745](https://github.com/vigilnetwork/vgl/pull/745)
- blockchain: use the time source [Vigil/vgld#747](https://github.com/vigilnetwork/vgl/pull/747)
- peer: Strictly enforce bloom filter service bit [Vigil/vgld#768](https://github.com/vigilnetwork/vgl/pull/768)
- wire/peer: Implement feefilter p2p message [Vigil/vgld#779](https://github.com/vigilnetwork/vgl/pull/779)
- chaincfg: update checkpoints for 1.0.7 release  [Vigil/vgld#816](https://github.com/vigilnetwork/vgl/pull/816)

### Transaction relay (memory pool):
- mempool: Break dependency on chain instance [Vigil/vgld#754](https://github.com/vigilnetwork/vgl/pull/754)
- mempool: unexport the mutex [Vigil/vgld#755](https://github.com/vigilnetwork/vgl/pull/755)
- mempool: Add basic test harness infrastructure [Vigil/vgld#756](https://github.com/vigilnetwork/vgl/pull/756)
- mempool: Improve tx input standard checks [Vigil/vgld#758](https://github.com/vigilnetwork/vgl/pull/758)
- mempool: Update comments for dust calcs [Vigil/vgld#764](https://github.com/vigilnetwork/vgl/pull/764)
- mempool: Only perform standard dust checks on regular transactions  [Vigil/vgld#806](https://github.com/vigilnetwork/vgl/pull/806)

### RPC:
- Fix gettxout includemempool handling [Vigil/vgld#738](https://github.com/vigilnetwork/vgl/pull/738)
- Improve help text for getmininginfo [Vigil/vgld#748](https://github.com/vigilnetwork/vgl/pull/748)
- rpcserverhelp: update TicketFeeInfo help [Vigil/vgld#801](https://github.com/vigilnetwork/vgl/pull/801)
- blockchain: Improve getstakeversions efficiency [Vigil/vgld#81](https://github.com/vigilnetwork/vgl/pull/813)

### vgld command-line flags:
- config: introduce new flags to accept/reject non-std transactions [Vigil/vgld#757](https://github.com/vigilnetwork/vgl/pull/757)
- config: Add --whitelist option [Vigil/vgld#352](https://github.com/vigilnetwork/vgl/pull/352)
- config: Improve config file handling [Vigil/vgld#802](https://github.com/vigilnetwork/vgl/pull/802)
- config: Improve blockmaxsize check [Vigil/vgld#810](https://github.com/vigilnetwork/vgl/pull/810)

### vglctl:
- Add --walletrpcserver option [Vigil/vgld#736](https://github.com/vigilnetwork/vgl/pull/736)

### Documentation
- docs: add commit prefix notes  [Vigil/vgld#788](https://github.com/vigilnetwork/vgl/pull/788)

### Developer-related package changes:
- blockchain: check errors and remove ineffectual assignments [Vigil/vgld#689](https://github.com/vigilnetwork/vgl/pull/689)
- stake: less casting [Vigil/vgld#705](https://github.com/vigilnetwork/vgl/pull/705)
- blockchain: chainstate only needs prev block hash [Vigil/vgld#706](https://github.com/vigilnetwork/vgl/pull/706)
- remove dead code [Vigil/vgld#715](https://github.com/vigilnetwork/vgl/pull/715)
- Use btclog for determining valid log levels [Vigil/vgld#738](https://github.com/vigilnetwork/vgl/pull/738)
- indexers: Minimize differences with upstream code [Vigil/vgld#742](https://github.com/vigilnetwork/vgl/pull/742)
- blockchain: Add median time to state snapshot [Vigil/vgld#753](https://github.com/vigilnetwork/vgl/pull/753)
- blockmanager: remove unused GetBlockFromHash function [Vigil/vgld#761](https://github.com/vigilnetwork/vgl/pull/761)
- mining: call CheckConnectBlock directly [Vigil/vgld#762](https://github.com/vigilnetwork/vgl/pull/762)
- blockchain: add missing error code entries [Vigil/vgld#763](https://github.com/vigilnetwork/vgl/pull/763)
- blockchain: Sync main chain flag on ProcessBlock [Vigil/vgld#767](https://github.com/vigilnetwork/vgl/pull/767)
- blockchain: Remove exported CalcPastTimeMedian func [Vigil/vgld#770](https://github.com/vigilnetwork/vgl/pull/770)
- blockchain: check for error [Vigil/vgld#772](https://github.com/vigilnetwork/vgl/pull/772)
- multi: Optimize by removing defers [Vigil/vgld#782](https://github.com/vigilnetwork/vgl/pull/782)
- blockmanager: remove unused logBlockHeight [Vigil/vgld#787](https://github.com/vigilnetwork/vgl/pull/787)
- VGLutil: Replace DecodeNetworkAddress with DecodeAddress [Vigil/vgld#746](https://github.com/vigilnetwork/vgl/pull/746)
- txscript: Force extracted addrs to compressed [Vigil/vgld#775](https://github.com/vigilnetwork/vgl/pull/775)
- wire: Remove legacy transaction decoding [Vigil/vgld#794](https://github.com/vigilnetwork/vgl/pull/794)
- wire: Remove dead legacy tx decoding code [Vigil/vgld#796](https://github.com/vigilnetwork/vgl/pull/796)
- mempool/wire: Don't make policy decisions in wire [Vigil/vgld#797](https://github.com/vigilnetwork/vgl/pull/797)
- VGLjson: Remove unused cmds & types [Vigil/vgld#795](https://github.com/vigilnetwork/vgl/pull/795)
- VGLjson: move cmd types [Vigil/vgld#799](https://github.com/vigilnetwork/vgl/pull/799)
- multi: Separate tx serialization type from version [Vigil/vgld#798](https://github.com/vigilnetwork/vgl/pull/798)
- VGLjson: add Unconfirmed field to VGLjson.GetAccountBalanceResult [Vigil/vgld#812](https://github.com/vigilnetwork/vgl/pull/812)
- multi: Error descriptions should be lowercase [Vigil/vgld#771](https://github.com/vigilnetwork/vgl/pull/771)
- blockchain: cast to int64  [Vigil/vgld#817](https://github.com/vigilnetwork/vgl/pull/817)

### Testing and Quality Assurance:
- rpcserver: Upstream sync to add basic RPC tests [Vigil/vgld#750](https://github.com/vigilnetwork/vgl/pull/750)
- rpctest: Correct several issues tests and joins [Vigil/vgld#751](https://github.com/vigilnetwork/vgl/pull/751)
- rpctest: prevent process leak due to panics [Vigil/vgld#752](https://github.com/vigilnetwork/vgl/pull/752)
- rpctest: Cleanup resources on failed setup [Vigil/vgld#759](https://github.com/vigilnetwork/vgl/pull/759)
- rpctest: Use ports based on the process id [Vigil/vgld#760](https://github.com/vigilnetwork/vgl/pull/760)
- rpctest/deps: Update dependencies and API [Vigil/vgld#765](https://github.com/vigilnetwork/vgl/pull/765)
- rpctest: Gate rpctest-based behind a build tag [Vigil/vgld#766](https://github.com/vigilnetwork/vgl/pull/766)
- mempool: Add test for max orphan entry eviction [Vigil/vgld#769](https://github.com/vigilnetwork/vgl/pull/769)
- fullblocktests: Add more consensus tests [Vigil/vgld#77](https://github.com/vigilnetwork/vgl/pull/773)
- fullblocktests: Sync upstream block validation [Vigil/vgld#774](https://github.com/vigilnetwork/vgl/pull/774)
- rpctest: fix a harness range bug in syncMempools [Vigil/vgld#778](https://github.com/vigilnetwork/vgl/pull/778)
- secp256k1: Add regression tests for field.go [Vigil/vgld#781](https://github.com/vigilnetwork/vgl/pull/781)
- secp256k1: Sync upstream test consolidation [Vigil/vgld#783](https://github.com/vigilnetwork/vgl/pull/783)
- txscript: Correct p2sh hashes in json test data  [Vigil/vgld#785](https://github.com/vigilnetwork/vgl/pull/785)
- txscript: Replace CODESEPARATOR json test data [Vigil/vgld#786](https://github.com/vigilnetwork/vgl/pull/786)
- txscript: Remove multisigdummy from json test data [Vigil/vgld#789](https://github.com/vigilnetwork/vgl/pull/789)
- txscript: Remove max money from json test data [Vigil/vgld#790](https://github.com/vigilnetwork/vgl/pull/790)
- txscript: Update signatures in json test data [Vigil/vgld#791](https://github.com/vigilnetwork/vgl/pull/791)
- txscript: Use native encoding in json test data [Vigil/vgld#792](https://github.com/vigilnetwork/vgl/pull/792)
- rpctest: Store logs and data in same path [Vigil/vgld#780](https://github.com/vigilnetwork/vgl/pull/780)
- txscript: Cleanup reference test code  [Vigil/vgld#793](https://github.com/vigilnetwork/vgl/pull/793)

### Misc:
- Update deps to pull in additional logging changes [Vigil/vgld#734](https://github.com/vigilnetwork/vgl/pull/734)
- Update markdown files for GFM changes [Vigil/vgld#744](https://github.com/vigilnetwork/vgl/pull/744)
- blocklogger: Show votes, tickets, & revocations [Vigil/vgld#784](https://github.com/vigilnetwork/vgl/pull/784)
- blocklogger: Remove STransactions from transactions calculation [Vigil/vgld#811](https://github.com/vigilnetwork/vgl/pull/811)

### Contributors (alphabetical order):

- Alex Yocomm-Piatt
- Atri Viss
- Chris Martin
- Dave Collins
- David Hill
- Donald Adu-Poku
- Jimmy Song
- John C. Vernaleo
- Jolan Luff
- Josh Rickmar
- Olaoluwa Osuntokun
- Marco Peereboom




