// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hdkeychain

// References:
//   [BIP32]: BIP0032 - Hierarchical Deterministic Wallets
//   https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/kdsmith18542/vigil/base58"
	
	"github.com/kdsmith18542/vigil/crypto/ripemd160"
	"github.com/kdsmith18542/vigil/VGLec/secp256k1/v4"
	"github.com/kdsmith18542/vigil/kawpow"
)

const (
	// RecommendedSeedLen is the recommended length in bytes for a seed
	// to a master node.
	RecommendedSeedLen = 32 // 256 bits

	// HardenedKeyStart is the index at which a hardened key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardened child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = 0x80000000 // 2^31

	// MinSeedBytes is the minimum number of bytes allowed for a seed to
	// a master node.
	MinSeedBytes = 16 // 128 bits

	// MaxSeedBytes is the maximum number of bytes allowed for a seed to
	// a master node.
	MaxSeedBytes = 64 // 512 bits

	// serializedKeyLen is the length of a serialized public or private
	// extended key.  It consists of 4 bytes version, 1 byte depth, 4 bytes
	// fingerprint, 4 bytes child number, 32 bytes chain code, and 33 bytes
	// public/private key data.
	serializedKeyLen = 4 + 1 + 4 + 4 + 32 + 33 // 78 bytes
)

var (
	// ErrDeriveHardFromPublic describes an error in which the caller
	// attempted to derive a hardened extended key from a public key.
	ErrDeriveHardFromPublic = errors.New("cannot derive a hardened key " +
		"from a public key")

	// ErrNotPrivExtKey describes an error in which the caller attempted
	// to extract a private key from a public extended key.
	ErrNotPrivExtKey = errors.New("unable to create private keys from a " +
		"public extended key")

	// ErrInvalidChild describes an error in which the child at a specific
	// index is invalid due to the derived key falling outside of the valid
	// range for secp256k1 private keys.  This error indicates the caller
	// should simply ignore the invalid child extended key at this index and
	// increment to the next index.
	ErrInvalidChild = errors.New("the extended key at this index is invalid")

	// ErrUnusableSeed describes an error in which the provided seed is not
	// usable due to the derived key falling outside of the valid range for
	// secp256k1 private keys.  This error indicates the caller must choose
	// another seed.
	ErrUnusableSeed = errors.New("unusable seed")

	// ErrInvalidSeedLen describes an error in which the provided seed or
	// seed length is not in the allowed range.
	ErrInvalidSeedLen = fmt.Errorf("seed length must be between %d and %d "+
		"bits", MinSeedBytes*8, MaxSeedBytes*8)

	// ErrBadChecksum describes an error in which the checksum encoded with
	// a serialized extended key does not match the calculated value.
	ErrBadChecksum = errors.New("bad extended key checksum")

	// ErrInvalidKeyLen describes an error in which the provided serialized
	// key is not the expected length.
	ErrInvalidKeyLen = errors.New("the provided serialized extended key " +
		"length is invalid")

	// ErrWrongNetwork describes an error in which the provided serialized
	// key is not for the expected network.
	ErrWrongNetwork = errors.New("the provided serialized extended key " +
		"is for the wrong network")
)

// masterKey is the master key used along with a random seed used to generate
// the master node in the hierarchical tree.
var masterKey = []byte("Bitcoin seed")

// NetworkParams defines an interface that is used throughout the package to
// access the hierarchical deterministic extended key magic versions that
// uniquely identify a network.
type NetworkParams interface {
	// HDPrivKeyVersion returns the hierarchical deterministic extended private
	// key magic version bytes.
	HDPrivKeyVersion() [4]byte

	// HDPubKeyVersion returns the hierarchical deterministic extended public
	// key magic version bytes.
	HDPubKeyVersion() [4]byte
}

// ExtendedKey houses all the information needed to support a hierarchical
// deterministic extended key.  See the package overview documentation for
// more details on how to use extended keys.
type ExtendedKey struct {
	privVer   [4]byte // Network version bytes for extended priv keys
	pubVer    [4]byte // Network version bytes for extended pub keys
	key       []byte  // This will be the pubkey for extended pub keys
	pubKey    []byte  // This will only be set for extended priv keys
	chainCode []byte
	parentFP  []byte
	childNum  uint32
	depth     uint16
	isPrivate bool
}

