package crypto

import (
	"hash"

	"github.com/Vigil-Labs/vgl/kawpow"
)

// Size is the size of a Keccak256 hash in bytes.
const Size = 32

// Hasher256 provides a zero-allocation implementation to compute a rolling
// Keccak256 checksum.
type Hasher256 struct {
	h hash.Hash
}

// Write adds the given bytes to the rolling hash.
func (h *Hasher256) Write(b []byte) (int, error) {
	return h.h.Write(b)
}

// WriteBytes adds the given bytes to the rolling hash.
func (h *Hasher256) WriteBytes(b []byte) {
	h.h.Write(b)
}

// WriteByte adds the given byte to the rolling hash.
func (h *Hasher256) WriteByte(b byte) {
	h.h.Write([]byte{b})
}

// WriteString adds the given string to the rolling hash.
func (h *Hasher256) WriteString(s string) {
	h.h.Write([]byte(s))
}

// Reset resets the state of the rolling hash.
func (h *Hasher256) Reset() {
	h.h.Reset()
}

// Size returns the size of a Keccak256 hash in bytes.
func (h *Hasher256) Size() int {
	return h.h.Size()
}

// BlockSize returns the underlying block size of the Keccak256 hashing
// algorithm.
func (h *Hasher256) BlockSize() int {
	return h.h.BlockSize()
}

// Sum finalizes the rolling hash, appends the resulting checksum to the
// provided slice and returns the resulting slice.
func (h Hasher256) Sum(b []byte) []byte {
	return h.h.Sum(b)
}

// Sum256 finalizes the rolling hash and returns the resulting checksum.
func (h Hasher256) Sum256() [Size]byte {
	var sum [Size]byte
	copy(sum[:], h.h.Sum(nil))
	return sum
}

// NewHasher256 returns a zero-allocation hasher for computing a rolling
// Keccak256 checksum.
func NewHasher256() *Hasher256 {
	return &Hasher256{kawpow.NewKeccak256Hasher()}
}

// Sum256 returns the Keccak256 checksum of the data.
func Sum256(data []byte) [Size]byte {
	var sum [Size]byte
	copy(sum[:], kawpow.Keccak256(data))
	return sum
}



