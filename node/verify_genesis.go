//go:build ignore
// +build ignore

// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	"github.com/kdsmith18542/vigil/chaincfg"
)

func main() {
	fmt.Println("Verifying Vigil genesis blocks...\n")

	// Check mainnet genesis block
	mainnetHash := chaincfg.MainNetParams.GenesisHash
	wantMainnetHash, _ := chainhash.NewHashFromStr("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	mainnetMerkle := chaincfg.MainNetParams.GenesisBlock.Header.MerkleRoot
	wantMainnetMerkle, _ := chainhash.NewHashFromStr("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b")

	fmt.Println("=== Mainnet ===")
	fmt.Printf("Genesis Hash: %s\n", mainnetHash)
	fmt.Printf("Expected:     %s\n", wantMainnetHash)
	fmt.Printf("Match:        %v\n", mainnetHash.IsEqual(wantMainnetHash))
	fmt.Printf("Merkle Root:  %s\n", mainnetMerkle)
	fmt.Printf("Expected:     %s\n", wantMainnetMerkle)
	fmt.Printf("Match:        %v\n\n", bytes.Equal(mainnetMerkle[:], wantMainnetMerkle[:]))

	// Check testnet3 genesis block
	testnet3Hash := chaincfg.TestNet3Params.GenesisHash
	wantTestnet3Hash, _ := chainhash.NewHashFromStr("000000000933ea01ad0ee984209779baaec3ced90fa3f408719526f8d77f4943")
	testnet3Merkle := chaincfg.TestNet3Params.GenesisBlock.Header.MerkleRoot
	wantTestnet3Merkle, _ := chainhash.NewHashFromStr("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b")

	fmt.Println("=== Testnet3 ===")
	fmt.Printf("Genesis Hash: %s\n", testnet3Hash)
	fmt.Printf("Expected:     %s\n", wantTestnet3Hash)
	fmt.Printf("Match:        %v\n", testnet3Hash.IsEqual(wantTestnet3Hash))
	fmt.Printf("Merkle Root:  %s\n", testnet3Merkle)
	fmt.Printf("Expected:     %s\n", wantTestnet3Merkle)
	fmt.Printf("Match:        %v\n", bytes.Equal(testnet3Merkle[:], wantTestnet3Merkle[:]))

	// Exit with error if any check fails
	if !mainnetHash.IsEqual(wantMainnetHash) || !bytes.Equal(mainnetMerkle[:], wantMainnetMerkle[:]) ||
		!testnet3Hash.IsEqual(wantTestnet3Hash) || !bytes.Equal(testnet3Merkle[:], wantTestnet3Merkle[:]) {
		os.Exit(1)
	}
}
