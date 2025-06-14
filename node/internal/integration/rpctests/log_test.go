// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//go:build rpctest

package rpctests

import (
	"testing"

	"github.com/Vigil/VGLtest/vgldtest"
	"github.com/vigilnetwork/vgl/slog"
)

type testLog struct {
	*testing.T
}

func (t *testLog) Write(b []byte) (int, error) {
	t.Logf("%s", b)
	return len(b), nil
}

// useTestLogger sets the vgldtest package-level logger to a backend that
// writes trace-level logs to the test log.  A function is returned to set the
// logger back to Disabled when finished.
//
// Due to vgldtest's use of a global logger variable that must write to test
// logs to individual test variables, it is not possible to parallelize tests.
func useTestLogger(t *testing.T) func() {
	backend := slog.NewBackend(&testLog{T: t})
	l := backend.Logger("TEST")
	l.SetLevel(slog.LevelTrace)
	vgldtest.UseLogger(l)
	return func() {
		vgldtest.UseLogger(slog.Disabled)
	}
}
