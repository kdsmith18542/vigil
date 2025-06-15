wire
====

[![Build Status](https://github.com/Vigil-Labs/vgl/workflows/Build%20and%20Test/badge.svg)](https://github.com/Vigil-Labs/vgl/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/Vigil-Labs/vgl/wire)

Package wire implements the Vigil wire protocol. A comprehensive suite of
tests with 100% test coverage is provided to ensure proper functionality.

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to interface with Vigil peers at the wire
protocol level.

## Installation and Updating

This package is part of the `github.com/Vigil-Labs/vgl/wire` module. Use the
standard go tooling for working with modules to incorporate it.

## Vigil Message Overview

The Vigil protocol consists of exchanging messages between peers. Each message
is preceded by a header which identifies information about it such as which
Vigil network it is a part of, its type, how big it is, and a checksum to
verify validity. All encoding and decoding of message headers is handled by this
package.

To accomplish this, there is a generic interface for Vigil messages named
`Message` which allows messages of any type to be read, written, or passed
around through channels, functions, etc. In addition, concrete implementations
of most of the currently supported Vigil messages are provided. For these
supported messages, all of the details of marshalling and unmarshalling to and
from the wire using Vigil encoding are handled so the caller doesn't have to
concern themselves with the specifics.

## Reading Messages Example

In order to unmarshal Vigil messages from the wire, use the `ReadMessage`
function. It accepts any `io.Reader`, but typically this will be a `net.Conn`
to a remote node running a Vigil peer. Example syntax is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main Vigil network.
	pver := wire.ProtocolVersion
	VGLnet := wire.MainNet

	// Reads and validates the next Vigil message from conn using the
	// protocol version pver and the Vigil network VGLnet. The returns
	// are a wire.Message, a []byte which contains the unmarshalled
	// raw payload, and a possible error.
	msg, rawPayload, err := wire.ReadMessage(conn, pver, VGLnet)
	if err != nil {
		// Log and handle the error
	}
```

See the package documentation for details on determining the message type.

## Writing Messages Example

In order to marshal Vigil messages to the wire, use the `WriteMessage`
function. It accepts any `io.Writer`, but typically this will be a `net.Conn`
to a remote node running a Vigil peer. Example syntax to request addresses
from a remote peer is:

```Go
	// Use the most recent protocol version supported by the package and the
	// main Vigil network.
	pver := wire.ProtocolVersion
	VGLnet := wire.MainNet

	// Create a new getaddr Vigil message.
	msg := wire.NewMsgGetAddr()

	// Writes a Vigil message msg to conn using the protocol version
	// pver, and the Vigil network VGLnet. The return is a possible
	// error.
	err := wire.WriteMessage(conn, msg, pver, VGLnet)
	if err != nil {
		// Log and handle the error
	}
```

## License

Package wire is licensed under the [copyfree](http://copyfree.org) ISC
License.




