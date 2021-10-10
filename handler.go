package tlsmux

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
)

// Handler is in charge of handling a raw connection.
type Handler interface {
	ServeConn(net.Conn) error
}

// HandlerFunc is an adapter to allow the use of a function as a Handler.
type HandlerFunc func(net.Conn) error

func (h HandlerFunc) ServeConn(conn net.Conn) error {
	return h(conn)
}

// TLSHandler is a Handler implementation handling TLS connection by using the configured tls.Config.
type TLSHandler struct {
	Handler

	Config *tls.Config
}

func (h TLSHandler) ServeConn(conn net.Conn) error {
	tlsConn := tls.Server(conn, h.Config)

	if err := tlsConn.Handshake(); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	return h.Handler.ServeConn(tlsConn)
}

// TLSHandlerFunc is an adapter to allow the use of a function as a TLSHandler.
func TLSHandlerFunc(config *tls.Config, handler HandlerFunc) TLSHandler {
	return TLSHandler{handler, config}
}

// ProxyHandler is a Handler implementation forwarding the connection bytes to the configured Addr.
type ProxyHandler struct {
	Addr string
}

// TODO handle errors.
func (p ProxyHandler) ServeConn(conn net.Conn) error {
	defer func() { _ = conn.Close() }()

	dstConn, err := net.Dial("tcp", p.Addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer func() { _ = dstConn.Close() }()

	go func() { _, _ = io.Copy(dstConn, conn) }()
	_, _ = io.Copy(conn, dstConn)

	return nil
}

// ProxyHandlerFunc is an adapter to allow the use of a ProxyHandler as a HandlerFunc.
func ProxyHandlerFunc(addr string) HandlerFunc {
	return ProxyHandler{Addr: addr}.ServeConn
}
