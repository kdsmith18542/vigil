// Copyright (c) 2015-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ffldb

import (
	"fmt"

	"github.com/kdsmith18542/vigil/database/v3"
	"github.com/kdsmith18542/vigil/wire"
	"github.com/kdsmith18542/vigil/slog"
)

var log = slog.Disabled

const (
	dbType = "ffldb"
)

// parseArgs parses the arguments from the database Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (string, wire.CurrencyNet, error) {
	if len(args) != 2 {
		return "", 0, fmt.Errorf("invalid arguments to %s.%s -- "+
			"expected database path and block network", dbType,
			funcName)
	}

	dbPath, ok := args[0].(string)
	if !ok {
		return "", 0, fmt.Errorf("first argument to %s.%s is invalid -- "+
			"expected database path string", dbType, funcName)
	}

	network, ok := args[1].(wire.CurrencyNet)
	if !ok {
		return "", 0, fmt.Errorf("second argument to %s.%s is invalid -- "+
			"expected block network", dbType, funcName)
	}

	return dbPath, network, nil
}

// openDBDriver is the callback provided during driver registration that opens
// an existing database for use.
func openDBDriver(args ...interface{}) (database.DB, error) {
	dbPath, network, err := parseArgs("Open", args...)
	if err != nil {
		return nil, err
	}

	return openDB(dbPath, network, false)
}

// createDBDriver is the callback provided during driver registration that
// creates, initializes, and opens a database for use.
func createDBDriver(args ...interface{}) (database.DB, error) {
	dbPath, network, err := parseArgs("Create", args...)
	if err != nil {
		return nil, err
	}

	return openDB(dbPath, network, true)
}

// useLogger is the callback provided during driver registration that sets the
// current logger to the provided one.
func useLogger(logger slog.Logger) {
	log = logger
}

func init() {
	// Register the driver.
	driver := database.Driver{
		DbType:    dbType,
		Create:    createDBDriver,
		Open:      openDBDriver,
		UseLogger: useLogger,
	}
	if err := database.RegisterDriver(driver); err != nil {
		panic(fmt.Sprintf("Failed to register database driver '%s': %v",
			dbType, err))
	}
}
