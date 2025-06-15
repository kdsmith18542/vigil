// Script to generate a new treasury key pair for Vigil
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

	// Create a hash of the public key (RIPEMD160(SHA256(pubkey)))
	sha := sha256.Sum256(compressedPubKey)
	ripe := ripemd160(sha[:])

	// Create a P2SH script (OP_HASH160 <hash> OP_EQUAL)
	script := make([]byte, 23)
	script[0] = 0xa9       // OP_HASH160
	script[1] = 0x14       // Push 20 bytes
	copy(script[2:22], ripe) // The hash
	script[22] = 0x87      // OP_EQUAL

	// Output the results
	fmt.Println("Vigil Treasury Key Generation")
	fmt.Println("==============================")
	fmt.Printf("Private Key (hex): %x\n", privateKey.D.Bytes())
	fmt.Printf("Public Key (hex):  %x\n", compressedPubKey)
	fmt.Printf("P2SH Script:       %x\n", script)
	fmt.Println("\nAdd the following to mainnetparams.go:")
	fmt.Printf("OrganizationPkScript:        %#v,\n", script)
	fmt.Println("OrganizationPkScriptVersion: 0,")
}

// Simple RIPEMD-160 implementation
func ripemd160(data []byte) []byte {
	h := NewRipemd160()
	h.Write(data)
	return h.Sum(nil)
}

// RIPEMD-160 implementation
type ripemd160Digest struct {
	s [5]uint32
	x [16]uint32
	nx int
	len uint64
}

func NewRipemd160() *ripemd160Digest {
	r := new(ripemd160Digest)
	r.Reset()
	return r
}

func (r *ripemd160Digest) Reset() {
	r.s = [5]uint32{0x67452301, 0xefcdab89, 0x98badcfe, 0x10325476, 0xc3d2e1f0}
	r.nx = 0
	r.len = 0
}

func (r *ripemd160Digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	r.len += uint64(nn)
	if r.nx > 0 {
		n := copy(r.x[r.nx:], p)
		r.nx += n
		if r.nx == 16 {
			block(r, r.x[:])
			r.nx = 0
		}
		p = p[n:]
	}
	for len(p) >= 16 {
		block(r, p[:16])
		p = p[16:]
	}
	if len(p) > 0 {
		r.nx = copy(r.x[:], p)
	}
	return
}

func (r *ripemd160Digest) Sum(in []byte) []byte {
	r0 := *r
	hash := r0.checkSum()
	return append(in, hash[:]...)
}

func (r *ripemd160Digest) checkSum() [20]byte {
	// Append padding
	tmp := [64]byte{0x80}
	if r.nx < 56 {
		r.Write(tmp[0 : 56-r.nx])
	} else {
		r.Write(tmp[0 : 64+56-r.nx])
	}

	// Append length in bits
	r.len <<= 3
	for i := uint(0); i < 8; i++ {
		tmp[i] = byte(r.len >> (8 * i))
	}
	r.Write(tmp[0:8])

	if r.nx != 0 {
		panic("r.nx != 0")
	}

	var digest [20]byte
	for i, s := range r.s {
		digest[i*4] = byte(s)
		digest[i*4+1] = byte(s >> 8)
		digest[i*4+2] = byte(s >> 16)
		digest[i*4+3] = byte(s >> 24)
	}
	return digest
}

