// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package netsync

import (
	"github.com/kdsmith18542/vigil/VGLutil/v4"
	"github.com/kdsmith18542/vigil/mixing"
)

// PeerNotifier provides an interface to notify peers of status changes related
// to blocks and transactions.
type PeerNotifier interface {
	// AnnounceNewTransactions generates and relays inventory vectors and
	// notifies websocket clients of the passed transactions.
	AnnounceNewTransactions(txns []*VGLutil.Tx)

	// AnnounceMixMessages generates and relays inventory vectors of the
	// passed messages.
	AnnounceMixMessages(msgs []mixing.Message)
}
