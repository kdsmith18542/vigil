// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package VGLutil

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	"github.com/kdsmith18542/vigil/wire"
)

// assertTransactionImmutability throws a panic when a transaction has been
// mutated.
var assertTransactionImmutability = false

// TxIndexUnknown is the value returned for a transaction index that is unknown.
// This is typically because the transaction has not been inserted into a block
// yet.
const TxIndexUnknown = -1

// Tx defines a transaction that provides easier and more efficient manipulation
// of raw transactions.  It also memoizes the hash for the transaction on its
// first access so subsequent accesses don't have to repeat the relatively
// expensive hashing operations.
type Tx struct {
	hash    chainhash.Hash // Cached transaction hash
	msgTx   *wire.MsgTx    // Underlying MsgTx
	txTree  int8           // Indicates which tx tree the tx is found in
	txIndex int            // Position within a block or TxIndexUnknown
}

// MsgTx returns the underlying wire.MsgTx for the transaction.
func (t *Tx) MsgTx() *wire.MsgTx {
	// Return the cached transaction.
	return t.msgTx
}

// Hash returns the hash of the transaction.  This is equivalent to
// calling TxHash on the underlying wire.MsgTx, however it caches the
// result so subsequent calls are more efficient.
func (t *Tx) Hash() *chainhash.Hash {
	if assertTransactionImmutability {
		hash := t.msgTx.TxHash()
		if !hash.IsEqual(&t.hash) {
			str := fmt.Sprintf("ASSERT: mutated util.tx detected, old hash %v, "+
				"new hash %v",
				t.hash,
				hash)
			panic(str)
		}
	}
	return &t.hash
}

// Index returns the saved index of the transaction within a block.  This value
// will be TxIndexUnknown if it hasn't already explicitly been set.
func (t *Tx) Index() int {
	return t.txIndex
}

// SetIndex sets the index of the transaction in within a block.
func (t *Tx) SetIndex(index int) {
	t.txIndex = index
}

// Tree returns the saved tree of the transaction within a block.  This value
// will be TxTreeUnknown if it hasn't already explicitly been set.
func (t *Tx) Tree() int8 {
	return t.txTree
}

// SetTree sets the tree of the transaction in within a block.
func (t *Tx) SetTree(tree int8) {
	t.txTree = tree
}

// Version returns the transaction version.
func (t *Tx) Version() uint16 {
	return t.msgTx.Version
}

// SetVersion sets the version of the transaction within a block.
func (t *Tx) SetVersion(version uint16) {
	t.msgTx.Version = version
}

// NewTx returns a new instance of a transaction given an underlying
// wire.MsgTx.  See Tx.
func NewTx(msgTx *wire.MsgTx) *Tx {
	return &Tx{
		hash:    msgTx.TxHash(),
		msgTx:   msgTx,
		txTree:  wire.TxTreeUnknown,
		txIndex: TxIndexUnknown,
	}
}

// NewTxDeep returns a new instance of a transaction given an underlying
// wire.MsgTx.  Until NewTx, it completely copies the data in the msgTx
// so that there are new memory allocations, in case you were to somewhere
// else modify the data assigned to these pointers.
func NewTxDeep(msgTx *wire.MsgTx) *Tx {
	txIns := make([]*wire.TxIn, len(msgTx.TxIn))
	txOuts := make([]*wire.TxOut, len(msgTx.TxOut))

	for i, txin := range msgTx.TxIn {
		sigScript := make([]byte, len(txin.SignatureScript))
		copy(sigScript[:], txin.SignatureScript[:])

		txIns[i] = &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  txin.PreviousOutPoint.Hash,
				Index: txin.PreviousOutPoint.Index,
				Tree:  txin.PreviousOutPoint.Tree,
			},
			Sequence:        txin.Sequence,
			ValueIn:         txin.ValueIn,
			BlockHeight:     txin.BlockHeight,
			BlockIndex:      txin.BlockIndex,
			SignatureScript: sigScript,
		}
	}

	for i, txout := range msgTx.TxOut {
		pkScript := make([]byte, len(txout.PkScript))
		copy(pkScript[:], txout.PkScript[:])

		txOuts[i] = &wire.TxOut{
			Value:    txout.Value,
			Version:  txout.Version,
			PkScript: pkScript,
		}
	}

	mtx := &wire.MsgTx{
		CachedHash: nil,
		SerType:    msgTx.SerType,
		Version:    msgTx.Version,
		TxIn:       txIns,
		TxOut:      txOuts,
		LockTime:   msgTx.LockTime,
		Expiry:     msgTx.Expiry,
	}

	return &Tx{
		hash:    mtx.TxHash(),
		msgTx:   mtx,
		txTree:  wire.TxTreeUnknown,
		txIndex: TxIndexUnknown,
	}
}

// NewTxDeepTxIns is used to deep copy a transaction, maintaining the old
// pointers to the TxOuts while replacing the old pointers to the TxIns with
// deep copies. This is to prevent races when the fraud proofs for the
// transactions are set by the miner.
func NewTxDeepTxIns(tx *Tx) *Tx {
	if tx == nil {
		return nil
	}

	msgTx := tx.msgTx
	txIns := make([]*wire.TxIn, len(msgTx.TxIn))
	txOuts := make([]*wire.TxOut, len(msgTx.TxOut))

	for i, txin := range msgTx.TxIn {
		sigScript := make([]byte, len(txin.SignatureScript))
		copy(sigScript[:], txin.SignatureScript[:])

		txIns[i] = &wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  txin.PreviousOutPoint.Hash,
				Index: txin.PreviousOutPoint.Index,
				Tree:  txin.PreviousOutPoint.Tree,
			},
			Sequence:        txin.Sequence,
			ValueIn:         txin.ValueIn,
			BlockHeight:     txin.BlockHeight,
			BlockIndex:      txin.BlockIndex,
			SignatureScript: sigScript,
		}
	}

	// Shallow copy the TxOuts.
	copy(txOuts, msgTx.TxOut)

	newTxMsg := &wire.MsgTx{
		CachedHash: nil,
		SerType:    msgTx.SerType,
		Version:    msgTx.Version,
		TxIn:       txIns,
		TxOut:      txOuts,
		LockTime:   msgTx.LockTime,
		Expiry:     msgTx.Expiry,
	}

	return &Tx{
		hash:    newTxMsg.TxHash(),
		msgTx:   newTxMsg,
		txTree:  tx.txTree,
		txIndex: tx.txIndex,
	}
}

// NewTxFromBytes returns a new instance of a transaction given the
// serialized bytes.  See Tx.
func NewTxFromBytes(serializedTx []byte) (*Tx, error) {
	br := bytes.NewReader(serializedTx)
	return NewTxFromReader(br)
}

// NewTxFromReader returns a new instance of a transaction given a
// Reader to deserialize the transaction.  See Tx.
func NewTxFromReader(r io.Reader) (*Tx, error) {
	// Deserialize the bytes into a MsgTx.
	var msgTx wire.MsgTx
	err := msgTx.Deserialize(r)
	if err != nil {
		return nil, err
	}

	t := Tx{
		hash:    msgTx.TxHash(),
		msgTx:   &msgTx,
		txTree:  wire.TxTreeUnknown,
		txIndex: TxIndexUnknown,
	}

	return &t, nil
}
