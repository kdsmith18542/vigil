// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wallet

import (
	"context"
	"math/big"
	"time"

	"github.com/kdsmith18542/vigil/wallet/deployments"
	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/wallet/txrules"
	"github.com/kdsmith18542/vigil/wallet/wallet/udb"
	"github.com/kdsmith18542/vigil/wallet/wallet/walletdb"
	"github.com/kdsmith18542/vigil/blockchain/stake/v5"
	blockchain "github.com/kdsmith18542/vigil/blockchain/standalone/v2"
	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
	"github.com/kdsmith18542/vigil/crypto/rand"
	"github.com/kdsmith18542/vigil/txscript/v4"
	"github.com/kdsmith18542/vigil/txscript/v4/stdaddr"
	"github.com/kdsmith18542/vigil/txscript/v4/stdscript"
	"github.com/kdsmith18542/vigil/wire"
)

func (w *Wallet) extendMainChain(ctx context.Context, op errors.Op, dbtx walletdb.ReadWriteTx,
	n *BlockNode, transactions []*wire.MsgTx) ([]wire.OutPoint, error) {
	txmgrNs := dbtx.ReadWriteBucket(wtxmgrNamespaceKey)

	blockHash := n.Hash

	// Enforce checkpoints
	height := int32(n.Header.Height)
	ckpt := CheckpointHash(w.chainParams.Net, height)
	if ckpt != nil && *blockHash != *ckpt {
		err := errors.Errorf("block hash %v does not satisify "+
			"checkpoint hash %v for height %v", blockHash,
			ckpt, height)
		return nil, errors.E(errors.Consensus, err)
	}

	// Propagate the error unless this block is already included in the main
	// chain.
	err := w.txStore.ExtendMainChain(txmgrNs, n.Header, blockHash, n.FilterV2)
	if err != nil && !errors.Is(err, errors.Exist) {
		return nil, errors.E(op, err)
	}

	// Notify interested clients of the connected block.
	w.NtfnServer.notifyAttachedBlock(n.Header, blockHash)

	blockMeta, err := w.txStore.GetBlockMetaForHash(txmgrNs, blockHash)
	if err != nil {
		return nil, errors.E(op, err)
	}

	var watch []wire.OutPoint
	for _, tx := range transactions {
		// In manual ticket mode, tickets are only ever added to the
		// wallet using AddTransaction.  Skip over any relevant tickets
		// seen in this block unless they already exist in the wallet.
		if w.manualTickets && stake.IsSStx(tx) {
			txHash := tx.TxHash()
			if !w.txStore.ExistsTx(txmgrNs, &txHash) {
				continue
			}
		}

		rec, err := udb.NewTxRecordFromMsgTx(tx, time.Now())
		if err != nil {
			return nil, errors.E(op, err)
		}
		ops, err := w.processTransactionRecord(ctx, dbtx, rec, n.Header, &blockMeta)
		if err != nil {
			return nil, errors.E(op, err)
		}
		watch = append(watch, ops...)
	}

	return watch, nil
}

