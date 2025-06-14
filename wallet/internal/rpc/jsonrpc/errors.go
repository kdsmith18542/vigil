// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package jsonrpc

import (
	"fmt"

	"github.com/vigilnetwork/vgl/wallet/errors"
	"github.com/vigilnetwork/vgl/VGLjson/v4"
	"github.com/jrick/wsrpc/v2"
)

func convertError(err error) *VGLjson.RPCError {
	switch err := err.(type) {
	case *VGLjson.RPCError:
		return err
	case *wsrpc.Error:
		return &VGLjson.RPCError{
			Code:    VGLjson.RPCErrorCode(err.Code),
			Message: err.Message,
		}
	}

	code := VGLjson.ErrRPCWallet
	var kind errors.Kind
	if errors.As(err, &kind) {
		switch kind {
		case errors.Bug:
			code = VGLjson.ErrRPCInternal.Code
		case errors.Encoding:
			code = VGLjson.ErrRPCInvalidParameter
		case errors.Locked:
			code = VGLjson.ErrRPCWalletUnlockNeeded
		case errors.Passphrase:
			code = VGLjson.ErrRPCWalletPassphraseIncorrect
		case errors.NoPeers:
			code = VGLjson.ErrRPCClientNotConnected
		case errors.InsufficientBalance:
			code = VGLjson.ErrRPCWalletInsufficientFunds
		}
	}
	return &VGLjson.RPCError{
		Code:    code,
		Message: err.Error(),
	}
}

func rpcError(code VGLjson.RPCErrorCode, err error) *VGLjson.RPCError {
	return &VGLjson.RPCError{
		Code:    code,
		Message: err.Error(),
	}
}

func rpcErrorf(code VGLjson.RPCErrorCode, format string, args ...any) *VGLjson.RPCError {
	return &VGLjson.RPCError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Errors variables that are defined once here to avoid duplication.
var (
	errUnloadedWallet = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCWallet,
		Message: "request requires a wallet but wallet has not loaded yet",
	}

	errRPCClientNotConnected = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCClientNotConnected,
		Message: "disconnected from consensus RPC",
	}

	errNoNetwork = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCClientNotConnected,
		Message: "disconnected from network",
	}

	errAccountNotFound = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCWalletInvalidAccountName,
		Message: "account not found",
	}

	errAddressNotInWallet = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCWallet,
		Message: "address not found in wallet",
	}

	errNotImportedAccount = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCWallet,
		Message: "imported addresses must belong to the imported account",
	}

	errNeedPositiveAmount = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCInvalidParameter,
		Message: "amount must be positive",
	}

	errWalletUnlockNeeded = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCWalletUnlockNeeded,
		Message: "wallet or account locked; use walletpassphrase or unlockaccount first",
	}

	errReservedAccountName = &VGLjson.RPCError{
		Code:    VGLjson.ErrRPCInvalidParameter,
		Message: "account name is reserved by RPC server",
	}
)