// newExtendedKey returns a new instance of an extended key with the given
// fields.  No error checking is performed here as it's only intended to be a
// convenience method used to create a populated struct.
func newExtendedKey(privVer, pubVer [4]byte, key, chainCode, parentFP []byte,
	depth uint16, childNum uint32, isPrivate bool) *ExtendedKey {

	// NOTE: The pubKey field is intentionally left nil so it is only
	// computed and memoized as required.
	return &ExtendedKey{
		privVer:   privVer,
		pubVer:    pubVer,
		key:       key,
		chainCode: chainCode,
		depth:     depth,
		parentFP:  parentFP,
		childNum:  childNum,
		isPrivate: isPrivate,
	}
}

// pubKeyBytes returns bytes for the serialized compressed public key associated
// with this extended key in an efficient manner including memoization as
// necessary.
//
// When the extended key is already a public key, the key is simply returned as
// is since it's already in the correct form.  However, when the extended key is
// a private key, the public key will be calculated and memoized so future
// accesses can simply return the cached result.
func (k *ExtendedKey) pubKeyBytes() []byte {
	// Just return the key if it's already an extended public key.
	if !k.isPrivate {
		return k.key
	}

	// This is a private extended key, so calculate and memoize the public
	// key if needed.
	if len(k.pubKey) == 0 {
		privKey := secp256k1.PrivKeyFromBytes(k.key)
		k.pubKey = privKey.PubKey().SerializeCompressed()
	}

	return k.pubKey
}

// ChildNum returns the child number of the extended key.
func (k *ExtendedKey) ChildNum() uint32 {
	return k.childNum
}

// Depth returns the depth of the extended key.
func (k *ExtendedKey) Depth() uint16 {
	return k.depth
}

// IsPrivate returns whether or not the extended key is a private extended key.
//
// A private extended key can be used to derive both hardened and non-hardened
// child private and public extended keys.  A public extended key can only be
// used to derive non-hardened child public extended keys.
func (k *ExtendedKey) IsPrivate() bool {
	return k.isPrivate
}

// ParentFingerprint returns a fingerprint of the parent extended key from which
// this one was derived.
func (k *ExtendedKey) ParentFingerprint() uint32 {
	return binary.BigEndian.Uint32(k.parentFP)
}

// hash160 returns RIPEMD160(BLAKE256(v)).
func hash160(v []byte) []byte {
	kawpowHash := kawpow.HashFunc(v)
	h := ripemd160.New()
	h.Write(kawpowHash[:])
	return h.Sum(nil)
}

// doubleBlake256Cksum returns the first four bytes of BLAKE256(BLAKE256(v)).
func doubleBlake256Cksum(v []byte) []byte {
	first := kawpow.HashFunc(v)
	second := kawpow.HashFunc(first[:])
	return second[:4]
}

