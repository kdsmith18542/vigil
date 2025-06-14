// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaingen_test

import (
	"fmt"

	"github.com/kdsmith18542/vigil/blockchain/v5/chaingen"
	"github.com/kdsmith18542/vigil/chaincfg/v3"
)

// This example demonstrates creating a new generator instance and using it to
// generate the required first block and enough blocks to have mature coinbase
// outputs to work with along with asserting the generator state along the way.
func Example_basicUsage() {
	params := chaincfg.RegNetParams()
	g, err := chaingen.MakeGenerator(params)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Shorter versions of useful params for convenience.
	coinbaseMaturity := params.CoinbaseMaturity

	// ---------------------------------------------------------------------
	// First block.
	// ---------------------------------------------------------------------

	// Add the required first block.
	//
	//   genesis -> bfb
	g.CreateBlockOne("bfb", 0)
	g.AssertTipHeight(1)
	fmt.Println(g.TipName())

	// ---------------------------------------------------------------------
	// Generate enough blocks to have mature coinbase outputs to work with.
	//
	//   genesis -> bp -> bm0 -> bm1 -> ... -> bm#
	// ---------------------------------------------------------------------

	for i := uint16(0); i < coinbaseMaturity; i++ {
		blockName := fmt.Sprintf("bm%d", i)
		g.NextBlock(blockName, nil, nil)
		g.SaveTipCoinbaseOuts()
		fmt.Println(g.TipName())
	}
	g.AssertTipHeight(uint32(coinbaseMaturity) + 1)

	// Output:
	// bfb
	// bm0
	// bm1
	// bm2
	// bm3
	// bm4
	// bm5
	// bm6
	// bm7
	// bm8
	// bm9
	// bm10
	// bm11
	// bm12
	// bm13
	// bm14
	// bm15
}
