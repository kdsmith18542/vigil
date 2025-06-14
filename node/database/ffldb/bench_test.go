// Copyright (c) 2015-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ffldb

import (
	"testing"

	"github.com/kdsmith18542/vigil/chaincfg/v3"
	"github.com/kdsmith18542/vigil/database/v3"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
)

// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis
// block header.
func BenchmarkBlockHeader(b *testing.B) {
	// Start by creating a new database and populating it with the mainnet
	// genesis block.
	dbPath := b.TempDir()
	db, err := database.Create("ffldb", dbPath, blockDataNet)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		db.Close()
	})
	mainNetParams := chaincfg.MainNetParams()
	err = db.Update(func(tx database.Tx) error {
		block := VGLutil.NewBlock(mainNetParams.GenesisBlock)
		return tx.StoreBlock(block)
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	err = db.View(func(tx database.Tx) error {
		blockHash := &mainNetParams.GenesisHash
		for i := 0; i < b.N; i++ {
			_, err := tx.FetchBlockHeader(blockHash)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}

	// Don't benchmark teardown.
	b.StopTimer()
}

// BenchmarkBlockHeader benchmarks how long it takes to load the mainnet genesis
// block.
func BenchmarkBlock(b *testing.B) {
	// Start by creating a new database and populating it with the mainnet
	// genesis block.
	dbPath := b.TempDir()
	db, err := database.Create("ffldb", dbPath, blockDataNet)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		db.Close()
	})
	mainNetParams := chaincfg.MainNetParams()
	err = db.Update(func(tx database.Tx) error {
		block := VGLutil.NewBlock(mainNetParams.GenesisBlock)
		return tx.StoreBlock(block)
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	err = db.View(func(tx database.Tx) error {
		blockHash := &mainNetParams.GenesisHash
		for i := 0; i < b.N; i++ {
			_, err := tx.FetchBlock(blockHash)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}

	// Don't benchmark teardown.
	b.StopTimer()
}
