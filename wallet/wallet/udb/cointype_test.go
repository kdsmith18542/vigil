// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package udb

import (
	"context"
	"testing"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/internal/compat"
	"github.com/kdsmith18542/vigil/wallet/wallet/walletdb"
	"github.com/kdsmith18542/vigil/chaincfg/v3"
	"github.com/kdsmith18542/vigil/hdkeychain/v3"
	"github.com/kdsmith18542/vigil/txscript/v4/stdaddr"
)

func TestCoinTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		params                           *chaincfg.Params
		legacyCoinType, slip0044CoinType uint32
	}{
		{chaincfg.MainNetParams(), 20, 42},
		{chaincfg.TestNet3Params(), 11, 1},
		{chaincfg.SimNetParams(), 115, 1},
	}
	for _, test := range tests {
		legacyCoinType, slip0044CoinType := CoinTypes(test.params)
		if legacyCoinType != test.legacyCoinType {
			t.Errorf("%s: got legacy coin type %d, expected %d", test.params.Name,
				legacyCoinType, test.legacyCoinType)
		}
		if slip0044CoinType != test.slip0044CoinType {
			t.Errorf("%s: got SLIP0044 coin type %d, expected %d", test.params.Name,
				slip0044CoinType, test.slip0044CoinType)
		}
	}
}

func deriveChildAddress(accountExtKey *hdkeychain.ExtendedKey, branch, child uint32, params *chaincfg.Params) (stdaddr.Address, error) {
	branchKey, err := accountExtKey.Child(branch)
	if err != nil {
		return nil, err
	}
	addressKey, err := branchKey.Child(child)
	if err != nil {
		return nil, err
	}
	return compat.HD2Address(addressKey, params)
}

func equalExtKeys(k0, k1 *hdkeychain.ExtendedKey) bool {
	return k0.String() == k1.String()
}

func TestCoinTypeUpgrade(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db, teardown := tempDB(t)
	defer teardown()

	params := chaincfg.TestNet3Params()

	err := Initialize(ctx, db, params, seed, pubPass, privPassphrase)
	if err != nil {
		t.Fatal(err)
	}

	m, _, err := Open(ctx, db, params, pubPass)
	if err != nil {
		t.Fatal(err)
	}

	legacyCoinType, slip0044CoinType := CoinTypes(params)

	masterExtKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		t.Fatal(err)
	}
	legacyCoinTypeExtKey, err := deriveCoinTypeKey(masterExtKey, legacyCoinType)
	if err != nil {
		t.Fatal(err)
	}
	slip0044CoinTypeExtKey, err := deriveCoinTypeKey(masterExtKey, slip0044CoinType)
	if err != nil {
		t.Fatal(err)
	}
	slip0044Account0ExtKey, err := deriveAccountKey(slip0044CoinTypeExtKey, 0)
	if err != nil {
		t.Fatal(err)
	}
	slip0044Account0ExtKey = slip0044Account0ExtKey.Neuter()
	slip0044Account1ExtKey, err := deriveAccountKey(slip0044CoinTypeExtKey, 1)
	if err != nil {
		t.Fatal(err)
	}
	slip0044Account1ExtKey = slip0044Account1ExtKey.Neuter()
	slip0044Account0Address0, err := deriveChildAddress(slip0044Account0ExtKey, 0, 0, params)
	if err != nil {
		t.Fatal(err)
	}
	slip0044Account1Address0, err := deriveChildAddress(slip0044Account1ExtKey, 0, 0, params)
	if err != nil {
		t.Fatal(err)
	}
	slip0044Account0Address0Hash160 := slip0044Account0Address0.(stdaddr.Hash160er).Hash160()
	slip0044Account1Address0Hash160 := slip0044Account1Address0.(stdaddr.Hash160er).Hash160()

	err = walletdb.Update(ctx, db, func(dbtx walletdb.ReadWriteTx) error {
		ns := dbtx.ReadWriteBucket(waddrmgrBucketKey)
		err := m.Unlock(ns, privPassphrase)
		if err != nil {
			t.Fatal(err)
		}

		// Check reported initial coin type and compare the key itself against
		// the expected value.
		coinType, err := m.CoinType(dbtx)
		if err != nil {
			t.Fatal(err)
		}
		if coinType != legacyCoinType {
			t.Fatalf("initialized database has wrong coin type %d", coinType)
		}
		coinTypeExtKey, err := m.CoinTypePrivKey(dbtx)
		if err != nil {
			t.Fatal(err)
		}
		if !equalExtKeys(coinTypeExtKey, legacyCoinTypeExtKey) {
			t.Fatalf("initialized database has wrong coin type key")
		}

		// Perform the upgrade
		err = m.UpgradeToSLIP0044CoinType(dbtx)
		if err != nil {
			t.Fatal(err)
		}

		// Check upgraded coin type and keys.
		coinType, err = m.CoinType(dbtx)
		if err != nil {
			t.Fatal(err)
		}
		if coinType != slip0044CoinType {
			t.Fatalf("upgraded database has wrong coin type %d", coinType)
		}
		coinTypeExtKey, err = m.CoinTypePrivKey(dbtx)
		if err != nil {
			t.Fatal(err)
		}
		if !equalExtKeys(coinTypeExtKey, slip0044CoinTypeExtKey) {
			t.Fatalf("upgraded database has wrong coin type key")
		}

		// Check the account 0 xpub matches the one derived from the SLIP0044
		// coin type.
		accountExtKey, err := m.AccountExtendedPubKey(dbtx, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !equalExtKeys(accountExtKey, slip0044Account0ExtKey) {
			t.Fatalf("upgraded database has wrong account xpub")
		}

		// Check that the SLIP0044-derived account 0's first address can be
		// created and is indexed.
		err = m.SyncAccountToAddrIndex(ns, 0, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !m.ExistsHash160(ns, slip0044Account0Address0Hash160[:]) {
			t.Fatalf("upgraded database does not record SLIP0044-derived account 0 branch 0 address 0")
		}

		// Create the next account, and perform all of the same checks on it as
		// the first account.
		_, err = m.NewAccount(ns, "account-1")
		if err != nil {
			t.Fatal(err)
		}
		accountExtKey, err = m.AccountExtendedPubKey(dbtx, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !equalExtKeys(accountExtKey, slip0044Account1ExtKey) {
			t.Fatal("upgraded database derived wrong account xpub")
		}
		err = m.SyncAccountToAddrIndex(ns, 1, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !m.ExistsHash160(ns, slip0044Account1Address0Hash160[:]) {
			t.Fatalf("upgraded database does not record SLIP0044-derived account 1 branch 0 address 0")
		}

		// Check that the upgrade can not be performed a second time.
		err = m.UpgradeToSLIP0044CoinType(dbtx)
		if !errors.Is(err, errors.Invalid) {
			t.Fatalf("upgrade database did not refuse second upgrade with errors.Invalid")
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
