// Copyright (c) 2024 The Vigil developers
// Copyright (c) 2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chainhash

import (
	"golang.org/x/crypto/sha3"
)

const (
	// HashBlockSize is the block size of the hash function in bytes.
	HashBlockSize = 64
)

// keccak256 calculates the Keccak256 hash of the input data.
func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

// HashFunc calculates hash(b) and returns the resulting bytes as an array.
func HashFunc(b []byte) [HashSize]byte {
	var result [HashSize]byte
	h := keccak256(b)
	copy(result[:], h[:HashSize])
	return result
}

// HashB calculates hash(b) and returns the resulting bytes.
func HashB(b []byte) []byte {
	hash := HashFunc(b)
	return hash[:]
}

// HashH calculates hash(b) and returns the resulting bytes as a Hash.
func HashH(b []byte) Hash {
	return Hash(HashFunc(b))
}




