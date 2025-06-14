// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

/*
Package secp256k1 implements optimized secp256k1 elliptic curve operations in
pure Go.

This package provides an optimized pure Go implementation of elliptic curve
cryptography operations over the secp256k1 curve as well as data structures and
functions for working with public and private secp256k1 keys.  See
https://www.secg.org/sec2-v2.pdf for details on the standard.

In addition, sub packages are provided to produce, verify, parse, and serialize
ECDSA signatures and EC-Schnorr-VGLv0 (a custom Schnorr-based signature scheme
specific to Vigil) signatures.  See the README.md files in the relevant sub
packages for more details about those aspects.

An overview of the features provided by this package are as follows:

  - Private key generation, serialization, and parsing
  - Public key generation, serialization and parsing per ANSI X9.62-1998
  - Parses uncompressed, compressed, and hybrid public keys
  - Serializes uncompressed and compressed public keys
  - Specialized types for performing optimized and constant time field operations
  - FieldVal type for working modulo the secp256k1 field prime
  - ModNScalar type for working modulo the secp256k1 group order
  - Elliptic curve operations in Jacobian projective coordinates
  - Point addition
  - Point doubling
  - Scalar multiplication with an arbitrary point
  - Scalar multiplication with the base point (group generator)
  - Point decompression from a given x coordinate
  - Nonce generation via RFC6979 with support for extra data and version
    information that can be used to prevent nonce reuse between signing
    algorithms

It also provides an implementation of the Go standard library crypto/elliptic
Curve interface via the S256 function so that it may be used with other packages
in the standard library such as crypto/tls, crypto/x509, and crypto/ecdsa.
However, in the case of ECDSA, it is highly recommended to use the ecdsa sub
package of this package instead since it is optimized specifically for secp256k1
and is significantly faster as a result.

Although this package was primarily written for vgld, it has intentionally been
designed so it can be used as a standalone package for any projects needing to
use optimized secp256k1 elliptic curve cryptography.

Finally, a comprehensive suite of tests is provided to provide a high level of
quality assurance.

# Use of secp256k1 in Vigil

At the time of this writing, the primary public key cryptography in widespread
use on the Vigil network used to secure coins is based on elliptic curves
defined by the secp256k1 domain parameters.
*/
package secp256k1
