vgld
====

[![Build Status](https://github.com/vigilnetwork/vgl/workflows/Build%20and%20Test/badge.svg)](https://github.com/vigilnetwork/vgl/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/vigilnetwork/vgl)
[![Go Report Card](https://goreportcard.com/badge/github.com/vigilnetwork/vgl)](https://goreportcard.com/report/github.com/vigilnetwork/vgl)

## Vigil Overview

Vigil is a blockchain-based cryptocurrency with a strong focus on community
input, open governance, and sustainable funding for development. It utilizes a
hybrid proof-of-work and proof-of-stake mining system to ensure that a small
group cannot dominate the flow of transactions or make changes to Vigil without
the input of the community.  A unit of the currency is called a `Vigil` (VGL).

https://vigil.network

## Latest Downloads

https://vigil.network/downloads/

Core software:

* vgld: a Vigil full node daemon (this)
* [vglwallet](https://github.com/Vigil/vglwallet): a CLI Vigil wallet daemon
* [vglctl](https://github.com/Vigil/vglctl): a CLI client for vgld and vglwallet

Bundles:

* [Vigiliton](https://github.com/Vigil/Vigiliton): a GUI bundle for `vgld`
  and `vglwallet`
* [CLI app suite](https://github.com/Vigil/Vigil-release/releases/latest):
  a CLI bundle for `vgld` and `vglwallet`

## What is vgld?

vgld is a full node implementation of Vigil written in Go (golang).

It acts as a fully-validating chain daemon for the Vigil cryptocurrency.  vgld
maintains the entire past transactional ledger of Vigil and allows relaying of
transactions to other Vigil nodes around the world.

This software is currently under active development.  It is extremely stable and
has been in production use since February 2016.

It important to note that vgld does *NOT* include wallet functionality.  Users
who desire a wallet will need to use [vglwallet(CLI)](https://github.com/Vigil/vglwallet)
or [Vigiliton(GUI)](https://github.com/Vigil/Vigiliton).

## What is a full node?

The term 'full node' is short for 'fully-validating node' and refers to software
that fully validates all transactions and blocks, as opposed to trusting a 3rd
party.  In addition to validating transactions and blocks, nearly all full nodes
also participate in relaying transactions and blocks to other full nodes around
the world, thus forming the peer-to-peer network that is the backbone of the
Vigil cryptocurrency.

The full node distinction is important, since full nodes are not the only type
of software participating in the Vigil peer network. For instance, there are
'lightweight nodes' which rely on full nodes to serve the transactions, blocks,
and cryptographic proofs they require to function, as well as relay their
transactions to the rest of the global network.

## Why run vgld?

As described in the previous section, the Vigil cryptocurrency relies on having
a peer-to-peer network of nodes that fully validate all transactions and blocks
and then relay them to other full nodes.

Running a full node with vgld contributes to the overall security of the
network, increases the available paths for transactions and blocks to relay,
and helps ensure there are an adequate number of nodes available to serve
lightweight clients, such as Simplified Payment Verification (SPV) wallets.

Without enough full nodes, the network could be unable to expediently serve
users of lightweight clients which could force them to have to rely on
centralized services that significantly reduce privacy and are vulnerable to
censorship.

In terms of individual benefits, since vgld fully validates every block and
transaction, it provides the highest security and privacy possible when used in
conjunction with a wallet that also supports directly connecting to it in full
validation mode, such as [vglwallet (CLI)](https://github.com/Vigil/vglwallet)
and [Vigiliton (GUI)](https://github.com/Vigil/Vigiliton).  It is also ideal
for businesses and services that need the most reliable and accurate data about
transactions.

## Minimum Recommended Specifications (vgld only)

* 16 GB disk space (as of April 2022, increases over time, ~2 GB/yr)
* 2 GB memory (RAM)
* ~150 MB/day download, ~1.5 GB/day upload
  * Plus one-time initial download of the entire block chain
* Windows 10 (server preferred), macOS, Linux
* High uptime

## Getting Started

So, you've decided to help the network by running a full node.  Great!  Running
vgld is simple.  All you need to do is install vgld on a machine that is
connected to the internet and meets the minimum recommended specifications, and
launch it.

Also, make sure your firewall is configured to allow inbound connections to port
9108.

<a name="Installation" />

## Installing and updating

### Binaries (Windows/Linux/macOS)

Binary releases are provided for common operating systems and architectures.
The easiest method is to download Vigiliton from the link below, which will
include vgld. Advanced users may prefer the Command-line app suite, which
includes vgld and vglwallet.

https://vigil.network/downloads/

* How to verify binaries before installing: https://docs.vigil.network/advanced/verifying-binaries/
* How to install the CLI Suite: https://docs.vigil.network/wallets/cli/cli-installation/
* How to install Vigiliton: https://docs.vigil.network/wallets/Vigiliton/Vigiliton-setup/

### Build from source (all platforms)

<details><summary><b>Install Dependencies</b></summary>

- **Go 1.23 or 1.24**

  Installation instructions can be found here: https://golang.org/doc/install.
  Ensure Go was installed properly and is a supported version:
  ```sh
  $ go version
  $ go env GOROOT GOPATH
  ```
  NOTE: `GOROOT` and `GOPATH` must not be on the same path. Since Go 1.8 (2016),
  `GOROOT` and `GOPATH` are set automatically, and you do not need to change
  them. However, you still need to add `$GOPATH/bin` to your `PATH` in order to
  run binaries installed by `go get` and `go install` (On Windows, this happens
  automatically).

  Unix example -- add these lines to .profile:

  ```
  PATH="$PATH:/usr/local/go/bin"  # main Go binaries ($GOROOT/bin)
  PATH="$PATH:$HOME/go/bin"       # installed Go projects ($GOPATH/bin)
  ```

- **Git**

  Installation instructions can be found at https://git-scm.com or
  https://gitforwindows.org.
  ```sh
  $ git version
  ```
</details>
<details><summary><b>Windows Example</b></summary>

  ```PowerShell
  PS> git clone https://github.com/vigilnetwork/vgl $env:USERPROFILE\src\vgld
  PS> cd $env:USERPROFILE\src\vgld
  PS> go install . .\cmd\...
  PS> vgld -V
  ```

  Run the `vgld` executable now installed in `"$(go env GOPATH)\bin"`.
</details>
<details><summary><b>Unix Example</b></summary>

  This assumes you have already added `$GOPATH/bin` to your `$PATH` as described
  in dependencies.

  ```sh
  $ git clone https://github.com/vigilnetwork/vgl $HOME/src/vgld
  $ git clone https://github.com/Vigil/vglctl $HOME/src/vglctl
  $ (cd $HOME/src/vgld && go install . ./...)
  $ (cd $HOME/src/vglctl && go install)
  $ vgld -V
  ```

  Run the `vgld` executable now installed in `$GOPATH/bin`.
</details>

## Building and Running OCI Containers (aka Docker/Podman)

The project does not officially provide container images.  However, all of the
necessary files to build your own lightweight non-root container image based on
`scratch` from the latest source code are available in
[contrib/docker](./contrib/docker/README.md).

It is also worth noting that, to date, most users typically prefer to run `vgld`
directly, without using a container, for at least a few reasons:

- `vgld` is a static binary that does not require root privileges and therefore
  does not suffer from the usual deployment issues that typically make
  containers attractive
- It is harder and more verbose to run `vgld` from a container as compared to
  normal:
  - `vgld` is designed to automatically create a working default configuration
    which means it just works out of the box without the need for additional
    configuration for almost all typical users
  - The blockchain data and configuration files need to be persistent which
    means configuring and managing a docker data volume
  - Running non-root containers with `docker` requires special care in regards
    to permissions

## Running Tests

All tests and linters may be run using the script `run_vgl_tests.sh`.  Generally,
Vigil only supports the current and previous major versions of Go.

```
./run_vgl_tests.sh
```

## Contact

If you have any further questions you can find us at:

https://vigil.network/community/

## Issue Tracker

The [integrated github issue tracker](https://github.com/vigilnetwork/vgl/issues)
is used for this project.

## Documentation

The documentation for vgld is a work-in-progress.  It is located in the
[docs](https://github.com/vigilnetwork/vgl/tree/master/docs) folder.

## License

vgld is licensed under the [copyfree](http://copyfree.org) ISC License.
