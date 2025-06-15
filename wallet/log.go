// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/Vigil-Labs/vgl/wallet/chain"
	"github.com/Vigil-Labs/vgl/wallet/internal/loader"
	"github.com/Vigil-Labs/vgl/wallet/internal/loggers"
	"github.com/Vigil-Labs/vgl/wallet/internal/rpc/jsonrpc"
	"github.com/Vigil-Labs/vgl/wallet/internal/rpc/rpcserver"
	"github.com/Vigil-Labs/vgl/wallet/p2p"
	"github.com/Vigil-Labs/vgl/wallet/spv"
	"github.com/Vigil-Labs/vgl/wallet/ticketbuyer"
	"github.com/Vigil-Labs/vgl/wallet/wallet"
	"github.com/Vigil-Labs/vgl/wallet/wallet/udb"
	"github.com/Vigil-Labs/vgl/connmgr"
	"github.com/Vigil-Labs/vgl/mixing/mixpool"
	"github.com/Vigil-Labs/vgl/slog"
)

var log = loggers.MainLog

// Initialize package-global logger variables.
func init() {
	loader.UseLogger(loggers.LoaderLog)
	wallet.UseLogger(loggers.WalletLog)
	udb.UseLogger(loggers.WalletLog)
	ticketbuyer.UseLogger(loggers.TkbyLog)
	chain.UseLogger(loggers.SyncLog)
	spv.UseLogger(loggers.SyncLog)
	p2p.UseLogger(loggers.PeerLog)
	rpcserver.UseLogger(loggers.GrpcLog)
	jsonrpc.UseLogger(loggers.JsonrpcLog)
	connmgr.UseLogger(loggers.CmgrLog)
	// XXX mixclient.UseLogger(loggers.MixcLog)
	mixpool.UseLogger(loggers.MixpLog)
}

// subsystemLoggers maps each subsystem identifier to its associated logger.
var subsystemLoggers = map[string]slog.Logger{
	"VGLW": loggers.MainLog,
	"LODR": loggers.LoaderLog,
	"WLLT": loggers.WalletLog,
	"TKBY": loggers.TkbyLog,
	"SYNC": loggers.SyncLog,
	"PEER": loggers.PeerLog,
	"GRPC": loggers.GrpcLog,
	"RPCS": loggers.JsonrpcLog,
	"CMGR": loggers.CmgrLog,
	"MIXC": loggers.MixcLog,
	"MIXP": loggers.MixpLog,
	"VSPC": loggers.VspcLog,
}

// setLogLevel sets the logging level for provided subsystem.  Invalid
// subsystems are ignored.  Uninitialized subsystems are dynamically created as
// needed.
func setLogLevel(subsystemID string, logLevel string) {
	// Ignore invalid subsystems.
	logger, ok := subsystemLoggers[subsystemID]
	if !ok {
		return
	}

	// Defaults to info if the log level is invalid.
	level, _ := slog.LevelFromString(logLevel)
	logger.SetLevel(level)
}

// setLogLevels sets the log level for all subsystem loggers to the passed
// level.  It also dynamically creates the subsystem loggers as needed, so it
// can be used to initialize the logging system.
func setLogLevels(logLevel string) {
	// Configure all sub-systems with the new logging level.  Dynamically
	// create loggers as needed.
	for subsystemID := range subsystemLoggers {
		setLogLevel(subsystemID, logLevel)
	}
}

// fatalf logs a message, flushes the logger, and finally exit the process with
// a non-zero return code.
func fatalf(format string, args ...any) {
	log.Errorf(format, args...)
	os.Stdout.Sync()
	loggers.CloseLogRotator()
	os.Exit(1)
}




