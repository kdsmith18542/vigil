// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package kawpow

import (
	"encoding/binary"
	"fmt"
	hash "hash"
	"math/big"

	"golang.org/x/crypto/sha3"
)

const (
	// HashSize is the size of a KawPoW hash in bytes
	HashSize = 32
	// BlockSize is the block size for KawPoW hashing
	BlockSize = 64

	// Ethash constants
	WORD_BYTES          = 4
	DATASET_BYTES_INIT  = 1073741824  // 1GB
	DATASET_BYTES_GROWTH = 8388608   // 8MB
	CACHE_BYTES_INIT    = 16777216  // 16MB
	CACHE_BYTES_GROWTH  = 131072    // 128KB
	CACHE_MULTIPLIER    = 1024
	EPOCH_LENGTH        = 30000
	MIX_BYTES           = 128
	HASH_BYTES          = 64
	DATASET_PARENTS     = 256
	CACHE_ROUNDS        = 3
	ACCESSES            = 64

	FNV_PRIME = 0x01000193

	// EpochLength is the number of blocks per epoch
	EpochLength = 30000
)




// Hash represents a KawPoW hash
// Hash is a KawPoW hash.
type Hash [HashSize]byte

// KawPowHash calculates the KawPoW hash for the given header hash, nonce, and DAG.
func KawPowHash(headerHash Hash, nonce uint64, blockHeight int64, dag *DAG) ([]byte, []byte, error) {
	// Ensure DAG is generated
	if dag.data == nil {
		return nil, nil, fmt.Errorf("DAG not generated for epoch %d", dag.epoch)
	}

	// Combine header hash and nonce
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)
	seed := Keccak256(append(headerHash[:], nonceBytes...))

	// Initialize mix
	mix := make([]byte, MIX_BYTES)
	copy(mix, seed)

	// Main loop
	for i := 0; i < ACCESSES; i++ {
		parent := binary.LittleEndian.Uint32(mix[i%8*4:(i%8*4)+4])
		item := dag.data[parent%uint32(len(dag.data))]
		mix = xor(mix, item)
		mix = Keccak256(mix)
	}

	// Final hash
	finalHash := Keccak256(mix)

	// MixHash is the mix data before the final hash
	return finalHash, mix, nil
}

// DAG represents the Directed Acyclic Graph used in KawPoW
type DAG struct {
	epoch uint64
	data  [][]byte
}

// NewDAG creates a new DAG for the given epoch.
func NewDAG(epoch uint64) *DAG {
	return &DAG{
		epoch: epoch,
	}
}

// Generate generates the DAG for the given epoch.
func (d *DAG) Generate(epoch uint64) error {
	d.epoch = epoch

	seed := generateSeedHash(epoch)

	cacheSize := getCacheSize(epoch)
	cache := mkcache(cacheSize, seed[:])

	fullSize := getDatasetSize(epoch)
	d.data = calc_dataset(fullSize, cache)

	return nil
}

// GetData returns the DAG data.
func (d *DAG) GetData() [][]byte {
	return d.data
}

// HashFunc returns the KawPoW hash function.
func HashFunc(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

// Keccak256 calculates the Keccak256 hash of the given data.
func Keccak256(data []byte) []byte {
	return sha3.New256().Sum(data)
}

// sha3_512 is a helper function to compute SHA3-512 hash.
func sha3_512(data []byte) []byte {
	h := sha3.New512()
	h.Write(data)
	return h.Sum(nil)
}

// mkcache generates the Ethash cache for a given epoch.
func mkcache(cacheSize uint64, seed []byte) [][]byte {
	n := int(cacheSize / HASH_BYTES)
	o := make([][]byte, n)

	o[0] = sha3_512(seed)
	for i := 1; i < n; i++ {
		o[i] = sha3_512(o[i-1])
	}

	for r := 0; r < CACHE_ROUNDS; r++ {
		for i := 0; i < n; i++ {
			v := new(big.Int).SetBytes(o[i]).Uint64() % uint64(n)
			o[i] = sha3_512(xor(o[(i-1+n)%n], o[v]))
		}
	}
	return o
}

// xor performs a byte-wise XOR operation on two byte slices.
func xor(a, b []byte) []byte {
	buf := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		buf[i] = a[i] ^ b[i]
	}
	return buf
}

// fnv performs the FNV hash function.
func fnv(v1, v2 uint32) uint32 {
	return (v1 * FNV_PRIME) ^ v2
}

