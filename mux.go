package tlsmux

import (
	"net"
	"strings"
	"sync"
)

// Muxer is a TCP connection multiplexer which reads the TLS ServerName indication
// to route the connection to the configured Handler.
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

// Serve reads the server name and forwards the given net.Conn to the corresponding Handler.
// TODO: handle panics
// TODO: client hello timeout
func (m *Muxer) Serve(c net.Conn) {
	defer func() { _ = c.Close() }()

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

// handler returns the Handler corresponding to the given serverName.
func (m *Muxer) handler(serverName string) (Handler, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.hs == nil {
		return nil, false
	}

	handler, exists := m.hs[strings.ToLower(serverName)]

	return handler, exists
}
