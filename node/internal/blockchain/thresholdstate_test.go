// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/kdsmith18542/vigil/blockchain/v5/chaingen"
	"github.com/kdsmith18542/vigil/chaincfg/v3"
)

const (
	// vbPrevBlockValid defines the vote bit necessary to vote yes to the
	// previous block being valid.
	vbPrevBlockValid = 0x01
)

// mockVote1 returns a voting agenda for use throughout the tests.
func mockVote1() chaincfg.Vote {
	return chaincfg.Vote{
		Id:          "mockvote1",
		Description: "Mock vote 1",
		Mask:        0x6, // 0b0110
		Choices: []chaincfg.Choice{{
			Id:          "abstain",
			Description: "abstain voting for change",
			Bits:        0x0000,
			IsAbstain:   true,
		}, {
			Id:          "no",
			Description: "vote no",
			Bits:        0x0002, // Bit 1 (1 << 1)
			IsNo:        true,
		}, {
			Id:          "yes",
			Description: "vote yes",
			Bits:        0x0004, // Bit 2 (2 << 1)
		}},
	}
}

// mockVote2 returns a voting agenda with a different mask and choice bits than
// those in mockVote1 for use throughout the tests.
func mockVote2() chaincfg.Vote {
	return chaincfg.Vote{
		Id:          "mockvote2",
		Description: "Mock vote 2",
		Mask:        0x18, // 0b11000
		Choices: []chaincfg.Choice{{
			Id:          "abstain",
			Description: "abstain voting for change",
			Bits:        0x0000,
			IsAbstain:   true,
		}, {
			Id:          "no",
			Description: "vote no",
			Bits:        0x0008, // Bit 3 (1 << 3)
			IsNo:        true,
		}, {
			Id:          "yes",
			Description: "vote yes",
			Bits:        0x0010, // Bit 4 (2 << 3)
		}},
	}
}

// TestCurrentDeploymentVersion ensures that the highest deployment version is
// returned based on given network parameters.
func TestCurrentDeploymentVersion(t *testing.T) {
	t.Parallel()

	// Default parameters to use for tests.  Clone the parameters so they can be
	// mutated.
	params := cloneParams(chaincfg.RegNetParams())

	tests := []struct {
		name        string
		deployments map[uint32][]chaincfg.ConsensusDeployment
		wantVersion uint32
	}{{
		name:        "no deployments defined",
		wantVersion: 0,
	}, {
		name: "single deployment defined",
		deployments: map[uint32][]chaincfg.ConsensusDeployment{
			7: {{
				Vote: mockVote1(),
			}},
		},
		wantVersion: 7,
	}, {
		name: "multiple deployments defined",
		deployments: map[uint32][]chaincfg.ConsensusDeployment{
			7: {{
				Vote: mockVote1(),
			}},
			8: {{
				Vote: mockVote2(),
			}},
		},
		wantVersion: 8,
	}}
	for _, test := range tests {
		// Set deployments based on the test parameter.
		params.Deployments = test.deployments

		// Ensure that the returned version matches the expected version.
		gotVersion := currentDeploymentVersion(params)
		if gotVersion != test.wantVersion {
			t.Errorf("%q: mismatched current deployment version:\nwant: %v\n "+
				"got: %v\n", test.name, test.wantVersion, gotVersion)
		}
	}
}

// TestNextDeploymentVersion ensures that the next deployment version is
// returned based on given network parameters.
func TestNextDeploymentVersion(t *testing.T) {
	t.Parallel()

	// Default parameters to use for tests.  Clone the parameters so they can be
	// mutated.
	params := cloneParams(chaincfg.RegNetParams())
	deployments := map[uint32][]chaincfg.ConsensusDeployment{
		7: {{
			Vote: mockVote1(),
		}},
		8: {{
			Vote: mockVote2(),
		}},
	}

	tests := []struct {
		name        string
		deployments map[uint32][]chaincfg.ConsensusDeployment
		version     uint32
		wantVersion uint32
	}{{
		name:        "no deployments defined",
		version:     0,
		wantVersion: 0,
	}, {
		name:        "provided version is less than all defined versions",
		deployments: deployments,
		version:     6,
		wantVersion: 7,
	}, {
		name:        "provided version is between defined versions",
		deployments: deployments,
		version:     7,
		wantVersion: 8,
	}, {
		name:        "provided version is the current version",
		deployments: deployments,
		version:     8,
		wantVersion: 0,
	}, {
		name:        "provided version is greater than the current version",
		deployments: deployments,
		version:     9,
		wantVersion: 0,
	}}

	for _, test := range tests {
		// Set deployments based on the test parameter.
		params.Deployments = test.deployments

		// Ensure that the returned version matches the expected version.
		gotVersion := nextDeploymentVersion(params, test.version)
		if gotVersion != test.wantVersion {
			t.Errorf("%q: mismatched next deployment version:\nwant: %v\n "+
				"got: %v\n", test.name, test.wantVersion, gotVersion)
		}
	}
}

