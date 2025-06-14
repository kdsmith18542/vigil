// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ticketbuyer

import (
	"context"
	"runtime/trace"
	"sync"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/wallet"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
	"github.com/kdsmith18542/vigil/wire"
)

const minconf = 1

// Config modifies the behavior of TB.
type Config struct {
	BuyTickets bool

	// Account to buy tickets from
	Account uint32

	// Account to derive voting addresses from
	VotingAccount uint32

	// Minimum amount to maintain in purchasing account
	Maintain VGLutil.Amount

	// Limit maximum number of purchased tickets per block
	Limit int

	// CSPP-related options
	Mixing             bool
	MixedAccount       uint32
	MixedAccountBranch uint32
	TicketSplitAccount uint32
	ChangeAccount      uint32
	MixChange          bool

	// VSP client
	VSP *wallet.VSPClient
}

// TB is an automated ticket buyer, buying as many tickets as possible given an
// account's available balance. TB may optionally be configured to register
// purchased tickets with a VSP.
type TB struct {
	wallet *wallet.Wallet

	cfg Config
	mu  sync.Mutex
}

// New returns a new TB to buy tickets from a wallet.
func New(w *wallet.Wallet, cfg Config) *TB {
	return &TB{wallet: w, cfg: cfg}
}

// Run executes the ticket buyer.  If the private passphrase is incorrect, or
// ever becomes incorrect due to a wallet passphrase change, Run exits with an
// errors.Passphrase error.
func (tb *TB) Run(ctx context.Context, passphrase []byte) error {
	if len(passphrase) > 0 {
		err := tb.wallet.Unlock(ctx, passphrase, nil)
		if err != nil {
			return err
		}
	}

	c := tb.wallet.NtfnServer.MainTipChangedNotifications()
	defer c.Done()

	ctx, outerCancel := context.WithCancel(ctx)
	defer outerCancel()
	var fatal error
	var fatalMu sync.Mutex

	var nextIntervalStart, expiry int32
	var cancels []func()
	for {
		select {
		case <-ctx.Done():
			defer outerCancel()
			fatalMu.Lock()
			err := fatal
			fatalMu.Unlock()
			if err != nil {
				return err
			}
			return ctx.Err()
		case n := <-c.C:
			if len(n.AttachedBlocks) == 0 {
				continue
			}

			tip := n.AttachedBlocks[len(n.AttachedBlocks)-1]
			w := tb.wallet

			// Don't perform any actions while transactions are not synced through
			// the tip block.
			rp, err := w.RescanPoint(ctx)
			if err != nil {
				log.Debugf("Skipping autobuyer actions: RescanPoint err: %v", err)
				continue
			}
			if rp != nil {
				log.Debugf("Skipping autobuyer actions: transactions are not synced")
				continue
			}

			tipHeader, err := w.BlockHeader(ctx, tip)
			if err != nil {
				log.Error(err)
				continue
			}
			height := int32(tipHeader.Height)

			// Cancel any ongoing ticket purchases which are buying
			// at an old ticket price or are no longer able to
			// create mined tickets the window.
			if height+2 >= nextIntervalStart {
				for i, cancel := range cancels {
					cancel()
					cancels[i] = nil
				}
				cancels = cancels[:0]

				intervalSize := int32(w.ChainParams().StakeDiffWindowSize)
				currentInterval := height / intervalSize
				nextIntervalStart = (currentInterval + 1) * intervalSize

				// Skip this purchase when no more tickets may be purchased in the interval and
				// the next sdiff is unknown.  The earliest any ticket may be mined is two
				// blocks from now, with the next block containing the split transaction
				// that the ticket purchase spends.
				if height+2 == nextIntervalStart {
					log.Debugf("Skipping purchase: next sdiff interval starts soon")
					continue
				}
				// Set expiry to prevent tickets from being mined in the next
				// sdiff interval.  When the next block begins the new interval,
				// the ticket is being purchased for the next interval; therefore
				// increment expiry by a full sdiff window size to prevent it
				// being mined in the interval after the next.
				expiry = nextIntervalStart
				if height+1 == nextIntervalStart {
					expiry += intervalSize
				}
			}

			// Read config
			tb.mu.Lock()
			cfg := tb.cfg
			tb.mu.Unlock()

			multiple := 1
			if cfg.Mixing {
				multiple = cfg.Limit
				cfg.Limit = 1
			}

			cancelCtx, cancel := context.WithCancel(ctx)
			cancels = append(cancels, cancel)
			buyTickets := func() {
				err := tb.buy(cancelCtx, passphrase, tipHeader, expiry, &cfg)
				if err != nil {
					switch {
					// silence these errors
					case errors.Is(err, errors.InsufficientBalance):
					case errors.Is(err, context.Canceled):
					case errors.Is(err, context.DeadlineExceeded):
					default:
						log.Errorf("Ticket purchasing failed: %v", err)
					}
					if errors.Is(err, errors.Passphrase) {
						fatalMu.Lock()
						fatal = err
						fatalMu.Unlock()
						outerCancel()
					}
				}
			}
			for i := 0; cfg.BuyTickets && i < multiple; i++ {
				go buyTickets()
			}
			go func() {
				err := tb.mixChange(ctx, &cfg)
				if err != nil {
					log.Error(err)
				}
			}()
		}
	}
}

