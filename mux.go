package tlsmux

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Handler is in charge of handling as connection.
type Handler interface {
	ServeConn(net.Conn) error
}

// Mux is a TCP connection multiplexer which reads the TLS server name indication to route the connection
// to the matching Handler.
type Mux struct {
	mu sync.RWMutex
	hs map[string]Handler
}

// Handle registers a Handler for the given server name.
func (m *Mux) Handle(serverName string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hs == nil {
		m.hs = make(map[string]Handler)
	}

	m.hs[strings.ToLower(serverName)] = handler
}

// Serve accepts incoming connections on the given listener and starts a go routine to serve each connection.
// TODO handle errors.
func (m *Mux) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accept: %w", err)
		}

		go func() { _ = m.ServeConn(conn) }()
	}
}

// ServeConn reads the TLS server name indication and forwards the net.Conn to the matching Handler.
// Handler implementations are responsible for closing the connection.
// TODO: handle panics
// TODO: client hello timeout.
func (m *Mux) ServeConn(c net.Conn) error {
	serverName, peeked := ClientHelloServerName(c)
	if serverName == "" {
		return errors.New("empty server name")
	}

	handler, exists := m.handler(serverName)
	if exists {
		return handler.ServeConn(&conn{c, peeked})
	}

	if err := c.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return fmt.Errorf("no handler for %s", serverName)
}

// handler returns the Handler matching the given server name value.
func (m *Mux) handler(serverName string) (Handler, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.hs == nil {
		return nil, false
	}

	handler, exists := m.hs[strings.ToLower(serverName)]

	return handler, exists
}