// ChainSwitch updates the wallet's main chain, either by extending the chain
// with new blocks, or switching to a better sidechain.  A sidechain for removed
// blocks (if any) is returned.  If relevantTxs is non-nil, the block marker for
// the latest block with processed transactions is updated for the new tip
// block.
func (w *Wallet) ChainSwitch(ctx context.Context, forest *SidechainForest, chain []*BlockNode,
	relevantTxs map[chainhash.Hash][]*wire.MsgTx) ([]*BlockNode, error) {
	const op errors.Op = "wallet.ChainSwitch"

	if len(chain) == 0 {
		return nil, errors.E(op, errors.Invalid, "zero-length chain")
	}

	chainTipChanges := &MainTipChangedNotification{
		AttachedBlocks: make([]*chainhash.Hash, 0, len(chain)),
		DetachedBlocks: nil,
		NewHeight:      int32(chain[len(chain)-1].Header.Height),
	}

	sideChainForkHeight := int32(chain[0].Header.Height)
	var prevChain []*BlockNode

	newWork := chain[len(chain)-1].workSum
	oldWork := new(big.Int)

	w.lockedOutpointMu.Lock()

	var watchOutPoints []wire.OutPoint
	err := walletdb.Update(ctx, w.db, func(dbtx walletdb.ReadWriteTx) error {
		txmgrNs := dbtx.ReadWriteBucket(wtxmgrNamespaceKey)

		tipHash, tipHeight := w.txStore.MainChainTip(dbtx)

		if tipHash == *chain[len(chain)-1].Hash {
			return nil
		}

		if sideChainForkHeight <= tipHeight {
			chainTipChanges.DetachedBlocks = make([]*chainhash.Hash, tipHeight-sideChainForkHeight+1)
			prevChain = make([]*BlockNode, tipHeight-sideChainForkHeight+1)
			for i := tipHeight; i >= sideChainForkHeight; i-- {
				hash, err := w.txStore.GetMainChainBlockHashForHeight(txmgrNs, i)
				if err != nil {
					return err
				}
				header, err := w.txStore.GetBlockHeader(dbtx, &hash)
				if err != nil {
					return err
				}
				_, filter, err := w.txStore.CFilterV2(dbtx, &hash)
				if err != nil {
					return err
				}

				// DetachedBlocks and prevChain are sorted in order of increasing heights.
				chainTipChanges.DetachedBlocks[i-sideChainForkHeight] = &hash
				prevChain[i-sideChainForkHeight] = NewBlockNode(header, &hash, filter)

				// For transaction notifications, the blocks are notified in reverse
				// height order.
				w.NtfnServer.notifyDetachedBlock(header)

				oldWork.Add(oldWork, blockchain.CalcWork(header.Bits))
			}

			if newWork.Cmp(oldWork) != 1 {
				return errors.Errorf("failed reorganize: sidechain ending at block %v has less total work "+
					"than the main chain tip block %v", chain[len(chain)-1].Hash, &tipHash)
			}

			// Remove blocks on the current main chain that are at or above the
			// height of the block that begins the side chain.
			err := w.txStore.Rollback(dbtx, sideChainForkHeight)
			if err != nil {
				return err
			}
		}

		birthState := udb.BirthState(dbtx)

		for _, n := range chain {
			if voteVersion(w.chainParams) < n.Header.StakeVersion {
				log.Warnf("Old vote version detected (v%v), please update your "+
					"wallet to the latest version.", voteVersion(w.chainParams))
			}

			watch, err := w.extendMainChain(ctx, op, dbtx, n, relevantTxs[*n.Hash])
			if err != nil {
				return err
			}
			watchOutPoints = append(watchOutPoints, watch...)

			// Add the block hash to the notification.
			chainTipChanges.AttachedBlocks = append(chainTipChanges.AttachedBlocks, n.Hash)

			if birthState != nil &&
				((birthState.SetFromTime && n.Header.Timestamp.After(birthState.Time)) ||
					(birthState.SetFromHeight && n.Header.Height == birthState.Height+1)) &&
				n.Header.Height != 0 {
				bh := n.Header.PrevBlock
				height := n.Header.Height - 1
				birthState.Hash = bh
				birthState.Height = height
				birthState.SetFromTime = false
				birthState.SetFromHeight = false
				if err := udb.SetBirthState(dbtx, birthState); err != nil {
					return err
				}
				// Do not store the tip hash as that will cause w.rescanPoint
				// to return nil. This is why we wait for one block pass the
				// height if set.
				if err := w.txStore.UpdateProcessedTxsBlockMarker(dbtx, &bh); err != nil {
					return err
				}
				log.Infof("Set wallet birthday to block %d (%v).",
					height, bh)
			}
			// NOTE: A birthday block or time set past the main tip
			// searches until it is passed.
		}

		if relevantTxs != nil {
			// To avoid skipped blocks, the marker is not advanced if there is a
			// gap between the existing rescan point (main chain fork point of
			// the current marker) and the first block attached in this chain
			// switch.
			r, err := w.rescanPoint(dbtx)
			if err != nil {
				return err
			}
			rHeader, err := w.txStore.GetBlockHeader(dbtx, r)
			if err != nil {
				return err
			}
			if !(rHeader.Height+1 < chain[0].Header.Height) {
				marker := chain[len(chain)-1].Hash
				log.Debugf("Updating processed txs block marker to %v", marker)
				err := w.txStore.UpdateProcessedTxsBlockMarker(dbtx, marker)
				if err != nil {
					return err
				}
			}
		}

		// Prune unmined transactions that don't belong on the extended chain.
		// An error here is not fatal and should just be logged.
		//
		// TODO: The stake difficulty passed here is not correct.  This must be
		// the difficulty of the next block, not the tip block.
		tip := chain[len(chain)-1]
		hashes, err := w.txStore.PruneUnmined(dbtx, tip.Header.SBits)
		if err != nil {
			log.Errorf("Failed to prune unmined transactions when "+
				"connecting block height %v: %v", tip.Header.Height, err)
		}

		for _, hash := range hashes {
			w.NtfnServer.notifyRemovedTransaction(*hash)
		}
		return nil
	})
	w.lockedOutpointMu.Unlock()
	if err != nil {
		return nil, errors.E(op, err)
	}

	if len(chainTipChanges.AttachedBlocks) != 0 {
		w.recentlyPublishedMu.Lock()
		for _, node := range chain {
			for _, tx := range relevantTxs[*node.Hash] {
				txHash := tx.TxHash()
				delete(w.recentlyPublished, txHash)
			}
		}
		w.recentlyPublishedMu.Unlock()
	}

	if n, err := w.NetworkBackend(); err == nil {
		_, err = w.watchHDAddrs(ctx, false, n)
		if err != nil {
			return nil, errors.E(op, err)
		}

		if len(watchOutPoints) > 0 {
			err = n.LoadTxFilter(ctx, false, nil, watchOutPoints)
			if err != nil {
				log.Errorf("Failed to watch outpoints: %v", err)
			}
		}
	}

	forest.PruneChain(chain)
	forest.Prune(int32(chain[len(chain)-1].Header.Height), w.chainParams)

	if w.mixingEnabled {
		w.mixClient.ExpireMessages(chain[len(chain)-1].Header.Height)
	}

	w.NtfnServer.notifyMainChainTipChanged(chainTipChanges)
	w.NtfnServer.sendAttachedBlockNotification(ctx)

	return prevChain, nil
}

