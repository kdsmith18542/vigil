// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

/*
Package hdkeychain provides an API for Vigil hierarchical deterministic
extended keys (based on BIP0032).

The ability to implement hierarchical deterministic wallets depends on the
ability to create and derive hierarchical deterministic extended keys.

At a high level, this package provides support for those hierarchical
deterministic extended keys by providing an ExtendedKey type and supporting
functions.  Each extended key can either be a private or public extended key
which itself is capable of deriving a child extended key.

# Determining the Extended Key Type

Whether an extended key is a private or public extended key can be determined
with the IsPrivate function.

# Transaction Signing Keys and Payment Addresses

In order to create and sign transactions, or provide others with addresses to
send funds to, the underlying key and address material must be accessible.  This
package provides the SerializedPubKey and SerializedPrivKey functions for this
purpose.  The caller may then create the desired address types.

# The Master Node

As previously mentioned, the extended keys are hierarchical meaning they are
used to form a tree.  The root of that tree is called the master node and this
package provides the NewMaster function to create it from a cryptographically
random seed.  The GenerateSeed function is provided as a convenient way to
create a random seed for use with the NewMaster function.

# Deriving Children

Once you have created a tree root (or have deserialized an extended key as
discussed later), the child extended keys can be derived by using either the
Child or ChildBIP32Std function.  The difference is described in the following
section.  These functions support deriving both normal (non-hardened) and
hardened child extended keys.  In order to derive a hardened extended key, use
the HardenedKeyStart constant + the hardened key number as the index to the
Child function.  This provides the ability to cascade the keys into a tree and
hence generate the hierarchical deterministic key chains.

# BIP0032 Conformity

The Child function derives extended keys with a modified scheme based on
BIP0032, whereas ChildBIP32Std produces keys that strictly conform to the
standard.  Specifically, the Vigil variation strips leading zeros of a private
key, causing subsequent child keys to differ from the keys expected by standard
BIP0032.  The ChildBIP32Std method retains leading zeros, ensuring the child
keys expected by BIP0032 are derived.  The Child function must be used for
Vigil wallet key derivation for legacy reasons.

# Normal vs Hardened Child Extended Keys

A private extended key can be used to derive both hardened and non-hardened
(normal) child private and public extended keys.  A public extended key can only
be used to derive non-hardened child public extended keys.  As enumerated in
BIP0032 "knowledge of the extended public key plus any non-hardened private key
descending from it is equivalent to knowing the extended private key (and thus
every private and public key descending from it).  This means that extended
public keys must be treated more carefully than regular public keys. It is also
the reason for the existence of hardened keys, and why they are used for the
account level in the tree. This way, a leak of an account-specific (or below)
private key never risks compromising the master or other accounts."

# Neutering a Private Extended Key

A private extended key can be converted to a new instance of the corresponding
public extended key with the Neuter function.  The original extended key is not
modified.  A public extended key is still capable of deriving non-hardened child
public extended keys.

# Serializing and Deserializing Extended Keys

Extended keys are serialized and deserialized with the String and
NewKeyFromString functions.  The serialized key is a Base58-encoded string which
looks like the following:

	public key:   dpubZCGVaKZBiMo7pMgLaZm1qmchjWenTeVcUdFQkTNsFGFEA6xs4EW8PKiqYqP7HBAitt9Hw16VQkQ1tjsZQSHNWFc6bEK6bLqrbco24FzBTY4
	private key:  dprv3kUQDBztdyjKuwnaL3hfKYpT7W6X2huYH5d61YSWFBebSYwEBHAXJkCpQ7rvMAxPzKqxVCGLvBqWvGxXjAyMJsV1XwKkfnQCM9KctC8k8bk

# Network

Extended keys are much like normal Vigil addresses in that they have version
bytes which tie them to a specific network.  The network that an extended key is
associated with is specified when creating and decoding the key.  In the case of
decoding, an error will be returned if a given encoded extended key is not for
the specified network.
*/
package hdkeychain
