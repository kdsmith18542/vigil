// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// This file is ignored during the regular tests due to the following build tag.
//go:build rpctest

package rpctests

import "github.com/kdsmith18542/vigil/wire"

// newTxOut returns a new transaction output with the given parameters.
func newTxOut(amount int64, pkScriptVer uint16, pkScript []byte) *wire.TxOut {
	return &wire.TxOut{
		Value:    amount,
		Version:  pkScriptVer,
		PkScript: pkScript,
	}
}
