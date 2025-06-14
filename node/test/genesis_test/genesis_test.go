// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package genesis_test

import (
	"bytes"
	"testing"

	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	)

// TestGenesisBlockHashes verifies the hashes of the genesis blocks.
func TestGenesisBlockHashes(t *testing.T) {
	// Mainnet genesis block hash
	mainnetHash := chaincfg.MainNetParams.GenesisHash
	wantMainnetHash, _ := chainhash.NewHashFromStr("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	if !mainnetHash.IsEqual(wantMainnetHash) {
		t.Errorf("mainnet genesis block hash does not match expected value - got %v, want %v",
			mainnetHash, wantMainnetHash)
	}

	// Testnet3 genesis block hash
	testnet3Hash := chaincfg.TestNet3Params.GenesisHash
	wantTestnet3Hash, _ := chainhash.NewHashFromStr("000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943")
	if !testnet3Hash.IsEqual(wantTestnet3Hash) {
		t.Errorf("testnet3 genesis block hash does not match expected value - got %v, want %v",
			testnet3Hash, wantTestnet3Hash)
	}
}

// TestGenesisBlockMerkleRoots verifies the merkle roots of the genesis blocks.
func TestGenesisBlockMerkleRoots(t *testing.T) {
	// Mainnet genesis merkle root
	mainnetMerkle := chaincfg.MainNetParams.GenesisBlock.Header.MerkleRoot
	wantMainnetMerkle, _ := chainhash.NewHashFromStr("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b")
	if !bytes.Equal(mainnetMerkle[:], wantMainnetMerkle[:]) {
		t.Errorf("mainnet genesis merkle root does not match expected value - got %v, want %v",
			mainnetMerkle, wantMainnetMerkle)
	}

	// Testnet3 genesis merkle root
	testnet3Merkle := chaincfg.TestNet3Params.GenesisBlock.Header.MerkleRoot
	wantTestnet3Merkle, _ := chainhash.NewHashFromStr("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b")
	if !bytes.Equal(testnet3Merkle[:], wantTestnet3Merkle[:]) {
		t.Errorf("testnet3 genesis merkle root does not match expected value - got %v, want %v",
			testnet3Merkle, wantTestnet3Merkle)
	}
}
