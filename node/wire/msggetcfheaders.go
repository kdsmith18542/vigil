// Copyright (c) 2017 The btcsuite developers
// Copyright (c) 2017 The Lightning Network Developers
// Copyright (c) 2018-2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"

	"github.com/kdsmith18542/vigil-Labs/vgl/node/chaincfg/chainhash"
)

// MsgGetCFHeaders is a message similar to MsgGetHeaders, but for committed
// filter headers. It allows to set the FilterType field to get headers in the
// chain of basic (0x00) or extended (0x01) headers.
//
// Deprecated: This message is no longer valid as of protocol version
// CFilterV2Version.
type MsgGetCFHeaders struct {
	BlockLocatorHashes []*chainhash.Hash
	HashStop           chainhash.Hash
	FilterType         FilterType
}

// AddBlockLocatorHash adds a new block locator hash to the message.
func (msg *MsgGetCFHeaders) AddBlockLocatorHash(hash *chainhash.Hash) error {
	const op = "MsgGetCFHeaders.AddBlockLocatorHash"
	if len(msg.BlockLocatorHashes)+1 > MaxBlockLocatorsPerMsg {
		msg := fmt.Sprintf("too many block locator hashes for message [max %v]",
			MaxBlockLocatorsPerMsg)
		return messageError(op, ErrTooManyLocators, msg)
	}

	msg.BlockLocatorHashes = append(msg.BlockLocatorHashes, hash)
	return nil
}

// BtcDecode decodes r using the Vigil protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgGetCFHeaders.BtcDecode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("getcfheaders message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	// Read num block locator hashes and limit to max.
	count, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}
	if count > MaxBlockLocatorsPerMsg {
		msg := fmt.Sprintf("too many block locator hashes for message "+
			"[count %v, max %v]", count, MaxBlockLocatorsPerMsg)
		return messageError(op, ErrTooManyLocators, msg)
	}

	// Create a contiguous slice of hashes to deserialize into in order to
	// reduce the number of allocations.
	locatorHashes := make([]chainhash.Hash, count)
	msg.BlockLocatorHashes = make([]*chainhash.Hash, 0, count)
	for i := uint64(0); i < count; i++ {
		hash := &locatorHashes[i]
		err := readElement(r, hash)
		if err != nil {
			return err
		}
		msg.AddBlockLocatorHash(hash)
	}

	err = readElement(r, &msg.HashStop)
	if err != nil {
		return err
	}

	return readElement(r, (*uint8)(&msg.FilterType))
}

// BtcEncode encodes the receiver to w using the Vigil protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgGetCFHeaders.BtcEncode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("getcfheaders message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	// Limit to max block locator hashes per message.
	count := len(msg.BlockLocatorHashes)
	if count > MaxBlockLocatorsPerMsg {
		msg := fmt.Sprintf("too many block locator hashes for message "+
			"[count %v, max %v]", count, MaxBlockLocatorsPerMsg)
		return messageError(op, ErrTooManyLocators, msg)
	}

	err := WriteVarInt(w, pver, uint64(count))
	if err != nil {
		return err
	}

	for _, hash := range msg.BlockLocatorHashes {
		err := writeElement(w, hash)
		if err != nil {
			return err
		}
	}

	err = writeElement(w, &msg.HashStop)
	if err != nil {
		return err
	}

	return binarySerializer.PutUint8(w, uint8(msg.FilterType))
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetCFHeaders) Command() string {
	return CmdGetCFHeaders
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFHeaders) MaxPayloadLength(pver uint32) uint32 {
	// Num block locator hashes (varInt) 3 bytes + max allowed
	// block locators + hash stop + filter type 1 byte.
	return uint32(VarIntSerializeSize(MaxBlockLocatorsPerMsg)) +
		(MaxBlockLocatorsPerMsg * chainhash.HashSize) + chainhash.HashSize + 1
}

// NewMsgGetCFHeaders returns a new Vigil getcfheader message that conforms to the
// Message interface using the passed parameters and defaults for the remaining
// fields.
//
// Deprecated: This message is no longer valid as of protocol version
// CFilterV2Version.
func NewMsgGetCFHeaders() *MsgGetCFHeaders {
	return &MsgGetCFHeaders{
		BlockLocatorHashes: make([]*chainhash.Hash, 0,
			MaxBlockLocatorsPerMsg),
	}
}
