package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"time"

	"github.com/kevinpollet/tlsmux"
)

func main() {
	m := tlsmux.Muxer{}
	m.Handle("foo.localhost", tlsmux.TLSHandlerFunc(tlsConfig("foo.localhost"), func(conn net.Conn) {
		defer func() { _ = conn.Close() }()

		_, _ = conn.Write([]byte("foo"))
	}))
	m.Handle("bar.localhost", tlsmux.TLSHandlerFunc(tlsConfig("bar.localhost"), func(conn net.Conn) {
		defer func() { _ = conn.Close() }()

		_, _ = conn.Write([]byte("bar"))
	}))

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go m.ServeConn(conn)
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
