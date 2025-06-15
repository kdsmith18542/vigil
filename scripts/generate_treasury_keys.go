// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// This script generates a new treasury address and outputs the corresponding
// P2SH script that should be set in the mainnet parameters.

package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/Vigil-Labs/vgl/chaincfg"
	"github.com/Vigil-Labs/vgl/chaincfg/chainec"
	"github.com/Vigil-Labs/vgl/chaincfg/chainhash"
	"github.com/Vigil-Labs/vgl/txscript"
	"github.com/Vigil-Labs/vgl/util"
)

func main() {
	// Use the mainnet parameters
	params := chaincfg.MainNetParams()

	// Generate a new secp256k1 private key
	privKey, err := chainec.Secp256k1.PrivateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the public key from the private key
	pubKey := privKey.Public()
	pubKeyBytes := pubKey.(*chainec.PublicKey).SerializeCompressed()

	// Create a pay-to-pubkey-hash address from the public key
	addr, err := util.NewAddressSecpPubKey(pubKeyBytes, params)
	if err != nil {
		log.Fatalf("Failed to create address: %v", err)
	}

	// Create a P2SH script for the treasury
	treasuryScript, err := txscript.PayToAddrScript(addr.AddressPubKeyHash())
	if err != nil {
		log.Fatalf("Failed to create P2SH script: %v", err)
	}

	// Output the results
	fmt.Println("TREASURY KEY GENERATION")
	fmt.Println("=======================")
	fmt.Printf("Network:            %s\n", params.Name)
	fmt.Printf("Private Key (WIF):  %s\n", privKey.String())
	fmt.Printf("Public Key (hex):   %x\n", pubKeyBytes)
	fmt.Printf("Address:            %s\n", addr.Address())
	fmt.Printf("P2SH Script (hex):  %x\n", treasuryScript)
	fmt.Println("\nIMPORTANT: Keep the private key secure and do not share it!")
	fmt.Println("This key will control the treasury funds.")

	// Format the script for the mainnet parameters
	fmt.Println("\nAdd the following to mainnetparams.go:")
	fmt.Printf("OrganizationPkScript:        %#v,\n", treasuryScript)
	fmt.Println("OrganizationPkScriptVersion: 0,")
}




