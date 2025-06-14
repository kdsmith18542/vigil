// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/kdsmith18542/vigil/wallet/errors"
	"github.com/kdsmith18542/vigil/wallet/internal/cfgutil"
	"github.com/kdsmith18542/vigil/wallet/internal/loader"
	"github.com/kdsmith18542/vigil/wallet/internal/loggers"
	"github.com/kdsmith18542/vigil/wallet/internal/rpc/jsonrpc"
	"github.com/kdsmith18542/vigil/wallet/internal/rpc/rpcserver"
	"github.com/kdsmith18542/vigil/crypto/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// openRPCKeyPair creates or loads the RPC TLS keypair specified by the
// application config.  This function respects the cfg.OneTimeTLSKey setting.
func openRPCKeyPair() (tls.Certificate, error) {
	// Check for existence of the TLS key file.  If one time TLS keys are
	// enabled but a key already exists, this function should error since
	// it's possible that a persistent certificate was copied to a remote
	// machine.  Otherwise, generate a new keypair when the key is missing.
	// When generating new persistent keys, overwriting an existing cert is
	// acceptable if the previous execution used a one time TLS key.
	// Otherwise, both the cert and key should be read from disk.  If the
	// cert is missing, the read error will occur in LoadX509KeyPair.
	_, e := os.Stat(cfg.RPCKey.Value)
	keyExists := !os.IsNotExist(e)
	switch {
	case cfg.OneTimeTLSKey && keyExists:
		err := errors.Errorf("one time TLS keys are enabled, but TLS key "+
			"`%s` already exists", cfg.RPCKey)
		return tls.Certificate{}, err
	case cfg.OneTimeTLSKey:
		return generateRPCKeyPair(false)
	case !keyExists:
		return generateRPCKeyPair(true)
	default:
		return tls.LoadX509KeyPair(cfg.RPCCert.Value, cfg.RPCKey.Value)
	}
}

// generateRPCKeyPair generates a new RPC TLS keypair and writes the cert and
// possibly also the key in PEM format to the paths specified by the config.  If
// successful, the new keypair is returned.
func generateRPCKeyPair(writeKey bool) (tls.Certificate, error) {
	log.Infof("Generating TLS certificates...")

	// Create directories for cert and key files if they do not yet exist.
	certDir, _ := filepath.Split(cfg.RPCCert.Value)
	keyDir, _ := filepath.Split(cfg.RPCKey.Value)
	err := os.MkdirAll(certDir, 0700)
	if err != nil {
		return tls.Certificate{}, err
	}
	err = os.MkdirAll(keyDir, 0700)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Generate cert pair.
	org := "vglwallet autogenerated cert"
	validUntil := time.Now().Add(time.Hour * 24 * 365 * 10)
	cert, key, err := cfg.TLSCurve.CertGen(org, validUntil, nil)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyPair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return tls.Certificate{}, err
	}

	// Write cert and (potentially) the key files.
	err = os.WriteFile(cfg.RPCCert.Value, cert, 0600)
	if err != nil {
		return tls.Certificate{}, err
	}
	if writeKey {
		err = os.WriteFile(cfg.RPCKey.Value, key, 0600)
		if err != nil {
			rmErr := os.Remove(cfg.RPCCert.Value)
			if rmErr != nil {
				log.Warnf("Cannot remove written certificates: %v",
					rmErr)
			}
			return tls.Certificate{}, err
		}
	}

	log.Info("Done generating TLS certificates")
	return keyPair, nil
}

func randomX509SerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber := rand.BigInt(serialNumberLimit)
	return serialNumber
}

// End of ASN.1 time
var endOfTime = time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)

type ClientCA struct {
	CertBlock  []byte
	Cert       *x509.Certificate
	PrivateKey any
}

