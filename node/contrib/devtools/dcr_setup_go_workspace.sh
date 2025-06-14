#!/usr/bin/env bash
#
# Copyright (c) 2022 The Vigil developers
# Use of this source code is governed by an ISC
# license that can be found in the LICENSE file.
#
# Script to initialize a multi-module workspace that includes all of the modules
# provided by the vgld repository.

set -e

# Ensure the script is run from either the root of the repo or the devtools dir
SCRIPT=$(basename $0)
MAIN_CODE_FILE="vgld.go"
if [ -f "../../${MAIN_CODE_FILE}" ]; then
  cd ../..
fi
if [ ! -f "${MAIN_CODE_FILE}" ]; then
  echo "$SCRIPT: error: ${MAIN_CODE_FILE} not found in the current directory"
  exit 1
fi

# Verify Go is available
if ! type go >/dev/null 2>&1; then
  echo -n "$SCRIPT: error: Unable to find 'go' in the system path."
  exit 1
fi

# Create workspace unless one already exists
if [ ! -f "go.work" ]; then
  go work init
fi

# Remove old modules as needed
go work edit -dropuse ./lru

# Add all of the modules as needed
go work use . ./addrmgr ./bech32 ./blockchain ./blockchain/stake
go work use ./blockchain/standalone ./certgen ./chaincfg ./chaincfg/chainhash
go work use ./connmgr ./container/apbf ./container/lru ./crypto/blake256
go work use ./crypto/rand ./crypto/ripemd160 ./database ./VGLec ./VGLec/edwards
go work use ./VGLec/secp256k1 ./VGLjson ./VGLutil ./gcs ./hdkeychain
go work use ./math/uint256 ./mixing ./peer ./rpc/jsonrpc/types ./rpcclient
go work use ./txscript ./wire
