// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package main

import (
	"os"
	"syscall"
)

func init() {
	signals = []os.Signal{
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
	}
}