// Block processing function
func block(dig *ripemd160Digest, p []byte) {
	var x [16]uint32
	j := 0
	for i := 0; i < 16; i++ {
		x[i] = uint32(p[j]) | uint32(p[j+1])<<8 | uint32(p[j+2])<<16 | uint32(p[j+3])<<24
		j += 4
	}

	a, b, c, d, e := dig.s[0], dig.s[1], dig.s[2], dig.s[3], dig.s[4]

	// Round 1
	a = rol(a+f1(b, c, d)+x[0], 11) + e
	c = rol(c, 10)
	e = rol(e+f1(a, b, c)+x[1], 14) + d
	b = rol(b, 10)
	d = rol(d+f1(e, a, b)+x[2], 15) + c
	a = rol(a, 10)
	c = rol(c+f1(d, e, a)+x[3], 12) + b
	e = rol(e, 10)
	b = rol(b+f1(c, d, e)+x[4], 5) + a
	d = rol(d, 10)
	a = rol(a+f1(b, c, d)+x[5], 8) + e
	c = rol(c, 10)
	e = rol(e+f1(a, b, c)+x[6], 7) + d
	b = rol(b, 10)
	d = rol(d+f1(e, a, b)+x[7], 9) + c
	a = rol(a, 10)
	c = rol(c+f1(d, e, a)+x[8], 11) + b
	e = rol(e, 10)
	b = rol(b+f1(c, d, e)+x[9], 13) + a
	d = rol(d, 10)
	a = rol(a+f1(b, c, d)+x[10], 14) + e
	c = rol(c, 10)
	e = rol(e+f1(a, b, c)+x[11], 15) + d
	b = rol(b, 10)
	d = rol(d+f1(e, a, b)+x[12], 6) + c
	a = rol(a, 10)
	c = rol(c+f1(d, e, a)+x[13], 7) + b
	e = rol(e, 10)
	b = rol(b+f1(c, d, e)+x[14], 9) + a
	d = rol(d, 10)
	a = rol(a+f1(b, c, d)+x[15], 8) + e
	c = rol(c, 10)

	// Round 2
	e = rol(e+f2(a, b, c)+x[7]+0x5a827999, 7) + d
	b = rol(b, 10)
	d = rol(d+f2(e, a, b)+x[4]+0x5a827999, 6) + c
	a = rol(a, 10)
	c = rol(c+f2(d, e, a)+x[13]+0x5a827999, 8) + b
	e = rol(e, 10)
	b = rol(b+f2(c, d, e)+x[1]+0x5a827999, 13) + a
	d = rol(d, 10)
	a = rol(a+f2(b, c, d)+x[10]+0x5a827999, 11) + e
	c = rol(c, 10)
	e = rol(e+f2(a, b, c)+x[6]+0x5a827999, 9) + d
	b = rol(b, 10)
	d = rol(d+f2(e, a, b)+x[15]+0x5a827999, 7) + c
	a = rol(a, 10)
	c = rol(c+f2(d, e, a)+x[3]+0x5a827999, 15) + b
	e = rol(e, 10)
	b = rol(b+f2(c, d, e)+x[12]+0x5a827999, 7) + a
	d = rol(d, 10)
	a = rol(a+f2(b, c, d)+x[0]+0x5a827999, 12) + e
	c = rol(c, 10)
	e = rol(e+f2(a, b, c)+x[9]+0x5a827999, 15) + d
	b = rol(b, 10)
	d = rol(d+f2(e, a, b)+x[5]+0x5a827999, 9) + c
	a = rol(a, 10)
	c = rol(c+f2(d, e, a)+x[2]+0x5a827999, 11) + b
	e = rol(e, 10)
	b = rol(b+f2(c, d, e)+x[14]+0x5a827999, 7) + a
	d = rol(d, 10)
	a = rol(a+f2(b, c, d)+x[11]+0x5a827999, 13) + e
	c = rol(c, 10)
	e = rol(e+f2(a, b, c)+x[8]+0x5a827999, 12) + d
	b = rol(b, 10)

	// Round 3
	d = rol(d+f3(e, a, b)+x[3]+0x6ed9eba1, 11) + c
	a = rol(a, 10)
	c = rol(c+f3(d, e, a)+x[10]+0x6ed9eba1, 13) + b
	e = rol(e, 10)
	b = rol(b+f3(c, d, e)+x[14]+0x6ed9eba1, 6) + a
	d = rol(d, 10)
	a = rol(a+f3(b, c, d)+x[4]+0x6ed9eba1, 7) + e
	c = rol(c, 10)
	e = rol(e+f3(a, b, c)+x[9]+0x6ed9eba1, 14) + d
	b = rol(b, 10)
	d = rol(d+f3(e, a, b)+x[15]+0x6ed9eba1, 9) + c
	a = rol(a, 10)
	c = rol(c+f3(d, e, a)+x[8]+0x6ed9eba1, 13) + b
	e = rol(e, 10)
	b = rol(b+f3(c, d, e)+x[1]+0x6ed9eba1, 15) + a
	d = rol(d, 10)
	a = rol(a+f3(b, c, d)+x[2]+0x6ed9eba1, 14) + e
	c = rol(c, 10)
	e = rol(e+f3(a, b, c)+x[7]+0x6ed9eba1, 8) + d
	b = rol(b, 10)
	d = rol(d+f3(e, a, b)+x[0]+0x6ed9eba1, 13) + c
	a = rol(a, 10)
	c = rol(c+f3(d, e, a)+x[6]+0x6ed9eba1, 6) + b
	e = rol(e, 10)
	b = rol(b+f3(c, d, e)+x[13]+0x6ed9eba1, 5) + a
	d = rol(d, 10)
	a = rol(a+f3(b, c, d)+x[11]+0x6ed9eba1, 12) + e
	c = rol(c, 10)
	e = rol(e+f3(a, b, c)+x[5]+0x6ed9eba1, 7) + d
	b = rol(b, 10)
	d = rol(d+f3(e, a, b)+x[12]+0x6ed9eba1, 5) + c
	a = rol(a, 10)

	// Round 4
	c = rol(c+f4(d, e, a)+x[1]+0x8f1bbcdc, 11) + b
	e = rol(e, 10)
	b = rol(b+f4(c, d, e)+x[9]+0x8f1bbcdc, 12) + a
	d = rol(d, 10)
	a = rol(a+f4(b, c, d)+x[11]+0x8f1bbcdc, 14) + e
	c = rol(c, 10)
	e = rol(e+f4(a, b, c)+x[10]+0x8f1bbcdc, 15) + d
	b = rol(b, 10)
	d = rol(d+f4(e, a, b)+x[0]+0x8f1bbcdc, 14) + c
	a = rol(a, 10)
	c = rol(c+f4(d, e, a)+x[8]+0x8f1bbcdc, 15) + b
	e = rol(e, 10)
	b = rol(b+f4(c, d, e)+x[12]+0x8f1bbcdc, 9) + a
	d = rol(d, 10)
	a = rol(a+f4(b, c, d)+x[4]+0x8f1bbcdc, 8) + e
	c = rol(c, 10)
	e = rol(e+f4(a, b, c)+x[13]+0x8f1bbcdc, 9) + d
	b = rol(b, 10)
	d = rol(d+f4(e, a, b)+x[3]+0x8f1bbcdc, 14) + c
	a = rol(a, 10)
	c = rol(c+f4(d, e, a)+x[7]+0x8f1bbcdc, 5) + b
	e = rol(e, 10)
	b = rol(b+f4(c, d, e)+x[15]+0x8f1bbcdc, 6) + a
	d = rol(d, 10)
	a = rol(a+f4(b, c, d)+x[14]+0x8f1bbcdc, 8) + e
	c = rol(c, 10)
	e = rol(e+f4(a, b, c)+x[5]+0x8f1bbcdc, 6) + d
	b = rol(b, 10)
	d = rol(d+f4(e, a, b)+x[6]+0x8f1bbcdc, 5) + c
	a = rol(a, 10)
	c = rol(c+f4(d, e, a)+x[2]+0x8f1bbcdc, 12) + b
	e = rol(e, 10)

	// Round 5
	b = rol(b+f5(c, d, e)+x[4]+0xa953fd4e, 9) + a
	d = rol(d, 10)
	a = rol(a+f5(b, c, d)+x[0]+0xa953fd4e, 15) + e
	c = rol(c, 10)
	e = rol(e+f5(a, b, c)+x[5]+0xa953fd4e, 5) + d
	b = rol(b, 10)
	d = rol(d+f5(e, a, b)+x[9]+0xa953fd4e, 11) + c
	a = rol(a, 10)
	c = rol(c+f5(d, e, a)+x[7]+0xa953fd4e, 6) + b
	e = rol(e, 10)
	b = rol(b+f5(c, d, e)+x[12]+0xa953fd4e, 8) + a
	d = rol(d, 10)
	a = rol(a+f5(b, c, d)+x[2]+0xa953fd4e, 13) + e
	c = rol(c, 10)
	e = rol(e+f5(a, b, c)+x[10]+0xa953fd4e, 12) + d
	b = rol(b, 10)
	d = rol(d+f5(e, a, b)+x[14]+0xa953fd4e, 5) + c
	a = rol(a, 10)
	c = rol(c+f5(d, e, a)+x[1]+0xa953fd4e, 12) + b
	e = rol(e, 10)
	b = rol(b+f5(c, d, e)+x[3]+0xa953fd4e, 13) + a
	d = rol(d, 10)
	a = rol(a+f5(b, c, d)+x[8]+0xa953fd4e, 14) + e
	c = rol(c, 10)
	e = rol(e+f5(a, b, c)+x[11]+0xa953fd4e, 11) + d
	b = rol(b, 10)
	d = rol(d+f5(e, a, b)+x[6]+0xa953fd4e, 8) + c
	a = rol(a, 10)
	c = rol(c+f5(d, e, a)+x[15]+0xa953fd4e, 5) + b
	e = rol(e, 10)
	b = rol(b+f5(c, d, e)+x[13]+0xa953fd4e, 6) + a
	d = rol(d, 10)

	dig.s[1] += c + e
	dig.s[2] += d + a
	dig.s[3] += e + b
	dig.s[4] += a + c
	dig.s[0] += b + d
}

// Helper functions for RIPEMD-160
func rol(x uint32, n uint) uint32 {
	return (x << n) | (x >> (32 - n))
}

func f1(x, y, z uint32) uint32 { return x ^ y ^ z }
func f2(x, y, z uint32) uint32 { return (x & y) | (^x & z) }
func f3(x, y, z uint32) uint32 { return (x | ^y) ^ z }
func f4(x, y, z uint32) uint32 { return (x & z) | (y & ^z) }
func f5(x, y, z uint32) uint32 { return x ^ (y | ^z) }