func generateAuthority(pub, priv any) (*ClientCA, error) {
	validUntil := time.Now().Add(time.Hour * 24 * 365 * 10)
	now := time.Now()
	if validUntil.After(endOfTime) {
		validUntil = endOfTime
	}
	if validUntil.Before(now) {
		return nil, fmt.Errorf("valid until date %v already elapsed", validUntil)
	}
	serialNumber := randomX509SerialNumber()
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         "vglwallet",
			Organization:       []string{"vglwallet"},
			OrganizationalUnit: []string{"vglwallet certificate authority"},
		},
		NotBefore:             now.Add(-time.Hour * 24),
		NotAfter:              validUntil,
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	cert, err := x509.CreateCertificate(rand.Reader(), template, template, pub, priv)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		return nil, fmt.Errorf("failed to encode certificate: %v", err)
	}
	certBlock := buf.Bytes()

	x509Cert, err := x509.ParseCertificate(cert)
	if err != nil {
		return nil, err
	}

	clientCA := &ClientCA{
		CertBlock:  certBlock,
		Cert:       x509Cert,
		PrivateKey: priv,
	}
	return clientCA, nil
}

func marshalPrivateKey(key any) ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal private key: %v", err)
	}
	buf := new(bytes.Buffer)
	err = pem.Encode(buf, &pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key: %v", err)
	}
	return buf.Bytes(), nil
}

func createSignedClientCert(pub, caPriv any, ca *x509.Certificate) ([]byte, error) {
	serialNumber := randomX509SerialNumber()
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		NotBefore:    time.Now().Add(-time.Hour * 24),
		NotAfter:     ca.NotAfter,
		Subject: pkix.Name{
			CommonName:         "vglwallet",
			Organization:       []string{"vglwallet"},
			OrganizationalUnit: []string{"vglwallet client certificate"},
		},
	}
	cert, err := x509.CreateCertificate(rand.Reader(), template, ca, pub, caPriv)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if err != nil {
		return nil, fmt.Errorf("failed to encode certificate: %v", err)
	}
	return buf.Bytes(), nil
}

func generateClientKeyPair(caPriv any, ca *x509.Certificate) (cert, key []byte, err error) {
	pub, priv, err := cfg.TLSCurve.GenerateKeyPair(rand.Reader())
	if err != nil {
		return
	}
	key, err = marshalPrivateKey(priv)
	if err != nil {
		return
	}
	cert, err = createSignedClientCert(pub, caPriv, ca)
	if err != nil {
		return
	}
	return cert, key, nil
}

type rpcLoggers struct{}

func (rpcLoggers) Subsystems() []string {
	s := make([]string, 0, len(subsystemLoggers))
	for name := range subsystemLoggers {
		s = append(s, name)
	}
	sort.Strings(s)
	return s
}

func (rpcLoggers) SetLevels(levelSpec string) error {
	return parseAndSetDebugLevels(levelSpec)
}

