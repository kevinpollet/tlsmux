package tlsmux

import (
	"net"
	"strings"
	"sync"
)

// Muxer is a TCP connection mux which reads the TLS server name indication
// to route the connection to the matching Handler.
type Muxer struct {
	mu sync.RWMutex
	hs map[string]Handler
}

// Handle registers a Handler for the given server name.
func (m *Muxer) Handle(serverName string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hs == nil {
		m.hs = make(map[string]Handler)
	}

	m.hs[strings.ToLower(serverName)] = handler
}

// ServeConn reads the TLS server name indication and forwards the net.Conn to the matching Handler.
// Handler implementations are responsible for closing the connection.
// TODO: handle panics?
// TODO: client hello timeout.
func (m *Muxer) ServeConn(c net.Conn) {
	serverName, peeked := ClientHelloServerName(c)
	if serverName == "" {
		return
	}

	handler, exists := m.handler(serverName)
	if !exists {
		return
	}

	handler.Serve(&conn{c, peeked})
}

// handler returns the Handler matching the given server name value.
func (m *Muxer) handler(serverName string) (Handler, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.hs == nil {
		return nil, false
	}

	handler, exists := m.hs[strings.ToLower(serverName)]

	return handler, exists
}
