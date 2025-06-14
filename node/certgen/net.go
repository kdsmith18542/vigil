// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

//go:build !appengine
// +build !appengine

package certgen

import (
	"net"
)

// interfaceAddrs returns a list of the system's network interface addresses.
// It is wrapped here so that we can substitute it for other functions when
// building for systems that do not allow access to net.InterfaceAddrs().
func interfaceAddrs() ([]net.Addr, error) {
	return net.InterfaceAddrs()
}