// child derives a child extended key at the given index. The derived key will
// retain any leading zeros of a private key if the strict BIP32 flag is true,
// otherwise they will be stripped.  Strict BIP32 derivation is not intended for
// Vigil wallets.  The derived extended key will be either public or private as
// determined by the IsPrivate function.
func (k *ExtendedKey) child(i uint32, strictBIP32 bool) (*ExtendedKey, error) {
	// There are four scenarios that could happen here:
	// 1) Private extended key -> Hardened child private extended key
	// 2) Private extended key -> Non-hardened child private extended key
	// 3) Public extended key -> Non-hardened child public extended key
	// 4) Public extended key -> Hardened child public extended key (INVALID!)

	// Case #4 is invalid, so error out early.
	// A hardened child extended key may not be created from a public
	// extended key.
	isChildHardened := i >= HardenedKeyStart
	if !k.isPrivate && isChildHardened {
		return nil, ErrDeriveHardFromPublic
	}

	// The data used to derive the child key depends on whether or not the
	// child is hardened per [BIP32].
	//
	// For hardened children:
	//   0x00 || ser256(parentKey) || ser32(i)
	//
	// For normal children:
	//   serP(parentPubKey) || ser32(i)
	keyLen := 33
	data := make([]byte, keyLen+4)
	if isChildHardened {
		// Case #1.
		// When the child is a hardened child, the key is known to be a
		// private key due to the above early return.  Pad it with a
		// leading zero as required by [BIP32] for deriving the child.
		copy(data[1:], k.key)
	} else {
		// Case #2 or #3.
		// This is either a public or private extended key, but in
		// either case, the data which is used to derive the child key
		// starts with the secp256k1 compressed public key bytes.
		copy(data, k.pubKeyBytes())
	}
	binary.BigEndian.PutUint32(data[keyLen:], i)

	// Take the HMAC-SHA512 of the current key's chain code and the derived
	// data:
	//   I = HMAC-SHA512(Key = chainCode, Data = data)
	hmac512 := hmac.New(sha512.New, k.chainCode)
	hmac512.Write(data)
	ilr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = intermediate key used to derive the child
	//   Ir = child chain code
	il := ilr[:len(ilr)/2]
	childChainCode := ilr[len(ilr)/2:]

	// Both derived public or private keys rely on treating the left 32-byte
	// sequence calculated above (Il) as a 256-bit integer that must be
	// within the valid range for a secp256k1 private key.  There is a small
	// chance (< 1 in 2^127) this condition will not hold, and in that case,
	// a child extended key can't be created for this index and the caller
	// should simply increment to the next index.
	var ilModN secp256k1.ModNScalar
	if overflow := ilModN.SetByteSlice(il); overflow || ilModN.IsZero() {
		return nil, ErrInvalidChild
	}

	// The algorithm used to derive the child key depends on whether or not
	// a private or public child is being derived.
	//
	// For private children:
	//   childKey = parse256(Il) + parentKey
	//
	// For public children:
	//   childKey = serP(point(parse256(Il)) + parentKey)
	var isPrivate bool
	var childKey []byte
	if k.isPrivate {
		// Case #1 or #2.
		// Add the parent private key to the intermediate private key to
		// derive the final child key.
		//
		// childKey = parse256(Il) + parentKey
		var parentPrivKeyModN secp256k1.ModNScalar
		parentPrivKeyModN.SetByteSlice(k.key)
		ilModN.Add(&parentPrivKeyModN)
		childKeyBytes := ilModN.Bytes()
		childKey = childKeyBytes[:]

		// Optionally strip leading zeroes to maintain legacy behavior.  Note
		// that per [BIP32] this should be the fully zero-padded 32-bytes,
		// however, the Vigil variation strips leading zeros for legacy reasons
		// and changing it now would break derivation for a lot of Vigil
		// wallets that rely on this behavior.
		for !strictBIP32 && len(childKey) > 0 && childKey[0] == 0x00 {
			childKey = childKey[1:]
		}
		isPrivate = true
	} else {
		// Case #3.
		// Calculate the corresponding intermediate public key for
		// intermediate private key.
		var imPubKey secp256k1.JacobianPoint
		secp256k1.ScalarBaseMultNonConst(&ilModN, &imPubKey)
		imPubKey.ToAffine()
		if imPubKey.X.IsZero() || imPubKey.Y.IsZero() {
			return nil, ErrInvalidChild
		}

		// Convert the serialized compressed parent public key into a
		// point so it can be added to the intermediate public key.
		var parentPubKey secp256k1.JacobianPoint
		pubKey, err := secp256k1.ParsePubKey(k.key)
		if err != nil {
			return nil, err
		}
		pubKey.AsJacobian(&parentPubKey)

		// Add the intermediate public key to the parent public key to
		// derive the final child key.
		//
		// childKey = serP(point(parse256(Il)) + parentKey)
		var child secp256k1.JacobianPoint
		secp256k1.AddNonConst(&imPubKey, &parentPubKey, &child)
		child.ToAffine()
		pk := secp256k1.NewPublicKey(&child.X, &child.Y)
		childKey = pk.SerializeCompressed()
	}

	// The fingerprint of the parent for the derived child is the first 4
	// bytes of the RIPEMD160(BLAKE256(parentPubKey)).
	parentFP := hash160(k.pubKeyBytes())[:4]
	return newExtendedKey(k.privVer, k.pubVer, childKey, childChainCode,
		parentFP, k.depth+1, i, isPrivate), nil
}

// Child returns a derived child extended key at the given index.  When this
// extended key is a private extended key (as determined by the IsPrivate
// function), a private extended key will be derived.  Otherwise, the derived
// extended key will be also be a public extended key.
//
// When the index is greater to or equal than the HardenedKeyStart constant, the
// derived extended key will be a hardened extended key.  It is only possible to
// derive a hardened extended key from a private extended key.  Consequently,
// this function will return ErrDeriveHardFromPublic if a hardened child
// extended key is requested from a public extended key.
//
// A hardened extended key is useful since, as previously mentioned, it requires
// a parent private extended key to derive.  In other words, normal child
// extended public keys can be derived from a parent public extended key (no
// knowledge of the parent private key) whereas hardened extended keys may not
// be.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the specific child
// index does not derive to a usable child.  The ErrInvalidChild error will be
// returned if this should occur, and the caller is expected to ignore the
// invalid child and simply increment to the next index.
//
// NOTE 2: Child keys derived from the returned extended key will follow the
// modified Vigil variation of the BIP32 derivation scheme such that any
// leading zero bytes of private keys are stripped, resulting in different
// subsequent child keys. This should be used for legacy compatibility purposes.
func (k *ExtendedKey) Child(i uint32) (*ExtendedKey, error) {
	return k.child(i, false)
}

