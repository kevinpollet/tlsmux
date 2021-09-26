package tlsmux

import (
	"crypto/tls"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientHelloServerName(t *testing.T) {
	c, s := net.Pipe()

	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)

		err := tls.Client(c, &tls.Config{ServerName: "foo"}).Handshake()
		require.Error(t, io.EOF, err)

		err = c.Close()
		require.NoError(t, err)
	}()

	serverName, peeked := ClientHelloServerName(s)

	err := s.Close()
	require.NoError(t, err)

	<-doneCh

	assert.NotEmpty(t, peeked)
	assert.Equal(t, "foo", serverName)
}
