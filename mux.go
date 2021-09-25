package tlsmux

import (
	"net"
	"strings"
	"sync"
)

type Muxer struct {
	mu sync.RWMutex
	hs map[string]Handler
}

func (m *Muxer) Handle(serverName string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.hs == nil {
		m.hs = make(map[string]Handler)
	}

	m.hs[strings.ToLower(serverName)] = handler
}

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
