// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"vigil.network/vgl/cspp/v2/solverrpc"
	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/internal/cfgutil"
	"github.com/kdsmith18542/vigil/wallet/internal/loggers"
	"github.com/kdsmith18542/vigil/wallet/internal/netparams"
	"github.com/kdsmith18542/vigil/wallet/version"
	"github.com/kdsmith18542/vigil/wallet/wallet"
	"github.com/kdsmith18542/vigil/wallet/wallet/txrules"
	"github.com/kdsmith18542/vigil/connmgr/v3"
	"github.com/kdsmith18542/vigil/VGLutil/v4"
	"github.com/kdsmith18542/vigil/go-socks/socks"
	"github.com/kdsmith18542/vigil/slog"
	flags "github.com/jessevdk/go-flags"
)

const (
	// Authorization types.
	authTypeBasic      = "basic"
	authTypeClientCert = "clientcert"
)

const (
	defaultCAFilename              = "vgld.cert"
	defaultConfigFilename          = "vglwallet.conf"
	defaultLogLevel                = "info"
	defaultLogDirname              = "logs"
	defaultLogFilename             = "vglwallet.log"
	defaultLogSize                 = "10M"
	defaultRPCMaxClients           = 10
	defaultRPCMaxWebsockets        = 25
	defaultAuthType                = authTypeBasic
	defaultEnableTicketBuyer       = false
	defaultEnableVoting            = false
	defaultPurchaseAccount         = "default"
	defaultPromptPass              = false
	defaultPass                    = ""
	defaultPromptPublicPass        = false
	defaultGapLimit                = wallet.DefaultGapLimit
	defaultAllowHighFees           = false
	defaultAccountGapLimit         = wallet.DefaultAccountGapLimit
	defaultDisableCoinTypeUpgrades = false
	defaultCircuitLimit            = 32
	defaultMixSplitLimit           = 10
	defaultVSPMaxFee               = VGLutil.Amount(0.2e8)

	// ticket buyer options
	defaultBalanceToMaintainAbsolute = 0
	defaultTicketbuyerLimit          = 1

	walletDbName = "wallet.db"
)

var (
	vgldDefaultCAFile         = filepath.Join(VGLutil.AppDataDir("vgld", false), "rpc.cert")
	defaultAppDataDir         = VGLutil.AppDataDir("vglwallet", false)
	defaultConfigFile         = filepath.Join(defaultAppDataDir, defaultConfigFilename)
	defaultRPCKeyFile         = filepath.Join(defaultAppDataDir, "rpc.key")
	defaultRPCCertFile        = filepath.Join(defaultAppDataDir, "rpc.cert")
	defaultvgldClientCertFile = filepath.Join(defaultAppDataDir, "vgld-client.cert")
	defaultvgldClientKeyFile  = filepath.Join(defaultAppDataDir, "vgld-client.key")
	defaultRPCClientCAFile    = filepath.Join(defaultAppDataDir, "clients.pem")
	defaultLogDir             = filepath.Join(defaultAppDataDir, defaultLogDirname)
)

