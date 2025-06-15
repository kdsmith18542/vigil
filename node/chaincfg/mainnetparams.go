// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chaincfg

import (
	"math/big"
	"time"

	"github.com/Vigil-Labs/vgl/chaincfg/chainhash"
	"github.com/Vigil-Labs/vgl/wire"
)

// MainNetParams returns the network parameters for the main Vigil network.
// This configuration represents the parameters for the production network.
func MainNetParams() *Params {
	// mainNetPowLimit is the highest proof of work value a Vigil block
	// can have for the main network. It is the value 2^224 - 1.
	mainNetPowLimit := new(big.Int).Sub(new(big.Int).Lsh(bigOne, 224), bigOne)

	// mainNetPowLimitBits is the main network proof of work limit in its
	// compact representation.
	const mainNetPowLimitBits = 0x1d00ffff // 486604799

	// genesisBlock defines the genesis block of the block chain which serves as
	// the public transaction ledger for the main network.
	genesisBlock := wire.MsgBlock{
		Header: wire.BlockHeader{
			Version:   1,
			PrevBlock: chainhash.Hash{},
			// MerkleRoot: Calculated below.
			Timestamp:    time.Unix(1735689600, 0), // 2025-01-01 00:00:00 +0000 UTC
			Bits:         mainNetPowLimitBits,
			SBits:        20000000, // 200 VGL (in atoms)
			Nonce:        0x00000000,
			StakeVersion: 1,
		},
		Transactions: []*wire.MsgTx{{
			SerType: wire.TxSerializeFull,
			Version: 1,
			TxIn: []*wire.TxIn{{
				// Fully null.
				PreviousOutPoint: wire.OutPoint{
					Hash:  chainhash.Hash{},
					Index: 0xffffffff,
					Tree:  0,
				},
				SignatureScript: hexDecode("0000"),
				Sequence:        0xffffffff,
				BlockHeight:     wire.NullBlockHeight,
				BlockIndex:      wire.NullBlockIndex,
				ValueIn:         wire.NullValueIn,
			}},
			TxOut: []*wire.TxOut{{
				Version: 0x0000,
				Value:   0x00000000,
				PkScript: hexDecode("801679e98561ada96caec2949a5d41c4cab3851e" +
					"b740d951c10ecbcf265c1fd9"),
			}},
			LockTime: 0,
			Expiry:   0,
		}},
	}
	genesisBlock.Header.MerkleRoot = genesisBlock.Transactions[0].TxHashFull()

	return &Params{
		Name:        "mainnet",
		Net:         wire.MainNet,
		DefaultPort: "9108",
		DNSSeeds: []DNSSeed{
			{"seed.vigil.network", true},
			{"seed1.vigil.network", true},
			{"seed2.vigil.network", true},
		},

		// Chain parameters
		GenesisBlock:         &genesisBlock,
		GenesisHash:          genesisBlock.BlockHash(),
		PowLimit:             mainNetPowLimit,
		PowLimitBits:         mainNetPowLimitBits,
		ReduceMinDifficulty:  false,
		GenerateSupported:    true,
		MaximumBlockSizes:    []int{1310720}, // 1.25MB
		MaxTxSize:            1000000,        // 1MB
		TargetTimePerBlock:   time.Minute * 2 + time.Second*30, // 2.5 minutes

		// Version 1 difficulty algorithm (EMA + BLAKE256) parameters.
		WorkDiffAlpha:            1,
		WorkDiffWindowSize:       144,
		WorkDiffWindows:          20,
		TargetTimespan:           time.Minute * 2 * 144, // TimePerBlock * WindowSize
		RetargetAdjustmentFactor: 4,

		// Version 2 difficulty algorithm (ASERT + KawPoW) parameters.
		WorkDiffV2Blake3StartBits: mainNetPowLimitBits,
		WorkDiffV2HalfLifeSecs:    8 * 60 * 60, // 8 hours for KawPoW ASERT

		// Subsidy parameters for VGL tokenomics
		BaseSubsidy:              2000000000, // 20 VGL (in atoms)
		MulSubsidy:               99,         // 1% reduction (99/100 = 0.99)
		DivSubsidy:               100,
		SubsidyReductionInterval: 6144,       // ~10.7 days at 2.5 min blocks
		WorkRewardProportion:     50,         // 50% to PoW miners (10 VGL)
		StakeRewardProportion:    40,         // 40% to PoS stakers (8 VGL)
		BlockTaxProportion:       10,         // 10% to treasury (2 VGL)

		// No V2 proportions needed as we're using the initial split

		// Checkpoints
		Checkpoints: []Checkpoint{
			// Add checkpoints as needed after mainnet launch
		},


		// AssumeValid is the hash of a block that has been externally verified
		// to be valid. This will be updated after mainnet launch.
		AssumeValid: *newHashFromStr("0000000000000000000000000000000000000000000000000000000000000000"),

		// MinKnownChainWork will be updated after mainnet launch
		MinKnownChainWork: hexToBigInt("0000000000000000000000000000000000000000000000000000000000000000"),

		// Consensus rule change deployments.
		RuleChangeActivationQuorum:     4032, // 10% of RuleChangeActivationInterval * TicketsPerBlock
		RuleChangeActivationMultiplier: 3,    // 75%
		RuleChangeActivationDivisor:    4,
		RuleChangeActivationInterval:   8064, // 2 weeks at 2.5 min blocks

		// Deployments define the consensus deployments for the network.
		Deployments: map[uint32][]ConsensusDeployment{
			// VGLP0005 (KawPoW) deployment
			5: {{
				Vote: Vote{
					Id:          VoteIDKawPoW,
					Description: "Activate KawPoW as the new proof of work algorithm",
					Mask:        0x0006, // Bits 1 and 2
					Choices: []Choice{{
						Id:          "abstain",
						Description: "abstain from voting for change",
						Bits:        0x0000,
						IsAbstain:   true,
						IsNo:        false,
					}, {
						Id:          "no",
						Description: "reject KawPoW",
						Bits:        0x0002, // Bit 1
						IsAbstain:   false,
						IsNo:        true,
					}, {
						Id:          "yes",
						Description: "support KawPoW",
						Bits:        0x0004, // Bit 2
						IsAbstain:   false,
						IsNo:        false,
					}},
				},
				// Activate from genesis
				StartTime:   1735689600, // 2025-01-01 00:00:00 +0000 UTC
				ExpireTime:  1767225600, // 2026-01-01 00:00:00 +0000 UTC
			}},

			// VGLP0011 (Blake3) deployment
			6: {{ // Use a new deployment ID, e.g., 6
				Vote: Vote{
					Id:          VoteIDBlake3Pow,
					Description: "Activate Blake3 as the new proof of work algorithm",
					Mask:        0x0006, // Bits 1 and 2
					Choices: []Choice{{
						Id:          "abstain",
						Description: "abstain from voting for change",
						Bits:        0x0000,
						IsAbstain:   true,
						IsNo:        false,
					}, {
						Id:          "no",
						Description: "reject Blake3",
						Bits:        0x0002, // Bit 1
						IsAbstain:   false,
						IsNo:        true,
					}, {
						Id:          "yes",
						Description: "support Blake3",
						Bits:        0x0004, // Bit 2
						IsAbstain:   false,
						IsNo:        false,
					}},
				},
				// Activate from genesis
				StartTime:   1735689600, // 2025-01-01 00:00:00 +0000 UTC
				ExpireTime:  1767225600, // 2026-01-01 00:00:00 +0000 UTC
			}},
		},

		// Enforce current block version once majority of the network has
		// upgraded.
		BlockEnforceNumRequired: 4032, // 75% of 5376
		BlockRejectNumRequired:  4608, // 85% of 5376
		BlockUpgradeNumToCheck:  5376, // ~9.3 days at 2.5 min blocks

		// AcceptNonStdTxs is a mempool param to either accept and relay
		// non standard txs to the network or reject them
		AcceptNonStdTxs: false,

		// Address encoding magics
		NetworkAddressPrefix: "vgl",
		PubKeyAddrID:        [2]byte{0x07, 0x3f}, // starts with Vk
		PubKeyHashAddrID:    [2]byte{0x0f, 0x12}, // starts with Vs
		PKHEdwardsAddrID:    [2]byte{0x0f, 0x01}, // starts with Ve
		PKHSchnorrAddrID:    [2]byte{0x0f, 0x1d}, // starts with VS
		ScriptHashAddrID:    [2]byte{0x0f, 0x1c}, // starts with Vc
		PrivateKeyID:        [2]byte{0x23, 0x0e}, // starts with Pv

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xprv
		HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

		// BIP44 coin type used in the hierarchical deterministic path for
		// address generation.
		SLIP0044CoinType: 1, // Testnet (all coins)
		LegacyCoinType:   1, // For backwards compatibility

		// PoS parameters
		MinimumStakeDiff:        20000000, // 200 VGL (in atoms)
		TicketPoolSize:          8192,     // ~2.5 days of tickets at 2.5 min blocks
		TicketsPerBlock:         5,        // 5 votes per block
		TicketMaturity:          256,      // ~10.7 hours at 2.5 min blocks
		TicketExpiry:            40960,    // ~71 days at 2.5 min blocks
		CoinbaseMaturity:        100,      // ~4.2 hours at 2.5 min blocks
		SStxChangeMaturity:      1,
		TicketPoolSizeWeight:     4,
		StakeDiffAlpha:           1,
		StakeDiffWindowSize:      144,
		StakeDiffWindows:         20,
		StakeVersionInterval:     144, // ~6 hours at 2.5 min blocks
		MaxFreshStakePerBlock:    20,  // 4*TicketsPerBlock
		StakeEnabledHeight:       256 + 256, // CoinbaseMaturity + TicketMaturity
		StakeValidationHeight:    4096, // ~1 week at 2.5 min blocks
		StakeBaseSigScript:      []byte{0x00, 0x00},
		StakeMajorityMultiplier: 3,
		StakeMajorityDivisor:    4,

		// Treasury parameters
		OrganizationPkScript:        []byte{0xa9, 0x14, 0x13, 0x5a, 0xb1, 0x88, 0xf3, 0x30, 0xce, 0x31, 0x9, 0x99, 0xe0, 0x30, 0x2b, 0x1d, 0x2b, 0xde, 0x1d, 0x58, 0xb4, 0x9d, 0x87},
		OrganizationPkScriptVersion: 0,
		BlockOneLedger:              []TokenPayout{},

		// Pi keys for block validation (to be set)
		PiKeys: [][]byte{},

		// Treasury voting parameters
		TreasuryVoteInterval:          144, // ~6 hours at 2.5 min blocks
		TreasuryVoteIntervalMultiplier: 2,
		TreasuryVoteQuorumMultiplier:   3,
		TreasuryVoteQuorumDivisor:     4,
		TreasuryVoteRequiredMultiplier: 3,
		TreasuryVoteRequiredDivisor:    4,
		TreasuryExpenditureWindow:     144, // ~6 hours at 2.5 min blocks
		TreasuryExpenditurePolicy:     5,   // 5 VGL per block max
		TreasuryExpenditureBootstrap:  5000000000, // 50,000 VGL

		// Seeders for DNS discovery
		seeders: []string{
			"seed.vigil.network:9108",
			"seed1.vigil.network:9108",
			"seed2.vigil.network:9108",
		},
	}
}




