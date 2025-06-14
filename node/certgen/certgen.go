// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2024 The Vigil developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package certgen

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	_ "crypto/sha512" // Needed for RegisterHash in init
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
	"unicode"

	"github.com/kdsmith18542/vigil/crypto/rand"
	"golang.org/x/net/idna"
)

func isASCII(s string) bool {
	for _, c := range s {
		if c > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// NewTLSCertPair returns a new PEM-encoded x.509 certificate pair with new
// ECDSA keys.  The machine's local interface addresses and all variants of IPv4
// and IPv6 localhost are included as valid IP addresses.
func NewTLSCertPair(curve elliptic.Curve, organization string, validUntil time.Time, extraHosts []string) (cert, key []byte, err error) {
	now := time.Now()
	if validUntil.Before(now) {
		return nil, nil, errors.New("validUntil would create an already-expired certificate")
	}

	priv, err := ecdsa.GenerateKey(curve, rand.Reader())
	if err != nil {
		return nil, nil, err
	}

	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	if validUntil.After(endOfTime) {
		validUntil = endOfTime
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber := rand.BigInt(serialNumberLimit)

	host, err := os.Hostname()
	if err != nil {
		return nil, nil, err
	}
	if !isASCII(host) {
		host, err = idna.ToASCII(host)
		if err != nil {
			return nil, nil, err
		}
	}

	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	dnsNames := []string{host}
	if host != "localhost" {
		dnsNames = append(dnsNames, "localhost")
	}

	addIP := func(ipAddr net.IP) {
		for _, ip := range ipAddresses {
			if ip.Equal(ipAddr) {
				return
			}
		}
		ipAddresses = append(ipAddresses, ipAddr)
	}
	addHost := func(host string) {
		for _, dnsName := range dnsNames {
			if host == dnsName {
				return
			}
		}
		if !isASCII(host) {
			var err error
			host, err = idna.ToASCII(host)
			if err != nil {
				return
			}
		}
		dnsNames = append(dnsNames, host)
	}

	addrs, err := interfaceAddrs()
	if err != nil {
		return nil, nil, err
	}
	for _, a := range addrs {
		ipAddr, _, err := net.ParseCIDR(a.String())
		if err == nil {
			addIP(ipAddr)
		}
	}

	for _, hostStr := range extraHosts {
		host, _, err := net.SplitHostPort(hostStr)
		if err != nil {
			host = hostStr
		}
		if ip := net.ParseIP(host); ip != nil {
			addIP(ip)
		} else {
			addHost(host)
		}
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   host,
		},
		NotBefore: now.Add(-time.Hour * 24),
		NotAfter:  validUntil,

		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature |
			x509.KeyUsageCertSign,
		IsCA:                  true, // so can sign self.
		BasicConstraintsValid: true,

		DNSNames:    dnsNames,
		IPAddresses: ipAddresses,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader(), &template,
		&template, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certBuf := &bytes.Buffer{}
	err = pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode certificate: %w", err)
	}

	keybytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	keyBuf := &bytes.Buffer{}
	err = pem.Encode(keyBuf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keybytes})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	return certBuf.Bytes(), keyBuf.Bytes(), nil
}
