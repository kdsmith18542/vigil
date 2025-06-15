// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wallet

import (
	"flag"
	"os"
	"testing"

	"github.com/kdsmith18542/vigil/slog"
)

var logFlag = flag.Bool("log", false, "enable package logger")

func TestMain(m *testing.M) {
	flag.Parse()
	if *logFlag {
		UseLogger(slog.NewBackend(os.Stderr).Logger("WLLT"))
	}

	os.Exit(m.Run())
}
