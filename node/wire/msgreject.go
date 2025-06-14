// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2023 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"

	"github.com/kdsmith18542/vigil-Labs/vigil/chaincfg/chainhash"
)

// RejectCode represents a numeric value by which a remote peer indicates
// why a message was rejected.
type RejectCode uint8

// These constants define the various supported reject codes.
const (
	RejectMalformed       RejectCode = 0x01
	RejectInvalid         RejectCode = 0x10
	RejectObsolete        RejectCode = 0x11
	RejectDuplicate       RejectCode = 0x12
	RejectNonstandard     RejectCode = 0x40
	RejectDust            RejectCode = 0x41
	RejectInsufficientFee RejectCode = 0x42
	RejectCheckpoint      RejectCode = 0x43
)

// Map of reject codes back strings for pretty printing.
var rejectCodeStrings = map[RejectCode]string{
	RejectMalformed:       "REJECT_MALFORMED",
	RejectInvalid:         "REJECT_INVALID",
	RejectObsolete:        "REJECT_OBSOLETE",
	RejectDuplicate:       "REJECT_DUPLICATE",
	RejectNonstandard:     "REJECT_NONSTANDARD",
	RejectDust:            "REJECT_DUST",
	RejectInsufficientFee: "REJECT_INSUFFICIENTFEE",
	RejectCheckpoint:      "REJECT_CHECKPOINT",
}

// String returns the RejectCode in human-readable form.
func (code RejectCode) String() string {
	if s, ok := rejectCodeStrings[code]; ok {
		return s
	}

	return fmt.Sprintf("Unknown RejectCode (%d)", uint8(code))
}

// MsgReject implements the Message interface and represents a Vigil reject
// message.
//
// Deprecated: This message is no longer valid as of protocol version
// RemoveRejectVersion.
type MsgReject struct {
	// Cmd is the command for the message which was rejected such as CmdBlock or
	// CmdTx.  This can be obtained from the Command function of a Message.
	Cmd string

	// RejectCode is a code indicating why the command was rejected.  It
	// is encoded as a uint8 on the wire.
	Code RejectCode

	// Reason is a human-readable string with specific details (over and
	// above the reject code) about why the command was rejected.
	Reason string

	// Hash identifies a specific block or transaction that was rejected
	// and therefore only applies the MsgBlock and MsgTx messages.
	Hash chainhash.Hash
}

// validateRejectCommand ensures the provided reject command conforms to the
// results imposed by the protocol.
func validateRejectCommand(cmd string) error {
	const op = "MsgReject.validateRejectCommand"
	if !isStrictAscii(cmd) {
		msg := "reject command is not strict ASCII"
		return messageError(op, ErrMalformedStrictString, msg)
	}

	return nil
}

// validateRejectReason ensures the provided reject reason conforms to the
// results imposed by the protocol.
func validateRejectReason(reason string) error {
	const op = "MsgReject.validateRejectReason"
	if !isStrictAscii(reason) {
		msg := "reject reason is not strict ASCII"
		return messageError(op, ErrMalformedStrictString, msg)
	}

	return nil
}

// BtcDecode decodes r using the Vigil protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgReject) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgReject.BtcDecode"
	if pver >= RemoveRejectVersion {
		str := fmt.Sprintf("%s message invalid for protocol version %d",
			msg.Command(), pver)
		return messageError(op, ErrMsgInvalidForPVer, str)
	}

	// Command that was rejected.
	cmd, err := ReadVarString(r, pver)
	if err != nil {
		return err
	}
	if err := validateRejectCommand(cmd); err != nil {
		return err
	}
	msg.Cmd = cmd

	// Code indicating why the command was rejected.
	err = readElement(r, &msg.Code)
	if err != nil {
		return err
	}

	// Human readable string with specific details (over and above the
	// reject code above) about why the command was rejected.
	reason, err := ReadVarString(r, pver)
	if err != nil {
		return err
	}
	if err := validateRejectReason(reason); err != nil {
		return err
	}
	msg.Reason = reason

	// CmdBlock and CmdTx messages have an additional hash field that
	// identifies the specific block or transaction.
	if msg.Cmd == CmdBlock || msg.Cmd == CmdTx {
		err := readElement(r, &msg.Hash)
		if err != nil {
			return err
		}
	}

	return nil
}

// BtcEncode encodes the receiver to w using the Vigil protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgReject) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgReject.BtcEncode"
	if pver >= RemoveRejectVersion {
		str := fmt.Sprintf("%s message invalid for protocol version %d",
			msg.Command(), pver)
		return messageError(op, ErrMsgInvalidForPVer, str)
	}

	if err := validateRejectCommand(msg.Cmd); err != nil {
		return err
	}
	if err := validateRejectReason(msg.Reason); err != nil {
		return err
	}

	// Command that was rejected.
	err := WriteVarString(w, pver, msg.Cmd)
	if err != nil {
		return err
	}

	// Code indicating why the command was rejected.
	err = writeElement(w, msg.Code)
	if err != nil {
		return err
	}

	// Human readable string with specific details (over and above the
	// reject code above) about why the command was rejected.
	err = WriteVarString(w, pver, msg.Reason)
	if err != nil {
		return err
	}

	// CmdBlock and CmdTx messages have an additional hash field that
	// identifies the specific block or transaction.
	if msg.Cmd == CmdBlock || msg.Cmd == CmdTx {
		err := writeElement(w, &msg.Hash)
		if err != nil {
			return err
		}
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgReject) Command() string {
	return CmdReject
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgReject) MaxPayloadLength(pver uint32) uint32 {
	if pver >= RemoveRejectVersion {
		return 0
	}

	// Unfortunately the Vigil protocol does not enforce a sane
	// limit on the length of the reason, so the max payload is the
	// overall maximum message payload.
	return uint32(MaxMessagePayload)
}

// NewMsgReject returns a new Vigil reject message that conforms to the
// Message interface.  See MsgReject for details.
func NewMsgReject(command string, code RejectCode, reason string) *MsgReject {
	return &MsgReject{
		Cmd:    command,
		Code:   code,
		Reason: reason,
	}
}
