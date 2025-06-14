// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2015-2016 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"
)

// MsgFeeFilter implements the Message interface and represents a feefilter
// message.  It is used to request the receiving peer does not announce any
// transactions below the specified minimum fee rate.
//
// This message was not added until protocol versions starting with
// FeeFilterVersion.
type MsgFeeFilter struct {
	MinFee int64
}

// BtcDecode decodes r using the Vigil protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgFeeFilter) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgFeeFilter.BtcDecode"
	if pver < FeeFilterVersion {
		msg := fmt.Sprintf("feefilter message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	return readElement(r, &msg.MinFee)
}

// BtcEncode encodes the receiver to w using the Vigil protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgFeeFilter) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgFeeFilter.BtcEncode"
	if pver < FeeFilterVersion {
		msg := fmt.Sprintf("feefilter message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	return writeElement(w, msg.MinFee)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgFeeFilter) Command() string {
	return CmdFeeFilter
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgFeeFilter) MaxPayloadLength(pver uint32) uint32 {
	// 8 bytes min fee.
	return 8
}

// NewMsgFeeFilter returns a new Vigil feefilter message that conforms to the Message
// interface.  See MsgFeeFilter for details.
func NewMsgFeeFilter(minfee int64) *MsgFeeFilter {
	return &MsgFeeFilter{
		MinFee: minfee,
	}
}
