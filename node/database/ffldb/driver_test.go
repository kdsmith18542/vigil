// Copyright (c) 2015-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ffldb_test

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/kdsmith18542/vigil/chaincfg/v3"
	"github.com/kdsmith18542/vigil/database/v3"
	"github.com/kdsmith18542/vigil/database/v3/ffldb"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
)

// dbType is the database type name for this driver.
const dbType = "ffldb"

// TestCreateOpenFail ensures that errors related to creating and opening a
// database are handled properly.
func TestCreateOpenFail(t *testing.T) {
	t.Parallel()

	// Ensure that attempting to open a database that doesn't exist returns
	// the expected error.
	wantErrKind := database.ErrDbDoesNotExist
	_, err := database.Open(dbType, "noexist", blockDataNet)
	if !checkDbError(t, "Open", err, wantErrKind) {
		return
	}

	// Ensure that attempting to open a database with the wrong number of
	// parameters returns the expected error.
	wantErr := fmt.Errorf("invalid arguments to %s.Open -- expected "+
		"database path and block network", dbType)
	_, err = database.Open(dbType, 1, 2, 3)
	if err.Error() != wantErr.Error() {
		t.Errorf("Open: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to open a database with an invalid type for
	// the first parameter returns the expected error.
	wantErr = fmt.Errorf("first argument to %s.Open is invalid -- "+
		"expected database path string", dbType)
	_, err = database.Open(dbType, 1, blockDataNet)
	if err.Error() != wantErr.Error() {
		t.Errorf("Open: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to open a database with an invalid type for
	// the second parameter returns the expected error.
	wantErr = fmt.Errorf("second argument to %s.Open is invalid -- "+
		"expected block network", dbType)
	_, err = database.Open(dbType, "noexist", "invalid")
	if err.Error() != wantErr.Error() {
		t.Errorf("Open: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to create a database with the wrong number of
	// parameters returns the expected error.
	wantErr = fmt.Errorf("invalid arguments to %s.Create -- expected "+
		"database path and block network", dbType)
	_, err = database.Create(dbType, 1, 2, 3)
	if err.Error() != wantErr.Error() {
		t.Errorf("Create: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to create a database with an invalid type for
	// the first parameter returns the expected error.
	wantErr = fmt.Errorf("first argument to %s.Create is invalid -- "+
		"expected database path string", dbType)
	_, err = database.Create(dbType, 1, blockDataNet)
	if err.Error() != wantErr.Error() {
		t.Errorf("Create: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure that attempting to create a database with an invalid type for
	// the second parameter returns the expected error.
	wantErr = fmt.Errorf("second argument to %s.Create is invalid -- "+
		"expected block network", dbType)
	_, err = database.Create(dbType, "noexist", "invalid")
	if err.Error() != wantErr.Error() {
		t.Errorf("Create: did not receive expected error - got %v, "+
			"want %v", err, wantErr)
		return
	}

	// Ensure operations against a closed database return the expected
	// error.
	dbPath := t.TempDir()
	db, err := database.Create(dbType, dbPath, blockDataNet)
	if err != nil {
		t.Errorf("Create: unexpected error: %v", err)
		return
	}
	db.Close()

	wantErrKind = database.ErrDbNotOpen
	err = db.View(func(tx database.Tx) error {
		return nil
	})
	if !checkDbError(t, "View", err, wantErrKind) {
		return
	}

	wantErrKind = database.ErrDbNotOpen
	err = db.Update(func(tx database.Tx) error {
		return nil
	})
	if !checkDbError(t, "Update", err, wantErrKind) {
		return
	}

	wantErrKind = database.ErrDbNotOpen
	_, err = db.Begin(false)
	if !checkDbError(t, "Begin(false)", err, wantErrKind) {
		return
	}

	wantErrKind = database.ErrDbNotOpen
	_, err = db.Begin(true)
	if !checkDbError(t, "Begin(true)", err, wantErrKind) {
		return
	}

	wantErrKind = database.ErrDbNotOpen
	err = db.Close()
	if !checkDbError(t, "Close", err, wantErrKind) {
		return
	}
}

// TestPersistence ensures that values stored are still valid after closing and
// reopening the database.
func TestPersistence(t *testing.T) {
	t.Parallel()

	// Create a new database to run tests against.
	dbPath := t.TempDir()
	db, err := database.Create(dbType, dbPath, blockDataNet)
	if err != nil {
		t.Errorf("Failed to create test database (%s) %v", dbType, err)
		return
	}
	t.Cleanup(func() {
		db.Close()
	})

	// Create a bucket, put some values into it, and store a block so they
	// can be tested for existence on re-open.
	bucket1Key := []byte("bucket1")
	storeValues := map[string]string{
		"b1key1": "foo1",
		"b1key2": "foo2",
		"b1key3": "foo3",
	}
	mainNetParams := chaincfg.MainNetParams()
	genesisBlock := VGLutil.NewBlock(mainNetParams.GenesisBlock)
	genesisHash := &mainNetParams.GenesisHash
	err = db.Update(func(tx database.Tx) error {
		metadataBucket := tx.Metadata()
		if metadataBucket == nil {
			return fmt.Errorf("Metadata: unexpected nil bucket")
		}

		bucket1, err := metadataBucket.CreateBucket(bucket1Key)
		if err != nil {
			return fmt.Errorf("CreateBucket: unexpected error: %w",
				err)
		}

		for k, v := range storeValues {
			err := bucket1.Put([]byte(k), []byte(v))
			if err != nil {
				return fmt.Errorf("Put: unexpected error: %w",
					err)
			}
		}

		if err := tx.StoreBlock(genesisBlock); err != nil {
			return fmt.Errorf("StoreBlock: unexpected error: %w",
				err)
		}

		return nil
	})
	if err != nil {
		t.Errorf("Update: unexpected error: %v", err)
		return
	}

	// Close and reopen the database to ensure the values persist.
	db.Close()
	db, err = database.Open(dbType, dbPath, blockDataNet)
	if err != nil {
		t.Errorf("failed to open test database (%s) %v", dbType, err)
		return
	}
	defer db.Close()

	// Ensure the values previously stored in the 3rd namespace still exist
	// and are correct.
	err = db.View(func(tx database.Tx) error {
		metadataBucket := tx.Metadata()
		if metadataBucket == nil {
			return fmt.Errorf("Metadata: unexpected nil bucket")
		}

		bucket1 := metadataBucket.Bucket(bucket1Key)
		if bucket1 == nil {
			return fmt.Errorf("bucket1: unexpected nil bucket")
		}

		for k, v := range storeValues {
			gotVal := bucket1.Get([]byte(k))
			if !reflect.DeepEqual(gotVal, []byte(v)) {
				return fmt.Errorf("Get: key '%s' does not "+
					"match expected value - got %s, want %s",
					k, gotVal, v)
			}
		}

		genesisBlockBytes, _ := genesisBlock.Bytes()
		gotBytes, err := tx.FetchBlock(genesisHash)
		if err != nil {
			return fmt.Errorf("FetchBlock: unexpected error: %w",
				err)
		}
		if !reflect.DeepEqual(gotBytes, genesisBlockBytes) {
			return fmt.Errorf("FetchBlock: stored block mismatch")
		}

		return nil
	})
	if err != nil {
		t.Errorf("View: unexpected error: %v", err)
		return
	}
}

// TestInterface performs all interfaces tests for this database driver.
func TestInterface(t *testing.T) {
	t.Parallel()

	// Create a new database to run tests against.
	dbPath := t.TempDir()
	db, err := database.Create(dbType, dbPath, blockDataNet)
	if err != nil {
		t.Errorf("failed to create test database (%s) %v", dbType, err)
		return
	}
	t.Cleanup(func() {
		db.Close()
	})

	// Ensure the driver type is the expected value.
	gotDbType := db.Type()
	if gotDbType != dbType {
		t.Errorf("Type: unexpected driver type - got %v, want %v",
			gotDbType, dbType)
		return
	}

	// Run all of the interface tests against the database.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Change the maximum file size to a small value to force multiple flat
	// files with the test data set.
	ffldb.TstRunWithMaxBlockFileSize(db, 2048, func() {
		testInterface(t, db)
	})
}