// ChildBIP32Std is like Child, except that derived keys will follow BIP32
// strictly by retaining leading zeros in the keys, always generating 32-byte
// keys, and thus different subsequently derived child keys.
func (k *ExtendedKey) ChildBIP32Std(i uint32) (*ExtendedKey, error) {
	return k.child(i, true)
}

// Neuter returns a new extended public key from this extended private key.  The
// same extended key will be returned unaltered if it is already an extended
// public key.
//
// As the name implies, an extended public key does not have access to the
// private key, so it is not capable of signing transactions or deriving
// child extended private keys.  However, it is capable of deriving further
// child extended public keys.
func (k *ExtendedKey) Neuter() *ExtendedKey {
	// Already an extended public key.
	if !k.isPrivate {
		return k
	}

	// Convert it to an extended public key.  The key for the new extended
	// key will simply be the pubkey of the current extended private key.
	//
	// This is the function N((k,c)) -> (K, c) from [BIP32].
	return newExtendedKey(k.privVer, k.pubVer, k.pubKeyBytes(), k.chainCode,
		k.parentFP, k.depth, k.childNum, false)
}

// SerializedPubKey returns the compressed serialization of the secp256k1 public
// key.  The bytes must not be modified.
func (k *ExtendedKey) SerializedPubKey() []byte {
	return k.pubKeyBytes()
}

// SerializedPrivKey converts the extended key to a secp256k1 private key and
// returns its serialization.  The returned bytes must not be modified.
//
// As you might imagine this is only possible if the extended key is a private
// extended key (as determined by the IsPrivate function).  The ErrNotPrivExtKey
// error will be returned if this function is called on a public extended key.
func (k *ExtendedKey) SerializedPrivKey() ([]byte, error) {
	if !k.isPrivate {
		return nil, ErrNotPrivExtKey
	}

	return k.key, nil
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
	for i := 0; i < int(size)-len(src); i++ {
		dst = append(dst, 0)
	}
	return append(dst, src...)
}

// String returns the extended key as a human-readable base58-encoded string.
func (k *ExtendedKey) String() string {
	if len(k.key) == 0 {
		return "zeroed extended key"
	}

	var childNumBytes [4]byte
	depthByte := byte(k.depth % 256)
	binary.BigEndian.PutUint32(childNumBytes[:], k.childNum)

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)
	serializedBytes := make([]byte, 0, serializedKeyLen+4)
	if k.isPrivate {
		serializedBytes = append(serializedBytes, k.privVer[:]...)
	} else {
		serializedBytes = append(serializedBytes, k.pubVer[:]...)
	}
	serializedBytes = append(serializedBytes, depthByte)
	serializedBytes = append(serializedBytes, k.parentFP...)
	serializedBytes = append(serializedBytes, childNumBytes[:]...)
	serializedBytes = append(serializedBytes, k.chainCode...)
	if k.isPrivate {
		serializedBytes = append(serializedBytes, 0x00)
		serializedBytes = paddedAppend(32, serializedBytes, k.key)
	} else {
		serializedBytes = append(serializedBytes, k.pubKeyBytes()...)
	}

	checkSum := doubleBlake256Cksum(serializedBytes)
	serializedBytes = append(serializedBytes, checkSum...)
	return base58.Encode(serializedBytes)
}

// zero sets all bytes in the passed slice to zero.  This is used to
// explicitly clear private key material from memory.
func zero(b []byte) {
	lenb := len(b)
	for i := 0; i < lenb; i++ {
		b[i] = 0
	}
}

// Zero manually clears all fields and bytes in the extended key.  This can be
// used to explicitly clear key material from memory for enhanced security
// against memory scraping.  This function only clears this particular key and
// not any children that have already been derived.
func (k *ExtendedKey) Zero() {
	zero(k.key)
	zero(k.pubKey)
	zero(k.chainCode)
	zero(k.parentFP)
	k.key = nil
	k.depth = 0
	k.childNum = 0
	k.isPrivate = false
}

