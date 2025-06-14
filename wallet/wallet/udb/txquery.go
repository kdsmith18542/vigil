// Copyright (c) 2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package udb

import (
	"bytes"
	"context"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/wallet/walletdb"
	"github.com/kdsmith18542/vigil/blockchain/stake/v5"
	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
	"github.com/kdsmith18542/vigil/wire"
)

// CreditRecord contains metadata regarding a transaction credit for a known
// transaction.  Further details may be looked up by indexing a wire.MsgTx.TxOut
// with the Index field.
type CreditRecord struct {
	Index      uint32
	Amount     VGLutil.Amount
	Spent      bool
	Change     bool
	OpCode     uint8
	IsCoinbase bool
	HasExpiry  bool
}

// DebitRecord contains metadata regarding a transaction debit for a known
// transaction.  Further details may be looked up by indexing a wire.MsgTx.TxIn
// with the Index field.
type DebitRecord struct {
	Amount VGLutil.Amount
	Index  uint32
}

// TxDetails is intended to provide callers with access to rich details
// regarding a relevant transaction and which inputs and outputs are credit or
// debits.
type TxDetails struct {
	TxRecord
	Block   BlockMeta
	Credits []CreditRecord
	Debits  []DebitRecord
}

// Height returns the height of a transaction according to the BlockMeta.
func (t *TxDetails) Height() int32 {
	return t.Block.Block.Height
}

// minedTxDetails fetches the TxDetails for the mined transaction with hash
// txHash and the passed tx record key and value.
func (s *Store) minedTxDetails(ns walletdb.ReadBucket, txHash *chainhash.Hash, recKey, recVal []byte) (*TxDetails, error) {
	var details TxDetails

	// Parse transaction record k/v, lookup the full block record for the
	// block time, and read all matching credits, debits.
	err := readRawTxRecord(txHash, recVal, &details.TxRecord)
	if err != nil {
		return nil, err
	}
	err = readRawTxRecordBlock(recKey, &details.Block.Block)
	if err != nil {
		return nil, err
	}
	details.Block.Time, err = fetchBlockTime(ns, details.Block.Height)
	if err != nil {
		return nil, err
	}

	credIter := makeReaVGLeditIterator(ns, recKey, DBVersion)
	for credIter.next() {
		if int(credIter.elem.Index) >= len(details.MsgTx.TxOut) {
			credIter.close()
			return nil, errors.E(errors.IO, "saved credit index exceeds number of outputs")
		}

		// The credit iterator does not record whether this credit was
		// spent by an unmined transaction, so check that here.
		if !credIter.elem.Spent {
			k := canonicalOutPoint(txHash, credIter.elem.Index)
			spent := existsRawUnminedInput(ns, k) != nil
			credIter.elem.Spent = spent
		}
		details.Credits = append(details.Credits, credIter.elem)
	}
	credIter.close()
	if credIter.err != nil {
		return nil, credIter.err
	}

	debIter := makeReadDebitIterator(ns, recKey)
	defer debIter.close()
	for debIter.next() {
		if int(debIter.elem.Index) >= len(details.MsgTx.TxIn) {
			return nil, errors.E(errors.IO, "saved debit index exceeds number of inputs")
		}

		details.Debits = append(details.Debits, debIter.elem)
	}
	return &details, debIter.err
}

