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

	"github.com/Vigil-Labs/vgl/addrmgr"
	"github.com/Vigil-Labs/vgl/wallet/chain"
	"github.com/Vigil-Labs/vgl/wallet/errors"
	ldr "github.com/Vigil-Labs/vgl/wallet/internal/loader"
	"github.com/Vigil-Labs/vgl/wallet/internal/loggers"
	"github.com/Vigil-Labs/vgl/wallet/internal/prompt"
	"github.com/Vigil-Labs/vgl/wallet/internal/rpc/rpcserver"
	"github.com/Vigil-Labs/vgl/wallet/p2p"
	"github.com/Vigil-Labs/vgl/wallet/spv"
	"github.com/Vigil-Labs/vgl/wallet/ticketbuyer"
	"github.com/Vigil-Labs/vgl/wallet/version"
	"github.com/Vigil-Labs/vgl/wallet/wallet"
	"github.com/Vigil-Labs/vgl/wire"
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

	// Run the vigilwallet until permanent failure or shutdown is requested.
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

func fatalf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

// readStdin reads a single line from stdin.
func readStdin(prompt string) (string, error) {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text(), nil
	}
	return "", scanner.Err()
}

// promptPrivPass prompts the user for a private passphrase and returns it.
func promptPrivPass() ([]byte, error) {
	privPass, err := prompt.PrivatePass(os.Stdin, os.Stdout, 100)
	if err != nil {
		return nil, err
	}
	return privPass, nil
}

// networkDir returns the directory name of a network given the application
// directory and a wallet net.
func networkDir(appDir string, netParams *chain.Params) string {
	return filepath.Join(appDir, netParams.Name)
}

// newChainService creates a rpc client for the specified protocol after
// ensuring the specified protocol is supported.
func newChainService(chainServer, protocol, cert, clientCert, clientKey string, mixClientCert, mixClientKey []string) (chain.Interface, error) {
	// ... existing code ...
	// The loader is used to create and load wallets.
	// This is the default --appdata dir, vigilwallet uses an explicit
	// appdata dir passed to the loader.
	loader := ldr.NewLoader(activeNet.Params, dbDir, cfg.EnableVoting,
		cfg.GapLimit, cfg.WatchLast, cfg.AllowHighFees, cfg.RelayFee.Amount,
		cfg.VSPOpts.MaxFee.Amount, cfg.AccountGapLimit,
		cfg.DisableCoinTypeUpgrades, cfg.MixingEnabled, cfg.ManualTickets,
		cfg.MixSplitLimit, cfg.dial)
	// ... existing code ...
}

