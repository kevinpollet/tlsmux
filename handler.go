package tlsmux

import (
	"crypto/tls"
	"net"
)

// Handler is in charge of handling raw TCP connection.
type Handler interface {
	Serve(conn net.Conn)
}

// HandlerFunc is an adapter allowing the use of plain func as a Handler.
type HandlerFunc func(conn net.Conn)

func (h HandlerFunc) Serve(conn net.Conn) {
	h(conn)
}

// TLSHandler is in charge of handling TLS connection by using the configured tls.Config.
type TLSHandler struct {
	Handler
	Config *tls.Config
}

func (t TLSHandler) Serve(conn net.Conn) {
	tlsConn := tls.Server(conn, t.Config)

	if err := tlsConn.Handshake(); err != nil {
		return
	}

	t.Handler.Serve(tlsConn)
}

// TLSHandlerFunc is an adapter allowing the use of plain func as a TLSHandler.
func TLSHandlerFunc(config *tls.Config, handler Handler) Handler {
	return TLSHandler{handler, config}
}
