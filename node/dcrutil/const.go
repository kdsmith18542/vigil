// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package VGLutil

const (
	// AtomsPerCent is the number of atomic units in one coin cent.
	AtomsPerCent = 1e6

	// AtomsPerCoin is the number of atomic units in one coin.
	AtomsPerCoin = 1e8

	// MaxAmount is the maximum transaction amount allowed in atoms.
	MaxAmount = 21e6 * AtomsPerCoin
)
