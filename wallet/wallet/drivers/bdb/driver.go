// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Package bdb registers the bdb driver at init time.  Importing bdb allows the
// wallet.OpenDB and wallet.CreateDB functions to be called with the following
// arguments:
//
//	var filename string
//	db, err := wallet.CreateDB("bdb", filename)
//	if err != nil { /* handle error */ }
//	db, err = wallet.OpenDB("bdb", filename)
//	if err != nil { /* handle error */ }
package bdb

import _ "github.com/kdsmith18542/vigil/wallet/wallet/internal/bdb" // Register bdb driver during init
