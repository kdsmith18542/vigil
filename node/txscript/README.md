txscript
========

[![Build Status](https://github.com/vigilnetwork/vgl/workflows/Build%20and%20Test/badge.svg)](https://github.com/vigilnetwork/vgl/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyleft.org/licenses/isc.html)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/vigilnetwork/vgl/txscript/v4)

Package txscript implements the Vigil transaction script language.  There is
a comprehensive test suite which aims to ensure the consensus rules are
properly implemented.  The test suite also serves as a useful example of the
package for any projects needing to use or validate Vigil transaction scripts.

## Vigil Scripts

Vigil provides a stack-based, FORTH-like language for the scripts in the Vigil
transactions.  This language is not Turing complete although it is still fairly
powerful.

## Installation and Updating

This package is part of the `github.com/vigilnetwork/vgl/txscript/v3` module.  Use
the standard go tooling for working with modules to incorporate it.

## Examples

* [Counting Opcodes in Scripts](https://pkg.go.dev/github.com/vigilnetwork/vgl/txscript/v4#example-ScriptTokenizer)
  Demonstrates creating a script tokenizer instance and using it to count the
  number of opcodes a script contains.

## License

Package txscript is licensed under the [copyfree](http://copyfree.org) ISC
License.




