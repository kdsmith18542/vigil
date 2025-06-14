// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/kdsmith18542/vigil/wallet/chain"
	"github.com/kdsmith18542/vigil/wallet/errors"
	ldr "github.com/kdsmith18542/vigil/wallet/internal/loader"
	"github.com/kdsmith18542/vigil/wallet/internal/loggers"
	"github.com/kdsmith18542/vigil/wallet/internal/prompt"
	"github.com/kdsmith18542/vigil/wallet/internal/rpc/rpcserver"
	"github.com/kdsmith18542/vigil/wallet/p2p"
	"github.com/kdsmith18542/vigil/wallet/spv"
	"github.com/kdsmith18542/vigil/wallet/ticketbuyer"
	"github.com/kdsmith18542/vigil/wallet/version"
	"github.com/kdsmith18542/vigil/wallet/wallet"
	"github.com/kdsmith18542/vigil/addrmgr/v2"
	"github.com/kdsmith18542/vigil/wire"
)

func init() {
	// Format nested errors without newlines (better for logs).
	errors.Separator = ":: "
}

var (
	cfg *config
)

func main() {
	// Create a context that is cancelled when a shutdown request is received
	// through an interrupt signal or an RPC request.
	ctx := withShutdownCancel(context.Background())
	go shutdownListener()

	// Run the wallet until permanent failure or shutdown is requested.
	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		os.Exit(1)
	}
}

