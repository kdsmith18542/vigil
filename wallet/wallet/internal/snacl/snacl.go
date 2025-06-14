// Copyright (c) 2014-2024 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package snacl

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"runtime/debug"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/crypto/rand"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
)

// Various constants needed for encryption scheme.
const (
	// Expose secretbox's Overhead const here for convenience.
	Overhead  = secretbox.Overhead
	KeySize   = 32
	NonceSize = 24
	DefaultN  = 16384 // 2^14
	DefaultR  = 8
	DefaultP  = 1
)

// CryptoKey represents a secret key which can be used to encrypt and decrypt
// data.
type CryptoKey [KeySize]byte

// Encrypt encrypts the passed data.
func (ck *CryptoKey) Encrypt(in []byte) ([]byte, error) {
	const op errors.Op = "cryptokey.Encrypt"
	var nonce [NonceSize]byte
	rand.Read(nonce[:])
	sealed := make([]byte, NonceSize, NonceSize+len(in)+Overhead)
	copy(sealed, nonce[:])
	sealed = secretbox.Seal(sealed, in, &nonce, (*[KeySize]byte)(ck))
	return sealed, nil
}

// Decrypt decrypts the passed data.  The must be the output of the Encrypt
// function.
func (ck *CryptoKey) Decrypt(in []byte) ([]byte, error) {
	const op errors.Op = "cryptokey.Decrypt"
	if len(in) < NonceSize {
		return nil, errors.E(op, errors.Invalid, "missing nonce")
	}

	var nonce [NonceSize]byte
	copy(nonce[:], in[:NonceSize])
	blob := in[NonceSize:]

	opened, ok := secretbox.Open(nil, blob, &nonce, (*[KeySize]byte)(ck))
	if !ok {
		return nil, errors.E(op, errors.Crypto)
	}

	return opened, nil
}

// Zero clears the key by manually zeroing all memory.  This is for security
// conscience application which wish to zero the memory after they've used it
// rather than waiting until it's reclaimed by the garbage collector.  The
// key is no longer usable after this call.
func (ck *CryptoKey) Zero() {
	*ck = [KeySize]byte{}
}

// GenerateCryptoKey generates a new crypotgraphically random key.
func GenerateCryptoKey() (*CryptoKey, error) {
	const op errors.Op = "snacl.GenerateCryptoKey"
	var key CryptoKey
	rand.Read(key[:])
	return &key, nil
}

// Parameters are not secret and can be stored in plain text.
type Parameters struct {
	Salt   [KeySize]byte
	Digest [sha256.Size]byte
	N      int
	R      int
	P      int
}

