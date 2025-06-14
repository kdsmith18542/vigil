standalone
==========

[![Build Status](https://github.com/vigilnetwork/vgl/workflows/Build%20and%20Test/badge.svg)](https://github.com/vigilnetwork/vgl/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/vigilnetwork/vgl/blockchain/standalone/v2)

Package standalone provides standalone functions useful for working with the
Vigil blockchain consensus rules.

The primary goal of offering these functions via a separate module is to reduce
the required dependencies to a minimum as compared to the blockchain module.

It is ideal for applications such as lightweight clients that need to ensure
basic security properties hold and calculate appropriate vote subsidies and
block explorers.

For example, some things an SPV wallet needs to prove are that the block headers
all connect together, that they satisfy the proof of work requirements, and that
a given transaction tree is valid for a given header.

The provided functions fall into the following categories:

- Proof-of-work
  - Converting to and from the compact target difficulty representation
  - Calculating work values based on the compact target difficulty
  - Checking a block hash satisfies a target difficulty and that target
    difficulty is within a valid range
  - Calculating required target difficulties using the ASERT algorithm
- Merkle root calculation
  - Calculation from individual leaf hashes
  - Calculation from a slice of transactions
- Subsidy calculation
  - Proof-of-work subsidy for a given height and number of votes
  - Stake vote subsidy for a given height
  - Treasury subsidy for a given height and number of votes
- Coinbase transaction identification
 - Merkle tree inclusion proofs
   - Generate an inclusion proof for a given tree and leaf index
   - Verify a leaf is a member of the tree at a given index via the proof
- Transaction sanity checking

## Installation and Updating

This package is part of the `github.com/vigilnetwork/vgl/blockchain/standalone/v2`
module.  Use the standard go tooling for working with modules to incorporate it.

## Examples

* [CompactToBig Example](https://pkg.go.dev/github.com/vigilnetwork/vgl/blockchain/standalone#example-CompactToBig)
  Demonstrates how to convert the compact "bits" in a block header which
  represent the target difficulty to a big integer and display it using the
  typical hex notation.

* [BigToCompact Example](https://pkg.go.dev/github.com/vigilnetwork/vgl/blockchain/standalone#example-BigToCompact)
  Demonstrates how to convert a target difficulty into the compact "bits" in a
  block header which represent that target difficulty.

* [CheckProofOfWork Example](https://pkg.go.dev/github.com/vigilnetwork/vgl/blockchain/standalone#example-CheckProofOfWork)
  Demonstrates checking the proof of work of a block hash against a target
  difficulty.

* [CalcMerkleRoot Example](https://pkg.go.dev/github.com/vigilnetwork/vgl/blockchain/standalone#example-CalcMerkleRoot)
  Demonstrates calculating a merkle root from a slice of leaf hashes.

## License

Package standalone is licensed under the [copyfree](http://copyfree.org) ISC
License.
