// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kdsmith18542/vigil/database/v3"
	"github.com/kdsmith18542/vigil/internal/blockchain"
	"github.com/kdsmith18542/vigil/internal/blockchain/indexers"
	"github.com/kdsmith18542/vigil/internal/limits"
	"github.com/kdsmith18542/vigil/slog"
)

const (
	// blockDbNamePrefix is the prefix for the vgld block database.
	blockDbNamePrefix = "blocks"
)

var (
	cfg *config
	log slog.Logger
)

// loadBlockDB opens the block database and returns a handle to it.
func loadBlockDB() (database.DB, error) {
	// The database name is based on the database type.
	dbName := blockDbNamePrefix + "_" + cfg.DbType
	dbPath := filepath.Join(cfg.DataDir, dbName)

	log.Infof("Loading block database from '%s'", dbPath)
	db, err := database.Open(cfg.DbType, dbPath, activeNetParams.Net)
	if err != nil {
		// Return the error if it's not because the database doesn't
		// exist.
		if !errors.Is(err, database.ErrDbDoesNotExist) {
			return nil, err
		}

		// Create the db if it does not exist.
		err = os.MkdirAll(cfg.DataDir, 0700)
		if err != nil {
			return nil, err
		}
		db, err = database.Create(cfg.DbType, dbPath, activeNetParams.Net)
		if err != nil {
			return nil, err
		}
	}

	log.Info("Block database loaded")
	return db, nil
}

// realMain is the real main function for the utility.  It is necessary to work
// around the fact that deferred functions do not run when os.Exit() is called.
func realMain() error {
	// Load configuration and parse command line.
	tcfg, _, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = tcfg

	// Setup logging.
	backendLogger := slog.NewBackend(os.Stdout)
	defer os.Stdout.Sync()
	log = backendLogger.Logger("MAIN")
	database.UseLogger(backendLogger.Logger("BCDB"))
	blockchain.UseLogger(backendLogger.Logger("CHAN"))
	indexers.UseLogger(backendLogger.Logger("INDX"))

	// Load the block database.
	db, err := loadBlockDB()
	if err != nil {
		log.Errorf("Failed to load database: %v", err)
		return err
	}
	defer db.Close()

	// Load the UTXO database.
	utxoDb, err := blockchain.LoadUtxoDB(context.Background(), activeNetParams,
		cfg.DataDir)
	if err != nil {
		log.Errorf("Failed to load UTXO database: %v", err)
		return err
	}
	defer utxoDb.Close()

	fi, err := os.Open(cfg.InFile)
	if err != nil {
		log.Errorf("Failed to open file %v: %v", cfg.InFile, err)
		return err
	}
	defer fi.Close()

	// Create a block importer for the database and input file and start it.
	// The done channel returned from start will contain an error if
	// anything went wrong.
	ctx, cancel := context.WithCancel(context.Background())
	importer, err := newBlockImporter(ctx, db, utxoDb, fi, cancel)
	if err != nil {
		log.Errorf("Failed create block importer: %v", err)
		return err
	}

	// Perform the import asynchronously.  This allows blocks to be
	// processed and read in parallel.  The results channel returned from
	// Import contains the statistics about the import including an error
	// if something went wrong.
	log.Info("Starting import")
	resultsChan := importer.Import(ctx)
	results := <-resultsChan
	if results.err != nil {
		log.Errorf("%v", results.err)
		return results.err
	}

	log.Infof("Processed a total of %d blocks (%d imported, %d already "+
		"known) in %v", results.blocksProcessed, results.blocksImported,
		results.blocksProcessed-results.blocksImported, results.duration)

	return nil
}

func main() {
	// Use all processor cores and up some limits.
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := limits.SetLimits(); err != nil {
		os.Exit(1)
	}

	// Work around defer not working after os.Exit()
	if err := realMain(); err != nil {
		os.Exit(1)
	}
}