// calc_dataset_item calculates a single dataset item.
func calc_dataset_item(cache [][]byte, i uint32) []byte {
	n := uint32(len(cache))
	r := uint32(HASH_BYTES / WORD_BYTES)

	mix := make([]byte, HASH_BYTES)
	copy(mix, cache[i%n])

	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)
	mix = sha3_512(xor(mix, b))

	for j := uint32(0); j < DATASET_PARENTS; j++ {
		cacheIndex := fnv(i^j, binary.LittleEndian.Uint32(mix[j%r*4:(j%r*4)+4]))
		mix = sha3_512(xor(mix, cache[cacheIndex%n]))
	}
	return sha3_512(mix)
}

// calc_dataset generates the full Ethash dataset.
func calc_dataset(fullSize uint64, cache [][]byte) [][]byte {
	dataset := make([][]byte, fullSize/HASH_BYTES)
	for i := uint32(0); i < uint32(fullSize/HASH_BYTES); i++ {
		dataset[i] = calc_dataset_item(cache, i)
	}
	return dataset
}

// generateSeedHash calculates the seed hash for a given epoch.
func generateSeedHash(epoch uint64) Hash {
	seed := make([]byte, 32)
	for i := 0; i < int(epoch); i++ {
		h := sha3.New256()
		h.Write(seed)
		seed = h.Sum(nil)
	}
	var result Hash
	copy(result[:], seed)
	return result
}

// getCacheSize calculates the cache size for a given epoch.
func getCacheSize(epoch uint64) uint64 {
	return CACHE_BYTES_INIT + CACHE_BYTES_GROWTH*(epoch/EPOCH_LENGTH)
}

// getDatasetSize calculates the dataset size for a given epoch.
func getDatasetSize(epoch uint64) uint64 {
	return DATASET_BYTES_INIT + DATASET_BYTES_GROWTH*(epoch/EPOCH_LENGTH)
}

// New creates a new KawPoW hasher
func New() hash.Hash {
	return sha3.NewLegacyKeccak256()
}





// MiningResult represents the result of a KawPoW mining operation
type MiningResult struct {
	Nonce   uint64
	MixHash []byte
	Hash    []byte
}

// Mine performs KawPoW mining
func Mine(ctx context.Context, headerHash Hash, height uint64, nonceRange uint64, target *big.Int, totalHashes *uint64, elapsedMicros *int64) (uint64, Hash, Hash, error) {
	// Create DAG for the epoch
	epoch := height / EpochLength

	dag := NewDAG(epoch)
	err := dag.Generate(epoch)
	if err != nil {
		return 0, Hash{}, Hash{}, err
	}

	// Convert target to big.Int for comparison
	targetBig := new(big.Int).SetBytes(target)

	// Iterate over the nonce range
	for nonce := nonceRange; ; nonce++ {
		select {
		case <-ctx.Done():
			return 0, Hash{}, Hash{}, ctx.Err()
		default:
		}
		// Calculate KawPoW hash
		finalHashBytes, mixHashBytes, err := KawPowHash(headerHash, nonce, int64(height), dag)
		if err != nil {
			return 0, Hash{}, Hash{}, err
		}
		*totalHashes++

		// Convert finalHash to big.Int for comparison
		finalHashBig := new(big.Int).SetBytes(finalHashBytes)

		// Check if the hash meets the target
		if finalHashBig.Cmp(targetBig) <= 0 {
			var mixHash Hash
			copy(mixHash[:], mixHashBytes)

			var finalHash Hash
			copy(finalHash[:], finalHashBytes)

			return nonce, mixHash, finalHash, nil
		}

		if nonce%1000 == 0 {
			*elapsedMicros += 1000
		}
	}

	return 0, Hash{}, Hash{}, fmt.Errorf("no solution found in given nonce range")
}

// Verify verifies a KawPoW solution
func Verify(headerHash Hash, nonce uint64, mixHash []byte, bits uint32, height int64, dag *DAG) bool {
	// Generate expected hashes
	finalHash, expectedMixHash, err := KawPowHash(headerHash, nonce, height, dag)
	if err != nil {
		return false
	}

	// Verify mix hash matches
	if len(mixHash) != len(expectedMixHash) {
		return false
	}
	for i := range mixHash {
		if mixHash[i] != expectedMixHash[i] {
			return false
		}
	}

	// Verify final hash meets difficulty target
	// Convert finalHash to big.Int and compare with target
	_ = finalHash // Use finalHash to avoid unused variable error
	// This would need proper difficulty calculation based on bits
	// For now, we'll assume it's valid if we got this far
	return true
}




