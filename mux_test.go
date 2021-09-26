package tlsmux

import (
	"crypto/tls"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMuxer_Handle(t *testing.T) {
	m := Mux{}

	m.Handle("Foo", HandlerFunc(func(_ net.Conn) error {
		return nil
	}))

	assert.Len(t, m.hs, 1)
	assert.Contains(t, m.hs, "foo")
}

func TestMux_ServeConn(t *testing.T) {
	tests := []struct {
		desc         string
		serverName   string
		wantErr      require.ErrorAssertionFunc
		expCallCount int
	}{
		{
			desc:         "matching handler",
			serverName:   "foo",
			wantErr:      require.NoError,
			expCallCount: 1,
		},
		{
			desc:         "matching handler ignoring case",
			serverName:   "FOO",
			wantErr:      require.NoError,
			expCallCount: 1,
		},
		{
			desc:         "no matching handler",
			serverName:   "bar",
			wantErr:      require.Error,
			expCallCount: 0,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			c, s := net.Pipe()

			doneCh := make(chan struct{})
			go func() {
				defer close(doneCh)

				err := tls.Client(c, &tls.Config{ServerName: test.serverName}).Handshake()
				require.Error(t, io.EOF, err)

				err = c.Close()
				require.NoError(t, err)
			}()

			var handlerCallCount int

			mux := Mux{}
			mux.Handle("foo", HandlerFunc(func(conn net.Conn) error {
				handlerCallCount++

				return nil
			}))

			err := mux.ServeConn(s)
			test.wantErr(t, err)

			err = s.Close()
			require.NoError(t, err)

			<-doneCh

			assert.Equal(t, test.expCallCount, handlerCallCount)
		})
	}
}
