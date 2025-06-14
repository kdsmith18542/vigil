// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
)

var (
	errVoterVersionMajorityNotFound = errors.New("voter version majority " +
		"not found")
	errStakeVersionMajorityNotFound = errors.New("stake version majority " +
		"not found")
)

// stakeMajorityCacheVersionKey creates a map key that is comprised of a stake
// version and a hash.  This is used for caches that require a version in
// addition to a simple hash.
func stakeMajorityCacheVersionKey(version uint32, hash *chainhash.Hash) [stakeMajorityCacheKeySize]byte {
	var key [stakeMajorityCacheKeySize]byte
	binary.LittleEndian.PutUint32(key[0:], version)
	copy(key[4:], hash[:])
	return key
}

// calcWantHeight calculates the height of the final block of the previous
// interval given a stake validation height, stake validation interval, and
// block height.
func calcWantHeight(stakeValidationHeight, interval, height int64) int64 {
	// The adjusted height accounts for the fact the starting validation
	// height does not necessarily start on an interval and thus the
	// intervals might not be zero-based.
	intervalOffset := stakeValidationHeight % interval
	adjustedHeight := height - intervalOffset - 1
	return (adjustedHeight - ((adjustedHeight + 1) % interval)) +
		intervalOffset
}

// CalcWantHeight calculates the height of the final block of the previous
// interval given a block height.
func (b *BlockChain) CalcWantHeight(interval, height int64) int64 {
	return calcWantHeight(b.chainParams.StakeValidationHeight, interval,
		height)
}

// findStakeVersionPriorNode walks the chain backwards from prevNode until it
// reaches the final block of the previous stake version interval and returns
// that node.  The returned node will be nil when the provided prevNode is too
// low such that there is no previous stake version interval due to falling
// prior to the stake validation interval.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) findStakeVersionPriorNode(prevNode *blockNode) *blockNode {
	// Check to see if the blockchain is high enough to begin accounting
	// stake versions.
	svh := b.chainParams.StakeValidationHeight
	svi := b.chainParams.StakeVersionInterval
	nextHeight := prevNode.height + 1
	if nextHeight < svh+svi {
		return nil
	}

	wantHeight := calcWantHeight(svh, svi, nextHeight)
	return prevNode.Ancestor(wantHeight)
}

// isStakeMajorityVersion determines if minVer requirement is met based on
// prevNode.  The function always uses the stake versions of the prior window.
// For example, if StakeVersionInterval = 11 and StakeValidationHeight = 13 the
// windows start at 13 + (11 * 2) 25 and are as follows: 24-34, 35-45, 46-56 ...
// If height comes in at 35 we use the 24-34 window, up to height 45.
// If height comes in at 46 we use the 35-45 window, up to height 56 etc.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) isStakeMajorityVersion(minVer uint32, prevNode *blockNode) bool {
	// Walk blockchain backwards to calculate version.
	node := b.findStakeVersionPriorNode(prevNode)
	if node == nil {
		return 0 == minVer
	}

	// Generate map key and look up cached result.
	key := stakeMajorityCacheVersionKey(minVer, &node.hash)
	if result, ok := b.isStakeMajorityVersionCache[key]; ok {
		return result
	}

	// Tally how many of the block headers in the previous stake version
	// validation interval have their stake version set to at least the
	// requested minimum version.
	versionCount := int32(0)
	iterNode := node
	for i := int64(0); i < b.chainParams.StakeVersionInterval && iterNode != nil; i++ {
		if iterNode.stakeVersion >= minVer {
			versionCount++
		}

		iterNode = iterNode.parent
	}

	// Determine the required amount of votes to reach supermajority.
	numRequired := int32(b.chainParams.StakeVersionInterval) *
		b.chainParams.StakeMajorityMultiplier /
		b.chainParams.StakeMajorityDivisor

	// Cache result.
	result := versionCount >= numRequired
	b.isStakeMajorityVersionCache[key] = result

	return result
}

// calcPriorStakeVersion calculates the header stake version of the prior
// interval.  The function walks the chain backwards by one interval and then
// it performs a standard majority calculation.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) calcPriorStakeVersion(prevNode *blockNode) (uint32, error) {
	// Walk blockchain backwards to calculate version.
	node := b.findStakeVersionPriorNode(prevNode)
	if node == nil {
		return 0, nil
	}

	// Check cache.
	if result, ok := b.calcPriorStakeVersionCache[node.hash]; ok {
		return result, nil
	}

	// Tally how many of each stake version the block headers in the previous stake
	// version validation interval have.
	versions := make(map[uint32]int32) // [version][count]
	iterNode := node
	for i := int64(0); i < b.chainParams.StakeVersionInterval && iterNode != nil; i++ {
		versions[iterNode.stakeVersion]++

		iterNode = iterNode.parent
	}

	// Determine the required amount of votes to reach supermajority.
	numRequired := int32(b.chainParams.StakeVersionInterval) *
		b.chainParams.StakeMajorityMultiplier /
		b.chainParams.StakeMajorityDivisor

	for version, count := range versions {
		if count >= numRequired {
			b.calcPriorStakeVersionCache[node.hash] = version
			return version, nil
		}
	}

	return 0, errStakeVersionMajorityNotFound
}

