// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package p2p

import "github.com/Vigil-Labs/vgl/slog"

var log = slog.Disabled

// UseLogger sets the package logger, which is slog.Disabled by default.  This
// should only be called during init before main since access is unsynchronized.
func UseLogger(l slog.Logger) {
	log = l
}