// done returns whether the context's Done channel was closed due to
// cancellation or exceeded deadline.
func done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// run is the main startup and teardown logic performed by the main package.  It
// is responsible for parsing the config, starting RPC servers, loading and
// syncing the wallet (if necessary), and stopping all started services when the
// context is cancelled.
func run(ctx context.Context) error {
	// Load configuration and parse command line.  This function also
	// initializes logging and configures it accordingly.
	tcfg, _, err := loadConfig(ctx)
	if err != nil {
		return err
	}
	cfg = tcfg
	defer loggers.CloseLogRotator()

	// Show version at startup.
	log.Infof("Version %s (Go version %s %s/%s)", version.String(), runtime.Version(),
		runtime.GOOS, runtime.GOARCH)
	if cfg.NoFileLogging {
		log.Info("File logging disabled")
	}

	// Read IPC messages from the read end of a pipe created and passed by the
	// parent process, if any.  When this pipe is closed, shutdown is
	// initialized.
	if cfg.PipeRx != nil {
		go serviceControlPipeRx(uintptr(*cfg.PipeRx))
	}
	if cfg.PipeTx != nil {
		go serviceControlPipeTx(uintptr(*cfg.PipeTx))
	} else {
		go drainOutgoingPipeMessages()
	}

	// Run the pprof profiler if enabled.
	if len(cfg.Profile) > 0 {
		if done(ctx) {
			return ctx.Err()
		}

		profileRedirect := http.RedirectHandler("/debug/pprof", http.StatusSeeOther)
		http.Handle("/", profileRedirect)
		for _, listenAddr := range cfg.Profile {
			listenAddr := listenAddr // copy for closure
			go func() {
				log.Infof("Starting profile server on %s", listenAddr)
				err := http.ListenAndServe(listenAddr, nil)
				if err != nil {
					fatalf("Unable to run profiler: %v", err)
				}
			}()
		}
	}

	// Write cpu profile if requested.
	if cfg.CPUProfile != "" {
		if done(ctx) {
			return ctx.Err()
		}

		f, err := os.Create(cfg.CPUProfile)
		if err != nil {
			log.Errorf("Unable to create cpu profile: %v", err.Error())
			return err
		}
		pprof.StartCPUProfile(f)
		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	// Write mem profile if requested.
	if cfg.MemProfile != "" {
		if done(ctx) {
			return ctx.Err()
		}

		f, err := os.Create(cfg.MemProfile)
		if err != nil {
			log.Errorf("Unable to create mem profile: %v", err)
			return err
		}
		defer func() {
			pprof.WriteHeapProfile(f)
			f.Close()
		}()
	}

	if done(ctx) {
		return ctx.Err()
	}

	// Create the loader which is used to load and unload the wallet.  If
	// --noinitialload is not set, this function is responsible for loading the
	// wallet.  Otherwise, loading is deferred so it can be performed over RPC.
	dbDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)

	loader := ldr.NewLoader(activeNet.Params, dbDir, cfg.EnableVoting,
		cfg.GapLimit, cfg.WatchLast, cfg.AllowHighFees, cfg.RelayFee.Amount,
		cfg.VSPOpts.MaxFee.Amount, cfg.AccountGapLimit,
		cfg.DisableCoinTypeUpgrades, cfg.MixingEnabled, cfg.ManualTickets,
		cfg.MixSplitLimit, cfg.dial)

	// Stop any services started by the loader after the shutdown procedure is
	// initialized and this function returns.
	defer func() {
		// When panicing, do not cleanly unload the wallet (by closing
		// the db).  If a panic occurred inside a bolt transaction, the
		// db mutex is still held and this causes a deadlock.
		if r := recover(); r != nil {
			panic(r)
		}
		err := loader.UnloadWallet()
		if err != nil && !errors.Is(err, errors.Invalid) {
			log.Errorf("Failed to close wallet: %v", err)
		} else if err == nil {
			log.Infof("Closed wallet")
		}
	}()

	// Open the wallet when --noinitialload was not set.
	var vspClient *wallet.VSPClient
	passphrase := []byte{}
	if !cfg.NoInitialLoad {
		walletPass := []byte(cfg.WalletPass)
		if cfg.PromptPublicPass {
			walletPass, _ = passPrompt(ctx, "Enter public wallet passphrase", false)
		}

		if done(ctx) {
			return ctx.Err()
		}

		// Load the wallet.  It must have been created already or this will
		// return an appropriate error.
		var w *wallet.Wallet
		errc := make(chan error, 1)
		go func() {
			defer zero(walletPass)
			var err error
			w, err = loader.OpenExistingWallet(ctx, walletPass)
			if err != nil {
				log.Errorf("Failed to open wallet: %v", err)
				if errors.Is(err, errors.Passphrase) {
					// walletpass not provided, advice using --walletpass or --promptpublicpass
					if cfg.WalletPass == wallet.InsecurePubPassphrase {
						log.Info("Configure public passphrase with walletpass or promptpublicpass options.")
					}
				}
			}
			errc <- err
		}()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errc:
			if err != nil {
				return err
			}
		}

		// TODO(jrick): I think that this prompt should be removed
		// entirely instead of enabling it when --noinitialload is
		// unset.  It can be replaced with an RPC request (either
		// providing the private passphrase as a parameter, or require
		// unlocking the wallet first) to trigger a full accounts
		// rescan.
		//
		// Until then, since --noinitialload users are expecting to use
		// the wallet only over RPC, disable this feature for them.
		if cfg.Pass != "" {
			passphrase = []byte(cfg.Pass)
			err = w.Unlock(ctx, passphrase, nil)
			if err != nil {
				log.Errorf("Incorrect passphrase in pass config setting.")
				return err
			}
		} else {
			passphrase = startPromptPass(ctx, w)
		}

		if cfg.VSPOpts.URL != "" {
			changeAccountName := cfg.ChangeAccount
			if changeAccountName == "" && !cfg.MixingEnabled {
				log.Warnf("Change account not set, using "+
					"purchase account %q", cfg.PurchaseAccount)
				changeAccountName = cfg.PurchaseAccount
			}
			changeAcct, err := w.AccountNumber(ctx, changeAccountName)
			if err != nil {
				log.Warnf("failed to get account number for "+
					"ticket change account %q: %v",
					changeAccountName, err)
				return err
			}
			purchaseAcct, err := w.AccountNumber(ctx, cfg.PurchaseAccount)
			if err != nil {
				log.Warnf("failed to get account number for "+
					"ticket purchase account %q: %v",
					cfg.PurchaseAccount, err)
				return err
			}
			vspCfg := wallet.VSPClientConfig{
				URL:    cfg.VSPOpts.URL,
				PubKey: cfg.VSPOpts.PubKey,
				Policy: &wallet.VSPPolicy{
					MaxFee:     cfg.VSPOpts.MaxFee.Amount,
					FeeAcct:    purchaseAcct,
					ChangeAcct: changeAcct,
				},
			}
			vspClient, err = w.VSP(vspCfg)
			if err != nil {
				log.Errorf("vsp: %v", err)
				return err
			}
		}

		if cfg.MixChange || cfg.EnableTicketBuyer {
			var err error
			var lastFlag, lastLookup string
			lookup := func(flag, name string) (account uint32) {
				if err == nil {
					lastFlag = flag
					lastLookup = name
					account, err = w.AccountNumber(ctx, name)
				}
				return
			}
			var (
				purchaseAccount    uint32 // enableticketbuyer
				votingAccount      uint32 // enableticketbuyer
				mixedAccount       uint32 // (enableticketbuyer && mixing) || mixchange
				changeAccount      uint32 // (enableticketbuyer && mixing) || mixchange
				ticketSplitAccount uint32 // enableticketbuyer && mixing
			)
			if cfg.EnableTicketBuyer {
				purchaseAccount = lookup("purchaseaccount", cfg.PurchaseAccount)

				if cfg.MixingEnabled && cfg.TBOpts.VotingAccount == "" {
					err := errors.New("cannot run mixed ticketbuyer without --votingaccount")
					log.Error(err)
					return err
				}
				if cfg.TBOpts.VotingAccount != "" {
					votingAccount = lookup("ticketbuyer.votingaccount", cfg.TBOpts.VotingAccount)
				} else {
					votingAccount = purchaseAccount
				}
			}
			if (cfg.EnableTicketBuyer && cfg.MixingEnabled) || cfg.MixChange {
				mixedAccount = lookup("mixedaccount", cfg.mixedAccount)
				changeAccount = lookup("changeaccount", cfg.ChangeAccount)
			}
			if cfg.EnableTicketBuyer && cfg.MixingEnabled {
				ticketSplitAccount = lookup("ticketsplitaccount", cfg.TicketSplitAccount)
			}

			// Check if any of the above calls to lookup() have failed.
			if err != nil {
				log.Errorf("%s: account %q does not exist", lastFlag, lastLookup)
				return err
			}

			// Start a ticket buyer.
			tb := ticketbuyer.New(w, ticketbuyer.Config{
				BuyTickets:         cfg.EnableTicketBuyer,
				Account:            purchaseAccount,
				Maintain:           cfg.TBOpts.BalanceToMaintainAbsolute.Amount,
				Limit:              int(cfg.TBOpts.Limit),
				VotingAccount:      votingAccount,
				Mixing:             cfg.MixingEnabled,
				MixChange:          cfg.MixChange,
				MixedAccount:       mixedAccount,
				MixedAccountBranch: cfg.mixedBranch,
				TicketSplitAccount: ticketSplitAccount,
				ChangeAccount:      changeAccount,
				VSP:                vspClient,
			})

			log.Infof("Starting auto transaction creator")
			tbdone := make(chan struct{})
			go func() {
				err := tb.Run(ctx, passphrase)
				if err != nil && !errors.Is(err, context.Canceled) {
					log.Errorf("Transaction creator ended: %v", err)
				}
				tbdone <- struct{}{}
			}()
			defer func() { <-tbdone }()
		}
	}

	if done(ctx) {
		return ctx.Err()
	}

	// Create and start the RPC servers to serve wallet client connections.  If
	// any of the servers can not be started, it will be nil.  If none of them
	// can be started, this errors since at least one server must run for the
	// wallet to be useful.
	//
	// Servers will be associated with a loaded wallet if it has already been
	// loaded, or after it is loaded later on.
	gRPCServer, jsonRPCServer, err := startRPCServers(ctx, loader)
	if err != nil {
		log.Errorf("Unable to create RPC servers: %v", err)
		return err
	}
	if gRPCServer != nil {
		// Start wallet, voting and network gRPC services after a
		// wallet is loaded.
		loader.RunAfterLoad(func(w *wallet.Wallet) {
			rpcserver.StartWalletService(gRPCServer, w)
			rpcserver.StartNetworkService(gRPCServer, w)
			rpcserver.StartVotingService(gRPCServer, w)
		})
		defer func() {
			log.Warn("Stopping gRPC server...")
			gRPCServer.Stop()
			log.Info("gRPC server shutdown")
		}()
	}
	if jsonRPCServer != nil {
		go func() {
			for range jsonRPCServer.RequestProcessShutdown() {
				requestShutdown()
			}
		}()
		defer func() {
			log.Warn("Stopping JSON-RPC server...")
			jsonRPCServer.Stop()
			log.Info("JSON-RPC server shutdown")
		}()
	}

	// When not running with --noinitialload, it is the main package's
	// responsibility to synchronize the wallet with the network through SPV or
	// the trusted vgld server.  This blocks until cancelled.
	if !cfg.NoInitialLoad {
		if done(ctx) {
			return ctx.Err()
		}

		loader.RunAfterLoad(func(w *wallet.Wallet) {
			if vspClient != nil && cfg.VSPOpts.Sync {
				tickets, err := w.ProcessedTickets(ctx)
				if err != nil {
					log.Errorf("Getting VSP tickets failed: %v", err)
				}
				err = vspClient.ProcessManagedTickets(ctx, tickets)
				if err != nil {
					log.Errorf("Adding tickets to VSP client failed: %v", err)
				}
			}

			switch {
			case cfg.Offline:
				w.SetNetworkBackend(wallet.OfflineNetworkBackend{})
			case cfg.SPV:
				spvLoop(ctx, w)
			default:
				rpcSyncLoop(ctx, w)
			}
		})
	}

	// Wait until shutdown is signaled before returning and running deferred
	// shutdown tasks.
	<-ctx.Done()
	return ctx.Err()
}

