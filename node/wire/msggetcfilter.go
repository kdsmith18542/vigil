// Copyright (c) 2017 The btcsuite developers
// Copyright (c) 2017 The Lightning Network Developers
// Copyright (c) 2017-2020 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"

	"github.com/kdsmith18542/vigil-Labs/vgl/node/chaincfg/chainhash"
)

// MsgGetCFilter implements the Message interface and represents a vigilgetcfilter
// message. It is used to request a committed filter for a block.
//
// Deprecated: This message is no longer valid as of protocol version
// CFilterV2Version.
type MsgGetCFilter struct {
	BlockHash  chainhash.Hash
	FilterType FilterType
}

// BtcDecode decodes r using the Vigil protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgGetCFilter) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgGetCFilter.BtcDecode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("getcfilter message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	err := readElement(r, &msg.BlockHash)
	if err != nil {
		return err
	}
	return readElement(r, (*uint8)(&msg.FilterType))
}

// BtcEncode encodes the receiver to w using the Vigil protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetCFilter) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgGetCFilter.BtcEncode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("getcfilter message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	err := writeElement(w, &msg.BlockHash)
	if err != nil {
		return err
	}
	return binarySerializer.PutUint8(w, uint8(msg.FilterType))
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetCFilter) Command() string {
	return CmdGetCFilter
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFilter) MaxPayloadLength(pver uint32) uint32 {
	// Block hash + filter type.
	return chainhash.HashSize + 1
}

// NewMsgGetCFilter returns a new Vigil getcfilter message that conforms to the
// Message interface using the passed parameters and defaults for the remaining
// fields.
//
// Deprecated: This message is no longer valid as of protocol version
// CFilterV2Version.
func NewMsgGetCFilter(blockHash *chainhash.Hash, filterType FilterType) *MsgGetCFilter {
	return &MsgGetCFilter{
		BlockHash:  *blockHash,
		FilterType: filterType,
	}
}
