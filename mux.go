package tlsmux

import (
	"net"
	"sync"
)

type Mux struct {
	mu sync.RWMutex
	hs map[string]Handler
}

func (m *Mux) Handle(serverName string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hs == nil {
		m.hs = make(map[string]Handler)
	}

	m.hs[serverName] = handler
}

func (m *Mux) Serve(c net.Conn) {
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
func (m *Mux) handler(serverName string) (Handler, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.hs == nil {
		return nil, false
	}

	handler, exists := m.hs[serverName]

	return handler, exists
}
