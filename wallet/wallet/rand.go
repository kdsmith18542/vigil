// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wallet

import "github.com/kdsmith18542/vigil/crypto/rand"

// Shuffle cryptographically shuffles a total of n items.
func Shuffle(n int, swap func(i, j int)) {
	rand.Shuffle(n, swap)
}
