// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//go:build tinygo

package secp256k1

// This file contains the variants suitable for
// memory or storage constrained environments.

func scalarBaseMultNonConst(k *ModNScalar, result *JacobianPoint) {
	scalarBaseMultNonConstSlow(k, result)
}
