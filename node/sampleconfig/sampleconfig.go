// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package sampleconfig

import (
	_ "embed"
)

// samplevgldConf is a string containing the commented example config for vgld.
//
//go:embed sample-vgld.conf
var samplevgldConf string

// samplevglctlConf is a string containing the commented example config for
// vglctl.
//
//go:embed sample-vglctl.conf
var samplevglctlConf string

// vgld returns a string containing the commented example config for vgld.
func vgld() string {
	return samplevgldConf
}

// FileContents returns a string containing the commented example config for
// vgld.
//
// Deprecated: Use the [vgld] function instead.
func FileContents() string {
	return vgld()
}

// vglctl returns a string containing the commented example config for vglctl.
func vglctl() string {
	return samplevglctlConf
}

// vglctlSampleConfig is a string containing the commented example config for
// vglctl.
//
// Deprecated: Use the [vglctl] function instead.
var vglctlSampleConfig = vglctl()
