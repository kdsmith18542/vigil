// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"bytes"
	"testing"
)

func TestMemPool(t *testing.T) {
	pver := ProtocolVersion

	// Ensure the command is expected value.
	wantCmd := "mempool"
	msg := NewMsgMemPool()
	if cmd := msg.Command(); cmd != wantCmd {
		t.Errorf("NewMsgMemPool: wrong command - got %v want %v",
			cmd, wantCmd)
	}

	// Ensure max payload is expected value.
	wantPayload := uint32(0)
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

	// Test encode with latest protocol version.
	var buf bytes.Buffer
	err := msg.BtcEncode(&buf, pver)
	if err != nil {
		t.Errorf("encode of MsgMemPool failed %v err <%v>", msg, err)
	}

	// Test decode with latest protocol version.
	readmsg := NewMsgMemPool()
	err = readmsg.BtcDecode(&buf, pver)
	if err != nil {
		t.Errorf("decode of MsgMemPool failed [%v] err <%v>", buf, err)
	}
}