// unminedTxDetails fetches the TxDetails for the unmined transaction with the
// hash txHash and the passed unmined record value.
func (s *Store) unminedTxDetails(ns walletdb.ReadBucket, txHash *chainhash.Hash, v []byte) (*TxDetails, error) {
	details := TxDetails{
		Block: BlockMeta{Block: Block{Height: -1}},
	}
	err := readRawTxRecord(txHash, v, &details.TxRecord)
	if err != nil {
		return nil, err
	}

	it := makeReadUnmineVGLeditIterator(ns, txHash, DBVersion)
	defer it.close()
	for it.next() {
		if int(it.elem.Index) >= len(details.MsgTx.TxOut) {
			return nil, errors.E(errors.IO, errors.Errorf("credit output index %d does not exist for tx %v", it.elem.Index, txHash))
		}

		// Set the Spent field since this is not done by the iterator.
		it.elem.Spent = existsRawUnminedInput(ns, it.ck) != nil
		details.Credits = append(details.Credits, it.elem)
	}
	if it.err != nil {
		return nil, it.err
	}

	// Debit records are not saved for unmined transactions.  Instead, they
	// must be looked up for each transaction input manually.  There are two
	// kinds of previous credits that may be debited by an unmined
	// transaction: mined unspent outputs (which remain marked unspent even
	// when spent by an unmined transaction), and credits from other unmined
	// transactions.  Both situations must be considered.
	for i, output := range details.MsgTx.TxIn {
		opKey := canonicalOutPoint(&output.PreviousOutPoint.Hash,
			output.PreviousOutPoint.Index)
		credKey := existsRawUnspent(ns, opKey)
		if credKey != nil {
			v := existsRawCredit(ns, credKey)
			amount, err := fetchRawCreditAmount(v)
			if err != nil {
				return nil, err
			}

			details.Debits = append(details.Debits, DebitRecord{
				Amount: amount,
				Index:  uint32(i),
			})
			continue
		}

		v := existsRawUnmineVGLedit(ns, opKey)
		if v == nil {
			continue
		}

		amount, err := fetchRawCreditAmount(v)
		if err != nil {
			return nil, err
		}
		details.Debits = append(details.Debits, DebitRecord{
			Amount: amount,
			Index:  uint32(i),
		})
	}

	return &details, nil
}

// TxDetails looks up all recorded details regarding a transaction with some
// hash.  In case of a hash collision, the most recent transaction with a
// matching hash is returned.
func (s *Store) TxDetails(ns walletdb.ReadBucket, txHash *chainhash.Hash) (*TxDetails, error) {
	// First, check whether there exists an unmined transaction with this
	// hash.  Use it if found.
	v := existsRawUnmined(ns, txHash[:])
	if v != nil {
		return s.unminedTxDetails(ns, txHash, v)
	}

	// Otherwise, if there exists a mined transaction with this matching
	// hash, skip over to the newest and begin fetching all details.
	k, v := latestTxRecord(ns, txHash[:])
	if v == nil {
		return nil, errors.E(errors.NotExist)
	}
	return s.minedTxDetails(ns, txHash, k, v)
}

// TicketDetails is intended to provide callers with access to rich details
// regarding a relevant transaction and which inputs and outputs are credit or
// debits.
type TicketDetails struct {
	Ticket  *TxDetails
	Spender *TxDetails
}

// TicketDetails looks up all recorded details regarding a ticket with some
// hash.
//
// Not finding a ticket with this hash is not an error.  In this case, a nil
// TicketDetails is returned.
func (s *Store) TicketDetails(ns walletdb.ReadBucket, txDetails *TxDetails) (*TicketDetails, error) {
	var ticketDetails = &TicketDetails{}
	if !stake.IsSStx(&txDetails.MsgTx) {
		return nil, nil
	}
	ticketDetails.Ticket = txDetails
	var spenderHash = chainhash.Hash{}
	// Check if the ticket is spent or not.  Look up the credit for output 0
	// and check if either a debit is recorded or the output is spent by an
	// unmined transaction.
	_, credVal := existsCredit(ns, &txDetails.Hash, 0, &txDetails.Block.Block)
	if credVal != nil {
		if extractRawCreditIsSpent(credVal) {
			debKey := extractRawCreditSpenderDebitKey(credVal)
			debHash := extractRawDebitHash(debKey)
			copy(spenderHash[:], debHash)
		}
	} else {
		opKey := canonicalOutPoint(&txDetails.Hash, 0)
		spenderVal := existsRawUnminedInput(ns, opKey)
		if spenderVal != nil {
			copy(spenderHash[:], spenderVal)
		}
	}
	spenderDetails, err := s.TxDetails(ns, &spenderHash)
	if (err != nil) && (!errors.Is(err, errors.NotExist)) {
		return nil, err
	}
	ticketDetails.Spender = spenderDetails
	return ticketDetails, nil
}

// parseTx deserializes a transaction into a MsgTx using the readRawTxRecord
// method.
func (s *Store) parseTx(txHash chainhash.Hash, v []byte) (*wire.MsgTx, error) {
	details := TxDetails{
		Block: BlockMeta{Block: Block{Height: -1}},
	}
	err := readRawTxRecord(&txHash, v, &details.TxRecord)
	if err != nil {
		return nil, err
	}

	return &details.MsgTx, nil
}