func startRPCServers(ctx context.Context, walletLoader *loader.Loader) (*grpc.Server, *jsonrpc.Server, error) {
	var jsonrpcAddrNotifier jsonrpcListenerEventServer
	var grpcAddrNotifier grpcListenerEventServer
	if cfg.RPCListenerEvents {
		jsonrpcAddrNotifier = newJSONRPCListenerEventServer(outgoingPipeMessages)
		grpcAddrNotifier = newGRPCListenerEventServer(outgoingPipeMessages)
	}

	var (
		server         *grpc.Server
		jsonrpcServer  *jsonrpc.Server
		jsonrpcListen  = net.Listen
		keyPair        tls.Certificate
		clientCAsExist bool
		err            error
	)
	if cfg.DisableServerTLS {
		log.Info("Server TLS is disabled.  Only JSON-RPC may be used")
	} else {
		keyPair, err = openRPCKeyPair()
		if err != nil {
			return nil, nil, err
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{keyPair},
			MinVersion:   tls.VersionTLS12,
			ClientCAs:    x509.NewCertPool(),
		}
		clientCAsExist, _ = cfgutil.FileExists(cfg.ClientCAFile.Value)
		if clientCAsExist {
			cafile, err := os.ReadFile(cfg.ClientCAFile.Value)
			if err != nil {
				return nil, nil, err
			}
			if !tlsConfig.ClientCAs.AppendCertsFromPEM(cafile) {
				log.Warnf("No certificates added from CA file %v",
					cfg.ClientCAFile)
			}
		}
		if cfg.JSONRPCAuthType == "clientcert" {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
		if cfg.IssueClientCert {
			pub, priv, err := cfg.TLSCurve.GenerateKeyPair(rand.Reader())
			if err != nil {
				return nil, nil, err
			}
			ca, err := generateAuthority(pub, priv)
			if err != nil {
				return nil, nil, err
			}
			certBlock, keyBlock, err := generateClientKeyPair(ca.PrivateKey, ca.Cert)
			if err != nil {
				return nil, nil, err
			}
			tlsConfig.ClientCAs.AddCert(ca.Cert)

			s := newIssuedClientCertEventServer(outgoingPipeMessages)
			s.notify(keyBlock, certBlock, ca.CertBlock)
		}

		// Change the standard net.Listen function to the tls one.
		jsonrpcListen = func(net string, laddr string) (net.Listener, error) {
			return tls.Listen(net, laddr, tlsConfig)
		}

		clientCAsExist = clientCAsExist || cfg.IssueClientCert
		if !clientCAsExist && len(cfg.GRPCListeners) != 0 {
			log.Warnf("gRPC server is configured with listeners, but no "+
				"trusted client certificates exist (looked in %v)",
				cfg.ClientCAFile)
		} else if clientCAsExist && len(cfg.GRPCListeners) != 0 {
			tlsConfig := tlsConfig.Clone()
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
			listeners := makeListeners(cfg.GRPCListeners, net.Listen)
			if len(listeners) == 0 {
				err := errors.New("failed to create listeners for RPC server")
				return nil, nil, err
			}
			server = grpc.NewServer(
				grpc.Creds(credentials.NewTLS(tlsConfig)),
				grpc.StreamInterceptor(interceptStreaming),
				grpc.UnaryInterceptor(interceptUnary),
			)
			rpcserver.RegisterServices(server)
			rpcserver.StartWalletLoaderService(server, walletLoader, activeNet)
			rpcserver.StartTicketBuyerService(server, walletLoader)
			rpcserver.StartAccountMixerService(server, walletLoader)
			rpcserver.StartAgendaService(server, activeNet.Params)
			rpcserver.StartDecodeMessageService(server, activeNet.Params)
			rpcserver.StartMessageVerificationService(server, activeNet.Params)
			for _, lis := range listeners {
				lis := lis
				go func() {
					laddr := lis.Addr().String()
					grpcAddrNotifier.notify(laddr)
					log.Infof("gRPC server listening on %s", laddr)
					err := server.Serve(lis)
					log.Tracef("Finished serving gRPC: %v", err)
				}()
			}
		}
	}

	if !cfg.DisableServerTLS && len(cfg.LegacyRPCListeners) != 0 &&
		cfg.JSONRPCAuthType == "clientcert" && !clientCAsExist {
		log.Warnf("JSON-RPC TLS server is configured with listeners and "+
			"client cert auth, but no trusted client certificates exist "+
			"(looked in %v)", cfg.ClientCAFile)
	} else if cfg.JSONRPCAuthType == "basic" && (cfg.Username == "" || cfg.Password == "") {
		log.Info("JSON-RPC server disabled (basic auth requires username and " +
			"password, and client cert authentication is not enabled)")
	} else if len(cfg.LegacyRPCListeners) != 0 {
		listeners := makeListeners(cfg.LegacyRPCListeners, jsonrpcListen)
		if len(listeners) == 0 {
			err := errors.New("failed to create listeners for JSON-RPC server")
			return nil, nil, err
		}
		var user, pass string
		if cfg.JSONRPCAuthType == "basic" {
			user, pass = cfg.Username, cfg.Password
		}
		opts := jsonrpc.Options{
			Username:            user,
			Password:            pass,
			MaxPOSTClients:      cfg.LegacyRPCMaxClients,
			MaxWebsocketClients: cfg.LegacyRPCMaxWebsockets,
			MixingEnabled:       cfg.MixingEnabled,
			MixAccount:          cfg.mixedAccount,
			MixBranch:           cfg.mixedBranch,
			MixChangeAccount:    cfg.ChangeAccount,
			VSPHost:             cfg.VSPOpts.URL,
			VSPPubKey:           cfg.VSPOpts.PubKey,
			TicketSplitAccount:  cfg.TicketSplitAccount,
			Dial:                cfg.dial,
			Loggers:             rpcLoggers{},
		}
		jsonrpcServer = jsonrpc.NewServer(ctx, &opts, activeNet.Params, walletLoader, listeners)
		for _, lis := range listeners {
			jsonrpcAddrNotifier.notify(lis.Addr().String())
		}
	}

	// Error when neither the GRPC nor JSON-RPC servers can be started.
	if server == nil && jsonrpcServer == nil {
		return nil, nil, errors.New("no suitable RPC services can be started")
	}

	return server, jsonrpcServer, nil
}

// serviceName returns the package.service segment from the full gRPC method
// name `/package.service/method`.
func serviceName(method string) string {
	// Slice off first /
	method = method[1:]
	// Keep everything before the next /
	return method[:strings.IndexRune(method, '/')]
}

func interceptStreaming(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	p, ok := peer.FromContext(ss.Context())
	if ok {
		loggers.GrpcLog.Debugf("Streaming method %s invoked by %s", info.FullMethod,
			p.Addr.String())
	}
	err := rpcserver.ServiceReady(serviceName(info.FullMethod))
	if err != nil {
		return err
	}
	err = handler(srv, ss)
	if err != nil && ok {
		logf := loggers.GrpcLog.Errorf
		if status.Code(err) == codes.Canceled && done(ss.Context()) {
			// Canceled contexts in streaming calls are expected
			// when client-initiated, so only log them with debug
			// level to reduce clutter.
			logf = loggers.GrpcLog.Debugf
		}

		logf("Streaming method %s invoked by %s errored: %v",
			info.FullMethod, p.Addr.String(), err)
	}
	return err
}

func interceptUnary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		loggers.GrpcLog.Debugf("Unary method %s invoked by %s", info.FullMethod,
			p.Addr.String())
	}
	err = rpcserver.ServiceReady(serviceName(info.FullMethod))
	if err != nil {
		return nil, err
	}
	resp, err = handler(ctx, req)
	if err != nil && ok {
		loggers.GrpcLog.Errorf("Unary method %s invoked by %s errored: %v",
			info.FullMethod, p.Addr.String(), err)
	}
	return resp, err
}