type config struct {
	// General application behavior
	ConfigFile         *cfgutil.ExplicitString `short:"C" long:"configfile" description:"Path to configuration file"`
	ShowVersion        bool                    `short:"V" long:"version" description:"Display version information and exit"`
	Create             bool                    `long:"create" description:"Create new wallet"`
	CreateTemp         bool                    `long:"createtemp" description:"Create simulation wallet in nonstandard --appdata; private passphrase is 'password'"`
	CreateWatchingOnly bool                    `long:"createwatchingonly" description:"Create watching wallet from account extended pubkey"`
	AppDataDir         *cfgutil.ExplicitString `short:"A" long:"appdata" description:"Application data directory for wallet config, databases and logs"`
	TestNet            bool                    `long:"testnet" description:"Use the test network"`
	SimNet             bool                    `long:"simnet" description:"Use the simulation test network"`
	NoInitialLoad      bool                    `long:"noinitialload" description:"Defer wallet creation/opening on startup and enable loading wallets over RPC"`
	DebugLevel         string                  `short:"d" long:"debuglevel" description:"Logging level {trace, debug, info, warn, error, critical}"`
	LogDir             *cfgutil.ExplicitString `long:"logdir" description:"Directory to log output."`
	LogSize            string                  `long:"logsize" description:"Maximum size of log file before it is rotated"`
	NoFileLogging      bool                    `long:"nofilelogging" description:"Disable file logging"`
	Profile            []string                `long:"profile" description:"Enable HTTP profiling this interface/port"`
	MemProfile         string                  `long:"memprofile" description:"Write mem profile to the specified file"`
	CPUProfile         string                  `long:"cpuprofile" description:"Write cpu profile to the specified file"`

	// Wallet options
	WalletPass              string              `long:"walletpass" default-mask:"-" description:"Public wallet password; required when created with one"`
	PromptPass              bool                `long:"promptpass" description:"Prompt for private passphase from terminal and unlock without timeout"`
	Pass                    string              `long:"pass" description:"Unlock with private passphrase"`
	PromptPublicPass        bool                `long:"promptpublicpass" description:"Prompt for public passphrase from terminal"`
	EnableTicketBuyer       bool                `long:"enableticketbuyer" description:"Enable the automatic ticket buyer"`
	EnableVoting            bool                `long:"enablevoting" description:"Automatically vote on winning tickets"`
	PurchaseAccount         string              `long:"purchaseaccount" description:"Account to autobuy tickets from"`
	GapLimit                uint32              `long:"gaplimit" description:"Allowed unused address gap between used addresses of accounts"`
	WatchLast               uint32              `long:"watchlast" description:"Limit watched previous addresses of each HD account branch"`
	ManualTickets           bool                `long:"manualtickets" description:"Do not discover new tickets through network synchronization"`
	AllowHighFees           bool                `long:"allowhighfees" description:"Do not perform high fee checks"`
	RelayFee                *cfgutil.AmountFlag `long:"txfee" description:"Transaction fee per kilobyte"`
	AccountGapLimit         int                 `long:"accountgaplimit" description:"Allowed gap of unused accounts"`
	DisableCoinTypeUpgrades bool                `long:"disablecointypeupgrades" description:"Never upgrade from legacy to SLIP0044 coin type keys"`

	// RPC client options
	RPCConnect       string                  `short:"c" long:"rpcconnect" description:"Network address of vgld RPC server"`
	CAFile           *cfgutil.ExplicitString `long:"cafile" description:"vgld RPC Certificate Authority"`
	ClientCAFile     *cfgutil.ExplicitString `long:"clientcafile" description:"Certficate Authority to verify TLS client certificates"`
	DisableClientTLS bool                    `long:"noclienttls" description:"Disable TLS for vgld RPC; only allowed when connecting to localhost"`
	vgldUsername     string                  `long:"vgldusername" description:"vgld RPC username; overrides --username"`
	vgldPassword     string                  `long:"vgldpassword" default-mask:"-" description:"vgld RPC password; overrides --password"`
	vgldClientCert   *cfgutil.ExplicitString `long:"vgldclientcert" description:"TLS client certificate to present to authenticate RPC connections to vgld"`
	vgldClientKey    *cfgutil.ExplicitString `long:"vgldclientkey" description:"Key for vgld RPC client certificate"`
	vgldAuthType     string                  `long:"vgldauthtype" description:"Method for vgld JSON-RPC client authentication (basic or clientcert)"`

	// Proxy and Tor settings
	Proxy        string `long:"proxy" description:"Establish network connections and DNS lookups through a SOCKS5 proxy (e.g. 127.0.0.1:9050)"`
	ProxyUser    string `long:"proxyuser" description:"Proxy server username"`
	ProxyPass    string `long:"proxypass" default-mask:"-" description:"Proxy server password"`
	CircuitLimit int    `long:"circuitlimit" description:"Set maximum number of open Tor circuits; used only when --torisolation is enabled"`
	TorIsolation bool   `long:"torisolation" description:"Enable Tor stream isolation by randomizing user credentials for each connection"`
	NovgldProxy  bool   `long:"novgldproxy" description:"Never use configured proxy to dial vgld websocket connectons"`
	dial         func(ctx context.Context, network, address string) (net.Conn, error)
	lookup       func(name string) ([]net.IP, error)

	// Offline mode.
	Offline bool `long:"offline" description:"Do not sync the wallet"`

	// SPV options
	SPV               bool     `long:"spv" description:"Sync using simplified payment verification"`
	SPVConnect        []string `long:"spvconnect" description:"SPV sync only with specified peers; disables DNS seeding"`
	SPVDisableRelayTx bool     `long:"spvdisablerelaytx" description:"Disable receiving mempool transactions when in SPV mode"`

	// RPC server options
	RPCCert                *cfgutil.ExplicitString `long:"rpccert" description:"RPC server TLS certificate"`
	RPCKey                 *cfgutil.ExplicitString `long:"rpckey" description:"RPC server TLS key"`
	TLSCurve               *cfgutil.CurveFlag      `long:"tlscurve" description:"Curve to use when generating TLS keypairs"`
	OneTimeTLSKey          bool                    `long:"onetimetlskey" description:"Generate self-signed TLS keypairs each startup; only write certificate file"`
	DisableServerTLS       bool                    `long:"noservertls" description:"Disable TLS for the RPC servers; only allowed when binding to localhost"`
	GRPCListeners          []string                `long:"grpclisten" description:"Listen for gRPC connections on this interface"`
	LegacyRPCListeners     []string                `long:"rpclisten" description:"Listen for JSON-RPC connections on this interface"`
	NoGRPC                 bool                    `long:"nogrpc" description:"Disable gRPC server"`
	NoLegacyRPC            bool                    `long:"nolegacyrpc" description:"Disable JSON-RPC server"`
	LegacyRPCMaxClients    int64                   `long:"rpcmaxclients" description:"Max JSON-RPC HTTP POST clients"`
	LegacyRPCMaxWebsockets int64                   `long:"rpcmaxwebsockets" description:"Max JSON-RPC websocket clients"`
	Username               string                  `short:"u" long:"username" description:"JSON-RPC username and default vgld RPC username"`
	Password               string                  `short:"P" long:"password" default-mask:"-" description:"JSON-RPC password and default vgld RPC password"`
	JSONRPCAuthType        string                  `long:"jsonrpcauthtype" description:"Method for JSON-RPC client authentication (basic or clientcert)"`

	// IPC options
	PipeTx            *uint `long:"pipetx" description:"File descriptor or handle of write end pipe to enable child -> parent process communication"`
	PipeRx            *uint `long:"piperx" description:"File descriptor or handle of read end pipe to enable parent -> child process communication"`
	RPCListenerEvents bool  `long:"rpclistenerevents" description:"Notify JSON-RPC and gRPC listener addresses over the TX pipe"`
	IssueClientCert   bool  `long:"issueclientcert" description:"Notify a client cert and key over the TX pipe for RPC authentication"`

	// CSPP
	MixingEnabled      bool                    `long:"mixing" description:"Enable creation of mixed transactions and participation in the peer-to-peer mixing network"`
	CSPPSolver         *cfgutil.ExplicitString `long:"csppsolver" description:"Path to CSPP solver executable (if not in PATH)"`
	MixedAccount       string                  `long:"mixedaccount" description:"Account/branch used to derive CoinShuffle++ mixed outputs and voting rewards"`
	mixedAccount       string
	mixedBranch        uint32
	TicketSplitAccount string `long:"ticketsplitaccount" description:"Account to derive fresh addresses from for mixed ticket splits; uses mixedaccount if unset"`
	ChangeAccount      string `long:"changeaccount" description:"Account used to derive unmixed CoinJoin outputs in CoinShuffle++ protocol"`
	MixChange          bool   `long:"mixchange" description:"Use CoinShuffle++ to mix change account outputs into mix account"`
	MixSplitLimit      int    `long:"mixsplitlimit" description:"Connection limit to CoinShuffle++ server per change amount"`

	TBOpts ticketBuyerOptions `group:"Ticket Buyer Options" namespace:"ticketbuyer"`

	VSPOpts vspOptions `group:"VSP Options" namespace:"vsp"`
}

