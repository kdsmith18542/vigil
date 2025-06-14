// Copyright (c) 2017 The btcsuite developers
// Copyright (c) 2017 The Lightning Network Developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"
)

// MaxFilterTypesPerMsg is the maximum number of filter types allowed per
// message.
const MaxFilterTypesPerMsg = 256

// FilterType is used to represent a filter type.
type FilterType uint8

const (
	// GCSFilterRegular is the regular filter type.
	GCSFilterRegular FilterType = iota

	// GCSFilterExtended is the extended filter type.
	GCSFilterExtended
)

// MsgCFTypes is the cftypes message.
//
// Deprecated: This message is no longer valid as of protocol version
// CFilterV2Version.
type MsgCFTypes struct {
	SupportedFilters []FilterType
}

// BtcDecode decodes r using the wire protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgCFTypes) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgCFTypes.BtcDecode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("cftypes message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	// Read the number of filter types supported.  The count may not exceed the
	// total number of filters that can be represented by a FilterType byte.
	count, err := ReadVarInt(r, pver)
	if err != nil {
		return err
	}

	if count > MaxFilterTypesPerMsg {
		msg := fmt.Sprintf("too many filter types for message "+
			"[count %v, max %v]", count, MaxFilterTypesPerMsg)
		return messageError(op, ErrTooManyFilterTypes, msg)
	}

	// Read each filter type.
	msg.SupportedFilters = make([]FilterType, count)
	for i := range msg.SupportedFilters {
		err = readElement(r, (*uint8)(&msg.SupportedFilters[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// BtcEncode encodes the receiver to w using the wire protocol encoding. This is
// part of the Message interface implementation.
func (msg *MsgCFTypes) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgCFTypes.BtcEncode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("cftypes message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	if len(msg.SupportedFilters) > MaxFilterTypesPerMsg {
		msg := fmt.Sprintf("too many filter types for message "+
			"[count %v, max %v]", len(msg.SupportedFilters), MaxFilterTypesPerMsg)
		return messageError(op, ErrTooManyFilterTypes, msg)
	}

	// Write length of supported filters slice.
	err := WriteVarInt(w, pver, uint64(len(msg.SupportedFilters)))
	if err != nil {
		return err
	}

	for i := range msg.SupportedFilters {
		err = binarySerializer.PutUint8(w, uint8(msg.SupportedFilters[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// Deserialize decodes a filter from r into the receiver using a format that is
// suitable for long-term storage such as a database. This function differs from
// BtcDecode in that BtcDecode decodes from the wire protocol as it was sent
// across the network.  The wire encoding can technically differ depending on
// the protocol version and doesn't even really need to match the format of a
// stored filter at all. As of the time this comment was written, the encoded
// filter is the same in both instances, but there is a distinct difference and
// separating the two allows the API to be flexible enough to deal with changes.
func (msg *MsgCFTypes) Deserialize(r io.Reader) error {
	// At the current time, there is no difference between the wire encoding
	// and the stable long-term storage format.  As a result, make use of
	// BtcDecode.
	return msg.BtcDecode(r, 0)
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgCFTypes) Command() string {
	return CmdCFTypes
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver. This is part of the Message interface implementation.
func (msg *MsgCFTypes) MaxPayloadLength(pver uint32) uint32 {
	// 3 bytes for filter count, 1 byte up to 256 bytes filter types.
	return uint32(VarIntSerializeSize(MaxFilterTypesPerMsg)) +
		MaxFilterTypesPerMsg
}

// NewMsgCFTypes returns a new cftypes message that conforms to the Message
// interface. See MsgCFTypes for details.
//
// Deprecated: This message is no longer valid as of protocol version
// CFilterV2Version.
func NewMsgCFTypes(filterTypes []FilterType) *MsgCFTypes {
	return &MsgCFTypes{
		SupportedFilters: filterTypes,
	}
}
