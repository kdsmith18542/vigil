// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2020 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/kdsmith18542/vigil/chaincfg/chainhash"
)

// TestGetBlocks tests the MsgGetBlocks API.
func TestGetBlocks(t *testing.T) {
	pver := ProtocolVersion

	// Block 99500 hash.
	hashStr := "000000000002e7ad7b9eef9479e4aabc65cb831269cc20d2632c13684406dee0"
	locatorHash, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Block 100000 hash.
	hashStr = "3ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"
	hashStop, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Ensure we get the same data back out.
	msg := NewMsgGetBlocks(hashStop)
	if !msg.HashStop.IsEqual(hashStop) {
		t.Errorf("NewMsgGetBlocks: wrong stop hash - got %v, want %v",
			msg.HashStop, hashStop)
	}

	// Ensure the command is expected value.
	wantCmd := "getblocks"
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgGetBlocks: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value for latest protocol version.
	// Protocol version 4 bytes + num hashes (varInt) 3 bytes + max block
	// locator hashes + hash stop.
	wantPayload := uint32(16039)
	maxPayload := msg.MaxPayloadLength(pver)
	if maxPayload != wantPayload {
		t.Errorf("MaxPayloadLength: wrong max payload length for "+
			"protocol version %d - got %v, want %v", pver,
			maxPayload, wantPayload)
	}

	// Ensure max payload length is not more than MaxMessagePayload.
	if maxPayload > MaxMessagePayload {
		t.Fatalf("MaxPayloadLength: payload length (%v) for protocol "+
			"version %d exceeds MaxMessagePayload (%v).", maxPayload, pver,
			MaxMessagePayload)
	}

	// Ensure block locator hashes are added properly.
	err = msg.AddBlockLocatorHash(locatorHash)
	if err != nil {
		t.Errorf("AddBlockLocatorHash: %v", err)
	}
	if msg.BlockLocatorHashes[0] != locatorHash {
		t.Errorf("AddBlockLocatorHash: wrong block locator added - "+
			"got %v, want %v",
			spew.Sprint(msg.BlockLocatorHashes[0]),
			spew.Sprint(locatorHash))
	}

	// Ensure adding more than the max allowed block locator hashes per
	// message returns an error.
	for i := 0; i < MaxBlockLocatorsPerMsg; i++ {
		err = msg.AddBlockLocatorHash(locatorHash)
	}
	if err == nil {
		t.Errorf("AddBlockLocatorHash: expected error on too many " +
			"block locator hashes not received")
	}
}