// AddTransaction stores tx, marking it as mined in the block described by
// blockHash, or recording it to the wallet's mempool when nil.
//
// This method will always add ticket transactions to the wallet, even when
// configured in manual ticket mode.  It is up to network syncers to avoid
// calling this method on unmined tickets.
func (w *Wallet) AddTransaction(ctx context.Context, tx *wire.MsgTx, blockHash *chainhash.Hash) error {
	const op errors.Op = "wallet.AddTransaction"

	w.recentlyPublishedMu.Lock()
	_, recent := w.recentlyPublished[tx.TxHash()]
	w.recentlyPublishedMu.Unlock()
	if recent {
		return nil
	}

	// Prevent recording unmined tspends since they need to go through
	// voting for potentially a long time.
	if isTreasurySpend(tx) && blockHash == nil {
		log.Debugf("Ignoring unmined TSPend %s", tx.TxHash())
		return nil
	}

	w.lockedOutpointMu.Lock()
	var watchOutPoints []wire.OutPoint
	err := walletdb.Update(ctx, w.db, func(dbtx walletdb.ReadWriteTx) error {
		txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)

		rec, err := udb.NewTxRecordFromMsgTx(tx, time.Now())
		if err != nil {
			return err
		}

		// Prevent orphan votes from entering the wallet's unmined transaction
		// set.
		if isVote(&rec.MsgTx) && blockHash == nil {
			votedBlock, _ := stake.SSGenBlockVotedOn(&rec.MsgTx)
			tipBlock, _ := w.txStore.MainChainTip(dbtx)
			if votedBlock != tipBlock {
				log.Debugf("Rejected unmined orphan vote %v which votes on block %v",
					&rec.Hash, &votedBlock)
				return nil
			}
		}

		var header *wire.BlockHeader
		var meta *udb.BlockMeta
		switch {
		case blockHash != nil:
			inChain, _ := w.txStore.BlockInMainChain(dbtx, blockHash)
			if !inChain {
				break
			}
			header, err = w.txStore.GetBlockHeader(dbtx, blockHash)
			if err != nil {
				return err
			}
			meta = new(udb.BlockMeta)
			*meta, err = w.txStore.GetBlockMetaForHash(txmgrNs, blockHash)
			if err != nil {
				return err
			}
		}

		watchOutPoints, err = w.processTransactionRecord(ctx, dbtx, rec, header, meta)
		return err
	})
	w.lockedOutpointMu.Unlock()
	if err != nil {
		return errors.E(op, err)
	}
	if n, err := w.NetworkBackend(); err == nil && len(watchOutPoints) > 0 {
		_, err := w.watchHDAddrs(ctx, false, n)
		if err != nil {
			return errors.E(op, err)
		}
		if len(watchOutPoints) > 0 {
			err = n.LoadTxFilter(ctx, false, nil, watchOutPoints)
			if err != nil {
				log.Errorf("Failed to watch outpoints: %v", err)
			}
		}
	}
	return nil
}

