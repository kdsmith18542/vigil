// Copyright (c) 2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package edwards

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestGenerateSharedSecret(t *testing.T) {
	privKey1, err := GeneratePrivateKey()
	if err != nil {
		t.Errorf("private key generation error: %s", err)
		return
	}
	privKey2, err := GeneratePrivateKey()
	if err != nil {
		t.Errorf("private key generation error: %s", err)
		return
	}

	pk1x, pk1y := privKey1.Public()
	pk1 := NewPublicKey(pk1x, pk1y)
	pk2x, pk2y := privKey2.Public()
	pk2 := NewPublicKey(pk2x, pk2y)
	secret1 := GenerateSharedSecret(privKey1, pk2)
	secret2 := GenerateSharedSecret(privKey2, pk1)

	if !bytes.Equal(secret1, secret2) {
		t.Errorf("ECDH failed, secrets mismatch - first: %x, second: %x",
			secret1, secret2)
	}
}

// Test 1: Encryption and decryption.
func TestCipheringBasic(t *testing.T) {
	privkey, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal("failed to generate private key")
	}

	in := []byte("Hey there dude. How are you doing? This is a test.")

	pk1x, pk1y := privkey.Public()
	pk1 := NewPublicKey(pk1x, pk1y)
	out, err := Encrypt(pk1, in)
	if err != nil {
		t.Fatal("failed to encrypt:", err)
	}

	dec, err := Decrypt(privkey, out)
	if err != nil {
		t.Fatal("failed to decrypt:", err)
	}

	if !bytes.Equal(in, dec) {
		t.Error("decrypted data doesn't match original")
	}
}

func TestCiphering(t *testing.T) {
	c := Edwards()

	pb, _ := hex.DecodeString("fe38240982f313ae5afb3e904fb8215fb11af1200592b" +
		"fca26c96c4738e4bf8f")
	pbBig := new(big.Int).SetBytes(pb)
	pbBig.Mod(pbBig, c.N)
	pb = pbBig.Bytes()
	pb = copyBytes(pb)[:]
	privkey, pubkey, err := PrivKeyFromScalar(pb)
	if err != nil {
		t.Error(err)
	}
	in := []byte("This is just a test.")
	localOut, err := Encrypt(pubkey, in)
	if err != nil {
		t.Error(err)
	}

	out, _ := hex.DecodeString("1ffcb6f11fb9dc57222382019ae710b2ffff0020503f4" +
		"117665f80b226961a4a0c0ae229f3b914d43e36238be05b0799623ae6ea0209d3095" +
		"04f86635c50baca78d11189d4dc02c2f32c4c11e9d50b04eb2d3ff4b9f95e7f2e90e" +
		"0f4a8d64a2a4149c27d21f88f2dedc200f4b609936c0d67ca98")

	_, err = Decrypt(privkey, out)
	if err != nil {
		t.Fatal("failed to decrypt:", err)
	}

	dec, err := Decrypt(privkey, localOut)
	if err != nil {
		t.Fatal("failed to decrypt:", err)
	}

	if !bytes.Equal(in, dec) {
		t.Error("decrypted data doesn't match original")
	}
}

func TestCipheringErrors(t *testing.T) {
	privkey, err := GeneratePrivateKey()
	if err != nil {
		t.Fatal("failed to generate private key")
	}

	tests1 := []struct {
		ciphertext []byte // input ciphertext
	}{
		{bytes.Repeat([]byte{0x00}, 133)},                   // errInputTooShort
		{bytes.Repeat([]byte{0x00}, 134)},                   // errUnsupportedCurve
		{bytes.Repeat([]byte{0xFF, 0xFF}, 134)},             // errInvalidXLength
		{bytes.Repeat([]byte{0xFF, 0xFF, 0x00, 0x20}, 134)}, // errInvalidYLength
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // IV
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF,
			0x00, 0x20, // Y length
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Y
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ciphertext
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // MAC
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}}, // invalid pubkey
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // IV
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF,
			0x00, 0x20, // Y length
			0x7E, 0x76, 0xDC, 0x58, 0xF6, 0x93, 0xBD, 0x7E, // Y
			0x70, 0x10, 0x35, 0x8C, 0xE6, 0xB1, 0x65, 0xE4,
			0x83, 0xA2, 0x92, 0x10, 0x10, 0xDB, 0x67, 0xAC,
			0x11, 0xB1, 0xB5, 0x1B, 0x65, 0x19, 0x53, 0xD2,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ciphertext
			// padding not aligned to 16 bytes
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // MAC
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}}, // errInvalidPadding
		{[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // IV
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0xFF, 0xFF,
			0x00, 0x20, // Y length
			0x7E, 0x76, 0xDC, 0x58, 0xF6, 0x93, 0xBD, 0x7E, // Y
			0x70, 0x10, 0x35, 0x8C, 0xE6, 0xB1, 0x65, 0xE4,
			0x83, 0xA2, 0x92, 0x10, 0x10, 0xDB, 0x67, 0xAC,
			0x11, 0xB1, 0xB5, 0x1B, 0x65, 0x19, 0x53, 0xD2,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // ciphertext
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // MAC
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		}}, // ErrInvalidMAC
	}

	for i, test := range tests1 {
		_, err = Decrypt(privkey, test.ciphertext)
		if err == nil {
			t.Errorf("Decrypt #%d did not get error", i)
		}
	}

	// test error from removePKCSPadding
	tests2 := []struct {
		in []byte // input data
	}{
		{bytes.Repeat([]byte{0x11}, 17)},
		{bytes.Repeat([]byte{0x07}, 15)},
	}
	for i, test := range tests2 {
		_, err = TstRemovePKCSPadding(test.in)
		if err == nil {
			t.Errorf("removePKCSPadding #%d did not get error", i)
		}
	}
}

// TstRemovePKCSPadding makes the internal removePKCSPadding function available
// to the test package.
func TstRemovePKCSPadding(src []byte) ([]byte, error) {
	return removePKCSPadding(src)
}