// Tx looks up all the stored wire.MsgTx for a transaction with some
// hash.  In case of a hash collision, the most recent transaction with a
// matching hash is returned.
func (s *Store) Tx(ns walletdb.ReadBucket, txHash *chainhash.Hash) (*wire.MsgTx, error) {
	// First, check whether there exists an unmined transaction with this
	// hash.  Use it if found.
	v := existsRawUnmined(ns, txHash[:])
	if v != nil {
		return s.parseTx(*txHash, v)
	}

	// Otherwise, if there exists a mined transaction with this matching
	// hash, skip over to the newest and begin fetching the msgTx.
	_, v = latestTxRecord(ns, txHash[:])
	if v == nil {
		return nil, errors.E(errors.NotExist,
			errors.Errorf("tx %s not found", txHash.String()))
	}
	return s.parseTx(*txHash, v)
}

// ExistsTx checks to see if a transaction exists in the database.
func (s *Store) ExistsTx(ns walletdb.ReadBucket, txHash *chainhash.Hash) bool {
	mined, unmined := s.ExistsTxMinedOrUnmined(ns, txHash)
	return mined || unmined
}

// ExistsTxMinedOrUnmined checks if a transaction is recorded as a mined or
// unmined transaction.
func (s *Store) ExistsTxMinedOrUnmined(ns walletdb.ReadBucket, txHash *chainhash.Hash) (mined, unmined bool) {
	v := existsRawUnmined(ns, txHash[:])
	if v != nil {
		return false, true
	}
	_, v = latestTxRecord(ns, txHash[:])
	return v != nil, false
}

// ExistsUTXO checks to see if op refers to an unspent transaction output or a
// credit spent by an unmined transaction.  This check is sufficient to
// determine whether a transaction input is relevant to the wallet by spending a
// UTXO or conflicting with another mempool transaction that double spends the
// output.
func (s *Store) ExistsUTXO(dbtx walletdb.ReadTx, op *wire.OutPoint) bool {
	ns := dbtx.ReadBucket(wtxmgrBucketKey)
	k, v := existsUnspent(ns, op)
	if v != nil {
		return true
	}
	return existsRawUnmineVGLedit(ns, k) != nil
}

// UniqueTxDetails looks up all recorded details for a transaction recorded
// mined in some particular block, or an unmined transaction if block is nil.
//
// Not finding a transaction with this hash from this block is not an error.  In
// this case, a nil TxDetails is returned.
func (s *Store) UniqueTxDetails(ns walletdb.ReadBucket, txHash *chainhash.Hash,
	block *Block) (*TxDetails, error) {

	if block == nil {
		v := existsRawUnmined(ns, txHash[:])
		if v == nil {
			return nil, nil
		}
		return s.unminedTxDetails(ns, txHash, v)
	}

	k, v := existsTxRecord(ns, txHash, block)
	if v == nil {
		return nil, nil
	}
	return s.minedTxDetails(ns, txHash, k, v)
}

// TxBlockHeight returns the block height of a mined transaction, or -1 for any
// unmined transactions.
func (s *Store) TxBlockHeight(dbtx walletdb.ReadTx, txHash *chainhash.Hash) (int32, error) {
	ns := dbtx.ReadBucket(wtxmgrBucketKey)
	v := existsRawUnmined(ns, txHash[:])
	if v != nil {
		return -1, nil
	}
	k, _ := latestTxRecord(ns, txHash[:])
	if k == nil {
		return 0, errors.E(errors.NotExist, errors.Errorf("no transaction %v", txHash))
	}
	var height int32
	err := readRawTxRecordBlockHeight(k, &height)
	return height, err
}

// rangeUnminedTransactions executes the function f with TxDetails for every
// unmined transaction.  f is not executed if no unmined transactions exist.
// Error returns from f (if any) are propigated to the caller.  Returns true
// (signaling breaking out of a RangeTransactions) iff f executes and returns
// true.
func (s *Store) rangeUnminedTransactions(ctx context.Context, ns walletdb.ReadBucket, f func([]TxDetails) (bool, error)) (bool, error) {
	var details []TxDetails
	err := ns.NestedReadBucket(bucketUnmined).ForEach(func(k, v []byte) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if len(k) < 32 {
			return errors.E(errors.IO, errors.Errorf("bad unmined tx key len %d", len(k)))
		}

		var txHash chainhash.Hash
		copy(txHash[:], k)
		detail, err := s.unminedTxDetails(ns, &txHash, v)
		if err != nil {
			return err
		}

		// Because the key was created while foreach-ing over the
		// bucket, it should be impossible for unminedTxDetails to ever
		// successfully return a nil details struct.
		details = append(details, *detail)
		return nil
	})
	if err == nil && len(details) > 0 {
		return f(details)
	}
	return false, err
}