// SecretKey houses a crypto key and the parameters needed to derive it from a
// passphrase.  It should only be used in memory.
type SecretKey struct {
	Key        *CryptoKey
	Parameters Parameters
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// deriveKey fills out the Key field.
func (sk *SecretKey) deriveKey(op errors.Op, password *[]byte) error {
	key, err := scrypt.Key(*password, sk.Parameters.Salt[:],
		sk.Parameters.N,
		sk.Parameters.R,
		sk.Parameters.P,
		len(sk.Key))
	if err != nil {
		return errors.E(op, err)
	}
	copy(sk.Key[:], key)
	zero(key)

	// I'm not a fan of forced garbage collections, but scrypt allocates a
	// ton of memory and calling it back to back without a GC cycle in
	// between means you end up needing twice the amount of memory.  For
	// example, if your scrypt parameters are such that you require 1GB and
	// you call it twice in a row, without this you end up allocating 2GB
	// since the first GB probably hasn't been released yet.
	debug.FreeOSMemory()

	return nil
}

// Marshal returns the Parameters field marshalled into a format suitable for
// storage.  This result of this can be stored in clear text.
func (sk *SecretKey) Marshal() []byte {
	params := &sk.Parameters

	// The marshalled format for the params is as follows:
	//   <salt><digest><N><R><P>
	//
	// KeySize + sha256.Size + N (8 bytes) + R (8 bytes) + P (8 bytes)
	marshalled := make([]byte, KeySize+sha256.Size+24)

	b := marshalled
	copy(b[:KeySize], params.Salt[:])
	b = b[KeySize:]
	copy(b[:sha256.Size], params.Digest[:])
	b = b[sha256.Size:]
	binary.LittleEndian.PutUint64(b[:8], uint64(params.N))
	b = b[8:]
	binary.LittleEndian.PutUint64(b[:8], uint64(params.R))
	b = b[8:]
	binary.LittleEndian.PutUint64(b[:8], uint64(params.P))

	return marshalled
}

// Unmarshal unmarshalls the parameters needed to derive the secret key from a
// passphrase into sk.
func (sk *SecretKey) Unmarshal(marshalled []byte) error {
	const op errors.Op = "secretkey.Unmarshal"
	if sk.Key == nil {
		sk.Key = (*CryptoKey)(&[KeySize]byte{})
	}

	// The marshalled format for the params is as follows:
	//   <salt><digest><N><R><P>
	//
	// KeySize + sha256.Size + N (8 bytes) + R (8 bytes) + P (8 bytes)
	if len(marshalled) != KeySize+sha256.Size+24 {
		return errors.E(op, errors.Encoding, errors.Errorf("bad marshalled data len %d", len(marshalled)))
	}

	params := &sk.Parameters
	copy(params.Salt[:], marshalled[:KeySize])
	marshalled = marshalled[KeySize:]
	copy(params.Digest[:], marshalled[:sha256.Size])
	marshalled = marshalled[sha256.Size:]
	params.N = int(binary.LittleEndian.Uint64(marshalled[:8]))
	marshalled = marshalled[8:]
	params.R = int(binary.LittleEndian.Uint64(marshalled[:8]))
	marshalled = marshalled[8:]
	params.P = int(binary.LittleEndian.Uint64(marshalled[:8]))

	return nil
}

// Zero zeroes the underlying secret key while leaving the parameters intact.
// This effectively makes the key unusable until it is derived again via the
// DeriveKey function.
func (sk *SecretKey) Zero() {
	sk.Key.Zero()
}

// DeriveKey derives the underlying secret key and ensures it matches the
// expected digest.  This should only be called after previously calling the
// Zero function or on an initial Unmarshal.
func (sk *SecretKey) DeriveKey(password *[]byte) error {
	const op errors.Op = "secretkey.DeriveKey"
	if err := sk.deriveKey(op, password); err != nil {
		return err
	}

	// verify password
	digest := sha256.Sum256(sk.Key[:])
	if subtle.ConstantTimeCompare(digest[:], sk.Parameters.Digest[:]) != 1 {
		return errors.E(op, errors.Passphrase)
	}

	return nil
}

// Encrypt encrypts in bytes and returns a JSON blob.
func (sk *SecretKey) Encrypt(in []byte) ([]byte, error) {
	const op errors.Op = "secretkey.Encrypt"
	out, err := sk.Key.Encrypt(in)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return out, nil
}

// Decrypt takes in a JSON blob and returns it's decrypted form.
func (sk *SecretKey) Decrypt(in []byte) ([]byte, error) {
	const op errors.Op = "secretkey.Decrypt"
	out, err := sk.Key.Decrypt(in)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return out, nil
}

// NewSecretKey returns a SecretKey structure based on the passed parameters.
func NewSecretKey(password *[]byte, N, r, p int) (*SecretKey, error) {
	const op errors.Op = "snacl.NewSecretKey"
	sk := SecretKey{
		Key: (*CryptoKey)(&[KeySize]byte{}),
	}
	// setup parameters
	sk.Parameters.N = N
	sk.Parameters.R = r
	sk.Parameters.P = p
	rand.Read(sk.Parameters.Salt[:])

	// derive key
	err := sk.deriveKey(op, password)
	if err != nil {
		return nil, err
	}

	// store digest
	sk.Parameters.Digest = sha256.Sum256(sk.Key[:])

	return &sk, nil
}
