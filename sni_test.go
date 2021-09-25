package tlsmux

import (
	"crypto/tls"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientHelloServerName(t *testing.T) {
	c, s := net.Pipe()

	defer func() { _ = s.Close() }()

	go func() {
		defer func() { _ = c.Close() }()

		err := tls.Client(c, &tls.Config{ServerName: "foo"}).Handshake()
		require.NoError(t, err)
	}()

	serverName, peeked := ClientHelloServerName(s)

	assert.NotEmpty(t, peeked)
	assert.Equal(t, "foo", serverName)
}