// rangeBlockTransactions executes the function f with TxDetails for every block
// between heights begin and end (reverse order when end > begin) until f
// returns true, or the transactions from block is processed.  Returns true iff
// f executes and returns true.
func (s *Store) rangeBlockTransactions(ctx context.Context, ns walletdb.ReadBucket, begin, end int32,
	f func([]TxDetails) (bool, error)) (bool, error) {

	// Mempool height is considered a high bound.
	if begin < 0 {
		begin = int32(^uint32(0) >> 1)
	}
	if end < 0 {
		end = int32(^uint32(0) >> 1)
	}

	var blockIter blockIterator
	var advance func(*blockIterator) bool
	if begin < end {
		// Iterate in forwards order
		blockIter = makeReadBlockIterator(ns, begin)
		defer blockIter.close()
		advance = func(it *blockIterator) bool {
			if !it.next() {
				return false
			}
			return it.elem.Height <= end
		}
	} else {
		// Iterate in backwards order, from begin -> end.
		blockIter = makeReadBlockIterator(ns, begin)
		defer blockIter.close()
		advance = func(it *blockIterator) bool {
			if !it.prev() {
				return false
			}
			return end <= it.elem.Height
		}
	}

	var details []TxDetails
	for advance(&blockIter) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		block := &blockIter.elem

		if cap(details) < len(block.transactions) {
			details = make([]TxDetails, 0, len(block.transactions))
		} else {
			details = details[:0]
		}

		for _, txHash := range block.transactions {
			k := keyTxRecord(&txHash, &block.Block)
			v := existsRawTxRecord(ns, k)
			if v == nil {
				return false, errors.E(errors.IO, errors.Errorf("missing transaction %v for block %v", txHash, block.Height))
			}
			detail := TxDetails{
				Block: BlockMeta{
					Block: block.Block,
					Time:  block.Time,
				},
			}
			err := readRawTxRecord(&txHash, v, &detail.TxRecord)
			if err != nil {
				return false, err
			}

			credIter := makeReaVGLeditIterator(ns, k, DBVersion)
			for credIter.next() {
				if int(credIter.elem.Index) >= len(detail.MsgTx.TxOut) {
					credIter.close()
					return false, errors.E(errors.IO, "saved credit index exceeds number of outputs")
				}

				// The credit iterator does not record whether
				// this credit was spent by an unmined
				// transaction, so check that here.
				if !credIter.elem.Spent {
					k := canonicalOutPoint(&txHash, credIter.elem.Index)
					spent := existsRawUnminedInput(ns, k) != nil
					credIter.elem.Spent = spent
				}
				detail.Credits = append(detail.Credits, credIter.elem)
			}
			credIter.close()
			if credIter.err != nil {
				return false, credIter.err
			}

			debIter := makeReadDebitIterator(ns, k)
			defer debIter.close()
			for debIter.next() {
				if int(debIter.elem.Index) >= len(detail.MsgTx.TxIn) {
					return false, errors.E(errors.IO, "saved debit index exceeds number of inputs")
				}

				detail.Debits = append(detail.Debits, debIter.elem)
			}
			if debIter.err != nil {
				return false, debIter.err
			}

			details = append(details, detail)
		}

		// Vigil: Block records are saved even when no transactions are
		// included.  This is used to save the votebits from every
		// block.  This differs from btcwallet where every block must
		// have one transaction.  Since f may only be called when
		// len(details) > 0, this must be explicitly tested.
		if len(details) == 0 {
			continue
		}
		brk, err := f(details)
		if err != nil || brk {
			return brk, err
		}
	}
	return false, blockIter.err
}

