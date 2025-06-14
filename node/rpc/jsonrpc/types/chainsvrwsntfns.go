// Copyright (c) 2014 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// NOTE: This file is intended to house the RPC websocket notifications that are
// supported by a chain server.

package types

import "github.com/kdsmith18542/vigil/VGLjson/v4"

const (
	// BlockConnectedNtfnMethod is the method used for notifications from
	// the chain server that a block has been connected.
	BlockConnectedNtfnMethod Method = "blockconnected"

	// BlockDisconnectedNtfnMethod is the method used for notifications from
	// the chain server that a block has been disconnected.
	BlockDisconnectedNtfnMethod Method = "blockdisconnected"

	// NewTicketsNtfnMethod is the method of the daemon newtickets notification.
	NewTicketsNtfnMethod Method = "newtickets"

	// WorkNtfnMethod is the method used for notifications from
	// the chain server that a new block template has been generated.
	WorkNtfnMethod Method = "work"

	// TSpendNtfnMethod is the method used for notifications from the chain
	// server that a new tspend has arrived in the mempool.
	TSpendNtfnMethod Method = "tspend"

	// ReorganizationNtfnMethod is the method used for notifications that the
	// block chain is in the process of a reorganization.
	ReorganizationNtfnMethod Method = "reorganization"

	// TxAcceptedNtfnMethod is the method used for notifications from the
	// chain server that a transaction has been accepted into the mempool.
	TxAcceptedNtfnMethod Method = "txaccepted"

	// TxAcceptedVerboseNtfnMethod is the method used for notifications from
	// the chain server that a transaction has been accepted into the
	// mempool.  This differs from TxAcceptedNtfnMethod in that it provides
	// more details in the notification.
	TxAcceptedVerboseNtfnMethod Method = "txacceptedverbose"

	// RelevantTxAcceptedNtfnMethod is the method used for notifications
	// from the chain server that inform a client that a relevant
	// transaction was accepted by the mempool.
	RelevantTxAcceptedNtfnMethod Method = "relevanttxaccepted"

	// WinningTicketsNtfnMethod is the method of the daemon winningtickets
	// notification.
	WinningTicketsNtfnMethod Method = "winningtickets"

	// MixMessageNtfnMethod is the method of the mixmessage notification.
	MixMessageNtfnMethod Method = "mixmessage"
)

// BlockConnectedNtfn defines the blockconnected JSON-RPC notification.
type BlockConnectedNtfn struct {
	Header        string   `json:"header"`
	SubscribedTxs []string `json:"subscribedtxs"`
}

// NewBlockConnectedNtfn returns a new instance which can be used to issue a
// blockconnected JSON-RPC notification.
func NewBlockConnectedNtfn(header string, subscribedTxs []string) *BlockConnectedNtfn {
	return &BlockConnectedNtfn{
		Header:        header,
		SubscribedTxs: subscribedTxs,
	}
}

// BlockDisconnectedNtfn defines the blockdisconnected JSON-RPC notification.
type BlockDisconnectedNtfn struct {
	Header string `json:"header"`
}

// NewBlockDisconnectedNtfn returns a new instance which can be used to issue a
// blockdisconnected JSON-RPC notification.
func NewBlockDisconnectedNtfn(header string) *BlockDisconnectedNtfn {
	return &BlockDisconnectedNtfn{
		Header: header,
	}
}

// NewTicketsNtfn is a type handling custom marshaling and
// unmarshaling of newtickets JSON websocket notifications.
type NewTicketsNtfn struct {
	Hash      string
	Height    int32
	StakeDiff int64
	Tickets   []string
}

// NewNewTicketsNtfn creates a new NewTicketsNtfn.
func NewNewTicketsNtfn(hash string, height int32, stakeDiff int64, tickets []string) *NewTicketsNtfn {
	return &NewTicketsNtfn{
		Hash:      hash,
		Height:    height,
		StakeDiff: stakeDiff,
		Tickets:   tickets,
	}
}

// WorkNtfn defines the work JSON-RPC notification.
type WorkNtfn struct {
	Data   string `json:"data"`
	Target string `json:"target"`
	Reason string `json:"reason"`
}

// NewWorkNtfn returns a new instance which can be used to issue a
// work JSON-RPC notification.
func NewWorkNtfn(data string, target string) *WorkNtfn {
	return &WorkNtfn{
		Data:   data,
		Target: target,
	}
}

// TSpendNtfn defines the tspend JSON-RPC notification.
type TSpendNtfn struct {
	TSpend string `json:"tspend"` // Hex string encoded tspend.
}

// NewTSpendNtfn returns a new instance which can be used to issue a tspend
// JSON-RPC notification.
func NewTSpendNtfn(tspend string) *TSpendNtfn {
	return &TSpendNtfn{
		TSpend: tspend,
	}
}

