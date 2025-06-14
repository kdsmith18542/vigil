// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package addrmgr

import (
	"testing"

	"github.com/vigilnetwork/vgl/slog"
)

var testLogger = slog.NewBackend(&testWriter{}).Logger("TEST")

type testWriter struct{}

// Required to create a Write function for the testWriter
func (tw *testWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestUseLogger(t *testing.T) {
	UseLogger(testLogger)

	if log != testLogger {
		t.Errorf("Expected log to be set to testLogger, got %v", log)
	}
}
