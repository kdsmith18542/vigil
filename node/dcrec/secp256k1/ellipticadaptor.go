// Copyright 2020-2022 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package secp256k1

// References:
//   [SECG]: Recommended Elliptic Curve Domain Parameters
//     https://www.secg.org/sec2-v2.pdf
//
//   [GECC]: Guide to Elliptic Curve Cryptography (Hankerson, Menezes, Vanstone)

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
)

// CurveParams contains the parameters for the secp256k1 curve.
type CurveParams struct {
	// P is the prime used in the secp256k1 field.
	P *big.Int

	// N is the order of the secp256k1 curve group generated by the base point.
	N *big.Int

	// Gx and Gy are the x and y coordinate of the base point, respectively.
	Gx, Gy *big.Int

	// BitSize is the size of the underlying secp256k1 field in bits.
	BitSize int

	// H is the cofactor of the secp256k1 curve.
	H int

	// ByteSize is simply the bit size / 8 and is provided for convenience
	// since it is calculated repeatedly.
	ByteSize int
}

// Curve parameters taken from [SECG] section 2.4.1.
var curveParams = CurveParams{
	P:        fromHex("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"),
	N:        fromHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"),
	Gx:       fromHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"),
	Gy:       fromHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"),
	BitSize:  256,
	H:        1,
	ByteSize: 256 / 8,
}

// Params returns the secp256k1 curve parameters for convenience.
func Params() *CurveParams {
	return &curveParams
}

// KoblitzCurve provides an implementation for secp256k1 that fits the ECC Curve
// interface from crypto/elliptic.
type KoblitzCurve struct {
	*elliptic.CurveParams
}

// bigAffineToJacobian takes an affine point (x, y) as big integers and converts
// it to Jacobian point with Z=1.
func bigAffineToJacobian(x, y *big.Int, result *JacobianPoint) {
	result.X.SetByteSlice(x.Bytes())
	result.Y.SetByteSlice(y.Bytes())
	result.Z.SetInt(1)
}

// jacobianToBigAffine takes a Jacobian point (x, y, z) as field values and
// converts it to an affine point as big integers.
func jacobianToBigAffine(point *JacobianPoint) (*big.Int, *big.Int) {
	point.ToAffine()

	// Convert the field values for the now affine point to big.Ints.
	x3, y3 := new(big.Int), new(big.Int)
	x3.SetBytes(point.X.Bytes()[:])
	y3.SetBytes(point.Y.Bytes()[:])
	return x3, y3
}

// Params returns the parameters for the curve.
//
// This is part of the elliptic.Curve interface implementation.
func (curve *KoblitzCurve) Params() *elliptic.CurveParams {
	return curve.CurveParams
}

// IsOnCurve returns whether or not the affine point (x,y) is on the curve.
//
// This is part of the elliptic.Curve interface implementation.  This function
// differs from the crypto/elliptic algorithm since a = 0 not -3.
func (curve *KoblitzCurve) IsOnCurve(x, y *big.Int) bool {
	// Convert big ints to a Jacobian point for faster arithmetic.
	var point JacobianPoint
	bigAffineToJacobian(x, y, &point)
	return isOnCurve(&point.X, &point.Y)
}

// Add returns the sum of (x1,y1) and (x2,y2).
//
// This is part of the elliptic.Curve interface implementation.
func (curve *KoblitzCurve) Add(x1, y1, x2, y2 *big.Int) (*big.Int, *big.Int) {
	// The point at infinity is the identity according to the group law for
	// elliptic curve cryptography.  Thus, ∞ + P = P and P + ∞ = P.
	if x1.Sign() == 0 && y1.Sign() == 0 {
		return x2, y2
	}
	if x2.Sign() == 0 && y2.Sign() == 0 {
		return x1, y1
	}

	// Convert the affine coordinates from big integers to Jacobian points,
	// do the point addition in Jacobian projective space, and convert the
	// Jacobian point back to affine big.Ints.
	var p1, p2, result JacobianPoint
	bigAffineToJacobian(x1, y1, &p1)
	bigAffineToJacobian(x2, y2, &p2)
	AddNonConst(&p1, &p2, &result)
	return jacobianToBigAffine(&result)
}