// RangeTransactions runs the function f on all transaction details between
// blocks on the best chain over the height range [begin,end].  The special
// height -1 may be used to also include unmined transactions.  If the end
// height comes before the begin height, blocks are iterated in reverse order
// and unmined transactions (if any) are processed first.
//
// The function f may return an error which, if non-nil, is propagated to the
// caller.  Additionally, a boolean return value allows exiting the function
// early without reading any additional transactions early when true.
//
// All calls to f are guaranteed to be passed a slice with more than zero
// elements.  The slice may be reused for multiple blocks, so it is not safe to
// use it after the loop iteration it was acquired.
func (s *Store) RangeTransactions(ctx context.Context, ns walletdb.ReadBucket, begin, end int32,
	f func([]TxDetails) (bool, error)) error {

	var addedUnmined bool
	if begin < 0 {
		brk, err := s.rangeUnminedTransactions(ctx, ns, f)
		if err != nil || brk {
			return err
		}
		addedUnmined = true
	}

	brk, err := s.rangeBlockTransactions(ctx, ns, begin, end, f)
	if err == nil && !brk && !addedUnmined && end < 0 {
		_, err = s.rangeUnminedTransactions(ctx, ns, f)
	}
	return err
}

// PreviousPkScripts returns a slice of previous output scripts for each credit
// output this transaction record debits from.
func (s *Store) PreviousPkScripts(ns walletdb.ReadBucket, rec *TxRecord, block *Block) ([][]byte, error) {
	var pkScripts [][]byte

	if block == nil {
		for _, input := range rec.MsgTx.TxIn {
			prevOut := &input.PreviousOutPoint

			// Input may spend a previous unmined output, a
			// mined output (which would still be marked
			// unspent), or neither.

			v := existsRawUnmined(ns, prevOut.Hash[:])
			if v != nil {
				// Ensure a credit exists for this
				// unmined transaction before including
				// the output script.
				k := canonicalOutPoint(&prevOut.Hash, prevOut.Index)
				vUC := existsRawUnmineVGLedit(ns, k)
				if vUC == nil {
					continue
				}

				// If we encounter an error here, it likely means
				// we have a legacy outpoint. Ignore the error and
				// just let the scrPos be 0, which will trigger
				// whole transaction deserialization to retrieve
				// the script.
				scrPos := fetchRawUnmineVGLeditScriptOffset(vUC)
				scrLen := fetchRawUnmineVGLeditScriptLength(vUC)

				pkScript, err := fetchRawTxRecordPkScript(
					prevOut.Hash[:], v, prevOut.Index, scrPos, scrLen)
				if err != nil {
					return nil, err
				}
				pkScripts = append(pkScripts, pkScript)
				continue
			}

			_, credKey := existsUnspent(ns, prevOut)
			if credKey != nil {
				credVal := existsRawCredit(ns, credKey)
				if credVal == nil {
					return nil, errors.E(errors.IO, errors.Errorf("missing credit value for key %x", credKey))
				}

				// Legacy outputs in the credit bucket may be of the
				// wrong size.
				scrPos := fetchRawCreditScriptOffset(credVal)
				scrLen := fetchRawCreditScriptLength(credVal)

				k := extractRawCreditTxRecordKey(credKey)
				v = existsRawTxRecord(ns, k)
				pkScript, err := fetchRawTxRecordPkScript(k, v,
					prevOut.Index, scrPos, scrLen)
				if err != nil {
					return nil, err
				}
				pkScripts = append(pkScripts, pkScript)
			}
		}
	}

	recKey := keyTxRecord(&rec.Hash, block)
	it := makeReadDebitIterator(ns, recKey)
	for it.next() {
		credKey := extractRawDebitCreditKey(it.cv)
		index := extractRawCreditIndex(credKey)

		credVal := existsRawCredit(ns, credKey)
		if credVal == nil {
			return nil, errors.E(errors.IO, errors.Errorf("missing credit val for key %x", credKey))
		}

		// Legacy credit output values may be of the wrong
		// size.
		scrPos := fetchRawCreditScriptOffset(credVal)
		scrLen := fetchRawCreditScriptLength(credVal)

		k := extractRawCreditTxRecordKey(credKey)
		v := existsRawTxRecord(ns, k)
		pkScript, err := fetchRawTxRecordPkScript(k, v, index,
			scrPos, scrLen)
		if err != nil {
			return nil, err
		}
		pkScripts = append(pkScripts, pkScript)
	}
	if it.err != nil {
		return nil, it.err
	}

	return pkScripts, nil
}

