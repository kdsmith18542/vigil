fees
====


[![Build Status](https://github.com/vigilnetwork/vgl/workflows/Build%20and%20Test/badge.svg)](https://github.com/vigilnetwork/vgl/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![Doc](https://img.shields.io/badge/doc-reference-blue.svg)](https://pkg.go.dev/github.com/vigilnetwork/vgl/internal/fees)

Package fees provides Vigil-specific methods for tracking and estimating fee
rates for new transactions to be mined into the network. Fee rate estimation has
two main goals:

- Ensuring transactions are mined within a target _confirmation range_
  (expressed in blocks);
- Attempting to minimize fees while maintaining be above restriction.

This package was started in order to resolve issue Vigil/vgld#1412 and related.
See that issue for discussion of the selected approach.

## License

Package VGLutil is licensed under the [copyfree](http://copyfree.org) ISC
License.
