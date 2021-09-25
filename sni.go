package tlsmux

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
)

// ClientHelloServerName reads the TLS server name from the given net.Conn and returns it with the peeked bytes.
func ClientHelloServerName(conn net.Conn) (string, []byte) {
	var (
		serverName string
		peeked     bytes.Buffer
	)

	cfg := &tls.Config{
		GetConfigForClient: func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
			serverName = hello.ServerName
			return nil, nil
		},
	}

	_ = tls.
		Server(readOnlyConn{Reader: io.TeeReader(conn, &peeked)}, cfg).
		Handshake()

	return serverName, peeked.Bytes()
}

type readOnlyConn struct {
	net.Conn // panic on unwanted method calls.
	io.Reader
}

// Read reads data from the connection.
func (r readOnlyConn) Read(b []byte) (n int, err error) {
	return r.Reader.Read(b)
}

// Write returns an error on call.
func (r readOnlyConn) Write(_ []byte) (n int, err error) {
	return -1, io.EOF
}