// vigilwallet.
type config struct {
	ShowVersion             bool           `short:"V" long:"version" description:"Display version information and exit"`
	ConfigFile              string         `short:"C" long:"configfile" description:"Path to configuration file"`
	AppDataDir              flags.Filename `long:"appdata" description:"Directory to store application data"`
	LogDir                  flags.Filename `long:"logdir" description:"Directory to log output."`
	NoFileLogging           bool           `long:"nofilelogging" description:"Disable logging to files.  Logging will only be printed to the console."`
	TestNet                 bool           `long:"testnet" description:"Use the test network"`
	SimNet                  bool           `long:"simnet" description:"Use the simulation test network"`
	RegressionTest          bool           `long:"regtest" description:"Use the regression test network"`
	DebugLevel              string         `short:"d" long:"debuglevel" description:"Logging level for all modules {trace, debug, info, warn, error, critical}"`
	Profile                 []string       `long:"profile" description:"Enable HTTP profiling for given comma separated listen IP:port and block profile on shutdown. Port defaults to 6060 if not specified."`
	CPUProfile              string         `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	MemProfile              string         `long:"memprofile" description:"Write mem profile to the specified file on shutdown"`
	NoInitialLoad           bool           `long:"noinitialload" description:"Don't load wallet until Wallet RPCs are used"`
	ClientTLS               bool           `long:"clienttls" description:"Enable TLS for the RPC client"`
	ClientRPCListeners      []string       `long:"clientrpclisten" description:"Comma-separated list of interfaces and/or :ports to listen for RPC connections"`
	GrpcMaxRecvMsgSize      uint           `long:"grpcmaxrecvmsgsize" description:"Maximum message size accepted by gRPC servers. This is limited to 4MiB, to prevent DoS attacks."`
	GrpcMaxSendMsgSize      uint           `long:"grpcmaxsendmsgsize" description:"Maximum message size the gRPC servers can send. This is limited to 4MiB, to prevent DoS attacks."`
	NoLegacyRPC             bool           `long:"nolegacyrpc" description:"Do not enable the legacy RPC server."`
	LegacyRPCListeners      []string       `long:"legacyservicelistener" description:"Comma-separated list of interfaces and/or :ports to listen for legacy RPC connections"`
	RPCMaxClients           int            `long:"rpcmaxclients" description:"Max number of RPC clients for legacy RPC servers"`
	DisableRPC              bool           `long:"disablerpc" description:"Disable RPC server (both gRPC and legacy)."`
	WalletPass              string         `long:"walletpass" description:"The public passphrase for the wallet. This is an insecure option and should only be used for testing."`
	PromptPublicPass        bool           `long:"promptpublicpass" description:"Prompt for the public wallet passphrase"`
	Pass                    string         `long:"pass" description:"The private passphrase for the wallet. This is an insecure option and should only be used for testing."`
	Create                  bool           `long:"create" description:"Create a new wallet with the provided passphrase and exit."`
	AllowErrors             bool           `long:"allowerrors" description:"Allow RPC to return internal errors. THIS IS INSECURE AND SHOULD NOT BE USED OUTSIDE OF TESTING."`
	EnableVoting            bool           `long:"enablevoting" description:"Enable ticket buyer voting.  Default true"`
	MaxFeeRate              float64        `long:"maxfeerate" description:"Maximum fee rate in BTC/kB accepted. Transactions with higher fee rates will not be relayed or mined. (deprecated)"`
	RelayFee                flags.Amount   `long:"relayfee" description:"Fee rate in DCR/kB for transactions that are relayed to other peers, will be converted to atom/kB and rounded up to the nearest multiple of 256. (default: 0.0001 DCR/kB)"`
	SigCacheMaxEntries      uint           `long:"sigcachemaxentries" default:"131072" description:"Maximum number of entries in the signature verification cache."`
	AccountGapLimit         uint32         `long:"accountgaplimit" description:"The maximum number of unused consecutive accounts that will be created and watched on-chain.  Used to be known as 'gaplimit'.  (default: 20)"`
	GapLimit                uint32         `long:"gaplimit" description:"The maximum number of unused consecutive addresses that will be created and watched on-chain. (default: 20)"`
	NoPeerCoinbase          bool           `long:"nopeercoinbase" description:"Disable coinbase transaction requests from peers."`
	ClientTLSKey            string         `long:"clienttlskey" description:"Path to client TLS key file"`
	ClientTLSCert           string         `long:"clienttlscert" description:"Path to client TLS certificate file"`
	CAFile                  string         `long:"cafile" description:"Path to trusted root certificate file for the RPC client"`
	WalletRPCMaxClients     int            `long:"walletrpccmaxclients" description:"Max number of RPC clients for the wallet RPC server"`
	P2PDisableTLS           bool           `long:"p2pdisabletls" description:"Disable TLS for the p2p client. THIS IS INSECURE AND SHOULD ONLY BE USED FOR TESTING."`
	EnableP2P               bool           `long:"enablep2p" description:"Enable the p2p client."`
	P2PListenAddrs          []string       `long:"p2plisten" description:"Comma-separated list of interfaces and/or :ports to listen for p2p connections."`
	P2PAddPeers             []string       `long:"p2paddpeer" description:"Add a peer to connect to at startup (can be used multiple times)"`
	P2PConnectPeers         []string       `long:"p2pconnect" description:"Connect only to the specified peers at startup (can be used multiple times)"`
	P2PTargetOutbound       int            `long:"p2ptargetoutbound" default:"8" description:"Target number of outbound peers to connect to."`
	P2PMaxPeers             int            `long:"p2pmaxpeers" default:"125" description:"Max number of p2p peers"`
	P2PDisableDNSSeed       bool           `long:"p2pdisablednsseed" description:"Disable DNS seeding for p2p connections."`
	P2PReadTimeout          time.Duration  `long:"p2preadtimeout" default:"30s" description:"Timeout for reads from p2p peers."`
	P2PWriteTimeout         time.Duration  `long:"p2pwritetimeout" default:"30s" description:"Timeout for writes to p2p peers."`
	NodeRPCServer           string         `long:"noderpcserver" description:"Host:port for the node RPC server"`
	RPCUser                 string         `long:"rpcuser" description:"Username for RPC connections"`
	RPCPass                 string         `long:"rpcpass" description:"Password for RPC connections"`
	PromptForValues         bool           `long:"promptforvalues" description:"Prompt for missing config values if needed"`
	Proxy                   string         `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser               string         `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass               string         `long:"proxypass" description:"Password for proxy server"`
	ProxySOCKS              int            `long:"proxysocks" default:"5" description:"SOCKS proxy version (4 or 5)"`
	MixClientCert           flags.Filename `long:"mixclientcert" description:"Path to mixer client TLS certificate file (can be used multiple times)"`
	MixClientKey            flags.Filename `long:"mixclientkey" description:"Path to mixer client TLS key file (can be used multiple times)"`
	MixerCert               flags.Filename `long:"mixercert" description:"Path to trusted mixer server TLS certificate file (can be used multiple times)"`
	MixerHost               string         `long:"mixerhost" description:"Host:port for the mixer server"`
	MixingEnabled           bool           `long:"mixingenabled" description:"Enable CoinShuffle++ mixing."`
	ManualTickets           bool           `long:"manualtickets" description:"Don't participate in automatic ticket buying"`
	MixChangeAccount        string         `long:"mixchangeaccount" description:"Account to send change from mixed transactions to (default: default account)"`
	MixOutputAccount        string         `long:"mixoutputaccount" description:"Account to send mixed outputs to (default: default account)"`
	MixMaxTickets           int            `long:"mixmaxtickets" description:"Maximum number of tickets to include in a single CoinShuffle++ mix."`
	MixSplitLimit           int            `long:"mixsplitlimit" description:"Maximum number of mixed outputs to create from a single CoinShuffle++ mix."`
	NoWalletRepair          bool           `long:"nowalletrepair" description:"Disable automatic wallet repair."`
	DeprecatedArgs          []string       `hidden:"true" description:"Catch-all for deprecated arguments to avoid breaking builds."`
	DisableCoinTypeUpgrades bool           `long:"disablecointypeupgrades" description:"Disable automatic cointypeupgrades. THIS IS INSECURE AND SHOULD NOT BE USED OUTSIDE OF TESTING."`
	VSPOpts                 struct {
		Host     string         `long:"host" description:"Address of the VSP to use for automatic ticket purchasing"`
		APIKey   string         `long:"apikey" description:"API key for the VSP"`
		Dial     string         `long:"dial" description:"Network address to dial (e.g. "tcp", "tcp4", "tcp6", "unix", "unixpacket")"`
		MaxFee   flags.Amount   `long:"maxfee" description:"Maximum fee the VSP can charge for automatic ticket purchasing."`
		MaxPrice flags.Amount   `long:"maxprice" description:"Maximum ticket price for automatic ticket purchasing."`
		PoolFees flags.Amount   `long:"poolfees" description:"The percentage of the ticket price that will be paid as fees to the VSP."`
		VSPPubs  []flags.PubKey `long:"vsppub" description:"Public key of the VSP. If this option is provided, the VSP's certificate will be validated against this public key."`
		CertFile string         `long:"certfile" description:"Path to the VSP's TLS certificate file."`
	} `group:"VSP Options"`
}