// Spender queries for the transaction and input index which spends a Credit.
// If the output is not a Credit, an error with code ErrInput is returned.  If
// the output is unspent, the ErrNoExist code is used.
func (s *Store) Spender(dbtx walletdb.ReadTx, out *wire.OutPoint) (*wire.MsgTx, uint32, error) {
	ns := dbtx.ReadBucket(wtxmgrBucketKey)

	var spender wire.MsgTx
	var spenderHash chainhash.Hash
	var spenderIndex uint32

	// Check mined txs
	k, v := latestTxRecord(ns, out.Hash[:])
	if v != nil {
		var block Block
		err := readRawTxRecordBlock(k, &block)
		if err != nil {
			return nil, 0, err
		}
		k = keyCredit(&out.Hash, out.Index, &block)
		v = existsRawCredit(ns, k)
		if v == nil {
			return nil, 0, errors.E(errors.Invalid, "output is not a credit")
		}
		if extractRawCreditIsSpent(v) {
			// Credit exists and is spent by a mined transaction.
			k = extractRawCreditSpenderDebitKey(v)
			copy(spenderHash[:], extractRawDebitHash(k))
			spenderIndex = extractRawDebitInputIndex(k)
			k = extractRawDebitTxRecordKey(k)
			v = existsRawTxRecord(ns, k)
			err = readRawTxRecordMsgTx(v, &spender)
			if err != nil {
				return nil, 0, err
			}
			return &spender, spenderIndex, nil
		}
		// Credit is not spent by a mined transaction, but may still be spent by
		// an unmined one.  Check whether it is spent by an unmined tx, and
		// record the spender hash if spent.
		k = canonicalOutPoint(&out.Hash, out.Index)
		v = existsRawUnminedInput(ns, k)
		if v == nil {
			return nil, 0, errors.E(errors.NotExist, "credit is unspent")
		}
		readRawUnminedInputSpenderHash(v, &spenderHash)
	}

	// If a spender exists at this point, it must be an unmined transaction.
	// The spender hash will not yet be known if the credit is also unmined, or
	// if there is no credit.
	if spenderHash == (chainhash.Hash{}) {
		k = canonicalOutPoint(&out.Hash, out.Index)
		v = existsRawUnmineVGLedit(ns, k)
		if v == nil {
			return nil, 0, errors.E(errors.Invalid, "output is not a credit")
		}
		v = existsRawUnminedInput(ns, k)
		if v == nil {
			return nil, 0, errors.E(errors.NotExist, "credit is unspent")
		}
		readRawUnminedInputSpenderHash(v, &spenderHash)
	}

	// Credit is spent by an unmined transaction.  Index is unknown so the
	// spending tx must be searched for a matching previous outpoint.
	v = existsRawUnmined(ns, spenderHash[:])
	if v == nil {
		return nil, 0, errors.E(errors.NotExist, "missing unmined spending tx")
	}
	err := spender.Deserialize(bytes.NewReader(extractRawUnminedTx(v)))
	if err != nil {
		return nil, 0, errors.E(errors.Bug, err)
	}
	found := false
	for i, in := range spender.TxIn {
		// Compare outpoints without comparing tree.
		if out.Hash == in.PreviousOutPoint.Hash && out.Index == in.PreviousOutPoint.Index {
			spenderIndex = uint32(i)
			found = true
			break
		}
	}
	if !found {
		return nil, 0, errors.E(errors.NotExist, "recorded spending tx does not spend credit")
	}
	return &spender, spenderIndex, nil
}

// RangeBlocks execute function `f` for all blocks within the given range of
// blocks in the main chain.
func (s *Store) RangeBlocks(ns walletdb.ReadBucket, begin, end int32,
	f func(*Block) (bool, error)) error {

	// Same convention as rangeTransactions: -1 means the full range.
	if begin < 0 {
		begin = int32(^uint32(0) >> 1)
	}
	if end < 0 {
		end = int32(^uint32(0) >> 1)
	}

	var blockIter blockIterator
	var advance func(*blockIterator) bool

	if begin < end {
		// Iterate in forwards order
		blockIter = makeReadBlockIterator(ns, begin)
		defer blockIter.close()
		advance = func(it *blockIterator) bool {
			if !it.next() {
				return false
			}
			return it.elem.Height <= end
		}
	} else {
		// Iterate in backwards order, from begin -> end.
		blockIter = makeReadBlockIterator(ns, begin)
		defer blockIter.close()
		advance = func(it *blockIterator) bool {
			if !it.prev() {
				return false
			}
			return end <= it.elem.Height
		}
	}

	for advance(&blockIter) {
		block := &blockIter.elem

		brk, err := f(&block.Block)
		if err != nil || brk {
			return err
		}
	}

	return nil
}
