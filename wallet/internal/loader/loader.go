// Copyright (c) 2015-2018 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package loader

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/wallet"
	_ "github.com/kdsmith18542/vigil/wallet/wallet/drivers/bdb" // driver loaded during init
	"github.com/kdsmith18542/vigil/chaincfg/v3"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
)

const (
	walletDbName = "wallet.db"
	driver       = "bdb"
)

// Loader implements the creating of new and opening of existing wallets, while
// providing a callback system for other subsystems to handle the loading of a
// wallet.  This is primarely intended for use by the RPC servers, to enable
// methods and services which require the wallet when the wallet is loaded by
// another subsystem.
//
// Loader is safe for concurrent access.
type Loader struct {
	callbacks   []func(*wallet.Wallet)
	chainParams *chaincfg.Params
	dbDirPath   string
	wallet      *wallet.Wallet
	db          wallet.DB

	votingEnabled           bool
	gapLimit                uint32
	watchLast               uint32
	accountGapLimit         int
	disableCoinTypeUpgrades bool
	mixingEnabled           bool
	allowHighFees           bool
	manualTickets           bool
	relayFee                VGLutil.Amount
	vspMaxFee               VGLutil.Amount
	mixSplitLimit           int
	dialer                  wallet.DialFunc

	mu sync.Mutex
}

// NewLoader constructs a Loader.
func NewLoader(chainParams *chaincfg.Params, dbDirPath string, votingEnabled bool, gapLimit uint32,
	watchLast uint32, allowHighFees bool, relayFee VGLutil.Amount, vspMaxFee VGLutil.Amount, accountGapLimit int,
	disableCoinTypeUpgrades bool, mixingEnabled bool, manualTickets bool, mixSplitLimit int, dialer wallet.DialFunc) *Loader {

	return &Loader{
		chainParams:             chainParams,
		dbDirPath:               dbDirPath,
		votingEnabled:           votingEnabled,
		gapLimit:                gapLimit,
		watchLast:               watchLast,
		accountGapLimit:         accountGapLimit,
		disableCoinTypeUpgrades: disableCoinTypeUpgrades,
		mixingEnabled:           mixingEnabled,
		allowHighFees:           allowHighFees,
		manualTickets:           manualTickets,
		relayFee:                relayFee,
		vspMaxFee:               vspMaxFee,
		mixSplitLimit:           mixSplitLimit,
		dialer:                  dialer,
	}
}

// onLoaded executes each added callback and prevents loader from loading any
// additional wallets.  Requires mutex to be locked.
func (l *Loader) onLoaded(w *wallet.Wallet, db wallet.DB) {
	for _, fn := range l.callbacks {
		fn(w)
	}

	l.wallet = w
	l.db = db
	l.callbacks = nil // not needed anymore
}

// RunAfterLoad adds a function to be executed when the loader creates or opens
// a wallet.  Functions are executed in a single goroutine in the order they are
// added.
func (l *Loader) RunAfterLoad(fn func(*wallet.Wallet)) {
	l.mu.Lock()
	if l.wallet != nil {
		w := l.wallet
		l.mu.Unlock()
		fn(w)
	} else {
		l.callbacks = append(l.callbacks, fn)
		l.mu.Unlock()
	}
}

