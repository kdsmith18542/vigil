// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
//
// Main Go code originally written and optimized by Dave Collins May 2020.
// Additional cleanup and comments added in July 2024.

package blake256

import (
	"hash"
)

// iv224 is the BLAKE-224 initialization vector.
var iv224 = [8]uint32{
	0xc1059ed8, 0x367cd507, 0x3070dd17, 0xf70e5939,
	0xffc00b31, 0x68581511, 0x64f98fa7, 0xbefa4fa4,
}

// statePrefix224 is the prefix used when serializing the intermediate state to
// identify the state as belonging to a BLAKE-224 rolling hash.  It is the
// second value in iv224.
const statePrefix224 = 0x367cd507

// Hasher224 provides a zero-allocation implementation to compute a rolling
// BLAKE-224 checksum.
//
// It can safely be copied at any point to save its intermediate state for use
// in additional processing later, without having to write the previously
// written data again.
//
// In addition to the aforementioned in-process state saving capability, it also
// supports serializing the intermediate state to enable sharing across process
// boundaries.
//
// It is effectively a mix of a [hash.Hash], [encoding.BinaryMarshaler], and
// [encoding.BinaryUnmarshaler] with a modified API that enables zero
// allocations and also provides additional convenience funcs for writing
// integers encoded with both big and little endian as well as writing
// individual bytes.
//
// However, it also implements [hash.Hash], [encoding.BinaryMarshaler], and
// [encoding.BinaryUnmarshaler] for callers that aren't as concerned about
// reducing allocations and would prefer to use it with the aforementioned
// standard library interfaces.
//
// NOTE: The zero value is NOT safe to use.  It must be initialized via
// NewHasher224 or NewHasher224Salt.
type Hasher224 struct {
	h hasher
}

// Write adds the given bytes to the rolling hash.
//
// NOTE: This method only returns an error in order to satisfy the [io.Writer]
// and [hash.Hash] interfaces.  However, it will never error, meaning the error
// will always be nil, so it is safe to ignore.
//
// Callers may optionally choose to call [WriteBytes] which does not return an
// error to make the fact writing can never fail.
func (h *Hasher224) Write(b []byte) (int, error) {
	return h.h.write(b)
}

// WriteByte adds the given byte to the rolling hash.
func (h *Hasher224) WriteByte(b byte) {
	h.h.writeByte(b)
}

// WriteBytes adds the given bytes to the rolling hash.
//
// This method is identical to [Write] except it does not return an error in
// order to make it clear that writing can never fail.
func (h *Hasher224) WriteBytes(b []byte) {
	h.h.write(b)
}

// WriteString adds the given string to the rolling hash.
func (h *Hasher224) WriteString(s string) {
	h.h.writeString(s)
}

// WriteUint16LE encodes the given unsigned 16-bit integer as a 2-byte
// little-endian byte sequence and adds it to the rolling hash.
func (h *Hasher224) WriteUint16LE(val uint16) {
	h.h.writeUint16LE(val)
}

// WriteUint16BE encodes the given unsigned 16-bit integer as a 2-byte
// big-endian byte sequence and adds it to the rolling hash.
func (h *Hasher224) WriteUint16BE(val uint16) {
	h.h.writeUint16BE(val)
}

// WriteUint32LE encodes the given unsigned 32-bit integer as a 4-byte
// little-endian byte sequence and adds it to the rolling hash.
func (h *Hasher224) WriteUint32LE(val uint32) {
	h.h.writeUint32LE(val)
}

// WriteUint32BE encodes the given unsigned 32-bit integer as a 4-byte
// big-endian byte sequence and adds it to the rolling hash.
func (h *Hasher224) WriteUint32BE(val uint32) {
	h.h.writeUint32BE(val)
}

// WriteUint64LE encodes the given unsigned 64-bit integer as an 8-byte
// little-endian byte sequence and adds it to the rolling hash.
func (h *Hasher224) WriteUint64LE(val uint64) {
	h.h.writeUint64LE(val)
}

// WriteUint64BE encodes the given unsigned 64-bit integer as an 8-byte
// big-endian byte sequence and adds it to the rolling hash.
func (h *Hasher224) WriteUint64BE(val uint64) {
	h.h.writeUint64BE(val)
}

// Reset resets the state of the rolling hash.
//
// This is part of the [hash.Hash] interface.
func (h *Hasher224) Reset() {
	h.h.reset(iv224)
}

// Size returns the size of a BLAKE-224 hash in bytes.
//
// This is part of the [hash.Hash] interface.
func (h *Hasher224) Size() int {
	return Size224
}

// BlockSize returns the underlying block size of the BLAKE-224 hashing
// algorithm.
//
// This is part of the [hash.Hash] interface.
func (h *Hasher224) BlockSize() int {
	return BlockSize
}

