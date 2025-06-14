// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wallet

import (
	"context"
	"os"
	"testing"

	_ "github.com/vigilnetwork/vgl/wallet/wallet/drivers/bdb"
	"github.com/vigilnetwork/vgl/wallet/wallet/walletdb"
	"github.com/vigilnetwork/vgl/chaincfg/v3"
	"github.com/vigilnetwork/vgl/VGLutil/v4"
)

var testPrivPass = []byte("private")

var basicWalletConfig = Config{
	PubPassphrase: []byte(InsecurePubPassphrase),
	GapLimit:      20,
	RelayFee:      VGLutil.Amount(1e5),
	Params:        chaincfg.SimNetParams(),
	MixingEnabled: true,
}

func testWallet(ctx context.Context, t *testing.T, cfg *Config, seed []byte) (w *Wallet, teardown func()) {
	f, err := os.CreateTemp(t.TempDir(), "vglwallet.testdb")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	db, err := walletdb.Create("bdb", f.Name())
	if err != nil {
		t.Fatal(err)
	}
	rm := func() {
		db.Close()
		os.Remove(f.Name())
	}
	err = Create(ctx, opaqueDB{db}, []byte(InsecurePubPassphrase), testPrivPass, seed, cfg.Params)
	if err != nil {
		rm()
		t.Fatal(err)
	}
	cfg.DB = opaqueDB{db}
	w, err = Open(ctx, cfg)
	if err != nil {
		rm()
		t.Fatal(err)
	}
	teardown = func() {
		rm()
	}
	return
}
