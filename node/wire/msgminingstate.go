// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2015-2016 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"

	"github.com/kdsmith18542/vigil-Labs/vgl/node/chaincfg/chainhash"
)

// MaxMSBlocksAtHeadPerMsg is the maximum number of block hashes allowed
// per message.
const MaxMSBlocksAtHeadPerMsg = 8

// MaxMSVotesAtHeadPerMsg is the maximum number of votes at head per message.
const MaxMSVotesAtHeadPerMsg = 40 // 8 * 5

// MsgMiningState implements the Message interface and represents a mining state
// message.  It is used to request a list of blocks located at the chain tip
// along with all votes for those blocks.  The list is returned is limited by
// the maximum number of blocks per message and the maximum number of votes per
// message.
type MsgMiningState struct {
	Version     uint32
	Height      uint32
	BlockHashes []*chainhash.Hash
	VoteHashes  []*chainhash.Hash
}

// AddBlockHash adds a new block hash to the message.
func (msg *MsgMiningState) AddBlockHash(hash *chainhash.Hash) error {
	const op = "MsgMiningState.AddBlockHash"
	if len(msg.BlockHashes)+1 > MaxMSBlocksAtHeadPerMsg {
		msg := fmt.Sprintf("too many block hashes for message [max %v]",
			MaxMSBlocksAtHeadPerMsg)
		return messageError(op, ErrTooManyHeaders, msg)
	}

	msg.BlockHashes = append(msg.BlockHashes, hash)
	return nil
}

// AddVoteHash adds a new vote hash to the message.
func (msg *MsgMiningState) AddVoteHash(hash *chainhash.Hash) error {
	const op = "MsgMiningState.AddVoteHash"
	if len(msg.VoteHashes)+1 > MaxMSVotesAtHeadPerMsg {
		msg := fmt.Sprintf("too many vote hashes for message [max %v]",
			MaxMSVotesAtHeadPerMsg)
		return messageError(op, ErrTooManyVotes, msg)
	}

	msg.VoteHashes = append(msg.VoteHashes, hash)
	return nil
}

// BtcDecode decodes r using the Vigil protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgMiningState) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgMiningState.BtcDecode"
	err := readElement(r, &msg.Version)
	if err != nil {
		return err
	}

	err = readElement(r, &msg.Height)
	if err != nil {
		return err
	}

	// Read num block hashes and limit to max.
	count, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > MaxMSBlocksAtHeadPerMsg {
		msg := fmt.Sprintf("too many block hashes for message "+
			"[count %v, max %v]", count, MaxMSBlocksAtHeadPerMsg)
		return messageError(op, ErrTooManyBlocks, msg)
	}

	msg.BlockHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := chainhash.Hash{}
		err := readElement(r, &hash)
		if err != nil {
			return err
		}
		msg.AddBlockHash(&hash)
	}

	// Read num vote hashes and limit to max.
	count, err = ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > MaxMSVotesAtHeadPerMsg {
		msg := fmt.Sprintf("too many vote hashes for message "+
			"[count %v, max %v]", count, MaxMSVotesAtHeadPerMsg)
		return messageError(op, ErrTooManyVotes, msg)
	}

	msg.VoteHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := chainhash.Hash{}
		err := readElement(r, &hash)
		if err != nil {
			return err
		}
		err = msg.AddVoteHash(&hash)
		if err != nil {
			return err
		}
	}

	return nil
}

// BtcEncode encodes the receiver to w using the Vigil protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgMiningState) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgMiningState.BtcEncode"
	err := writeElement(w, msg.Version)
	if err != nil {
		return err
	}

	err = writeElement(w, msg.Height)
	if err != nil {
		return err
	}

	// Write block hashes.
	count := len(msg.BlockHashes)
	if count > MaxMSBlocksAtHeadPerMsg {
		msg := fmt.Sprintf("too many block hashes for message "+
			"[count %v, max %v]", count, MaxMSBlocksAtHeadPerMsg)
		return messageError(op, ErrTooManyBlocks, msg)
	}

	err = WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.BlockHashes {
		err = writeElement(w, hash)
		if err != nil {
			return err
		}
	}

	// Write vote hashes.
	count = len(msg.VoteHashes)
	if count > MaxMSVotesAtHeadPerMsg {
		msg := fmt.Sprintf("too many vote hashes for message "+
			"[count %v, max %v]", count, MaxMSVotesAtHeadPerMsg)
		return messageError(op, ErrTooManyVotes, msg)
	}

	err = WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.VoteHashes {
		err = writeElement(w, hash)
		if err != nil {
			return err
		}
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgMiningState) Command() string {
	return CmdMiningState
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgMiningState) MaxPayloadLength(pver uint32) uint32 {
	// Protocol version 4 bytes + Height 4 bytes + num block hashes (varInt)
	// 1 byte + block hashes + num vote hashes (varInt) 1 byte + vote hashes
	return 4 + 4 + uint32(VarIntSerializeSize(MaxMSBlocksAtHeadPerMsg)) +
		(MaxMSBlocksAtHeadPerMsg * chainhash.HashSize) +
		uint32(VarIntSerializeSize(MaxMSVotesAtHeadPerMsg)) +
		(MaxMSVotesAtHeadPerMsg * chainhash.HashSize)
}

// NewMsgMiningState returns a new Vigil miningstate message that conforms to
// the Message interface using the defaults for the fields.
func NewMsgMiningState() *MsgMiningState {
	return &MsgMiningState{
		Version:     1,
		Height:      0,
		BlockHashes: make([]*chainhash.Hash, 0, MaxMSBlocksAtHeadPerMsg),
		VoteHashes:  make([]*chainhash.Hash, 0, MaxMSVotesAtHeadPerMsg),
	}
}