// Sum finalizes the rolling hash, appends the resulting checksum to the
// provided slice and returns the resulting slice.  It does not change the
// underlying hash state.
//
// Note that allocations can often be avoided by providing a slice that has
// enough capacity to house the resulting checksum.  For example:
//
//	digest := make([]byte, blake256.Size224)
//	h := blake256.NewHasher224()
//	h.WriteUint64LE(1)
//	digest = h.Sum(digest[:0])
//
// This is part of the [hash.Hash] interface.
func (h Hasher224) Sum(b []byte) []byte {
	// Note h is a copy so that the caller can keep writing and summing.
	sum := h.h.finalize224()
	return append(b, sum[:]...)
}

// Sum224 finalizes the rolling hash and returns the resulting checksum.  It
// does not change the underlying hash state.
func (h Hasher224) Sum224() [Size224]byte {
	// Note h is a copy so that the caller can keep writing and summing.
	return h.h.finalize224()
}

// SaveState appends the current intermediate state of the rolling hash, as
// generated by [Hasher224.MarshalBinary], to the provided slice and returns the
// resulting slice.  It does not change the underlying hash state.
//
// The resulting serialized data may be used to resume from the current
// intermediate state later without having to write the previously written data
// again by providing it to [Hasher224.UnmarshalBinary].
//
// As described by the [Hasher224] documentation, the hasher instance can simply
// be copied to achieve the same result much more efficiently when the caller is
// able to keep a copy.  Therefore, that approach should be preferred when
// possible.
//
// However, the ability to serialize the state is also provided to enable
// sharing it across process boundaries.
//
// Note that allocations can typically be avoided by providing a slice that has
// enough capacity to house the resulting state as defined by the
// [SavedStateSize] constant.  For example:
//
//	state := make([]byte, blake256.SavedStateSize)
//	h := blake256.NewHasher224()
//	h.WriteUint64LE(1)
//	state = h.SaveState(state[:0])
func (h *Hasher224) SaveState(target []byte) []byte {
	return h.h.saveState(target, statePrefix224)
}

// MarshalBinary returns the intermediate state of the rolling hash serialized
// into a binary form that may be used to resume from the current state later
// without having to write the previously written data again.  It does not
// change the underlying hash state.
//
// As described by the [Hasher224] documentation, the hasher instance can simply
// be copied to achieve the same result much more efficiently when the caller is
// able to keep a copy.  Therefore, that approach should be preferred when
// possible.
//
// However, the ability to serialize the state is also provided to enable
// sharing it across process boundaries.
//
// NOTE: This method only returns an error in order to satisfy the
// [encoding.BinaryMarshaler] interface.  However, it will never error, meaning
// the error will always be nil, so it is safe to ignore.
//
// Callers that wish to avoid allocations should prefer [Hasher224.SaveState]
// instead.
func (h *Hasher224) MarshalBinary() ([]byte, error) {
	var state [SavedStateSize]byte
	h.h.putSavedState(state[:], statePrefix224)
	return state[:], nil
}

// UnmarshalBinary restores the rolling hash to the provided serialized
// intermediate state.  See [Hasher224.MarshalBinary] for more details.
//
// [ErrMalformedState] will be returned when the provided serialized state is
// not at least the required [SavedStateSize] number of bytes.
//
// [ErrMismatchedState] will be returned if the provided state is not for a
// BLAKE-224 hash.  For example, it will be returned when attempting to restore
// a BLAKE-256 intermediate state.
//
// This implements the [encoding.BinaryUnmarshaler] interface.
func (h *Hasher224) UnmarshalBinary(state []byte) error {
	return h.h.loadState(state, statePrefix224)
}

// NewHasher224 returns a zero-allocation hasher for computing a rolling
// BLAKE-224 checksum.
func NewHasher224() *Hasher224 {
	h := Hasher224{makeHasher(iv224)}
	return &h
}

// NewHasher224Salt returns a zero-allocation hasher for computing a rolling
// BLAKE-224 checksum initialized with the given 16-byte salt slice.
//
// It will panic if the provided salt is not 16 bytes.
func NewHasher224Salt(salt []byte) *Hasher224 {
	h := Hasher224{makeHasher(iv224)}
	h.h.initializeSalt(salt)
	return &h
}

// New224 returns a new [hash.Hash] computing the BLAKE-224 checksum.
//
// Callers should prefer [NewHasher224] instead since it returns a concrete type
// that has more functionality and allows avoiding additional allocations.  It
// can also be used as a [hash.Hash] if desired.
func New224() hash.Hash {
	return NewHasher224()
}

// New224Salt returns a new [hash.Hash] computing the BLAKE-224 checksum
// initialized with the given 16-byte salt.
//
// It will panic if the provided salt is not 16 bytes.
//
// Callers should prefer [NewHasher224Salt] instead since it returns a concrete
// type that has more functionality and allows avoiding additional allocations.
// It can also be used as a [hash.Hash] if desired.
func New224Salt(salt []byte) hash.Hash {
	return NewHasher224Salt(salt)
}

// Sum224 returns the BLAKE-224 checksum of the data.
func Sum224(data []byte) [Size224]byte {
	h := makeHasher(iv224)
	h.write(data)
	return h.finalize224()
}