// Double returns 2*(x1,y1).
//
// This is part of the elliptic.Curve interface implementation.
func (curve *KoblitzCurve) Double(x1, y1 *big.Int) (*big.Int, *big.Int) {
	if y1.Sign() == 0 {
		return new(big.Int), new(big.Int)
	}

	// Convert the affine coordinates from big integers to Jacobian points,
	// do the point doubling in Jacobian projective space, and convert the
	// Jacobian point back to affine big.Ints.
	var point, result JacobianPoint
	bigAffineToJacobian(x1, y1, &point)
	DoubleNonConst(&point, &result)
	return jacobianToBigAffine(&result)
}

// moduloReduce reduces k from more than 32 bytes to 32 bytes and under.  This
// is done by doing a simple modulo curve.N.  We can do this since G^N = 1 and
// thus any other valid point on the elliptic curve has the same order.
func moduloReduce(k []byte) []byte {
	// Since the order of G is curve.N, we can use a much smaller number by
	// doing modulo curve.N
	if len(k) > curveParams.ByteSize {
		tmpK := new(big.Int).SetBytes(k)
		tmpK.Mod(tmpK, curveParams.N)
		return tmpK.Bytes()
	}

	return k
}

// ScalarMult returns k*(bx, by) where k is a big endian integer.
//
// This is part of the elliptic.Curve interface implementation.
func (curve *KoblitzCurve) ScalarMult(bx, by *big.Int, k []byte) (*big.Int, *big.Int) {
	// Convert the affine coordinates from big integers to Jacobian points,
	// do the multiplication in Jacobian projective space, and convert the
	// Jacobian point back to affine big.Ints.
	var kModN ModNScalar
	kModN.SetByteSlice(moduloReduce(k))
	var point, result JacobianPoint
	bigAffineToJacobian(bx, by, &point)
	ScalarMultNonConst(&kModN, &point, &result)
	return jacobianToBigAffine(&result)
}

// ScalarBaseMult returns k*G where G is the base point of the group and k is a
// big endian integer.
//
// This is part of the elliptic.Curve interface implementation.
func (curve *KoblitzCurve) ScalarBaseMult(k []byte) (*big.Int, *big.Int) {
	// Perform the multiplication and convert the Jacobian point back to affine
	// big.Ints.
	var kModN ModNScalar
	kModN.SetByteSlice(moduloReduce(k))
	var result JacobianPoint
	ScalarBaseMultNonConst(&kModN, &result)
	return jacobianToBigAffine(&result)
}

// X returns the x coordinate of the public key.
func (p *PublicKey) X() *big.Int {
	return new(big.Int).SetBytes(p.x.Bytes()[:])
}

// Y returns the y coordinate of the public key.
func (p *PublicKey) Y() *big.Int {
	return new(big.Int).SetBytes(p.y.Bytes()[:])
}

// ToECDSA returns the public key as a *ecdsa.PublicKey.
func (p *PublicKey) ToECDSA() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: S256(),
		X:     p.X(),
		Y:     p.Y(),
	}
}

// ToECDSA returns the private key as a *ecdsa.PrivateKey.
func (p *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	var privKeyBytes [PrivKeyBytesLen]byte
	p.Key.PutBytes(&privKeyBytes)
	var result JacobianPoint
	ScalarBaseMultNonConst(&p.Key, &result)
	x, y := jacobianToBigAffine(&result)
	newPrivKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: S256(),
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(privKeyBytes[:]),
	}
	zeroArray32(&privKeyBytes)
	return newPrivKey
}

// fromHex converts the passed hex string into a big integer pointer and will
// panic is there is an error.  This is only provided for the hard-coded
// constants so errors in the source code can bet detected. It will only (and
// must only) be called for initialization purposes.
func fromHex(s string) *big.Int {
	if s == "" {
		return big.NewInt(0)
	}
	r, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("invalid hex in source file: " + s)
	}
	return r
}

// secp256k1 is a global instance of the KoblitzCurve implementation which in
// turn embeds and implements elliptic.CurveParams.
var secp256k1 = &KoblitzCurve{
	CurveParams: &elliptic.CurveParams{
		P:       curveParams.P,
		N:       curveParams.N,
		B:       fromHex("0000000000000000000000000000000000000000000000000000000000000007"),
		Gx:      curveParams.Gx,
		Gy:      curveParams.Gy,
		BitSize: curveParams.BitSize,
		Name:    "secp256k1",
	},
}

// S256 returns an elliptic.Curve which implements secp256k1.
func S256() *KoblitzCurve {
	return secp256k1
}
