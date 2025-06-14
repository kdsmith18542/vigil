// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package VGLutil

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/kdsmith18542/vigil/base58"
	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	"github.com/kdsmith18542/vigil/VGLec"
	"github.com/kdsmith18542/vigil/VGLec/edwards/v2"
	"github.com/kdsmith18542/vigil/VGLec/secp256k1/v4"
)

var (
	// ErrMalformedPrivateKey describes an error where a WIF-encoded private key
	// cannot be decoded due to being improperly formatted.  This may occur if
	// the byte length is incorrect or an unexpected magic number was
	// encountered.
	ErrMalformedPrivateKey = errors.New("malformed private key")

	// ErrChecksumMismatch describes an error where decoding failed due to a bad
	// checksum.
	ErrChecksumMismatch = errors.New("checksum mismatch")
)

// ErrWrongWIFNetwork describes an error in which the provided WIF is not for
// the expected network.
type ErrWrongWIFNetwork [2]byte

// Error implements the error interface.
func (e ErrWrongWIFNetwork) Error() string {
	return fmt.Sprintf("WIF is not for the network identified by %#04x",
		[2]byte(e))
}

// WIF contains the individual components described by the Wallet Import Format
// (WIF).  A WIF string is typically used to represent a private key and its
// associated address in a way that may be easily copied and imported into or
// exported from wallet software.  WIF strings may be decoded into this
// structure by calling DecodeWIF or created with a user-provided private key
// by calling NewWIF.
type WIF struct {
	// scheme is the type of signature scheme used.
	scheme VGLec.SignatureType

	// privKey is the private key being imported or exported.
	privKey []byte

	// pubKey is the public key of the privKey
	pubKey []byte

	// netID is the network identifier byte used when
	// WIF encoding the private key.
	netID [2]byte
}

// NewWIF creates a new WIF structure to export an address and its private key
// as a string encoded in the Wallet Import Format.  The net parameter specifies
// the magic bytes of the network for which the WIF string is intended.
func NewWIF(privKey []byte, net [2]byte, scheme VGLec.SignatureType) (*WIF, error) {
	var pubBytes []byte
	switch scheme {
	case VGLec.STEcdsaSecp256k1, VGLec.STSchnorrSecp256k1:
		priv := secp256k1.PrivKeyFromBytes(privKey)
		pubBytes = priv.PubKey().SerializeCompressed()
	case VGLec.STEd25519:
		_, pub, err := edwards.PrivKeyFromScalar(privKey)
		if err != nil {
			return nil, err
		}
		pubBytes = pub.SerializeCompressed()
	default:
		return nil, fmt.Errorf("unsupported signature type '%v'", scheme)
	}

	return &WIF{scheme, privKey, pubBytes, net}, nil
}

// DecodeWIF creates a new WIF structure by decoding the string encoding of
// the import format which is required to be for the provided network.
//
// The WIF string must be a base58-encoded string of the following byte
// sequence:
//
//   - 2 bytes to identify the network
//   - 1 byte for ECDSA type
//   - 32 bytes of a binary-encoded, big-endian, zero-padded private key
//   - 4 bytes of checksum, must equal the first four bytes of the double SHA256
//     of every byte before the checksum in this sequence
//
// If the base58-decoded byte sequence does not match this, DecodeWIF will
// return a non-nil error.  ErrMalformedPrivateKey is returned when the WIF
// is of an impossible length.  ErrChecksumMismatch is returned if the
// expected WIF checksum does not match the calculated checksum.
func DecodeWIF(wif string, net [2]byte) (*WIF, error) {
	// The provided encoded WIF must not be larger than the maximum possible
	// encoded size.
	//
	// Since the encoding converts from base256 to base58, the max possible
	// number of bytes of output per input byte is log_58(256) ~= 1.37.  Thus, a
	// reasonable estimate for the max possible encoded size is
	// ceil(decodedDataLen * 1.37).
	//
	// Note that the actual max size in practice is one less than this value due
	// to the network prefixes in use, however, this uses the theoretical max so
	// the code works properly with all prefixes since they are parameterized.
	const decodedDataLen = 39
	const maxWIFLen = 54
	if len(wif) > maxWIFLen {
		return nil, ErrMalformedPrivateKey
	}

	// Decode and ensure the decoded data is the expected length.
	decoded := base58.Decode(wif)
	decodedLen := len(decoded)
	if decodedLen != decodedDataLen {
		return nil, ErrMalformedPrivateKey
	}

	// Checksum is first four bytes of hash of the identifier byte
	// and privKey.  Verify this matches the final 4 bytes of the decoded
	// private key.
	cksum := chainhash.HashB(decoded[:decodedLen-4])
	if !bytes.Equal(cksum[:4], decoded[decodedLen-4:]) {
		return nil, ErrChecksumMismatch
	}

	netID := [2]byte{decoded[0], decoded[1]}
	if netID != net {
		return nil, ErrWrongWIFNetwork(net)
	}
	var privKeyBytes, pubKeyBytes []byte
	var scheme VGLec.SignatureType
	switch VGLec.SignatureType(decoded[2]) {
	case VGLec.STEcdsaSecp256k1:
		privKeyBytes = decoded[3 : 3+secp256k1.PrivKeyBytesLen]
		privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
		pubKeyBytes = privKey.PubKey().SerializeCompressed()
		scheme = VGLec.STEcdsaSecp256k1
	case VGLec.STEd25519:
		privKeyBytes = decoded[3 : 3+edwards.PrivScalarSize]
		_, pubKey, err := edwards.PrivKeyFromScalar(privKeyBytes)
		if err != nil {
			return nil, err
		}
		pubKeyBytes = pubKey.SerializeCompressed()
		scheme = VGLec.STEd25519
	case VGLec.STSchnorrSecp256k1:
		privKeyBytes = decoded[3 : 3+secp256k1.PrivKeyBytesLen]
		privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
		pubKeyBytes = privKey.PubKey().SerializeCompressed()
		scheme = VGLec.STSchnorrSecp256k1
	}

	return &WIF{scheme, privKeyBytes, pubKeyBytes, netID}, nil
}

// String creates the Wallet Import Format string encoding of a WIF structure.
// See DecodeWIF for a detailed breakdown of the format and requirements of
// a valid WIF string.
func (w *WIF) String() string {
	// Precalculate size.  Maximum number of bytes before base58 encoding
	// is two bytes for the network, one byte for the ECDSA type, 32 bytes
	// of private key and finally four bytes of checksum.
	encodeLen := 2 + 1 + 32 + 4

	a := make([]byte, 0, encodeLen)
	a = append(a, w.netID[:]...)
	a = append(a, byte(w.scheme))
	a = append(a, w.privKey...)

	cksum := chainhash.HashB(a)
	a = append(a, cksum[:4]...)
	return base58.Encode(a)
}

// PrivKey returns the serialized private key described by the WIF.  The bytes
// must not be modified.
func (w *WIF) PrivKey() []byte {
	return w.privKey
}

// PubKey returns the compressed serialization of the associated public key for
// the WIF's private key.
func (w *WIF) PubKey() []byte {
	return w.pubKey
}

// DSA describes the signature scheme of the key.
func (w *WIF) DSA() VGLec.SignatureType {
	return w.scheme
}