func (w *Wallet) processTransactionRecord(ctx context.Context, dbtx walletdb.ReadWriteTx, rec *udb.TxRecord,
	header *wire.BlockHeader, blockMeta *udb.BlockMeta) (watchOutPoints []wire.OutPoint, err error) {

	const op errors.Op = "wallet.processTransactionRecord"

	addrmgrNs := dbtx.ReadWriteBucket(waddrmgrNamespaceKey)
	txmgrNs := dbtx.ReadWriteBucket(wtxmgrNamespaceKey)

	// At the moment all notified transactions are assumed to actually be
	// relevant.  This assumption will not hold true when SPV support is
	// added, but until then, simply insert the transaction because there
	// should either be one or more relevant inputs or outputs.
	if header == nil {
		err = w.txStore.InsertMemPoolTx(dbtx, rec)
		if errors.Is(err, errors.Exist) {
			log.Warnf("Refusing to add unmined transaction %v since same "+
				"transaction already exists mined", &rec.Hash)
			return nil, nil
		}
	} else {
		err = w.txStore.InsertMinedTx(dbtx, rec, &blockMeta.Hash)
	}
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Skip unlocking outpoints if the transaction is a vote or revocation as the lock
	// is not held.
	skipOutpoints := rec.TxType == stake.TxTypeSSGen || rec.TxType == stake.TxTypeSSRtx

	// Handle input scripts that contain P2PKs that we care about.
	for i, input := range rec.MsgTx.TxIn {
		if !skipOutpoints {
			prev := input.PreviousOutPoint
			delete(w.lockedOutpoints, outpoint{prev.Hash, prev.Index})
		}
		// TODO: the prevout's actual pkScript version is needed.
		if stdscript.IsMultiSigSigScript(scriptVersionAssumed, input.SignatureScript) {
			rs := stdscript.MultiSigRedeemScriptFromScriptSigV0(input.SignatureScript)

			class, addrs := stdscript.ExtractAddrs(scriptVersionAssumed, rs, w.chainParams)
			if class != stdscript.STMultiSig {
				// This should never happen, but be paranoid.
				continue
			}

			isRelevant := false
			for _, addr := range addrs {
				ma, err := w.manager.Address(addrmgrNs, addr)
				if err != nil {
					// Missing addresses are skipped.  Other errors should be
					// propagated.
					if errors.Is(err, errors.NotExist) {
						continue
					}
					return nil, errors.E(op, err)
				}
				isRelevant = true
				err = w.markUsedAddress(op, dbtx, ma)
				if err != nil {
					return nil, err
				}
				log.Debugf("Marked address %v used", addr)
			}

			// Add the script to the script databases.
			// TODO Markused script address? cj
			if isRelevant {
				n, _ := w.NetworkBackend()
				addr, err := w.manager.ImportScript(addrmgrNs, rs)
				switch {
				case errors.Is(err, errors.Exist):
				case err != nil:
					return nil, errors.E(op, err)
				case n != nil:
					addrs := []stdaddr.Address{addr.Address()}
					err := n.LoadTxFilter(ctx, false, addrs, nil)
					if err != nil {
						return nil, errors.E(op, err)
					}
				}
			}

			// If we're spending a multisig outpoint we know about,
			// update the outpoint. Inefficient because you deserialize
			// the entire multisig output info. Consider a specific
			// exists function in udb. The error here is skipped
			// because the absence of an multisignature output for
			// some script can not always be considered an error. For
			// example, the wallet might be rescanning as called from
			// the above function and so does not have the output
			// included yet.
			mso, err := w.txStore.GetMultisigOutput(txmgrNs, &input.PreviousOutPoint)
			if mso != nil && err == nil {
				err = w.txStore.SpendMultisigOut(txmgrNs, &input.PreviousOutPoint,
					rec.Hash, uint32(i))
				if err != nil {
					return nil, errors.E(op, err)
				}
			}
		}
	}

	// Check every output to determine whether it is controlled by a
	// wallet key.  If so, mark the output as a credit and mark
	// outpoints to watch.
	for i, output := range rec.MsgTx.TxOut {
		class, addrs := stdscript.ExtractAddrs(output.Version, output.PkScript, w.chainParams)
		if class == stdscript.STNonStandard {
			// Non-standard outputs are skipped.
			continue
		}
		subClass, isStakeType := txrules.StakeSubScriptType(class)
		if isStakeType {
			class = subClass
		}

		isTicketCommit := rec.TxType == stake.TxTypeSStx && i%2 == 1
		watchOutPoint := true
		if isTicketCommit {
			// For ticket commitments, decode the address stored in the pkscript
			// and evaluate ownership of that.
			addr, err := stake.AddrFromSStxPkScrCommitment(output.PkScript,
				w.chainParams)
			if err != nil {
				log.Warnf("failed to decode ticket commitment script of %s:%d",
					rec.Hash, i)
				continue
			}
			addrs = []stdaddr.Address{addr}
			watchOutPoint = false
		} else if output.Value == 0 {
			// The only case of outputs with 0 value that we need to handle are
			// ticket commitments. All other outputs can be ignored.
			continue
		}

		var tree int8
		if isStakeType {
			tree = 1
		}
		outpoint := wire.OutPoint{Hash: rec.Hash, Tree: tree}
		for _, addr := range addrs {
			ma, err := w.manager.Address(addrmgrNs, addr)
			// Missing addresses are skipped.  Other errors should
			// be propagated.
			if errors.Is(err, errors.NotExist) {
				continue
			}
			if err != nil {
				return nil, errors.E(op, err)
			}
			if isTicketCommit {
				err = w.txStore.AddTicketCommitment(txmgrNs, rec, uint32(i),
					ma.Account())
			} else {
				err = w.txStore.AdVGLedit(dbtx, rec, blockMeta,
					uint32(i), ma.Internal(), ma.Account())
			}
			if err != nil {
				return nil, errors.E(op, err)
			}
			err = w.markUsedAddress(op, dbtx, ma)
			if err != nil {
				return nil, err
			}
			if watchOutPoint {
				outpoint.Index = uint32(i)
				watchOutPoints = append(watchOutPoints, outpoint)
			}
			log.Debugf("Marked address %v used", addr)
		}

		// Handle P2SH addresses that are multisignature scripts
		// with keys that we own.
		if class == stdscript.STScriptHash {
			var expandedScript []byte
			for _, addr := range addrs {
				expandedScript, err = w.manager.RedeemScript(addrmgrNs, addr)
				if err != nil {
					log.Debugf("failed to find redeemscript for "+
						"address %v in address manager: %v",
						addr, err)
					continue
				}
			}

			// Otherwise, extract the actual addresses and see if any are ours.
			expClass, multisigAddrs := stdscript.ExtractAddrs(scriptVersionAssumed, expandedScript, w.chainParams)

			// Skip non-multisig scripts.
			if expClass != stdscript.STMultiSig {
				continue
			}

			for _, maddr := range multisigAddrs {
				_, err := w.manager.Address(addrmgrNs, maddr)
				// An address we own; handle accordingly.
				if err == nil {
					err := w.txStore.AddMultisigOut(
						dbtx, rec, blockMeta, uint32(i))
					if err != nil {
						// This will throw if there are multiple private keys
						// for this multisignature output owned by the wallet,
						// so it's routed to debug.
						log.Debugf("unable to add multisignature output: %v", err)
					}
				}
			}
		}
	}

	if (rec.TxType == stake.TxTypeSSGen) || (rec.TxType == stake.TxTypeSSRtx) {
		err = w.txStore.RedeemTicketCommitments(txmgrNs, rec, blockMeta)
		if err != nil {
			log.Errorf("Error redeeming ticket commitments: %v", err)
		}
	}

	// Send notification of mined or unmined transaction to any interested
	// clients.
	//
	// TODO: Avoid the extra db hits.
	if header == nil {
		details, err := w.txStore.UniqueTxDetails(txmgrNs, &rec.Hash, nil)
		if err != nil {
			log.Errorf("Cannot query transaction details for notifiation: %v", err)
		} else {
			w.NtfnServer.notifyUnminedTransaction(dbtx, details)
		}
	} else {
		details, err := w.txStore.UniqueTxDetails(txmgrNs, &rec.Hash, &blockMeta.Block)
		if err != nil {
			log.Errorf("Cannot query transaction details for notifiation: %v", err)
		} else {
			w.NtfnServer.notifyMinedTransaction(dbtx, details, blockMeta)
		}
	}

	return watchOutPoints, nil
}

