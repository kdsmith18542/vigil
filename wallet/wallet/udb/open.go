// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package udb

import (
	"context"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/wallet/walletdb"
	"github.com/kdsmith18542/vigil/chaincfg/v3"
)

// Open opens the database and returns various "manager" types that must be used
// to access and modify data in the database.
//
// A NotExist error will be returned if the database has not been initialized.
// The recorded database version must match exactly with DBVersion.  If the
// version does not match, an Invalid error is returned.
func Open(ctx context.Context, db walletdb.DB, params *chaincfg.Params, pubPass []byte) (addrMgr *Manager, txStore *Store, err error) {
	err = walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		// Verify the database exists and the recorded version is supported by
		// this software version.
		metadataBucket := tx.ReadBucket(unifiedDBMetadata{}.rootBucketKey())
		if metadataBucket == nil {
			return errors.E(errors.NotExist, "database has not been initialized")
		}
		dbVersion, err := unifiedDBMetadata{}.getVersion(metadataBucket)
		if err != nil {
			return err
		}
		if dbVersion < DBVersion {
			return errors.E(errors.Invalid, "database upgrade required")
		}
		if dbVersion > DBVersion {
			return errors.E(errors.Invalid, "database has been upgraded to an unknown newer version")
		}

		addrmgrNs := tx.ReadBucket(waddrmgrBucketKey)

		addrMgr, err = loadManager(addrmgrNs, pubPass, params)
		if err != nil {
			return err
		}
		txStore = &Store{
			chainParams:    params,
			acctLookupFunc: addrMgr.AddrAccount,
			manager:        addrMgr,
		}
		return err
	})
	return
}
