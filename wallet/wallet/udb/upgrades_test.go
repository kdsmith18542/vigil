// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package udb

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/kdsmith18542/vigil/wallet/wallet/drivers/bdb"
	"github.com/kdsmith18542/vigil/wallet/wallet/walletdb"
	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	"github.com/kdsmith18542/vigil/chaincfg/v3"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
	"github.com/kdsmith18542/vigil/wire"
)

var dbUpgradeTests = [...]struct {
	verify   func(context.Context, *testing.T, walletdb.DB)
	filename string // in testdata directory
}{
	{verifyV2Upgrade, "v1.db.gz"},
	{verifyV3Upgrade, "v2.db.gz"},
	{verifyV4Upgrade, "v3.db.gz"},
	{verifyV5Upgrade, "v4.db.gz"},
	{verifyV6Upgrade, "v5.db.gz"},
	// No upgrade test for V7, it is a backwards-compatible upgrade
	{verifyV8Upgrade, "v7.db.gz"},
	// No upgrade test for V9, it is a fix for V8 and the previous test still applies
	// TODO: V10 upgrade test
	{verifyV12Upgrade, "v11.db.gz"},
}

var pubPass = []byte("public")

func TestUpgrades(t *testing.T) {
	ctx := context.Background()
	d := t.TempDir()

	t.Run("group", func(t *testing.T) {
		for i, test := range dbUpgradeTests {
			test := test
			name := fmt.Sprintf("test%d", i)
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				testFile, err := os.Open(filepath.Join("testdata", test.filename))
				if err != nil {
					t.Fatal(err)
				}
				defer testFile.Close()
				r, err := gzip.NewReader(testFile)
				if err != nil {
					t.Fatal(err)
				}
				dbPath := filepath.Join(d, name+".db")
				fi, err := os.Create(dbPath)
				if err != nil {
					t.Fatal(err)
				}
				_, err = io.Copy(fi, r)
				fi.Close()
				if err != nil {
					t.Fatal(err)
				}
				db, err := walletdb.Open("bdb", dbPath)
				if err != nil {
					t.Fatal(err)
				}
				defer db.Close()
				err = Upgrade(ctx, db, pubPass, chaincfg.TestNet3Params())
				if err != nil {
					t.Fatalf("Upgrade failed: %v", err)
				}
				test.verify(ctx, t, db)
			})
		}
	})

}