// CreateWatchingOnlyWallet creates a new watch-only wallet using the provided
// extended public key and public passphrase.
func (l *Loader) CreateWatchingOnlyWallet(ctx context.Context, extendedPubKey string, pubPass []byte) (w *wallet.Wallet, err error) {
	const op errors.Op = "loader.CreateWatchingOnlyWallet"

	defer l.mu.Unlock()
	l.mu.Lock()

	if l.wallet != nil {
		return nil, errors.E(op, errors.Exist, "wallet already loaded")
	}

	// Ensure that the network directory exists.
	if fi, err := os.Stat(l.dbDirPath); err != nil {
		if os.IsNotExist(err) {
			// Attempt data directory creation
			if err = os.MkdirAll(l.dbDirPath, 0700); err != nil {
				return nil, errors.E(op, err)
			}
		} else {
			return nil, errors.E(op, err)
		}
	} else {
		if !fi.IsDir() {
			return nil, errors.E(op, errors.Invalid, errors.Errorf("%q is not a directory", l.dbDirPath))
		}
	}

	dbPath := filepath.Join(l.dbDirPath, walletDbName)
	exists, err := fileExists(dbPath)
	if err != nil {
		return nil, errors.E(op, err)
	}
	if exists {
		return nil, errors.E(op, errors.Exist, "wallet already exists")
	}

	// At this point it is asserted that there is no existing database file, and
	// deleting anything won't destroy a wallet in use.  Defer a function that
	// attempts to remove any written database file if this function errors.
	defer func() {
		if err != nil {
			_ = os.Remove(dbPath)
		}
	}()

	// Create the wallet database backed by bolt db.
	err = os.MkdirAll(l.dbDirPath, 0700)
	if err != nil {
		return nil, errors.E(op, err)
	}
	db, err := wallet.CreateDB(driver, dbPath)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Initialize the watch-only database for the wallet before opening.
	err = wallet.CreateWatchOnly(ctx, db, extendedPubKey, pubPass, l.chainParams)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Open the watch-only wallet.
	cfg := &wallet.Config{
		DB:                      db,
		PubPassphrase:           pubPass,
		VotingEnabled:           l.votingEnabled,
		GapLimit:                l.gapLimit,
		WatchLast:               l.watchLast,
		AccountGapLimit:         l.accountGapLimit,
		DisableCoinTypeUpgrades: l.disableCoinTypeUpgrades,
		MixingEnabled:           l.mixingEnabled,
		ManualTickets:           l.manualTickets,
		AllowHighFees:           l.allowHighFees,
		RelayFee:                l.relayFee,
		VSPMaxFee:               l.vspMaxFee,
		MixSplitLimit:           l.mixSplitLimit,
		Params:                  l.chainParams,
		Dialer:                  l.dialer,
	}
	w, err = wallet.Open(ctx, cfg)
	if err != nil {
		return nil, errors.E(op, err)
	}

	l.onLoaded(w, db)
	return w, nil
}

// CreateNewWallet creates a new wallet using the provided public and private
// passphrases.  The seed is optional.  If non-nil, addresses are derived from
// this seed.  If nil, a secure random seed is generated.
func (l *Loader) CreateNewWallet(ctx context.Context, pubPassphrase, privPassphrase, seed []byte) (w *wallet.Wallet, err error) {
	const op errors.Op = "loader.CreateNewWallet"

	defer l.mu.Unlock()
	l.mu.Lock()

	if l.wallet != nil {
		return nil, errors.E(op, errors.Exist, "wallet already opened")
	}

	// Ensure that the network directory exists.
	if fi, err := os.Stat(l.dbDirPath); err != nil {
		if os.IsNotExist(err) {
			// Attempt data directory creation
			if err = os.MkdirAll(l.dbDirPath, 0700); err != nil {
				return nil, errors.E(op, err)
			}
		} else {
			return nil, errors.E(op, err)
		}
	} else {
		if !fi.IsDir() {
			return nil, errors.E(op, errors.Errorf("%q is not a directory", l.dbDirPath))
		}
	}

	dbPath := filepath.Join(l.dbDirPath, walletDbName)
	exists, err := fileExists(dbPath)
	if err != nil {
		return nil, errors.E(op, err)
	}
	if exists {
		return nil, errors.E(op, errors.Exist, "wallet DB exists")
	}

	// At this point it is asserted that there is no existing database file, and
	// deleting anything won't destroy a wallet in use.  Defer a function that
	// attempts to remove any written database file if this function errors.
	defer func() {
		if err != nil {
			_ = os.Remove(dbPath)
		}
	}()

	// Create the wallet database backed by bolt db.
	err = os.MkdirAll(l.dbDirPath, 0700)
	if err != nil {
		return nil, errors.E(op, err)
	}
	db, err := wallet.CreateDB(driver, dbPath)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Initialize the newly created database for the wallet before opening.
	err = wallet.Create(ctx, db, pubPassphrase, privPassphrase, seed, l.chainParams)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// Open the newly-created wallet.
	cfg := &wallet.Config{
		DB:                      db,
		PubPassphrase:           pubPassphrase,
		VotingEnabled:           l.votingEnabled,
		GapLimit:                l.gapLimit,
		WatchLast:               l.watchLast,
		AccountGapLimit:         l.accountGapLimit,
		DisableCoinTypeUpgrades: l.disableCoinTypeUpgrades,
		MixingEnabled:           l.mixingEnabled,
		ManualTickets:           l.manualTickets,
		AllowHighFees:           l.allowHighFees,
		RelayFee:                l.relayFee,
		VSPMaxFee:               l.vspMaxFee,
		Params:                  l.chainParams,
		Dialer:                  l.dialer,
	}
	w, err = wallet.Open(ctx, cfg)
	if err != nil {
		return nil, errors.E(op, err)
	}

	l.onLoaded(w, db)
	return w, nil
}