type listenFunc func(net string, laddr string) (net.Listener, error)

// makeListeners splits the normalized listen addresses into IPv4 and IPv6
// addresses and creates new net.Listeners for each with the passed listen func.
// Invalid addresses are logged and skipped.
func makeListeners(normalizedListenAddrs []string, listen listenFunc) []net.Listener {
	ipv4Addrs := make([]string, 0, len(normalizedListenAddrs)*2)
	ipv6Addrs := make([]string, 0, len(normalizedListenAddrs)*2)
	for _, addr := range normalizedListenAddrs {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			// Shouldn't happen due to already being normalized.
			log.Errorf("`%s` is not a normalized "+
				"listener address", addr)
			continue
		}

		// Empty host or host of * on plan9 is both IPv4 and IPv6.
		if host == "" || (host == "*" && runtime.GOOS == "plan9") {
			ipv4Addrs = append(ipv4Addrs, addr)
			ipv6Addrs = append(ipv6Addrs, addr)
			continue
		}

		// Remove the IPv6 zone from the host, if present.  The zone
		// prevents ParseIP from correctly parsing the IP address.
		// ResolveIPAddr is intentionally not used here due to the
		// possibility of leaking a DNS query over Tor if the host is a
		// hostname and not an IP address.
		zoneIndex := strings.Index(host, "%")
		if zoneIndex != -1 {
			host = host[:zoneIndex]
		}

		ip := net.ParseIP(host)
		switch {
		case ip == nil:
			log.Warnf("`%s` is not a valid IP address", host)
		case ip.To4() == nil:
			ipv6Addrs = append(ipv6Addrs, addr)
		default:
			ipv4Addrs = append(ipv4Addrs, addr)
		}
	}
	listeners := make([]net.Listener, 0, len(ipv6Addrs)+len(ipv4Addrs))
	for _, addr := range ipv4Addrs {
		listener, err := listen("tcp4", addr)
		if err != nil {
			log.Warnf("Can't listen on %s: %v", addr, err)
			continue
		}
		listeners = append(listeners, listener)
	}
	for _, addr := range ipv6Addrs {
		listener, err := listen("tcp6", addr)
		if err != nil {
			log.Warnf("Can't listen on %s: %v", addr, err)
			continue
		}
		listeners = append(listeners, listener)
	}
	return listeners
}
