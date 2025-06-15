// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// This file exists to prevent go mod tidy from removing requires on tools.
//
// It is excluded from the build to avoid including unused modules in the final
// binary.

//go:build require
// +build require

package main

import (
	_ "github.com/kdsmith18542/vigil/bech32"
	_ "github.com/kdsmith18542/vigil/mixing/mixclient"
)
