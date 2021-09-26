package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"math/big"
	"net"
	"time"

	"github.com/kevinpollet/tlsmux"
)

func main() {
	mux := tlsmux.Mux{}

	mux.Handle("httpbin.org", tlsmux.HandlerFunc(func(conn net.Conn) error {
		defer func() { _ = conn.Close() }()

		dst, err := net.Dial("tcp", "httpbin.org:443")
		if err != nil {
			return fmt.Errorf("dial: %w", err)
		}
		defer func() { _ = dst.Close() }()

		go func() { _, _ = io.Copy(dst, conn) }()
		_, _ = io.Copy(conn, dst)

		return nil
	}))

	mux.Handle("foo.localhost", tlsmux.TLSHandlerFunc(tlsConfig("foo.localhost"), func(conn net.Conn) error {
		defer func() { _ = conn.Close() }()

		_, err := io.WriteString(conn, "foo")

		return err
	}))

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	if err := mux.Serve(l); err != nil {
		panic(err)
	}
}

func tlsConfig(dnsName string) *tls.Config {
	keyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour)

	tpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"tlsmux"},
			CommonName:   "Self signed cert",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              []string{dnsName},
		BasicConstraintsValid: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, &keyPair.PublicKey, keyPair)
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{der},
			PrivateKey:  keyPair,
		}},
	}
}
