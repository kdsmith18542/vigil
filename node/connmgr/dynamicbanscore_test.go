// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package connmgr

import (
	"fmt"
	"math"
	"testing"
	"time"
)

// TestDynamicBanScoreDecay tests the exponential decay implemented in
// DynamicBanScore.
func TestDynamicBanScoreDecay(t *testing.T) {
	var bs DynamicBanScore
	base := time.Now()

	r := bs.increase(100, 50, base)
	if r != 150 {
		t.Errorf("Unexpected result %d after ban score increase.", r)
	}

	r = bs.int(base.Add(time.Minute))
	if r != 125 {
		t.Errorf("Halflife check failed - %d instead of 125", r)
	}

	r = bs.int(base.Add(7 * time.Minute))
	if r != 100 {
		t.Errorf("Decay after 7m - %d instead of 100", r)
	}
}

// TestDynamicBanScoreLifetime tests that DynamicBanScore properly yields zero
// once the maximum age is reached.
func TestDynamicBanScoreLifetime(t *testing.T) {
	var bs DynamicBanScore
	base := time.Now()

	_ = bs.increase(0, math.MaxUint32, base)
	r := bs.int(base.Add(Lifetime * time.Second))
	if r != 3 { // 3, not 4 due to precision loss and truncating 3.999...
		t.Errorf("Pre max age check with MaxUint32 failed - %d", r)
	}
	r = bs.int(base.Add((Lifetime + 1) * time.Second))
	if r != 0 {
		t.Errorf("Zero after max age check failed - %d instead of 0", r)
	}
}

// TestDynamicBanScore tests exported functions of DynamicBanScore. Exponential
// decay or other time based behavior is tested by other functions.
func TestDynamicBanScoreReset(t *testing.T) {
	var bs DynamicBanScore
	if bs.Int() != 0 {
		t.Errorf("Initial state is not zero.")
	}
	bs.Increase(100, 0)
	r := bs.Int()
	if r != 100 {
		t.Errorf("Unexpected result %d after ban score increase.", r)
	}
	bs.Reset()
	if bs.Int() != 0 {
		t.Errorf("Failed to reset ban score.")
	}
}

// TestDynamicBanScoreStringer ensure the stringer for a dynamic ban score
// produces the expected result.
func TestDynamicBanScoreStringer(t *testing.T) {
	var bs DynamicBanScore
	persistent := uint32(100)
	transient := uint32(50)
	now := time.Now()
	score := bs.increase(persistent, transient, now)
	wantScore := persistent + transient
	if score != wantScore {
		t.Fatalf("Unexpected ban score increase -- got %v, want %v", score,
			wantScore)
	}

	result := bs.String()
	wantStr := fmt.Sprintf("persistent %v + transient %v at %v = %v as of now",
		persistent, transient, now.Unix(), wantScore)
	if result != wantStr {
		t.Fatalf("Unexpected string result -- got %v, want %v", result, wantStr)
	}
}