type ticketBuyerOptions struct {
	BalanceToMaintainAbsolute *cfgutil.AmountFlag `long:"balancetomaintainabsolute" description:"Amount of funds to keep in wallet when purchasing tickets"`
	Limit                     uint                `long:"limit" description:"Buy no more than specified number of tickets per block"`
	VotingAccount             string              `long:"votingaccount" description:"Account used to derive addresses specifying voting rights"`
}

type vspOptions struct {
	// VSP - TODO: VSPServer to a []string to support multiple VSPs
	URL    string              `long:"url" description:"Base URL of the VSP server"`
	PubKey string              `long:"pubkey" description:"VSP server pubkey"`
	Sync   bool                `long:"sync" description:"sync tickets to vsp"`
	MaxFee *cfgutil.AmountFlag `long:"maxfee" description:"Maximum VSP fee"`
}

// cleanAndExpandPath expands environement variables and leading ~ in the
// passed path, cleans the result, and returns it.
func cleanAndExpandPath(path string) string {
	// Do not try to clean the empty string
	if path == "" {
		return ""
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows cmd.exe-style
	// %VARIABLE%, but they variables can still be expanded via POSIX-style
	// $VARIABLE.
	path = os.ExpandEnv(path)

	if !strings.HasPrefix(path, "~") {
		return filepath.Clean(path)
	}

	// Expand initial ~ to the current user's home directory, or ~otheruser
	// to otheruser's home directory.  On Windows, both forward and backward
	// slashes can be used.
	path = path[1:]

	var pathSeparators string
	if runtime.GOOS == "windows" {
		pathSeparators = string(os.PathSeparator) + "/"
	} else {
		pathSeparators = string(os.PathSeparator)
	}

	userName := ""
	if i := strings.IndexAny(path, pathSeparators); i != -1 {
		userName = path[:i]
		path = path[i:]
	}

	homeDir := ""
	var u *user.User
	var err error
	if userName == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(userName)
	}
	if err == nil {
		homeDir = u.HomeDir
	}
	// Fallback to CWD if user lookup fails or user has no home directory.
	if homeDir == "" {
		homeDir = "."
	}

	return filepath.Join(homeDir, path)
}

// validLogLevel returns whether or not logLevel is a valid debug log level.
func validLogLevel(logLevel string) bool {
	_, ok := slog.LevelFromString(logLevel)
	return ok
}

// supportedSubsystems returns a sorted slice of the supported subsystems for
// logging purposes.
func supportedSubsystems() []string {
	// Convert the subsystemLoggers map keys to a slice.
	subsystems := make([]string, 0, len(subsystemLoggers))
	for subsysID := range subsystemLoggers {
		subsystems = append(subsystems, subsysID)
	}

	// Sort the subsytems for stable display.
	sort.Strings(subsystems)
	return subsystems
}