// TestGetBlocksWire tests the MsgGetBlocks wire encode and decode for various
// numbers of block locator hashes and protocol versions.
func TestGetBlocksWire(t *testing.T) {
	// Set protocol inside getblocks message.
	pver := uint32(60002)

	// Block 99499 hash.
	hashStr := "2710f40c87ec93d010a6fd95f42c59a2cbacc60b18cf6b7957535"
	hashLocator, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Block 99500 hash.
	hashStr = "2e7ad7b9eef9479e4aabc65cb831269cc20d2632c13684406dee0"
	hashLocator2, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Block 100000 hash.
	hashStr = "3ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"
	hashStop, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// MsgGetBlocks message with no block locators or stop hash.
	noLocators := NewMsgGetBlocks(&chainhash.Hash{})
	noLocators.ProtocolVersion = pver
	noLocatorsEncoded := []byte{
		0x62, 0xea, 0x00, 0x00, // Protocol version 60002
		0x00, // Varint for number of block locator hashes
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Hash stop
	}

	// MsgGetBlocks message with multiple block locators and a stop hash.
	multiLocators := NewMsgGetBlocks(hashStop)
	multiLocators.AddBlockLocatorHash(hashLocator2)
	multiLocators.AddBlockLocatorHash(hashLocator)
	multiLocators.ProtocolVersion = pver
	multiLocatorsEncoded := []byte{
		0x62, 0xea, 0x00, 0x00, // Protocol version 60002
		0x02, // Varint for number of block locator hashes
		0xe0, 0xde, 0x06, 0x44, 0x68, 0x13, 0x2c, 0x63,
		0xd2, 0x20, 0xcc, 0x69, 0x12, 0x83, 0xcb, 0x65,
		0xbc, 0xaa, 0xe4, 0x79, 0x94, 0xef, 0x9e, 0x7b,
		0xad, 0xe7, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, // Block 99500 hash
		0x35, 0x75, 0x95, 0xb7, 0xf6, 0x8c, 0xb1, 0x60,
		0xcc, 0xba, 0x2c, 0x9a, 0xc5, 0x42, 0x5f, 0xd9,
		0x6f, 0x0a, 0x01, 0x3d, 0xc9, 0x7e, 0xc8, 0x40,
		0x0f, 0x71, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, // Block 99499 hash
		0x06, 0xe5, 0x33, 0xfd, 0x1a, 0xda, 0x86, 0x39,
		0x1f, 0x3f, 0x6c, 0x34, 0x32, 0x04, 0xb0, 0xd2,
		0x78, 0xd4, 0xaa, 0xec, 0x1c, 0x0b, 0x20, 0xaa,
		0x27, 0xba, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, // Hash stop
	}

	tests := []struct {
		in   *MsgGetBlocks // Message to encode
		out  *MsgGetBlocks // Expected decoded message
		buf  []byte        // Wire encoding
		pver uint32        // Protocol version for wire encoding
	}{
		// Latest protocol version with no block locators.
		{
			noLocators,
			noLocators,
			noLocatorsEncoded,
			ProtocolVersion,
		},

		// Latest protocol version with multiple block locators.
		{
			multiLocators,
			multiLocators,
			multiLocatorsEncoded,
			ProtocolVersion,
		},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode the message to wire format.
		var buf bytes.Buffer
		err := test.in.BtcEncode(&buf, test.pver)
		if err != nil {
			t.Errorf("BtcEncode #%d error %v", i, err)
			continue
		}
		if !bytes.Equal(buf.Bytes(), test.buf) {
			t.Errorf("BtcEncode #%d\n got: %s want: %s", i,
				spew.Sdump(buf.Bytes()), spew.Sdump(test.buf))
			continue
		}

		// Decode the message from wire format.
		var msg MsgGetBlocks
		rbuf := bytes.NewReader(test.buf)
		err = msg.BtcDecode(rbuf, test.pver)
		if err != nil {
			t.Errorf("BtcDecode #%d error %v", i, err)
			continue
		}
		if !reflect.DeepEqual(&msg, test.out) {
			t.Errorf("BtcDecode #%d\n got: %s want: %s", i,
				spew.Sdump(&msg), spew.Sdump(test.out))
			continue
		}
	}
}

