// Simple script to generate a new treasury key pair for Vigil
package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	// Generate a new private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// Get the public key in compressed format
	pubKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	compressedPubKey := make([]byte, 33)
	if pubKey[63]%2 == 0 {
		compressedPubKey[0] = 0x02 // Even Y value
	} else {
		compressedPubKey[0] = 0x03 // Odd Y value
	}
	copy(compressedPubKey[1:], pubKey[:32])

	// Create a SHA-256 hash of the public key
	sha256Hash := sha256.Sum256(compressedPubKey)
	
	// For demonstration, we'll create a simple P2SH script
	// In a real implementation, you would use a proper RIPEMD160 hash here
	// and create a proper P2SH script
	
	// Create a simple script (this is just an example - use proper P2SH in production)
	script := make([]byte, 23)
	script[0] = 0xa9       // OP_HASH160
	script[1] = 0x14       // Push 20 bytes
	// For demonstration, we'll use the first 20 bytes of the SHA-256 hash
	// In a real implementation, use RIPEMD160(SHA256(pubkey))
	copy(script[2:22], sha256Hash[:20])
	script[22] = 0x87      // OP_EQUAL

	// Output the results
	fmt.Println("Vigil Treasury Key Generation")
	fmt.Println("==============================")
	fmt.Printf("Private Key (hex): %x\n", privateKey.D.Bytes())
	fmt.Printf("Public Key (hex):  %x\n", compressedPubKey)
	fmt.Printf("P2SH Script:       %x\n", script)
	
	// Format for mainnetparams.go
	fmt.Println("\nAdd the following to mainnetparams.go:")
	fmt.Printf("OrganizationPkScript:        %#v,\n", script)
	fmt.Println("OrganizationPkScriptVersion: 0,")
}