// parseAndSetDebugLevels attempts to parse the specified debug level and set
// the levels accordingly.  An appropriate error is returned if anything is
// invalid.
func parseAndSetDebugLevels(debugLevel string) error {
	// When the specified string doesn't have any delimters, treat it as
	// the log level for all subsystems.
	if !strings.Contains(debugLevel, ",") && !strings.Contains(debugLevel, "=") {
		// Validate debug log level.
		if !validLogLevel(debugLevel) {
			str := "The specified debug level [%v] is invalid"
			return errors.Errorf(str, debugLevel)
		}

		// Change the logging level for all subsystems.
		setLogLevels(debugLevel)

		return nil
	}

	// Split the specified string into subsystem/level pairs while detecting
	// issues and update the log levels accordingly.
	for _, logLevelPair := range strings.Split(debugLevel, ",") {
		if !strings.Contains(logLevelPair, "=") {
			str := "The specified debug level contains an invalid " +
				"subsystem/level pair [%v]"
			return errors.Errorf(str, logLevelPair)
		}

		// Extract the specified subsystem and log level.
		fields := strings.Split(logLevelPair, "=")
		subsysID, logLevel := fields[0], fields[1]

		// Validate subsystem.
		if _, exists := subsystemLoggers[subsysID]; !exists {
			str := "The specified subsystem [%v] is invalid -- " +
				"supported subsytems %v"
			return errors.Errorf(str, subsysID, supportedSubsystems())
		}

		// Validate log level.
		if !validLogLevel(logLevel) {
			str := "The specified debug level [%v] is invalid"
			return errors.Errorf(str, logLevel)
		}

		setLogLevel(subsysID, logLevel)
	}

	return nil
}