// NewMaster creates a new master node for use in creating a hierarchical
// deterministic key chain.  The seed must be between 128 and 512 bits and
// should be generated by a cryptographically secure random generation source.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the provided seed
// will derive to an unusable secret key.  The ErrUnusable error will be
// returned if this should occur, so the caller must check for it and generate a
// new seed accordingly.
func NewMaster(seed []byte, net NetworkParams) (*ExtendedKey, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if len(seed) < MinSeedBytes || len(seed) > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	// First take the HMAC-SHA512 of the master key and the seed data:
	//   I = HMAC-SHA512(Key = "Bitcoin seed", Data = S)
	hmac512 := hmac.New(sha512.New, masterKey)
	hmac512.Write(seed)
	lr := hmac512.Sum(nil)

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = master secret key
	//   Ir = master chain code
	secretKey := lr[:len(lr)/2]
	chainCode := lr[len(lr)/2:]

	// Ensure the key is usable.
	var priv secp256k1.ModNScalar
	if overflow := priv.SetByteSlice(secretKey); overflow || priv.IsZero() {
		return nil, ErrUnusableSeed
	}

	parentFP := []byte{0x00, 0x00, 0x00, 0x00}
	return newExtendedKey(net.HDPrivKeyVersion(), net.HDPubKeyVersion(),
		secretKey, chainCode, parentFP, 0, 0, true), nil
}

// NewKeyFromString returns a new extended key instance from a base58-encoded
// extended key which is required to be for the provided network.
func NewKeyFromString(key string, net NetworkParams) (*ExtendedKey, error) {
	// The provided encoded extended key must not be larger than the maximum
	// possible encoded size.  The base58-decoded extended key consists of the
	// serialized payload plus an additional 4 bytes for the checksum.
	//
	// Since the encoding converts from base256 to base58, the max possible
	// number of bytes of output per input byte is log_58(256) ~= 1.37.  Thus, a
	// reasonable estimate for the max possible encoded size is
	// ceil(decodedDataLen * 1.37).
	//
	// Note that the actual max size in practice is two less than this value due
	// to rounding and the network prefixes in use, however, this uses the
	// theoretical max so the code works properly with all prefixes since they
	// are parameterized.
	const decodedDataLen = serializedKeyLen + 4
	const maxKeyLen = (decodedDataLen * 137 / 100) + 1
	if len(key) > maxKeyLen {
		return nil, ErrInvalidKeyLen
	}

	// Decode the extended key and ensure it is the expected length.
	decoded := base58.Decode(key)
	if len(decoded) != decodedDataLen {
		return nil, ErrInvalidKeyLen
	}

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)

	// Split the payload and checksum up and ensure the checksum matches.
	payload := decoded[:len(decoded)-4]
	checkSum := decoded[len(decoded)-4:]
	expectedCheckSum := doubleBlake256Cksum(payload)
	if !bytes.Equal(checkSum, expectedCheckSum) {
		return nil, ErrBadChecksum
	}

	// Ensure the version encoded in the payload matches the provided network.
	privVersion := net.HDPrivKeyVersion()
	pubVersion := net.HDPubKeyVersion()
	version := payload[:4]
	if !bytes.Equal(version, privVersion[:]) &&
		!bytes.Equal(version, pubVersion[:]) {

		return nil, ErrWrongNetwork
	}

	// Deserialize the remaining payload fields.
	depth := uint16(payload[4:5][0])
	parentFP := payload[5:9]
	childNum := binary.BigEndian.Uint32(payload[9:13])
	chainCode := payload[13:45]
	keyData := payload[45:78]

	// The key data is a private key if it starts with 0x00.  Serialized
	// compressed pubkeys either start with 0x02 or 0x03.
	isPrivate := keyData[0] == 0x00
	if isPrivate {
		// Ensure the private key is valid.  It must be within the range
		// of the order of the secp256k1 curve and not be 0.
		keyData = keyData[1:]
		var priv secp256k1.ModNScalar
		if overflow := priv.SetByteSlice(keyData); overflow || priv.IsZero() {
			return nil, ErrUnusableSeed
		}
	} else {
		// Ensure the public key parses correctly and is actually on the
		// secp256k1 curve.
		_, err := secp256k1.ParsePubKey(keyData)
		if err != nil {
			return nil, err
		}
	}

	return newExtendedKey(privVersion, pubVersion, keyData, chainCode, parentFP,
		depth, childNum, isPrivate), nil
}

// GenerateSeed returns a cryptographically secure random seed that can be used
// as the input for the NewMaster function to generate a new master node.
//
// The length is in bytes and it must be between 16 and 64 (128 to 512 bits).
// The recommended length is 32 (256 bits) as defined by the RecommendedSeedLen
// constant.
func GenerateSeed(length uint8) ([]byte, error) {
	// Per [BIP32], the seed must be in range [MinSeedBytes, MaxSeedBytes].
	if length < MinSeedBytes || length > MaxSeedBytes {
		return nil, ErrInvalidSeedLen
	}

	buf := make([]byte, length)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