func passPrompt(ctx context.Context, prefix string, confirm bool) (passphrase []byte, err error) {
	os.Stdout.Sync()
	c := make(chan struct{}, 1)
	go func() {
		passphrase, err = prompt.PassPrompt(bufio.NewReader(os.Stdin), prefix, confirm)
		c <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c:
		return passphrase, err
	}
}

// startPromptPass prompts the user for a password to unlock their wallet in
// the event that it was restored from seed or --promptpass flag is set.
func startPromptPass(ctx context.Context, w *wallet.Wallet) []byte {
	promptPass := cfg.PromptPass

	// Watching only wallets never require a password.
	if w.WatchingOnly() {
		return nil
	}

	// The wallet is totally desynced, so we need to resync accounts.
	// Prompt for the password. Then, set the flag it wallet so it
	// knows which address functions to call when resyncing.
	needSync, err := w.NeedsAccountsSync(ctx)
	if err != nil {
		log.Errorf("Error determining whether an accounts sync is necessary: %v", err)
	}
	if err == nil && needSync {
		fmt.Println("*** ATTENTION ***")
		fmt.Println("Since this is your first time running we need to sync accounts. Please enter")
		fmt.Println("the private wallet passphrase. This will complete syncing of the wallet")
		fmt.Println("accounts and then leave your wallet unlocked. You may relock wallet after by")
		fmt.Println("calling 'walletlock' through the RPC.")
		fmt.Println("*****************")
		promptPass = true
	}
	if cfg.EnableTicketBuyer {
		promptPass = true
	}

	if !promptPass {
		return nil
	}

	// We need to rescan accounts for the initial sync. Unlock the
	// wallet after prompting for the passphrase. The special case
	// of a --createtemp simnet wallet is handled by first
	// attempting to automatically open it with the default
	// passphrase. The wallet should also request to be unlocked
	// if stake mining is currently on, so users with this flag
	// are prompted here as well.
	for {
		if w.ChainParams().Net == wire.SimNet {
			err := w.Unlock(ctx, wallet.SimulationPassphrase, nil)
			if err == nil {
				// Unlock success with the default password.
				return wallet.SimulationPassphrase
			}
		}

		passphrase, err := passPrompt(ctx, "Enter private passphrase", false)
		if err != nil {
			return nil
		}

		err = w.Unlock(ctx, passphrase, nil)
		if err != nil {
			fmt.Println("Incorrect password entered. Please " +
				"try again.")
			continue
		}
		return passphrase
	}
}