func verifyV2Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	amgr, _, err := Open(ctx, db, chaincfg.TestNet3Params(), pubPass)
	if err != nil {
		t.Fatalf("Open after Upgrade failed: %v", err)
	}

	err = walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(waddrmgrBucketKey)
		nsMetaBucket := ns.NestedReadBucket(metaBucketName)

		accounts := []struct {
			totalAddrs uint32
			lastUsed   uint32
		}{
			{^uint32(0), ^uint32(0)},
			{20, 18},
			{20, 19},
			{20, 19},
			{30, 25},
			{30, 29},
			{30, 29},
			{200, 185},
			{200, 199},
		}

		switch lastAccount, err := fetchLastAccount(ns); {
		case err != nil:
			t.Errorf("fetchLastAccount: %v", err)
		case lastAccount != uint32(len(accounts)-1):
			t.Errorf("Number of BIP0044 accounts got %v want %v",
				lastAccount+1, uint32(len(accounts)))
		}

		for i, a := range accounts {
			account := uint32(i)

			if nsMetaBucket.Get(accountNumberToAddrPoolKey(false, account)) != nil {
				t.Errorf("Account %v external address pool bucket still exists", account)
			}
			if nsMetaBucket.Get(accountNumberToAddrPoolKey(true, account)) != nil {
				t.Errorf("Account %v external address pool bucket still exists", account)
			}

			props, err := amgr.AccountProperties(ns, account)
			if err != nil {
				t.Errorf("AccountProperties: %v", err)
				continue
			}
			if props.LastUsedExternalIndex != a.lastUsed {
				t.Errorf("Account %v last used ext index got %v want %v",
					account, props.LastUsedExternalIndex, a.lastUsed)
			}
			if props.LastUsedInternalIndex != a.lastUsed {
				t.Errorf("Account %v last used int index got %v want %v",
					account, props.LastUsedInternalIndex, a.lastUsed)
			}
		}

		if ns.NestedReadBucket(usedAddrBucketName) != nil {
			t.Error("Used address bucket still exists")
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func verifyV3Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	_, _, err := Open(ctx, db, chaincfg.TestNet3Params(), pubPass)
	if err != nil {
		t.Fatalf("Open after Upgrade failed: %v", err)
	}

	err = walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		// This upgrade originally updated every ticket in the "wstakemgr" bucket,
		// however that bucket became redundant due to the removal of legacy
		// stakepool functionality. Therefore, that part of the upgrade has been
		// removed, and the only thing remaining is the creation of a new
		// "ticketsagendaprefs" bucket.

		// Verify that the agenda preferences bucket was created.
		if tx.ReadBucket(agendaPreferences.defaultBucketKey()) == nil {
			t.Errorf("Agenda preferences bucket was not created")
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func verifyV4Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	err := walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(waddrmgrBucketKey)
		mainBucket := ns.NestedReadBucket(mainBucketName)
		if mainBucket.Get(seedName) != nil {
			t.Errorf("Seed was not deleted")
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func verifyV5Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	err := walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(waddrmgrBucketKey)

		data := []struct {
			acct             uint32
			lastUsedExtChild uint32
			lastUsedIntChild uint32
		}{
			{0, ^uint32(0), ^uint32(0)},
			{1, 0, 0},
			{2, 9, 9},
			{3, 5, 15},
			{4, 19, 20},
			{5, 20, 19},
			{6, 29, 30},
			{7, 30, 29},
			{8, 1<<31 - 1, 1<<31 - 1},
			{ImportedAddrAccount, 0, 0},
		}

		const dbVersion = 5

		for _, d := range data {
			acct, err := fetchDBAccount(ns, d.acct, dbVersion)
			if err != nil {
				return err
			}
			a, ok := acct.(*dbBIP0044Account)
			if !ok {
				return fmt.Errorf("unknown account type %T", acct)
			}
			if a.lastUsedExternalIndex != d.lastUsedExtChild {
				t.Errorf("Account %d last used ext child mismatch %d != %d",
					d.acct, a.lastUsedExternalIndex, d.lastUsedExtChild)
			}
			if a.lastReturnedExternalIndex != d.lastUsedExtChild {
				t.Errorf("Account %d last returned ext child mismatch %d != %d",
					d.acct, a.lastReturnedExternalIndex, d.lastUsedExtChild)
			}
			if a.lastUsedInternalIndex != d.lastUsedIntChild {
				t.Errorf("Account %d last used int child mismatch %d != %d",
					d.acct, a.lastUsedInternalIndex, d.lastUsedIntChild)
			}
			if a.lastReturnedInternalIndex != d.lastUsedIntChild {
				t.Errorf("Account %d last returned int child mismatch %d != %d",
					d.acct, a.lastReturnedInternalIndex, d.lastUsedIntChild)
			}
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func verifyV6Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	err := walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(wtxmgrBucketKey)

		data := []*chainhash.Hash{
			decodeHash("7bc19eb0bf3a57be73d6879b6c411404b14b0156353dd47c5e0456768704bfd1"),
			decodeHash("a6abeb0127c347b5f38ebc2401134b324612d5b1ad9a9b8bdf6a91521842b7b1"),
			decodeHash("1107757fa4f238803192c617c7b60bf35bdc57bc0fc94b408c71239ff9eaeb98"),
			decodeHash("3fd00cda28c4d148e0cd38e1d646ba1365116b3ddd9a49aca4483bef80513ff9"),
			decodeHash("f4bdebefaa174470182960046fa53f554108b8ea09a86de5306a14c3a0124566"),
			decodeHash("bca8c2649860585f10b27d774b354ea7b80007e9ad79c090ea05596d63995cf5"),
		}

		c := ns.NestedReadBucket(bucketTickets).ReadCursor()
		found := 0
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var hash chainhash.Hash
			copy(hash[:], k)
			var foundHash *chainhash.Hash
			for _, foundHash = range data {
				if hash == *foundHash {
					goto Found
				}
			}
			t.Errorf("tickets bucket records %v as a ticket", &hash)
			continue
		Found:
			found++
			if extractRawTicketPickedHeight(v) != -1 {
				t.Errorf("ticket purchase %v was not set with picked height -1", foundHash)
			}
		}
		if found != len(data) {
			t.Errorf("missing ticket purchase transactions from tickets bucket")
		}

		// Ensure that the stakebase input recorded for an unmined vote was
		// removed.
		stakebaseKey := canonicalOutPoint(&chainhash.Hash{}, ^uint32(0))
		if ns.NestedReadBucket(bucketUnminedInputs).Get(stakebaseKey) != nil {
			t.Errorf("stakebase input for unmined vote was not removed")
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func verifyV8Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	err := walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		ns := tx.ReadBucket(wtxmgrBucketKey)
		creditBucket := ns.NestedReadBucket(bucketCredits)
		err := creditBucket.ForEach(func(k []byte, v []byte) error {
			hasExpiry := fetchRawCreditHasExpiry(v, DBVersion)
			if !hasExpiry {
				t.Errorf("expected expiry to be set")
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}

		unmineVGLeditBucket := ns.NestedReadBucket(bucketUnmineVGLedits)
		err = unmineVGLeditBucket.ForEach(func(k []byte, v []byte) error {
			hasExpiry := fetchRawCreditHasExpiry(v, DBVersion)

			if !hasExpiry {
				t.Errorf("expected expiry to be set")
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}

		txBucket := ns.NestedReadBucket(bucketTxRecords)
		minedTxWithExpiryCount := 0
		minedTxWithoutExpiryCount := 0
		err = txBucket.ForEach(func(k []byte, v []byte) error {
			var txHash chainhash.Hash
			var rec TxRecord
			err := readRawTxRecordHash(k, &txHash)
			if err != nil {
				t.Error(err)
			}
			err = readRawTxRecord(&txHash, v, &rec)
			if err != nil {
				t.Error(err)
			}

			if rec.MsgTx.Expiry != wire.NoExpiryValue {
				minedTxWithExpiryCount++
			} else {
				minedTxWithoutExpiryCount++
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}

		if minedTxWithExpiryCount != 3 {
			t.Errorf("expected 3 txs with expiries set, got %d", minedTxWithExpiryCount)
		}
		if minedTxWithoutExpiryCount != 3 {
			t.Errorf("expected 3 txs without expiries set, got %d", minedTxWithoutExpiryCount)
		}
		return err
	})
	if err != nil {
		t.Error(err)
	}
}

// verifyV12Upgrade tests whether the upgrade to the v12 database was
// successful, using the v11 test database.
//
// See the v11.db.go file for an explanation of the database layout and test
// plan.
func verifyV12Upgrade(ctx context.Context, t *testing.T, db walletdb.DB) {
	_, txmgr, err := Open(ctx, db, chaincfg.TestNet3Params(), pubPass)
	if err != nil {
		t.Fatalf("Open after Upgrade failed: %v", err)
	}

	err = walletdb.View(ctx, db, func(tx walletdb.ReadTx) error {
		txmgrns := tx.ReadBucket(wtxmgrBucketKey)

		if b := txmgrns.NestedReadBucket(bucketTicketCommitments); b == nil {
			t.Fatalf("upgrade should have created bucketTicketCommitments")
		}

		if b := txmgrns.NestedReadBucket(bucketTicketCommitmentsUsp); b == nil {
			t.Fatalf("upgrade should have created bucketTicketCommitmentsUsp")
		}

		balances, err := txmgr.AccountBalances(tx, 0)
		if err != nil {
			t.Fatal(err)
		}

		expectedBalances := []struct {
			acct             uint32
			spendable        VGLutil.Amount
			votingAuth       VGLutil.Amount
			total            VGLutil.Amount
			unconfirmed      VGLutil.Amount
			locked           VGLutil.Amount
			immatureStakeGen VGLutil.Amount
			empty            bool
		}{
			// unmined ticket
			{acct: 1, votingAuth: 1100},
			{acct: 2, locked: 1000, total: 1000},

			// mined ticket
			{acct: 3, votingAuth: 1100},
			{acct: 4, locked: 1000, total: 1000},

			// mined ticket + unmined vote
			{acct: 5, empty: true},
			{acct: 6, total: 1300, immatureStakeGen: 1300},

			// mined ticket + mined vote
			{acct: 7, empty: true},
			{acct: 8, total: 1300, immatureStakeGen: 1300},

			// mined ticket + unmined revocation
			{acct: 9, empty: true},
			{acct: 10, total: 700, immatureStakeGen: 700},

			// mined ticket + mined revocation
			{acct: 11, empty: true},
			{acct: 12, total: 700, immatureStakeGen: 700},
		}

		testFunc := func(testIdx int) func(t *testing.T) {
			return func(t *testing.T) {
				expected := expectedBalances[testIdx]
				actual, has := balances[expected.acct]

				if expected.empty {
					if !has {
						// this account was actually supposed to be empty
						return
					}
					t.Fatalf("Balance should have been empty")
				}
				if !has {
					t.Fatalf("Database does not have balance for expected account")
				}

				if actual.Spendable != expected.spendable {
					t.Errorf("Actual spendable (%d) different than expected (%d)",
						actual.Spendable, expected.spendable)
				}
				if actual.Unconfirmed != expected.unconfirmed {
					t.Errorf("Actual unconfirmed (%d) different than expected (%d)",
						actual.Unconfirmed, expected.unconfirmed)
				}
				if actual.LockedByTickets != expected.locked {
					t.Errorf("Actual locked by tickets (%d) different than expected (%d)",
						actual.LockedByTickets, expected.locked)
				}
				if actual.ImmatureStakeGeneration != expected.immatureStakeGen {
					t.Errorf("Actual immature stake gen (%d) different than expected (%d)",
						actual.ImmatureStakeGeneration, expected.immatureStakeGen)
				}
				if actual.VotingAuthority != expected.votingAuth {
					t.Errorf("Actual voting authority (%d) different than expected (%d)",
						actual.VotingAuthority, expected.votingAuth)
				}
				if actual.Total != expected.total {
					t.Errorf("Actual total (%d) different than expected (%d)",
						actual.Total, expected.total)
				}
			}
		}

		for i, e := range expectedBalances {
			t.Run(fmt.Sprintf("acct=%d", e.acct), testFunc(i))
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