// TestGetBlocksWireErrors performs negative tests against wire encode and
// decode of MsgGetBlocks to confirm error paths work correctly.
func TestGetBlocksWireErrors(t *testing.T) {
	// Set protocol inside getheaders message.  Use protocol version 60002
	// specifically here instead of the latest because the test data is
	// using bytes encoded with that protocol version.
	pver := uint32(60002)

	// Block 99499 hash.
	hashStr := "2710f40c87ec93d010a6fd95f42c59a2cbacc60b18cf6b7957535"
	hashLocator, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Block 99500 hash.
	hashStr = "2e7ad7b9eef9479e4aabc65cb831269cc20d2632c13684406dee0"
	hashLocator2, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// Block 100000 hash.
	hashStr = "3ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"
	hashStop, err := chainhash.NewHashFromStr(hashStr)
	if err != nil {
		t.Errorf("NewHashFromStr: %v", err)
	}

	// MsgGetBlocks message with multiple block locators and a stop hash.
	baseGetBlocks := NewMsgGetBlocks(hashStop)
	baseGetBlocks.ProtocolVersion = pver
	baseGetBlocks.AddBlockLocatorHash(hashLocator2)
	baseGetBlocks.AddBlockLocatorHash(hashLocator)
	baseGetBlocksEncoded := []byte{
		0x62, 0xea, 0x00, 0x00, // Protocol version 60002
		0x02, // Varint for number of block locator hashes
		0xe0, 0xde, 0x06, 0x44, 0x68, 0x13, 0x2c, 0x63,
		0xd2, 0x20, 0xcc, 0x69, 0x12, 0x83, 0xcb, 0x65,
		0xbc, 0xaa, 0xe4, 0x79, 0x94, 0xef, 0x9e, 0x7b,
		0xad, 0xe7, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, // Block 99500 hash
		0x35, 0x75, 0x95, 0xb7, 0xf6, 0x8c, 0xb1, 0x60,
		0xcc, 0xba, 0x2c, 0x9a, 0xc5, 0x42, 0x5f, 0xd9,
		0x6f, 0x0a, 0x01, 0x3d, 0xc9, 0x7e, 0xc8, 0x40,
		0x0f, 0x71, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, // Block 99499 hash
		0x06, 0xe5, 0x33, 0xfd, 0x1a, 0xda, 0x86, 0x39,
		0x1f, 0x3f, 0x6c, 0x34, 0x32, 0x04, 0xb0, 0xd2,
		0x78, 0xd4, 0xaa, 0xec, 0x1c, 0x0b, 0x20, 0xaa,
		0x27, 0xba, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, // Hash stop
	}

	// Message that forces an error by having more than the max allowed
	// block locator hashes.
	maxGetBlocks := NewMsgGetBlocks(hashStop)
	for i := 0; i < MaxBlockLocatorsPerMsg; i++ {
		maxGetBlocks.AddBlockLocatorHash(&mainNetGenesisHash)
	}
	maxGetBlocks.BlockLocatorHashes = append(maxGetBlocks.BlockLocatorHashes,
		&mainNetGenesisHash)
	maxGetBlocksEncoded := []byte{
		0x62, 0xea, 0x00, 0x00, // Protocol version 60002
		0xfd, 0xf5, 0x01, // Varint for number of block loc hashes (501)
	}

	tests := []struct {
		in       *MsgGetBlocks // Value to encode
		buf      []byte        // Wire encoding
		pver     uint32        // Protocol version for wire encoding
		max      int           // Max size of fixed buffer to induce errors
		writeErr error         // Expected write error
		readErr  error         // Expected read error
	}{
		// Force error in protocol version.
		{baseGetBlocks, baseGetBlocksEncoded, pver, 0, io.ErrShortWrite, io.EOF},
		// Force error in block locator hash count.
		{baseGetBlocks, baseGetBlocksEncoded, pver, 4, io.ErrShortWrite, io.EOF},
		// Force error in block locator hashes.
		{baseGetBlocks, baseGetBlocksEncoded, pver, 5, io.ErrShortWrite, io.EOF},
		// Force error in stop hash.
		{baseGetBlocks, baseGetBlocksEncoded, pver, 69, io.ErrShortWrite, io.EOF},
		// Force error with greater than max block locator hashes.
		{maxGetBlocks, maxGetBlocksEncoded, pver, 7, ErrTooManyLocators, ErrTooManyLocators},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		// Encode to wire format.
		w := newFixedWriter(test.max)
		err := test.in.BtcEncode(w, test.pver)
		if !errors.Is(err, test.writeErr) {
			t.Errorf("BtcEncode #%d wrong error got: %v, want: %v", i, err,
				test.writeErr)
			continue
		}

		// Decode from wire format.
		var msg MsgGetBlocks
		r := newFixedReader(test.max, test.buf)
		err = msg.BtcDecode(r, test.pver)
		if !errors.Is(err, test.readErr) {
			t.Errorf("BtcDecode #%d wrong error got: %v, want: %v", i, err,
				test.readErr)
			continue
		}
	}
}