func spvLoop(ctx context.Context, w *wallet.Wallet) {
	addr := &net.TCPAddr{IP: net.ParseIP("::1"), Port: 0}
	amgrDir := filepath.Join(cfg.AppDataDir.Value, w.ChainParams().Name)
	amgr := addrmgr.New(amgrDir, cfg.lookup)
	for {
		lp := p2p.NewLocalPeer(w.ChainParams(), addr, amgr)
		lp.SetDialFunc(cfg.dial)
		lp.SetDisableRelayTx(cfg.SPVDisableRelayTx)
		syncer := spv.NewSyncer(w, lp)
		if len(cfg.SPVConnect) > 0 {
			syncer.SetPersistentPeers(cfg.SPVConnect)
		}
		err := syncer.Run(ctx)
		if err == nil || done(ctx) {
			loggers.SyncLog.Infof("SPV synchronization stopped")
			return
		}
		loggers.SyncLog.Errorf("SPV synchronization stopped: %v", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
		}
	}
}

// rpcSyncLoop loops forever, attempting to create a connection to the
// consensus RPC server.  If this connection succeeds, the RPC client is used as
// the loaded wallet's network backend and used to keep the wallet synchronized
// to the network.  If/when the RPC connection is lost, the wallet is
// disassociated from the client and a new connection is attempmted.
func rpcSyncLoop(ctx context.Context, w *wallet.Wallet) {
	certs := readCAFile()
	clientCert, clientKey := readClientCertKey()
	dial := cfg.dial
	if cfg.NovgldProxy {
		dial = new(net.Dialer).DialContext
	}
	for {
		rpcOptions := &chain.RPCOptions{
			Address:     cfg.RPCConnect,
			DefaultPort: activeNet.JSONRPCClientPort,
			User:        cfg.vgldUsername,
			Pass:        cfg.vgldPassword,
			Dial:        dial,
			CA:          certs,
			Insecure:    cfg.DisableClientTLS,
		}
		if len(clientCert) != 0 {
			rpcOptions.User = ""
			rpcOptions.Pass = ""
			rpcOptions.ClientCert = clientCert
			rpcOptions.ClientKey = clientKey
		}
		syncer := chain.NewSyncer(w, rpcOptions)
		err := syncer.Run(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				loggers.SyncLog.Infof("RPC synchronization stopped")
				return
			}
			loggers.SyncLog.Errorf("RPC synchronization stopped: %v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}
}

func readCAFile() []byte {
	// Read certificate file if TLS is not disabled.
	var certs []byte
	if !cfg.DisableClientTLS {
		var err error
		certs, err = os.ReadFile(cfg.CAFile.Value)
		if err != nil {
			log.Warnf("Cannot open CA file: %v", err)
			// If there's an error reading the CA file, continue
			// with nil certs and without the client connection.
			certs = nil
		}
	} else {
		log.Info("Chain server RPC TLS is disabled")
	}

	return certs
}

func readClientCertKey() ([]byte, []byte) {
	if cfg.vgldAuthType != authTypeClientCert {
		return nil, nil
	}
	cert, err := os.ReadFile(cfg.vgldClientCert.Value)
	if err != nil {
		log.Warnf("Cannot open vgld RPC client certificate: %v", err)
		cert = nil
	}
	key, err := os.ReadFile(cfg.vgldClientKey.Value)
	if err != nil {
		log.Warnf("Cannot open vgld RPC client key: %v", err)
		key = nil
	}
	return cert, key
}