// selectOwnedTickets returns a slice of tickets hashes from the tickets
// argument that are owned by the wallet.
//
// Because votes must be created for tickets tracked by both the transaction
// manager and the stake manager, this function checks both.
func selectOwnedTickets(w *Wallet, dbtx walletdb.ReadTx, tickets []*chainhash.Hash) []*chainhash.Hash {
	var owned []*chainhash.Hash
	for _, ticketHash := range tickets {
		if w.txStore.OwnTicket(dbtx, ticketHash) {
			owned = append(owned, ticketHash)
		}
	}
	return owned
}

// VoteOnOwnedTickets creates and publishes vote transactions for all owned
// tickets in the winningTicketHashes slice if wallet voting is enabled.  The
// vote is only valid when voting on the block described by the passed block
// hash and height.  When a network backend is associated with the wallet,
// relevant commitment outputs are loaded as watched data.
func (w *Wallet) VoteOnOwnedTickets(ctx context.Context, winningTicketHashes []*chainhash.Hash, blockHash *chainhash.Hash, blockHeight int32) error {
	const op errors.Op = "wallet.VoteOnOwnedTickets"

	if !w.votingEnabled || blockHeight < int32(w.chainParams.StakeValidationHeight)-1 {
		return nil
	}

	n, err := w.NetworkBackend()
	if err != nil {
		return errors.E(op, err)
	}
	dq, ok := n.(deployments.Querier)
	if !ok {
		return errors.E(op, "network backend does not provide deployment information")
	}
	VGLP0010Active, err := deployments.VGLP0010Active(ctx, blockHeight,
		w.chainParams, dq)
	if err != nil {
		return errors.E(op, err)
	}
	VGLP0012Active, err := deployments.VGLP0012Active(ctx, blockHeight,
		w.chainParams, dq)
	if err != nil {
		return errors.E(op, err)
	}

	// TODO The behavior of this is not quite right if tons of blocks
	// are coming in quickly, because the transaction store will end up
	// out of sync with the voting channel here. This should probably
	// be fixed somehow, but this should be stable for networks that
	// are voting at normal block speeds.

	var ticketHashes []*chainhash.Hash
	var votes []*wire.MsgTx
	var usedVoteBits []stake.VoteBits
	defaultVoteBits := w.VoteBits()
	var watchOutPoints []wire.OutPoint
	err = walletdb.View(ctx, w.db, func(dbtx walletdb.ReadTx) error {
		txmgrNs := dbtx.ReadBucket(wtxmgrNamespaceKey)

		// Only consider tickets owned by this wallet.
		ticketHashes = selectOwnedTickets(w, dbtx, winningTicketHashes)
		if len(ticketHashes) == 0 {
			return nil
		}

		votes = make([]*wire.MsgTx, len(ticketHashes))
		usedVoteBits = make([]stake.VoteBits, len(ticketHashes))

		addrmgrNs := dbtx.ReadBucket(waddrmgrNamespaceKey)

		for i, ticketHash := range ticketHashes {
			ticketPurchase, err := w.txStore.Tx(txmgrNs, ticketHash)
			if err != nil {
				log.Errorf("Failed to read ticket purchase transaction for "+
					"owned winning ticket %v: %v", ticketHash, err)
				continue
			}

			// Don't create votes when this wallet doesn't have voting
			// authority or the private key to vote.
			owned, haveKey, err := w.hasVotingAuthority(addrmgrNs, ticketPurchase)
			if err != nil {
				return err
			}
			if !(owned && haveKey) {
				continue
			}

			ticketVoteBits := defaultVoteBits
			// Check for and use per-ticket votebits if set for this ticket.
			if tvb, found := w.readDBTicketVoteBits(dbtx, ticketHash); found {
				ticketVoteBits = tvb
			}

			// When not on mainnet, randomly disapprove blocks based
			// on the disapprove percent.
			dp := w.DisapprovePercent()
			if dp > 0 {
				if w.chainParams.Net == wire.MainNet {
					log.Warnf("block disapprove percent set on mainnet")
				} else if int64(dp) > rand.Int64N(100) {
					log.Infof("Disapproving block %v voted with ticket %v",
						blockHash, ticketHash)
					// Set the BlockValid bit to zero,
					// disapproving the block.
					const blockIsValidBit = uint16(0x01)
					ticketVoteBits.Bits &= ^blockIsValidBit
				}
			}

			// Deal with treasury votes
			tspends := w.GetAllTSpends(ctx)

			// Dealwith consensus votes
			vote, err := createUnsignedVote(ticketHash, ticketPurchase,
				blockHeight, blockHash, ticketVoteBits, w.subsidyCache,
				w.chainParams, VGLP0010Active, VGLP0012Active)
			if err != nil {
				log.Errorf("Failed to create vote transaction for ticket "+
					"hash %v: %v", ticketHash, err)
				continue
			}

			// Iterate over all tpends and determine if they are
			// within the voting window.
			tVotes := make([]byte, 0, 256)
			tVotes = append(tVotes, 'T', 'V')
			for _, v := range tspends {
				if !blockchain.InsideTSpendWindow(int64(blockHeight),
					v.Expiry, w.chainParams.TreasuryVoteInterval,
					w.chainParams.TreasuryVoteIntervalMultiplier) {
					continue
				}

				// Get policy for tspend, falling back to any
				// policy for the Pi key.
				tspendHash := v.TxHash()
				tspendVote := w.TSpendPolicy(&tspendHash, ticketHash)
				if tspendVote == stake.TreasuryVoteInvalid {
					continue
				}

				// Append tspend hash and vote bits
				tVotes = append(tVotes, tspendHash[:]...)
				tVotes = append(tVotes, byte(tspendVote))
			}
			if len(tVotes) > 2 {
				// Vote was appended. Create output and flip
				// script version.
				var b txscript.ScriptBuilder
				b.AddOp(txscript.OP_RETURN)
				b.AddData(tVotes)
				tspendVoteScript, err := b.Script()
				if err != nil {
					// Log error and continue.
					log.Errorf("Failed to create treasury "+
						"vote for ticket hash %v: %v",
						ticketHash, err)
				} else {
					// Success.
					vote.AddTxOut(wire.NewTxOut(0, tspendVoteScript))
					vote.Version = 3
				}
			}

			// Sign vote and sumit.
			err = w.signVote(addrmgrNs, ticketPurchase, vote)
			if err != nil {
				log.Errorf("Failed to sign vote for ticket hash %v: %v",
					ticketHash, err)
				continue
			}
			votes[i] = vote
			usedVoteBits[i] = ticketVoteBits

			watchOutPoints = w.appendRelevantOutpoints(watchOutPoints, dbtx, vote)
		}
		return nil
	})
	if err != nil {
		log.Errorf("View failed: %v", errors.E(op, err))
	}

	// Remove nil votes without preserving order.
	for i := 0; i < len(votes); {
		if votes[i] == nil {
			votes[i], votes[len(votes)-1] = votes[len(votes)-1], votes[i]
			votes = votes[:len(votes)-1]
			continue
		}
		i++
	}

	voteRecords := make([]*udb.TxRecord, 0, len(votes))
	for i := range votes {
		rec, err := udb.NewTxRecordFromMsgTx(votes[i], time.Now())
		if err != nil {
			log.Errorf("Failed to create transaction record: %v", err)
			continue
		}
		voteRecords = append(voteRecords, rec)
	}
	w.recentlyPublishedMu.Lock()
	for i := range voteRecords {
		w.recentlyPublished[voteRecords[i].Hash] = struct{}{}

		log.Infof("Voting on block %v (height %v) using ticket %v "+
			"(vote hash: %v bits: %v)", blockHash, blockHeight,
			ticketHashes[i], &voteRecords[i].Hash, usedVoteBits[i].Bits)
	}
	w.recentlyPublishedMu.Unlock()

	// Publish before recording votes in database to slightly reduce latency.
	err = n.PublishTransactions(ctx, votes...)
	if err != nil {
		log.Errorf("Failed to send one or more votes: %v", err)
	}

	if len(watchOutPoints) > 0 {
		err := n.LoadTxFilter(ctx, false, nil, watchOutPoints)
		if err != nil {
			log.Errorf("Failed to watch outpoints: %v", err)
		}
	}

	// w.lockedOutpointMu is intentionally not locked.
	err = walletdb.Update(ctx, w.db, func(dbtx walletdb.ReadWriteTx) error {
		for i := range voteRecords {
			_, err := w.processTransactionRecord(ctx, dbtx, voteRecords[i], nil, nil)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if n, err := w.NetworkBackend(); err == nil {
		_, err := w.watchHDAddrs(ctx, false, n)
		if err != nil {
			return err
		}
	}
	return nil
}

// RevokeOwnedTickets no longer revokes any tickets since revocations are now
// automatically created per VGLP0009.
//
// Deprecated: this method will be removed in the next major version.
func (w *Wallet) RevokeOwnedTickets(ctx context.Context, missedTicketHashes []*chainhash.Hash) error {
	return nil
}
