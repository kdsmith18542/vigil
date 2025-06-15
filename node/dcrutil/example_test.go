// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package VGLutil_test

import (
	"fmt"
	"math"

	"github.com/Vigil-Labs/vgl/VGLutil"
)

func ExampleAmount() {
	a := VGLutil.Amount(0)
	fmt.Println("Zero Atom:", a)

	a = VGLutil.Amount(1e8)
	fmt.Println("100,000,000 Atoms:", a)

	a = VGLutil.Amount(1e5)
	fmt.Println("100,000 Atoms:", a)
	// Output:
	// Zero Atom: 0 VGL
	// 100,000,000 Atoms: 1 VGL
	// 100,000 Atoms: 0.001 VGL
}

func ExampleNewAmount() {
	amountOne, err := VGLutil.NewAmount(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := VGLutil.NewAmount(0.01234567)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := VGLutil.NewAmount(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := VGLutil.NewAmount(math.NaN())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 VGL
	// 0.01234567 VGL
	// 0 VGL
	// invalid coin amount
}

func ExampleAmount_unitConversions() {
	amount := VGLutil.Amount(44433322211100)

	fmt.Println("Atom to kCoin:", amount.Format(VGLutil.AmountKiloCoin))
	fmt.Println("Atom to Coin:", amount)
	fmt.Println("Atom to MilliCoin:", amount.Format(VGLutil.AmountMilliCoin))
	fmt.Println("Atom to MicroCoin:", amount.Format(VGLutil.AmountMicroCoin))
	fmt.Println("Atom to Atom:", amount.Format(VGLutil.AmountAtom))

	// Output:
	// Atom to kCoin: 444.333222111 kVGL
	// Atom to Coin: 444333.222111 VGL
	// Atom to MilliCoin: 444333222.111 mVGL
	// Atom to MicroCoin: 444333222111 μVGL
	// Atom to Atom: 44433322211100 Atom
}