// ReorganizationNtfn defines the reorganization JSON-RPC notification.
type ReorganizationNtfn struct {
	OldHash   string `json:"oldhash"`
	OldHeight int32  `json:"oldheight"`
	NewHash   string `json:"newhash"`
	NewHeight int32  `json:"newheight"`
}

// NewReorganizationNtfn returns a new instance which can be used to issue a
// blockdisconnected JSON-RPC notification.
func NewReorganizationNtfn(oldHash string, oldHeight int32, newHash string,
	newHeight int32) *ReorganizationNtfn {
	return &ReorganizationNtfn{
		OldHash:   oldHash,
		OldHeight: oldHeight,
		NewHash:   newHash,
		NewHeight: newHeight,
	}
}

// TxAcceptedNtfn defines the txaccepted JSON-RPC notification.
type TxAcceptedNtfn struct {
	TxID   string  `json:"txid"`
	Amount float64 `json:"amount"`
}

// NewTxAcceptedNtfn returns a new instance which can be used to issue a
// txaccepted JSON-RPC notification.
func NewTxAcceptedNtfn(txHash string, amount float64) *TxAcceptedNtfn {
	return &TxAcceptedNtfn{
		TxID:   txHash,
		Amount: amount,
	}
}

// TxAcceptedVerboseNtfn defines the txacceptedverbose JSON-RPC notification.
type TxAcceptedVerboseNtfn struct {
	RawTx TxRawResult `json:"rawtx"`
}

// NewTxAcceptedVerboseNtfn returns a new instance which can be used to issue a
// txacceptedverbose JSON-RPC notification.
func NewTxAcceptedVerboseNtfn(rawTx TxRawResult) *TxAcceptedVerboseNtfn {
	return &TxAcceptedVerboseNtfn{
		RawTx: rawTx,
	}
}

// RelevantTxAcceptedNtfn defines the parameters to the relevanttxaccepted
// JSON-RPC notification.
type RelevantTxAcceptedNtfn struct {
	Transaction string `json:"transaction"`
}

// NewRelevantTxAcceptedNtfn returns a new instance which can be used to issue a
// relevantxaccepted JSON-RPC notification.
func NewRelevantTxAcceptedNtfn(txHex string) *RelevantTxAcceptedNtfn {
	return &RelevantTxAcceptedNtfn{Transaction: txHex}
}

// WinningTicketsNtfn is a type handling custom marshaling and
// unmarshaling of blockconnected JSON websocket notifications.
type WinningTicketsNtfn struct {
	BlockHash   string
	BlockHeight int32
	Tickets     map[string]string
}

// NewWinningTicketsNtfn creates a new WinningTicketsNtfn.
func NewWinningTicketsNtfn(hash string, height int32, tickets map[string]string) *WinningTicketsNtfn {
	return &WinningTicketsNtfn{
		BlockHash:   hash,
		BlockHeight: height,
		Tickets:     tickets,
	}
}

// MixMessageNtfn defines the mixmessage JSON-RPC notification.
type MixMessageNtfn struct {
	Command string `json:"command"`
	Payload string `json:"payload"`
}

// NewMixMessageNtfn returns a new instance which can be used to issue a
// mixmessage JSON-RPC notification.
func NewMixMessageNtfn(command, payload string) *MixMessageNtfn {
	return &MixMessageNtfn{
		Command: command,
		Payload: payload,
	}
}
func init() {
	// The commands in this file are only usable by websockets and are
	// notifications.
	flags := VGLjson.UFWebsocketOnly | VGLjson.UFNotification

	VGLjson.MustRegister(BlockConnectedNtfnMethod, (*BlockConnectedNtfn)(nil), flags)
	VGLjson.MustRegister(BlockDisconnectedNtfnMethod, (*BlockDisconnectedNtfn)(nil), flags)
	VGLjson.MustRegister(WorkNtfnMethod, (*WorkNtfn)(nil), flags)
	VGLjson.MustRegister(TSpendNtfnMethod, (*TSpendNtfn)(nil), flags)
	VGLjson.MustRegister(NewTicketsNtfnMethod, (*NewTicketsNtfn)(nil), flags)
	VGLjson.MustRegister(ReorganizationNtfnMethod, (*ReorganizationNtfn)(nil), flags)
	VGLjson.MustRegister(TxAcceptedNtfnMethod, (*TxAcceptedNtfn)(nil), flags)
	VGLjson.MustRegister(TxAcceptedVerboseNtfnMethod, (*TxAcceptedVerboseNtfn)(nil), flags)
	VGLjson.MustRegister(RelevantTxAcceptedNtfnMethod, (*RelevantTxAcceptedNtfn)(nil), flags)
	VGLjson.MustRegister(WinningTicketsNtfnMethod, (*WinningTicketsNtfn)(nil), flags)
	VGLjson.MustRegister(MixMessageNtfnMethod, (*MixMessageNtfn)(nil), flags)
}
