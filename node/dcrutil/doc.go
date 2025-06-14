// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

/*
Package VGLutil provides Vigil-specific convenience functions and types.

# Block Overview

A Block defines a Vigil block that provides easier and more efficient
manipulation of raw wire protocol blocks.  It also memoizes hashes for the
block and its transactions on their first access so subsequent accesses don't
have to repeat the relatively expensive hashing operations.

# Tx Overview

A Tx defines a Vigil transaction that provides more efficient manipulation of
raw wire protocol transactions.  It memoizes the hash for the transaction on its
first access so subsequent accesses don't have to repeat the relatively
expensive hashing operations.

# Address Overview

The Address interface provides an abstraction for a Vigil address.  While the
most common type is a pay-to-pubkey-hash, Vigil already supports others and
may well support more in the future.  This package currently provides
implementations for the pay-to-pubkey, pay-to-pubkey-hash, and
pay-to-script-hash address types.
*/
package VGLutil
