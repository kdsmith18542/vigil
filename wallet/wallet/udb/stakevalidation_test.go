// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package udb

import (
	"bytes"
	"context"
	"testing"
	"time"

	_ "github.com/kdsmith18542/vigil/wallet/wallet/internal/bdb"
	"github.com/kdsmith18542/vigil/wallet/wallet/walletdb"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
	gcs2 "github.com/kdsmith18542/vigil/gcs/v4"
	"github.com/kdsmith18542/vigil/wire"
)

func insertMainChainHeaders(s *Store, dbtx walletdb.ReadWriteTx,
	headerData []BlockHeaderData, filters []*gcs2.FilterV2) error {

	ns := dbtx.ReadWriteBucket(wtxmgrBucketKey)

	for i := range headerData {
		h := &headerData[i]
		f := filters[i]
		header := new(wire.BlockHeader)
		err := header.Deserialize(bytes.NewReader(h.SerializedHeader[:]))
		if err != nil {
			return err
		}
		blockHash := header.BlockHash()
		err = s.ExtendMainChain(ns, header, &blockHash, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestStakeInvalidationOfTip(t *testing.T) {
	ctx := context.Background()
	db, _, s, teardown, err := cloneDB(ctx, "stake_inv_of_tip.kv")
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	g := makeBlockGenerator()
	block1Header := g.generate(VGLutil.BlockValid)
	block2Header := g.generate(VGLutil.BlockValid)
	block3Header := g.generate(0)

	block1Tx := wire.MsgTx{
		TxOut: []*wire.TxOut{{Value: 2e8}},
	}
	block2Tx := wire.MsgTx{
		TxIn: []*wire.TxIn{
			{
				PreviousOutPoint: wire.OutPoint{
					Hash:  block1Tx.TxHash(),
					Index: 0,
					Tree:  0,
				},
			},
		},
		TxOut: []*wire.TxOut{{Value: 1e8}},
	}
	block1TxRec, err := NewTxRecordFromMsgTx(&block1Tx, time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	block2TxRec, err := NewTxRecordFromMsgTx(&block2Tx, time.Time{})
	if err != nil {
		t.Fatal(err)
	}

	err = walletdb.Update(ctx, db, func(dbtx walletdb.ReadWriteTx) error {
		err := s.InsertMemPoolTx(dbtx, block1TxRec)
		if err != nil {
			return err
		}
		err = s.AdVGLedit(dbtx, block1TxRec, nil, 0, false, 0)
		if err != nil {
			return err
		}
		err = s.InsertMemPoolTx(dbtx, block2TxRec)
		if err != nil {
			return err
		}
		err = s.AdVGLedit(dbtx, block2TxRec, nil, 0, false, 0)
		if err != nil {
			return err
		}

		bal, err := s.AccountBalance(dbtx, 0, 0)
		if err != nil {
			return err
		}
		if bal.Total != 1e8 {
			t.Errorf("Wrong balance before mining either transaction: %v", bal)
		}

		headerData := makeHeaderDataSlice(block1Header, block2Header)
		filters := emptyFilters(2)
		err = insertMainChainHeaders(s, dbtx, headerData, filters)
		if err != nil {
			return err
		}

		err = s.InsertMinedTx(dbtx, block1TxRec, &headerData[0].BlockHash)
		if err != nil {
			return err
		}
		err = s.InsertMinedTx(dbtx, block2TxRec, &headerData[1].BlockHash)
		if err != nil {
			return err
		}

		// At this point there should only be one credit for the tx in block 2.
		bal, err = s.AccountBalance(dbtx, 1, 0)
		if err != nil {
			return err
		}
		if bal.Total != VGLutil.Amount(block2Tx.TxOut[0].Value) {
			t.Errorf("Wrong balance: expected %v got %v",
				VGLutil.Amount(block2Tx.TxOut[0].Value), bal)
		}
		credits, err := s.UnspentOutputs(dbtx)
		if err != nil {
			return err
		}
		if len(credits) != 1 {
			t.Errorf("Expected only 1 credit, got %v", len(credits))
			return nil
		}
		if credits[0].Hash != block2Tx.TxHash() {
			t.Errorf("Credit hash does match tx from block 2")
			return nil
		}
		if credits[0].Amount != VGLutil.Amount(block2Tx.TxOut[0].Value) {
			t.Errorf("Credit value does not match tx output 0 from block 2")
			return nil
		}

		// Add the next block header which invalidates the regular tx tree of
		// block 2.
		t.Log("Invalidating block 2")
		headerData = makeHeaderDataSlice(block3Header)
		filters = emptyFilters(1)
		err = insertMainChainHeaders(s, dbtx, headerData, filters)
		if err != nil {
			return err
		}

		// Now the transaction in block 2 is invalidated.  There should only be
		// one unspent output, from block 1.
		bal, err = s.AccountBalance(dbtx, 1, 0)
		if err != nil {
			return err
		}
		if bal.Total != VGLutil.Amount(block1Tx.TxOut[0].Value) {
			t.Errorf("Wrong balance: expected %v got %v", VGLutil.Amount(block1Tx.TxOut[0].Value), bal)
		}
		credits, err = s.UnspentOutputs(dbtx)
		if err != nil {
			return err
		}
		if len(credits) != 1 {
			t.Errorf("Expected only 1 credit, got %v", len(credits))
			return nil
		}
		if credits[0].Hash != block1Tx.TxHash() {
			t.Errorf("Credit hash does not match tx from block 1")
			return nil
		}
		if credits[0].Amount != VGLutil.Amount(block1Tx.TxOut[0].Value) {
			t.Errorf("Credit value does not match tx output 0 from block 1")
			return nil
		}

		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