// loadConfig initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
//  1. Start with a default config with sane settings
//  2. Pre-parse the command line to check for an alternative config file
//  3. Load configuration file overwriting defaults with any specified options
//  4. Parse CLI options and overwrite/add any specified options
//
// The above results in vglwallet functioning properly without any config
// settings while still allowing the user to override settings with config files
// and command line options.  Command line options always take precedence.
// The bool returned indicates whether or not the wallet was recreated from a
// seed and needs to perform the initial resync. The []byte is the private
// passphrase required to do the sync for this special case.
func loadConfig(ctx context.Context) (*config, []string, error) {
	loadConfigError := func(err error) (*config, []string, error) {
		return nil, nil, err
	}

	// Default config.
	cfg := config{
		DebugLevel:              defaultLogLevel,
		ConfigFile:              cfgutil.NewExplicitString(defaultConfigFile),
		AppDataDir:              cfgutil.NewExplicitString(defaultAppDataDir),
		LogDir:                  cfgutil.NewExplicitString(defaultLogDir),
		LogSize:                 defaultLogSize,
		WalletPass:              wallet.InsecurePubPassphrase,
		CAFile:                  cfgutil.NewExplicitString(""),
		ClientCAFile:            cfgutil.NewExplicitString(defaultRPCClientCAFile),
		vgldClientCert:          cfgutil.NewExplicitString(defaultvgldClientCertFile),
		vgldClientKey:           cfgutil.NewExplicitString(defaultvgldClientKeyFile),
		dial:                    new(net.Dialer).DialContext,
		lookup:                  net.LookupIP,
		PromptPass:              defaultPromptPass,
		Pass:                    defaultPass,
		PromptPublicPass:        defaultPromptPublicPass,
		RPCKey:                  cfgutil.NewExplicitString(defaultRPCKeyFile),
		RPCCert:                 cfgutil.NewExplicitString(defaultRPCCertFile),
		TLSCurve:                cfgutil.NewCurveFlag(cfgutil.PreferredCurve),
		LegacyRPCMaxClients:     defaultRPCMaxClients,
		LegacyRPCMaxWebsockets:  defaultRPCMaxWebsockets,
		JSONRPCAuthType:         defaultAuthType,
		vgldAuthType:            defaultAuthType,
		EnableTicketBuyer:       defaultEnableTicketBuyer,
		EnableVoting:            defaultEnableVoting,
		PurchaseAccount:         defaultPurchaseAccount,
		GapLimit:                defaultGapLimit,
		AllowHighFees:           defaultAllowHighFees,
		RelayFee:                cfgutil.NewAmountFlag(txrules.DefaultRelayFeePerKb),
		AccountGapLimit:         defaultAccountGapLimit,
		DisableCoinTypeUpgrades: defaultDisableCoinTypeUpgrades,
		CircuitLimit:            defaultCircuitLimit,
		MixSplitLimit:           defaultMixSplitLimit,
		CSPPSolver:              cfgutil.NewExplicitString(solverrpc.SolverProcess),

		// Ticket Buyer Options
		TBOpts: ticketBuyerOptions{
			BalanceToMaintainAbsolute: cfgutil.NewAmountFlag(defaultBalanceToMaintainAbsolute),
			Limit:                     defaultTicketbuyerLimit,
		},

		VSPOpts: vspOptions{
			MaxFee: cfgutil.NewAmountFlag(defaultVSPMaxFee),
		},
	}

	// Pre-parse the command line options to see if an alternative config
	// file or the version flag was specified.
	preCfg := cfg
	preParser := flags.NewParser(&preCfg, flags.Default)
	_, err := preParser.Parse()
	if err != nil {
		var e *flags.Error
		if errors.As(err, &e) && e.Type == flags.ErrHelp {
			os.Exit(0)
		}
		preParser.WriteHelp(os.Stderr)
		return loadConfigError(err)
	}

	// Show the version and exit if the version flag was specified.
	funcName := "loadConfig"
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	usageMessage := fmt.Sprintf("Use %s -h to show usage", appName)
	if preCfg.ShowVersion {
		fmt.Printf("%s version %s (Go version %s %s/%s)\n", appName,
			version.String(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Load additional config from file.
	var configFileError error
	parser := flags.NewParser(&cfg, flags.Default)
	configFilePath := preCfg.ConfigFile.Value
	if preCfg.ConfigFile.ExplicitlySet() {
		configFilePath = cleanAndExpandPath(configFilePath)
	} else {
		appDataDir := preCfg.AppDataDir.Value
		if appDataDir != defaultAppDataDir {
			configFilePath = filepath.Join(appDataDir, defaultConfigFilename)
		}
	}
	err = flags.NewIniParser(parser).ParseFile(configFilePath)
	if err != nil {
		var e *os.PathError
		if !errors.As(err, &e) {
			fmt.Fprintln(os.Stderr, err)
			parser.WriteHelp(os.Stderr)
			return loadConfigError(err)
		}
		configFileError = err
	}

	// Parse command line options again to ensure they take precedence.
	remainingArgs, err := parser.Parse()
	if err != nil {
		var e *flags.Error
		if !errors.As(err, &e) || e.Type != flags.ErrHelp {
			parser.WriteHelp(os.Stderr)
		}
		return loadConfigError(err)
	}

	// If an alternate data directory was specified, and paths with defaults
	// relative to the data dir are unchanged, modify each path to be
	// relative to the new data dir.
	if cfg.AppDataDir.ExplicitlySet() {
		cfg.AppDataDir.Value = cleanAndExpandPath(cfg.AppDataDir.Value)
		if !cfg.RPCKey.ExplicitlySet() {
			cfg.RPCKey.Value = filepath.Join(cfg.AppDataDir.Value, "rpc.key")
		}
		if !cfg.RPCCert.ExplicitlySet() {
			cfg.RPCCert.Value = filepath.Join(cfg.AppDataDir.Value, "rpc.cert")
		}
		if !cfg.ClientCAFile.ExplicitlySet() {
			cfg.ClientCAFile.Value = filepath.Join(cfg.AppDataDir.Value, "clients.pem")
		}
		if !cfg.vgldClientCert.ExplicitlySet() {
			cfg.vgldClientCert.Value = filepath.Join(cfg.AppDataDir.Value, "vgld-client.cert")
		}
		if !cfg.vgldClientKey.ExplicitlySet() {
			cfg.vgldClientKey.Value = filepath.Join(cfg.AppDataDir.Value, "vgld-client.key")
		}
		if !cfg.LogDir.ExplicitlySet() {
			cfg.LogDir.Value = filepath.Join(cfg.AppDataDir.Value, defaultLogDirname)
		}
	}

	// Choose the active network params based on the selected network.
	// Multiple networks can't be selected simultaneously.
	numNets := 0
	if cfg.TestNet {
		activeNet = &netparams.TestNet3Params
		numNets++
	}
	if cfg.SimNet {
		activeNet = &netparams.SimNetParams
		numNets++
	}
	if numNets > 1 {
		str := "%s: The testnet and simnet params can't be used " +
			"together -- choose one"
		err := errors.Errorf(str, "loadConfig")
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}

	if !cfg.NoFileLogging {
		// Append the network type to the log directory so it is
		// "namespaced" per network.
		cfg.LogDir.Value = cleanAndExpandPath(cfg.LogDir.Value)
		cfg.LogDir.Value = filepath.Join(cfg.LogDir.Value,
			activeNet.Params.Name)

		var units int
		for i, r := range cfg.LogSize {
			if r < '0' || r > '9' {
				units = i
				break
			}
		}
		invalidSize := func() error {
			str := "%s: Invalid logsize: %v "
			err := errors.Errorf(str, funcName, cfg.LogSize)
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		if units == 0 {
			return loadConfigError(invalidSize())
		}
		// Parsing a 32-bit number prevents 64-bit overflow after unit
		// multiplication.
		logsize, err := strconv.ParseInt(cfg.LogSize[:units], 10, 32)
		if err != nil {
			return loadConfigError(invalidSize())
		}
		switch cfg.LogSize[units:] {
		case "k", "K", "KiB":
		case "m", "M", "MiB":
			logsize <<= 10
		case "g", "G", "GiB":
			logsize <<= 20
		default:
			return loadConfigError(invalidSize())
		}

		// Initialize log rotation.  After log rotation has been initialized, the
		// logger variables may be used.
		loggers.InitLogRotator(filepath.Join(cfg.LogDir.Value, defaultLogFilename), logsize)
	}

	// Special show command to list supported subsystems and exit.
	if cfg.DebugLevel == "show" {
		fmt.Println("Supported subsystems", supportedSubsystems())
		os.Exit(0)
	}

	// Parse, validate, and set debug log level(s).
	if err := parseAndSetDebugLevels(cfg.DebugLevel); err != nil {
		err := errors.Errorf("%s: %v", "loadConfig", err.Error())
		fmt.Fprintln(os.Stderr, err)
		parser.WriteHelp(os.Stderr)
		return loadConfigError(err)
	}

	// Error and shutdown if config file is specified on the command line
	// but cannot be found.
	if configFileError != nil && cfg.ConfigFile.ExplicitlySet() {
		if preCfg.ConfigFile.ExplicitlySet() || cfg.ConfigFile.ExplicitlySet() {
			log.Errorf("%v", configFileError)
			return loadConfigError(configFileError)
		}
	}

	// Warn about missing config file after the final command line parse
	// succeeds.  This prevents the warning on help messages and invalid
	// options.
	if configFileError != nil {
		log.Warnf("%v", configFileError)
	}

	// Sanity check BalanceToMaintainAbsolute
	if cfg.TBOpts.BalanceToMaintainAbsolute.ToCoin() < 0 {
		str := "%s: balancetomaintainabsolute cannot be negative: %v"
		err := errors.Errorf(str, funcName, cfg.TBOpts.BalanceToMaintainAbsolute)
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}

	// Exit if you try to use a simulation wallet with a standard
	// data directory.
	if !cfg.AppDataDir.ExplicitlySet() && cfg.CreateTemp {
		fmt.Fprintln(os.Stderr, "Tried to create a temporary simulation "+
			"wallet, but failed to specify data directory!")
		os.Exit(0)
	}

	// Exit if you try to use a simulation wallet on anything other than
	// simnet or testnet.
	if !cfg.SimNet && cfg.CreateTemp {
		fmt.Fprintln(os.Stderr, "Tried to create a temporary simulation "+
			"wallet for network other than simnet!")
		os.Exit(0)
	}

	// Ensure the wallet exists or create it when the create flag is set.
	netDir := networkDir(cfg.AppDataDir.Value, activeNet.Params)
	dbPath := filepath.Join(netDir, walletDbName)

	if cfg.CreateTemp && cfg.Create {
		err := errors.Errorf("The flags --create and --createtemp can not " +
			"be specified together. Use --help for more information.")
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}

	dbFileExists, err := cfgutil.FileExists(dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}

	if cfg.CreateTemp {
		tempWalletExists := false

		if dbFileExists {
			str := fmt.Sprintf("The wallet already exists. Loading this " +
				"wallet instead.")
			fmt.Fprintln(os.Stdout, str)
			tempWalletExists = true
		}

		// Ensure the data directory for the network exists.
		if err := checkCreateDir(netDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}

		if !tempWalletExists {
			// Perform the initial wallet creation wizard.
			if err := createSimulationWallet(ctx, &cfg); err != nil {
				fmt.Fprintln(os.Stderr, "Unable to create wallet:", err)
				return loadConfigError(err)
			}
		}
	} else if cfg.Create || cfg.CreateWatchingOnly {
		// Error if the create flag is set and the wallet already
		// exists.
		if dbFileExists {
			err := errors.Errorf("The wallet database file `%v` "+
				"already exists.", dbPath)
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}

		// Ensure the data directory for the network exists.
		if err := checkCreateDir(netDir); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}

		// Perform the initial wallet creation wizard.
		os.Stdout.Sync()
		if cfg.CreateWatchingOnly {
			err = createWatchingOnlyWallet(ctx, &cfg)
		} else {
			err = createWallet(ctx, &cfg)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to create wallet:", err)
			return loadConfigError(err)
		}

		// Created successfully, so exit now with success.
		os.Exit(0)
	} else if !dbFileExists && !cfg.NoInitialLoad {
		err := errors.Errorf("The wallet does not exist.  Run with the " +
			"--create option to initialize and create it.")
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}

	ipNet := func(cidr string) net.IPNet {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		return *ipNet
	}
	privNets := []net.IPNet{
		// IPv4 loopback
		ipNet("127.0.0.0/8"),

		// IPv6 loopback
		ipNet("::1/128"),

		// RFC 1918
		ipNet("10.0.0.0/8"),
		ipNet("172.16.0.0/12"),
		ipNet("192.168.0.0/16"),

		// RFC 4193
		ipNet("fc00::/7"),
	}

	// Set dialer and DNS lookup functions if proxy settings are provided.
	if cfg.Proxy != "" {
		proxy := socks.Proxy{
			Addr:         cfg.Proxy,
			Username:     cfg.ProxyUser,
			Password:     cfg.ProxyPass,
			TorIsolation: cfg.TorIsolation,
		}

		var proxyDialer func(context.Context, string, string) (net.Conn, error)
		var noproxyDialer net.Dialer
		if cfg.TorIsolation {
			proxyDialer = socks.NewPool(proxy, uint32(cfg.CircuitLimit)).DialContext
		} else {
			proxyDialer = proxy.DialContext
		}

		cfg.dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				host = address
			}
			if host == "localhost" {
				return noproxyDialer.DialContext(ctx, network, address)
			}
			ip := net.ParseIP(host)
			if len(ip) == 4 || len(ip) == 16 {
				for i := range privNets {
					if privNets[i].Contains(ip) {
						return noproxyDialer.DialContext(ctx, network, address)
					}
				}
			}
			conn, err := proxyDialer(ctx, network, address)
			if err != nil {
				return nil, errors.Errorf("proxy dial %v %v: %w", network, address, err)
			}
			return conn, nil
		}
		cfg.lookup = func(host string) ([]net.IP, error) {
			ip, err := connmgr.TorLookupIP(context.Background(), host, cfg.Proxy)
			if err != nil {
				return nil, errors.Errorf("proxy lookup for %v: %w", host, err)
			}
			return ip, nil
		}
	}

	var solverMustWork bool
	if cfg.MixingEnabled {
		if cfg.CSPPSolver.ExplicitlySet() {
			solverrpc.SolverProcess = cfg.CSPPSolver.Value
			solverMustWork = true
		} else if err := solverrpc.StartSolver(); err == nil {
			solverMustWork = true
		} else {
			log.Warnf("Unable to start csppsolver; must rely on " +
				"other peers publishing results")
		}
	}
	if solverMustWork {
		if err := testStartedSolverWorks(); err != nil {
			err := errors.Errorf("csppsolver process is not operating properly: %v", err)
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
	}

	// Parse mixedaccount account/branch
	if cfg.MixedAccount != "" {
		indexSlash := strings.LastIndex(cfg.MixedAccount, "/")
		if indexSlash == -1 {
			err := errors.Errorf("--mixedaccount must have form 'accountname/branch'")
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		cfg.mixedAccount = cfg.MixedAccount[:indexSlash]
		switch cfg.MixedAccount[indexSlash+1:] {
		case "0":
			cfg.mixedBranch = 0
		case "1":
			cfg.mixedBranch = 1
		default:
			err := errors.Errorf("--mixedaccount branch must be 0 or 1")
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
	}
	// Use mixedaccount as default ticketsplitaccount if unset.
	if cfg.TicketSplitAccount == "" {
		cfg.TicketSplitAccount = cfg.mixedAccount
	}

	if cfg.RPCConnect == "" {
		cfg.RPCConnect = net.JoinHostPort("localhost", activeNet.JSONRPCClientPort)
	}

	// Add default port to connect flag if missing.
	cfg.RPCConnect, err = cfgutil.NormalizeAddress(cfg.RPCConnect,
		activeNet.JSONRPCClientPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid rpcconnect network address: %v\n", err)
		return loadConfigError(err)
	}

	localhostListeners := map[string]struct{}{
		"localhost": {},
		"127.0.0.1": {},
		"::1":       {},
	}
	RPCHost, _, err := net.SplitHostPort(cfg.RPCConnect)
	if err != nil {
		return loadConfigError(err)
	}
	if cfg.DisableClientTLS {
		if _, ok := localhostListeners[RPCHost]; !ok {
			str := "%s: the --noclienttls option may not be used " +
				"when connecting RPC to non localhost " +
				"addresses: %s"
			err := errors.Errorf(str, funcName, cfg.RPCConnect)
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, usageMessage)
			return loadConfigError(err)
		}
	} else {
		// If CAFile is unset, choose either the copy or local vgld cert.
		if !cfg.CAFile.ExplicitlySet() {
			cfg.CAFile.Value = filepath.Join(cfg.AppDataDir.Value, defaultCAFilename)

			// If the CA copy does not exist, check if we're connecting to
			// a local vgld and switch to its RPC cert if it exists.
			certExists, err := cfgutil.FileExists(cfg.CAFile.Value)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return loadConfigError(err)
			}
			if !certExists {
				if _, ok := localhostListeners[RPCHost]; ok {
					vgldCertExists, err := cfgutil.FileExists(
						vgldDefaultCAFile)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return loadConfigError(err)
					}
					if vgldCertExists {
						cfg.CAFile.Value = vgldDefaultCAFile
					}
				}
			}
		}
	}

	if cfg.SPV && cfg.Offline {
		err := errors.E("SPV and Offline mode cannot be specified at the same time")
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}

	if cfg.SPV && cfg.EnableVoting {
		err := errors.E("SPV voting is not possible: disable --spv or --enablevoting")
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}
	if !cfg.SPV && len(cfg.SPVConnect) > 0 {
		err := errors.E("--spvconnect requires --spv")
		fmt.Fprintln(os.Stderr, err)
		return loadConfigError(err)
	}
	for i, p := range cfg.SPVConnect {
		cfg.SPVConnect[i], err = cfgutil.NormalizeAddress(p, activeNet.Params.DefaultPort)
		if err != nil {
			return loadConfigError(err)
		}
	}

	// Default to localhost listen addresses if no listeners were manually
	// specified.  When the RPC server is configured to be disabled, remove all
	// listeners so it is not started.
	localhostAddrs, err := net.LookupHost("localhost")
	if err != nil {
		return loadConfigError(err)
	}
	if len(cfg.GRPCListeners) == 0 && !cfg.NoGRPC {
		cfg.GRPCListeners = make([]string, 0, len(localhostAddrs))
		for _, addr := range localhostAddrs {
			cfg.GRPCListeners = append(cfg.GRPCListeners,
				net.JoinHostPort(addr, activeNet.GRPCServerPort))
		}
	} else if cfg.NoGRPC {
		cfg.GRPCListeners = nil
	}
	if len(cfg.LegacyRPCListeners) == 0 && !cfg.NoLegacyRPC {
		cfg.LegacyRPCListeners = make([]string, 0, len(localhostAddrs))
		for _, addr := range localhostAddrs {
			cfg.LegacyRPCListeners = append(cfg.LegacyRPCListeners,
				net.JoinHostPort(addr, activeNet.JSONRPCServerPort))
		}
	} else if cfg.NoLegacyRPC {
		cfg.LegacyRPCListeners = nil
	}

	// Add default port to all rpc listener addresses if needed and remove
	// duplicate addresses.
	cfg.LegacyRPCListeners, err = cfgutil.NormalizeAddresses(
		cfg.LegacyRPCListeners, activeNet.JSONRPCServerPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid network address in legacy RPC listeners: %v\n", err)
		return loadConfigError(err)
	}
	cfg.GRPCListeners, err = cfgutil.NormalizeAddresses(
		cfg.GRPCListeners, activeNet.GRPCServerPort)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"Invalid network address in RPC listeners: %v\n", err)
		return loadConfigError(err)
	}

	// Both RPC servers may not listen on the same interface/port, with the
	// exception of listeners using port 0.
	if len(cfg.LegacyRPCListeners) > 0 && len(cfg.GRPCListeners) > 0 {
		seenAddresses := make(map[string]struct{}, len(cfg.LegacyRPCListeners))
		for _, addr := range cfg.LegacyRPCListeners {
			seenAddresses[addr] = struct{}{}
		}
		for _, addr := range cfg.GRPCListeners {
			_, seen := seenAddresses[addr]
			if seen && !strings.HasSuffix(addr, ":0") {
				err := errors.Errorf("Address `%s` may not be "+
					"used as a listener address for both "+
					"RPC servers", addr)
				fmt.Fprintln(os.Stderr, err)
				return loadConfigError(err)
			}
		}
	}

	// Only allow server TLS to be disabled if the RPC server is bound to
	// localhost addresses.
	if cfg.DisableServerTLS {
		allListeners := append(cfg.LegacyRPCListeners, cfg.GRPCListeners...)
		for _, addr := range allListeners {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				str := "%s: RPC listen interface '%s' is " +
					"invalid: %v"
				err := errors.Errorf(str, funcName, addr, err)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return loadConfigError(err)
			}
			if _, ok := localhostListeners[host]; !ok {
				str := "%s: the --noservertls option may not be used " +
					"when binding RPC to non localhost " +
					"addresses: %s"
				err := errors.Errorf(str, funcName, addr)
				fmt.Fprintln(os.Stderr, err)
				fmt.Fprintln(os.Stderr, usageMessage)
				return loadConfigError(err)
			}
		}
	}

	// If either VSP pubkey or URL are specified, validate VSP options.
	if cfg.VSPOpts.PubKey != "" || cfg.VSPOpts.URL != "" {
		if cfg.VSPOpts.PubKey == "" {
			err := errors.New("vsp pubkey can not be null")
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		if cfg.VSPOpts.URL == "" {
			err := errors.New("vsp URL can not be null")
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		if cfg.VSPOpts.MaxFee.Amount == 0 {
			err := errors.New("vsp max fee must be greater than zero")
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
	}

	// Expand environment variable and leading ~ for filepaths.
	cfg.CAFile.Value = cleanAndExpandPath(cfg.CAFile.Value)
	cfg.RPCCert.Value = cleanAndExpandPath(cfg.RPCCert.Value)
	cfg.RPCKey.Value = cleanAndExpandPath(cfg.RPCKey.Value)
	cfg.vgldClientCert.Value = cleanAndExpandPath(cfg.vgldClientCert.Value)
	cfg.vgldClientKey.Value = cleanAndExpandPath(cfg.vgldClientKey.Value)
	cfg.ClientCAFile.Value = cleanAndExpandPath(cfg.ClientCAFile.Value)

	// If the vgld username or password are unset, use the same auth as for
	// the client.  The two settings were previously shared for vgld and
	// client auth, so this avoids breaking backwards compatibility while
	// allowing users to use different auth settings for vgld and wallet.
	if cfg.vgldUsername == "" {
		cfg.vgldUsername = cfg.Username
	}
	if cfg.vgldPassword == "" {
		cfg.vgldPassword = cfg.Password
	}

	switch cfg.vgldAuthType {
	case authTypeBasic:
	case authTypeClientCert:
		if cfg.DisableClientTLS {
			err := fmt.Errorf("vgldauthtype=clientcert is " +
				"incompatible with disableclienttls")
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		vgldClientCertExists, _ := cfgutil.FileExists(
			cfg.vgldClientCert.Value)
		if !vgldClientCertExists {
			err := fmt.Errorf("vgldclientcert %q is required "+
				"by vgldauthtype=clientcert but does not exist",
				cfg.vgldClientCert.Value)
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
		vgldClientKeyExists, _ := cfgutil.FileExists(
			cfg.vgldClientKey.Value)
		if !vgldClientKeyExists {
			err := fmt.Errorf("vgldclientkey %q is required "+
				"by vgldauthtype=clientcert but does not exist",
				cfg.vgldClientKey.Value)
			fmt.Fprintln(os.Stderr, err)
			return loadConfigError(err)
		}
	default:
		err := fmt.Errorf("unknown vgld authtype %q", cfg.vgldAuthType)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return loadConfigError(err)
	}

	switch cfg.JSONRPCAuthType {
	case authTypeBasic, authTypeClientCert:
	default:
		err := fmt.Errorf("unknown JSON-RPC authtype %q", cfg.JSONRPCAuthType)
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, usageMessage)
		return loadConfigError(err)
	}

	// Make list of old versions of testnet directories.
	var oldTestNets []string
	oldTestNets = append(oldTestNets, filepath.Join(cfg.AppDataDir.Value, "testnet"))
	// Warn if old testnet directory is present.
	for _, oldDir := range oldTestNets {
		oldDirExists, _ := cfgutil.FileExists(oldDir)
		if oldDirExists {
			log.Warnf("Wallet data from previous testnet"+
				" found (%v) and can probably be removed.",
				oldDir)
		}
	}

	return &cfg, remainingArgs, nil
}