func (tb *TB) buy(ctx context.Context, passphrase []byte, tip *wire.BlockHeader, expiry int32,
	cfg *Config) error {
	ctx, task := trace.NewTask(ctx, "ticketbuyer.buy")
	defer task.End()

	tb.mu.Lock()
	buyTickets := tb.cfg.BuyTickets
	tb.mu.Unlock()
	if !buyTickets {
		return nil
	}

	w := tb.wallet

	// Unable to publish any transactions if the network backend is unset.
	n, err := w.NetworkBackend()
	if err != nil {
		return err
	}
	ctx, cancel := wallet.WrapNetworkBackendContext(n, ctx)
	defer cancel()

	if len(passphrase) > 0 {
		// Ensure wallet is unlocked with the current passphrase.  If the passphase
		// is changed, the Run exits and TB must be restarted with the new
		// passphrase.
		err = w.Unlock(ctx, passphrase, nil)
		if err != nil {
			return err
		}
	}

	// Read config
	account := cfg.Account
	maintain := cfg.Maintain
	limit := cfg.Limit
	mixing := cfg.Mixing
	votingAccount := cfg.VotingAccount
	mixedAccount := cfg.MixedAccount
	mixedBranch := cfg.MixedAccountBranch
	splitAccount := cfg.TicketSplitAccount
	changeAccount := cfg.ChangeAccount

	minconf := int32(minconf)
	if mixing {
		minconf = 2
	}

	sdiff, err := w.NextStakeDifficultyAfterHeader(ctx, tip)
	if err != nil {
		return err
	}

	// Determine how many tickets to buy
	var buy int
	if maintain != 0 {
		bal, err := w.AccountBalance(ctx, account, minconf)
		if err != nil {
			return err
		}
		spendable := bal.Spendable
		if spendable < maintain {
			log.Debugf("Skipping purchase: low available balance")
			return nil
		}
		spendable -= maintain
		buy = int(spendable / sdiff)
		if buy == 0 {
			log.Debugf("Skipping purchase: low available balance")
			return nil
		}
		max := int(w.ChainParams().MaxFreshStakePerBlock)
		if buy > max {
			buy = max
		}
	} else {
		buy = int(w.ChainParams().MaxFreshStakePerBlock)
	}
	if limit == 0 && mixing {
		buy = 1
	} else if limit > 0 && buy > limit {
		buy = limit
	}

	purchaseTicketReq := &wallet.PurchaseTicketsRequest{
		Count:         buy,
		SourceAccount: account,
		VotingAccount: votingAccount,
		MinConf:       minconf,
		Expiry:        expiry,

		// CSPP
		Mixing:             mixing,
		MixedAccount:       mixedAccount,
		MixedAccountBranch: mixedBranch,
		MixedSplitAccount:  splitAccount,
		ChangeAccount:      changeAccount,

		VSPClient: tb.cfg.VSP,
	}

	tix, err := w.PurchaseTickets(ctx, n, purchaseTicketReq)
	if tix != nil {
		for _, hash := range tix.TicketHashes {
			log.Infof("Purchased ticket %v at stake difficulty %v", hash, sdiff)
		}
	}
	return err
}

// AccessConfig runs f with the current config passed as a parameter.  The
// config is protected by a mutex and this function is safe for concurrent
// access to read or modify the config.  It is unsafe to leak a pointer to the
// config, but a copy of *cfg is legal.
func (tb *TB) AccessConfig(f func(cfg *Config)) {
	tb.mu.Lock()
	f(&tb.cfg)
	tb.mu.Unlock()
}

func (tb *TB) mixChange(ctx context.Context, cfg *Config) error {
	// Read config
	mixing := cfg.Mixing
	mixedAccount := cfg.MixedAccount
	mixedBranch := cfg.MixedAccountBranch
	changeAccount := cfg.ChangeAccount
	mixChange := cfg.MixChange

	if !mixChange || !mixing {
		return nil
	}

	ctx, task := trace.NewTask(ctx, "ticketbuyer.mixChange")
	defer task.End()

	return tb.wallet.MixAccount(ctx, changeAccount, mixedAccount, mixedBranch)
}
