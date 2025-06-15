package txrules_test

import (
	"testing"

	. "github.com/kdsmith18542/vigil/wallet/wallet/txrules"
	"github.com/kdsmith18542/vigil/chaincfg/v3"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
)

func TestStakePoolTicketFee(t *testing.T) {
	params := chaincfg.MainNetParams()
	tests := []struct {
		StakeDiff       VGLutil.Amount
		Fee             VGLutil.Amount
		Height          int32
		PoolFee         float64
		Expected        VGLutil.Amount
		IsVGLP0010Active bool
		IsVGLP0012Active bool
	}{
		0: {
			StakeDiff:       10 * 1e8,
			Fee:             0.01 * 1e8,
			Height:          25000,
			PoolFee:         1.00,
			Expected:        0.01500463 * 1e8,
			IsVGLP0010Active: false,
			IsVGLP0012Active: false,
		},
		1: {
			StakeDiff:       20 * 1e8,
			Fee:             0.01 * 1e8,
			Height:          25000,
			PoolFee:         1.00,
			Expected:        0.01621221 * 1e8,
			IsVGLP0010Active: false,
			IsVGLP0012Active: false,
		},
		2: {
			StakeDiff:       5 * 1e8,
			Fee:             0.05 * 1e8,
			Height:          50000,
			PoolFee:         2.59,
			Expected:        0.03310616 * 1e8,
			IsVGLP0010Active: false,
			IsVGLP0012Active: false,
		},
		3: {
			StakeDiff:       15 * 1e8,
			Fee:             0.05 * 1e8,
			Height:          50000,
			PoolFee:         2.59,
			Expected:        0.03956376 * 1e8,
			IsVGLP0010Active: false,
			IsVGLP0012Active: false,
		},
		4: {
			StakeDiff:       15 * 1e8,
			Fee:             0.05 * 1e8,
			Height:          50000,
			PoolFee:         2.59,
			Expected:        0.09023823 * 1e8,
			IsVGLP0010Active: true,
			IsVGLP0012Active: false,
		},
		5: {
			StakeDiff:       15 * 1e8,
			Fee:             0.05 * 1e8,
			Height:          50000,
			PoolFee:         2.59,
			Expected:        0.09784185 * 1e8,
			IsVGLP0010Active: false,
			IsVGLP0012Active: true,
		},
		6: {
			StakeDiff:       15 * 1e8,
			Fee:             0.05 * 1e8,
			Height:          50000,
			PoolFee:         2.59,
			Expected:        0.09784185 * 1e8,
			IsVGLP0010Active: true,
			IsVGLP0012Active: true,
		},
	}
	for i, test := range tests {
		poolFeeAmt := StakePoolTicketFee(test.StakeDiff, test.Fee, test.Height,
			test.PoolFee, params, test.IsVGLP0010Active, test.IsVGLP0012Active)
		if poolFeeAmt != test.Expected {
			t.Errorf("Test %d: Got %v: Want %v", i, poolFeeAmt, test.Expected)
		}
	}
}