// calcVoterVersionInterval tallies all voter versions in an interval and
// returns a version that has reached 75% majority.  This function MUST be
// called with a node that is the final node in a valid stake version interval
// and greater than or equal to the stake validation height or it will result in
// an assertion error.
//
// This function is really meant to be called internally only from this file.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) calcVoterVersionInterval(prevNode *blockNode) (uint32, error) {
	// Ensure the provided node is the final node in a valid stake version
	// interval and is greater than or equal to the stake validation height
	// since the logic below relies on these assumptions.
	svh := b.chainParams.StakeValidationHeight
	svi := b.chainParams.StakeVersionInterval
	expectedHeight := calcWantHeight(svh, svi, prevNode.height+1)
	if prevNode.height != expectedHeight || expectedHeight < svh {
		return 0, AssertError(fmt.Sprintf("calcVoterVersionInterval "+
			"must be called with a node that is the final node "+
			"in a stake version interval -- called with node %s "+
			"(height %d)", prevNode.hash, prevNode.height))
	}

	// See if we have cached results.
	if result, ok := b.calcVoterVersionIntervalCache[prevNode.hash]; ok {
		return result, nil
	}

	// Tally both the total number of votes in the previous stake version validation
	// interval and how many of each version those votes have.
	versions := make(map[uint32]int32) // [version][count]
	totalVotesFound := int32(0)
	iterNode := prevNode
	for i := int64(0); i < svi && iterNode != nil; i++ {
		totalVotesFound += int32(len(iterNode.votes))
		for _, v := range iterNode.votes {
			versions[v.Version]++
		}

		iterNode = iterNode.parent
	}

	// Determine the required amount of votes to reach supermajority.
	numRequired := totalVotesFound * b.chainParams.StakeMajorityMultiplier /
		b.chainParams.StakeMajorityDivisor

	for version, count := range versions {
		if count >= numRequired {
			b.calcVoterVersionIntervalCache[prevNode.hash] = version
			return version, nil
		}
	}

	return 0, errVoterVersionMajorityNotFound
}

// calcVoterVersion calculates the last prior valid majority stake version.  If
// the current interval does not have a majority stake version it'll go back to
// the prior interval.  It'll keep going back up to the minimum height at which
// point we know the version was 0.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) calcVoterVersion(prevNode *blockNode) (uint32, *blockNode) {
	// Walk blockchain backwards to find interval.
	node := b.findStakeVersionPriorNode(prevNode)

	// Iterate over versions until a majority is found.  Don't try to count
	// votes before the stake validation height since there could not
	// possibly have been any.
	for node != nil && node.height >= b.chainParams.StakeValidationHeight {
		version, err := b.calcVoterVersionInterval(node)
		if err == nil {
			return version, node
		}
		if !errors.Is(err, errVoterVersionMajorityNotFound) {
			break
		}

		node = node.RelativeAncestor(b.chainParams.StakeVersionInterval)
	}

	// No majority version found.
	return 0, nil
}

// calcStakeVersion calculates the header stake version based on voter versions.
// If there is a majority of voter versions it uses the header stake version to
// prevent reverting to a prior version.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) calcStakeVersion(prevNode *blockNode) uint32 {
	version, node := b.calcVoterVersion(prevNode)
	if version == 0 || node == nil {
		// Short circuit.
		return 0
	}

	// Check cache.
	if result, ok := b.calcStakeVersionCache[node.hash]; ok {
		return result
	}

	// Walk chain backwards to start of node interval (start of current
	// period) Note that calcWantHeight returns the LAST height of the
	// prior interval; hence the + 1.
	startIntervalHeight := calcWantHeight(b.chainParams.StakeValidationHeight,
		b.chainParams.StakeVersionInterval, node.height) + 1
	startNode := node.Ancestor(startIntervalHeight)
	if startNode == nil {
		// Note that should this not be possible to hit because a
		// majority voter version was obtained above, which means there
		// is at least an interval of nodes.  However, be paranoid.
		b.calcStakeVersionCache[node.hash] = 0
	}

	// See if we are enforcing V3 blocks yet.  Just return V0 since it
	// wasn't enforced and therefore irrelevant.
	if !b.isMajorityVersion(3, startNode,
		b.chainParams.BlockRejectNumRequired) {
		b.calcStakeVersionCache[node.hash] = 0
		return 0
	}

	// Don't allow the stake version to go backwards once it has been locked
	// in by a previous majority, even if the majority of votes are now a
	// lower version.
	if b.isStakeMajorityVersion(version, node) {
		priorVersion, _ := b.calcPriorStakeVersion(node)
		if priorVersion > version {
			version = priorVersion
		}
	}

	b.calcStakeVersionCache[node.hash] = version
	return version
}

// calcStakeVersionByHash calculates the last prior valid majority stake
// version.  If the current interval does not have a majority stake version
// it'll go back to the prior interval.  It'll keep going back up to the
// minimum height at which point we know the version was 0.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) calcStakeVersionByHash(hash *chainhash.Hash) (uint32, error) {
	prevNode := b.index.LookupNode(hash)
	if prevNode == nil || !b.index.CanValidate(prevNode) {
		return 0, fmt.Errorf("block %s is not known", hash)
	}

	return b.calcStakeVersion(prevNode), nil
}

// CalcStakeVersionByHash calculates the expected stake version for the block
// AFTER the provided block hash.
//
// This function is safe for concurrent access.
func (b *BlockChain) CalcStakeVersionByHash(hash *chainhash.Hash) (uint32, error) {
	b.chainLock.Lock()
	version, err := b.calcStakeVersionByHash(hash)
	b.chainLock.Unlock()
	return version, err
}
