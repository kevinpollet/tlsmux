package tlsmux

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConn_Read(t *testing.T) {
	tests := []struct {
		desc          string
		peeked        []byte
		readCallCount int
	}{
		{
			desc:   "should read peeked bytes first",
			peeked: []byte{0, 1, 2, 3, 4, 5, 6},
		},
		{
			desc:          "should read from the net.Conn if peeked bytes are empty",
			readCallCount: 1,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			var readCallCount int
			mock := connMock{read: func(_ []byte) (int, error) {
				readCallCount++

				return 0, nil
			}}

			c := conn{mock, test.peeked}
			buffer := make([]byte, len(test.peeked))

			n, err := c.Read(buffer)
			require.NoError(t, err)

			assert.Equal(t, 0, len(c.peeked))
			assert.Equal(t, len(test.peeked), n)
			assert.Equal(t, len(test.peeked), len(buffer))
			assert.Equal(t, test.readCallCount, readCallCount)
		})
	}
}

type connMock struct {
	net.Conn

	read func([]byte) (int, error)
}

func (c connMock) Read(b []byte) (int, error) {
	return c.read(b)
}
