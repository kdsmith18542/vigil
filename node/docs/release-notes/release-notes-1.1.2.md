# vgld v1.1.2

This release of vgld primarily contains performance enhancements, infrastructure
improvements, and other quality assurance changes.

While it is not visible in this release, significant infrastructure work has
also been done this release cycle towards porting the Lightning Network (LN)
daemon which will ultimately allow LN payments to be backed by Vigil.

## Notable Changes

### Faster Block Validation

A significant portion of block validation involves handling the stake tickets
which form an integral part of Vigil's hybrid proof-of-work and proof-of-stake
system.  The code which handles this portion of validation has been
significantly optimized in this release such that overall block validation is
up to approximately 3 times faster depending on the specific underlying hardware
configuration.  This also has a noticeable impact on the speed of the initial
block download process as well as how quickly votes for winning tickets are
submitted to the network.

### Data Carrier Transaction Standardness Policy

The standard policy for transaction relay of data carrier transaction outputs
has been modified to support canonically-encoded small data pushes.  These
outputs are also known as `OP_RETURN` or `nulldata` outputs.  In particular,
single byte small integers data pushes (0-16) are now supported.

## Changelog

All commits since the last release may be viewed on GitHub [here](https://github.com/vigilnetwork/vgl/compare/v1.1.0...v1.1.2).

### Protocol and network:
- chaincfg: update checkpoints for 1.1.2 release [Vigil/vgld#946](https://github.com/vigilnetwork/vgl/pull/946)
- chaincfg: Rename one of the testnet seeders [Vigil/vgld#873](https://github.com/vigilnetwork/vgl/pull/873)
- stake: treap index perf improvement [Vigil/vgld#853](https://github.com/vigilnetwork/vgl/pull/853)
- stake: ticket expiry perf improvement [Vigil/vgld#853](https://github.com/vigilnetwork/vgl/pull/853)

### Transaction relay (memory pool):

- txscript: Correct nulldata standardness check [Vigil/vgld#935](https://github.com/vigilnetwork/vgl/pull/935)

### RPC:

- rpcserver: searchrawtransactions skip first input for vote tx [Vigil/vgld#859](https://github.com/vigilnetwork/vgl/pull/859)
- multi: update stakebase tx vin[0] structure [Vigil/vgld#859](https://github.com/vigilnetwork/vgl/pull/859)
- rpcserver: Fix empty ssgen verbose results [Vigil/vgld#871](https://github.com/vigilnetwork/vgl/pull/871)
- rpcserver: check for error in getwork request [Vigil/vgld#898](https://github.com/vigilnetwork/vgl/pull/898)
- multi: Add NoSplitTransaction to purchaseticket [Vigil/vgld#904](https://github.com/vigilnetwork/vgl/pull/904)
- rpcserver: avoid nested decodescript p2sh addrs [Vigil/vgld#929](https://github.com/vigilnetwork/vgl/pull/929)
- rpcserver: skip generating certs when nolisten set [Vigil/vgld#932](https://github.com/vigilnetwork/vgl/pull/932)
- rpc: Add localaddr and relaytxes to getpeerinfo [Vigil/vgld#933](https://github.com/vigilnetwork/vgl/pull/933)
- rpcserver: update handleSendRawTransaction error handling [Vigil/vgld#939](https://github.com/vigilnetwork/vgl/pull/939)

### vgld command-line flags:

- config: add --nofilelogging option [Vigil/vgld#872](https://github.com/vigilnetwork/vgl/pull/872)

### Documentation:

- rpcclient: Remove docker info from README.md [Vigil/vgld#886](https://github.com/vigilnetwork/vgl/pull/886)
- bloom: Fix link in README [Vigil/vgld#922](https://github.com/vigilnetwork/vgl/pull/922)
- doc: tiny fix url [Vigil/vgld#928](https://github.com/vigilnetwork/vgl/pull/928)
- doc: update go version for example test run in readme [Vigil/vgld#936](https://github.com/vigilnetwork/vgl/pull/936)

### Developer-related package changes:

- multi: Drop glide, use dep [Vigil/vgld#818](https://github.com/vigilnetwork/vgl/pull/818)
- txsort: Implement stable tx sorting package  [Vigil/vgld#940](https://github.com/vigilnetwork/vgl/pull/940)
- coinset: Remove package [Vigil/vgld#888](https://github.com/vigilnetwork/vgl/pull/888)
- base58: Use new github.com/Vigil/base58 package [Vigil/vgld#888](https://github.com/vigilnetwork/vgl/pull/888)
- certgen: Move self signed certificate code into package [Vigil/vgld#879](https://github.com/vigilnetwork/vgl/pull/879)
- certgen: Add doc.go and README.md [Vigil/vgld#883](https://github.com/vigilnetwork/vgl/pull/883)
- rpcclient: Allow request-scoped cancellation during Connect [Vigil/vgld#880](https://github.com/vigilnetwork/vgl/pull/880)
- rpcclient: Import VGLrpcclient repo into rpcclient directory [Vigil/vgld#880](https://github.com/vigilnetwork/vgl/pull/880)
- rpcclient: json unmarshal into unexported embedded pointer  [Vigil/vgld#941](https://github.com/vigilnetwork/vgl/pull/941)
- bloom: Copy github.com/Vigil/VGLutil/bloom to bloom package [Vigil/vgld#881](https://github.com/vigilnetwork/vgl/pull/881)
- Improve gitignore [Vigil/vgld#887](https://github.com/vigilnetwork/vgl/pull/887)
- VGLutil: Import VGLutil repo under VGLutil directory [Vigil/vgld#888](https://github.com/vigilnetwork/vgl/pull/888)
- hdkeychain: Move to github.com/vigilnetwork/vgl/hdkeychain [Vigil/vgld#892](https://github.com/vigilnetwork/vgl/pull/892)
- stake: Add IsStakeSubmission [Vigil/vgld#907](https://github.com/vigilnetwork/vgl/pull/907)
- txscript: Require SHA256 secret hashes for atomic swaps [Vigil/vgld#930](https://github.com/vigilnetwork/vgl/pull/930)

### Testing and Quality Assurance:

- gometalinter: run on subpkgs too [Vigil/vgld#878](https://github.com/vigilnetwork/vgl/pull/878)
- travis: test Gopkg.lock [Vigil/vgld#889](https://github.com/vigilnetwork/vgl/pull/889)
- hdkeychain: Work around go vet issue with examples [Vigil/vgld#890](https://github.com/vigilnetwork/vgl/pull/890)
- bloom: Add missing import to examples [Vigil/vgld#891](https://github.com/vigilnetwork/vgl/pull/891)
- bloom: workaround go vet issue in example [Vigil/vgld#895](https://github.com/vigilnetwork/vgl/pull/895)
- tests: make lockfile test work locally [Vigil/vgld#894](https://github.com/vigilnetwork/vgl/pull/894)
- peer: Avoid goroutine leaking during handshake timeout [Vigil/vgld#909](https://github.com/vigilnetwork/vgl/pull/909)
- travis: add gosimple linter [Vigil/vgld#897](https://github.com/vigilnetwork/vgl/pull/897)
- multi: Handle detected data race conditions [Vigil/vgld#920](https://github.com/vigilnetwork/vgl/pull/920)
- travis: add ineffassign linter [Vigil/vgld#896](https://github.com/vigilnetwork/vgl/pull/896)
- rpctest: Choose flags based on provided params [Vigil/vgld#937](https://github.com/vigilnetwork/vgl/pull/937)

### Misc:

- gofmt [Vigil/vgld#876](https://github.com/vigilnetwork/vgl/pull/876)
- dep: sync third-party deps [Vigil/vgld#877](https://github.com/vigilnetwork/vgl/pull/877)
- Bump for v1.1.2 [Vigil/vgld#916](https://github.com/vigilnetwork/vgl/pull/916)
- dep: Use upstream jrick/bitset [Vigil/vgld#899](https://github.com/vigilnetwork/vgl/pull/899)
- blockchain: removed unused funcs and vars [Vigil/vgld#900](https://github.com/vigilnetwork/vgl/pull/900)
- blockchain: remove unused file [Vigil/vgld#900](https://github.com/vigilnetwork/vgl/pull/900)
- rpcserver: nil pointer dereference when submit orphan block [Vigil/vgld#906](https://github.com/vigilnetwork/vgl/pull/906)
- multi: remove unused funcs and vars [Vigil/vgld#901](https://github.com/vigilnetwork/vgl/pull/901)

### Code Contributors (alphabetical order):

- Alex Yocom-Piatt
- Dave Collins
- David Hill
- detailyang
- Donald Adu-Poku
- Federico Gimenez
- Jason Zavaglia
- John C. Vernaleo
- Jonathan Chappelow
- Jolan Luff
- Josh Rickmar
- Maninder Lall
- Matheus Degiovani
- Nicola Larosa
- Samarth Hattangady
- Ugwueze Onyekachi Michael




