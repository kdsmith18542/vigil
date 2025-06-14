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

// Initialize prepares an empty database for usage by initializing all buckets
// and key/value pairs.  The database is initialized with the latest version and
// does not require any upgrades to use.
func Initialize(ctx context.Context, db walletdb.DB, params *chaincfg.Params, seed, pubPass, privPass []byte) error {
	err := walletdb.Update(ctx, db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, err := tx.CreateTopLevelBucket(waddrmgrBucketKey)
		if err != nil {
			return errors.E(errors.IO, err)
		}
		txmgrNs, err := tx.CreateTopLevelBucket(wtxmgrBucketKey)
		if err != nil {
			return errors.E(errors.IO, err)
		}

		// Create the address manager and transaction store.
		err = createAddressManager(addrmgrNs, seed, pubPass, privPass, params)
		if err != nil {
			return err
		}
		err = createStore(txmgrNs, params)
		if err != nil {
			return err
		}

		// Create the metadata bucket and write the current database version to
		// it.
		metadataBucket, err := tx.CreateTopLevelBucket(unifiedDBMetadata{}.rootBucketKey())
		if err != nil {
			return errors.E(errors.IO, err)
		}
		return unifiedDBMetadata{}.putVersion(metadataBucket, initialVersion)
	})
	if err != nil {
		return err
	}
	return Upgrade(ctx, db, pubPass, params)
}

// InitializeWatchOnly prepares an empty database for watching-only wallet usage
// by initializing all buckets and key/value pairs.  The database is initialized
// with the latest version and does not require any upgrades to use.
func InitializeWatchOnly(ctx context.Context, db walletdb.DB, params *chaincfg.Params, hdPubKey string, pubPass []byte) error {
	err := walletdb.Update(ctx, db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs, err := tx.CreateTopLevelBucket(waddrmgrBucketKey)
		if err != nil {
			return errors.E(errors.IO, err)
		}
		txmgrNs, err := tx.CreateTopLevelBucket(wtxmgrBucketKey)
		if err != nil {
			return errors.E(errors.IO, err)
		}

		// Create the address manager and transaction store.
		err = createWatchOnly(addrmgrNs, hdPubKey, pubPass, params)
		if err != nil {
			return err
		}
		err = createStore(txmgrNs, params)
		if err != nil {
			return err
		}

		// Create the metadata bucket and write the current database version to
		// it.
		metadataBucket, err := tx.CreateTopLevelBucket(unifiedDBMetadata{}.rootBucketKey())
		if err != nil {
			return errors.E(errors.IO, err)
		}
		return unifiedDBMetadata{}.putVersion(metadataBucket, initialVersion)
	})
	if err != nil {
		return err
	}
	return Upgrade(ctx, db, pubPass, params)
}