// OpenExistingWallet opens the wallet from the loader's wallet database path
// and the public passphrase.  If the loader is being called by a context where
// standard input prompts may be used during wallet upgrades, setting
// canConsolePrompt will enable these prompts.
func (l *Loader) OpenExistingWallet(ctx context.Context, pubPassphrase []byte) (w *wallet.Wallet, rerr error) {
	const op errors.Op = "loader.OpenExistingWallet"

	defer l.mu.Unlock()
	l.mu.Lock()

	if l.wallet != nil {
		return nil, errors.E(op, errors.Exist, "wallet already opened")
	}

	// Open the database using the boltdb backend.
	dbPath := filepath.Join(l.dbDirPath, walletDbName)
	l.mu.Unlock()
	db, err := wallet.OpenDB(driver, dbPath)
	l.mu.Lock()

	if err != nil {
		log.Errorf("Failed to open database: %v", err)
		return nil, errors.E(op, err)
	}
	// If this function does not return to completion the database must be
	// closed.  Otherwise, because the database is locked on opens, any
	// other attempts to open the wallet will hang, and there is no way to
	// recover since this db handle would be leaked.
	defer func() {
		if rerr != nil {
			db.Close()
		}
	}()

	cfg := &wallet.Config{
		DB:                      db,
		PubPassphrase:           pubPassphrase,
		VotingEnabled:           l.votingEnabled,
		GapLimit:                l.gapLimit,
		WatchLast:               l.watchLast,
		AccountGapLimit:         l.accountGapLimit,
		DisableCoinTypeUpgrades: l.disableCoinTypeUpgrades,
		MixingEnabled:           l.mixingEnabled,
		ManualTickets:           l.manualTickets,
		AllowHighFees:           l.allowHighFees,
		RelayFee:                l.relayFee,
		VSPMaxFee:               l.vspMaxFee,
		MixSplitLimit:           l.mixSplitLimit,
		Params:                  l.chainParams,
		Dialer:                  l.dialer,
	}
	w, err = wallet.Open(ctx, cfg)
	if err != nil {
		return nil, errors.E(op, err)
	}

	l.onLoaded(w, db)
	return w, nil
}

// DbDirPath returns the Loader's database directory path
func (l *Loader) DbDirPath() string {
	return l.dbDirPath
}

// WalletExists returns whether a file exists at the loader's database path.
// This may return an error for unexpected I/O failures.
func (l *Loader) WalletExists() (bool, error) {
	const op errors.Op = "loader.WalletExists"
	dbPath := filepath.Join(l.dbDirPath, walletDbName)
	exists, err := fileExists(dbPath)
	if err != nil {
		return false, errors.E(op, err)
	}
	return exists, nil
}

// LoadedWallet returns the loaded wallet, if any, and a bool for whether the
// wallet has been loaded or not.  If true, the wallet pointer should be safe to
// dereference.
func (l *Loader) LoadedWallet() (*wallet.Wallet, bool) {
	l.mu.Lock()
	w := l.wallet
	l.mu.Unlock()
	return w, w != nil
}

// UnloadWallet stops the loaded wallet, if any, and closes the wallet database.
// Returns with errors.Invalid if the wallet has not been loaded with
// CreateNewWallet or LoadExistingWallet.  The Loader may be reused if this
// function returns without error.
func (l *Loader) UnloadWallet() error {
	const op errors.Op = "loader.UnloadWallet"

	defer l.mu.Unlock()
	l.mu.Lock()

	if l.wallet == nil {
		return errors.E(op, errors.Invalid, "wallet is unopened")
	}

	err := l.db.Close()
	if err != nil {
		return errors.E(op, err)
	}

	l.wallet = nil
	l.db = nil
	return nil
}

// NetworkBackend returns the associated wallet network backend, if any, and a
// bool describing whether a non-nil network backend was set.
func (l *Loader) NetworkBackend() (n wallet.NetworkBackend, ok bool) {
	l.mu.Lock()
	if l.wallet != nil {
		n, _ = l.wallet.NetworkBackend()
	}
	l.mu.Unlock()
	return n, n != nil
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
