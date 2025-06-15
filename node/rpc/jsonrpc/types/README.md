jsonrpc/types
=============

[![Build Status](https://github.com/vigilnetwork/vgl/workflows/Build%20and%20Test/badge.svg)](https://github.com/vigilnetwork/vgl/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/vigilnetwork/vgl/rpc/jsonrpc/types/v4)

Package types implements concrete types for marshalling to and from the vgld
JSON-RPC commands, return values, and notifications.  A comprehensive suite of
tests is provided to ensure proper functionality.

The provided types are automatically registered with
[VGLjson](https://github.com/vigilnetwork/vgl/tree/master/VGLjson) when the package
is imported.  Although this package was primarily written for vgld, it has
intentionally been designed so it can be used as a standalone package for any
projects needing to marshal to and from vgld JSON-RPC requests and responses.

## Installation and Updating

This package is part of the `github.com/vigilnetwork/vgl/rpc/jsonrpc/types/v2`
module.  Use the standard go tooling for working with modules to incorporate it.

## License

Package types is licensed under the [copyfree](http://copyfree.org) ISC License.