// TestThresholdState ensures that the threshold state function progresses
// through the states correctly.
func TestThresholdState(t *testing.T) {
	t.Parallel()

	// Create chain params based on regnet params, but add a specific mock votes
	// and set the proof-of-work difficulty readjustment size to a really large
	// number so that the test chain can be generated more quickly.
	const posVersion = 4
	params := chaincfg.RegNetParams()
	params.WorkDiffWindowSize = 200000
	params.WorkDiffWindows = 1
	params.TargetTimespan = params.TargetTimePerBlock *
		time.Duration(params.WorkDiffWindowSize)
	if params.Deployments == nil {
		params.Deployments = make(map[uint32][]chaincfg.ConsensusDeployment)
	}
	params.Deployments[posVersion] = append(params.Deployments[posVersion],
		chaincfg.ConsensusDeployment{
			Vote:       mockVote1(),
			StartTime:  0,
			ExpireTime: math.MaxUint64,
		})
	params.Deployments[posVersion] = append(params.Deployments[posVersion],
		chaincfg.ConsensusDeployment{
			Vote:       mockVote2(),
			StartTime:  0,
			ExpireTime: math.MaxUint64,
		})
	reassignVoteMaskAndChoiceBits(params.Deployments[posVersion])
	numDeployments := len(params.Deployments[posVersion])

	// Convenient references to the mock parameter votes and choices.
	vote1 := &params.Deployments[posVersion][numDeployments-2].Vote
	vote1Yes := findVoteChoice(t, vote1, "yes")
	vote1No := findVoteChoice(t, vote1, "no")
	vote2 := &params.Deployments[posVersion][numDeployments-1].Vote
	vote2Yes := findVoteChoice(t, vote2, "yes")
	vote2No := findVoteChoice(t, vote2, "no")

	// Create a test harness initialized with the genesis block as the tip.
	g := newChaingenHarness(t, params)

	// Shorter versions of useful params for convenience.
	ticketsPerBlock := int64(params.TicketsPerBlock)
	coinbaseMaturity := params.CoinbaseMaturity
	stakeEnabledHeight := params.StakeEnabledHeight
	stakeValidationHeight := params.StakeValidationHeight
	stakeVerInterval := params.StakeVersionInterval
	ruleChangeInterval := int64(params.RuleChangeActivationInterval)
	powNumToCheck := int64(params.BlockUpgradeNumToCheck)
	ruleChangeQuorum := int64(params.RuleChangeActivationQuorum)
	ruleChangeMult := int64(params.RuleChangeActivationMultiplier)
	ruleChangeDiv := int64(params.RuleChangeActivationDivisor)

	// ---------------------------------------------------------------------
	// Block One.
	//
	// NOTE: The advance funcs on the harness are intentionally not used in
	// these tests since they need to manually test the threshold state at
	// all key heights.
	// ---------------------------------------------------------------------

	// Add the required initial block.
	//
	//   genesis -> bfb
	g.CreateBlockOne("bfb", 0)
	g.AssertTipHeight(1)
	g.AcceptTipBlock()
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to have mature coinbase outputs to work with.
	//
	//   genesis -> bfb -> bm0 -> bm1 -> ... -> bm#
	// ---------------------------------------------------------------------

	for i := uint16(0); i < coinbaseMaturity; i++ {
		blockName := fmt.Sprintf("bm%d", i)
		g.NextBlock(blockName, nil, nil)
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(coinbaseMaturity) + 1)
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the stake enabled height while
	// creating ticket purchases that spend from the coinbases matured
	// above.  This will also populate the pool of immature tickets.
	//
	//   ... -> bm# ... -> bse0 -> bse1 -> ... -> bse#
	// ---------------------------------------------------------------------

	var ticketsPurchased int
	for i := int64(0); int64(g.Tip().Header.Height) < stakeEnabledHeight; i++ {
		outs := g.OldestCoinbaseOuts()
		ticketOuts := outs[1:]
		ticketsPurchased += len(ticketOuts)
		blockName := fmt.Sprintf("bse%d", i)
		g.NextBlock(blockName, nil, ticketOuts)
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeEnabledHeight))
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the stake validation height while
	// continuing to purchase tickets using the coinbases matured above and
	// allowing the immature tickets to mature and thus become live.
	//
	// The blocks are also generated with version 3 to ensure stake version
	// enforcement is reached.
	// ---------------------------------------------------------------------

	targetPoolSize := int64(g.Params().TicketPoolSize) * ticketsPerBlock
	for i := int64(0); int64(g.Tip().Header.Height) < stakeValidationHeight; i++ {
		// Only purchase tickets until the target ticket pool size is
		// reached.
		outs := g.OldestCoinbaseOuts()
		ticketOuts := outs[1:]
		if ticketsPurchased+len(ticketOuts) > int(targetPoolSize) {
			ticketsNeeded := int(targetPoolSize) - ticketsPurchased
			if ticketsNeeded > 0 {
				ticketOuts = ticketOuts[1 : ticketsNeeded+1]
			} else {
				ticketOuts = nil
			}
		}
		ticketsPurchased += len(ticketOuts)

		blockName := fmt.Sprintf("bsv%d", i)
		g.NextBlock(blockName, nil, ticketOuts,
			chaingen.ReplaceBlockVersion(3))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight))
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach one block before the next stake
	// version interval with block version 3, stake version 0, and vote
	// version 3.
	//
	// This will result in triggering enforcement of the stake version and
	// that the stake version is 3.  The threshold state for the test dummy
	// deployments must still be defined since a v4 majority proof-of-work
	// and proof-of-stake upgrade are required before moving to started.
	// ---------------------------------------------------------------------

	blocksNeeded := stakeValidationHeight + stakeVerInterval - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtA%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(3),
			chaingen.ReplaceStakeVersion(0),
			chaingen.ReplaceVoteVersions(3))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + stakeVerInterval - 1))
	g.AssertBlockVersion(3)
	g.AssertStakeVersion(0)
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach one block before the next rule change
	// interval with block version 3, stake version 3, and vote version 3.
	//
	// The threshold state for the dummy deployments must still be defined
	// since it can only change on a rule change boundary and it requires a
	// v4 majority proof-of-work and proof-of-stake upgrade before moving to
	// started.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval - 2 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtB%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(3),
			chaingen.ReplaceStakeVersion(3),
			chaingen.ReplaceVoteVersions(3))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval - 2))
	g.AssertBlockVersion(3)
	g.AssertStakeVersion(3)
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach one block before the next stake
	// version interval with block version 3, stake version 3, and vote
	// version 4.
	//
	// This will result in achieving stake version 4 enforcement.
	//
	// The threshold state for the dummy deployments must still be defined
	// since it can only change on a rule change boundary and it still
	// requires a v4 majority proof-of-work upgrade before moving to
	// started.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + stakeVerInterval*4 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtC%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(3),
			chaingen.ReplaceStakeVersion(3),
			chaingen.ReplaceVoteVersions(4))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + stakeVerInterval*4 - 1))
	g.AssertBlockVersion(3)
	g.AssertStakeVersion(3)
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 3 majority, stake version 4, and vote version 4.  Set
	// the final two blocks to block version 4 so that majority version 4
	// is not achieved, but the final block in the interval is version 4.
	//
	// The threshold state for the dummy deployments must still be defined
	// since it still requires a v4 majority proof-of-work upgrade before
	// moving to started.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval*2 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtD%d", i)
		blockVersion := int32(3)
		if i >= blocksNeeded-2 {
			blockVersion = 4
		}
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(blockVersion),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVoteVersions(4))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*2 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to achieve proof-of-work block version lockin
	// with block version 4, stake version 4, and vote version 4.  Also, set
	// the vote bits to include yes votes for the first test dummy agenda
	// and no for the second test dummy agenda for an upcoming test.
	//
	// Since v4 majority proof-of-stake upgrade has been already been
	// achieved and this will achieve v4 majority proof-of-work upgrade,
	// voting can begin at the next rule change interval.
	//
	// The threshold state for the dummy deployments must still be defined
	// since even though all required upgrade conditions are met, the state
	// change must not happen until the start of the next rule change
	// interval.
	// ---------------------------------------------------------------------

	for i := int64(0); i < powNumToCheck; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtE%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(vbPrevBlockValid|vote1Yes.Bits|vote2No.Bits,
				4))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*2 -
		1 + powNumToCheck))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdDefined, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdDefined, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 4, stake version 4, and vote version 4.  Also, set the
	// vote bits to include yes votes for the first test dummy agenda and
	// no for the second test dummy agenda to ensure they aren't counted.
	//
	// The threshold state for the dummy deployments must move to started.
	// Even though the majority of the votes have already been voting yes
	// for the first test dummy agenda, and no for the second one, they must
	// not count, otherwise it would move straight to lockedin or failed,
	// respectively.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval*3 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtF%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(vbPrevBlockValid|vote1Yes.Bits|vote2No.Bits,
				4))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*3 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdStarted, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdStarted, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 4, stake version 4, and vote version 3.  Also, set the
	// vote bits to include yes votes for the first test dummy agenda and
	// no for the second test dummy agenda to ensure they aren't counted.
	//
	// The threshold state for the dummy deployments must remain in started
	// because the votes are an old version and thus have a different
	// definition and don't apply to version 4.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval*4 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtG%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(vbPrevBlockValid|vote1Yes.Bits|vote2No.Bits,
				3))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*4 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdStarted, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdStarted, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 4, stake version 4, and vote version 4.  Set the vote
	// bits such that quorum is not reached, but there is a majority yes
	// votes for the first test dummy agenda and a majority no for the
	// second test dummy agenda.
	//
	// The threshold state for the dummy deployments must remain in started
	// because quorum was not reached.
	// ---------------------------------------------------------------------

	var totalVotes int64
	blocksNeeded = stakeValidationHeight + ruleChangeInterval*5 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtH%d", i)
		voteBits := uint16(vbPrevBlockValid) // Abstain both test dummy
		if totalVotes+ticketsPerBlock < ruleChangeQuorum {
			voteBits = vbPrevBlockValid | vote1Yes.Bits | vote2No.Bits
		}
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(voteBits, 4))
		totalVotes += ticketsPerBlock
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*5 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdStarted, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdStarted, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 4, stake version 4, and vote version 4.  Set the vote
	// bits such that quorum is reached, but there are a few votes shy of a
	// majority yes for the first test dummy agenda and a few votes shy of a
	// majority no for the second test dummy agenda.
	//
	// The threshold state for the dummy deployments must remain in started
	// because even though quorum was reached, a required majority was not.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval*6 - 1 -
		int64(g.Tip().Header.Height)
	totalVotes = 0
	numActiveNeeded := ruleChangeQuorum * 2
	numMinorityNeeded := numActiveNeeded*ruleChangeMult/ruleChangeDiv - 1
	if numActiveNeeded > ticketsPerBlock*blocksNeeded {
		numActiveNeeded = ticketsPerBlock * blocksNeeded
	}
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtI%d", i)
		voteBits := uint16(vbPrevBlockValid) // Abstain both test dummy
		if totalVotes+ticketsPerBlock < numMinorityNeeded {
			voteBits = vbPrevBlockValid | vote1Yes.Bits | vote2No.Bits
		} else if totalVotes+ticketsPerBlock <= numActiveNeeded {
			voteBits = vbPrevBlockValid | vote1No.Bits | vote2Yes.Bits
		}
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(voteBits, 4))
		totalVotes += ticketsPerBlock
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*6 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdStarted, nil)
	g.TestThresholdStateChoice(vote2.Id, ThresholdStarted, nil)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 4, stake version 4, and vote version 4.  Also, set the
	// vote bits to yes for the first test dummy agenda and no to the second
	// one.
	//
	// The threshold state for the first dummy deployment must move to
	// lockedin since a majority yes vote was achieved while the second
	// dummy deployment must move to failed since a majority no vote was
	// achieved.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval*7 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtJ%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(vbPrevBlockValid|vote1Yes.Bits|vote2No.Bits,
				4))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*7 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdLockedIn, vote1Yes)
	g.TestThresholdStateChoice(vote2.Id, ThresholdFailed, vote2No)

	// ---------------------------------------------------------------------
	// Generate enough blocks to reach the next rule change interval with
	// block version 4, stake version 4, and vote version 4.  Also, set the
	// vote bits to include no votes for the first test dummy agenda and
	// yes votes for the second one.
	//
	// The threshold state for the first dummy deployment must move to
	// active since even though the interval had a majority no votes,
	// lockedin status has already been achieved and can't be undone without
	// a new agenda.  Similarly, the second one must remain in failed even
	// though the interval had a majority yes votes since a failed state
	// can't be undone.
	// ---------------------------------------------------------------------

	blocksNeeded = stakeValidationHeight + ruleChangeInterval*8 - 1 -
		int64(g.Tip().Header.Height)
	for i := int64(0); i < blocksNeeded; i++ {
		outs := g.OldestCoinbaseOuts()
		blockName := fmt.Sprintf("bsvtK%d", i)
		g.NextBlock(blockName, nil, outs[1:],
			chaingen.ReplaceBlockVersion(4),
			chaingen.ReplaceStakeVersion(4),
			chaingen.ReplaceVotes(vbPrevBlockValid|vote1No.Bits|vote2Yes.Bits,
				4))
		g.SaveTipCoinbaseOuts()
		g.AcceptTipBlock()
	}
	g.AssertTipHeight(uint32(stakeValidationHeight + ruleChangeInterval*8 - 1))
	g.AssertBlockVersion(4)
	g.AssertStakeVersion(4)
	g.TestThresholdStateChoice(vote1.Id, ThresholdActive, vote1Yes)
	g.TestThresholdStateChoice(vote2.Id, ThresholdFailed, vote2No)
}
